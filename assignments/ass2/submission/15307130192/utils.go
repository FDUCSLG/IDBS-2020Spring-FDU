package main

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	// YOUR CODE BEGIN remove the follow packages if you don't need them

	// YOUR CODE END

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/jmoiron/sqlx"
)

var (
	// YOUR CODE BELOW
	EvaluatorID   = "15307130192" // your student id, e.g. 18307130177
	SubmissionDir = "../../../ass1/submission" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
	User          = "root" // the user name to connect the database, e.g. root
	Password      = "qwer1234" // the password for the user name, e.g. xxx
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
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
	if err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	cap_goroutine := make(chan int, 50)
	for submitter, sub := range subs {
		for comparer, sub2 := range subs {
			for i := 0; i < NumSQL; i++ {
				cap_goroutine <- 0
				wg.Add(1)
				go func(submitter string, comparer string, sub *Submission, sub2 *Submission, i int) {
					defer wg.Done()
					var equal int
					if reflect.DeepEqual(sub.sqlResults[i], sub2.sqlResults[i]) {
						equal = 1
					} else {
						equal = 0
					}
					s := fmt.Sprintf("INSERT INTO comparison_result VALUES ('%s', '%s', %d, %d)", submitter, comparer, i+1, equal)
					_, err := db.Exec(s)

					if err != nil {
						fmt.Println(s)
						panic(err)
					}
					<-cap_goroutine
				}(submitter, comparer, sub, sub2, i)
			}
		}
	}
	wg.Wait()
	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	SQL = "SELECT 1" // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN
	SQL = `INSERT INTO score
(WITH score_raw (submitter, item, vote) AS (
SELECT submitter, item, sum(is_equal) AS vote FROM comparison_result
GROUP BY submitter, item),
max_vote (item, vote_max) AS (
SELECT  item, MAX(vote) FROM score_raw
GROUP BY item)
(SELECT submitter, item, '0' AS score, vote FROM score_raw AS X
WHERE vote < (SELECT vote_max FROM max_vote AS Y WHERE X.item = Y.item))
UNION
(SELECT submitter, item, '1' AS score, vote FROM score_raw AS X
WHERE vote = (SELECT vote_max FROM max_vote AS Y WHERE X.item = Y.item)))`
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	rows, err := db.Query("SELECT submitter, item, score FROM score")
	if err != nil {
		panic(err)
	}
	var submitter string
	var item, score int
	for rows.Next() {
		rows.Scan(&submitter, &item, &score)
		subs[submitter].score[item] = score
	}
	// YOUR CODE END
}
