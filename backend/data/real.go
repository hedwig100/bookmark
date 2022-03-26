package data

import (
	"context"
	"fmt"
	"os"

	"github.com/hedwig100/bookmark/backend/slog"
	"github.com/jackc/pgx/v5/pgxpool"
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
	user_id := id()

	// hash password
	hashedPassword := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(hashedPassword, 10)
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

func (db *DbReal) ReadCreate(username string, read Read) error {
	tx, err := db.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// insert author
	authorId := id()
	_, err = tx.Exec(context.Background(),
		"INSERT INTO authors (author_id,name) SELECT $1,$2 WHERE NOT EXISTS (SELECT * FROM authors WHERE authors.name = $3)",
		authorId, read.AuthorName, read.AuthorName)
	if err != nil {
		return err
	}
	row := tx.QueryRow(context.Background(), "SELECT author_id FROM authors WHERE name = $1", read.AuthorName)
	if err = row.Scan(&authorId); err != nil {
		return err
	}
	slog.Info(authorId)

	// insert book
	bookId := id()
	_, err = tx.Exec(context.Background(),
		"INSERT INTO books (book_id,author_id,name) SELECT $1,$2,$3 WHERE NOT EXISTS (SELECT * FROM books WHERE name = $4)",
		bookId, authorId, read.BookName, read.BookName)
	row = tx.QueryRow(context.Background(), "SELECT book_id FROM books WHERE name = $1", read.BookName)
	if err = row.Scan(&bookId); err != nil {
		return err
	}
	slog.Info(bookId)

	// insert genre
	var genreId string
	for _, genre := range read.Genres {
		genreId = id()
		_, err = tx.Exec(context.Background(),
			"INSERT INTO genres (genre_id,name) SELECT $1,$2 WHERE NOT EXISTS (SELECT * FROM genres WHERE name = $3)",
			genreId, genre, genre)
		if err != nil {
			return err
		}
		row = tx.QueryRow(context.Background(), "SELECT genre_id FROM genres WHERE name = $1", genre)
		if err = row.Scan(&genreId); err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(),
			"INSERT INTO books_genres (book_id,genre_id) SELECT $1,$2 WHERE NOT EXISTS (SELECT * FROM books_genres WHERE book_id = $3 AND genre_id = $4)",
			bookId, genreId, bookId, genreId)
		if err != nil {
			return err
		}
	}

	var userId string
	row = tx.QueryRow(context.Background(), "SELECT user_id FROM users WHERE username = $1", username)
	if err := row.Scan(&userId); err != nil {
		return err
	}
	// insert read
	readId := id()
	_, err = tx.Exec(context.Background(),
		"INSERT INTO reads (read_id,user_id,book_id,thoughts,read_at) VALUES ($1,$2,$3,$4,$5)",
		readId, userId, bookId, read.Thoughts, read.ReadAt)
	if err != nil {
		return err
	}

	// commit
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}
