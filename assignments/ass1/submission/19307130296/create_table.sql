CREATE TABLE employee(
	id int(32) not null auto_increment primary key,
    name varchar(32),
    office varchar(32),
    age int(3)
)

CREATE TABLE book(
	id int(32) not null auto_increment primary key,
    name varchar(32),
    author varchar(32),
    publisher varchar(32)
)

CREATE TABLE record(
	time DATE,
    book_id int(32) not null,
    employee_id int(32) not null,
    foreign key(book_id) references book(id),
    foreign key(employee_id) references employee(id)
)