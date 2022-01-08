package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/jchen42703/crud/internal/templates"
	"github.com/labstack/echo/v4"
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

		return c.JSON(http.StatusCreated, nil)
	}
}
