package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const letter = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const other = "0123456789+-*/=,.!?~<>;:$#@!^&()[]{}|"
const TestTimes = 100

func TestCreateTables(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CreateTables()
	if err != nil {
		t.Errorf("can't create tables")
	}
}

func RandomString(t, minlen, maxlen int) string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	var runes string
	if t != 0 {
		runes = letter + other
	} else {
		runes = letter
	}
	res := make([]byte, minlen+r.Intn(maxlen-minlen))
	for i := range res {
		res[i] = runes[r.Intn(len(runes))]
	}
	return string(res)
}

func TestAddStu(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	sqlExecute(lib.db, []string{
		fmt.Sprintf(" SET GLOBAL max_connections = %d", 3*TestTimes),
	})
	err := lib.CreateTables()
	if err != nil {
		t.Errorf("can't create tables")
	}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < TestTimes; i++ {
		account := fmt.Sprintf("stu%.7d", 1+r.Intn(9999999))
		password := RandomString(1, 4, 16)
		lib.AddStu(account, string(password))
		/*Successfully add ?*/
		row, err := lib.db.Query(fmt.Sprintf("SELECT count(*) FROM Student WHERE stu_id = '%s' ", account))
		if err != nil {
			panic(err)
		}
		var cnt int
		row.Next()
		_ = row.Scan(&cnt)
		if cnt != 1 {
			t.Errorf("account: '%s', password: '%s', failed to add.", account, string(password))
		}

	}

}

func TestAddLib(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CreateTables()
	if err != nil {
		t.Errorf("can't create tables")
	}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < TestTimes; i++ {
		account := fmt.Sprintf("%.5d", 1+r.Intn(99999))
		password := RandomString(1, 4, 16)
		lib.AddLib(account, string(password))
		/*Successfully add ?*/
		row, err := lib.db.Query(fmt.Sprintf("SELECT count(*) FROM Librarian WHERE lib_id = '%s'", account))
		if err != nil {
			panic(err)
		}
		var cnt int
		row.Next()
		_ = row.Scan(&cnt)
		if cnt != 1 {
			t.Errorf("account: '%s', password: '%s', failed to add.", account, string(password))
		}

	}

}

func TestAddBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	sqlExecute(lib.db, []string{
		fmt.Sprintf(" SET GLOBAL max_connections = %d", 10*TestTimes),
	})
	err := lib.CreateTables()
	if err != nil {
		t.Errorf("can't create tables")
	}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	CUR_ACCOUNT = "lib01"
	for i := 0; i < TestTimes; i++ {
		ISBN := fmt.Sprintf("987%.10d", 1+r.Intn(99999999))
		title := RandomString(0, 4, 16)
		auther := RandomString(0, 4, 16)
		same_cnt := r.Intn(10)
		for j := 0; j < same_cnt; j++ {
			lib.AddBook(ISBN, title, auther)
		}
	}
}

func TestRemoveBook(t *testing.T) {
	/*AddBook*/
	lib := Library{}
	lib.ConnectDB()
	sqlExecute(lib.db, []string{
		fmt.Sprintf(" SET GLOBAL max_connections = %d", 10*TestTimes),
	})
	err := lib.CreateTables()
	if err != nil {
		t.Errorf("can't create tables")
	}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	CUR_ACCOUNT = "lib01"
	for i := 0; i < TestTimes; i++ {
		ISBN := fmt.Sprintf("987%.10d", 1+r.Intn(99999999))
		title := RandomString(0, 4, 16)
		auther := RandomString(0, 4, 16)
		same_cnt := r.Intn(10)
		for j := 0; j < same_cnt; j++ {
			lib.AddBook(ISBN, title, auther)
		}
	}
	/*RemoveBook*/
	remv_id := make([]string, TestTimes)
	for i := 0; i < TestTimes; i++ {
		remv_id[i] = fmt.Sprintf("987%.10dn%.3d", 1+r.Intn(99999999), 1+r.Intn(10))
	}
	/* Randomly Select */
	rows, err := lib.db.Query(fmt.Sprintf("SELECT book_id FROM Book"))
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&remv_id[r.Intn(TestTimes)])
	}

	for i := 0; i < TestTimes; i++ {
		lib.RemoveBook(remv_id[i], "test")
	}
}
