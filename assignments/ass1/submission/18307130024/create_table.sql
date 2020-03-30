create table employee
(
    id int primary key,
    name varchar(32),
    office varchar(32),
    age int constraint cstage check(age <= 100 and age >= 0)
);

create table book
(
    id int primary key,
    name varchar(32),
    author varchar(32),
    publisher varchar(32)
);

create table record
(
    book_id int,
    employee_id int,
    time date,
    foreign key(book_id) references book(id) on delete cascade,
    foreign key(employee_id) references employee(id) on delete cascade
);
