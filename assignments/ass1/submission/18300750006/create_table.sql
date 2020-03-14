CREATE TABLE employee
		( id VARCHAR(32) NOT NULL PRIMARY KEY,
		  name CHAR(32) NOT NULL,
          office CHAR(32) NOT NULL,
          age INT NOT NULL CHECK(age >= 0 AND age <= 100));
          
CREATE TABLE book
		( id VARCHAR(32) NOT NULL PRIMARY KEY,
		  name CHAR(32) NOT NULL,
          author CHAR(32) NOT NULL,
          publisher CHAR(32) NOT NULL);
          
CREATE TABLE record
		( book_id VARCHAR(32) NOT NULL,
		  employee_id CHAR(32) NOT NULL,
          time DATE NOT NULL,
          FOREIGN KEY (book_id) REFERENCES book(id),
          FOREIGN KEY (employee_id) REFERENCES employee(id));