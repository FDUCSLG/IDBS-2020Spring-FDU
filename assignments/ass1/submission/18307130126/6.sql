SELECT DISTINCT X.publisher
FROM book AS X, book AS Y
WHERE X.publisher=Y.publisher AND X.id!=Y.id ORDER BY publisher;