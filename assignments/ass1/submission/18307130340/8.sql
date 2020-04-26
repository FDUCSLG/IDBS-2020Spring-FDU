#Query for each employee who has borrowed book more than once, output the id, name, and number of borrow record (name the field num), ordered by num in descending order
select id, name, count(*) as num
from employee join record on id = employee_id
group by id
	having num > 1
order by num desc;