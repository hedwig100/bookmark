package data

import (
	"github.com/hedwig100/bookmark/backend/slog"
	"golang.org/x/crypto/bcrypt"
)

func NewDbMock() *DbMock {
	return &DbMock{
		users: make([]DbUser, 0),
	}
}

func (db *DbMock) UserCreate(user User) (string, error) {
	user_id, err := id()
	if err != nil {
		return "", err
	}

	// hash password
	hashedPassword := []byte(user.Password)
	hashedPassword, err = bcrypt.GenerateFromPassword(hashedPassword, 10)
	if err != nil {
		slog.Infof("internal error: %v", err)
		return "", err
	}

	db.users = append(db.users, DbUser{
		UserId:   user_id,
		Username: user.Username,
		Password: string(hashedPassword),
	})

	return user_id, nil
}
