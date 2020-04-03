CREATE TABLE employee (	id INT NOT NULL,
 						name VARCHAR(32) NOT NULL,
 						office VARCHAR(32) NOT NULL,
 						age SMALLINT NOT NULL,
 						PRIMARY KEY(id),
 						CHECK(age BETWEEN 0 AND 100));

CREATE TABLE book (	id INT NOT NULL,
					name VARCHAR(32) NOT NULL,
					author VARCHAR(32) NOT NULL,
					publisher VARCHAR(32) NOT NULL,
					PRIMARY KEY(id));

CREATE TABLE record (	book_id INT NOT NULL,
						employee_id INT NOT NULL,
						time DATE NOT NULL,
						FOREIGN KEY(book_id) REFERENCES book(id),
						FOREIGN KEY(employee_id) REFERENCES employee(id),
						PRIMARY KEY(book_id,employee_id));
