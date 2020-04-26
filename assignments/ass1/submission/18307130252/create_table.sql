create table employee(
	id varchar(32) not null,
	name varchar(32),
	office varchar(32),
	age smallint,
	primary key(id)
);

create table book(
	id varchar(32) not null,
	name varchar(32),
	author varchar(32),
	publisher varchar(32),
	primary key(id)
);

create table record(
	book_id varchar(32),
	employee_id varchar(32),
	time date,
	foreign key(book_id) references book(id),
	foreign key(employee_id) references employee(id)
);


