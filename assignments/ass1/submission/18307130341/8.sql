SELECT employee.id,employee.name,COUNT(*) AS num
FROM employee,record
WHERE employee.id=record.employee_id
GROUP BY employee.id
 HAVING COUNT(*)>1
ORDER BY num DESC;
