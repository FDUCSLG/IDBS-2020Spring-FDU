package main

import (
	"fmt"
	"testing"
)


func TestCreateTables(t *testing.T) {
	lib := Library{}
	lib.ConnectDB(User,Password)
	err := lib.CreateTables()
	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestCreateTables: \n")
	if err != nil {
		t.Errorf("Error: Failure at TestCreateTables")
	}else{
		fmt.Printf("PASS THIS\n")
	}
}

func TestAddbook(t *testing.T) {
	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestAddBook: \n")
	lib := Library{}
	lib.ConnectDB(User,Password)
	err := lib.AddBook("Gone with the wind","Margaret Mitchell" , "06813","978-1-4165-4894-2")
	if err != nil {
		t.Errorf("Error: Failure at TestAddBook \n")
	}else{
		fmt.Printf("PASS THIS\n")
	}
	fmt.Printf("------------------------------------- \n ")
	errb := lib.AddBook("Gone with the wind","Margaret Mitchell" , "70826","978-1-4165-4894-2")
	if errb != nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestAddBook \n")
	}
}

func TestRemoveBook(t *testing.T) {

	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestRemoveBook: \n")

	lib := Library{}
	lib.ConnectDB(User,Password)

	//input id, reason

	err := lib.RemoveBook("71364","lost")
	if err != nil {
		t.Errorf("Error: Failure at TestRemoveBook \n")
	}else{
		fmt.Printf("PASS THIS \n")
	}
	fmt.Printf("------------------------------------- \n ")
	errb := lib.RemoveBook("58168","lost")
	if errb != nil {
		fmt.Printf("PASS THIS \n")
	}else{
		t.Errorf("Error: Failure at TestRemoveBook \n")
	}

}

func TestAddAccount(t *testing.T) {
	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestAddAccount: \n")

	lib := Library{}
	lib.ConnectDB(User,Password)
	err := lib.AddAccount("hanxy","18300750006","Male","20","CS","qwertyuiop")
	if err == Erraccount("18300750006") {
		fmt.Printf("PASS THIS\n")
	}else if err == nil{
		t.Errorf("Error: Failure at TestAddAccount \n")
	}else{
		t.Errorf("Error: Failure at TestAddAccount \n")
	}
	fmt.Printf("------------------------------------- \n ")
	errb := lib.AddAccount("hatsune","17301257636","Male","18","EE","qwertyuiop")
	if errb == Erraccount("17301257636") {
		t.Errorf("Error: Failure at TestAddAccount \n")
	}else if errb == nil{
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestAddAccount \n")
	}
}

func TestQueryBook(t *testing.T) {
	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestQueryBook: \n")

	lib := Library{}
	lib.ConnectDB(User,Password)

	fmt.Printf("You search the books whose title is Miser \n ")
	err := lib.QueryBook("title","Miser")

	if err == Errbook("Miser")  {
		t.Errorf("Error: Failure at TestQueryBook \n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestQueryBook \n")
	}

	fmt.Printf("------------------------------------- \n ")
	fmt.Printf("You search the books whose author is Charles Dickens \n ")
	err = lib.QueryBook("author","Charles Dickens")

	if err == Errbook("Charles Dickens" )  {
		t.Errorf("Error: Failure at TestQueryBook \n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestQueryBook \n")
	}

	fmt.Printf("------------------------------------- \n ")
	fmt.Printf("You search the books whose ISBN is 978-1-1206-4684-2 \n ")
	err = lib.QueryBook("ISBN","978-1-1206-4684-2")

	if err == Errbook("978-1-1206-4684-2")  {
		t.Errorf("Error: Failure at TestQueryBook \n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestQueryBook \n")
	}

	fmt.Printf("------------------------------------- \n ")
	fmt.Printf("You search the books whose ABCD is 978-1-1206-4684-2 \n ")
	err = lib.QueryBook("ABCD","978-1-1206-4684-2")

	if err == Errbook("978-1-1206-4684-2")  {
		fmt.Printf("PASS THIS\n")
	}else if err == nil {
		t.Errorf("Error: Failure at TestQueryBook \n")
	}else{
		t.Errorf("Error: Failure at TestQueryBook \n")
	}

	fmt.Printf("------------------------------------- \n ")
	fmt.Printf("You search the books whose ISBN is 928-1-1236-4689-2 \n ")
	err = lib.QueryBook("ISBN","928-1-1236-4689-2")

	if err == Errbook("928-1-1236-4689-2")  {
		t.Errorf("Error: Failure at TestQueryBook \n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestQueryBook \n")
	}

}

