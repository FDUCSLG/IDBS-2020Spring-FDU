CREATE TABLE if not exists employee
		( id int ,
		  name CHAR(32) NOT NULL,
          office CHAR(32) NOT NULL,
          age INT NOT NULL CHECK(age >= 0 AND age <= 100),
          primary key(id)
);
          
CREATE TABLE if not exists book
		( id int not null,
		  name CHAR(32) NOT NULL,
          author CHAR(32) NOT NULL,
          publisher CHAR(32) NOT NULL,
          primary key(id));
          
CREATE TABLE if not exists record
		( book_id VARCHAR(32) NOT NULL,
		  employee_id CHAR(32) NOT NULL,
          time DATE NOT NULL,
          FOREIGN KEY (book_id) REFERENCES book(id),
          FOREIGN KEY (employee_id) REFERENCES employee(id));
