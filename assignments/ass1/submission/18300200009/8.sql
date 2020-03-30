SELECT id,name,count(book_id) AS num
FROM record,employee
WHERE id=employee_id
GROUP BY id
	HAVING count(book_id)>1
ORDER BY num;