-- Query the name of all employees except the one whose ID is 1 (ordered by ID)
SELECT
    name
FROM
    employee
WHERE
    id <> 1
ORDER BY
    id;