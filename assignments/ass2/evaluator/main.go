package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // this is the driver for mysql
	sql "github.com/jmoiron/sqlx" // this is the connector, both package are external packages that you need to use `go install` to install before use
)

const (
	NumSQL        = 8
)

type Submission struct {
	db         *sql.DB
	submitter  string
	workingDir string
	score      []int
	sqls       []string
	sqlResults [][]map[string]interface{}
}

func mustExecute(db *sql.DB, SQLs []string) {
	for _, s := range SQLs {
		_, err := db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
}

func newSubmission(submitter string) (sub *Submission, err error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", User, Password))
	if err != nil {
		return
	}
	sub = &Submission{
		db:         db,
		submitter:  submitter,
		workingDir: filepath.Join(SubmissionDir, submitter),
		score:      make([]int, NumSQL + 1),
		sqls:       make([]string, NumSQL),
		sqlResults: make([][]map[string]interface{}, NumSQL),
	}
	mustExecute(db, []string{
		fmt.Sprintf("DROP DATABASE IF EXISTS ass1_%s", submitter),
		fmt.Sprintf("CREATE DATABASE ass1_%s", submitter),
		fmt.Sprintf("USE ass1_%s", submitter),
	})
	return
}

func executeSQLsFromFile(filePath string, db *sql.DB) error {
	binData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	SQLs := string(binData)
	for _, s := range strings.Split(SQLs, ";") {
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		_, err := db.Exec(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sub *Submission) close() {
	_, _ = sub.db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS ass1_%s", sub.submitter))
}

func (sub *Submission) testCreateTable() {
	var err error
	err = executeSQLsFromFile(filepath.Join(sub.workingDir, "create_table.sql"), sub.db)
	if err != nil || executeSQLsFromFile(filepath.Join(sub.workingDir, "../../insert_data.sql"), sub.db) != nil {
		if err == nil {
			err = executeSQLsFromFile(filepath.Join(sub.workingDir, "../../insert_data.sql"), sub.db)
		}
		fmt.Println(sub.submitter, " ", 0, " ", err)
		tmp, _ := newSubmission(sub.submitter)
		*sub = *tmp
		err = executeSQLsFromFile(filepath.Join(sub.workingDir, "../../create_table.sql"), sub.db)
		if err != nil {
			panic(err)
		}
		err = executeSQLsFromFile(filepath.Join(sub.workingDir, "../../insert_data.sql"), sub.db)
		if err != nil {
			panic(err)
		}
		sub.score[0] = 0
	} else {
		sub.score[0] = 1
	}
}

func createSubmissions() map[string]*Submission {
	dirs, err := ioutil.ReadDir(SubmissionDir)
	if err != nil {
		panic(err)
	}
	subs := make(map[string]*Submission)
	wg := sync.WaitGroup{}
	for _, di := range dirs {
		if di.IsDir() {
			sub, err := newSubmission(di.Name())
			if err == nil {
				subs[di.Name()] = sub
				wg.Add(1)
				go func() {
					sub.testCreateTable()
					sub.querySQLs()
					sub.close()
					wg.Done()
				}()
			} else {
				fmt.Println(err)
			}
		}
	}
	wg.Wait()
	fmt.Println("submission created")
	return subs
}

func (sub *Submission) querySQLs() {
	for i := 0; i < NumSQL; i++ {
		file := filepath.Join(sub.workingDir, fmt.Sprintf("%d.sql", i+1))
		binData, err := ioutil.ReadFile(file)
		if err == nil {
			sub.sqls[i] = string(binData)
			s := string(binData)
			rows, err := sub.db.Queryx(s)
			if err == nil {
				results := make([]map[string]interface{}, 0)
				for rows.Next() {
					dest := make(map[string]interface{})
					_ = rows.MapScan(dest)
					results = append(results, dest)
				}
				sub.sqlResults[i] = results
			} else {
				fmt.Println(sub.submitter, " ", i + 1, " ", err)
			}
		} else {
			fmt.Println(sub.submitter, " ", i + 1, " ", err)
		}
	}
}

func compareAndInsert(subs map[string]*Submission) {
	start := time.Now()
	defer func() {
		fmt.Println("CompareAndInsert takes ", time.Since(start))
	}()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
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
				s := fmt.Sprintf("INSERT INTO result VALUES ('%s', '%s', %d, %d)", submitter, comparer, i+1, equal)
				_, err := db.Exec(s)
				if err != nil {
					fmt.Println(s)
					panic(err)
				}
			}
		}
	}
}

func insertComparisonResult(subs map[string]*Submission) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", User, Password))
	if err != nil {
		panic(err)
	}
	mustExecute(db, []string{
		fmt.Sprintf("DROP DATABASE IF EXISTS ass1_result_evaluated_by_%s", EvaluatorID),
		fmt.Sprintf("CREATE DATABASE ass1_result_evaluated_by_%s", EvaluatorID),
		fmt.Sprintf("USE ass1_result_evaluated_by_%s", EvaluatorID),
		"CREATE TABLE comparision_result (submitter CHAR(11), comparer CHAR(11), item INT, is_equal INT)",
	})
	ConcurrentCompareAndInsert(subs)
	// uncomment the following line to see how slow the single-thread insert is
	// compareAndInsert(subs)
}


func createScoreTable(subs map[string]*Submission) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
	if err != nil {
		panic(err)
	}
	mustExecute(db, []string{
		"CREATE TABLE score (submitter CHAR(11), item int, score int, vote int)",
		GetScoreSQL(),
	})
	rows, err := db.Query("SELECT COUNT(*) FROM score")
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
		panic("GetScoreSQL Not Implemented")
	}
	fmt.Printf("GetScoreSQL inserted %d records into score\n", cnt)

	GetScore(db, subs)
	for _, sub := range subs {
		fmt.Print(sub.submitter, " ", sub.score[0])
		for i := 1; i <= NumSQL; i++ {
			fmt.Print(" ", sub.score[i])
		}
		fmt.Println()
	}
}

func main() {
	subs := createSubmissions()
	insertComparisonResult(subs)
	createScoreTable(subs)
}
