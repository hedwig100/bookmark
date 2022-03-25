
-- users
SELECT username,password FROM users where user_id = '0';

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
WHERE r.user_id = '0'
ORDER BY r.read_at;

-- follows
SELECT followee_id FROM follows WHERE follower_id = '0';

-- messages
SELECT sender_id,receiver_id,content,send_at FROM messages 
WHERE sender_id = '0' AND receiver_id = '1';

-- open messages
SELECT sender_id,content,send_at FROM open_messages 
WHERE sender_id IN (SELECT followee_id FROM follows WHERE follower_id = '0');
