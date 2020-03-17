CREATE TABLE employee
    (
        id INT NOT NULL, 
        name VARCHAR(32),
        office VARCHAR(32),
        age INT,
        check (age <= 100 and age >= 0),
        PRIMARY KEY (id)
    );

CREATE TABLE book
    (
        id INT NOT NULL, 
        name VARCHAR(32),
        author VARCHAR(32),
        publisher VARCHAR(32),
        PRIMARY KEY (id)
    );

CREATE TABLE record
    (
        book_id INT,
        employee_id INT,
        time DATE,
        FOREIGN KEY(book_id) references book(id),
        FOREIGN KEY(employee_id) references employee(id)
    );
