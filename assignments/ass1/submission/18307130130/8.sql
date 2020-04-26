SELECT employee.id,name,COUNT(*) as num
FROM employee LEFT OUTER JOIN record ON employee.id=record.employee_id
GROUP BY employee.id
  HAVING COUNT(*)>1
ORDER BY 3 DESC;

