
CREATE TABLE users (
    user_id VARCHAR(128) PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL
);

CREATE TABLE authors (
    author_id VARCHAR(128) PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE
);

CREATE TABLE genres (
    genre_id VARCHAR(128) PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE
);

CREATE TABLE books (
    book_id VARCHAR(128) PRIMARY KEY,
    author_id VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL UNIQUE,
    FOREIGN KEY (author_id) REFERENCES authors(author_id) ON DELETE CASCADE
);

CREATE TABLE books_genres (
    book_id VARCHAR(128) NOT NULL,
    genre_id VARCHAR(128) NOT NULL,
    PRIMARY KEY(book_id,genre_id),
    FOREIGN KEY (book_id) REFERENCES books(book_id) ON DELETE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES genres(genre_id) ON DELETE CASCADE
);

CREATE TABLE reads (
    read_id VARCHAR(128) PRIMARY KEY,
    user_id VARCHAR(128) NOT NULL,
    book_id VARCHAR(128) NOT NULL,
    thoughts VARCHAR(1024),
    read_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(book_id) ON DELETE CASCADE
);

CREATE TABLE follows (
    follower_id VARCHAR(128) NOT NULL,
    followee_id VARCHAR(128) NOT NULL,
    PRIMARY KEY(follower_id,followee_id),
    CHECK (followee_id <> follower_id),
    FOREIGN KEY (follower_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (followee_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE messages (
    message_id VARCHAR(128) PRIMARY KEY,
    sender_id VARCHAR(128) NOT NULL,
    receiver_id VARCHAR(128) NOT NULL,
    content VARCHAR(1024) NOT NULL,
    send_at TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (sender_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE open_messages (
    message_id VARCHAR(128) PRIMARY KEY,
    sender_id VARCHAR(128) NOT NULL,
    content VARCHAR(1024) NOT NULL,
    send_at TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (sender_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- view
CREATE VIEW book_author (
    book_id,
    book_name,
    author_name
) AS 
SELECT
    b.book_id,
    b.name AS book_name,
    (SELECT a.name FROM authors AS a WHERE a.author_id = b.author_id) AS author_name
FROM books AS b;

-- index
CREATE INDEX books_01 ON books (author_id);
CREATE INDEX books_genres_01 ON books_genres (book_id);
CREATE INDEX books_genres_02 ON books_genres (genre_id);
CREATE INDEX reads_01 ON reads (user_id);
CREATE INDEX reads_02 ON reads (book_id);
CREATE INDEX follows_01 ON follows (follower_id);
CREATE INDEX follows_02 ON follows (followee_id);
CREATE INDEX messages_01 ON messages (sender_id);
CREATE INDEX messages_02 ON messages (receiver_id);
CREATE INDEX open_messages_01 ON open_messages (sender_id);
