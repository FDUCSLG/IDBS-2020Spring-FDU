select B.id,B.name,A.num
from (select record.employee_id,count(*) as num
	from record
	group by record.employee_id
	having count(*)>1) as A
	inner join 
	employee as B
	on A.employee_id=B.id
order by num desc;
