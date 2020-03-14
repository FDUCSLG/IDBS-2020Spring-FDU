select distinct publisher
from book
group by publisher
having count(*) > 2
order by publisher
;