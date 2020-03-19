SELECT DISTINCT id FROM book, record
WHERE book.id = record.book_id AND record.time > to_date('2016-10-31')
ORDER BY id;
