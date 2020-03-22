CREATE TABLE employee(
	id INT NOT NULL PRIMARY KEY,
	name VARCHAR(32),
	office VARCHAR(32),
	age INT,
	CHECK (0 < age AND age < 100)
);

CREATE TABLE book(
	id INT NOT NULL PRIMARY KEY,
	name VARCHAR(32),
	author VARCHAR(32),
	publisher VARCHAR(32)
);

CREATE TABLE record(
	book_id INT NOT NULL,
	employee_id INT NOT NULL,
	time DATE,
	FOREIGN KEY (book_id) REFERENCES book(id),
	FOREIGN KEY (employee_id) REFERENCES employee(id)
);
