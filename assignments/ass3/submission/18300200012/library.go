package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	// mysql connector
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

//const
const (
	User     = "root"
	Password = "llq0215llq"
	DBName   = "Library"
	/*Library Constant*/
	VT             bool = true
	BorrowPeriod   int  = 40
	KeepDay        int  = 10
	MaxExtendDDL   int  = 3
	StuAccLen      int  = 10
	LibAccLen      int  = 5
	PasswordMaxlen int  = 16
	PasswordMinlen int  = 4
	MaxOverdue     int  = 3
	EnPWLen        int  = 60
	BookIdLen      int  = 17
	ISBNLen        int  = 13
	MaxBorBookNum  int  = 15
	/*Account authority*/
	Administrator int = 100
	Librarian     int = 101
	Student       int = 102
	/*Student & Librarian Status tag*/
	InvalidAcc   int = 200 //InvalidAcc < InitialState is a must
	InitialState int = 201
	/*Book Status*/
	OnShelf   int = 301 /*min(book_status) == OnShelf && max(book_status) == NonLoan  is a must*/
	OnLoan    int = 302
	OnReserve int = 303
	NonLoan   int = 305
	/*Record Type*/
	Borrow        int = 306
	Reserve       int = 307
	Return        int = 308
	CancelReserve int = 309
	/*ManageLog Type*/
	AddBook       int = 311
	RemoveBook    int = 312
	AddStu        int = 313
	RemoveStu     int = 314
	SuspendStu    int = 315
	DissuspendStu int = 316
)

var CUR_ACCOUNT string
var Authority int
var Days int

//Library library
type Library struct {
	db *sqlx.DB
}

//TESTDATA  test
func (lib *Library) TESTDATA() error {
	fmt.Println("Add test date...")
	lib.AddLib("lib01", "1111")
	lib.AddLib("lib02", "1111")

	/*4 students*/
	lib.AddStu("stu0000001", "2222")
	lib.AddStu("stu0000002", "2222")
	lib.AddStu("stu0000003", "2222")
	lib.AddStu("stu0000004", "2222")

	/*add books*/
	lib.AddBook("9870000000002", "English", "Jhon")
	lib.AddBook("9870000000003", "Chinese", "Lisa")
	lib.AddBook("9870000000004", "Physics", "Frank")
	lib.AddBook("9870000000005", "Biology", "Sam")
	lib.AddBook("9870000000006", "Geography", "Aaron")
	lib.AddBook("9870000000007", "Economics", "Frank")
	lib.AddBook("9870000000008", "Computer", "Barton")
	lib.AddBook("9870000000009", "History", "Clark")
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := 0; i < 1000; i++ {
		lib.AddBook(fmt.Sprintf("9870000000%.3d", 1+r.Intn(999)), "xxx", "xxx")
	}

	/*Randomly Student Activities*/

	for i := 0; i < 1000; i++ {
		todo := r.Intn(8)
		/*P(borrow) : P(return): P (exddl) = 5 : 2 : 1 */
		if todo < 5 { // borrow
			lib.BorrowBook(fmt.Sprintf("stu000000%d", 1+r.Intn(4)), fmt.Sprintf("9870000000%.3dn001", 1+r.Intn(999)))
		} else if todo < 7 { //return
			lib.ReturnBook(fmt.Sprintf("stu000000%d", 1+r.Intn(4)), fmt.Sprintf("9870000000%.3dn001", 1+r.Intn(999)))
		} else if todo < 8 {
			book_id := fmt.Sprintf("9870000000%.3dn001", 1+r.Intn(999))
			var stu_id string
			err := lib.db.QueryRow(fmt.Sprintf(
				"SELECT count stu_id "+
					"FROM Record "+
					"WHERE rec_id IN (SELECT max(rec_id) FROM Record WHERE type = %v AND book_id = '%s') AND "+
					"1 = (SELECT count(*) FROM Book WHERE book_id = '%s' AND (book_status = %v OR book_status = %v)) ", Borrow, book_id, book_id, OnLoan, OnReserve)).Scan(&stu_id)
			if err == nil {
				lib.ExtendDDL(stu_id, book_id)
			}
		}
	}
	fmt.Println("---------------------------------------------------------------------------------------------------------")
	return nil

}

//VirtualTime  test
func VirtualTime() string {

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	Days = Days + r.Intn(2) //add 0 or 1
	aDay, _ := time.ParseDuration("24h")
	virtual_today := time.Now()
	for i := 0; i < Days && VT; i++ {
		virtual_today = virtual_today.Add(aDay)
	}
	return virtual_today.Format("2006-01-02")
}

