with tmp(publisher,num)as(
    select publisher,count(*) from book
    group by publisher
)
select publisher from tmp where num>2 order by publisher;
