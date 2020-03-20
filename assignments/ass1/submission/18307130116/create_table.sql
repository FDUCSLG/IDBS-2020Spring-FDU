CREATE TABLE employee
	(id	INTEGER,
	 name	VARCHAR(32) NOT NULL,
	 office VARCHAR(32) NOT NULL,
	 age SMALLINT,
	 PRIMARY KEY(id));
CREATE TABLE book
	(id	INTEGER,
	 name	VARCHAR(32) NOT NULL,
	 author	VARCHAR(32) NOT NULL,
	 publisher VARCHAR(32) NOT NULL,
	 PRIMARY KEY(id));
CREATE TABLE record
	(book_id INTEGER,
	 employee_id INTEGER,
	 time	 DATE,
	 FOREIGN KEY(book_id) REFERENCES book(id),
	 FOREIGN KEY(employee_id) REFERENCES employee(id));

