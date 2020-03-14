select id, name, count(employee_id) as num
from record, employee
where employee_id = id
group by employee_id
having num > 1
order by num desc
;