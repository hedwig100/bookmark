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
		return "", InternalServerError
	}

	_, err = db.selectUsers(user.Username)
	if err == nil {
		return "", UserAlreadyRegistered
	}

	db.users = append(db.users, DbUser{
		UserId:   user_id,
		Username: user.Username,
		Password: string(hashedPassword),
	})

	return user_id, nil
}

func (db *DbMock) Login(user User) (string, error) {
	for _, registeredUser := range db.users {
		if registeredUser.Username == user.Username {
			if err := bcrypt.CompareHashAndPassword([]byte(registeredUser.Password), []byte(user.Password)); err != nil {
				return "", err
			}
			return registeredUser.UserId, nil
		}
	}
	return "", fmt.Errorf("user(%s) not found", user.Username)
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

func (db *DbMock) ReadGet(username string) ([]Read, error) {
	user, err := db.selectUsers(username)
	if err != nil {
		return nil, err
	}
	userId := user.UserId
	reads := db.selectRead(userId)
	ret := make([]Read, 0)
	for _, read := range reads {
		bookName, authorName, genres, err := db.selectBook(read.BookId)
		if err != nil {
			return nil, err
		}
		ret = append(ret, Read{
			BookName:   bookName,
			AuthorName: authorName,
			Genres:     genres,
			Thoughts:   read.Thoughts,
			ReadAt:     read.ReadAt,
		})
	}
	return ret, nil
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

func (db *DbMock) selectRead(userId string) []DbRead {
	ret := make([]DbRead, 0)
	for _, read := range db.reads {
		if read.UserId == userId {
			ret = append(ret, read)
		}
	}
	return ret
}

func (db *DbMock) selectBook(bookId string) (bookName string, authorName string, genres []string, err error) {

	var book DbBook
	var find bool
	for _, bk := range db.books {
		if bk.BookId == bookId {
			book = bk
			find = true
			break
		}
	}
	if !find {
		return "", "", nil, fmt.Errorf("the book isn't found.")
	}
	bookName = book.Name

	find = false
	for _, at := range db.authors {
		if at.AuthorId == book.AuthorId {
			authorName = at.Name
			find = true
			break
		}
	}
	if !find {
		return "", "", nil, fmt.Errorf("the author isn't found.")
	}

	genres = make([]string, 0)
	for _, bookGenre := range db.booksGenres {
		if bookGenre.BookId == bookId {
			for _, genre := range db.genres {
				if genre.GenreId == bookGenre.GenreId {
					genres = append(genres, genre.Name)
				}
			}
		}
	}

	return bookName, authorName, genres, nil
}
