package main

import (
	"fmt"
	"time"

	// YOUR CODE BEGIN remove the follow packages if you don't need them
	"sync"
	"reflect"
	// YOUR CODE END

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/jmoiron/sqlx"
)

var (
	// YOUR CODE BELOW
	EvaluatorID   = "18307130163" // your student id, e.g. 18307130177
	SubmissionDir = "D:/study/大二下/数据库引论/lab/IDBS-Spring20-Fudan/assignments/ass1/submission" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
	User          = "root" // the user name to connect the database, e.g. root
	Password      = "123456" // the password for the user name, e.g. xxx
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
	db.SetMaxOpenConns(100)
	defer db.Close()

	wg := sync.WaitGroup{}
	defer wg.Wait()
	for submitter, sub := range subs {
		for comparer, sub2 := range subs {
			wg.Add(1)
			go func() {
				for i := 0; i < NumSQL; i++ {
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
				}
				wg.Done()
			}()
		}
	}
	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	SQL = "SELECT 1" // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN
	SQL = `
	INSERT INTO score
	WITH votes AS (
		SELECT submitter, item, COUNT(comparer) AS vote
		FROM comparison_result
		WHERE is_equal = 1
		GROUP BY submitter, item
	), standard AS (
		SELECT item, MAX(vote) AS highest
		FROM votes
		GROUP BY item
	)
	SELECT submitter, item, IF(vote = highest, 1, 0) AS score, vote
	FROM votes JOIN standard USING(item)
	`
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	rows, rowsErr := db.Query("SELECT submitter, item, score FROM score")
	if rowsErr != nil {
		panic(rowsErr)
	}

	for {
		ok := rows.Next()
		if !ok {
			break
		}

		var uid string
		var item, score int
		rows.Scan(&uid, &item, &score)
		subs[uid].score[item] = score
	}
	// YOUR CODE END
}
