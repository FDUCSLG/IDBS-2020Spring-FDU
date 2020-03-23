CREATE TABLE employee
	(id INT NOT NULL,
     name CHAR (32),
     office CHAR (32),
     age INT CHECK(age between 0 and 100),
     PRIMARY KEY (id));
     
CREATE TABLE book
	(id INT NOT NULL,
     name CHAR (32),
     author CHAR (32),
     publisher CHAR (32),
     PRIMARY KEY (id));
     
CREATE TABLE record
	(book_id INT,
     employee_id INT,
     time DATE,
     FOREIGN KEY (book_id) REFERENCES book(id),
     FOREIGN KEY (employee_id) REFERENCES employee(id));
     
    