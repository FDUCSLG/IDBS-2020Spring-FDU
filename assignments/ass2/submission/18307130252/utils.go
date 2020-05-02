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
	EvaluatorID   = "18307130252"               // your student id, e.g. 18307130177
	SubmissionDir = "../../../ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
	User          = "kleinercubs"               // the user name to connect the database, e.g. root
	Password      = "kleinercubs"               // the password for the user name, e.g. xxx
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
		panic(nil)
	}

	type data struct {
		sub  string
		cmp  string
		id   int
		same int
	}

	bufferedChan := make(chan data, 128)
	wg := sync.WaitGroup{}

	for thread := 0; thread < 16; thread++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for tmp := range bufferedChan {
				s := fmt.Sprintf("INSERT INTO comparison_result VALUES ('%s', '%s', %d, %d)", tmp.sub, tmp.cmp, tmp.id, tmp.same)
				_, err := db.Exec(s)
				if err != nil {
					fmt.Println(s)
					panic(err)
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for submitter, sub := range subs {
			for comparer, sub2 := range subs {
				for i := 0; i < NumSQL; i++ {
					var equal int
					if reflect.DeepEqual(sub.sqlResults[i], sub2.sqlResults[i]) {
						equal = 1
					} else {
						equal = 0
					}
					bufferedChan <- data{submitter, comparer, i + 1, equal}
				}
			}
		}
		close(bufferedChan)
	}()

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
	WITH RESULT AS(
		SELECT submitter, item, count(is_equal) AS vote
		FROM comparison_result
		WHERE is_equal = 1
		GROUP BY submitter, item
	)
	SELECT submitter, item, IF (vote IN (SELECT MAX(Tmp.vote) FROM RESULT AS Tmp WHERE Tmp.item = RESULT.item), 1, 0),vote
	FROM RESULT`
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN

	rows, err := db.Query("SELECT * FROM score")

	if err != nil {
		panic(err)
	}

	var submitter string
	var item, score, vote int

	for rows.Next() {
		rows.Scan(&submitter, &item, &score, &vote)
		subs[submitter].score[item] = score
	}

	// YOUR CODE END
}
