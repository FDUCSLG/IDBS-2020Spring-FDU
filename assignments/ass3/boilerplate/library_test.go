package main

import (
	"testing"
)

func TestCreateTables(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CreateTables()
	if err != nil {
		t.Errorf("can't create tables")
	}
	lib.init()
}

func TestQueryALLUser(t *testing.T){

}
func TestQueryALLBook(t *testing.T){

}
func TestQueryALLBorrowing(t *testing.T){

}
func TestQueryALLBorrowHis(t *testing.T){

}
func TestAddAdm(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.init()
	tables := []struct {
		id string
		password string
		flag int
	}{
		{"1", "a",0},
		{"adm99","111",1},
	}
	for _, table := range tables {
		f, err := lib.AddAdm(table.id, table.password)
		if err != nil {
			t.Errorf("AddAdm Error")
		}
		if f != table.flag{
			t.Errorf("got: %d, want: %d.", f, table.flag)
		}
	}
}
func TestAddStu(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.init()
	tables := []struct {
		id string
		password string
		flag int
	}{
		{"stu01", "123456",0},
		{"student99","111",1},
	}
	for _, table := range tables {
		f, err := lib.AddStu(table.id, table.password)
		if err != nil {
			t.Errorf("AddStu Error")
		}
		if f != table.flag{
			t.Errorf("got: %d, want: %d.", f, table.flag)
		}
	}
}
func TestAddBook(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.init()
	tables := []struct {
		title string
		author string
		isbn string
	}{
		{"title007","author05","isbn0013"},
		{"title020","author12","isbn1202"},
	}
	for _, table := range tables {
		err := lib.AddBook(table.title,table.author,table.isbn)
		if err != nil {
			t.Errorf("AddBook Error")
		}
	}
}
func TestRemoveBook(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.init()
	tables := []struct {
		bookid string
	}{
		{"isbn000101"},
		{"isbn111101"},
	}
	for _, table := range tables {
		err := lib.RemoveBook(table.bookid)
		if err != nil {
			t.Errorf("RemoveBook Error")
		}
	}
}
func TestBorrowBook(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.init()
	tables := []struct {
		studentid string
		isbn string
		flag int
	}{
		{"stu03","isbn0003",1},
		{"stu01","isbn0001",2},
		{"stu01","isbn0006",3},
		{"stu01","isbn0003",0},
	}
	lib.QueryALLBorrowing()
	for _, table := range tables {
		f, err := lib.BorrowBook(table.isbn, table.studentid)
		if err != nil {
			t.Errorf("BorrowBook Error")
		}
		if f != table.flag {
			t.Errorf("got: %d, want: %d.", f, table.flag)

		}
	}
}
func TestQueryBook(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.init()
	tables := []struct {
		bookinfo string
		swit string
		flag int
	}{
		{"title001","0",0},
		{"title055","0",1},
		{"isbn0002","2",0},
		{"author01","1",0},
	}
	for _, table := range tables {
		f, err := lib.QueryBook(table.bookinfo, table.swit)
		if err != nil {
			t.Errorf("QueryBook Error")
		}
		if f != table.flag {
			t.Errorf("got: %d, want: %d.", f, table.flag)
		}
	}
}
func TestQueryHis(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.init()
	tables := []struct {
		studentid string
		flag int
	}{
		{"stu01",0},
		{"stu04",1},
	}
	for _, table := range tables {
		f, err := lib.QueryHis(table.studentid)
		if err != nil {
			t.Errorf("QueryHistory Error")
		}
		if f != table.flag {
			t.Errorf("got: %d, want: %d.", f, table.flag)
		}
	}
}
func TestQueryBorrowing(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.init()
	tables := []struct {
		studentid string
		flag int
	}{
		{"stu01",0},
		{"stu04",1},
	}
	for _, table := range tables {
		f, err := lib.QueryBorrowing(table.studentid)
		if err != nil {
			t.Errorf("QueryBorrowing Error")
		}
		if f != table.flag {
			t.Errorf("got: %d, want: %d.", f, table.flag)
		}
	}
}
func TestCheckDDL(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.init()
	tables := []struct {
		studentid string
		isbn string
		flag int
	}{
		{"stu01","isbn0001",0},
		{"stu04","isbn0002",1},
	}
	for _, table := range tables {
		f, err := lib.CheckDDL(table.studentid, table.isbn)
		if err != nil {
			t.Errorf("CheckDDL Error")
		}
		if f != table.flag {
			t.Errorf("got: %d, want: %d.", f, table.flag)
		}
	}
}
func TestExtendDDL(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.init()
	tables := []struct {
		studentid string
		isbn string
		flag int
	}{
		{"stu01","isbn0001",2},
		{"stu04","isbn0002",1},
		{"stu02","isbn0002",0},
	}
	for _, table := range tables {
		f, err := lib.ExtendDDL(table.studentid, table.isbn)
		if err != nil {
			t.Errorf("ExtendDDL Error")
		}
		if f != table.flag {
			t.Errorf("got: %d, want: %d.", f, table.flag)
		}
	}
}
func TestCheckOver(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.init()
	tables := []struct {
		studentid string
		flag int
	}{
		{"stu01",1},
		{"stu03",0},
	}
	for _, table := range tables {
		f, err := lib.CheckOver(table.studentid)
		if err != nil {
			t.Errorf("CheckOver Error")
		}
		if f != table.flag {
			t.Errorf("got: %d, want: %d.", f, table.flag)
		}
	}
}
func TestRetBook(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.init()
	tables := []struct {
		studentid string
		isbn string
		flag int
	}{
		{"stu01","isbn0001",0},
		{"stu04","isbn0002",1},
	}
	for _, table := range tables {
		f, err := lib.RetBook(table.isbn, table.studentid)
		if err != nil {
			t.Errorf("RetBook Error")
		}
		if f != table.flag {
			t.Errorf("got: %d, want: %d.", f, table.flag)
		}
	}
}
func TestCheckborrow(t *testing.T){

}
