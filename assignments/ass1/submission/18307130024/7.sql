select distinct book_id
from record
where unix_timestamp(time) > unix_timestamp('2016-10-31')
order by book_id
;