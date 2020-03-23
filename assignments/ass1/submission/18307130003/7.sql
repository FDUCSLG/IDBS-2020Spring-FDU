-- Query the id of all books that is borrowed after 2016-10-31, also the IDs
-- should be distinct (ordered by id)
SELECT
    DISTINCT book_id
FROM
    record
WHERE
    time > '2016-10-31'
ORDER BY
    book_id;