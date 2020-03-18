SELECT id , name , COUNT(*) AS num
FROM employee, record
WHERE id = employee_id
GROUP BY id
	HAVING COUNT(*) >1
ORDER BY COUNT(*) DESC;
