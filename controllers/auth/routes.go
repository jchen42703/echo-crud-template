package auth

import (
	"github.com/jchen42703/crud/db"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, conn *db.Connections) {
	authGroup := g.Group("/auth")
	authGroup.POST("/signup", SignUp(conn.DB))
	authGroup.POST("/login", Login(conn.DB, conn.Cache))
	authGroup.POST("/logout", Logout(conn.Cache))

	// jwtMiddleware := middleware.JWT(utils.JWTSecret)
	// guestUsers := v1.Group("/users")
	// guestUsers.POST("", h.SignUp)
	// guestUsers.POST("/login", h.Login)

	// user := v1.Group("/user", jwtMiddleware)
	// user.GET("", h.CurrentUser)
	// user.PUT("", h.UpdateUser)

	// profiles := v1.Group("/profiles", jwtMiddleware)
	// profiles.GET("/:username", h.GetProfile)
	// profiles.POST("/:username/follow", h.Follow)
	// profiles.DELETE("/:username/follow", h.Unfollow)

	// articles := v1.Group("/articles", middleware.JWTWithConfig(
	// 	middleware.JWTConfig{
	// 		Skipper: func(c echo.Context) bool {
	// 			if c.Request().Method == "GET" && c.Path() != "/api/articles/feed" {
	// 				return true
	// 			}
	// 			return false
	// 		},
	// 		SigningKey: utils.JWTSecret,
	// 	},
	// ))
	// articles.POST("", h.CreateArticle)
	// articles.GET("/feed", h.Feed)
	// articles.PUT("/:slug", h.UpdateArticle)
	// articles.DELETE("/:slug", h.DeleteArticle)
	// articles.POST("/:slug/comments", h.AddComment)
	// articles.DELETE("/:slug/comments/:id", h.DeleteComment)
	// articles.POST("/:slug/favorite", h.Favorite)
	// articles.DELETE("/:slug/favorite", h.Unfavorite)
	// articles.GET("", h.Articles)
	// articles.GET("/:slug", h.GetArticle)
	// articles.GET("/:slug/comments", h.GetComments)

	// tags := v1.Group("/tags")
	// tags.GET("", h.Tags)
}
