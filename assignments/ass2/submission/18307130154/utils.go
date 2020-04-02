package main

import (
	"fmt"
	"time"
	"reflect"
	"sync"
	// YOUR CODE BEGIN remove the follow packages if you don't need them

	// YOUR CODE END

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/jmoiron/sqlx"
)

var (
	// YOUR CODE BELOW
	EvaluatorID   = "18307130154" // your student id, e.g. 18307130177
	SubmissionDir = "../../../ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
	User          = "root" // the user name to connect the database, e.g. root
	Password      = "080090" // the password for the user name, e.g. xxx
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

	stmt, stmtErr := db.Preparex("INSERT INTO comparison_result VALUES (?,?,?,?)")//YuChuLi
	if stmtErr != nil {
		panic(stmtErr)
	}
	
	type result struct{
		submitter string
		comparer string
		question int
		equal int
	}//result for one question: used in "INSERT INTO comparison_result VALUES (?,?,?,?)"
	ch := make(chan result, 32)//suppose at most 32 students
	var w sync.WaitGroup
	Insert := func() { //Insert operation//concurrent
		defer w.Done()
		w.Add(1)
		tx, err := db.Beginx()//make a  shiwu 
		if err != nil {
			panic(err)
		}
		defer func() {
			err := tx.Commit()
			if err != nil {
				panic(err)
			}
		}()

		exec := tx.Stmtx(stmt)//execute one insert from channel
		defer exec.Close()

		for insert := range ch {
			_, err := exec.Exec(insert.submitter, insert.comparer, insert.question, insert.equal)
			if err != nil {
				panic(err)
			}
		}
	}

	for i := 0; i < 32; i++ {		
		go Insert()
	}
	for submitter, sub := range subs {
		for comparer, comp := range subs {
			for i := 0; i < NumSQL; i++ {//like main.go/CompareAndInsert
				var equal int
				if reflect.DeepEqual(sub.sqlResults[i], comp.sqlResults[i]) {//compare
					equal = 1
				} else {
					equal = 0
				}

				ch <- result{submitter, comparer, i + 1, equal}
			}
		}
	}

	close(ch)
	w.Wait()

	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	SQL = "SELECT 1" // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN

	// using with for SQL8.* //////

	/*SQL = "INSERt INTO score (submitter,item,score,vote) "+
	"WITH result AS (SELECT submitter, item, COUNT(DISTINCT comparer) AS vote " +
			"FROM comparison_result WHERE is_equal = 1 GROUP BY submitter, item), " +
	"     refer AS (SELECT item, MAX(vote) AS max_vote FROM result GROUP BY item) " +
	"SELECT submitter, result.item, IF(vote = max_vote,1,0) AS score, vote " +
	"FROM result ,refer WHERE refer.item = result.item" */

	//without 'with' for mysql 5.*//////
	SQL = "INSERt INTO score (submitter,item,score,vote) "+
	"SELECT submitter, result.item, IF(vote = max_vote,1,0) AS score, vote " +
	"FROM (SELECT submitter, item, COUNT(DISTINCT comparer) AS vote FROM comparison_result WHERE is_equal = 1 GROUP BY submitter, item) AS result , (SELECT item, MAX(vote) AS max_vote FROM (SELECT submitter, item, COUNT(DISTINCT comparer) AS vote FROM comparison_result WHERE is_equal = 1 GROUP BY submitter, item) as a GROUP BY item) AS refer " +
	"WHERE refer.item = result.item" 
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	row,err := db.Query("SELECT submitter, item, score FROM score")
	if err != nil{
		panic(err)
	}
	for {
		if !row.Next() {//end of table
			break
		}
		var submitter string
		var item, score int
		row.Scan(&submitter, &item, &score)
		subs[submitter].score[item] = score
	}
	// YOUR CODE END
}
