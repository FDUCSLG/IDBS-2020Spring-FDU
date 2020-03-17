select employee.id, employee.name, count(*) as num
from employee inner join record on employee.id = record.employee_id
group by employee.id having num > 1 order by num DESC;