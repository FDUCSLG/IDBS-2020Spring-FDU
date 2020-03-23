CREATE TABLE `employee`
	(`id` VARCHAR(32) NOT NULL,
	`name` VARCHAR(32) NOT NULL,
	`office` VARCHAR(32) NOT NULL,
	`age` INTEGER NOT NULL,
	PRIMARY KEY(`id`));
CREATE TABLE `book`
	(`id` VARCHAR(32) NOT NULL,
	`name` VARCHAR(32) NOT NULL,
	`author` VARCHAR(32) NOT NULL,
	`publisher` VARCHAR(32) NOT NULL,
	PRIMARY KEY(`id`));
CREATE TABLE `record` 
	(book_id VARCHAR(32) NOT NULL,
	employee_id VARCHAR(32) NOT NULL,
	`time` DATE NOT NULL,
	FOREIGN KEY (book_id) REFERENCES book (`id`),
	FOREIGN KEY (employee_id) REFERENCES employee (`id`));
