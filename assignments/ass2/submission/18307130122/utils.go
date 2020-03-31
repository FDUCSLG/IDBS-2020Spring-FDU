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
	EvaluatorID   = "18307130122" // your student id, e.g. 18307130177
	SubmissionDir = "/home/chenyue/IDBS-Spring20-Fudan/assignments/ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
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
	wg := sync.WaitGroup{}
	wg.Add(36)
	for submitter, sub := range subs {
		go func(submitter string,sub *Submission){
			for comparer, sub2 := range subs {
				for i := 0; i < NumSQL ; i++{
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
			}
			wg.Done()
		}(submitter,sub)
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
	WITH Choose1 AS (
			SELECT submitter,item,COUNT(*) vote
			FROM comparison_result
			GROUP BY submitter,item,is_equal	
			HAVING SUM(is_equal)>0),
	     Choose2 AS(
			SELECT item,MAX(vote) maxvote
			FROM Choose1
			GROUP BY item)
	SELECT submitter,Choose1.item,(CASE WHEN vote=maxvote THEN 1 ELSE 0 END) AS score,vote
	FROM Choose2,Choose1
	WHERE Choose2.item=Choose1.item`
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	rows,err := db.Query("SELECT submitter,item,score FROM score")
	if err != nil {
		panic(err)
	}
	var a string
	var b int
	var c int
	for rows.Next(){
	rows.Scan(&a,&b,&c)
	subs[a].score[b]=c
	}
	// YOUR CODE END
}
