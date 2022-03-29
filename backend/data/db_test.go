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

func TestDb(t *testing.T) {
	t.Run("UserCreate", testUserCreate)
	t.Run("Login", testLogin)
	t.Run("ReadCreate", testReadCreate)
	t.Run("ReadGet", testReadGet)
}

func testUserCreate(t *testing.T) {
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

func testLogin(t *testing.T) {
	l := []struct {
		user        data.User
		expectError bool
	}{
		{user: data.User{Username: "hedwig100", Password: "abcde12345"}, expectError: false},
		{user: data.User{Username: "python39", Password: "1234567890"}, expectError: false},
		{user: data.User{Username: "hedwig100", Password: "9f83o"}, expectError: true},
		{user: data.User{Username: "fio", Password: "u3"}, expectError: true},
	}
	for i, td := range l {
		t.Run(fmt.Sprintf("Login%d", i), func(t *testing.T) {
			_, err := db.Login(td.user)
			if td.expectError && err == nil {
				t.Fatal("err expected")
			}
			if !td.expectError && err != nil {
				t.Fatal(err)
			}
		})
	}
}

func testReadCreate(t *testing.T) {
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

func compare[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func testReadGet(t *testing.T) {
	rg := []struct {
		username    string
		reads       []data.Read
		expectError bool
	}{
		{
			username: "hedwig100",
			reads: []data.Read{
				{
					BookName:   "ABC",
					AuthorName: "J.K.Rowling",
					Genres: []string{
						"fantasy",
						"textbook",
					},
					Thoughts: "fantastic!",
					ReadAt:   data.Timef(time.Now()),
				},
				{
					BookName:   "The Little Prince",
					AuthorName: "Antoine Marie Jean-Baptiste Roger, comte de Saint-Exupery",
					Genres: []string{
						"fantasy",
						"for children",
					},
					Thoughts: "It makes me think seriously.",
					ReadAt:   data.Timef(time.Now()),
				},
			},
			expectError: false,
		},
		{
			username: "python39",
			reads: []data.Read{
				{
					BookName:   "The Little Prince",
					AuthorName: "Antoine Marie Jean-Baptiste Roger, comte de Saint-Exupery",
					Genres: []string{
						"fantasy",
						"thought-provoking",
					},
					Thoughts: "",
					ReadAt:   data.Timef(time.Now()),
				},
			},
			expectError: false,
		},
	}
	for i, td := range rg {
		t.Run(fmt.Sprintf("ReadGet%d", i), func(t *testing.T) {
			reads, err := db.ReadGet(td.username)
			if td.expectError && err == nil {
				t.Fatal("err expected")
			}
			if !td.expectError && err != nil {
				t.Fatal(err)
			}
			_ = reads
		})
	}
}
