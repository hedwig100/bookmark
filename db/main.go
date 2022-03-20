package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStr(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)

	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func testData(N int, filename string) {

	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer file.Close()

	users := N
	authors := N / 10
	books := N / 10
	genres := N / 100
	books_genres := N / 10
	reads := N * 10
	follows := N * 10
	messages := N * 10
	open_messages := N * 10

	// users
	for i := 0; i < users; i++ {
		fmt.Fprintln(file, fmt.Sprintf("INSERT INTO users (user_id,username,password) VALUES (%d,'%s','%s');", i, randStr(10), randStr(10)))
	}

	// authors
	for i := 0; i < authors; i++ {
		fmt.Fprintln(file, fmt.Sprintf("INSERT INTO authors (author_id,name) VALUES (%d,'%s');", i, randStr(10)))
	}

	// books
	for i := 0; i < books; i++ {
		fmt.Fprintln(file, fmt.Sprintf("INSERT INTO books (book_id,author_id,name) VALUES (%d,%d,'%s');", i, rand.Intn(authors), randStr(10)))
	}

	// genres
	for i := 0; i < genres; i++ {
		fmt.Fprintln(file, fmt.Sprintf("INSERT INTO genres (genre_id,name) VALUES (%d,'%s');", i, randStr(10)))
	}

	// books_genres
	for i := 0; i < books_genres; i++ {
		fmt.Fprintln(file, fmt.Sprintf("INSERT INTO books_genres (book_id,genre_id) VALUES (%d,%d);", rand.Intn(books), rand.Intn(genres)))
	}

	// reads
	for i := 0; i < reads; i++ {
		fmt.Fprintln(file, fmt.Sprintf("INSERT INTO reads (read_id,user_id,book_id,thoughts,read_at) VALUES (%d,%d,%d,'%s',CURRENT_DATE);", i, rand.Intn(users), rand.Intn(books), randStr(30)))
	}

	// follows
	for i := 0; i < follows; i++ {
		fmt.Fprintln(file, fmt.Sprintf("INSERT INTO follows (follower_id,followee_id) VALUES (%d,%d);", rand.Intn(users), rand.Intn(users)))
	}

	// messages
	for i := 0; i < messages; i++ {
		fmt.Fprintln(file, fmt.Sprintf("INSERT INTO messages (message_id,sender_id,receiver_id,content,send_at) VALUES (%d,%d,%d,'%s',CURRENT_TIMESTAMP);", i, rand.Intn(users), rand.Intn(users), randStr(20)))
	}

	// open_messages
	for i := 0; i < open_messages; i++ {
		fmt.Fprintln(file, fmt.Sprintf("INSERT INTO open_messages (message_id,sender_id,content,send_at) VALUES (%d,%d,'%s',CURRENT_TIMESTAMP);", i, rand.Intn(users), randStr(20)))
	}
}

func testQuery(N int, filename string) {

	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer file.Close()

	users := N
	books := N / 10

	// users
	for i := 0; i < N; i++ {
		fmt.Fprintln(file, fmt.Sprintf("SELECT username,password FROM users where user_id = %d;", rand.Intn(users)))
	}

	// books
	for i := 0; i < N; i++ {
		fmt.Fprintln(file, fmt.Sprintf("SELECT book_name,author_name FROM book_author WHERE book_id = %d;", rand.Intn(books)))
	}

	// genre
	for i := 0; i < N; i++ {
		fmt.Fprintln(file, fmt.Sprintf(`SELECT g.name FROM books AS b
INNER JOIN books_genres as bg ON b.book_id = bg.book_id
INNER JOIN genres as g ON bg.genre_id = g.genre_id
WHERE b.book_id = %d;`, rand.Intn(books)))
	}

	// read
	for i := 0; i < N; i++ {
		fmt.Fprintln(file, fmt.Sprintf(`SELECT ba.book_name,ba.author_name,r.thoughts,r.read_at
FROM book_author AS ba
INNER JOIN reads AS r ON r.book_id = ba.book_id
WHERE r.user_id = %d
ORDER BY r.read_at;`, rand.Intn(users)))
	}

	// follows
	for i := 0; i < N; i++ {
		fmt.Fprintln(file, fmt.Sprintf(`SELECT followee_id FROM follows WHERE follower_id = %d;`, rand.Intn(users)))
	}

	// messsages
	for i := 0; i < N; i++ {
		fmt.Fprintln(file, fmt.Sprintf("SELECT sender_id,receiver_id,content,send_at FROM messages WHERE sender_id = %d AND receiver_id = %d;", rand.Intn(users), rand.Intn(users)))
	}

	// open_messages
	for i := 0; i < N; i++ {
		fmt.Fprintln(file, fmt.Sprintf("SELECT sender_id,content,send_at FROM open_messages WHERE sender_id IN (SELECT followee_id FROM follows WHERE follower_id = %d);", rand.Intn(users)))
	}
}

