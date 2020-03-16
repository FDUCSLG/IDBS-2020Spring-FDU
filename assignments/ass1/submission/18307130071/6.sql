SELECT `publisher`
FROM book
GROUP BY `publisher`
HAVING count(*) > 1
ORDER BY `name`;
