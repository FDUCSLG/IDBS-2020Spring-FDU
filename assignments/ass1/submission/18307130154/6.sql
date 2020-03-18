SELECT DISTINCT publisher
FROM book AS X
WHERE 2 <= (SELECT COUNT(*) FROM book AS Y
	    WHERE Y.publisher = X.publisher)
ORDER BY publisher ;
