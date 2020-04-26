SELECT `id`, `name`, num
FROM ( SELECT `id`, `name`, COUNT(*) AS num FROM record GROUP BY employee_id HAVING num > 1 ) AS A, employee
WHERE `id` = `employee_id`
ORDER BY num DESC;