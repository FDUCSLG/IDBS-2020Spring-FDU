CREATE TABLE employee(
	id INTEGER,
	name CHAR(32) NOT NULL,
	office CHAR(32) NOT NULL,
	age CHAR(32) NOT NULL,
	PRIMARY KEY(id)
);
CREATE TABLE book(
	id INTEGER,
	name CHAR(32) NOT NULL,
	author CHAR(32) NOT NULL,
	publisher CHAR(32) NOT NULL,
	PRIMARY KEY(id)
);
CREATE TABLE record(
	book_id INTEGER,
	employee_id INTEGER,
	time DATE,
	-- PRIMARY KEY (book_id),
	FOREIGN KEY (book_id) REFERENCES book(id),
	FOREIGN KEY (employee_id) REFERENCES employee(id)
);