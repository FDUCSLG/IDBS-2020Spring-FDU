-- Query all fields for employees whose name started with J (ordered by age)
SELECT
    *
FROM
    employee
WHERE
    name LIKE 'J%'
ORDER BY
    age;