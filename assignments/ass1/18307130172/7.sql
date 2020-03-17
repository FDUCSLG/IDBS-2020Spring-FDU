select distinct id
from book inner join record on id = book_id
where time > '2016-10-31'
order by id;
