use ass1;

create table employee(
	id int primary key,
	name varchar(32),
	office varchar(32),
	age int
);

create table book(
	id int primary key,
	name varchar(32),
	author varchar(32),
	publisher varchar(32)
);

create table record(
	book_id int,
	employee_id int,
	time date,
	primary key (book_id, employee_id),
	foreign key (book_id)
		references book (id),
	foreign key (employee_id)
		references employee (id)
);
