select distinct name 
from (select author, count(name) from book group by author) as pubs(name, sumbook)
where sumbook >= 2 order by name;