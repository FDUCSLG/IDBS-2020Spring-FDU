package main

import (
	"fmt"
	"time"

	// YOUR CODE BEGIN remove the follow packages if you don't need them
	"reflect"
	"sync"
	// YOUR CODE END
	_ "github.com/go-sql-driver/mysql"
	sql "github.com/jmoiron/sqlx"
)

var (
	// YOUR CODE BELOW
	EvaluatorID   = "19307130296"               // your student id, e.g. 18307130177
	SubmissionDir = "../../../ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
	User          = "orz"                      // the user name to connect the database, e.g. root
	Password      = "123"                    // the password for the user name, e.g. xxx
	// YOUR CODE END
)

// ConcurrentCompareAndInsert is similar with compareAndInsert in `main.go`, but it is concurrent and faster!
func ConcurrentCompareAndInsert(subs map[string]*Submission) {
	start := time.Now()
	defer func() {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
		if err != nil {
			panic(nil)
		}
		rows, err := db.Query("SELECT COUNT(*) FROM comparison_result")
		if err != nil {
			panic(err)
		}
		rows.Next()
		var cnt int
		err = rows.Scan(&cnt)
		if err != nil {
			panic(err)
		}
		if cnt == 0 {
			panic("ConcurrentCompareAndInsert Not Implemented")
		}
		fmt.Println("ConcurrentCompareAndInsert takes ", time.Since(start))
	}()
	// YOUR CODE BEGIN
	check := func(err error) {
		if(err != nil) {
			panic(err)
		}
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
	check(err)
	stmt, err := db.Preparex("INSERT INTO comparison_result VALUES (?,?,?,?)")
	check(err)
	defer stmt.Close()

	type node struct {
		submitter string
		comparer string
		item int
		is_equal int
	}
	ch := make(chan node, 20)

	var w sync.WaitGroup
	solve := func() {
		defer w.Done()
		for temp := range ch {
			_, err := stmt.Exec(temp.submitter, temp.comparer, temp.item, temp.is_equal)
			check(err)
		}
	}

	for i := 0; i < 50; i++ {
		w.Add(1)
		go solve()
	}

	for submitter, sub := range subs {
		for comparer, sub2 := range subs {
			for i := 0; i < NumSQL; i++ {
				var equal int
				if reflect.DeepEqual(sub.sqlResults[i], sub2.sqlResults[i]) {
					equal = 1
				} else {
					equal = 0
				}
				ch <- node{submitter, comparer, i + 1, equal}
			}
		}
	}

	close(ch)
	w.Wait()
	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	SQL = "SELECT 1" // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN
	SQL =`INSERT INTO score
			(WITH scoreCTE (submitter,item,vote) AS (
				SELECT distinct submitter, item, SUM(is_equal)
				FROM comparison_result
				GROUP BY submitter, item
			)
			((SELECT submitter, item, '1' AS score, vote
			FROM scoreCTE AS X
			WHERE vote=(SELECT MAX(vote)
						FROM scoreCTE AS Y
						WHERE Y.item=X.item))
			UNION
			(SELECT submitter, item, '0' AS score, vote
			FROM scoreCTE AS X
			WHERE vote!=(SELECT MAX(vote)
						 FROM scoreCTE AS Y
						 WHERE Y.item=X.item))))`
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	rows, err := db.Query("SELECT submitter, item, score From score")
	if err != nil {
		panic(err)
	}
	var item, score int
	var submitter string
	for rows.Next() {
		rows.Scan(&submitter, &item, &score)
		subs[submitter].score[item] = score
	}
	// YOUR CODE END
}
