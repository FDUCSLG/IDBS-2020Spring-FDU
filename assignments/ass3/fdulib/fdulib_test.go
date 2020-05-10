package fdulib

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"testing"
)

var lib = Library{User: "fan", Password: "", DBName: "ass3"}

var admins = []struct {
	user, password string
	}{
		{"admin1", "iloveadmin2"},
		{"admin2", "ilovemyself"},
		{"admin2", "ilovemyself"},
}
var users = []struct {
	user, password string
	}{
		{"user1", "iloveuser2"},
		{"user2", "ilovemyself"},
		{"user2", "ilovemyself"},
}
var books = []Book{
	{
		Book_id:        "QA76.9.D3 R237 2003-1",
		Book_stat:      0,
		Title:          "Database Management Systems",
		Edition:        3,
		Author:         sql.NullString{"Raghu Ramakrishnan, Johannes Gehrke", true},
		ISBN:           sql.NullString{"9780072465631", true},
		Added_date:     "2020-01-02",
		Added_by_admin: "admin1",
		Expire_days:    14,
	},
	{
		Book_id:        "QA76.9.D3 R237 2003-2",
		Book_stat:      0,
		Title:          "Database Management Systems",
		Edition:        3,
		Author:         sql.NullString{"Raghu Ramakrishnan, Johannes Gehrke", true},
		ISBN:           sql.NullString{"9780072465631", true},
		Added_date:     "2020-04-02",
		Added_by_admin: "admin1",
		Expire_days:    14,
	},
	{
		Book_id:        "QA76.9.D3 S5637 2002-1",
		Book_stat:      0,
		Title:          "Database System Concepts",
		Edition:        4,
		Author:         sql.NullString{"Abraham Silberschatz, Henry F. Korth, S. Sudarshan", true},
		ISBN:           sql.NullString{"9780072283631", true},
		Added_date:     "2020-04-02",
		Added_by_admin: "admin2",
		Expire_days:    14,
	},
	{
		Book_id:        "QA76.5.B795 2016-1",
		Book_stat:      0,
		Title:          "Computer Systems: A Programmer's Perspective",
		Edition:        3,
		Author:         sql.NullString{"Randal E. Bryant, David R. O'Hallaron", true},
		ISBN:           sql.NullString{"9780134092669", true},
		Added_date:     "2020-04-02",
		Added_by_admin: "admin2",
		Expire_days:    14,
	},
	{
		Book_id:        "QA76.5.B795 2016-2",
		Book_stat:      0,
		Title:          "Computer Systems: A Programmer's Perspective",
		Edition:        3,
		Author:         sql.NullString{"Randal E. Bryant, David R. O'Hallaron", true},
		ISBN:           sql.NullString{"9780134092669", true},
		Added_date:     "2020-04-02",
		Added_by_admin: "admin2",
		Expire_days:    14,
	},
	{
		Book_id:        "QA76.9.C643 P37 2014-1",
		Book_stat:      0,
		Title:          "Computer Organization and Design: The Hardware/Software Interface",
		Edition:        5,
		Author:         sql.NullString{"David A. Patterson, John L. Hennessy", true},
		ISBN:           sql.NullString{"9780124077263", true},
		Added_date:     "2020-04-02",
		Added_by_admin: "admin2",
		Expire_days:    14,
	},
}
var histories = []UserHistory{
	{
		Record_id:   "20200405-1",
		Lent_date:   "2020-04-05",
		Expire_date: "2020-04-19",
		Book_id:     "QA76.9.D3 R237 2003-1",
		Title:       "Database Management Systems",
		Author:      sql.NullString{"Raghu Ramakrishnan, Johannes Gehrke", true},
	},
	{
		Record_id:   "20200405-2",
		Lent_date:   "2020-04-05",
		Expire_date: "2020-04-19",
		Book_id:     "QA76.9.D3 R237 2003-2",
		Title:       "Database Management Systems",
		Author:      sql.NullString{"Raghu Ramakrishnan, Johannes Gehrke", true},
	},
	{
		Record_id:   "20200405-3",
		Lent_date:   "2020-04-05",
		Expire_date: "2020-04-19",
		Book_id:     "QA76.5.B795 2016-1",
		Title:       "Computer Systems: A Programmer's Perspective",
		Author:      sql.NullString{"Randal E. Bryant, David R. O'Hallaron", true},
	},
	{
		Record_id:   "20200405-4",
		Lent_date:   "2020-04-05",
		Expire_date: "2020-04-19",
		Book_id:     "QA76.9.C643 P37 2014-1",
		Title:       "Computer Organization and Design: The Hardware/Software Interface",
		Author:      sql.NullString{"David A. Patterson, John L. Hennessy", true},
	},
}

