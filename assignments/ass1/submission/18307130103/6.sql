select distinct X.publisher
	from book as X
	where 2<=(select count(*) 
		from book as Y 
		where X.publisher=Y.publisher)
	order by X.publisher;
