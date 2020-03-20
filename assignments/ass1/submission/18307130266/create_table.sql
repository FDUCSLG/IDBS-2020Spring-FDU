create table employee(    
    id     int not null,
    name   varchar(32),
    office varchar(32),
    age    int(8),
    primary key(id)
);

create table book(
    id        int not null,
    name      varchar(32),
    author    varchar(32),
    publisher varchar(32),
    primary key(id)
);

create table record(
    book_id     int not null,
    employee_id int not null,
    time        date,
    foreign key(book_id) references book(id),
    foreign key(employee_id) references employee(id)
);
