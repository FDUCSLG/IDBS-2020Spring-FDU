SELECT DISTINCT book_id
FROM book,record
WHERE book.id = record.book_id AND record.time > '2016-10-31'
ORDER BY book_id 
