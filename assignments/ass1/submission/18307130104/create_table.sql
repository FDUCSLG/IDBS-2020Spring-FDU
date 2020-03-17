CREATE TABLE employee
	(id CHAR(4) NOT NULL,
	 name VARCHAR(32),
	 office VARCHAR(32),
	 age SMALLINT,
	 PRIMARY KEY(id));
CREATE TABLE book
	(id CHAR(4) NOT NULL,
	 name VARCHAR(32),
	 author VARCHAR(32),
	 publisher VARCHAR(32),
	 PRIMARY KEY(id));
CREATE TABLE record
	(book_id CHAR(4) NOT NULL,
	 employee_id CHAR(4) NOT NULL,
	 time DATE,
	 FOREIGN KEY (book_id) REFERENCES book(id),
	 FOREIGN KEY (employee_id) REFERENCES employee(id));