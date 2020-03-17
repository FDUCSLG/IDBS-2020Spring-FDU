SELECT *
FROM employee AS X
WHERE X.name REGEXP '^J'
ORDER BY X.age;