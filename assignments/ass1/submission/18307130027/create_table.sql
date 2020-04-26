-- CREATE DATABASE IF NOT EXISTS ass1;
-- USE ass1;
CREATE TABLE employee (
  id INT NOT NULL,
  name VARCHAR(32),
  office VARCHAR(32),
  age INT,
  PRIMARY KEY (id));
CREATE TABLE book (
    id INT NOT NULL,
    name VARCHAR(32),
    author VARCHAR(32),
    publisher VARCHAR(32),
    primary key (id)
);
CREATE TABLE record (
	book_id INT,
    employee_id INT,
    time DATE,
    FOREIGN KEY (book_id)
		REFERENCES book(id),
	FOREIGN KEY (employee_id)
		REFERENCES employee(id)
);