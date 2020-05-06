package main

import (
	"fmt"
	//"time"
	// mysql connector
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	StudentID = "18307130112"
	User      = "root"
	Password  = "123456"
	DBName    = "ass3"
	AdminInit = "1"
	AdminInitPassword = "a"
)

type Library struct {
	db *sqlx.DB
}
func (lib *Library) ConnectDB() {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	if err != nil {
		panic(err)
	}
	lib.db = db
}
func mustExecute(db *sqlx.DB, SQLs []string) {
	for _, s := range SQLs {
		_, err := db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
}
// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error {
	mustExecute(lib.db,[]string{"DROP TABLE IF EXISTS Admins","DROP TABLE IF EXISTS Student","DROP TABLE IF EXISTS Book","DROP TABLE IF EXISTS BorrowedBook","DROP TABLE IF EXISTS BorrowHistory"});
	mustExecute(lib.db, []string{
		"CREATE TABLE IF NOT EXISTS Admins (id CHAR(11) NOT NULL,password CHAR(15) NOT NULL, primary key(id));",
		"CREATE TABLE IF NOT EXISTS Student (id CHAR(11), password CHAR(15), borrowright INT, primary key(id));",
		"CREATE TABLE IF NOT EXISTS Book (title CHAR(32), author CHAR(20), isbn CHAR(13), bookid CHAR(15), borrowflag INT, primary key(bookid));",
		"CREATE TABLE IF NOT EXISTS BorrowedBook (isbn CHAR(13), studentid CHAR(11), bookid CHAR(15), rettime DATE, extendtimes INT, primary key(bookid));",
		"CREATE TABLE IF NOT EXISTS BorrowHistory (isbn CHAR(13), studentid CHAR(11), bookid CHAR(15), rettime DATE, brwtime DATE);",
	})
	return nil
}
//test init
func (lib *Library) init() {
	mustExecute(lib.db,[]string{"DELETE FROM Admins","DELETE FROM Student","DELETE FROM Book","DELETE FROM BorrowedBook","DELETE FROM BorrowHistory"});
	_, _ = lib.db.Exec("INSERT INTO Book(isbn, author, title, bookid, borrowflag) " +
		"VALUES(\"isbn0001\", \"author01\", \"title001\", \"isbn000101\", 1), " +
		"(\"isbn0001\", \"author01\", \"title001\", \"isbn000102\", 0), " +
		"(\"isbn0002\", \"author02\", \"title002\", \"isbn000201\", 1), " +
		"(\"isbn0003\", \"author01\", \"title003\", \"isbn000301\", 0), " +
		"(\"isbn0004\", \"author04\", \"title003\", \"isbn000401\", 1), " +
		"(\"isbn0005\", \"author05\", \"title003\", \"isbn000501\", 1), " +
		"(\"isbn0006\", \"author06\", \"title003\", \"isbn000601\", 1)")
	_, _ = lib.db.Exec("INSERT INTO Student(id, password, borrowright) " +
		"VALUES(\"stu01\", \"123456\", 1), (\"stu02\", \"147258\", 1), (\"stu03\", \"147258\", 0), (\"stu04\", \"147258\", 1)")
	_, _ = lib.db.Exec(fmt.Sprintf("INSERT INTO Admins(id, password) " +
		"VALUES(\"12345678\", \"asdfghjk\"), (\"%s\", \"%s\")", AdminInit, AdminInitPassword),)
	_, _ = lib.db.Exec("INSERT INTO BorrowedBook(isbn, studentid, bookid, rettime, extendtimes) " +
		"VALUES(\"isbn0001\", \"stu01\", \"isbn000101\", '2020-06-01', 3), " +
		"(\"isbn0002\", \"stu02\", \"isbn000201\", '2020-06-01', 0)," +
		" (\"isbn0004\", \"stu03\", \"isbn000401\", '2020-05-01', 0), " +
		"(\"isbn0005\", \"stu03\", \"isbn000501\", '2020-05-01', 0), " +
		"(\"isbn0006\", \"stu03\", \"isbn000601\", '2020-05-01', 0)")
	_, _ = lib.db.Exec("INSERT INTO BorrowHistory(isbn, studentid, bookid, rettime, brwtime) " +
		"VALUES(\"isbn0001\", \"stu01\", \"isbn000101\", '2020-06-01', '2020-02-01')")
	_, _ = lib.db.Exec("INSERT INTO BorrowHistory(isbn, studentid, bookid, rettime, brwtime) " +
		"VALUES(\"isbn0002\", \"stu02\", \"isbn000201\", '2020-06-01', '2020-05-01')")
	_, _ = lib.db.Exec("INSERT INTO BorrowHistory(isbn, studentid, bookid, rettime, brwtime) " +
		"VALUES(\"isbn0002\", \"stu01\", \"isbn000201\", '2020-04-01', '2020-03-01')")
	_, _ = lib.db.Exec("INSERT INTO BorrowHistory(isbn, studentid, bookid, rettime, brwtime) " +
		"VALUE(\"isbn0004\", \"stu03\", \"isbn000401\", '2020-05-01', '2020-04-01'), " +
		"(\"isbn0005\", \"stu03\", \"isbn000501\", '2020-05-01', '2020-04-01'), " +
		"(\"isbn0006\", \"stu03\", \"isbn000601\", '2020-05-01', '2020-04-01')")
}
//Query all the student and administrator
func (lib *Library) QueryALLUser() error {
	rows, err := lib.db.Query("Select * From Admins order by id asc")
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	count := 0
	for rows.Next(){
		var id, pass string
		rows.Scan(&id, &pass)
		count += 1
		fmt.Println(fmt.Sprintf("Administrator %d: %s",count, id))
	}
	rows, err = lib.db.Query("Select * From Student order by id asc")
	if err != nil {
		panic(err)
	}
	count = 0
	for rows.Next(){
		var id, pass string
		var flag int
		rows.Scan(&id, &pass, &flag)
		count += 1
		stat := "can"
		if flag == 0{
			stat = "cannot"
		}
		fmt.Println(fmt.Sprintf("Student %d: %s, he/she %s borrow books.",count, id, stat))
	}
	return nil
}
//Query all the books
func (lib *Library) QueryALLBook() error {
	rows, err := lib.db.Query("Select * From Book order by bookid asc")
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next(){
		var tit, aut, isb, bid string
		var bfg int
		rows.Scan(&tit, &aut, &isb, &bid, &bfg)
		stat := "is"
		if bfg == 0{
			stat = "isn`t"
		}
		fmt.Println(fmt.Sprintf("title: %s, author: %s, ISBN: %s, bookid: %s, the book %s borrowed.",tit, aut, isb, bid, stat))
	}
	return nil
}
//Query all the borrowing books
func (lib *Library) QueryALLBorrowing() error {
	rows, err := lib.db.Query("Select * From BorrowedBook order by bookid asc")
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next(){
		var isb, sid, bid, d1 string
		var ett int
		rows.Scan(&isb, &sid, &bid, &d1, &ett)
		fmt.Println(fmt.Sprintf("ISBN: %s, studentid: %s, bookid: %s, extending times: %d, expected return date: %s", isb, sid, bid, ett, d1))
	}
	return nil
}
//Query all the borrow history
func (lib *Library) QueryALLBorrowHis() error {
	rows, err := lib.db.Query("Select * From BorrowHistory ORDER BY BOOKID ASC")
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next(){
		var isb, sid, bid, d1, d2 string
		rows.Scan(&isb, &sid, &bid, &d1, &d2)
		fmt.Println(fmt.Sprintf("ISBN: %s, studentid: %s, bookid: %s, borrow date: %s, expected return/returned date: %s", isb, sid, bid, d2, d1))
	}
	return nil
}
//check admin account
func (lib *Library) CheckAdmin(id, password string) bool {
	rows, err := lib.db.Query("SELECT password FROM Admins WHERE id = ?", id)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var k string
		rows.Scan(&k)
		if password == k {
			return true
		}
	}
	return false
}
//check student account
func (lib *Library) CheckStudent(id, password string) bool{
	rows, err := lib.db.Query("SELECT password FROM Student WHERE id = ?", id)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var k string
		rows.Scan(&k)
		if password == k {
			return true
		}
	}
	return false
}
//add administation user
func (lib *Library) AddAdm(id, password string) (int, error) {
	rows, err := lib.db.Query("select * from Admins where id = ?", id)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	if rows.Next(){
		fmt.Println("This administrator id has already had its account.")
		return 0, nil
	}
	_, _ = lib.db.Exec(fmt.Sprintf("INSERT INTO Admins(id, password) value (\"%s\",\"%s\")", id, password))
	return 1, nil
}
//add student account by the administrator's account
func (lib *Library) AddStu(id, password string) (int, error) {
	rows, err := lib.db.Query("select * from Student where id = ?", id)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	if rows.Next(){
		fmt.Println("This student id has already had its account.")
		return 0, nil
	}

	_, _ = lib.db.Exec(fmt.Sprintf("INSERT INTO Student(id, password, borrowright) value (\"%s\", \"%s\", \"%d\")", id, password, 1))
	return 1, nil
}
// AddBook add a book into the library
func (lib *Library) AddBook(title, author, ISBN string) error {
	rows, err := lib.db.Query("SELECT * From Book where isbn = ?", ISBN)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	count := 1
	for rows.Next(){
		count += 1
	}
	st := ""
	if count <= 10{
		st = fmt.Sprintf("%s0%d", ISBN, count)
	}else{
		st = fmt.Sprintf("%s%d", ISBN, count)
	}
	_, _ = lib.db.Exec("INSERT INTO BOOK(title, author, isbn, bookid, borrowflag) value (?,?,?,?,?)", title, author, ISBN, st, 0)
	return nil
}
//remove a book from the library with explanation (e.g. book is lost)
func (lib *Library) RemoveBook(bookid string) error {
	_, _ = lib.db.Exec("DELETE From Book where bookid = ?", bookid)
	_, _ = lib.db.Exec("DELETE From BorrowedBook where bookid = ?", bookid)
	return nil
}
//suspend student's account if the student has more than 3 overdue books
func (lib *Library) Checkborrow(studentid string) bool {
	rows, err := lib.db.Query("Select * From BorrowedBook where studentid = ? and rettime < CURRENT_DATE()", studentid)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	count := 0
	for rows.Next(){
		count += 1
	}
	if count < 3{
		return true
	}
	_, _ = lib.db.Exec("Update Student set borrowright = 0 where studentid = ?", studentid)
	return false
}
//borrow a book from the library with a student account
func (lib *Library) BorrowBook(isbn, studentid string) (int, error) {
	stat := lib.Checkborrow(studentid)
	if stat == false{
		fmt.Println("You need to return the overdue books before borrowing new one!")
		return 1, nil
	}
	rows, err := lib.db.Query("Select * From BorrowedBook where studentid = ?", studentid)
	if err != nil {
		panic(err)
	}
	for rows.Next(){
		var isb, sid, bid, rtt string
		var ett int
		rows.Scan(&isb, &sid, &bid, &rtt, &ett)
		if isb == isbn {
			fmt.Printf("You have already borrowed this book!\n")
			return 2, nil
		}
	}
	rows, err = lib.db.Query("Select * From Book where isbn = ?", isbn)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next(){
		var ttl, aut, isb, bid string
		var bfg int
		rows.Scan(&ttl, &aut, &isb, &bid, &bfg)
		if bfg == 1 {
			continue
		}else {
			_, _ = lib.db.Exec("INSERT INTO BorrowedBook(isbn, studentid, bookid, rettime, extendtimes) VALUES(?, ?, ?, date_add(CURRENT_DATE(), interval 1 month), 0)", isb, studentid, bid)
			_, _ = lib.db.Exec("INSERT INTO BorrowHistory(isbn, studentid, bookid, rettime, brwtime) VALUES(?, ?, ?, date_add(CURRENT_DATE(), interval 1 month), CURRENT_DATE())", isb, studentid, bid)
			_, _ = lib.db.Exec("update Book set borrowflag = 1 where bookid = ?", bid)
			fmt.Printf("You have Succeeded Borrowing this book!\n")
			return 0, nil
		}
	}
	fmt.Printf("Fail Borrowing! Book doesn`t exist or all the same books have been borrowed.\n")
	return 3, nil
}
//query books by title, author or ISBN
func (lib *Library) QueryBook(bookinfo, swit string) (int, error) {
	rows, err := lib.db.Query("Select * From Book where isbn = ?", bookinfo)
	if swit == "0" {
		rows, err = lib.db.Query("Select * From Book where title = ?", bookinfo)
	} else if swit == "1" {
		rows, err = lib.db.Query("Select * From Book where author = ?", bookinfo)
	}
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	count := 0
	for rows.Next() {
		var ttl, aut, isb, bid string
		var bfg int
		rows.Scan(&ttl, &aut, &isb, &bid, &bfg)
		state := "is"
		if bfg == 0 {
			state = "isn`t"
		}
		fmt.Println(fmt.Sprintf("title: %s, author: %s, ISBN: %s, bookid: %s, the book %s borrowed.", ttl, aut, isb, bid, state))
		count += 1
	}
	if count != 0{
		fmt.Println(fmt.Sprintf("The number of valid infomation is %d.", count))
		return 0, nil
	}
	fmt.Println("The book you search is inexistent.")
	return 1, nil
}
//query the borrow history of a student account
func (lib *Library) QueryHis(studentid string) (int, error) {
	rows, err := lib.db.Query("Select * From BorrowHistory where studentid = ?", studentid)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	count := 0
	for rows.Next(){
		count += 1
		var isb, sid, bid, d1, d2 string
		rows.Scan(&isb, &sid, &bid, &d1, &d2)
		fmt.Println(fmt.Sprintf("ISBN: %s, studentid: %s, bookid: %s, borrow date: %s, expected return/returned date: %s", isb, sid, bid, d2, d1))
	}
	if count == 0{
		fmt.Println(fmt.Sprintf("The student %s hasn`t borrowed any book yet.", studentid))
		return 1, nil
	}else{
		fmt.Println("All information is as above.")
		return 0, nil
	}
}
//query the books a student has borrowed and not returned yet
func (lib *Library) QueryBorrowing(studentid string) (int, error) {
	rows, err := lib.db.Query("Select * From BorrowedBook where studentid = ?", studentid)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	count := 0
	for rows.Next(){
		count += 1
		var isb, sid, bid, d1 string
		var ett int
		rows.Scan(&isb, &sid, &bid, &d1, &ett)
		fmt.Println(fmt.Sprintf("ISBN: %s, studentid: %s, bookid: %s, extending times: %d, expected return date: %s", isb, sid, bid, ett, d1))
	}
	if count == 0{
		fmt.Println(fmt.Sprintf("The student %s isn`t borrowing any book now.", studentid))
		return 1, nil
	}else{
		fmt.Println("All information is as above.")
		return 0, nil
	}
}
//check the deadline of returning a borrowed book
func (lib *Library) CheckDDL(studentid, isbn string) (int, error) {
	rows, err := lib.db.Query("Select * From BorrowedBook where studentid = ?", studentid)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next(){
		var isb, sid, bid, d1 string
		var ett int
		rows.Scan(&isb, &sid, &bid, &d1, &ett)
		if isb == isbn {
			fmt.Printf(fmt.Sprintf("Deadline of returning book(ISBN code: %s) is %s.\n", isb, d1))
			return 0, nil
		}
	}
	fmt.Printf(fmt.Sprintf("You didn`t borrow this book.\n"))
	return 1, nil
}
//extend the deadline of returning a book, at most 3 times
func (lib *Library) ExtendDDL(studentid, isbn string) (int, error) {
	rows, err := lib.db.Query("Select * From BorrowedBook where studentid = ?", studentid)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	var isb, sid, bid, d1 string
	var ett int
	for rows.Next() {
		rows.Scan(&isb, &sid, &bid, &d1, &ett)
		if isb == isbn {
			break
		}
	}
	if isb != isbn {
		fmt.Printf("You didn`t borrow this book!\n")
		return 1, nil
	}else if ett >= 3{
		fmt.Printf("You have extend this book for 3 times!\n")
		return 2, nil
	}else{
		_, _ = lib.db.Exec("Update BorrowedBook SET rettime = date_add(rettime, interval 1 month), extendtimes += 1 where studentid = ? and isbn = ?",studentid, isbn)
		_, _ = lib.db.Exec("Update BorrowedHistory SET rettime = rettime = date_add(rettime, interval 1 month) where studentid = ? and isbn = ?",studentid, isbn)
		fmt.Printf("Deadline of returning book is extended.\n")
		return 0, nil
	}
}
//check if a student has any overdue books that needs to be returned
func (lib *Library) CheckOver(studentid string) (int, error) {
	rows, err := lib.db.Query("Select * From BorrowedBook where studentid = ? and rettime < CURRENT_DATE()", studentid)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	count := 0
	for rows.Next(){
		var isb, sid, bid, d1 string
		var ett int
		rows.Scan(&isb, &sid, &bid, &d1, &ett)
		fmt.Printf("Book which ISBN = %s is overdue.\n", isb)
		count += 1
	}
	if count == 0{
		fmt.Printf("The student %s doesn`t have overdue books.\n", studentid)
		return 1, nil
	}
	return 0, nil
}
//return a book to the library by a student account
func (lib *Library) RetBook(isbn, studentid string) (int, error) {
	rows, err := lib.db.Query("Select * From BorrowedBook where studentid = ?", studentid)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next(){
		var isb, sid, bid, d1 string
		var ett int
		rows.Scan(&isb, &sid, &bid, &d1, &ett)
		if isb != isbn {
			continue
		}
		_, _ = lib.db.Exec("delete from BorrowedBook where  bookid = ?", bid)
		_, _ = lib.db.Exec("update Book set borrowflag = 0 where bookid = ?", bid)
		_, _ = lib.db.Exec("update BorrowHistory set rettime = CURRENT_DATE() where bookid = ?", bid)
		fmt.Printf("Succeed Returning\n")
		return 0, nil
	}
	fmt.Printf("The student did not borrow this book!\n")
	return 1, nil
}
func main() {
	fmt.Println("Welcome to the Library Management System!")
	var lib Library
	lib.ConnectDB()
	lib.CreateTables()
	lib.init()
	var initid, initpass string
	Menu:= []string{"1: Add Administrator (Admin Only)",
					"2: Add Student (Admin Only)",
					"3: Add Book (Admin Only)",
					"4: Remove Book (Admin Only)",
					"5: Borrow Book (Student Only)",
					"6: Query Book",
					"7: Query Student`s Borrowing History",
					"8: Query Student`s Borrowing Books",
					"9: Check the Deadline of Returing Book (Student Only)",
					"10: Extend the Deadline of Returing Book (Student Only)",
					"11: Query If a Student Has Overdue Books",
					"12: Return Book (Student Only)",
					"13: Check and Suspend Student's Account (Admin Only)",
					"14: Check all the administrators and students (Admin Only)",
					"15: Check all the book infomation (Admin Only)",
					"16: Check all the borrowing book infomation (Admin Only)",
					"17: Check all the borrow history (Admin Only)",
					"0: Exit Account"}
	fmt.Println("Please enter the activation account:")
	fmt.Scanln(&initid)
	fmt.Println("Please enter the activation password:")
	fmt.Scanln(&initpass)
	if initid != AdminInit || initpass != AdminInitPassword{
		fmt.Println("The account or password is invalid, please restart the program.")
	} else {
		for true {
			var modetest, mode string
			flag := false
			fmt.Println("Please enter the usage mode:\n 0: Student mode  1: Administrator mode  2: Exit")
			fmt.Scanln(&modetest)
			if modetest == "2" {
				fmt.Println("Thank you for using the Library Management System of Fairy Union of Defense and Attack Nebula University!")
				break
			}else if modetest == "1"{
				fmt.Println("Please enter the account:")
				fmt.Scanln(&initid)
				fmt.Println("Please enter the password:")
				fmt.Scanln(&initpass)
				if lib.CheckAdmin(initid, initpass) == false{
					fmt.Println("The account or password is invalid.")
				}else{
					flag = true
				}
			}else if modetest == "0"{
				fmt.Println("Please enter the account:")
				fmt.Scanln(&initid)
				fmt.Println("Please enter the password:")
				fmt.Scanln(&initpass)
				if lib.CheckStudent(initid, initpass) == false {
					fmt.Println("The account or password is invalid.")
				}else{
					flag = true
				}
			}else{
				fmt.Println("Mode is inexistent.")
			}
			for flag {
				for _, putstring := range Menu {
					fmt.Println(putstring)
				}
				fmt.Println("Please enter the function mode:")
				fmt.Scanln(&mode)
				if mode == "0" {
					break
				}
				switch mode {
				case "1":
					if modetest == "0" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						for true {
							var id, pass, repass string
							fmt.Println("Please enter the account number (11 digits):")
							fmt.Scanln(&id)
							fmt.Println("Please enter the password (no more than 15 characters):")
							fmt.Scanln(&pass)
							fmt.Println("Please enter the password again:")
							fmt.Scanln(&repass)
							if pass != repass {
								fmt.Println("The password entered twice is different, please re-enter:")
							} else {
								lib.AddAdm(id, pass)
								fmt.Println("The account is created!")
								break
							}
						}
					}
				case "2":
					if modetest == "0" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						for true {
							var id, pass, repass string
							fmt.Println("Please enter the account number (11 digits):")
							fmt.Scanln(&id)
							fmt.Println("Please enter the password (no more than 15 characters):")
							fmt.Scanln(&pass)
							fmt.Println("Please enter the password again:")
							fmt.Scanln(&repass)
							if pass != repass {
								fmt.Println("The password entered twice is different, please re-enter:")
							} else {
								lib.AddStu(id, pass)
								fmt.Println("The account is created!")
								break
							}
						}
					}
				case "3":
					if modetest == "0" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						var title, author, isbn string
						fmt.Println("Please enter the book title, author, and isbn code, separated by newlines:")
						fmt.Scanln(&title)
						fmt.Scanln(&author)
						fmt.Scanln(&isbn)
						lib.AddBook(title, author, isbn)
					}
				case "4":
					if modetest == "0" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						var bookid string
						fmt.Println("Please enter the bookid code:")
						fmt.Scanln(&bookid)
						lib.RemoveBook(bookid)
					}
				case "5":
					if modetest == "1" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						var bookinfo string
						fmt.Println("Please enter ISBN code of the book that you want:")
						fmt.Scanln(&bookinfo)
						lib.BorrowBook(bookinfo, initid)
					}
				case "6":
					var searchcode, bookinfo string
					fmt.Println("You can search the infomation you want by title(0), author(1) or ISBN code(2):")
					for true {
						fmt.Scanln(&searchcode)
						if searchcode == "0" || searchcode == "1" || searchcode == "2"{
							break
						}
						fmt.Println("Mode is inexistent.")
					}
					fmt.Println("Please enter the book information that matches your options:")
					fmt.Scanln(&bookinfo)
					lib.QueryBook(bookinfo, searchcode)
				case "7":
					var stuid string
					if modetest == "1"{
						fmt.Println("Please enter the student id :")
						fmt.Scanln(&stuid)
					}else{
						stuid = initid
					}
					lib.QueryHis(stuid)
				case "8":
					var stuid string
					if modetest == "1"{
						fmt.Println("Please enter the student id :")
						fmt.Scanln(&stuid)
					}else{
						stuid = initid
					}
					lib.QueryBorrowing(stuid)
				case "9":
					if modetest == "1" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						var bookinfo string
						fmt.Println("Please enter ISBN code of the book that you want to check:")
						fmt.Scanln(&bookinfo)
						lib.CheckDDL(initid, bookinfo)
					}
				case "10":
					if modetest == "1" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						var bookinfo string
						fmt.Println("Please enter ISBN code of the book that you want to extend return deadline:")
						fmt.Scanln(&bookinfo)
						lib.ExtendDDL(initid, bookinfo)
					}
				case "11":
					var stuid string
					if modetest == "1"{
						fmt.Println("Please enter the student id :")
						fmt.Scanln(&stuid)
					}else{
						stuid = initid
					}
					lib.CheckOver(stuid)
				case "12":
					if modetest == "1" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						var bookinfo string
						fmt.Println("Please enter ISBN code of the book that you want to return:")
						fmt.Scanln(&bookinfo)
						lib.RetBook(bookinfo, initid)
					}
				case "13":
					if modetest == "0" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						var studentid string
						fmt.Println("Please enter the student id:")
						fmt.Scanln(&studentid)
						borrowflag := lib.Checkborrow(studentid)
						if borrowflag == false{
							fmt.Println("This student has over 3 overdue books, and his/her account is suspended.")
						}else{
							fmt.Println("This student has less than 3 overdue books, and his/her account is normal.")
						}
					}
				case "14":
					if modetest == "0" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						lib.QueryALLUser()
					}
				case "15":
					if modetest == "0" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						lib.QueryALLBook()
					}
				case "16":
					if modetest == "0" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						lib.QueryALLBorrowing()
					}
				case "17":
					if modetest == "0" {
						fmt.Println("You do not have permission to execute this command.")
					} else {
						lib.QueryALLBorrowHis()
					}
				default:
					fmt.Println("Mode is inexistent.")
				}
			}
		}
	}
}