CREATE TABLE employee(
	id INT,
	office CHAR(32),
	name CHAR(32),
	age INTEGER,
	PRIMARY KEY (id)
);

CREATE TABLE book(
	id INT,
	name CHAR(32),
	author CHAR(32),
	publisher CHAR(32),
	PRIMARY KEY (id)
);

CREATE TABLE record(
	time DATE,
	book_id INT,
	employee_id INT,
	FOREIGN KEY (book_id) REFERENCES book(id),
	FOREIGN KEY (employee_id) REFERENCES employee(id)
);
