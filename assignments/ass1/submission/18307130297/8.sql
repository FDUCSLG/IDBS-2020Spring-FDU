select id,name,count(*)
from employee,record
where id=employee_id
group by employee_id
	having sum(book_id)>1
order by count(*) desc;