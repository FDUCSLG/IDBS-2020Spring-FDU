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
	EvaluatorID   = "18300750006" // your student id, e.g. 18307130177
	SubmissionDir = "../../../ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
	User          = "root" // the user name to connect the database, e.g. root
	Password      = ""//omit // the password for the user name, e.g. xxx
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
	var wg sync.WaitGroup
	wg.Add(len(subs))
	
	for submitter, sub := range subs {
		//var lock sync.Mutex	
		go func(submitter string , sub *Submission){
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
		
				//wg.Add(1)
				//go func(p int){		
				//lock.Lock()
				//defer lock.Unlock()	
				//wg.Done()
				//}(i)
				//wg.Wait()
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
	SQL = "INSERT INTO score(submitter,item,score,vote)" +
		"(SELECT A.submitter,A.item,1 AS score,vote FROM " +
		"(SELECT submitter,item,SUM(is_equal) AS vote FROM comparison_result GROUP BY submitter,item) as A " +
		"WHERE vote >= all(select sum(B.is_equal) AS vote FROM comparison_result as B " +
		"WHERE A.item = B.item GROUP BY B.submitter,B.item) ORDER BY A.submitter,A.item )" +
		" union " +
		"(SELECT C.submitter,C.item,0 AS score,vote FROM " +
		"(SELECT submitter,item,SUM(is_equal) AS vote FROM comparison_result GROUP BY submitter,item) as C " +
		"WHERE vote < some(select sum(D.is_equal) AS vote FROM comparison_result as D " +
		"WHERE C.item = D.item GROUP BY D.submitter,D.item) ORDER BY C.submitter,C.item) "

	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	var score int

	for _,sub := range subs {

		for i := 1; i < 9; i++ {

			errTables := db.QueryRow("SELECT score FROM score WHERE submitter = ? AND item = ?", sub.submitter,i).Scan(&score)
			sub.score[i] = score
			if errTables != nil {
			}


		}
	}

	// YOUR CODE END
}
