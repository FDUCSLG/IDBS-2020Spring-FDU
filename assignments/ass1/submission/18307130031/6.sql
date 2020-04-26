SELECT DISTINCT X.publisher
FROM book AS X, book AS Y, book AS Z
WHERE X.publisher = Y.publisher AND X.publisher = Z.publisher AND X.name != Y.name AND Y.name != Z.name AND Z.name != X.name
ORDER BY X.publisher ASC;
