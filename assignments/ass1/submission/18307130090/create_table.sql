create table employee(
    id int,
    name char(32),
    office char(32),
    age int,
    primary key(id),
    check (id between 0 and 100)
);
create table book(
    id int,
    name char(32),
    author char(32),
    publisher char(32),
    primary key(id)
);
create table record(
    book_id int,
    employee_id int,
    time date,
    primary key(book_id,employee_id),
    foreign key(book_id)references book(id),
    foreign key(employee_id)references employee(id)
);
