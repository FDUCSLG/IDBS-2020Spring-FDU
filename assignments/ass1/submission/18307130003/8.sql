-- Query for each employee who has borrowed book more than once, output the id,
-- name, and number of borrow record (name the field num), ordered by num in
-- descending order.
SELECT
    id,
    name,
    count(book_id) AS num
FROM
    employee,
    record
WHERE
    id = employee_id
GROUP BY
    employee_id
HAVING
    num > 1
ORDER BY
    num DESC;