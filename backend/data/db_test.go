package data_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hedwig100/bookmark/backend/data"
)

var db data.Db

func TestMain(m *testing.M) {
	db = data.NewDbReal()
	m.Run()
}

func TestUserCreate(t *testing.T) {
	uc := []struct {
		user        data.User
		expectError bool
	}{
		{
			user: data.User{
				Username: "hedwig100",
				Password: "abcde12345",
			},
			expectError: false,
		},
		{
			user: data.User{
				Username: "python39",
				Password: "1234567890",
			},
			expectError: false,
		},
		{
			user: data.User{
				Username: "hedwig100",
				Password: "1234567890",
			},
			expectError: true,
		},
	}
	for i, td := range uc {
		t.Run(fmt.Sprintf("UserCreate%d", i), func(t *testing.T) {
			_, err := db.UserCreate(td.user)
			if td.expectError && err == nil {
				t.Fatal("err expected")
			}
			if !td.expectError && err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestReadCreate(t *testing.T) {
	rc := []struct {
		username    string
		read        data.Read
		expectError bool
	}{
		{
			username: "hedwig100",
			read: data.Read{
				BookName:   "ABC",
				AuthorName: "J.K.Rowling",
				Genres: []string{
					"fantasy",
					"textbook",
				},
				Thoughts: "fantastic!",
				ReadAt:   data.Timef(time.Now()),
			},
			expectError: false,
		},
		{
			username: "hedwig100",
			read: data.Read{
				BookName:   "The Little Prince",
				AuthorName: "Antoine Marie Jean-Baptiste Roger, comte de Saint-Exupery",
				Genres: []string{
					"fantasy",
					"for children",
				},
				Thoughts: "It makes me think seriously.",
				ReadAt:   data.Timef(time.Now()),
			},
			expectError: false,
		},
		{
			username: "python39",
			read: data.Read{
				BookName:   "The Little Prince",
				AuthorName: "Antoine Marie Jean-Baptiste Roger, comte de Saint-Exupery",
				Genres: []string{
					"fantasy",
					"thought-provoking",
				},
				Thoughts: "",
				ReadAt:   data.Timef(time.Now()),
			},
			expectError: false,
		},
	}
	for i, td := range rc {
		t.Run(fmt.Sprintf("ReadCreate%d", i), func(t *testing.T) {
			err := db.ReadCreate(td.username, td.read)
			if td.expectError && err == nil {
				t.Fatal("err expected")
			}
			if !td.expectError && err != nil {
				t.Fatal(err)
			}
		})
	}
}
