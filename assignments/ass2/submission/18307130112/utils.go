package main

import (
	"fmt"
	//"go/ast"
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
	EvaluatorID   = "18307130112" // your student id, e.g. 18307130177
	SubmissionDir = "F:/Github_Clone/IDBS-Spring20-Fudan/assignments/ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
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
	//init
	type answer struct{
		submitters string
		comparers string
		ids int
		equals int
	}
	const channalSize = 40
	channal := make(chan answer, channalSize)
	var wg sync.WaitGroup
	memNum := 0
	for range subs{
		memNum ++
	}
	//fileopen
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
	defer db.Close()
	if err != nil{
		panic(err)
	}
	//insert goroutine
	functions := func() {
		defer wg.Done()
		Tx, err := db.Begin()
		if err != nil{
			panic(err)
		}
		defer func(){
			err := Tx.Commit()
			if err != nil{
				panic(err)
			}
		}()
		stmt, err := Tx.Prepare("insert into comparison_result values (?,?,?,?)")
		defer stmt.Close()
		if err != nil{
			panic(err)
		}
		for ans := range channal{
			_, err := stmt.Exec(ans.submitters, ans.comparers, ans.ids, ans.equals)
			if err != nil{
				panic(err)
			}
		}
	}
	//run goroutine & compare & pass value
	for i := 1; i <= memNum; i++{
		wg.Add(1)
		go functions()
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
				channal <- answer{submitter, comparer, i+1, equal}
			}
		}
	}
	close(channal)
	wg.Wait()
	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	SQL = "SELECT 1" // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN
	SQL =
		"insert into score(submitter, item,  score, vote) with " +
			"submission as (select submitter, item, sum(is_equal) as vote from comparison_result group by submitter, item), " +
			"standard as (select item, max(vote) as maximum from submission group by item) " +
		"select submitter, submission.item, if(submission.vote = standard.maximum, 1, 0) as score, vote from submission, standard where submission.item = standard.item"
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	row, err := db.Query("select submitter, item, score from score")
	if err != nil{
		panic(err)
	}
	defer row.Close()
	for {
		flag := row.Next()
		if(!flag){
			break
		}
		var submitters string
		var items, scores int
		row.Scan(&submitters, &items, &scores)
		subs[submitters].score[items] = scores
	}
	// YOUR CODE END
}
