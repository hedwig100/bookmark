package data

import (
	"context"
	"fmt"
	"os"
	"time"

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
		return "", InternalServerError
	}

	_, err = db.pool.Exec(context.Background(),
		"INSERT INTO users (user_id,username,password) VALUES ($1,$2,$3);",
		user_id, user.Username, hashedPassword)

	if err != nil {
		slog.Infof("db insert error: %v", err)
		return "", ErrUserAlreadyRegistered
	}

	return user_id, nil
}

func (db *DbReal) Login(user User) (string, error) {
	rows, err := db.pool.Query(context.Background(), "SELECT user_id,password FROM users WHERE username = $1", user.Username)
	if err != nil {
		slog.Err(err)
		return "", InternalServerError
	}
	if !rows.Next() {
		return "", ErrUserNotFound
	}
	var userId, password string
	if err = rows.Scan(&userId, &password); err != nil {
		slog.Err(err)
		return "", InternalServerError
	}
	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)); err != nil {
		return "", ErrPasswordInvalid
	}
	return userId, nil
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
		readId, userId, bookId, read.Thoughts, time.Time(read.ReadAt))
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

func (db *DbReal) ReadsGet(username string) ([]ReadWithId, error) {
	userId, err := db.getUserId(username)
	if err != nil {
		return []ReadWithId{}, err
	}
	rows, err := db.pool.Query(context.Background(), `SELECT r.read_id,ba.book_id,ba.book_name,ba.author_name,r.thoughts,r.read_at
FROM reads AS r
INNER JOIN book_author AS ba
ON r.book_id = ba.book_id 
WHERE r.user_id = $1`, userId)
	defer rows.Close()

	resp := make([]ReadWithId, 0)
	respBookId := make([]string, 0)
	for rows.Next() {
		var read ReadWithId
		var bookId string
		var tmp_time time.Time
		if err := rows.Scan(&read.ReadId, &bookId, &read.BookName, &read.AuthorName, &read.Thoughts, &tmp_time); err != nil {
			return nil, err
		}
		read.ReadAt = Timef(tmp_time)
		resp = append(resp, read)
		respBookId = append(respBookId, bookId)
	}

	for i := range resp {
		rows, _ = db.pool.Query(context.Background(), `SELECT (
	g.name
) FROM books_genres as bg
INNER JOIN genres as g ON bg.genre_id = g.genre_id
WHERE bg.book_id = $1`, respBookId[i])
		resp[i].Genres = make([]string, 0)
		for rows.Next() {
			var genre string
			if err := rows.Scan(&genre); err != nil {
				return nil, err
			}
			resp[i].Genres = append(resp[i].Genres, genre)
		}
	}

	return resp, nil
}

func (db *DbReal) ReadGet(readId string) (Read, error) {

	var read Read
	var bookId string
	var readAt time.Time
	err := db.pool.QueryRow(context.Background(), `SELECT book_id,thoughts,read_at FROM reads WHERE read_id = $1`, readId).
		Scan(&bookId, &read.Thoughts, &readAt)
	if err != nil {
		slog.Err(err)
		return Read{}, ErrReadNotFound
	}
	read.ReadAt = Timef(readAt)

	err = db.pool.QueryRow(context.Background(), `SELECT book_name,author_name FROM book_author WHERE book_id = $1`, bookId).
		Scan(&read.BookName, &read.AuthorName)
	if err != nil {
		slog.Err(err)
		return Read{}, ErrBookNotFound
	}

	rows, _ := db.pool.Query(context.Background(), `SELECT (
		g.name
	) FROM books_genres as bg
	INNER JOIN genres as g ON bg.genre_id = g.genre_id
	WHERE bg.book_id = $1`, bookId)
	read.Genres = make([]string, 0)
	for rows.Next() {
		var genre string
		if err := rows.Scan(&genre); err != nil {
			slog.Err(err)
			return Read{}, InternalServerError
		}
		read.Genres = append(read.Genres, genre)
	}

	return read, nil
}

func (db *DbReal) getUserId(username string) (string, error) {
	var userId string
	row := db.pool.QueryRow(context.Background(), "SELECT user_id FROM users WHERE username = $1", username)
	if err := row.Scan(&userId); err != nil {
		return "", err
	}
	return userId, nil
}
