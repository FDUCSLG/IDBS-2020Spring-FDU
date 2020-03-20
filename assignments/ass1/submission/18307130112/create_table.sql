create table if not exists employee
	(id int not null,
	name varchar(32) not null,
	office varchar(32) not null,
	age int not null ,check(age>=0 and age<=100),
	primary key(id));
create table if not exists book
	(id int not null,
	name varchar(32) not null,
	author varchar(32) not null,
	publisher varchar(32) not null,
	primary key(id));
create table if not exists record
	(book_id int not null,
	employee_id int not null,
	time date,
	foreign key (book_id) references book(id),
	foreign key (employee_id) references employee(id));