func checkQueryBooks(t *testing.T, gots []Book, expecteds ...*Book) {
	if len(expecteds) != len(gots) {
		t.Errorf("Expected %d result(s), but got %d result(s) instead.", len(expecteds), len(gots))
	}
	for _, expected := range expecteds {
		matched := false
		for _, got := range gots {
			if reflect.DeepEqual(*expected, got) {
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("\nExpected result unfound: %+v", *expected)
		}
	}
}

func checkQueryHistories(t *testing.T, gots []UserHistory, expecteds ...*UserHistory) {
	if len(expecteds) != len(gots) {
		t.Errorf("Expected %d result(s), but got %d result(s) instead.", len(expecteds), len(gots))
	}
	for i, expected := range expecteds {
		if !reflect.DeepEqual(*expected, gots[i]) {
			t.Errorf("\nExpected: %+v\nGot: %+v", *expected, gots[i])
		}
	}
}

func init() {
	inits := []string{
		"DROP DATABASE IF EXISTS ass3",
		"CREATE DATABASE ass3",
		"USE ass3",
	}
	err := lib.ConnectDB()
	if err != nil {
		panic(err)
	}
	for _, s := range inits {
		_, err := lib.db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
}

func TestLibrary_CreateTables(t *testing.T) {
	err := lib.CreateTables()
	if err != nil {
		t.Errorf("CreateTables() failed:")
		log.Fatal(err)
	}
}

func TestLibrary_CreateAdmin(t *testing.T) {
	for i, admin := range admins {
		err := lib.CreateAdmin(admin.user, admin.password)
		if me, ok := err.(*mysql.MySQLError); ok && me.Number == 1062 && i == 2 {
			// Error 1062: Duplicate entry 'admin2' for key 'admins.PRIMARY'
			t.Log(err)
		} else if err != nil {
			t.Errorf("CreateAdmin() failed: %s", err)
		}
	}
}

func TestLibrary_GetAdminPassword(t *testing.T) {
	for _, a := range admins {
		pw, err := lib.GetAdminPassword(a.user)
		if err != nil {
			t.Errorf("GetAdminPassword('%s') failed: %s", a.user, err)
		}
		if pw != a.password {
			t.Errorf("Expected password for admin '%s': '%s', got: '%s'", a.user, a.password, pw)
		}
	}
}

func TestLibrary_CreateUser(t *testing.T) {
	for i, user := range users {
		err := lib.CreateUser(user.user, user.password)
		if me, ok := err.(*mysql.MySQLError); ok && me.Number == 1062 && i == 2 {
			// Error 1062: Duplicate entry 'user2' for key 'users.PRIMARY'
			t.Log(err)
		} else if err != nil {
			t.Errorf("CreateUser() failed: %s", err)
		}
	}
}

func TestLibrary_GetUserPassword(t *testing.T) {
	for _, u := range users {
		pw, err := lib.GetUserPassword(u.user)
		if err != nil {
			t.Errorf("GetUserPassword('%s') failed: %s", u.user, err)
		}
		if pw != u.password {
			t.Errorf("Expected password for user '%s': '%s', got: '%s'", u.user, u.password, pw)
		}
	}
}

func TestLibrary_AddBook(t *testing.T) {
	for _, book := range books {
		err := lib.AddBook(&book)
		if err != nil {
			t.Errorf("AddBook() failed:")
			fmt.Print(err)
		}
	}
}

func TestLibrary_QueryBook(t *testing.T) {
	bs, err := lib.QueryBook("title", "database")
	if err != nil {
		t.Errorf(`QueryBook("title", "database") failed: %s`, err)
	}
	checkQueryBooks(t, bs, &books[0], &books[1], &books[2])
	bs, err = lib.QueryBook("titled", "atabase")
	if me, ok := err.(*mysql.MySQLError); ok && me.Number == 1054 {
		// Error 1054: Unknown column 'titled' in 'where clause'
		t.Log(err)
	} else if err != nil {
		t.Errorf(`QueryBook("titled", "atabase") failed: %s`, err)
	}
}

func TestLibrary_RemoveBookByID(t *testing.T) {
	removes := []Book{
		{
			Book_id:          "QA76.5.B795 2016-2",
			Book_stat:        2,
			Title:            "Computer Systems: A Programmer's Perspective",
			Edition:          3,
			Author:           sql.NullString{"Randal E. Bryant, David R. O'Hallaron", true},
			ISBN:             sql.NullString{"9780134092669", true},
			Added_date:       "2020-04-02",
			Added_by_admin:   "admin2",
			Expire_days:      14,
			Removed_date:     sql.NullString{"2020-05-03", true},
			Removed_by_admin: sql.NullString{"admin1", true},
			Removed_msg:      sql.NullString{"Donated to Hubei children.", true},
		},
	}
	for _, r := range removes {
		err := lib.RemoveBookByID(r.Book_id, r.Removed_date.String, r.Removed_by_admin.String , r.Removed_msg.String)
		if err != nil {
			t.Errorf("RemoveBookByID() failed: %s", err)
		}
		bs, err := lib.QueryBook("book_id", r.Book_id)
		checkQueryBooks(t, bs, &removes[0])
	}
	err := lib.RemoveBookByID("QA76.5.B795 2016-2", "2020-05-03", "admin1" , "Trying to remove a removed book.")
	if err != nil && err.Error() == "wrong book ID or book already removed" {
		t.Log(err)
	} else if err != nil {
		t.Errorf("RemoveBookByID() failed: %s", err)
	}
	err = lib.RemoveBookByID("Non-existent Book", "2020-05-03", "admin1" , "Trying to remove an invisible book.")
	if err != nil && err.Error() == "wrong book ID or book already removed" {
		t.Log(err)
	} else if err != nil {
		t.Errorf("RemoveBookByID() failed: %s", err)
	}
}

func TestLibrary_LendBook(t *testing.T) {
	lends := []struct {
		record_id string
		user string
		book_id string
		lent_date string
	}{
		{"20200405-1", "user1", "QA76.9.D3 R237 2003-1", "2020-04-05"},
		{"20200405-2", "user1", "QA76.9.D3 R237 2003-2", "2020-04-05"},
		{"20200405-3", "user1", "QA76.5.B795 2016-1", "2020-04-05"},
		{"20200405-4", "user1", "QA76.9.C643 P37 2014-1", "2020-04-05"},
		{"20200503-5", "user2", "QA76.9.D3 S5637 2002-1", "2020-05-03"},
	}
	for _, lend := range lends {
		err := lib.LendBook(lend.record_id, lend.user, lend.book_id, lend.lent_date)
		if err != nil {
			t.Errorf("LendBook('%s', '%s', '%s', '%s') failed: %s", lend.record_id, lend.user, lend.book_id, lend.lent_date, err)
		}
	}
	// Trying to borrow book using a suspended account
	err := lib.LendBook("20200503-6", "user1", "QA76.9.D3 S5637 2002-1", "2020-05-03")
	if err != nil && err.Error() == "user account suspended, unable to borrow book" {
		t.Log(err)
	} else if err != nil {
		t.Errorf(`LendBook() failed: %s`, err)
	}
	// Trying to borrow an unavailable book
	err = lib.LendBook("20200503-7", "user2", "QA76.9.C643 P37 2014-1", "2020-05-03")
	if err != nil && err.Error() == "book with id 'QA76.9.C643 P37 2014-1' is currently unavailable" {
		t.Log(err)
	} else if err != nil {
		t.Errorf(`LendBook() failed: %s`, err)
	}
}

func TestLibrary_QueryUserHistory(t *testing.T) {
	hs, err := lib.QueryUserHistory("user1")
	if err != nil {
		t.Errorf("QueryUserHistory('user1') failed: %s", err)
	}
	checkQueryHistories(t, hs, &histories[0], &histories[1], &histories[2], &histories[3])
}

func TestLibrary_QueryPendingBook(t *testing.T) {
	hs, err := lib.QueryPendingBook("user1")
	if err != nil {
		t.Errorf("QueryPendingBook('user1') failed: %s", err)
	}
	checkQueryHistories(t, hs, &histories[3], &histories[2], &histories[0], &histories[1])
}

func TestLibrary_ExtendDeadline(t *testing.T) {
	extends := []struct{
		user, record_id string
	}{
		{"user2", "20200503-5"},
		{"user2", "20200503-5"},
		{"user2", "20200503-5"},
	}
	for _, e := range extends {
		err := lib.ExtendDeadline(e.user, e.record_id)
		if err != nil {
			t.Errorf("ExtendDeadline('%s', '%s') failed: %s", e.user, e.record_id, err)
		}
	}
	// Extends for too many times
	err := lib.ExtendDeadline("user2", "20200503-5")
	if err != nil && err.Error() == "book unavailable for further extension" {
		t.Log(err)
	} else if err != nil {
		t.Errorf(`ExtendDeadline() failed: %s`, err)
	}
	// Suspended user tries to extend
	err = lib.ExtendDeadline("user1", "20200405-1")
	if err != nil && err.Error() == "user account suspended, unable to extend deadline" {
		t.Log(err)
	} else if err != nil {
		t.Errorf(`ExtendDeadline() failed: %s`, err)
	}
}

func TestLibrary_QueryOverdue(t *testing.T) {
	hs, err := lib.QueryOverdue("user1")
	if err != nil {
		t.Errorf("QueryOverdue('user1') failed: %s", err)
	}
	checkQueryHistories(t, hs, &histories[0], &histories[1], &histories[2], &histories[3])
}

func TestLibrary_ReturnBook(t *testing.T) {
	returns := []struct {
		user string
		record_id string
		returned_date string
	}{
		{"user1", "20200405-1", "2020-05-08"},
		{"user1", "20200405-2", "2020-05-08"},
		{"user1", "20200405-3", "2020-05-08"},
		{"user1", "20200405-4", "2020-05-08"},
	}
	for _, r := range returns {
		err := lib.ReturnBook(r.user, r.record_id, r.returned_date)
		if err != nil {
			t.Errorf("ReturnBook('%s', '%s', '%s') failed: %s", r.user, r.record_id, r.returned_date, err)
		}
	}
	// Trying to return a returned book
	err := lib.ReturnBook("user1", "20200405-1", "2020-05-08")
	if err != nil && err.Error() == "book already returned or user did not borrow book" {
		t.Log(err)
	} else if err != nil {
		t.Errorf(`ReturnBook() failed: %s`, err)
	}
	// Trying to return an unborrowed book
	err = lib.ReturnBook("user1", "20200405-9", "2020-05-08")
	if err != nil && err.Error() == "book already returned or user did not borrow book" {
		t.Log(err)
	} else if err != nil {
		t.Errorf(`ReturnBook() failed: %s`, err)
	}
}