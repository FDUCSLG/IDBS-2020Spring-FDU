package main

import (
	"fmt"
	"testing"
)

func TestCreateTables(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CreateTables()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log("Tables created.")
}

func TestCreateAccount(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CreateTables()
	if err != nil {
		t.Error(err.Error())
	}
	{
		err := lib.CreateAccount("studentid1", "123")
		if err != nil {
			t.Error(err.Error())
		}
	}
	{
		err := lib.CreateAccount("s@#$%^&*", "123")
		if err == nil {
			t.Error("Invalid user name test failed")
		} else {
			t.Log(err.Error())
		}
	}
}

func TestBooktable(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CreateTables()
	if err != nil {
		t.Errorf(err.Error())
	}
	{
		err := lib.AddBook("How to not Go", "some author", "000-0-000-00000-0")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	{
		err := lib.AddBook("How to everything", "authorA", "100-0-000-00000-1")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	{
		err := lib.AddBook("How to not Go", "some author", "000-0-000-00000-0")
		if err == nil {
			t.Error("same book multiple insertion test failed")
		} else {
			t.Log(err.Error())
		}
	}
	{
		err := lib.AddBook("How to Go", "authorB", "-0-000-00000-0")
		if err == nil {
			t.Error("invalid isbn test failed")
		} else {
			t.Log(err.Error())
		}
	}
	{
		q1, err := lib.QueryBooks("author", "authorA")
		if err != nil {
			t.Error(err.Error())
		}
		for _, v := range q1 {
			t.Log(v)
		}
		if len(q1) == 0 {
			t.Error("query 1 fail with no result")
		}
		errr := lib.RemoveBookByISBN("000-0-000-00000-0", "burnt")
		if errr != nil {
			t.Error(errr.Error())
		}
		q2, err2 := lib.QueryBooks("Title", "How to not Go")
		if err2 != nil {
			t.Error(err2.Error())
		}
		for _, v := range q2 {
			t.Log(v)
		}
		if len(q2) != 0 {
			t.Error("query 2 fail book returned after removal")
		}
	}
}

func TestBorrowReturn(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CreateTables()
	if err != nil {
		t.Errorf(err.Error())
	}
	{
		err := lib.AddBook("book1", "author1", "000-0-000-00001-1")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	{
		err := lib.AddBook("book2", "author2", "000-0-000-00001-2")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	{
		err := lib.AddBook("book3", "author3", "000-0-000-00001-3")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	{
		err := lib.AddBook("book4", "author4", "000-0-000-00001-4")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	{
		err := lib.RemoveBookByISBN("000-0-000-00001-4", "remark_removebook4")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	{
		err := lib.BorrowBook("000-0-000-00001-4")
		if err == nil {
			t.Error("borrow test1 fail borrowed missing book")
		} else {
			t.Log(err.Error())
		}
	}
	{
		err := lib.BorrowBook("000-0-000-00001-3")
		if err != nil {
			t.Error(err.Error())
		}
	}
	{
		orig := "anotherguy"
		orig, User = User, orig
		err := lib.BorrowBook("000-0-000-00001-3")
		if err == nil {
			t.Error("borrow test2 fail borrowed others book")
		} else {
			t.Log(err.Error())
		}
		User = orig
	}
	for i := 0; i != 3; i += 1 {
		err := lib.BorrowBook("000-0-000-00001-3")
		if err != nil {
			t.Error(err.Error())
		}
	}
	{
		err := lib.BorrowBook("000-0-000-00001-3")
		if err == nil {
			t.Error("borrow test3 fail borrowed more than 4 times")
		} else {
			t.Log(err.Error())
		}
	}
	{
		books, err := lib.HeldBooks()
		if err != nil || len(books) != 1 {
			t.FailNow()
		}
		for _, v := range books {
			t.Log(v)
		}
	}
	{
		err := lib.ReturnBook("000-0-000-00001-3")
		if err != nil {
			t.Error(err.Error())
		}
	}
	{
		books, err := lib.HeldBooks()
		if err != nil || len(books) != 0 {
			t.FailNow()
		}
	}
	{
		records, err := lib.BorrowHistory()
		if err != nil {
			t.Error(err.Error())
		}
		for _, v := range records {
			fmt.Println(v)
		}
	}
}
