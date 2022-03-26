package data

import (
	"fmt"

	"github.com/hedwig100/bookmark/backend/slog"
	"golang.org/x/crypto/bcrypt"
)

func NewDbMock() *DbMock {
	return &DbMock{
		users:       make([]DbUser, 0),
		authors:     make([]DbAuthor, 0),
		genres:      make([]DbGenre, 0),
		books:       make([]DbBook, 0),
		booksGenres: make([]DbBooksGenres, 0),
		reads:       make([]DbRead, 0),
	}
}

func (db *DbMock) UserCreate(user User) (string, error) {
	user_id := id()

	// hash password
	hashedPassword := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(hashedPassword, 10)
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

func (db *DbMock) ReadCreate(username string, read Read) error {
	authorId := db.insertAuthor(read.AuthorName)
	bookId := db.insertBook(authorId, read.BookName)
	var genreId string
	for _, genre := range read.Genres {
		genreId = db.insertGenre(genre)
		db.insertBooksGenres(bookId, genreId)
	}

	user, err := db.selectUsers(username)
	if err != nil {
		return err
	}
	db.reads = append(db.reads, DbRead{
		ReadId:   id(),
		UserId:   user.UserId,
		BookId:   bookId,
		Thoughts: read.Thoughts,
		ReadAt:   read.ReadAt,
	})
	return nil
}

func (db *DbMock) insertAuthor(name string) string {
	var already bool
	for _, author := range db.authors {
		if author.Name == name {
			already = true
			return author.AuthorId
		}
	}
	var authorId string
	if !already {
		authorId = id()
		db.authors = append(db.authors, DbAuthor{
			AuthorId: authorId,
			Name:     name,
		})
	}
	return authorId
}

func (db *DbMock) insertGenre(name string) string {
	var already bool
	for _, genre := range db.genres {
		if genre.Name == name {
			already = true
			return genre.GenreId
		}
	}
	var genreId string
	if !already {
		genreId = id()
		db.genres = append(db.genres, DbGenre{
			GenreId: genreId,
			Name:    name,
		})
	}
	return genreId
}

func (db *DbMock) insertBook(authorId string, name string) string {
	var already bool
	for _, book := range db.books {
		if book.Name == name {
			already = true
			return book.BookId
		}
	}
	var bookId string
	if !already {
		bookId = id()
		db.books = append(db.books, DbBook{
			BookId:   bookId,
			AuthorId: authorId,
			Name:     name,
		})
	}
	return bookId
}

func (db *DbMock) insertBooksGenres(bookId string, genreId string) {
	var already bool
	for _, bookGenre := range db.booksGenres {
		if bookGenre.BookId == bookId && bookGenre.GenreId == genreId {
			already = true
			return
		}
	}
	if !already {
		db.booksGenres = append(db.booksGenres, DbBooksGenres{
			BookId:  bookId,
			GenreId: genreId,
		})
	}
	return
}

func (db *DbMock) selectUsers(username string) (DbUser, error) {
	for _, user := range db.users {
		if user.Username == username {
			return user, nil
		}
	}
	return DbUser{}, fmt.Errorf("user '%s' not found", username)
}
