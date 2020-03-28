-- Query the names of all publishers, if one publisher has more than two books
-- in the database, output the name of the publisher only once (ordered by name,
-- ascii order)
SELECT
    DISTINCT publisher
FROM
    book
ORDER BY
    publisher;