func TestBorrowBook(t *testing.T) {
	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestBorrowBook: \n")

	lib := Library{}
	lib.ConnectDB("18300750006","qwertyuiop")

	fmt.Printf("Student 18300750006 wants to borrow the book whose id is 71425 \n ")

	err := lib.BorrowBook("18300750006","71425")
	if err == nil{
		fmt.Printf("PASS THIS\n")
	}else if err == Errstate("71425") {
		t.Errorf("Error: Failure at TestBorrowBook \n")
	}else{
		t.Errorf("Error: Failure at TestBorrowBook \n")
	}

	fmt.Printf("------------------------------------- \n ")

	lib.ConnectDB("17301257636","qwertyuiop")
	fmt.Printf("Student 17301257636 wants to borrow the book whose id is 71406 \n ")

	err = lib.BorrowBook("17301257636","71406")
	if err == nil {
		t.Errorf("Error: Failure at TestBorrowBook \n")
	}else if err == Errstate("71406") {
		t.Errorf("Error: Failure at TestBorrowBook \n")
	}else{
		fmt.Printf("PASS THIS\n")
	}

	fmt.Printf("------------------------------------- \n ")

	fmt.Printf("Student 17301257636 wants to borrow the book whose id is 70826 \n ")

	err = lib.BorrowBook("17301257636","70826")
	if err == nil {
		t.Errorf("Error: Failure at TestBorrowBook \n")
	}else if err == Errstate("70826"){
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestBorrowBook \n")
	}
}

func TestBorrowHistory(t *testing.T) {
	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestBorrowHistory: \n")

	lib := Library{}
	lib.ConnectDB("18300750006","qwertyuiop")

	fmt.Printf("Query the student 18300750006's records \n")

	err := lib.BorrowHistory("18300750006")
	if err != nil {
		t.Errorf("Error: Failure at TestBorrowHistory \n")
	}else{
		fmt.Printf("PASS THIS\n")
	}

	fmt.Printf("------------------------------------- \n ")

	lib.ConnectDB("17301257636","qwertyuiop")

	fmt.Printf("Query the student 17301257636's records \n")

	err = lib.BorrowHistory("17301257636")
	if err != nil {
		t.Errorf("Error: Failure at TestBorrowHistory \n")
	}else{
		fmt.Printf("PASS THIS\n")
	}

	fmt.Printf("------------------------------------- \n ")

	lib.ConnectDB("12301257696","qwertyuiop")

	fmt.Printf("Query the student 12301257696's records \n")

	err = lib.BorrowHistory("12301257696")
	if err != nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestBorrowHistory \n")
	}

}

func TestBorrowNow(t *testing.T) {
	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestBorrowNow: \n")

	fmt.Printf("------------------------------------- \n ")
	lib := Library{}
	lib.ConnectDB("18300750006","qwertyuiop")

	fmt.Printf("Query the student 18300750006's records \n")

	err := lib.BorrowNow("18300750006")
	if err == Errhistory("18300750006") {
		t.Errorf("Error: Failure at TestBorrowNow \n")
	}else if err == nil {
		fmt.Printf("PASS THIS \n")
	}else{
		t.Errorf("Error: Failure at TestBorrowNow \n")
	}

	fmt.Printf("------------------------------------- \n ")

	lib.ConnectDB("17301257636","qwertyuiop")

	fmt.Printf("Query the student 17301257636's records \n")

	err = lib.BorrowNow("17301257636")
	if err == Errhistory("17301257636") {
		t.Errorf("Error: Failure at TestBorrowNow \n")
	}else if err == nil {
		fmt.Printf("PASS THIS \n")
	}else{
		t.Errorf("Error: Failure at TestBorrowNow \n")
	}

	fmt.Printf("------------------------------------- \n ")

	lib.ConnectDB("12301257696","qwertyuiop")

	fmt.Printf("Query the student 12301257696's records \n")

	err = lib.BorrowNow("12301257696")
	if err == Errhistory("12301257696") {
		fmt.Printf("PASS THIS \n")
	}else if err == nil {
		t.Errorf("Error: Failure at TestBorrowNow \n")
	}else{
		t.Errorf("Error: Failure at TestBorrowNow \n")
	}
}

