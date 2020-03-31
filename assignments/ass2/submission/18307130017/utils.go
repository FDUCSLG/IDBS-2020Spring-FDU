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
	EvaluatorID   = "18307130017" // your student id, e.g. 18307130177
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
	detecterr:=func(err error){
		if err != nil { 
			panic(err) 
		}
	}
	type comp_result struct{
		submitter,comparer string
		id int
		equal int
	}
	ch:=make(chan comp_result,32)
	wg := sync.WaitGroup{}
	max:=0
	for range subs{
		max++
	}
	db,err:= sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
	detecterr(err)
	for i:=1;i<=max;i++{
		wg.Add(1)
		go func(){
			defer wg.Done()
			tx,err:=db.Begin()
			detecterr(err)
			defer func(){
				err:=tx.Commit()
				detecterr(err)
			}()
			stmt,err:=tx.Prepare("INSERT INTO comparison_result VALUES (?,?,?,?)")
			defer stmt.Close()
			detecterr(err)
			for result:=range ch{
				_,err:=stmt.Exec(result.submitter,result.comparer,result.id,result.equal)
				detecterr(err)
			}
		}()
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
				ch<-comp_result{submitter,comparer,i+1,equal}
			}
		}
	}
	close(ch)
	wg.Wait()
	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	SQL = "SELECT 1" // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN
	SQL = `INSERT INTO score(submitter,item,score,vote)
		 WITH
		 	submissions AS(
				SELECT submitter,item,SUM(is_equal) AS vote
				FROM comparison_result
				GROUP BY submitter,item
			),
			highest AS (
				SELECT item,MAX(vote) AS maxvote
				FROM submissions
				GROUP BY item
			)
		 SELECT 
			submitter,
			submissions.item,
			if(submissions.vote = highest.maxvote,1,0) AS score,
			vote
		 FROM 
			submissions,
			highest
		 WHERE submissions.item = highest.item`
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	scores,err:=db.Query("SELECT submitter,item,score FROM score")
	if err != nil { 
		panic(err) 
	}
	defer scores.Close()
	for {
		nxt:=scores.Next()
		if !nxt{
			break
		}
		var submitter string
		var score,item int
		scores.Scan(&submitter,&item,&score)
		subs[submitter].score[item]=score
	}
	// YOUR CODE END
}
