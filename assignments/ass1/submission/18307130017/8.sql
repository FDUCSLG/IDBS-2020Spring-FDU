SELECT employee_id,X.name,COUNT(book_id) AS num
FROM record,employee AS X
WHERE X.id = employee_id
GROUP BY employee_id
ORDER BY num desc;