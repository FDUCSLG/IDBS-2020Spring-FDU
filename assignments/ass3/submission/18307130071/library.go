package main

import (
	"fmt"
	"time"

	// mysql connector
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	User     = "YourUserName"
	Password = "YourPassword"
	DBName   = "ass3"
)

type Library struct {
	db *sqlx.DB
}

func ConnectDB (lib *Library) {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	if err != nil {
		panic(err)
	}
	lib.db = db
}

// CreateTables created the tables in MySQL
func  CreateTables(lib *Library) error {
	_, err := lib.db.Exec("create table BOOKS (title varchar(100) not null, author varchar(30)not null, ISBN char(13)not null, status smallint not null, primary key(ISBN))")
	if err != nil {
		fmt.Println("Illegal operation, tables creation failed!")
		panic(err)
	}
	_, err = lib.db.Exec("create table ACCOUNTS (priority smallint not null, accID varchar(11)not null, pwd char(6), overdue_records smallint not null, primary key(accID))")
	if err != nil {
		fmt.Println("Illegal operation, tables creation failed!")
		panic(err)
	}
	_, err = lib.db.Exec("create table Borrow_records (ISBN char(13)not null, accID varchar(11)not null, DDL DATETIME(0), extend smallint not null, status smallint not null, primary key(ISBN, accID, status), foreign key(ISBN) references BOOKS(ISBN), foreign key(accID) references ACCOUNTS(accID))")
	if err != nil {
		fmt.Println("Illegal operation, tables creation failed!")
		panic(err)
	}
	_, err = lib.db.Exec("create table REMOVAL (title varchar(100) not null, author varchar(30)not null, original_ISBN char(13)not null, remove_date DATETIME(0)not null, explanation char(13)not null, primary key(original_ISBN, remove_date))")
	if err != nil {
		fmt.Println("Illegal operation, tables creation failed!")
		panic(err)
	}
	fmt.Println("Tables created")
	return nil
}

