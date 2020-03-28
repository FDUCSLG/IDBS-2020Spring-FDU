with tmp(id,num)as(
    select employee_id,count(*)
    from record group by employee_id
)
select id,name,num
from tmp natural inner join employee
order by num desc;
