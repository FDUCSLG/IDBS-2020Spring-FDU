SELECT DISTINCT id,name,COUNT(*) AS num
FROM employee,record
WHERE employee.id=record.employee_id
GROUP BY id
	HAVING num > 1
ORDER BY num DESC;