func TestCheckDeadline(t *testing.T) {
	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestCheckDeadline: \n")

	fmt.Printf("------------------------------------- \n ")
	lib := Library{}
	lib.ConnectDB("18300750006","qwertyuiop")

	fmt.Printf("Check the deadline of 70826 with the student 18300750006's account \n")

	err := lib.CheckDeadline("18300750006","70826")

	if err == ErrDeadline("70826") {
		t.Errorf("Error: Failure at TestCheckDeadline\n")
	}else if  err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestCheckDeadline\n")
	}

	fmt.Printf("------------------------------------- \n ")

	lib.ConnectDB("17301257636","qwertyuiop")

	fmt.Printf("Check the deadline of 70826 with the student 17301257636's account \n")

	err = lib.CheckDeadline("17301257636","70826")
	if err == ErrDeadline("70826") {
		fmt.Printf("PASS THIS\n")
	}else if  err == nil {
		t.Errorf("Error: Failure at TestCheckDeadline\n")
	}else{
		t.Errorf("Error: Failure at TestCheckDeadline\n")
	}

}

func TestExtendDeadline(t *testing.T) {
	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestExtendDeadline: \n")

	fmt.Printf("------------------------------------- \n ")
	lib := Library{}
	lib.ConnectDB("18300750006","qwertyuiop")

	fmt.Printf("Extend the deadline of 71425 with the student 18300750006's account \n")

	err := lib.ExtendDeadline("18300750006","71425")
	if err == ErrDeadline("71425") {
		t.Errorf("Error: Failure at TestExtendDeadline\n")
	}else if err == ErrExtend("71425") {
		t.Errorf("Error: Failure at TestExtendDeadline\n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestExtendDeadline\n")
	}

	fmt.Printf("------------------------------------- \n ")

	fmt.Printf("Extend the deadline of 70826 with the student 18300750006's account \n")

	err = lib.ExtendDeadline("18300750006","70826")
	if err == ErrDeadline("70826") {
		t.Errorf("Error: Failure at TestExtendDeadline\n")
	}else if err == ErrExtend("70826") {
		fmt.Printf("PASS THIS\n")
	}else if err == nil {
		t.Errorf("Error: Failure at TestExtendDeadline\n")
	}else{
		t.Errorf("Error: Failure at TestExtendDeadline\n")
	}

	fmt.Printf("------------------------------------- \n ")

	lib.ConnectDB("17301257636","qwertyuiop")

	fmt.Printf("Extend the deadline of 71425 with the student 17301257636's account \n")

	err = lib.ExtendDeadline("17301257636","71425")
	if err == ErrDeadline("71425") {
		fmt.Printf("PASS THIS\n")
	}else if err == ErrExtend("71425") {
		t.Errorf("Error: Failure at TestExtendDeadline\n")
	}else if err == nil {
		t.Errorf("Error: Failure at TestExtendDeadline\n")
	}else{
		t.Errorf("Error: Failure at TestExtendDeadline\n")
	}

}

func TestCheckOverdue(t *testing.T) {
	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestCheckOverdue: \n")

	lib := Library{}
	lib.ConnectDB(User,Password)
	fmt.Printf("------------------------------------- \n ")
	fmt.Printf("Check the student 18300750006's overdue records \n")

	err := lib.CheckOverdue("18300750006")
	if err == ErrOverdue("18300750006") {
		t.Errorf("Error: Failure at TestCheckOverdue\n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestCheckOverdue\n")
	}

	fmt.Printf("------------------------------------- \n ")

	lib.ConnectDB(User,Password)

	fmt.Printf("Check the student 17301257636's overdue records \n")

	err = lib.CheckOverdue("17301257636")
	if err == ErrOverdue("17301257636") {
		t.Errorf("Error: Failure at TestCheckOverdue\n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestCheckOverdue\n")
	}

	fmt.Printf("------------------------------------- \n ")

	lib.ConnectDB(User,Password)

	fmt.Printf("Check the student 12301255636's overdue records \n")

	err = lib.CheckOverdue("12301255636")
	if err == ErrOverdue("12301255636") {
		fmt.Printf("PASS THIS\n")
	}else if err == nil{
		t.Errorf("Error: Failure at TestCheckOverdue\n")
	}else{
		t.Errorf("Error: Failure at TestCheckOverdue\n")
	}

}

func TestReturnBook(t *testing.T) {

	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestReturnBook: \n")

	lib := Library{}
	lib.ConnectDB("18300750006","qwertyuiop")
	fmt.Printf("------------------------------------- \n ")
	fmt.Printf("Return the student 18300750006's book 60842\n")

	err := lib.ReturnBook("18300750006","60842")
	if err == ErrReturn("60842") {
		t.Errorf("Error: Failure at TestReturnBook\n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Fail to return \n")
	}

	fmt.Printf("------------------------------------- \n ")

	lib.ConnectDB("17301257636","qwertyuiop")

	fmt.Printf("Return the student 17301257636's book 60842\n")

	err = lib.ReturnBook("17301257636","60842")
	if err == ErrReturn("60842") {
		fmt.Printf("PASS THIS\n")
	}else if err == nil {
		t.Errorf("Error: Failure at TestReturnBook\n")
	}else{
		t.Errorf("Error: Failure at TestReturnBook\n")
	}



}

