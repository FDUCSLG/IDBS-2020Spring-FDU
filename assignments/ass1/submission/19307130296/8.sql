SELECT id, name, num
FROM employee
GROUP BY id
HAVING COUNT(
	SELECT *
    FROM record
    WHERE record.employee_id=employee.id
) AS num>1
ORDER BY num