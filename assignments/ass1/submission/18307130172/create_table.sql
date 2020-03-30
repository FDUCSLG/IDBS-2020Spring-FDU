create database if not exists ass1;
use ass1;

create table if not exists employee (
    id int auto_increment primary key,
    name varchar(32),
    office varchar(32),
    age int,
    check (age between 0 and 100)
);
create table if not exists book (
    id int auto_increment primary key,
    name varchar(32),
    author varchar(32),
    publisher varchar(32)
);
create table if not exists record (
    book_id int,
    employee_id int,
    time date,
    primary key (book_id, employee_id),
    foreign key (book_id)
        references book (id),
    foreign key (employee_id)
        references employee (id)
);
