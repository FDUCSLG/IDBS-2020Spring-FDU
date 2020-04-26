package main

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	// YOUR CODE BEGIN remove the follow packages if you don't need them

	// YOUR CODE END

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/jmoiron/sqlx"
)

var (
	// YOUR CODE BELOW
	EvaluatorID   = "18307130172"               // your student id, e.g. 18307130177
	SubmissionDir = "../../../ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
	User          = "whofucksit"                // the user name to connect the database, e.g. root
	Password      = "19260817"                  // the password for the user name, e.g. xxx
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
		fmt.Println("(error) in database connection")
		panic(err)
	}

	stmt, stmtErr := db.Preparex("INSERT INTO comparison_result VALUES (?,?,?,?)")
	if stmtErr != nil {
		fmt.Println("(error) in statement preparation")
		panic(stmtErr)
	}

	type result struct {
		uid    string
		vid    string
		caseId int
		result int
	}

	const numWorkers = 32
	jobs := make(chan result, numWorkers)
	var w sync.WaitGroup
	worker := func() {
		defer w.Done()

		tx, err := db.Beginx()
		if err != nil {
			fmt.Println("(error) in transaction initialization")
			panic(err)
		}

		defer func() {
			err := tx.Commit()
			if err != nil {
				fmt.Println("(error) in transaction commit")
				panic(err)
			}
		}()

		exec := tx.Stmtx(stmt)
		defer exec.Close()

		for ret := range jobs {
			_, execErr := exec.Exec(ret.uid, ret.vid, ret.caseId, ret.result)
			if execErr != nil {
				fmt.Println("(error) in performing insertions")
				panic(execErr)
			}
		}
	}

	for i := 0; i < numWorkers; i++ {
		w.Add(1)
		go worker()
	}

	for uid, u := range subs {
		for vid, v := range subs {
			for i := 0; i < NumSQL; i++ {
				var ret int
				if reflect.DeepEqual(u.sqlResults[i], v.sqlResults[i]) {
					ret = 1
				} else {
					ret = 0
				}

				jobs <- result{uid, vid, i + 1, ret}
			}
		}
	}

	close(jobs)
	w.Wait()
	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	SQL = "SELECT 1" // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN
	SQL = `INSERT INTO score
	WITH Results AS (
		SELECT submitter, item, COUNT(comparer) AS vote
		FROM comparison_result
		WHERE is_equal = 1
		GROUP BY submitter, item
	), Standing AS (
		SELECT item, MAX(vote) AS highest
		FROM Results
		GROUP BY item
	)
	SELECT
		submitter, item,
		IF(vote = highest, 1, 0) AS score,
		vote
	FROM Results JOIN Standing USING(item)`
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	rows, rowsErr := db.Query("SELECT submitter, item, score FROM score")
	if rowsErr != nil {
		panic(rowsErr)
	}

	for {
		ok := rows.Next()
		if !ok {
			break
		}

		var uid string
		var item, score int
		rows.Scan(&uid, &item, &score)
		subs[uid].score[item] = score
	}
	// YOUR CODE END
}
