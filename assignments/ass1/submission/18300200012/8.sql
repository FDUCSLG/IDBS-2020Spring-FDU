SELECT id,name,count(id) AS num
FROM employee JOIN record ON id = employee_id
GROUP BY id,name
HAVING num>1
ORDER BY num DESC
