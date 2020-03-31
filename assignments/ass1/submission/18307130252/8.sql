with X as(
	select employee_id as id, count(*) as num
	from record
	group by employee_id)
select employee.id, name, X.num
from employee left outer join X
on employee.id = X.id
where X.num > 1
order by X.num desc;
