SELECT DISTINCT id, name, COUNT(id) AS num
FROM employee, record
WHERE id = employee_id
GROUP BY employee_id
	HAVING num >= '2'
ORDER BY num DESC;
