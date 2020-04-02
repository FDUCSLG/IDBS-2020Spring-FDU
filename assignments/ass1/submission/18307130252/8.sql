select employee.id, name, X.num
from employee left outer join ((select employee_id as id, count(*) as num
			    				from record
								group by employee_id) as X )
on employee.id = X.id
where X.num > 1
order by X.num desc;
