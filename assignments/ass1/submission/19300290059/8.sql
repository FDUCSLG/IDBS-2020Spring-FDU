SELECT id, name, COUNT(id) AS num FROM employee INNER JOIN record ON id=employee_id GROUP BY id HAVING num > 1 ORDER BY num DESC;