//InitializeDB, add 20 books and an administrator account
func InitializeDB(lib *Library, admin_pwd string) error {
	stmt, err := lib.db.Prepare("INSERT into BOOKS(title, author, ISBN, status) values (?, ?, ?, ?)")
	_, err = stmt.Exec("Constitution_of_FDU", "FudanUniversity", "0000000000001", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("Computer_age_statistical", "Efron", "2642176672401", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("Recent_Advances_in_NLP: The Case of Arabic Language", "Abd_Elaziz", "8975051686568", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("Essentials_of_computer_architecture", "Comer", "4690135027004", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("Financial_mathematics", "Bai_Dongjie", "4087013757311", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("The_foundations_of_mathematics", "Bradley", "7803189487455", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("Discrete_mathematics_and_its_applications", "Fu_yan", "8894393285891", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("Mathematics_for_computer_science", "Lehman", "8264711307003", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("A_global_history:_from_prehistory_to_the_21st_century", "Stavrianos", "6037958370961", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("A_history_of_god", "Armstrong", "3178651024459", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("The_gunpowder_age", "Andrade", "9509479575053", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("Natural_histories", "Baione", "7050747211446", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("Technopoly:_the_surrender_of_culture_to_technology", "Postman", "9791776101900", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("Bio-inspired_artificial_intelligence", "Floreano", "6849222079941", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("FOOD_BIO-CHEMISTRY", "Ning_Zhengxiang", "1713267825040", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("The_complete_classical_music_guide", "Burrows", "9625788155322", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("How_to_hear_classical_music", "Caddy", "5340531960820", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("The_essential_canon_of_classical_music", "Dubal", "7209576214405", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("Affluenza:_the_all-consuming_epidemic", "Graaf", "7530544598350", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	_, err = stmt.Exec("Epidemic_and_civilization", "Xin_Zhengren", "5006868242797", 1)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	

	Set_admin := fmt.Sprintf("INSERT into ACCOUNTS(priority, accID, pwd, overdue_records) values(0, 'admin', '%s', 0)", admin_pwd)
	_, err = lib.db.Exec(Set_admin)
	if err != nil {
		fmt.Println("Illegal operation, initialization failed!")
		panic(err)
	}
	fmt.Println("Initialization Done")
	return nil

}

// AccountUpdate update overdue_records for each account
func AccountUpdate(lib *Library) error {
	Date := time.Now().Format("2006-01-02 15:04:05") 
	check_due := fmt.Sprintf("SELECT accID, status from Borrow_records where DDL < '%s'", Date)
	out_of_due, err := lib.db.Query(check_due)
	for out_of_due.Next() {
		var ID string
		var status int
		err = out_of_due.Scan(&ID, &status)
		if err != nil {
			fmt.Println("Failed to check overdue!")
			panic(err)
		}
		if status==0 {
		update_account := fmt.Sprintf("UPDATE ACCOUNTS SET overdue_records=overdue_records+1 where accID = '%s'", ID)
			_, err = lib.db.Exec(update_account)
			if err != nil {
				fmt.Println("Failed to update overdue_records!")
				panic(err)
			}
			update_borrow_records := fmt.Sprintf("UPDATE Borrow_records SET status=-1 where accID = '%s' and status = 0 and DDL < '%s'", ID, Date)
			_, err = lib.db.Exec(update_borrow_records) 
			if err != nil {
				fmt.Println("Failed to update borrow_records_status!")
				panic(err)
			}
			//-1 means overdue, not returned and recorded
		}
	}
	_, err = lib.db.Exec("UPDATE ACCOUNTS SET priority = -1 WHERE overdue_records > 3")
	if err != nil {
			fmt.Println("Failed to update account status!")
			panic(err)
		}
	fmt.Println("Account status updated.")
	return nil
}



// AddBook add a book into the library
func AddBook(lib *Library, title, author, ISBN string) error {
	addBook := fmt.Sprintf("INSERT into BOOKS(title, author, ISBN, status) values('%s', '%s', '%s', 1)", title, author, ISBN)
	_, err := lib.db.Exec(addBook)
	if err != nil {
		fmt.Println("Illegal operation, insertion failed!")
		panic(err)
	}
	return nil
}

// DeleteBook delete a book from library with explanations
func DeleteBook(lib *Library, ISBN string) error {
	var del_info, cnt int
	query_cnt := fmt.Sprintf("SELECT count(*) from BOOKS where ISBN = '%s'", ISBN)
	cnt_result, err := lib.db.Query(query_cnt)
	if err != nil {
		fmt.Println("Invalid operation!")
		panic(err)
	}
	fmt.Println("Choose an expalanation: ")
	fmt.Println("--1.Book is lost.")
	fmt.Println("--2.Book damaged.")
	fmt.Println("--3.Book too old.")
	cnt_result.Next()
	cnt_result.Scan(&cnt)
	if cnt == 0 {
		fmt.Printf("No book match with ISBN %s\n", ISBN)
		return nil
	}
	fmt.Scanln(&del_info)
	query_book := fmt.Sprintf("SELECT title, author from BOOKS where ISBN = '%s'", ISBN)
	rows, err := lib.db.Query(query_book)
	if err != nil {
		fmt.Println("Invalid operation!")
		panic(err)
	}
	for rows.Next() {
		Date := time.Now()
		var title string
		var author string
		var explanation string		
		err = rows.Scan(&title, &author)
		if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
		}
		switch del_info {
			case 1: explanation = "Book is lost."
			case 2: explanation = "Book damaged."
			case 3: explanation = "Book too old."
		}
		removal, err := lib.db.Prepare("INSERT into REMOVAL (title, author, original_ISBN, remove_date, explanation) values(?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println("Invalid removal!")
			panic(err)
		}
		removal.Exec(title, author, ISBN, Date, explanation)
		del_borrow_records := fmt.Sprintf("DELETE from Borrow_records where ISBN = '%s'", ISBN)
		_, err = lib.db.Exec(del_borrow_records)
		if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
		}
		del := fmt.Sprintf("DELETE from BOOKS where ISBN = '%s'", ISBN)
		_, err = lib.db.Exec(del)
		if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
		}
	}
	fmt.Printf("Book ISBN %s removed from library.\n", ISBN)
	return nil
}

// AddStu add student account
func AddStu(lib *Library, ID, pwd string) error {
	addstu := fmt.Sprintf("INSERT into ACCOUNTS(priority, accID, pwd, overdue_records) values(1, '%s', '%s', 0)", ID, pwd)
	_, err := lib.db.Exec(addstu)
	if err != nil {
		fmt.Println("Illegal operation, add student account failed!")
		panic(err)
	}
	fmt.Printf("Student account %s added.\n", ID)
	return nil
}

// SearchBooks
func SearchBooks(lib *Library) error {
	fmt.Println("--1.Search with title.")
	fmt.Println("--2.Search with author.")
	fmt.Println("--3.Search with ISBN.")
	var choice int
	var key string
	var search, count string
	fmt.Scanln(&choice)
	switch choice {
		case 1: {fmt.Println("Input title:")
			 fmt.Scanln(&key)
			 search = fmt.Sprintf("SELECT title, author, ISBN, status from BOOKS where title = '%s'", key)
			 count = fmt.Sprintf("SELECT count(*) from BOOKS where title = '%s'", key)}
		case 2: {fmt.Println("Input author:")
			 fmt.Scanln(&key)
			 search = fmt.Sprintf("SELECT title, author, ISBN, status from BOOKS where author = '%s'", key)
			 count = fmt.Sprintf("SELECT count(*) from BOOKS where author = '%s'", key)}
		case 3: {fmt.Println("Input ISBN:")
			 fmt.Scanln(&key)
			 search = fmt.Sprintf("SELECT title, author, ISBN, status from BOOKS where ISBN = '%s'", key)
			 count = fmt.Sprintf("SELECT count(*) from BOOKS where ISBN = '%s'", key)}
	}
	Checkcnt, err := lib.db.Query(count)
	if err != nil {
		fmt.Println("Illegal operation, Book count failed!")
		panic(err)
	}
	Checkcnt.Next()
	var cnt int
	Checkcnt.Scan(&cnt)
	if cnt==0 {
		fmt.Println("No suitable books found.")
		return nil
	}
	SearchBook, err := lib.db.Query(search)
	if err != nil {
		fmt.Println("Illegal operation, Book search failed!")
		panic(err)
	}
	for SearchBook.Next() {
		var title, author, ISBN, borrow_status string
		var status int
		SearchBook.Scan(&title, &author, &ISBN, &status)
		borrow_status = "Borrowed"
		if status == 1 {
			borrow_status = "Available"
		}
		fmt.Printf("  Title: %s, Author: %s, ISBN: %s, Book status: %s\n", title, author, ISBN, borrow_status)
		
	}
	return nil
} 

// Borrowbooks with ISBN (you can search for ISBN first if you don't know ISBN)
// DDL is 2 weeks originally
func Borrowbooks(lib *Library, ISBN , usrID string) error {
	var cnt int
	Query_num := fmt.Sprintf("SELECT count(*) from BOOKS where ISBN = '%s'", ISBN)
	num_result, err := lib.db.Query(Query_num)
	if err != nil {
		fmt.Println("Invalid operation!")
		panic(err)
	}
	num_result.Next()
	num_result.Scan(&cnt)
	if cnt==0 {
		fmt.Println("Book not found, no such book.")
		return nil
	}
	borrow := fmt.Sprintf("SELECT status from BOOKS where ISBN = '%s'", ISBN)
	var status int
	borrow_result, err := lib.db.Query(borrow)
	if err != nil {
		fmt.Println("Invalid operation!")
		panic(err)
	}
	borrow_result.Next()
	err = borrow_result.Scan(&status)
	if err != nil {
		fmt.Println("Invalid operation!")
		panic(err)
	}
	if status == 0 {
		fmt.Println("Book already borrowed.")
		return nil
	}
	//add to Borrow_records and update BOOKS.status
	DDL := time.Now().AddDate(0, 0, 14).Format("2006-01-02 15:04:05") 
	borrow = fmt.Sprintf("INSERT into Borrow_records(ISBN, accID, DDL, extend, status) values('%s', '%s', '%s', 0, 0)", ISBN, usrID, DDL)
	_, err = lib.db.Exec(borrow)
	if err != nil {
		fmt.Println("Invalid operation!")
		panic(err)
	}
	update_status := fmt.Sprintf("update BOOKS set status = 0 where ISBN = '%s'", ISBN)
	_, err = lib.db.Exec(update_status)
	if err != nil {
		fmt.Println("Invalid operation!")
		panic(err)
	}
	fmt.Println("Book borrowed successfully.")
	return nil
}

// Query_Borrow_History gets the borrow history of a student 
// It also tells whether the book is now returned according to Borrow_records.status
func Query_Borrow_History(lib *Library, usrID string) error {
	var cnt int
	Query_num := fmt.Sprintf("SELECT count(*) from Borrow_records where accID = '%s'", usrID)
	num_result, err := lib.db.Query(Query_num)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	num_result.Next()
	num_result.Scan(&cnt)
	if cnt==0 {
		fmt.Println("No borrow history for this account")
		return nil
	}
	QueryHistory := fmt.Sprintf("SELECT title, author, Borrow_records.ISBN, Borrow_records.status from Borrow_records, BOOKS where accID='%s' and BOOKS.ISBN = Borrow_records.ISBN", usrID)
	History, err := lib.db.Query(QueryHistory)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	var title, author, ISBN, return_status string
	var status int
	for History.Next() {
		return_status = "Not returned"
		History.Scan(&title, &author, &ISBN, &status)
		if status==1 {
			return_status = "Returned"
		}
		fmt.Printf("	Title: %s, author: %s, ISBN: %s, %s\n", title, author, ISBN, return_status)
	}
	return nil
}

//NotReturned query books a student has borrowed but not returned.
func NotReturned(lib *Library, stuID string) error {
	var cnt int
	Query_num := fmt.Sprintf("SELECT count(*) from Borrow_records, BOOKS where accID = '%s' and Borrow_records.ISBN = BOOKS.ISBN and BOOKS.status=0 and Borrow_records.status <> 1", stuID)
	num_result, err := lib.db.Query(Query_num)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	num_result.Next()
	num_result.Scan(&cnt)
	if cnt==0 {
		fmt.Println("No books not returned for this account.")
		return nil
	}
	Query_not_returned := fmt.Sprintf("SELECT title, author, BOOKS.ISBN from Borrow_records, BOOKS where Borrow_records.accID='%s' and Borrow_records.ISBN = BOOKS.ISBN and BOOKS.status=0", stuID)
	not_returned, err := lib.db.Query(Query_not_returned)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	var title, author, ISBN string
	for not_returned.Next() {
		not_returned.Scan(&title, &author, &ISBN)
		fmt.Printf("	Title: %s, author: %s, ISBN: %s, Not returned\n", title, author, ISBN)
	}
	return nil
}

//CheckDDL book should be borrowed and not returned yet.
func CheckDDL(lib *Library, ISBN string) error {
	var cnt int
	Query_num := fmt.Sprintf("SELECT count(*) from Borrow_records where ISBN = '%s' and status <> 1", ISBN)
	num_result, err := lib.db.Query(Query_num)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	num_result.Next()
	num_result.Scan(&cnt)
	if cnt==0 {
		fmt.Println("Book not borrowed yet")
		return nil
	}
	Check := fmt.Sprintf("SELECT DDL from Borrow_records where ISBN = '%s' and status <> 1", ISBN)
	DDL, err := lib.db.Query(Check)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	DDL.Next()
	var Date string
	DDL.Scan(&Date)
	fmt.Printf("Returning DDL for this book is %s\n", Date)
	return nil
}

//Extend DDL for 7 days each time
//overdue books can't extend
func Extend_ddl(lib *Library, ISBN, usrID string) error {
	var cnt int
	Date := time.Now().Format("2006-01-02 15:04:05")
	Query_num := fmt.Sprintf("SELECT count(*) from Borrow_records where ISBN = '%s' and status <> 1 and accID = '%s' and DDL > '%s'", ISBN, usrID, Date)
	num_result, err := lib.db.Query(Query_num)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	num_result.Next()
	num_result.Scan(&cnt)
	if cnt == 0 {
		fmt.Println("Book not borrowed by this account or book already overdue.")
		return nil
	}
	Check := fmt.Sprintf("SELECT extend from Borrow_records where ISBN = '%s' and status = 0", ISBN)
	DDL, err := lib.db.Query(Check)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	DDL.Next()
	var extend int
	DDL.Scan(&extend)
	if extend > 2 {
		fmt.Println("Extened too many times, operation failed.")
		return nil
	}
	update_ddl := fmt.Sprintf("UPDATE Borrow_records SET DDL = date_add(DDL, interval 7 day), extend = extend + 1 where ISBN = '%s' and status = 0", ISBN)
	_, err = lib.db.Exec(update_ddl)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	fmt.Println("Deadline extended for 7 days.")
	return nil
}

//Check books overdue for one student
func OverDue(lib *Library, stuID string) error {
	Date := time.Now().Format("2006-01-02 15:04:05")
	var cnt int
	Query_num := fmt.Sprintf("SELECT overdue_records from ACCOUNTS where accID = '%s'", stuID)
	num_result, err := lib.db.Query(Query_num)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	num_result.Next()
	num_result.Scan(&cnt)
	if cnt==0 {
		fmt.Println("No book overdue for this account.")
		return nil
	}
	check_due := fmt.Sprintf("SELECT title, BOOKS.author, BOOKS.ISBN from Borrow_records, BOOKS where DDL < '%s' and BOOKS.ISBN = Borrow_records.ISBN and Borrow_records.accID = '%s' and Borrow_records.status <> 1", Date, stuID)
	out_of_due, err := lib.db.Query(check_due)
	for out_of_due.Next() {
		var title, author, ISBN string
		err = out_of_due.Scan(&title, &author, &ISBN)
		if err != nil {
			fmt.Println("Failed to check overdue!")
			panic(err)
		}
		fmt.Printf("	Title: %s, author: %s, ISBN: %s\n", title, author, ISBN)
	}
	return nil
}

//ReturnBooks and if the book is overdue, returning it  will also update overdue_records for your account
//every time overdue_records changes, we check if we need to update thepriority of this account
//the book you return have to be borrowed by yourself using the same account
func ReturnBooks(lib *Library, ISBN , stuID string, priority *int) error {
	var cnt int
	Query_num := fmt.Sprintf("SELECT count(*) from Borrow_records where ISBN = '%s' and accID = '%s' and status <> 1", ISBN, stuID)
	num_result, err := lib.db.Query(Query_num)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	num_result.Next()
	num_result.Scan(&cnt)
	if cnt == 0 {
		fmt.Println("No borrow record match.")
		return nil
	}
	Confirm_borrow := fmt.Sprintf("SELECT status from Borrow_records where ISBN = '%s' and accID = '%s' and status <> 1", ISBN, stuID)
	confirm, err := lib.db.Query(Confirm_borrow)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	confirm.Next()
	var status int
	confirm.Scan(&status)
	Return1 := fmt.Sprintf("UPDATE Borrow_records SET status = 1 where ISBN = '%s'", ISBN)
	_, err = lib.db.Exec(Return1)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	Return2 := fmt.Sprintf("UPDATE BOOKS SET status = 1 where ISBN = '%s'", ISBN)
	_, err = lib.db.Exec(Return2)
	if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
	}
	if status == -1 {
		Return3 := fmt.Sprintf("UPDATE ACCOUNTS SET overdue_records = overdue_records - 1 where accID = '%s'", stuID)
		_, err = lib.db.Exec(Return3)
		if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
		}
		_, err = lib.db.Exec("UPDATE ACCOUNTS SET priority = 1 WHERE overdue_records < 4 and priority = -1")
		if err != nil {
			fmt.Println("Failed to update account status!")
			panic(err)
		}
		query_priority := fmt.Sprintf("SELECT priority from ACCOUNTS where accID = '%s'", stuID)
		prio, err := lib.db.Query(query_priority)
		if err != nil {
			fmt.Println("Invalid operation!")
			panic(err)
		}
		prio.Next()
		prio.Scan(priority)
	}
	return nil
	
}


func main() {
	
	fmt.Println("Welcome to the Library Management System!")
	var lib *Library = new(Library)
	lib.db = nil
	ConnectDB(lib)
	//Check DB status, should have 4 tables in DB, if no tables, intialize DB, if DB status is wrong, quit DB.
	countTables := fmt.Sprintf("SELECT count(TABLE_NAME) FROM information_schema.TABLES WHERE TABLE_SCHEMA='%s'", DBName)
	Tablenum, err := lib.db.Query(countTables)
	if err != nil {
			fmt.Println("Database check failed!")
			panic(err)
		}
	Tablenum.Next()
	var Tablecnt int
	err = Tablenum.Scan(&Tablecnt)
	if err != nil {
		fmt.Println("Database check failed!")
		panic(err)
	}
	
	if Tablecnt==0 {
		var pwd string
		fmt.Println("Initialize library system...")
		fmt.Println("Please set a password for administrator account (Input a six-digit integer):")
		fmt.Scanln(&pwd)
		CreateTables(lib)
		InitializeDB(lib, pwd)
	} else if Tablecnt!=4 {
		fmt.Println("Errors in library DB, contact administrator.")
	}
	//Update overdue_records for each account
	AccountUpdate(lib)
	//Log in interface and Menu part
	var choice int = 0;
	for choice != 2 {
		fmt.Println("--1.Log in with an account.")
		fmt.Println("--2.Quit system.")
		fmt.Scanln(&choice)
		if choice == 1 {
			var usrID string
			var usrPWD string
			fmt.Println("Input your account ID:")
			fmt.Scanln(&usrID)
			fmt.Println("Input your password:")
			fmt.Scanln(&usrPWD)
			var accountnum int
			Query_num := fmt.Sprintf("SELECT count(*) from ACCOUNTS where accID='%s' and pwd = '%s'", usrID, usrPWD)
			num_result, err := lib.db.Query(Query_num)
			if err != nil {
					fmt.Println("Log in failed!")
					panic(err)
			}
			num_result.Next()
			num_result.Scan(&accountnum)
			if accountnum != 1 {
				fmt.Println("Incorrect password or account does not exist.")
				continue
			}
			loginCmd := fmt.Sprintf("SELECT priority from ACCOUNTS where accID='%s' and pwd = '%s'", usrID, usrPWD)
			login, err := lib.db.Query(loginCmd)
			if err != nil {
				fmt.Println("Log in failed!")
				panic(err)
			}
			login.Next()
			var priority int
			login.Scan(&priority)
			if priority == 0 {
				//admin account, Menu
				fmt.Println("Welcome, administrator!")
				choice = 0
				for choice != 9 {
					fmt.Println("--1.Add a book to library.")
					fmt.Println("--2.Remove a book from library with explanation.")
					fmt.Println("--3.Add student account.")
					fmt.Println("--4.Query a book from library.")
					fmt.Println("--5.Query borrow history of a student account.")
					fmt.Println("--6.Query books a student has borrowed and not returned yet.")
					fmt.Println("--7.Check the deadline of returning a borrowed book.")
					fmt.Println("--8.Check if a student has any overdue books.")
					fmt.Println("--9.Log out.")
					fmt.Scanln(&choice)
					switch choice {
						case 1: {var title, author, ISBN string
							 fmt.Println("Input title")
							 fmt.Scanln(&title)
							 fmt.Println("Input author")
							 fmt.Scanln(&author)
							 fmt.Println("Input ISBN")
							 fmt.Scanln(&ISBN)
							 AddBook(lib, title, author, ISBN)}
						case 2: {var ISBN string
							 fmt.Println("Input ISBN you want to remove")
							 fmt.Scanln(&ISBN)
							 DeleteBook(lib, ISBN)}
						case 3: {var ID, pwd string
							 fmt.Println("Input student ID")
							 fmt.Scanln(&ID)
							 fmt.Println("Set password")
							 fmt.Scanln(&pwd)
							 AddStu(lib, ID, pwd)}
						case 4: {SearchBooks(lib)}
						case 5: {var ID string
							 fmt.Println("Input student ID you want to query")
							 fmt.Scanln(&ID)
							 Query_Borrow_History(lib, ID)}
						case 6: {var ID string
							 fmt.Println("Input student ID you want to query")
							 fmt.Scanln(&ID)
							 NotReturned(lib, ID)}
						case 7: {var ISBN string
							 fmt.Println("Input ISBN you want to check")
							 fmt.Scanln(&ISBN)
							 CheckDDL(lib, ISBN)}
						case 8: {var ID string
							 fmt.Println("Input student ID you want to query")
							 fmt.Scanln(&ID)
							 OverDue(lib, ID)}
					}
				}
			} else{
				//stu account, Menu
				fmt.Printf("Welcome, %s!\n", usrID)
				if priority == -1 {
					fmt.Println("You have more than 3 book overdue, access limited.")
				}
				choice = 0
				for choice != 9 {
					fmt.Println("--1.Borrow book.")
					fmt.Println("--2.Return book.")
					fmt.Println("--3.Query books you have borrowed but not returned yet.")
					fmt.Println("--4.Check the deadline of returning a borrowed book.")
					fmt.Println("--5.Extend your deadline of returning a book.")
					fmt.Println("--6.Check if you have any overdue books")
					fmt.Println("--7.Query a book from library.")
					fmt.Println("--8.Query borrow history of your account.")
					fmt.Println("--9.Log out.")
					fmt.Scanln(&choice)
					switch choice {
						case 1: {if priority == -1 {
							 fmt.Printf("You have more than 3 book overdue, unable to borrow new books.\nPlease return overdue books first.\n")
							 continue
							 }
							 var ISBN string
							 fmt.Println("Input ISBN you want to borrow")
							 fmt.Scanln(&ISBN)
							 Borrowbooks(lib, ISBN , usrID)}
						case 2: {var ISBN string
							 fmt.Println("Input ISBN you want to return")
							 fmt.Scanln(&ISBN)
							 ReturnBooks(lib, ISBN , usrID, &priority)}
						case 3: {NotReturned(lib, usrID)}
						case 4: {var ISBN string
							 fmt.Println("Input ISBN you want to check")
							 fmt.Scanln(&ISBN)
							 CheckDDL(lib, ISBN)}
						case 5: {var ISBN string
							 fmt.Println("Input ISBN you want to extend")
							 fmt.Scanln(&ISBN)
							 Extend_ddl(lib, ISBN, usrID)}
						case 6: {OverDue(lib, usrID)}
						case 7: {SearchBooks(lib)}
						case 8: {Query_Borrow_History(lib, usrID)}
					}
				}
			}
		}else {
			fmt.Println("Undefined command.")
		}
	}

}
