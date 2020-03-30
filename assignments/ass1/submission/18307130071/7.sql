SELECT DISTINCT id
FROM book, record
WHERE `id` = `book_id` AND `time`.AFTER (2016-10-31);