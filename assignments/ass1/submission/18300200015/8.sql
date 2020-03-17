select employee.id,employee.name,count(record.employee_id)
from record,employee
where record.employee_id=employee.id
group by employee.id
having count(record.employee_id)>=2
order by count(record.employee_id) desc;