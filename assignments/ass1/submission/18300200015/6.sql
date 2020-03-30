select distinct x.publisher
from book as x,book as y
where x.publisher=y.publisher and x.name!=y.name
order by publisher;