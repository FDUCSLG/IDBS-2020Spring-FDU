SELECT id, name, count(id) AS num
FROM employee, record
WHERE employee.id = record.employee_id
GROUP BY id
   HAVING count(id) > 1
ORDER BY num
