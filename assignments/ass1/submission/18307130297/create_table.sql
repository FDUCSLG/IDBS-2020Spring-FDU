create table employee(
	id int,
    name varchar(32),
    office varchar(32),
    age int,
    primary key(id)
);

create table book(
	id int,
    name varchar(32),
    author varchar(32),
    publisher varchar(32),
    primary key(id)
);

create table record(
	book_id int,
    employee_id int,
    time DATE,
    foreign key(book_id)references book(id),
    foreign key(employee_id)references employee(id)
);