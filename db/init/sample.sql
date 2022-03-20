
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
INSERT INTO open_messages (message_id,sender_id,content,send_at) VALUES (2,2,'I have no books I want to read.',CURRENT_TIMESTAMP);
