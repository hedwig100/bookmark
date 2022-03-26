package data

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Db is a interface for representing database connection.
type Db interface {
	// UserCreate receives user and return user_id and error (if any).
	UserCreate(user User) (string, error)
}

// DbReal connects to a real databaes. The pointer of this object implements Db interface.
type DbReal struct {
	pool *pgxpool.Pool
}

// DbMock doesn't connect ot a real database, but simulate a database. The pointer of this object implements Db interface.
type DbMock struct {
	users []DbUser
}

// id create uuid.
func id() (string, error) {
	ret, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return ret.String(), nil
}
