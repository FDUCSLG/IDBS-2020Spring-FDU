SELECT DISTINCT id 
FROM book, record
WHERE id = book_id AND to_days(time) - to_days('2016-10-31') > 0
ORDER BY id ASC;
