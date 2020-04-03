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
	EvaluatorID   = "18307130104" // your student id, e.g. 18307130177
	SubmissionDir = "../../../ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
	User          = "" // the user name to connect the database, e.g. root
	Password      = "" // the password for the user name, e.g. xxx
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
	type safedb struct{
		db *sql.DB
		mux sync.Mutex
	}
	const tot = 30
	var dbarr[tot] safedb
	for i := 0; i < tot; i++ {
		dbarr[i] = safedb{ db: db }
	}
	var wg sync.WaitGroup
	cur := 0
	for submitter, sub := range subs {
		for comparer, sub2 := range subs {
			for i := 0; i < NumSQL; i++ {
				var equal int
				if reflect.DeepEqual(sub.sqlResults[i], sub2.sqlResults[i]) {
					equal = 1
				} else {
					equal = 0
				}
				wg.Add(1)
				cur ++
				if cur == tot {
					cur = 0
				}
				go func(wg *sync.WaitGroup, cur1 int, submitter string, comparer string, equal int, i int){
					dbarr[cur1].mux.Lock()
					defer dbarr[cur1].mux.Unlock()
					s := fmt.Sprintf("INSERT INTO comparison_result VALUES ('%s', '%s', %d, %d)", submitter, comparer, i+1, equal)
					_, err := dbarr[cur1].db.Exec(s)
					if err != nil {
						fmt.Println(s)
						panic(err)
					}
					wg.Done()
				}(&wg, cur, submitter, comparer, equal, i)
			}
		}
		wg.Wait()
	}
	/*db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
	if err != nil {
		panic(err)
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
				s := fmt.Sprintf("INSERT INTO comparison_result VALUES ('%s', '%s', %d, %d)", submitter, comparer, i+1, equal)
				_, err := db.Exec(s)
				if err != nil {
					fmt.Println(s)
					panic(err)
				}
			}
		}
	}*/
	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	SQL = "SELECT 1" // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN
	SQL = `INSERT INTO score (submitter, item, score, vote) 
			With equal_comp(submitter, item, vote) AS 
		   (SELECT submitter, item, COUNT(DISTINCT comparer) AS vote 
		    FROM comparison_result 
		    WHERE  is_equal = 1 
		    GROUP BY submitter, item)

		  	(SELECT submitter, item, 1 AS score, vote 		
		  	 FROM equal_comp AS X
		   	 WHERE vote >= ALL(SELECT Y.vote
		   					   FROM equal_comp AS Y 
		   					   WHERE Y.item = X.item))
		  	UNION
		  	(SELECT submitter, item, 0 AS score, vote 
		  	 FROM equal_comp AS X
		   	 WHERE vote < SOME(SELECT Y.vote
		   					   FROM equal_comp AS Y
		   					   WHERE Y.item = X.item)) `
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	score, scoreerr := db.Query("SELECT submitter, item, score FROM score")
	if scoreerr != nil {
		panic(scoreerr)
	}

	for score.Next() {
		var id string
		var itid, sc int
		score.Scan(&id, &itid, &sc)
		subs[id].score[itid] = sc
	}
	// YOUR CODE END
}
