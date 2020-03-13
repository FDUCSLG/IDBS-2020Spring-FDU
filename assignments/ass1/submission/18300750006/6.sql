SELECT DISTINCT X.publisher FROM book AS X , book AS Y 
	WHERE X.name != Y.name AND X.publisher = Y.publisher ORDER BY publisher;