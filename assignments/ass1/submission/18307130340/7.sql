#Query the id of all books that is borrowed after 2016-10-31, also the IDs should be distinct (ordered by id)
select distinct book_id
from record
where time > '2016-10-31'
order by book_id;