create table employee(
	id varchar(20) not null,
    name varchar(32),
    office varchar(32),
    age int check(age between 0 and 100),
    primary key(id)
);
create table book(
	id varchar(20) not null,
    name varchar(32),
    author varchar(32),
    publisher varchar(32),
    primary key(id)
);
create table record(
	book_id varchar(20),
    employee_id varchar(20),
    time DATE,
    foreign key(book_id) references book(id),
    foreign key(employee_id) references employee(id)
);