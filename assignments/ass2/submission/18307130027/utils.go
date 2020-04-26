package main

import (
	"fmt"
	"time"

	// YOUR CODE BEGIN remove the follow packages if you don't need them
	"reflect"
	// "sync"
	// YOUR CODE END

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/jmoiron/sqlx"
)

var (
	// YOUR CODE BELOW
	EvaluatorID   = "18307130027" // your student id, e.g. 18307130177
	SubmissionDir = "../../../ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
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
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(0)
	// newch := make(chan bool)
	for submitter, sub := range subs {
		for comparer, sub2 := range subs {
			for i := 0; i < NumSQL; i++ {
				var equal int
				if reflect.DeepEqual(sub.sqlResults[i], sub2.sqlResults[i]) {
					equal = 1
				} else {
					equal = 0
				}
				s := fmt.Sprintf("INSERT INTO comparison_result VALUES ('%s', '%s', %d, %d)", submitter, comparer, i+1, equal)
				for db.Stats().InUse > 64{
				}
				go func (sqlst string) {
					_, err := db.Exec(sqlst)
					if err != nil {
						fmt.Println(sqlst)
						panic(err)
					}	
				} (s)
			}
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
insert into score (submitter, item, score, vote)
with scr as (select submitter, item, sum(is_equal)as cnt from comparison_result group by submitter, item)
(SELECT submitter, item, 1, cnt		
FROM scr AS X
WHERE cnt >= ALL(SELECT cnt
				FROM scr 
				WHERE item = X.item))
				union
(SELECT submitter, item, 0, cnt
FROM scr AS X
WHERE cnt < SOME(SELECT cnt
				FROM scr 
				WHERE item = X.item))
	`
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	rows, err := db.Query("select submitter, item, score from score")
	if err != nil {
		panic(err)
	}
	var sid string
	var item, score int
	// rows.Next()
	for rows.Next() {
		rows.Scan(&sid, &item, &score)
		// fmt.Println(sid, item, score)
		subs[sid].score[item] = score
	}
	// YOUR CODE END
}
