SELECT employee.id, employee.name, COUNT(*) AS num
FROM employee INNER JOIN record ON employee.id=record.employee_id
GROUP BY employee.id
HAVING num>1
ORDER BY num DESC;