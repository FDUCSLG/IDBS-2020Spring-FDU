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
	EvaluatorID   = "18307130128" // your student id, e.g. 18307130177
	SubmissionDir = "/home/hang_y/IDBS-Spring20-Fudan/assignments/ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
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
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
        if err != nil {
                panic(err)
        }
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(100)
	wg := sync.WaitGroup{}
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
				_, err = db.Exec(s)
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
	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	 // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN
	SQL = "INSERT INTO score(submitter,item,score,vote) SELECT V1.submitter, V1.item, 1 AS score, V1.vote AS vote FROM (SELECT submitter,item, SUM(is_equal) AS vote FROM comparison_result GROUP BY submitter, item) V1 WHERE V1.vote = (SELECT MAX(V2.vote) FROM (SELECT submitter,item, SUM(is_equal) AS vote FROM comparison_result GROUP BY submitter, item) V2 WHERE V1.item = V2.item) UNION ALL SELECT V3.submitter, V3.item, 0 AS score, V3.vote AS vote FROM (SELECT submitter,item, SUM(is_equal) AS vote FROM comparison_result GROUP BY submitter, item) V3 WHERE V3.vote < (SELECT MAX(V4.vote) FROM (SELECT submitter,item, SUM(is_equal) AS vote FROM comparison_result GROUP BY submitter, item) V4 WHERE V3.item = V4.item)"
	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
	type scoreInfo struct{
                Submitter string `db:"submitter"`
                Item int `db:"item"`
                Score int `db:"score"`
                Vote int `db:"vote"`
        }
        for submitter, sub := range subs{
                for i := 0; i < NumSQL; i++ {
                        var scoredata *scoreInfo = new(scoreInfo)
                        s := fmt.Sprintf("SELECT * FROM score WHERE submitter = '%s' AND item = %d", submitter, i + 1)
			err := db.Get(scoredata, s)
                        if err == nil {
                                sub.score[i + 1] = scoredata.Score
                        }else{
                                fmt.Println(s)
                                panic(err)
                        }
                }
        }
	// YOUR CODE END
}
