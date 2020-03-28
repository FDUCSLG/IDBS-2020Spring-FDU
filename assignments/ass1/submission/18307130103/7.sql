select employee.name
	from employee,record
	where record.time>'2016-10-31' and record.employee_id=employee.id
	group by employee.name
	order by max(employee.id);
