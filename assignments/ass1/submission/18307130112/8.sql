select id, name, sum from employee,
(select employee_id, count(employee_id) from record group by employee_id) 
as borrows(employee_id, sum) 
where sum > 0 and id = employee_id
order by sum desc;