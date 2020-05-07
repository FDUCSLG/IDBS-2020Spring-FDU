package main

import (
	"fmt"
	"reflect"
	"time"

	// YOUR CODE BEGIN remove the follow packages if you don't need them

	// YOUR CODE END

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/jmoiron/sqlx"
)

var (
	// YOUR CODE BELOW
	EvaluatorID   = "18300200012" // your student id, e.g. 18307130177
	SubmissionDir = "/home/l1q/DB/IDBS-Spring20-Fudan/assignments/ass1/submission/"
	User          = "root"       // the user name to connect the database, e.g. root
	Password      = "llq0215llq" // the password for the user name, e.g. xxx
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
		panic(nil)
	}

	ch := make(chan int, NumSQL)
	quit := make(chan int)
	for submitter, sub := range subs {
		for comparer, sub2 := range subs {
			for j := 0; j < NumSQL; j++ {
				ch <- j
				go func() {
					i := <-ch
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

					if i == NumSQL-1 {
						quit <- 0
					}
				}()
			}
			<-quit
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

	SQL = "INSERT INTO score SELECT X.submitter,X.item,IF(X.vote >= Y.highest,1,0),X.vote FROM (SELECT submitter,item,count(is_equal) AS vote FROM comparison_result WHERE is_equal = 1 GROUP BY 1,2) AS X,(SELECT item,max(vote) AS highest FROM (SELECT submitter,item,count(is_equal) AS vote FROM comparison_result WHERE is_equal = 1 GROUP BY 1,2) AS Z GROUP BY 1) AS Y WHERE X.item = Y.item;"

	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	for _, sub := range subs {
		for i := 1; i <= NumSQL; i++ {
			row, err := db.Query(fmt.Sprintf("SELECT score FROM score WHERE submitter = '%s' AND item = %d", sub.submitter, i))
			if err != nil {
				panic(err)
			}
			for row.Next() {
				_ = row.Scan(&sub.score[i])
			}
		}
	}
	// YOUR CODE END
}
