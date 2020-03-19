CREATE TABLE employee (
    PRIMARY KEY (id),
    id INT NOT NULL,
    name VARCHAR(32) NOT NULL,
    office VARCHAR(32),
    age INT,
    CONSTRAINT age_range CHECK (
        age BETWEEN 0
        AND 100
    )
);

CREATE TABLE book (
    PRIMARY KEY (id),
    id INT NOT NULL,
    name VARCHAR(32) NOT NULL,
    author VARCHAR(32),
    publisher VARCHAR(32)
);

CREATE TABLE record (
    PRIMARY KEY (book_id, employee_id),
    book_id INT NOT NULL,
    employee_id INT NOT NULL,
    time DATE,
    FOREIGN KEY (book_id) REFERENCES book(id),
    FOREIGN KEY (employee_id) REFERENCES employee(id)
);