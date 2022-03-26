package data

import (
	"encoding/json"
	"time"
)

type Timef time.Time

func (t *Timef) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	var err error
	t_, err := time.Parse("2006-01-02T15:04", s)
	*t = Timef(t_)
	return err
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Read struct {
	BookName   string   `json:"bookName"`
	AuthorName string   `json:"authorName"`
	Genres     []string `json:"genres"`
	Thoughts   string   `json:"thoughts"`
	ReadAt     Timef    `json:"readAt"`
}

type DbUser struct {
	UserId   string
	Username string
	Password string
}

type DbAuthor struct {
	AuthorId string
	Name     string
}

type DbGenre struct {
	GenreId string
	Name    string
}

type DbBook struct {
	BookId   string
	AuthorId string
	Name     string
}

type DbBooksGenres struct {
	BookId  string
	GenreId string
}

type DbRead struct {
	ReadId   string
	UserId   string
	BookId   string
	Thoughts string
	ReadAt   Timef
}
