create table employee(
	id integer not null,
    name varchar(32),
    office varchar(32),
    age smallint,
    primary key(id)
);
    
create table book(
	id integer not null,
    name varchar(32),
    author varchar(32),
    publisher varchar(32),
    primary key(id)
);

create table record(
	book_id integer,
    employee_id integer,
    time date,
    foreign key(book_id) references book(id),
    foreign key(employee_id) references employee(id),
    primary key(book_id, employee_id)
);