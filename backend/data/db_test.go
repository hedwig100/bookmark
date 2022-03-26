package data_test

import (
	"testing"

	"github.com/hedwig100/bookmark/backend/data"
)

func TestDbReal(t *testing.T) {
	db := data.NewDbReal()

	user := data.User{
		Username: "hedwig100",
		Password: "abcde12345",
	}
	_, err := db.UserCreate(user)
	if err != nil {
		t.Fatalf("failure on UserCreate: %v", err)
	}
}