func sqlExecute(db *sqlx.DB, SQLs []string) {
	for _, s := range SQLs {
		_, err := db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
}

//ConnectDB connect to the database
func (lib *Library) ConnectDB() {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	if err != nil {
		panic(err)
	}
	sqlExecute(db, []string{
		fmt.Sprintf("DROP DATABASE IF EXISTS %s", DBName),
		fmt.Sprintf("CREATE DATABASE %s", DBName),
		fmt.Sprintf("USE %s", DBName),
	})
	lib.db = db

}

//CmdHelp  command line interpretation
func CmdHelp() error {

	admin := []string{
		fmt.Sprintf("al\t[lib_id]\t[password]\tAdd a librarian account,len(lib_id) = %d, %d<len(password)<%d", LibAccLen, PasswordMinlen, PasswordMinlen),
		fmt.Sprintf("as\t[stu_id]\t[password]\tAdd a student account,len(stu_id) = %d, %d<len(password)<%d", StuAccLen, PasswordMinlen, PasswordMinlen),
		fmt.Sprintf("rl\t[lib_id]\t\t\tRemove a librarian account,len(lib_id) = %d", LibAccLen),
		fmt.Sprintf("rs\t[stu_id]\t\t\tRemove a student account,len(stu_id) = %d", StuAccLen),
		"cpl\t[lib_id]\t\t\tChange the password of Librarian",
		"cps\t[stu_id]\t\t\tChange the password of Student",
	}
	lib := []string{
		"ab\t[ISBN]\t[title]\t[author]\tAdd a book.",
		"rb\t[book_id]\t[notes]\t\tRemovedd a book.",
		"hs\t[stu_id]\t[tag]\t\tQuery a student's borrow history.",
		"hb\t[book_id]\t\t\tQuery a book's borrow history.",
		"bi\t[stu_id]\t\t\tDisplay BorrowInfo.",
		"ss\t[stu_id]\t\t\tSuspend student account.",
		"cp\t\t\t\t\tChange the password.",
	}
	stu := []string{
		"hs\t\t[tag]\t\tQuery borrow history.",
		"bi\t\t\t\tDisplay BorrowInfo.",
		"cp\t\t\t\t\tChange the password.",
	}
	help := []string{
		"qu\t[kw1]\t[kw2]\t[kw3]\t\tQuery a book,at least one keyword.",
		"qui\t[ISBN]\t\t\t\tQuery a book by ISBN.",
		"qua\t[author]\t\t\tQuery a book by author.",
		"qut\t[title]\t\t\t\tQuery a book by title.",

		"h\t\t\t\t\tDisplay this help.",
		"q\t\t\t\t\tQuit the Program.",
		"lo\t\t\t\t\tLog out.",
	}

	if Authority == Administrator {
		for _, val := range admin {
			fmt.Println(val)
		}
	} else if Authority == Librarian {
		for _, val := range lib {
			fmt.Println(val)
		}
	} else if Authority == Student {
		for _, val := range stu {
			fmt.Println(val)
		}
	}
	for _, val := range help {
		fmt.Println(val)
	}
	return nil
}

//Initialization create database, tables
func Initialization() (fdu Library) {
	fdu.ConnectDB()
	fdu.CreateTables()
	fdu.TESTDATA()
	fmt.Println("Welcome to the Library Management System!For Help, Type 'h'\nInput your account and password")
	return fdu
}

//LogIn login account.
func (lib *Library) LogIn() bool {
	done := true
	var account, password string
	for done {
		fmt.Printf("Account: ")
		_, err := fmt.Scanf("%s", &account)
		if err != nil {
			panic(err)
		}
		if account == "admin" { //Administrators
			Authority = Administrator
			break
		}
		fmt.Printf("Password: ")
		_, err = fmt.Scanf("%s", &password)
		if err != nil {
			panic(err)
		}
		var acc_status int
		var enpw string
		if len(account) == LibAccLen { //Librarian
			_ = lib.db.QueryRow(fmt.Sprintf("SELECT lib_status, password FROM Librarian WHERE lib_id = '%s'", account)).Scan(&acc_status, &enpw)
			err = bcrypt.CompareHashAndPassword([]byte(enpw), []byte(password))
			if err == nil {
				if acc_status != InvalidAcc {
					fmt.Printf("Login successful! Account: %s, Type: Librarian\n", account)
					Authority = Librarian
					break
				} else if acc_status == InvalidAcc {
					fmt.Println("This account has been suspended.")
				}
			}
		} else if len(account) == StuAccLen { //Student
			_ = lib.db.QueryRow(fmt.Sprintf("SELECT stu_status,password FROM Student WHERE stu_id = '%s'", account)).Scan(&acc_status, &enpw)
			err = bcrypt.CompareHashAndPassword([]byte(enpw), []byte(password))
			if err == nil {
				if acc_status != InvalidAcc {
					fmt.Printf("Login successful! Account: %s, Type: Student\n", account)
					Authority = Student
					break
				} else if acc_status == InvalidAcc {
					fmt.Println("This account has been suspended.")
				}
			}

		}
		fmt.Printf("Account login failed! Try again?[Y/N]")
		var c string
		_, err = fmt.Scanf("%s", &c)
		if err != nil {
			panic(err)
		}
		if c != "y" && c != "Y" {
			fmt.Println("Bye")
			done = false
		}

	}
	CUR_ACCOUNT = account
	return done
}

//RUN command window
func RUN(fdu Library) {
	done := fdu.LogIn()
	in := bufio.NewReader(os.Stdin)
	for done {
		fmt.Printf("\033[1;33m%s@Library->\033[0m ", CUR_ACCOUNT)
		cmd, oper1, oper2, oper3 := "null", "null", "null", "null"
		s, err := in.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Sscanln(s, &cmd, &oper1, &oper2, &oper3)
		switch cmd {
		case "al":
			if Authority == Administrator {
				fdu.AddLib(oper1, oper2)
			} else {
				fmt.Println("Permission Denied!")
			}
		case "as":
			if Authority == Administrator {
				fdu.AddStu(oper1, oper2)
			} else {
				fmt.Println("Permission Denied!")
			}
		case "rl":
			if Authority == Administrator {
				fdu.RemoveLib(oper1)
			} else {
				fmt.Println("Permission Denied!")
			}
		case "rs":
			if Authority == Administrator {
				fdu.RemoveStu(oper1)
			} else {
				fmt.Println("Permission Denied!")
			}
		case "h":
			CmdHelp()
		case "lo":
			RUN(fdu)
			done = false
		case "q":
			done = false
			fmt.Println("Bye")
		case "ab":
			if Authority == Librarian {
				fdu.AddBook(oper1, oper2, oper3)
			} else {
				fmt.Println("Permission Denied!")
			}
		case "rb":
			if Authority == Librarian {
				fdu.RemoveBook(oper1, oper2)
			} else {
				fmt.Println("Permission Denied!")
			}
		case "hs":
			if Authority == Librarian {
				fdu.HistoryStu(oper1, oper2)
			} else if Authority == Student {
				fdu.HistoryStu(CUR_ACCOUNT, oper1)
			} else {
				fmt.Println("Permission Denied!")
			}
		case "hb":
			if Authority == Librarian {
				fdu.HistoryBook(oper1)
			} else {
				fmt.Println("Permission Denied!")
			}
		case "ss":
			if Authority == Librarian {
				fdu.SuspendStuAcc(CUR_ACCOUNT, oper1, oper2)
			} else {
				fmt.Println("Permission Denied!")
			}
		case "qu":
			fdu.QueryBook(oper1, oper2, oper3)
		case "qui":
			fdu.QueryBook(oper1, "$", "ISBN")
		case "qua":
			fdu.QueryBook(oper1, "$", "author")
		case "qut":
			fdu.QueryBook(oper1, "$", "title")
		case "bi":
			if Authority == Student {
				fdu.BorrowInfo(CUR_ACCOUNT)
			} else if Authority == Librarian {
				fdu.BorrowInfo(oper1)
			} else {
				fmt.Println("Permission Denied!")
			}
		case "cp":
			if Authority == Student || Authority == Librarian {
				fdu.ChangePassword(CUR_ACCOUNT, Authority)
			} else {
				fmt.Println("Permission Denied!")
			}
		case "cps":
			if Authority == Administrator {
				fdu.ChangePassword(oper1, Student)
			} else {
				fmt.Println("Permission Denied!")
			}
		case "cpl":
			if Authority == Administrator {
				fdu.ChangePassword(oper1, Librarian)
			} else {
				fmt.Println("Permission Denied!")
			}
		default:
			fmt.Printf("Unkonwn command '%s'\n", cmd)
		}

	}
	fdu.db.Close()
}

/*******************************/
/* The Administrator functions */
/*******************************/

// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error {
	sqlExecute(lib.db, []string{
		fmt.Sprintf("CREATE TABLE Librarian (lib_id CHAR(%v) NOT NULL PRIMARY KEY, password CHAR(%v) NOT NULL , lib_status INT NOT NULL , in_date DATE)", LibAccLen, EnPWLen),
		fmt.Sprintf("CREATE TABLE Student (stu_id CHAR(%v) NOT NULL PRIMARY KEY, password CHAR(%v) NOT NULL, stu_status INT NOT NULL, in_date DATE)", StuAccLen, EnPWLen),
		fmt.Sprintf("CREATE TABLE BookInfo (ISBN CHAR(%v) NOT NULL PRIMARY KEY, title VARCHAR(100) NOT NULL, author VARCHAR(100) NOT NULL )", ISBNLen),
		fmt.Sprintf("CREATE TABLE Book (book_id CHAR(%v) NOT NULL PRIMARY KEY, ISBN CHAR(%v) NOT NULL , book_status INT NOT NULL ,exddl INT NOT NULL, in_date DATE, FOREIGN KEY (ISBN)  REFERENCES BookInfo(ISBN))", BookIdLen, ISBNLen),
		fmt.Sprintf("CREATE TABLE Record (rec_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT, stu_id CHAR(%v) NOT NULL , book_id CHAR(%v) NOT NULL , type INT NOT NULL, rec_date DATE NOT NULL , FOREIGN KEY (stu_id)  REFERENCES Student(stu_id), FOREIGN KEY (book_id)  REFERENCES Book(book_id))", StuAccLen, BookIdLen),
		fmt.Sprintf("CREATE TABLE ManageLog (log_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT, lib_id CHAR(%v) NOT NULL, tar_id VARCHAR(%v) NOT NULL, type INT NOT NULL, log_date DATE NOT NULL , notes VARCHAR(100), FOREIGN KEY (lib_id)  REFERENCES Librarian(lib_id))", LibAccLen, BookIdLen),
	})
	return nil
}

