SELECT DISTINCT id
FROM book
WHERE id IN(SELECT book_id FROM record 
	    WHERE  time > '2016-10-31');
