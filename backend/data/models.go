package data

import "time"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Read struct {
	BookName   string    `json:"bookName"`
	AuthorName string    `json:"authorName"`
	Genres     []string  `json:"genres"`
	Thoughts   string    `json:"thoughts"`
	ReadAt     time.Time `json:"readAt"`
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
	ReadAt   time.Time
}
