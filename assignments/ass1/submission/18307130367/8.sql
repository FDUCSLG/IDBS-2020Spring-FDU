select id,name,count(*) as num

from employee join record on employee.id=record.employee_id

group by id

   having num>1

order by num DESC;
