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
	EvaluatorID   = "18307130130" // your student id, e.g. 18307130177
	SubmissionDir = "../../../ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
	User          = "root" // the user name to connect the database, e.g. root
	Password      = "lrc2521686" // the password for the user name, e.g. xxx
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
	//fmt.Println("1")
	for submitter, sub := range subs {
		for comparer, sub2 := range subs {
			wg.Add(1)
			submitter := submitter
			comparer := comparer
			sub := sub
			sub2 := sub2
			go func(){
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
		wg.Wait()
	}
	//fmt.Println("2")
	
	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	SQL = "SELECT 1" // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN
	SQL = fmt.Sprintf( "INSERT INTO score (submitter,item,score,vote) "+
	"SELECT submitter,item,is_equal,RESULT.VOTE "+
	"FROM comparison_result, "+
	"(SELECT comparer as com,item as ite,COUNT(*) as VOTE "+
	"FROM comparison_result "+
	"WHERE is_equal=1 "+
	"GROUP BY comparer,item) AS RESULT WHERE comparer=%s AND comparison_result.submitter=RESULT.com AND comparison_result.item=RESULT.ite;",EvaluatorID)

	 
	// YOUR CODE END
	return SQL
}
func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	for suber,sub := range subs{		
			var s string
			suber := suber
			s = fmt.Sprintf("SELECT item,score FROM score WHERE submitter=%s",suber)
			rows,err := db.Query(s)
			if err != nil {
				panic(err)
			}
			var value int
			var flag int
			for rows.Next(){
				err = rows.Scan(&flag,&value)
				if err != nil {
					panic(err)
				}
				sub.score[flag] = value
			}
		
	}
	// YOUR CODE END
}

