package main

import (
	"bufio"
	"database/sql"
	"fdulib"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"time"
)

const (
	User		= "fan"
	Password	= ""
	DBName		= "ass3"
)

var funcMap = []string {
	"Quit",
	"Create an admin account",
	"Create a user account",
	"Add a book into library",
	"Remove a book from library",
	"Query books by a specified field",
	"Query books by a specified field",
	"Borrow a book",
	"Query user history",
	"Query pending books",
	"Extend the deadline of a book",
	"Query overdue books",
	"Return a book",
}

var adminFunc = []int {1, 2, 3, 4, 5, 8, 9, 10, 11, 0}
var userFunc = []int {6, 7, 8, 9, 10, 11, 12, 0}

var recordIDCounter int = 0

//Scans only one element
func scanln(a interface{}) {
	_, err := fmt.Scanln(a)
	if err != nil {
		fmt.Printf("fmt.Scanln error: %s\n", err)
	}
}

func readString(s *string) {
	var err error
	reader := bufio.NewReader(os.Stdin)
	*s, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("reader.ReadString error: %s\n", err)
	}
}

func readNullString(s *sql.NullString) {
	var err error
	reader := bufio.NewReader(os.Stdin)
	s.String, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("reader.ReadString error: %s\n", err)
	}
	s.Valid = s.String != ""
}

func showFunc(userType string) {
	if userType == "admin" {
		for _, i := range adminFunc {
			fmt.Printf("%d: %s\n", i, funcMap[i])
		}
	} else if userType == "user" {
		for _, i := range userFunc {
			fmt.Printf("%d: %s\n", i, funcMap[i])
		}
	}
}

