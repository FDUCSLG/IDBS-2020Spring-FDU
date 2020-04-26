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
	EvaluatorID   = "18307130123" // your student id, e.g. 18307130177
	SubmissionDir = "../../../ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
	User          = "root" // the user name to connect the database, e.g. root
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
	wg:=sync.WaitGroup{}
        for submitter, sub := range subs {
		wg.Add(1)
		go func(submitter string, sub *Submission){
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
		wg.Done()
                }(submitter, sub)
        }
	wg.Wait()
	//The number of open files (-n) is no more than 1024
	//Actually we can change the limited number into a bigger one, but I give it up finally since I'm not sure if it'll work on another computer.
	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	SQL = "SELECT 1" // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN
	SQL = `INSERT INTO score
	SELECT submitter,VoteResult.item,(CASE WHEN VoteResult.vote = MaxVoteResult.max_vote THEN 1 ELSE 0 END) AS score,VoteResult.vote
	FROM (SELECT submitter,item,SUM(is_equal) AS vote
		FROM comparison_result
		GROUP BY submitter,item)AS VoteResult,
		
		(SELECT item,MAX(VoteResult.vote) AS max_vote 
		FROM (SELECT submitter,item,SUM(is_equal) AS vote FROM comparison_result GROUP BY submitter,item)AS VoteResult 
		GROUP BY item) AS MaxVoteResult
	WHERE VoteResult.item = MaxVoteResult.item`
	//The version of MySQL is too old, but I don't want to update now since I need to copy the database.So I export to table VoteResult twice. 
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	rows,err:=db.Query(`SELECT submitter,item,score FROM score`)
	if err!=nil{
		panic(err)
	}
	for rows.Next(){
		var sub string
		var it int
		var sc int
		err=rows.Scan(&sub,&it,&sc)
		if err!=nil{
			panic(err)
		}
		subs[sub].score[it]=sc
	}
	rows.Close()
	if err = rows.Err();err != nil{
		panic(err)
	}
	// YOUR CODE END
}
