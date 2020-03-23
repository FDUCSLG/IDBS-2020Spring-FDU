SELECT employee.id,name,count(*) AS num
FROM employee LEFT JOIN record ON employee.id = record.employee_id
GROUP BY employee.id
  HAVING count(*) > 1
ORDER BY num DESC;
