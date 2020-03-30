#Query the name of all employees except the one whose ID is 1 (ordered by ID)
select name
from employee
where id <> 1
order by id;