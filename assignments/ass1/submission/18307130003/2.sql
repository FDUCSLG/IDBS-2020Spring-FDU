-- Query the name of employees with ID equals to 1 or 2 (order does not matter)
SELECT
    name
FROM
    employee
WHERE
    id IN (1, 2);