SELECT DISTINCT publisher
FROM book
GROUP BY publisher
  HAVING COUNT(*)>1
ORDER BY publisher ;
