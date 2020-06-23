package main

import (
	"regexp"
	// "go/reader"
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	// mysql connector

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const DBName = "ass3"

var (
	User     = "root"
	Password = "123456"
	Isadmin  = true
	// reader   *bufio.Reader
	scanner *bufio.Scanner
)

type Library struct {
	db *sqlx.DB
}
type Book struct {
	ISBN,
	Title,
	Author string
	Timer          int
	Holder, Remark sql.NullString
	Date           mysql.NullTime
}

func (book *Book) available() bool {
	return (book.Timer == 0) || (book.Timer < 4 && book.Holder.Valid && book.Holder.String == User)
}
func (book Book) String() string {
	return fmt.Sprintf("Book %s, %s, ISBN %s", book.Title, book.Author, book.ISBN)
}

// this helper function execute a sql, panic any error
func (lib *Library) Exec(que string, args ...interface{}) sql.Result {
	result, err := lib.db.Exec(que, args...)
	if err != nil {
		panic(err)
	}
	return result
}

// this helper function execute a sql, panic any error, return the result as *sqlx.rows
func (lib *Library) Query(que string, args ...interface{}) *sqlx.Rows {
	results, err := lib.db.Queryx(que, args...)
	if err != nil {
		panic(err)
	}
	return results
}
func (lib *Library) initStudentRole() {
	lib.Exec(`drop role if exists 'student'@'localhost';`)
	lib.Exec(`create role 'student'@'localhost';`)
	lib.Exec(fmt.Sprintf(`grant select, update on %s.books to 'student'@'localhost';`, DBName))
	lib.Exec(fmt.Sprintf(`grant select, insert on %s.history to 'student'@'localhost';`, DBName))
}

func (lib *Library) ConnectDB() {
	{
		db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", User, Password))
		defer db.Close()
		if err != nil {
			panic(err)
		}
		lib.db = db
		lib.Exec(fmt.Sprintf("create database if not exists `%s`;", DBName))
		lib.Exec(fmt.Sprintf("use `%s`;", DBName))
		lib.CreateTables()
	}
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	if err != nil {
		panic(err)
	}
	lib.db = db
}

// add student account by the administrator's account so that the student could borrow books from the library
func (lib *Library) CreateAccount(name, pass string) error {
	if !Isadmin {
		return fmt.Errorf("CreateAccount: require privilege escalation.")
	}
	if !checkname(name) {
		return fmt.Errorf("CreateAccount: Invalid user name.")
	}
	if !checkname(pass) {
		return fmt.Errorf("CreateAccount: Invalid user password.")
	}
	lib.Exec(fmt.Sprintf("drop user if exists %s@'localhost';", name))
	lib.Exec(fmt.Sprintf(`CREATE USER %s@'localhost' identified by '%s';`, name, pass))
	lib.Exec(fmt.Sprintf("grant 'student'@'localhost' to %s@'localhost'", name))

	return nil
}
func createAccount() {
	name := readword()
	pass := readword()
	err := mlib.CreateAccount(name, pass)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error {

	if lib.Query(`SHOW TABLES LIKE 'books';`).Next() {
		return nil
	}
	if !Isadmin {
		panic(fmt.Errorf("CreateTables: Library system not initiated\n\tlog in as administrator to initiate the system"))
	}
	lib.Exec(`
CREATE TABLE if not exists books (
ISBN CHAR(17) NOT NULL,
Title VARCHAR(45) NOT NULL,
Author VARCHAR(45) NOT NULL,
timer TINYINT NOT NULL,
holder VARCHAR(45) NULL,
remark VARCHAR(45) NULL,
date DATE NULL,
PRIMARY KEY (ISBN),
UNIQUE INDEX ISBN_UNIQUE (ISBN ASC) VISIBLE);
	`)
	lib.Exec(`
CREATE TABLE ass3.history (
	ISBN CHAR(17) NOT NULL,
	ID VARCHAR(45) NULL,
	action VARCHAR(45) NOT NULL,
	date DATE NOT NULL,
	CONSTRAINT ISBN
		FOREIGN KEY (ISBN)
		REFERENCES ass3.books (ISBN));
	  `)
	lib.initStudentRole()
	return nil
}

// AddBook add a book into the library
func (lib *Library) AddBook(title, author, ISBN string) error {
	if !Isadmin {
		return fmt.Errorf("AddBook: require privilege escalation.")
	}
	if !checkISBN(ISBN) {
		return fmt.Errorf("AddBook: %s is not a valid ISBN.", ISBN)
	}
	{
		results := lib.Query("select * from books where ISBN = ?", ISBN)
		if results.Next() {
			return fmt.Errorf("AddBook: Book %s is already in the library.", ISBN)
		}
	}
	lib.Exec("insert into books(ISBN, Title, Author, timer) values (?,?,?,0)", ISBN, title, author)
	return nil
}
func addBook() {
	title := readword()
	author := readword()
	isbn := readword()
	err := mlib.AddBook(title, author, isbn)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// remove a book from the library with explanation (e.g. book is lost)
func (lib *Library) RemoveBookByISBN(ISBN, remarks string) error {
	if !Isadmin {
		return fmt.Errorf("RemoveBookByISBN: require privilege escalation.")
	}
	result := lib.Exec("update books set remark = ?,timer = 4, holder = NULL where ISBN = ?", remarks, ISBN)
	if rowsc, _ := result.RowsAffected(); rowsc == 0 {
		return fmt.Errorf("RemoveBookByISBN: Book %s is not in the library.", ISBN)
	}
	return nil
}
func removeBookByISBN() {
	isbn := readword()
	remarks := readword()
	err := mlib.RemoveBookByISBN(isbn, remarks)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// query books by title, author or ISBN
func (lib *Library) QueryBooks(Qtyp, refname string) ([]Book, error) {
	QtypU := strings.ToUpper(Qtyp)
	if QtypU != "ISBN" && QtypU != "TITLE" && QtypU != "AUTHOR" && QtypU != "HOLDER" {
		return nil, fmt.Errorf("QueryBooks: invalid Query type %s.", Qtyp)
	}
	//
	// que := fmt.Sprintf(`select * from books where %s = %s`, Qtyp, refname)
	que := fmt.Sprintf(`select * from books where %s = ? and (timer < 4 or holder is not null)`, Qtyp)
	res := lib.Query(que, refname)
	books := make([]Book, 0)
	for res.Next() {
		books = append(books, Book{})
		lbook := &(books[len(books)-1])
		res.Scan(&lbook.ISBN, &lbook.Title, &lbook.Author, &lbook.Timer, &lbook.Holder, &lbook.Remark, &lbook.Date)
	}
	return books, nil
}
func queryBooks() {
	Qtyp := readword()
	refname := readword()
	clc, err := mlib.QueryBooks(Qtyp, refname)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range clc {
		fmt.Println(v)
	}
}

// borrow a book from the library with a student account
func (lib *Library) BorrowBook(isbn string) error {
	books, err := lib.QueryBooks("ISBN", isbn)
	if err != nil {
		return err
	}
	if len(books) == 0 {
		return fmt.Errorf("BorrowBook: no such book %s.", isbn)
	}
	if !books[0].available() {
		return fmt.Errorf("cannot borrow/Renew this book")
	}
	lib.Exec("update books set timer = timer + 1, holder = ?, date = curdate() where ISBN = ?", User, isbn)
	lib.Exec("insert into history(ISBN, ID, action, date) values (?,?,'borrow',curdate())", isbn, User)
	return nil
}
func borrowBook() {
	isbn := readword()
	err := mlib.BorrowBook(isbn)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// query the borrow history of a student account
func (lib *Library) BorrowHistory() ([]string, error) {
	res := lib.Query(`select books.ISBN, Title, history.date
	from history inner join books on history.ISBN = books.ISBN
	 where ID = ? and action = 'borrow' order by history.date`, User)
	ss := make([]string, 0)
	for res.Next() {
		var (
			isbn, title string
			date        mysql.NullTime
		)
		res.Scan(&isbn, &title, &date)
		ss = append(ss, fmt.Sprintf("ISBN %s, %s, %s", isbn, title, date.Time.Format("2006-01-02")))
	}
	return ss, nil
}

func borrowHistory() {
	clc, err := mlib.BorrowHistory()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range clc {
		fmt.Println(v)
	}
}

// query the books a student has borrowed and not returned yet
func (lib *Library) HeldBooks() ([]Book, error) {
	return lib.QueryBooks("HOLDER", User)
}
func heldBooks() {
	clc, err := mlib.HeldBooks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range clc {
		fmt.Println(v)
	}
}

// check the deadline of returning a borrowed book
func (lib *Library) CheckBookddl(isbn string) (time.Time, error) {
	books, err := lib.QueryBooks("ISBN", isbn)
	if err != nil {
		return time.Time{}, err
	}
	if len(books) == 0 {
		return time.Time{}, fmt.Errorf("CheckBookddl: no such book %s.", isbn)
	}
	book := books[0]
	if book.Timer == 0 || !book.Holder.Valid || book.Holder.String != User {
		return time.Time{}, fmt.Errorf("CheckBookddl: You are not possessing the book.")
	}
	return book.Date.Time.AddDate(0, 0, 30), nil
}
func checkBookddl() {
	isbn := readword()
	t, err := mlib.CheckBookddl(isbn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(t)
}

// check if a student has any overdue books that needs to be returned
func (lib *Library) CheckIfHeldOverdue() bool {
	res := lib.Query(`select 1 from books where
	holder = ? and date > DATE_ADD(CURDATE(), INTERVAL -30 DAY)`, User)
	return res.Next()
}
func checkIfHeldOverdue() {
	t := mlib.CheckIfHeldOverdue()
	fmt.Println(t)
}

// return a book to the library by a student account (any student account will do)
func (lib *Library) ReturnBook(isbn string) error {
	result := lib.Exec(`update books set timer = 0,holder=NULL,date=null where ISBN = ?`, isbn)
	if rc, _ := result.RowsAffected(); rc == 0 {
		return fmt.Errorf("ReturnBook: Book %s not found", isbn)
	}
	return nil
}
func returnBook() {
	isbn := readword()
	err := mlib.ReturnBook(isbn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// etc...

func checkISBN(ISBN string) bool {
	res, err := regexp.MatchString(`\d{3}-\d-\d{3}-\d{5}-\d$`, ISBN)
	if err != nil {
		panic(err)
	}
	return res
}
func checkname(name string) bool {
	res, err := regexp.MatchString(`^\w+$`, name)
	if err != nil {
		panic(err)
	}
	return res
}

var fmap = map[string]func(){
	"addbook":            addBook,
	"borrowbook":         borrowBook,
	"borrowhistory":      borrowHistory,
	"checkbookddl":       checkBookddl,
	"checkifheldoverdue": checkIfHeldOverdue,
	"createaccount":      createAccount,
	"heldbooks":          heldBooks,
	"querybooks":         queryBooks,
	"removebookbyisbn":   removeBookByISBN,
	"returnbook":         returnBook,
}
var mlib = Library{nil}

func main() {
	// for {
	fmt.Println("Welcome to the Library Management System!")

	// file, err := os.Open("input")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// scanner = bufio.NewScanner(file)
	scanner = bufio.NewScanner(os.Stdin)

	scanner.Split(bufio.ScanWords)
	// os.exec("clear")
	// fmt.Print("Plz input your user name: ")
	// User = readword()
	// fmt.Print("Plz input your account password: ")
	// Password = readword()
	mlib.ConnectDB()
	fmt.Printf("logged in as %s.\n", User)
	// mlib.CreateTables()
	{
		res := mlib.Query("select CURRENT_ROLE();")

		if !res.Next() {
			Isadmin = true
		} else {
			s := ""
			res.Scan(&s)
			Isadmin = s != "student"
		}
	}
	for {
		cmd := readword()
		if cmd == "logout" || cmd == "" {
			break
		}
		if cmd == "help" {
			for v, _ := range fmap {
				fmt.Print(v, " ")
			}
			fmt.Println()
		} else if fmap[cmd] == nil {
			fmt.Printf("unknown command %s.\n", cmd)
		} else {
			fmap[cmd]()
		}

	}
	fmt.Println("Goodbye!!!!")
	// }
}
func readword() string {
	if scanner.Scan() {
		return scanner.Text()
	} else {
		return ""
	}
	// tex, err := reader.ReadString(' ')
	// if err != nil {
	// 	panic(err)
	// }
	// return strings.TrimSpace(tex)
}
