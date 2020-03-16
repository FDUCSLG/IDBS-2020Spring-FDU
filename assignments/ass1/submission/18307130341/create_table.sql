CREATE TABLE employee
	(id INTEGER NOT NULL ,
	 name CHAR(32) NOT NULL,
	 office CHAR(32) NOT NULL,
	 age SMALLINT NOT NULL,
	 PRIMARY KEY(id));

CREATE TABLE book
	(id INTEGER NOT NULL,
	 name CHAR(32) NOT NULL,
	 author CHAR(32) NOT NULL,
	 publisher CHAR(32) NOT NULL,
	 PRIMARY KEY(id));

CREATE TABLE record
	(book_id INTEGER NOT NULL,
	 employee_id INTEGER NOT NULL,
	 time DATE NOT NULL,
	 FOREIGN KEY (employee_id) REFERENCES employee(id),
	 FOREIGN KEY (book_id) REFERENCES book(id));
