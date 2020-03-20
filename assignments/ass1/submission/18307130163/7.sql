SELECT DISTINCT id
FROM book, record
WHERE id = book_id AND time > '2016-10-31'
ORDER BY id;
