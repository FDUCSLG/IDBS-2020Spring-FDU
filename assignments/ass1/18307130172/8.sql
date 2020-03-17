select id, name, count(*) as num
from employee inner join record on id = employee_id
group by id
having num > 1
order by num desc;
