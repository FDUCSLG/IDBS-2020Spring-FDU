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

