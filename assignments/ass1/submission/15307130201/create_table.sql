CREATE TABLE employee (
	id INT(16),
	name VARCHAR(32),
	office VARCHAR(32),
	age INT(8),
	PRIMARY KEY (id)
);
CREATE TABLE book (
	id INT(16),
	name VARCHAR(32),
	author VARCHAR(32),
	publisher VARCHAR(32),
	PRIMARY KEY (id)
);
CREATE TABLE record (
	book_id INT(16),
	employee_id INT(16),
	time DATE,
	FOREIGN KEY (book_id) REFERENCES book (id),
	FOREIGN KEY (employee_id) REFERENCES employee (id)
);
