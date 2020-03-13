select distinct A.publisher
from book as A, book as B, book as C
where A.publisher = B.publisher
    and B.publisher = C.publisher
    and A.id != B.id
    and B.id != C.id
    and A.id != C.id
order by A.publisher;