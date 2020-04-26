
CREATE TABLE employee (
	`id` int ,
    `name` char(32) ,
    `office` char(32) ,
    `age` int CHECK (age>0&&age<100) ,
    PRIMARY KEY (`id`)
);
CREATE TABLE book (
	`id` char(32) ,
    `name` char(32) ,
    `author` char(32) ,
    `publisher` char(32) ,
    PRIMARY KEY (`id`)
);
CREATE TABLE record (
	`book_id` char(32) ,
    `employee_id` int ,
    `time` DATE,
    FOREIGN KEY (`book_id`) REFERENCES book(`id`),
    FOREIGN KEY (`employee_id`) REFERENCES employee(`id`)
)