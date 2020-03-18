select distinct name
from employee,record
where id=employee_id
group by id
	having sum(book_id)>2
order by name;