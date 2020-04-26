SELECT `name`
FROM employee
WHERE NOT EXISTS ( SELECT *
				   FROM employee
                   WHERE id = '1' )
ORDER BY `id`;