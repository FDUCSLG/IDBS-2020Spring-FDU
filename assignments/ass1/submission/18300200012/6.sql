SELECT publisher 
FROM book
GROUP BY publisher
HAVING count(id)>2
ORDER BY publisher; 
