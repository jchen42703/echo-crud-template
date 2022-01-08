package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq" // <------------ here
)

// The "db" package level variable will hold the reference to our database instance
func initDB(dbInfo string) (*sql.DB, error) {
	var err error
	// Connect to the postgres db
	//you might have to change the connection string to add your database credentials
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}

func initCache() (redis.Conn, error) {
	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		return nil, err
	}
	// Assign the connection to the package level `cache` variable
	return conn, nil
}

type Connections struct {
	DB    *sql.DB
	Cache redis.Conn
}

// Initializes the connections to the cache and the db
func NewConnections() (*Connections, error) {
	// initialize our database connection
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	sqlDB, err := initDB(dbinfo)
	if err != nil {
		return nil, err
	}

	// Initialize redis cache
	cache, err := initCache()
	if err != nil {
		return nil, err
	}

	return &Connections{
		DB:    sqlDB,
		Cache: cache,
	}, nil
}
