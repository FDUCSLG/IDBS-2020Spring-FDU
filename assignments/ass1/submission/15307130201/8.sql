SELECT e.id, e.name, COUNT(*) as num
FROM employee AS e, record
WHERE e.id = record.employee_id
GROUP BY e.id
HAVING num > 1
ORDER BY num DESC;