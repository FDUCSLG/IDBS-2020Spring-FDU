select employee.id, name, count(*) as num
from record join employee
on record.employee_id = employee.id
group by employee.id
	having count(*) > 1
order by num desc
