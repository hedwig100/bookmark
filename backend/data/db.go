package data

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hedwig100/bookmark/backend/slog"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserAlreadyRegistered = fmt.Errorf("the username is already registered.")
	InternalServerError      = fmt.Errorf("internal server error")
	ErrUserNotFound          = fmt.Errorf("the username is not found.")
	ErrPasswordInvalid       = fmt.Errorf("the password is invalid.")
)

// Db is a interface for representing database connection.
type Db interface {
	// UserCreate receives user and return user_id and error (if any).
	UserCreate(User) (string, error)

	// Login authenticate User using password
	Login(User) (string, error)

	// ReadCreate receives username,read and return error (if any).
	ReadCreate(string, Read) error

	// ReadGet receives username and return user's reading log
	ReadGet(string) ([]Read, error)
}

// DbReal connects to a real databaes. The pointer of this object implements Db interface.
type DbReal struct {
	pool *pgxpool.Pool
}

// DbMock doesn't connect ot a real database, but simulate a database. The pointer of this object implements Db interface.
type DbMock struct {
	users       []DbUser
	authors     []DbAuthor
	genres      []DbGenre
	books       []DbBook
	booksGenres []DbBooksGenres
	reads       []DbRead
}

// id create uuid.
func id() string {
	ret, err := uuid.NewRandom()
	if err != nil {
		slog.Fatalf("failed to generate uuid: %v", err)
	}
	return ret.String()
}
