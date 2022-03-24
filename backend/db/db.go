package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Pool *pgxpool.Pool

func init() {

	// connect to database
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	// url := "postgres://username:password@hostname:port/database_name"
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, username)

	var err error
	Pool, err = pgxpool.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		os.Exit(1)
	}
	log.Println("Db connection successful!")
}