// AddStu add a student account
func (lib *Library) AddStu(stu_id, password string) error {
	if len(stu_id) != StuAccLen || len(password) > PasswordMaxlen || len(password) < PasswordMinlen {
		fmt.Println("Illegal Command!")
		return nil
	}
	var cnt int
	err := lib.db.QueryRow(fmt.Sprintf("SELECT count(*) FROM Student WHERE stu_id = '%s'", stu_id)).Scan(&cnt)
	if err != nil {
		panic(err)
	}

	if cnt != 0 {
		fmt.Printf("Account '%s' already exists\n", stu_id)
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	password = string(hash)
	sqlExecute(lib.db, []string{
		fmt.Sprintf("INSERT INTO Student VALUES('%s', '%s', %v, '%s')", stu_id, password, InitialState, VirtualTime()),
	})

	return nil
}

// AddLib add a librarian account
func (lib *Library) AddLib(lib_id, password string) error {
	if len(lib_id) != LibAccLen || len(password) > PasswordMaxlen || len(password) < PasswordMinlen {
		fmt.Println("Illegal Command!")
	}
	var cnt int
	err := lib.db.QueryRow(fmt.Sprintf("SELECT count(*) FROM Librarian WHERE lib_id = '%s'", lib_id)).Scan(&cnt)
	if err != nil {
		panic(err)
	}
	if cnt != 0 {
		fmt.Printf("Account '%s' already exists\n", lib_id)
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	password = string(hash)
	sqlExecute(lib.db, []string{
		fmt.Sprintf("INSERT INTO Librarian VALUES('%s', '%s', %v, '%s')", lib_id, password, InitialState, VirtualTime()),
	})

	return nil
}

// RemoveStu remove a student account
func (lib *Library) RemoveStu(stu_id string) error {
	sqlExecute(lib.db, []string{
		fmt.Sprintf("SET FOREIGN_KEY_CHECKS = 0"),
		fmt.Sprintf("DELETE FROM Student WHERE stu_id = '%s'", stu_id),
		fmt.Sprintf("SET FOREIGN_KEY_CHECKS = 1"),
	})
	return nil
}

// RemoveLib remove a librarian account
func (lib *Library) RemoveLib(lib_id string) error {
	sqlExecute(lib.db, []string{
		fmt.Sprintf("SET FOREIGN_KEY_CHECKS = 0"),
		fmt.Sprintf("DELETE FROM Librarian WHERE lib_id = '%s'", lib_id),
		fmt.Sprintf("SET FOREIGN_KEY_CHECKS = 1"),
	})

	return nil
}

//ChangePassword change the password
func (lib *Library) ChangePassword(id string, t int) error {
	done := true
	var err error
	var old1, old2, new1, new2 string
	if Authority != Administrator {
		if t == Librarian {
			err = lib.db.QueryRow(fmt.Sprintf("SELECT password FROM Librarian WHERE lib_id = '%s' ", id)).Scan(&old1)
		} else if t == Student {
			err = lib.db.QueryRow(fmt.Sprintf("SELECT password FROM Student WHERE stu_id = '%s' ", id)).Scan(&old1)
		}
		if err != nil {
			fmt.Printf("Account '%s' does not exist!\n", id)
			return nil
		}
	}
	for done {
		if Authority != Administrator {
			fmt.Printf("Enter the old password: ")
			fmt.Scanf("%s", &old2)
			err = bcrypt.CompareHashAndPassword([]byte(old1), []byte(old2))
		}
		fmt.Printf("Enter the new password: ")
		fmt.Scanf("%s", &new1)
		fmt.Printf("Enter the new password again:")
		fmt.Scanf("%s", &new2)
		if (Authority != Administrator && err != nil) || new1 != new2 || len(new1) < PasswordMinlen || len(new1) > PasswordMaxlen {
			if err != nil {
				fmt.Printf("Old password is wrong! ")
			} else if new1 != new2 {
				fmt.Printf("The two passwords are inconsistent! ")
			} else {
				fmt.Printf("The length of the password should be between %v and %v! ", PasswordMinlen, PasswordMaxlen)
			}
			fmt.Println("Try again?[Y/N]")
			var c string
			fmt.Scanf("%s", &c)
			if c != "Y" && c != "y" {
				done = false
			}
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(new2), bcrypt.DefaultCost)
			if err != nil {
				panic(err)
			}
			new2 = string(hash)
			if t == Librarian {
				sqlExecute(lib.db, []string{
					fmt.Sprintf("UPDATE Librarian SET password = '%s' WHERE lib_id = '%s' ", new2, id),
				})
			} else if t == Student {
				sqlExecute(lib.db, []string{
					fmt.Sprintf("UPDATE Student SET password = '%s' WHERE stu_id = '%s' ", new2, id),
				})
			}
			return nil
		}
	}
	return nil
}

/*******************************/
/* The Librarians functions    */
/*******************************/

//HistoryStu query the borrow history of a student account
//tag: 302:borrowed 303:reserving 304:overdue other :all
func (lib *Library) HistoryStu(stu_id, tag string) error {
	if len(tag) == 0 {
		tag = "0"
	}
	fmt.Println("---------------Borrow Record-------------------------------------------")
	fmt.Println("id\tbook_id\t\t\ttitle\t\tauther\t\ttype\t\tborrow_date\tbook_status")
	row, err := lib.db.Query(fmt.Sprintf("SELECT book_id, title, author, book_status, type, rec_date "+
		"FROM Record NATURAL JOIN (Book NATURAL JOIN BookInfo)"+
		"WHERE stu_id = '%s' AND (%v < %v OR %v > %v OR type = %v) "+
		"ORDER BY rec_date ", stu_id, tag, Borrow, tag, CancelReserve, tag))
	defer row.Close()
	if err != nil {
		panic(err)
	}
	cnt := 1
	overdue_cnt := 0
	period, _ := time.ParseDuration("-24h")
	now := time.Now()
	for i := 0; i < BorrowPeriod; i++ {
		now = now.Add(period)
	}
	for row.Next() {
		var rec_book_id, rec_title, rec_auther, rec_date string
		var rec_book_status, rec_type int
		_ = row.Scan(&rec_book_id, &rec_title, &rec_auther, &rec_book_status, &rec_type, &rec_date)

		if now.Format("2006-01-02") > rec_date && rec_type == Borrow {
			overdue_cnt++
			fmt.Printf("\033[1;31m%v\t%v\t%v\t\t%v\t\t%v\t\t%v\t%v\033[0m\n", cnt, rec_book_id, rec_title, rec_auther, rec_type, rec_date, rec_book_status)
		} else {
			fmt.Printf("%v\t%v\t%v\t\t%v\t\t%v\t\t%v\t%v\n", cnt, rec_book_id, rec_title, rec_auther, rec_type, rec_date, rec_book_status)
		}
		cnt++
	}
	fmt.Println("type: [306:borrow]  [307:reserve]  [308:return]  [309:cance_reserve]")

	// if overdue_cnt >= 3 {
	// 	lib.SuspendStuAcc(CUR_ACCOUNT, stu_id, "More than three times overdue.")
	// }
	return nil
}

//HistoryBook query the borrow history of a book
func (lib *Library) HistoryBook(book_id string) error {
	row, err := lib.db.Query(fmt.Sprintf("SELECT ISBN, title, author,book_status FROM Book NATURAL JOIN BookInfo WHERE book_id = '%s'", book_id))
	defer row.Close()
	if err != nil {
		panic(err)
	}
	for row.Next() {
		var ISBN, title, author string
		var book_status int
		_ = row.Scan(&ISBN, &title, &author, &book_status)
		fmt.Println("---------------BooKInfo-------------------------------")
		fmt.Println("ISBN\t\ttitle\tauther,\tbook_status")
		fmt.Println(ISBN, "\t", title, "\t", author, "\t", book_status)
		fmt.Println("---------------Borrow History-------------------------")
	}

	cnt := 1
	fmt.Println("id\tstu_id\t\ttype\tdate")
	row, err = lib.db.Query(fmt.Sprintf("SELECT stu_id,type,rec_date FROM Record WHERE book_id = '%s'", book_id))
	defer row.Close()
	for row.Next() {
		var stu_id, date string
		var t int
		_ = row.Scan(&stu_id, &t, &date)
		fmt.Println(cnt, "\t", stu_id, "\t", t, "\t", date)
		cnt++
	}

	return nil
}

//ManageLog add a log
func (lib *Library) ManageLog(lib_id, tar_id, notes string, t int) error {
	sqlExecute(lib.db, []string{
		fmt.Sprintf("SET FOREIGN_KEY_CHECKS = 0"),
		fmt.Sprintf("INSERT INTO ManageLog(lib_id,tar_id,type,log_date,notes) VALUES('%s','%s',%v,'%s','%s')", lib_id, tar_id, t, VirtualTime(), notes),
		fmt.Sprintf("SET FOREIGN_KEY_CHECKS = 1"),
	})
	return nil
}

//SuspendStuAcc suspend student's account if the student has more than 3 overdue books
func (lib *Library) SuspendStuAcc(lib_id, stu_id, notes string) error {
	/*Check: Does this student account exist and be valid?*/
	var status int
	err := lib.db.QueryRow(fmt.Sprintf("SELECT stu_status FROM Student WHERE stu_id = '%s' ", stu_id)).Scan(&status)
	if err != nil {
		fmt.Println("Failed", stu_id)
		return nil
	}
	if status == InvalidAcc {
		sqlExecute(lib.db, []string{
			fmt.Sprintf("UPDATE Student SET stu_status = %v WHERE stu_id = '%s' ", InitialState, stu_id),
		})
	} else {
		sqlExecute(lib.db, []string{
			fmt.Sprintf("UPDATE Student SET stu_status = %v WHERE stu_id = '%s' ", InvalidAcc, stu_id),
		})
	}
	lib.ManageLog(lib_id, stu_id, notes, SuspendStu)
	return nil
}

// AddBook add a book into the library
func (lib *Library) AddBook(ISBN, title, author string) error {
	if len(ISBN) != ISBNLen {
		fmt.Println("Illegal Command")
		return nil
	}
	var cnt int
	err := lib.db.QueryRow(fmt.Sprintf("SELECT count(*) FROM Book WHERE ISBN = '%s'", ISBN)).Scan(&cnt)
	if err != nil {
		panic(err)
	}
	if cnt == 0 { //The first add
		sqlExecute(lib.db, []string{
			fmt.Sprintf("INSERT INTO BookInfo VALUES('%s','%s','%s')", ISBN, title, author),
		})
	}
	book_id := fmt.Sprintf("%sn%.3v", ISBN, cnt+1)
	sqlExecute(lib.db, []string{
		fmt.Sprintf("INSERT INTO Book VALUES('%s','%s',%v,0,'%s')", book_id, ISBN, OnShelf, VirtualTime()),
	})
	lib.ManageLog(CUR_ACCOUNT, book_id, "null", AddBook)
	return nil
}

// RemoveBook remove a book from the library
func (lib *Library) RemoveBook(book_id, notes string) error {

	var isbn string
	err := lib.db.QueryRow(fmt.Sprintf("SELECT ISBN FROM Book WHERE book_id = '%s'", book_id)).Scan(&isbn)
	if err != nil {
		fmt.Printf("Not Found: book_id = %v\n", book_id)
		return nil
	}
	var cnt int
	err = lib.db.QueryRow(fmt.Sprintf("SELECT count(*) FROM Book WHERE ISBN = '%s'", isbn)).Scan(&cnt)
	if err != nil {
		panic(err)
	}

	if cnt == 1 {
		sqlExecute(lib.db, []string{
			fmt.Sprintf("SET FOREIGN_KEY_CHECKS = 0"),
			fmt.Sprintf("DELETE FROM BookInfo WHERE ISBN = '%s'", isbn),
		})
	} else if cnt == 0 {
		fmt.Println("removebook failed")
		return nil
	}
	sqlExecute(lib.db, []string{
		fmt.Sprintf("DELETE FROM Book WHERE book_id = '%s'", book_id),
		fmt.Sprintf("SET FOREIGN_KEY_CHECKS = 1"),
	})
	lib.ManageLog(CUR_ACCOUNT, book_id, notes, RemoveBook)
	return nil
}

/*******************************/
/* The Students functions      */
/*******************************/

// QueryBook query books by title,author or ISBN
func (lib *Library) QueryBook(s1, s2, s3 string) error {
	row, err := lib.db.Query("")
	if s2 == "$" {
		row, err = lib.db.Query(fmt.Sprintf(
			"SELECT book_id,ISBN,title,author,book_status "+
				"FROM Book NATURAL JOIN "+
				"(SELECT * FROM BookInfo "+
				"WHERE %v LIKE '%%%s%%') AS X ", s3, s1))
	} else {
		row, err = lib.db.Query(fmt.Sprintf(
			"SELECT book_id,ISBN,title,author,book_status "+
				"FROM Book NATURAL JOIN "+
				"(SELECT * FROM BookInfo "+
				"WHERE ISBN LIKE '%%%s%%' OR ISBN LIKE '%%%s%%' OR ISBN LIKE '%%%s%%' OR "+
				"title LIKE '%%%s%%' OR title LIKE '%%%s%%' OR title LIKE '%%%s%%' OR "+
				"author LIKE '%%%s%%' OR author LIKE '%%%s%%' OR author LIKE '%%%s%%') AS X ",
			s1, s2, s3, s1, s2, s3, s1, s2, s3))
	}
	defer row.Close()
	if err != nil {
		fmt.Println("Not found")
		return nil
	}
	result := make(map[string]string)
	cnt := 0
	fmt.Println("id\tbook_id\t\t\tISBN\t\ttitle\t\tauther\tbook_status")
	for row.Next() {
		cnt++
		var book_id, ISBN, title, author string
		var book_status int
		_ = row.Scan(&book_id, &ISBN, &title, &author, &book_status)
		fmt.Printf("%v\t%s\t%s\t%s\t\t%s\t%v\n", cnt, book_id, ISBN, title, author, book_status)
		result[fmt.Sprintf("%v", cnt)] = book_id
	}

	if Authority == Student {

		fmt.Println("Borrow book: b [id1] [id2] ... \nReserve book: r [id1] [id2] ... ")
		in := bufio.NewReader(os.Stdin)
		s, _, _ := in.ReadLine()
		cmds := strings.Split(string(s), " ")
		if cmds[0] != "b" && cmds[0] != "r" {
			return nil
		}

		var ids []string
		for _, key := range cmds {
			if _, ok := result[key]; ok {
				ids = append(ids, result[key])
			}
		}
		if cmds[0] == "b" {
			for _, id := range ids {
				lib.BorrowBook(CUR_ACCOUNT, id)
			}
		} else {
			for _, id := range ids {
				lib.ReserveBook(CUR_ACCOUNT, id)
			}
		}
	}
	return nil
}

// BorrowBook Borrow a book by book_id
func (lib *Library) BorrowBook(stu_id, book_id string) error {
	var cnt int
	/*Does the book_id exist?*/
	err := lib.db.QueryRow(fmt.Sprintf("SELECT count(*) FROM Book WHERE book_id = '%s' AND book_status = %v", book_id, OnShelf)).Scan(&cnt)
	if err != nil {
		panic(err)
	}
	if cnt == 0 {
		fmt.Println("Failed borrow:", book_id)
		return nil
	} else {

		/* MaxBorBookNun ?*/
		err = lib.db.QueryRow(fmt.Sprintf("SELECT stu_status FROM Student WHERE stu_id = '%s' ", stu_id)).Scan(&cnt)
		if err != nil {
			panic(err)
		}
		if cnt-InitialState >= MaxBorBookNum {
			fmt.Printf("The number of books you are allowed to borrow is the maximum : %v.\n", MaxBorBookNum)
			return nil
		} else {
			/*Successed*/
			sqlExecute(lib.db, []string{
				fmt.Sprintf("UPDATE Book SET book_status = %v WHERE book_id = '%s'", OnLoan, book_id),
				fmt.Sprintf("UPDATE Book SET exddl = 1 WHERE book_id = '%s'", book_id),
				fmt.Sprintf("UPDATE Student SET stu_status = stu_status + 1 WHERE stu_id = '%s' ", stu_id),
			})

			lib.AddRecord(stu_id, book_id, Borrow)
			fmt.Printf("%s borrowed:%s successfully\n", stu_id, book_id)

		}
	}
	return nil
}

//AddRecord add a record
func (lib *Library) AddRecord(stu_id, book_id string, t int) error {
	sqlExecute(lib.db, []string{
		fmt.Sprintf("SET FOREIGN_KEY_CHECKS = 0"),
		fmt.Sprintf("INSERT INTO Record(stu_id, book_id, type, rec_date) VALUES('%s','%s',%v,'%s')", stu_id, book_id, t, VirtualTime()),
		fmt.Sprintf("SET FOREIGN_KEY_CHECKS = 1"),
	})
	return nil
}

// ReturnBook return a book to the library
func (lib *Library) ReturnBook(stu_id, book_id string) error {
	/*Check: Is this the book you already borrowed ?*/
	var cnt int
	err := lib.db.QueryRow(fmt.Sprintf(
		"SELECT count(*) "+
			"FROM Record "+
			"WHERE rec_id IN (SELECT max(rec_id) FROM Record WHERE type = %v AND book_id = '%s') AND "+
			"stu_id = '%s' AND "+
			"1 = (SELECT count(*) FROM Book WHERE book_id = '%s' AND (book_status = %v OR book_status = %v)) ", Borrow, book_id, stu_id, book_id, OnLoan, OnReserve)).Scan(&cnt)
	if err != nil {
		panic(err)
	}
	if cnt == 0 {
		fmt.Printf("Failed return: %s\ns", book_id)
		return nil
	} else {
		/*Successed return*/
		sqlExecute(lib.db, []string{
			fmt.Sprintf("UPDATE Book SET book_status = %v WHERE book_id = '%s' AND book_status = %v ", OnShelf, book_id, OnLoan),
			fmt.Sprintf("UPDATE Book SET exddl = 0 WHERE book_id = '%s'", book_id),
			fmt.Sprintf("UPDATE Student SET stu_status = stu_status - 1 WHERE stu_id = '%s' ", stu_id),
		})

		lib.AddRecord(stu_id, book_id, Return)
		fmt.Printf("%s returned:%s successfully \n", stu_id, book_id)

	}
	return nil
}

//ReserveBook reserve a book
func (lib *Library) ReserveBook(stu_id, book_id string) error {
	/*Check :Is this a book that has been borrowed by others ?*/
	var cnt int
	err := lib.db.QueryRow(fmt.Sprintf(
		"SELECT count(*) "+
			"FROM Record "+
			"WHERE rec_id = (SELECT max(rec_id) FROM Record WHERE type = %d AND book_id = '%s') AND "+
			"stu_id != '%s' AND "+
			"1 = (SELECT count(*) FROM Book WHERE book_id = '%s' AND book_status = %d )", Borrow, book_id, stu_id, book_id, OnLoan)).Scan(&cnt)
	if err != nil {
		panic(err)
	}
	if cnt == 0 {
		fmt.Println("Failed reserve:", book_id)
		return nil
	}

	/*Reserve Successed */
	sqlExecute(lib.db, []string{
		fmt.Sprintf("UPDATE Book SET book_status = %v WHERE book_id = '%s'", OnReserve, book_id),
	})

	lib.AddRecord(stu_id, book_id, Reserve)
	fmt.Printf("%s reservedd:%s successfully\n", stu_id, book_id)
	return nil
}

//CancelReserve cancel the reservation
func (lib *Library) CancelReserve(stu_id, book_id string) error {
	var cnt int
	err := lib.db.QueryRow(fmt.Sprintf("SELECT count(*) FROM Book WHERE book_id = '%s' AND book_status = %v", book_id, OnReserve)).Scan(&cnt)
	if err != nil {
		panic(err)
	}
	if cnt == 0 {
		fmt.Println("Failed to cancle reserve:", book_id)
		return nil
	}
	sqlExecute(lib.db, []string{ //Successed
		fmt.Sprintf("UPDATE Book SET book_status = %v WHERE book_id = '%s'", OnLoan, book_id),
	})

	lib.AddRecord(stu_id, book_id, CancelReserve)
	fmt.Printf("%s cancle reserve:%s successfully\n", stu_id, book_id)
	return nil
}

//ExtendDDL extend the deadline of returning a book, at most 3 times
func (lib *Library) ExtendDDL(stu_id, book_id string) error {
	/*Check :exceeded MaxExtendDDL and this book is reserved by other?*/
	var cnt, book_status int
	err := lib.db.QueryRow(fmt.Sprintf("SELECT exddl,book_status FROM Book WHERE book_id = '%s'", book_id)).Scan(&cnt, &book_status)
	if err != nil {
		fmt.Printf("Not Found in BorrowInfo : book_id = %v\n", book_id)
		return nil
	}
	if cnt > MaxExtendDDL {
		fmt.Println("No more ExtendDDL.")
		return nil
	}
	if book_status == OnReserve {
		fmt.Println("This book is reserved by other.")
		return nil
	}

	/*Reserve ExtendDDL*/
	sqlExecute(lib.db, []string{ //Successed
		fmt.Sprintf("UPDATE Book SET exddl = %d WHERE book_id = '%s'", cnt+1, book_id),
	})

	lib.AddRecord(stu_id, book_id, Borrow)
	fmt.Printf("%s exddl:%s successfully\n", stu_id, book_id)
	return nil
}

//GetReservedBook
func (lib *Library) GetReservedBook(stu_id, book_id string) error {
	/*MaxBorBookNum ?*/
	var cnt int
	err := lib.db.QueryRow(fmt.Sprintf("SELECT stu_status FROM Student WHERE stu_id = '%s' ", stu_id)).Scan(&cnt)
	if err != nil {
		panic(err)
	}
	if cnt-InitialState >= MaxBorBookNum {
		fmt.Println("Too many books you borrowed, get reserved book failed")
		sqlExecute(lib.db, []string{
			fmt.Sprintf("UPDATE Book SET book_status = %v WHERE book_id = '%s'", OnShelf, book_id),
		})
		return nil
	}
	/*Successed*/
	sqlExecute(lib.db, []string{
		fmt.Sprintf("UPDATE Book SET book_status = %v WHERE book_id = '%s'", OnLoan, book_id),
		fmt.Sprintf("UPDATE Book SET exddl = 1 WHERE book_id = '%s'", book_id),
		fmt.Sprintf("UPDATE Student SET stu_status = stu_status + 1 WHERE stu_id = '%s' ", stu_id),
	})
	lib.AddRecord(stu_id, book_id, Borrow)
	fmt.Printf("%s get reserve:% s successfully\n", stu_id, book_id)
	return nil
}

//BorrowInfo View Borrow info(ddl, overdue ,reserving...)
func (lib *Library) BorrowInfo(stu_id string) error {
	/*Query in SQL*/

	fmt.Println("----------------------------------------------BorrowInfo-----------------------------------------------------------")
	fmt.Println("id\tbook_id\t\t\ttitle\t\tauther\t\tborrow_date\treturn_date")
	row, err := lib.db.Query(fmt.Sprintf(
		"WITH X (book_id,ISBN,book_status,rec_id,exddl) AS "+
			"(SELECT book_id,ISBN,book_status,max(rec_id) ,exddl "+
			"FROM Book NATURAL JOIN Record "+
			"WHERE book_status = %v AND type = %v "+
			"GROUP BY 1,2 "+
			"UNION "+
			"SELECT book_id,ISBN,book_status,max(rec_id),exddl "+
			"FROM Book NATURAL JOIN Record "+
			"WHERE book_status = %v AND type = %v "+
			"GROUP BY 1,2 "+
			"UNION "+
			"SELECT book_id,ISBN,book_status,max(rec_id),exddl "+
			"FROM Book NATURAL JOIN Record "+
			"WHERE book_status = %v AND type = %v "+
			"GROUP BY 1,2) "+
			"SELECT X.book_id, title, author, book_status,rec_date AS bo_date,date_add(rec_date, interval %d day) AS re_date ,exddl ,stu_id,type "+
			"FROM (BookInfo NATURAL JOIN X) JOIN Record USING(rec_id) "+
			"WHERE book_status = %v OR stu_id = '%s'", OnLoan, Borrow, OnReserve, Borrow, OnReserve, Reserve, BorrowPeriod, OnReserve, stu_id))
	defer row.Close()
	if err != nil {
		panic(err)
	}
	/*Print*/
	result := make(map[string]string)
	reserve_id := make(map[string]string)
	var get_id []string
	cnt := 0
	for row.Next() {
		cnt++
		var book_id, title, author, bo_date, re_date, obj_id string
		var book_status, exddl, t int
		_ = row.Scan(&book_id, &title, &author, &book_status, &bo_date, &re_date, &exddl, &obj_id, &t)

		if obj_id == stu_id {

			if VirtualTime() > re_date { //overdue color:front-red
				fmt.Printf("\033[1;31m")
			}
			if book_status == OnReserve {
				if t == Borrow {
					fmt.Printf("\033[1;43m") //The book is reserved by other color: back-yellow
				} else if exddl == 0 { //getreserve color:front-green
					get_id = append(get_id, book_id)
					fmt.Printf("\033[1;32m")
					re_date = ""
				} else { //waiting for being returned color: front-blue
					fmt.Printf("\033[1;34m")
					re_date = ""
				}
			}
			fmt.Printf("%v\t%v\t%v\t\t%v\t\t%v\t%v\n\033[0m", cnt, book_id, title, author, bo_date, re_date)
			if t == Reserve && exddl != 0 {
				reserve_id[fmt.Sprintf("%v", cnt)] = book_id
			} else if exddl != 0 {
				result[fmt.Sprintf("%v", cnt)] = book_id
			}
		}

	}
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("Today is:", VirtualTime())
	fmt.Printf("\033[1;31mThe book is overdue\033[0m\t\033[1;32mThe book is free now\033[0m\t\033[1;34mwaiting for the book to be returned\033[0m\t\033[1;43mThe book is reserved by other\033[0m\n")
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------")

	if Authority == Student {
		/*Retrun or Cancel Reserve or extend DDL*/
		fmt.Println("some books you reserved is free now, ready to borrow...")
		for _, val := range get_id {
			lib.GetReservedBook(CUR_ACCOUNT, val)
		}
		/*continue*/
		fmt.Println("Return book: r [id1] [id2] ...(Or 'r a' return all) \nCancle reserve book: c [id1] [id2] ... (Or 'c a' cancel all)\nExtend DDL: e [id1] [id2] ... (Or 'e a' extendDdl all) \n'q'to continue.")
		in := bufio.NewReader(os.Stdin)
		s, _, _ := in.ReadLine()
		cmds := strings.Split(string(s), " ")
		var last string
		for _, key := range cmds {
			switch key {
			case "q":
				return nil
			case "r":
				last = key
				break
			case "c":
				last = key
				break
			case "e":
				last = key
				break
			case "a":
				switch last {
				case "r":
					for _, val := range result {
						lib.ReturnBook(CUR_ACCOUNT, val)
					}
					break
				case "c":
					for _, val := range reserve_id {
						lib.CancelReserve(CUR_ACCOUNT, val)
					}
					break
				case "e":
					for _, val := range result {
						lib.ExtendDDL(CUR_ACCOUNT, val)
					}
					break
				default:
					fmt.Println("An input error is detected.")
					return nil
				}
				last = ""
				break
			default:
				switch last {
				case "r":
					if _, ok := result[key]; ok {
						lib.ReturnBook(CUR_ACCOUNT, result[key])
					} else {
						fmt.Println("An input error is detected.")
						return nil
					}
					break
				case "c":
					if _, ok := reserve_id[key]; ok {
						lib.CancelReserve(CUR_ACCOUNT, reserve_id[key])
					} else {
						fmt.Println("An input error is detected.")
						return nil
					}
					break
				case "e":
					if _, ok := result[key]; ok {
						lib.ExtendDDL(CUR_ACCOUNT, result[key])
					} else {
						fmt.Println("An input error is detected.")
						return nil
					}
					break
				default:
					fmt.Println("An input error is detected.")
					return nil
				}
			}
		}
	}
	return nil

}

func main() {
	RUN(Initialization())
}
