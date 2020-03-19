CREATE TABLE employee
  (id SMALLINT NOT NULL,
  name VARCHAR(32),
  office VARCHAR(32),
  age SMALLINT,
  PRIMARY KEY(id));
CREATE TABLE book
  (id SMALLINT NOT NULL,
  name VARCHAR(32),
  author VARCHAR(32),
  publisher VARCHAR(32),
  PRIMARY KEY(id));
CREATE TABLE record
  (book_id SMALLINT,
  employee_id SMALLINT,
  time DATE,
  PRIMARY KEY(book_id,employee_id),
  FOREIGN KEY(book_id)REFERENCES book(id),
  FOREIGN KEY(employee_id)REFERENCES employee(id));
