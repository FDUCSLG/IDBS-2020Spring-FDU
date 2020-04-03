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
	EvaluatorID   = "15307130201" // your student id, e.g. 18307130177
	SubmissionDir = "../../../ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
	User          = "root" // the user name to connect the database, e.g. root
	Password      = "123456" // the password for the user name, e.g. xxx
	// YOUR CODE END
)

// ConcurrentCompareAndInsert is similar with compareAndInsert in `main.go`, but it is concurrent and faster!
func ConcurrentCompareAndInsert(subs map[string]*Submission) {
	start := time.Now()
	// *** From the Golang docs *** "A defer statement postpones the execution of a function until the surrounding function returns, either normally or through a panic."
	defer func() {
		// Execute an SQL phrase to log in as user, and check error (same as in compareAndInsert)
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
		if err != nil {
			panic(nil)
		}
		// Execute an SQL phrase to query the number of rows in the table "comparison_result"
		rows, err := db.Query("SELECT COUNT(*) FROM comparison_result")
		if err != nil {
			panic(err)
		}
		// *** From the Golang docs *** "Every call to Scan, even the first one, must be preceded by a call to Next."
		rows.Next() // Moves the "row pointer" to the first row of "comparison_result"
		var cnt int  // Initialize an integer variable "cnt"
		err = rows.Scan(&cnt) // Try to copy the first row to the variable "cnt"
		if err != nil {
			panic(err)
		}
		// If the first row does not exist (i.e. the table has 0 rows now), clearly this function is not yet implemented
		if cnt == 0 {
			panic("ConcurrentCompareAndInsert Not Implemented")
		}
		fmt.Println("ConcurrentCompareAndInsert takes ", time.Since(start))
	}()
	// YOUR CODE BEGIN
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
	if err != nil {
		panic(err)
	} // If there is error, ...
	wg := sync.WaitGroup{}
	for submitter, sub := range subs {
		for comparer, sub2 := range subs { // For each submitter j and their submission in the array subs
			for i := 0; i < NumSQL; i++ { // For each SQL phrase (known a priori that there is 8 of them)
				var equal int // Initialize an integer variable named "equal"
				if reflect.DeepEqual(sub.sqlResults[i], sub2.sqlResults[i]) {
					equal = 1 // If the results of the i-th SQL phrase are "DeepEqual", ...
				} else {
					equal = 0
				}
				wg.Add(1)
				go func(submitter string, sub *Submission, comparer string, sub2 *Submission, i int) {
					s := fmt.Sprintf("INSERT INTO comparison_result VALUES ('%s', '%s', %d, %d)", submitter, comparer, i+1, equal)
					_, err := db.Exec(s) // Execute the SQL phrase
					if err != nil {
						fmt.Println(s)
						panic(err)
					}
					wg.Done()
				}(submitter, sub, comparer, sub2, i)
			}
			wg.Wait()
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
	// 1. Transform "comparison_result" (submitter, comparer, item, is_equal) to "V" (submitter, item, vote)
	// 2. Do a sub-query yielding the table "IV" (item, max_vote)
	// 3. Join V and IV to obtain a new table (submitter, item, vote, max_vote)
	// 4. The desired table "score" is simply (submitter, item, vote == max_vote, vote)
	SQL = fmt.Sprintf( "INSERT INTO score (submitter, item, score, vote) " +
		"WITH V (submitter, item, vote) " +
		"AS (" +
		"SELECT submitter, item, SUM(is_equal) " +
		"FROM comparison_result " +
		"GROUP BY submitter, item" +
		") " +
		"SELECT V.submitter, V.item, if(V.vote = max_vote, 1, 0), V.vote " +
		"FROM V, (" +
		"SELECT V.item, MAX(V.vote) AS max_vote " +
		"FROM V " +
		"GROUP BY V.item" +
		") AS IV " +
		"WHERE V.item = IV.item")

	//SQL = fmt.Sprintf()
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	rows, err := db.Query("SELECT submitter, item, score FROM score")
	if err != nil{
		panic(err)
	}
	defer rows.Close()
	for {
		next_row := rows.Next()
		if(!next_row){
			break
		}
		var (
			submitter string
			item int
			score int
		)
		rows.Scan(&submitter, &item, &score)
		subs[submitter].score[item] = score
	}
	// YOUR CODE END
}