func sample() {
	file, err := os.Create("db/init/sample.sql")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer file.Close()

	fmt.Fprintln(file, `
-- users
INSERT INTO users (user_id,username,password) VALUES (0,'Alice','password1'); -- actually password is encrypted
INSERT INTO users (user_id,username,password) VALUES (1,'Bob','password2'); -- actually password is encrypted
INSERT INTO users (user_id,username,password) VALUES (2,'Cate','password3'); -- actually password is encrypted

-- authors
INSERT INTO authors (author_id,name) VALUES (0,'Agatha Christie');
INSERT INTO authors (author_id,name) VALUES (1,'J.K. Rowling');
INSERT INTO authors (author_id,name) VALUES (2,'Minato Kanae');

-- books
INSERT INTO books (book_id,author_id,name) VALUES (0,0,'And Then There Were None');
INSERT INTO books (book_id,author_id,name) VALUES (1,1,'Harry Potter and the Philosopher''s Stone');
INSERT INTO books (book_id,author_id,name) VALUES (2,1,'Harry Potter and the Chamber of Secrets');
INSERT INTO books (book_id,author_id,name) VALUES (3,2,'Kokuhaku');

-- genres
INSERT INTO genres (genre_id,name) VALUES (0,'mystery');
INSERT INTO genres (genre_id,name) VALUES (1,'fantasy');
INSERT INTO genres (genre_id,name) VALUES (2,'horror');
INSERT INTO genres (genre_id,name) VALUES (3,'for children');

-- books_genres
INSERT INTO books_genres (book_id,genre_id) VALUES (0,0);
INSERT INTO books_genres (book_id,genre_id) VALUES (1,1);
INSERT INTO books_genres (book_id,genre_id) VALUES (1,3);
INSERT INTO books_genres (book_id,genre_id) VALUES (2,1);
INSERT INTO books_genres (book_id,genre_id) VALUES (2,3);
INSERT INTO books_genres (book_id,genre_id) VALUES (3,0);
INSERT INTO books_genres (book_id,genre_id) VALUES (3,2);

-- reads
INSERT INTO reads (read_id,user_id,book_id,thoughts,read_at) VALUES (0,0,0,'Very suprised with the last.',CURRENT_DATE);
INSERT INTO reads (read_id,user_id,book_id,thoughts,read_at) VALUES (1,0,1,'Harry cheered me up.',CURRENT_DATE);
INSERT INTO reads (read_id,user_id,book_id,thoughts,read_at) VALUES (2,1,2,'Voldemort scared me.',CURRENT_DATE);
INSERT INTO reads (read_id,user_id,book_id,thoughts,read_at) VALUES (3,1,3,'I felt the madness of the teacher.',CURRENT_DATE);

-- follows
INSERT INTO follows (follower_id,followee_id) VALUES (0,1);
INSERT INTO follows (follower_id,followee_id) VALUES (1,0);
INSERT INTO follows (follower_id,followee_id) VALUES (0,2);
INSERT INTO follows (follower_id,followee_id) VALUES (2,1);

-- messages
INSERT INTO messages (message_id,sender_id,receiver_id,content,send_at) VALUES (0,0,1,'What are you doing now?',CURRENT_TIMESTAMP);
INSERT INTO messages (message_id,sender_id,receiver_id,content,send_at) VALUES (1,1,0,'I''m reading a book.',CURRENT_TIMESTAMP);
INSERT INTO messages (message_id,sender_id,receiver_id,content,send_at) VALUES (2,0,1,'Which book are you reading?',CURRENT_TIMESTAMP);
INSERT INTO messages (message_id,sender_id,receiver_id,content,send_at) VALUES (3,1,0,'I read Harry Potter series. It''s interesting!',CURRENT_TIMESTAMP);

-- open_messages
INSERT INTO open_messages (message_id,sender_id,content,send_at) VALUES (0,0,'I have read Kokuhaku.',CURRENT_TIMESTAMP);
INSERT INTO open_messages (message_id,sender_id,content,send_at) VALUES (1,1,'This book scares me,',CURRENT_TIMESTAMP);
INSERT INTO open_messages (message_id,sender_id,content,send_at) VALUES (2,2,'I have no books I want to read.',CURRENT_TIMESTAMP);`)

	file, err = os.Create("db/query.sql")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer file.Close()

	fmt.Fprintln(file, `
-- users
SELECT username,password FROM users where user_id = 0;

-- books
SELECT book_name,author_name FROM book_author WHERE book_name = 'And Then There Were None';

-- genres
SELECT (
	g.name
) FROM books AS b 
INNER JOIN books_genres as bg ON b.book_id = bg.book_id
INNER JOIN genres as g ON bg.genre_id = g.genre_id
WHERE b.name = 'Harry Potter and the Philosopher''s Stone';

-- reads
SELECT ba.book_name,ba.author_name,r.thoughts,r.read_at
FROM book_author AS ba
INNER JOIN reads AS r ON r.book_id = ba.book_id
WHERE r.user_id = 0
ORDER BY r.read_at;

-- follows
SELECT followee_id FROM follows WHERE follower_id = 0;

-- messages
SELECT sender_id,receiver_id,content,send_at FROM messages 
WHERE sender_id = 0 AND receiver_id = 1;

-- open messages
SELECT sender_id,content,send_at FROM open_messages 
WHERE sender_id IN (SELECT followee_id FROM follows WHERE follower_id = 0);`)
}

var (
	is_sample bool
	is_test   bool
	N         int
)

func init() {
	flag.BoolVar(&is_sample, "sample", false, "create sample data")
	flag.BoolVar(&is_sample, "s", false, "create sample data (shorthand)")
	flag.BoolVar(&is_test, "test", false, "create test data for benchmark")
	flag.BoolVar(&is_test, "t", false, "create test data for benchmark (shorthand)")
	flag.IntVar(&N, "number", 10000, "number of test data")
	flag.IntVar(&N, "n", 10000, "number of test data (shorthand)")
}

func main() {

	// flag parse
	flag.Parse()
	if is_sample && is_test {
		fmt.Fprintln(os.Stderr, "err: both sample and test flag are set.")
		os.Exit(1)
	}

	// make data
	if is_sample {
		sample()
	}

	if is_test {
		testData(N, "db/init/sample.sql")
		testQuery(N, "db/query.sql")
	}
}
