create table employee(
	id int,
    name varchar(32),
    office varchar(32),
    age int check (age between 0 and 100),
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
    time date,
    foreign key(book_id) references book(id),
    foreign key(employee_id) references employee(id)
);