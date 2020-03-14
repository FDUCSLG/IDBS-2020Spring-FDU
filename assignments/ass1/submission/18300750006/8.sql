SELECT employee.id, employee.name, COUNT(record.employee_id) 
	AS num FROM (employee INNER JOIN record
		ON employee.id = record.employee_id )
			GROUP BY employee.id
				HAVING COUNT(record.employee_id) >=2 
					ORDER BY COUNT(record.employee_id) DESC;