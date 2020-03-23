SELECT id, name, count(*) AS num
FROM employee JOIN record ON id = employee_id
GROUP BY id
HAVING num >= 2
ORDER BY num DESC;
