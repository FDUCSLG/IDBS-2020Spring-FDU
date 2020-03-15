select id, name, count(*) as num
from employee, record
where employee.id = record.employee_id
group by id
	having num > 1
order by 3;
