#Qeury all fields for employees whose name started with J (ordered by age)
select *
from employee
where name like 'J%'
order by age;