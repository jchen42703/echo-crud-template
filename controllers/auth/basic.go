package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jchen42703/crud/internal/templates"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"username"`
}

func SignUp(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse and decode the request body into a new `Credentials` instance
		creds := &Credentials{}
		if err := c.Validate(creds); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, templates.NewError(err))
		}

		err := json.NewDecoder(c.Request().Body).Decode(creds)
		if err != nil {
			// If there is something wrong with the request body, return a 400 status
			// .WriteHeader(http.StatusBadRequest)
			return c.JSON(http.StatusBadRequest, templates.NewError(err))
		}

		// Salt and hash the password using the bcrypt algorithm
		// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, templates.NewError(err))
		}

		// Next, insert the username, along with the hashed password into the database
		if _, err = db.Query("insert into users values ($1, $2)", creds.Username, string(hashedPassword)); err != nil {
			// If there is any issue with inserting into the database, return a 500 error
			return c.JSON(http.StatusInternalServerError, templates.NewError(err))
		}

		return c.JSON(http.StatusCreated, "registered successfully")
	}
}

type LoginReq struct {
	Credentials

	RememberMe bool `json:"rememberMe"`
}

func Login(db *sql.DB, cache redis.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse and decode the request body into a new `Credentials` instance
		req := &LoginReq{}
		if err := c.Validate(req); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, templates.NewError(err))
		}

		err := json.NewDecoder(c.Request().Body).Decode(req)
		if err != nil {
			// If there is something wrong with the request body, return a 400 status
			return c.JSON(http.StatusBadRequest, templates.NewError(err))
		}

		// Get the existing entry present in the database for the given username
		result := db.QueryRow("select password from users where username=$1", req.Username)
		if err != nil {
			// If there is an issue with the database, return a 500 error
			return c.JSON(http.StatusInternalServerError, templates.NewError(err))
		}

		// We create another instance of `Credentials` to store the credentials we get from the database
		storedCreds := &Credentials{}
		// Store the obtained password in `storedCreds`
		err = result.Scan(&storedCreds.Password)
		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusUnauthorized, templates.NewError(err))
			}
			// If the error is of any other type, send a 500 status
			return c.JSON(http.StatusInternalServerError, templates.NewError(err))
		}

		// Compare the stored hashed password, with the hashed version of the password that was received
		if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(req.Password)); err != nil {
			// If the two passwords don't match, return a 401 status
			return c.JSON(http.StatusUnauthorized, templates.NewError(err))
		}

		// If we reach this point, that means the users password was correct, and that they are authorized
		// The default 200 status is sent
		// Create a new random session token
		sessionToken := uuid.NewV4().String()
		// Set the token in the cache, along with the user whom it represents
		// The token has an expiry time of 1 day
		// NOT THE SAFEST because if rememberMe == false and cookie deleted when browser session closes,
		// the session cookie is still in the cache, so technically the same session cookie is still usable.
		// Alternative: When making the authentication middlware, change the cache TTL to something low and refresh
		// the cache key when needed
		_, err = cache.Do("SETEX", sessionToken, templates.DayInSeconds(), req.Username)
		if err != nil {
			// If there is an error in setting the cache, return an internal server error
			return c.JSON(http.StatusInternalServerError, templates.NewError(err))
		}

		// Finally, we set the client cookie for "session_token" as the session token we just generated
		// we also set an expiry time of 120 seconds, the same as the cache
		sessionCookie := &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  templates.DayFromNow(),
			HttpOnly: true,
		}
		c.SetCookie(sessionCookie)

		// Set username cookie
		var usernameCookie *http.Cookie
		if req.RememberMe {
			// Saves the username in a cookie
			usernameCookie = &http.Cookie{
				Name:     "rememberMe",
				Value:    req.Username,
				Expires:  templates.MonthFromNow(),
				HttpOnly: true,
			}
		} else {
			// Expires the remember cookie if you don't want to be remembered
			usernameCookie = &http.Cookie{
				Name:     "rememberMe",
				Value:    "",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			}
		}
		c.SetCookie(usernameCookie)

		return c.String(http.StatusOK, "logged in successfully")
	}
}

type LogoutCredentials struct {
	SessionToken string `json:"sessionToken"`
}

func Logout(cache redis.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse and decode the request body into a new `Credentials` instance
		req := &LogoutCredentials{}
		if err := c.Validate(req); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, templates.NewError(err))
		}

		err := json.NewDecoder(c.Request().Body).Decode(req)
		if err != nil {
			// If there is something wrong with the request body, return a 400 status
			return c.JSON(http.StatusBadRequest, templates.NewError(err))
		}

		// deletes the session token
		_, err = cache.Do("DEL", req.SessionToken)
		if err != nil {
			// If there is an error in setting the cache, return an internal server error
			return c.JSON(http.StatusInternalServerError, templates.NewError(err))
		}

		// overwrites previous cookie
		cookie := &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),

			HttpOnly: true,
		}

		c.SetCookie(cookie)

		return c.String(http.StatusOK, "logged out successfully")
	}
}
