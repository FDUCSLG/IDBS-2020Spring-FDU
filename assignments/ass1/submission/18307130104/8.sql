SELECT id, name, COUNT(DISTINCT book_id) AS num
FROM employee, record
WHERE id = employee_id
GROUP BY id
	HAVING num > 1
ORDER BY num DESC;