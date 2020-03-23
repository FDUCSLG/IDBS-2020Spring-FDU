select id 

from book

where id in(select book_id

from record

where time>2016-10-31

);

order by id 
