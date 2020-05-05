package main

import (
	"testing"
)

func TestCreateTables(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	err := CreateTables(lib)
	if err != nil {
		t.Errorf("can't create tables")
	}
}

func TestInitializeDB(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	err := InitializeDB(lib, "666666")
	if err != nil {
		t.Errorf("can't initialize database")
	}
}

func TestAccountUpdate(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	err := AccountUpdate(lib)
	if err != nil {
		t.Errorf("can't update account status")
	}
}

func TestAddBook(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	err := AddBook(lib, "test_book_1", "Farland233", "t000000000001")
	if err != nil {
		t.Errorf("can't add books")
	}
}

func TestDeleteBook(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	err := DeleteBook(lib, "t000000000001")
	if err != nil {
		t.Errorf("can't delete books")
	}
}

func TestAddStu(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	err := AddStu(lib, "18307130071", "000000")
	if err != nil {
		t.Errorf("can't add student accounts")
	}
}

func TestSearchBooks(t *testing.T) {
	/*Because this function need to get your input to finish Query, 
	I'm not sure whether it can be tested properly here, so just run library.go
	and log in with accID: admin, pwd: 666666; or accID: 18307130071, pwd: 000000
	(Those two accounts were added to DB using TestInitializeDB and TestAddStu)
	then follow the instrutions to finish function SearchBooks in "Query a book from library" option.*/
}

func TestBorrowbooks(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	err := Borrowbooks(lib, "0000000000001", "18307130071")
	if err != nil {
		t.Errorf("can't borrow books")
	}
}

func TestQuery_Borrow_History(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	err := Query_Borrow_History(lib, "18307130071")
	if err != nil {
		t.Errorf("can't query borrow history")
	}
}

func TestNotReturned(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	err := NotReturned(lib, "18307130071")
	if err != nil {
		t.Errorf("can't query not returned books for an account")
	}
}

func TestCheckDDL(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	err := CheckDDL(lib, "0000000000001")
	if err != nil {
		t.Errorf("can't check DDL for a book")
	}
}

func TestExtend_ddl(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	err := Extend_ddl(lib, "0000000000001", "18307130071")
	if err != nil {
		t.Errorf("can't extend DDL for a book")
	}
}

func TestOverDue(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	err := OverDue(lib, "18307130071")
	if err != nil {
		t.Errorf("can't check overdue books for an account")
	}
}

func TestReturnBooks(t *testing.T) {
	var lib *Library = new(Library)
	ConnectDB(lib)
	var priority int
	priority = 1
	err := ReturnBooks(lib, "0000000000001", "18307130071", &priority)
	if err != nil {
		t.Errorf("can't return books")
	}
}