func TestSuspendAccount(t *testing.T) {

	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestSuspendAccount: \n")

	lib := Library{}
	lib.ConnectDB(User,Password)
	fmt.Printf("------------------------------------- \n ")
	fmt.Printf("Check whether to suspend the student 18300750006's account\n")

	err := lib.SuspendAccount("18300750006")
	if err == ErrSuspend("18300750006") {
		t.Errorf("Error: Failure at TestSuspendAccount\n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestSuspendAccount\n")
	}
	//
	fmt.Printf("------------------------------------- \n ")

	fmt.Printf("Check whether to suspend the student 17301257639's account\n")

	err = lib.SuspendAccount("17301257639")
	if err == ErrSuspend("17301257639") {
		fmt.Printf("PASS THIS\n")
	}else if err == nil {
		t.Errorf("Error: Failure at TestSuspendAccount\n")
	}else{
		t.Errorf("Error: Failure at TestSuspendAccount\n")
	}

	fmt.Printf("------------------------------------- \n ")

	fmt.Printf("Check whether to suspend the student 17307130028's account\n")

	err = lib.SuspendAccount("17307130028")
	if err == ErrSuspend("17307130028") {
		t.Errorf("Error: Failure at TestSuspendAccount\n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestSuspendAccount\n")
	}

	fmt.Printf("------------------------------------- \n ")

}

func TestFreeAccount(t *testing.T) {

	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestFreeAccount: \n")

	lib := Library{}
	lib.ConnectDB(User,Password)
	fmt.Printf("------------------------------------- \n ")
	fmt.Printf("Check whether to free the student 17301257639's account\n")

	err := lib.FreeAccount("17301257639")
	if err == ErrSuspend("17301257639") {
		fmt.Printf("PASS THIS\n")
	}else if err == nil {
		t.Errorf("Error: Failure at TestFreeAccount\n")
	}else{
		t.Errorf("Error: Failure at TestFreeAccount\n")
	}
	//
	fmt.Printf("------------------------------------- \n ")

	fmt.Printf("Check whether to free the student 16301820009's account\n")

	err = lib.FreeAccount("16301820009")
	if err == ErrSuspend("16301820009") {
		t.Errorf("Error: Failure at TestFreeAccount\n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestFreeAccount\n")
	}

	fmt.Printf("------------------------------------- \n ")

	fmt.Printf("Check whether to free the student 17307130028's account\n")

	err = lib.FreeAccount("17307130028")
	if err == ErrSuspend("17307130028") {
		t.Errorf("Error: Failure at TestFreeAccount\n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at TestFreeAccount\n")
	}

	fmt.Printf("------------------------------------- \n ")

}


func TestGetarrears(t *testing.T) {

	fmt.Printf("************************************* \n ")
	fmt.Printf("Now TestGetarrears: \n")

	lib := Library{}
	lib.ConnectDB(User,Password)
	fmt.Printf("------------------------------------- \n ")
	fmt.Printf("Check the money the student 17301257639's should pay\n")

	err := lib.Getarrears("17301257639")

	if err == ErrArrear("17301257639") {
		fmt.Printf("PASS THIS\n")
	}else if err == nil {
		t.Errorf("Error: Failure at Testgetarrears\n")
	}else{
		t.Errorf("Error: Failure at Testgetarrears\n")
	}

	fmt.Printf("------------------------------------- \n ")
	fmt.Printf("Check the money the student 17301257639's should pay\n")

	err = lib.Getarrears("17307130028")

	if err == ErrArrear("17307130028") {
		t.Errorf("Error: Failure at Testgetarrears\n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at Testgetarrears\n")
	}

	fmt.Printf("------------------------------------- \n ")
	fmt.Printf("Check the money the student 18300750006's should pay\n")

	err = lib.Getarrears("18300750006")

	if err == ErrArrear("18300750006") {
		t.Errorf("Error: Failure at Testgetarrears\n")
	}else if err == nil {
		fmt.Printf("PASS THIS\n")
	}else{
		t.Errorf("Error: Failure at Testgetarrears\n")
	}



}