func resolveFunc(lib *fdulib.Library, user string, id int) {
	var args [5]string
	var err error
	if id == 1 {
		//CreateAdmin(user, password string) error
		fmt.Print("Username:")
		scanln(&args[0])
		fmt.Print("Password:")
		scanln(&args[1])

		err = lib.CreateAdmin(args[0], args[1])
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	} else if id == 2 {
		//CreateUser(user, password string) error
		fmt.Print("Username:")
		scanln(&args[0])
		fmt.Print("Password:")
		scanln(&args[1])

		err = lib.CreateUser(args[0], args[1])
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	} else if id == 3 {
		//AddBook(book *Book) error
		var b fdulib.Book
		fmt.Print("ID:")
		readString(&b.Book_id)
		fmt.Print("Title:")
		readString(&b.Title)
		fmt.Print("Edition:")
		scanln(&b.Edition)
		fmt.Print("Author:")
		readNullString(&b.Author)
		fmt.Print("ISBN:")
		scanln(&b.ISBN)
		b.Added_date = time.Now().Format("2006-01-02")
		b.Added_by_admin = user
		fmt.Print("Longest kept days:")
		scanln(&b.Expire_days)

		err = lib.AddBook(&b)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	} else if id == 4 {
		//RemoveBookByID(book_id, removed_date, removed_by_admin, removed_msg string) error
		fmt.Print("Book ID:")
		scanln(&args[0])
		args[1] = time.Now().Format("2006-01-02")
		args[2] = user
		fmt.Print("Message upon removal:")
		scanln(&args[3])

		err = lib.RemoveBookByID(args[0], args[1], args[2], args[3])
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	} else if id == 5 {
		//QueryBook(field string, content string) ([]Book, error)
		//Admin mode
		fmt.Println("Fields available to search in: book_id, title, author, ISBN, added_date, added_by_admin, removed_date, removed_by_admin, removed_msg")
		fmt.Print("Field:")
		scanln(&args[0])
		fmt.Print("Content:")
		scanln(&args[1])

		bs, err := lib.QueryBook(args[0], args[1])
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		for _, b := range bs {
			fmt.Printf("ID:%s, ", b.Book_id)
			fmt.Printf("Status:%d, ", b.Book_stat)
			fmt.Printf("Title:%s, ", b.Title)
			fmt.Printf("Edition:%d, ", b.Edition)
			fmt.Printf("Author:%s, ", b.Author.String)
			fmt.Printf("ISBN:%s, ", b.ISBN.String)
			fmt.Printf("Added:%s, ", b.Added_date)
			fmt.Printf("Added by admin:%s, ", b.Added_by_admin)
			fmt.Printf("Longest kept days:%d, ", b.Expire_days)
			fmt.Printf("Removed:%s, ", b.Removed_date.String)
			fmt.Printf("Removed by admin:%s, ", b.Removed_by_admin.String)
			fmt.Printf("Message upon removal:%s\n", b.Removed_msg.String)
		}
	} else if id == 6 {
		//QueryBook(field string, content string) ([]Book, error)
		//User mode
		fmt.Println("Fields available to search in: book_id, title, author, ISBN")
		fmt.Print("Field:")
		scanln(&args[0])
		fmt.Print("Content:")
		scanln(&args[1])

		bs, err := lib.QueryBook(args[0], args[1])
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		for _, b := range bs {
			// Hide removed books
			if b.Book_stat == 2 {
				continue
			}
			fmt.Printf("ID:%s, ", b.Book_id)
			fmt.Printf("Status:%d, ", b.Book_stat)
			fmt.Printf("Title:%s, ", b.Title)
			fmt.Printf("Edition:%d, ", b.Edition)
			fmt.Printf("Author:%s, ", b.Author.String)
			fmt.Printf("ISBN:%s, ", b.ISBN.String)
			fmt.Printf("Added:%s, ", b.Added_date)
			fmt.Printf("Longest kept days:%d\n", b.Expire_days)
		}
	} else if id == 7 {
		//LendBook(record_id string, user string, book_id string, lent_date string) error
		args[0] = strconv.Itoa(recordIDCounter)
		recordIDCounter++
		args[1] = user
		fmt.Print("Book ID:")
		scanln(&args[2])
		args[3] = time.Now().Format("2006-01-02")

		err = lib.LendBook(args[0], args[1], args[2], args[3])
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	} else if id == 8 {
		//QueryUserHistory(user string) ([]UserHistory, error)
		hs, err := lib.QueryUserHistory(user)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		for _, h := range hs {
			fmt.Printf("Record ID:%s, ", h.Record_id)
			fmt.Printf("Lent date:%s, ", h.Lent_date)
			fmt.Printf("Expire date:%s, ", h.Expire_date)
			fmt.Printf("Returned date:%s, ", h.Returned_date.String)
			fmt.Printf("Book ID:%s, ", h.Book_id)
			fmt.Printf("Title:%s, ", h.Title)
			fmt.Printf("Author:%s\n", h.Author.String)
		}
	} else if id == 9 {
		//QueryPendingBook(user string) ([]UserHistory, error)
		hs, err := lib.QueryPendingBook(user)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		for _, h := range hs {
			fmt.Printf("Record ID:%s, ", h.Record_id)
			fmt.Printf("Lent date:%s, ", h.Lent_date)
			fmt.Printf("Expire date:%s, ", h.Expire_date)
			fmt.Printf("Returned date:%s, ", h.Returned_date.String)
			fmt.Printf("Book ID:%s, ", h.Book_id)
			fmt.Printf("Title:%s, ", h.Title)
			fmt.Printf("Author:%s\n", h.Author.String)
		}
	} else if id == 10 {
		//ExtendDeadline(user string, record_id string) error
		args[0] = user
		fmt.Print("Record ID:")
		scanln(&args[1])

		err = lib.ExtendDeadline(args[0], args[1])
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	} else if id == 11 {
		//QueryOverdue(user string) ([]UserHistory, error)
		hs, err := lib.QueryOverdue(user)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		for _, h := range hs {
			fmt.Printf("Record ID:%s, ", h.Record_id)
			fmt.Printf("Lent date:%s, ", h.Lent_date)
			fmt.Printf("Expire date:%s, ", h.Expire_date)
			fmt.Printf("Returned date:%s, ", h.Returned_date.String)
			fmt.Printf("Book ID:%s, ", h.Book_id)
			fmt.Printf("Title:%s, ", h.Title)
			fmt.Printf("Author:%s\n", h.Author.String)
		}
	} else if id == 12 {
		//ReturnBook(user string, record_id string, returned_date string) error
		args[0] = user
		fmt.Print("Record ID:")
		scanln(&args[1])
		args[2] = time.Now().Format("2006-01-02")

		err = lib.ReturnBook(args[0], args[1], args[2])
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}

func adminMode(lib *fdulib.Library) error {
	var user, pw string
	fmt.Print("Admin username:")
	scanln(&user)
	fmt.Print("Password:")
	scanln(&pw)
	correctPW, err := lib.GetAdminPassword(user)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	if pw != correctPW {
		return fmt.Errorf("incorrect username or password for admin")
	}

	fmt.Printf("Welcome, Administrator %s\n", user)
	showFunc("admin")
	for {
		var id int
		valid := false
		fmt.Print("> ")
		_, err := fmt.Scanln(&id)
		//Check for non-integer inputs
		if err != nil {
			showFunc("admin")
			continue
		}
		//Check if id is valid
		for _, v := range adminFunc {
			if v == id {
				valid = true
				break
			}
		}
		if !valid {
			showFunc("admin")
			continue
		}
		//Normal function handling
		if id == 0 {
			fmt.Println("Bye!")
			break
		} else if id == 8 || id == 9 || id == 11 {
			//In these functions the admin checks the data of a user other than himself
			var user string
			fmt.Print("User:")
			scanln(&user)
			resolveFunc(lib, user, id)
		} else {
			resolveFunc(lib, user, id)
		}
	}
	return nil
}

func userMode(lib *fdulib.Library) error {
	var user, pw string
	fmt.Print("Username:")
	scanln(&user)
	fmt.Print("Password:")
	scanln(&pw)
	correctPW, err := lib.GetUserPassword(user)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	if pw != correctPW {
		return fmt.Errorf("incorrect username or password for user")
	}

	fmt.Printf("Welcome, User %s\n", user)
	showFunc("user")
	for {
		var id int
		valid := false
		fmt.Print("> ")
		_, err := fmt.Scanln(&id)
		//Check for non-integer inputs
		if err != nil {
			showFunc("user")
			continue
		}
		//Check if id is valid
		for _, v := range userFunc {
			if v == id {
				valid = true
				break
			}
		}
		if !valid {
			showFunc("user")
			continue
		}
		//Normal function handling
		if id == 0 {
			fmt.Println("Bye!")
			break
		} else {
			resolveFunc(lib, user, id)
		}
	}
	return nil
}

func main() {
	fmt.Println("Welcome to the FUDAN Library Management System!")

	lib := fdulib.Library{User: User, Password: Password, DBName: DBName}
	fmt.Println("Establishing connection with database...")
	err := lib.ConnectDB()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println("Initializing database tables...")
	err = lib.CreateTables()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println("Adding root user...")
	err = lib.CreateAdmin("root", "1905")
	if me, ok := err.(*mysql.MySQLError); ok && me.Number == 1062 {
		//Ignore Error 1062: Duplicate entry 'root' for key 'admins.PRIMARY'
	} else if err != nil  {
		fmt.Printf("%s\n", err)
	}

	for {
		var t int
		fmt.Printf("0: Admin Mode\n1: User Mode\n2: Quit\n")
		fmt.Printf("> ")
		_, err = fmt.Scan(&t)
		if err != nil {
			fmt.Printf("fmt.Scan error: %s\n", err)
		}
		if t == 2 {
			fmt.Println("Bye!")
			break
		} else if t == 0 {
			err = adminMode(&lib)
		} else if t == 1 {
			err = userMode(&lib)
		}
		if err != nil {
			fmt.Printf("%s\n", err)
		} else {
			break
		}
	}

}