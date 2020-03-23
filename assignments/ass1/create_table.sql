CREATE TABLE employee (
    id INT PRIMARY KEY,
    name VARCHAR(32),
    office VARCHAR(32),
    age INT
);

CREATE TABLE book (
    id INT PRIMARY KEY,
    name VARCHAR(32),
    author VARCHAR(32),
    publisher VARCHAR(32)
);

CREATE TABLE record (
    book_id INT,
    employee_id INT,
    time DATE,
    FOREIGN KEY (book_id) REFERENCES book(id),
    FOREIGN KEY (employee_id) REFERENCES employee(id)
);
