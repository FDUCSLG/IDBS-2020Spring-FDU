INSERT INTO employee (id, name, office, age) VALUES
    (1, 'Jones', 'CA', 30),
    (2, 'Alice', 'SJ', 33),
    (3, 'Bob', 'NY', 29),
    (4, 'Jack', 'CN', 50);

INSERT INTO book (id, name, author, publisher) VALUES
    (1, 'T1', 'A1', 'McGraw-Hill'),
    (2, 'T1', 'A1', 'McGraw-Hill'),
    (3, 'T2', 'A2', 'McGraw-Hill'),
    (4, 'T2', 'A2', 'McGraw-Hill'),
    (5, 'T3', 'A3', 'McGraw-Hill'),
    (6, 'T3', 'A3', 'McGraw-Hill'),
    (11, 'T4', 'A4', 'Fudan'),
    (12, 'T4', 'A4', 'Fudan'),
    (13, 'T5', 'A5', 'Fudan'),
    (14, 'T6', 'A6', 'Fudan'),
    (15, 'T7', 'A7', 'Fudan'),
    (16, 'T8', 'A8', 'Fudan');

INSERT INTO record (employee_id, book_id, time) VALUES
    (1, 1, '2016-3-24'),
    (2, 2, '2016-3-24'),
    (2, 3, '2016-3-24'),
    (1, 4, '2016-3-25'),
    (3, 5, '2016-3-25'),
    (1, 6, '2016-3-25'),
    (1, 11, '2018-4-18'),
    (4, 12, '2018-4-18'),
    (4, 13, '2018-4-18'),
    (1, 14, '2018-4-18'),
    (1, 15, '2018-4-18'),
    (2, 16, '2018-4-18'),
    (2, 1, '2018-3-18');
