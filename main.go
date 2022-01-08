package main

import (
	"github.com/jchen42703/crud/db"
	"github.com/jchen42703/crud/router"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description Conduit API
// @title Conduit API

// @host 127.0.0.1:8585
// @BasePath /api

// @schemes http https
// @produce	application/json
// @consumes application/json

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	r := router.New()

	r.GET("/swagger/*", echoSwagger.WrapHandler)

	// v1 := r.Group("/api")

	connections, err := db.NewConnections()
	if err != nil {
		r.Logger.Fatalf("failed to initialize db or cache: %s", err)
	}
	defer connections.DB.Close()

	// us := store.NewUserStore(d)
	// as := store.NewArticleStore(d)
	// h := handler.NewHandler(us, as)
	// h.Register(v1)
	r.Logger.Fatal(r.Start("127.0.0.1:8585"))
}
