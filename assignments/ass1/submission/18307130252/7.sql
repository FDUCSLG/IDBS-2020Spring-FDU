select distinct book_id
from record
where time > str_to_date('2016-10-31', "%Y-%m-%d")
order by book_id;

