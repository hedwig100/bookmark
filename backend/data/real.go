package data

import (
	"context"
	"fmt"
	"os"

	"github.com/hedwig100/bookmark/backend/slog"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func NewDbReal() *DbReal {
	// connect to database
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	// url := "postgres://username:password@hostname:port/database_name"
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, username)

	var err error
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		slog.Fatalf("Unable to connect to database: %v", err)
	}
	slog.Info("Db connection successful!")
	return &DbReal{
		pool: pool,
	}
}

func (db *DbReal) UserCreate(user User) (string, error) {
	// generate uuid
	user_id, err := id()
	if err != nil {
		slog.Infof("uuid generation failure: %v", err)
		return "", err
	}

	// hash password
	hashedPassword := []byte(user.Password)
	hashedPassword, err = bcrypt.GenerateFromPassword(hashedPassword, 10)
	if err != nil {
		slog.Infof("internal error: %v", err)
		return "", err
	}

	_, err = db.pool.Exec(context.Background(),
		"INSERT INTO users (user_id,username,password) VALUES ($1,$2,$3);",
		user_id, user.Username, hashedPassword)

	if err != nil {
		slog.Infof("db insert error: %v", err)
		return "", err
	}

	return user_id, nil
}
