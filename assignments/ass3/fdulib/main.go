package fdulib

import (
	"database/sql"
	"fmt"

	// MySQL connector
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func (lib *Library) ConnectDB() error {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", lib.User, lib.Password, lib.DBName))
	if err != nil {
		return err
	}
	lib.db = db

	return nil
}

// Create tables in MySql
func (lib *Library) CreateTables() error {
	creates := []string{
		`CREATE TABLE IF NOT EXISTS admins(
			user VARCHAR(255) PRIMARY KEY,
			password VARCHAR(255)
		)`,
		`CREATE TABLE IF NOT EXISTS users(
			user VARCHAR(255) PRIMARY KEY,
			password VARCHAR(255),
			user_stat INT DEFAULT 0,
			overdue INT DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS books(
			book_id VARCHAR(255) PRIMARY KEY,
			book_stat INT DEFAULT 0,
			title VARCHAR(255) NOT NULL,
			edition INT DEFAULT 1,
			author VARCHAR(255),
			ISBN VARCHAR(20),
			added_date DATE NOT NULL,
			added_by_admin VARCHAR(255) NOT NULL,
			expire_days INT NOT NULL,
			removed_date DATE,
			removed_by_admin VARCHAR(255),
			removed_msg VARCHAR(1023),
			FOREIGN KEY (added_by_admin) REFERENCES admins(user),
			FOREIGN KEY (removed_by_admin) REFERENCES admins(user)
		)`,
		`CREATE TABLE IF NOT EXISTS records(
			record_id VARCHAR(255) PRIMARY KEY,
			user VARCHAR(255) NOT NULL,
			book_id VARCHAR(255) NOT NULL,
			lent_date DATE NOT NULL,
			expire_date DATE NOT NULL,
			extended INT DEFAULT 0,
			returned_date DATE,
			FOREIGN KEY (user) REFERENCES users(user),
			FOREIGN KEY (book_id) REFERENCES books(book_id)
		)`,
	}
	for _, s := range creates {
		_, err := lib.db.Exec(s)
		if err != nil {
			return err
		}
	}

	return nil
}

//1 Create an administrator account (by admin)
func (lib *Library) CreateAdmin(user, password string) error {
	s := fmt.Sprintf("INSERT INTO admins(user, password) VALUES ('%s', '%s')", user, password)
	_, err := lib.db.Exec(s)

	return err
}

//2 Create a user account (by admin)
func (lib *Library) CreateUser(user, password string) error {
	s := fmt.Sprintf("INSERT INTO users(user, password) VALUES ('%s', '%s')", user, password)
	_, err := lib.db.Exec(s)

	return err
}

//3 Add book into library (by admin)
func (lib *Library) AddBook(book *Book) error {
	s := "INSERT INTO books(book_id, title, edition, author, ISBN, added_date, added_by_admin, expire_days) VALUES "
	s += fmt.Sprintf("('%s', \"%s\", %d, \"%s\", '%s', '%s', '%s', %d)", book.Book_id, book.Title, book.Edition, book.Author.String, book.ISBN.String, book.Added_date, book.Added_by_admin, book.Expire_days)
	_, err := lib.db.Exec(s)

	return err
}

//4 Remove book by ID (by admin)
func (lib *Library) RemoveBookByID(book_id, removed_date, removed_by_admin, removed_msg string) error {
	s := fmt.Sprintf("UPDATE books SET book_stat = 2, removed_date = '%s', removed_by_admin = '%s', removed_msg = '%s' WHERE book_id = '%s' AND book_stat <> 2", removed_date, removed_by_admin, removed_msg, book_id)
	res, err := lib.db.Exec(s)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("wrong book ID or book already removed")
	}
	return nil
}

//5-6 Query book by specified field (by admin or user)
func (lib *Library) QueryBook(field string, content string) ([]Book, error) {
	var books []Book
	s := fmt.Sprintf("SELECT * FROM books WHERE %s LIKE '%%%s%%'", field, content)
	rows, err := lib.db.Queryx(s)
	if err != nil {
		return books, err
	}
	for rows.Next() {
		var b Book
		err = rows.Scan(&b.Book_id, &b.Book_stat, &b.Title, &b.Edition, &b.Author, &b.ISBN, &b.Added_date, &b.Added_by_admin, &b.Expire_days, &b.Removed_date, &b.Removed_by_admin, &b.Removed_msg)
		if err != nil {
			return books, err
		}
		books = append(books, b)
	}
	err = rows.Close()

	return books, nil
}

//7 Lend a book to user (by user)
func (lib *Library) LendBook(record_id, user, book_id, lent_date string) error {
	// Update and check user status
	user_stat, err := lib.UpdateUserStatus(user, lent_date)
	if err != nil {
		return err
	}
	if user_stat == 1 {
		return fmt.Errorf("user account suspended, unable to borrow book")
	}

	// Check and update book status
	res, err := lib.db.Exec(fmt.Sprintf("UPDATE books SET book_stat = 1 WHERE book_id = '%s' AND book_stat = 0", book_id))
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("book with id '%s' is currently unavailable", book_id)
	}

	// Insert corresponding record
	var expire_date string
	s := fmt.Sprintf("SELECT DATE_ADD('%s', INTERVAL expire_days DAY) FROM books WHERE book_id = '%s'", lent_date, book_id)
	err = lib.db.Get(&expire_date, s)
	if err != nil {
		return err
	}
	s = fmt.Sprintf("INSERT INTO records(record_id, user, book_id, lent_date, expire_date) VALUES ('%s', '%s', '%s', '%s', '%s')", record_id, user, book_id, lent_date, expire_date)
	_, err = lib.db.Exec(s)
	if err != nil {
		return err
	}

	// Success
	return nil
}

//8 Query user history (by admin or user)
func (lib *Library) QueryUserHistory(user string) ([]UserHistory, error) {
	var historys []UserHistory
	s := "SELECT record_id, lent_date, expire_date, returned_date, book_id, title, author FROM records NATURAL JOIN books "
	s += fmt.Sprintf("WHERE user = '%s' ORDER BY lent_date ASC", user)
	rows, err := lib.db.Queryx(s)
	if err != nil {
		return historys, err
	}
	for rows.Next() {
		var h UserHistory
		err = rows.Scan(&h.Record_id, &h.Lent_date, &h.Expire_date, &h.Returned_date, &h.Book_id, &h.Title, &h.Author)
		if err != nil {
			return historys, err
		}
		historys = append(historys, h)
	}
	err = rows.Close()

	return historys, nil
}

//9 Query pending books (by admin or user)
func (lib *Library) QueryPendingBook(user string) ([]UserHistory, error) {
	var historys []UserHistory
	s := "SELECT record_id, lent_date, expire_date, returned_date, book_id, title, author FROM records NATURAL JOIN books " +
		fmt.Sprintf("WHERE user = '%s' AND book_stat = 1 ORDER BY expire_date ASC, title ASC", user)
	rows, err := lib.db.Queryx(s)
	if err != nil {
		return historys, err
	}
	for rows.Next() {
		var h UserHistory
		err = rows.Scan(&h.Record_id, &h.Lent_date, &h.Expire_date, &h.Returned_date, &h.Book_id, &h.Title, &h.Author)
		if err != nil {
			return historys, err
		}
		historys = append(historys, h)
	}
	err = rows.Close()

	return historys, nil
}

//10 Extend deadline (by user)
func (lib *Library) ExtendDeadline(user string, record_id string) error {
	// Check user status
	var user_stat int
	err := lib.db.Get(&user_stat, fmt.Sprintf("SELECT user_stat FROM users WHERE user = '%s'", user))
	if err != nil {
		return err
	}
	if user_stat == 1 {
		return fmt.Errorf("user account suspended, unable to extend deadline")
	}

	// Check and extend deadline
	var expire_days int
	s := fmt.Sprintf("SELECT expire_days FROM records NATURAL JOIN books WHERE record_id = '%s'", record_id)
	err = lib.db.Get(&expire_days, s)
	if err != nil {
		return err
	}
	s = fmt.Sprintf("UPDATE records SET extended = extended + 1, expire_date = DATE_ADD(expire_date, INTERVAL %d DAY) ", expire_days)
	s += fmt.Sprintf("WHERE record_id = '%s' AND extended <= 2 AND returned_date IS NULL AND user = '%s'", record_id, user)
	res, err := lib.db.Exec(s)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("book unavailable for further extension")
	}

	return nil
}

//11 Query overdue books (by admin or user)
func (lib *Library) QueryOverdue(user string) ([]UserHistory, error) {
	var historys []UserHistory
	s := "SELECT record_id, lent_date, expire_date, returned_date, book_id, title, author " +
		 "FROM records NATURAL JOIN books WHERE expire_date < CURDATE() AND returned_date IS NULL ORDER BY expire_date ASC"
	rows, err := lib.db.Queryx(s)
	if err != nil {
		return historys, err
	}
	for rows.Next() {
		var h UserHistory
		err = rows.Scan(&h.Record_id, &h.Lent_date, &h.Expire_date, &h.Returned_date, &h.Book_id, &h.Title, &h.Author)
		if err != nil {
			return historys, err
		}
		historys = append(historys, h)
	}
	err = rows.Close()

	return historys, nil
}

//12 Return a book to library (by user)
func (lib *Library) ReturnBook(user, record_id, returned_date string) error {
	// Check and return book
	s := fmt.Sprintf("UPDATE records SET returned_date = '%s' WHERE record_id = '%s' AND returned_date IS NULL AND user = '%s'", returned_date, record_id, user)
	res, err := lib.db.Exec(s)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("book already returned or user did not borrow book")
	}

	// Update book status
	var book_id string
	s = fmt.Sprintf("SELECT book_id FROM books NATURAL JOIN records WHERE record_id = '%s'", record_id)
	err = lib.db.Get(&book_id, s)
	if err != nil {
		return err
	}
	s = fmt.Sprintf("UPDATE books SET book_stat = 0 WHERE book_id = '%s'", book_id)
	_, err = lib.db.Exec(s)
	if err != nil {
		return err
	}

	// Unsuspend user if applicable
	var row string
	s = fmt.Sprintf("SELECT expire_date FROM records WHERE record_id = '%s' AND expire_date < returned_date", record_id)
	err = lib.db.QueryRow(s).Scan(&row)
	switch err {
	case nil:
		// Returned an overdue book
		s = fmt.Sprintf("UPDATE users SET overdue = overdue - 1 WHERE user = '%s'", user)
		_, err = lib.db.Exec(s)
		if err != nil {
			return err
		}
		s = fmt.Sprintf("UPDATE users SET user_stat = 0 WHERE user = '%s' AND user_stat = 1 AND overdue <= 3", user)
		_, err = lib.db.Exec(s)
		if err != nil {
			return err
		}
	case sql.ErrNoRows:
		// Returned an ordinary book: do nothing
	default:
		return err
	}

	// Success
	return nil
}

func (lib *Library) GetAdminPassword(user string) (pw string, err error) {
	s := fmt.Sprintf("SELECT password FROM admins WHERE user = '%s'", user)
	err = lib.db.Get(&pw, s)
	return
}

func (lib *Library) GetUserPassword(user string) (pw string, err error) {
	s := fmt.Sprintf("SELECT password FROM users WHERE user = '%s'", user)
	err = lib.db.Get(&pw, s)
	return
}

func (lib *Library) UpdateUserStatus(user string, date string) (user_stat int, err error) {
	var overdue int
	s := fmt.Sprintf("SELECT COUNT(*) FROM records WHERE user = '%s' AND returned_date IS NULL AND expire_date < '%s'", user, date)
	err = lib.db.Get(&overdue, s)
	if err != nil {
		return user_stat, err
	}
	if overdue > 3 {
		user_stat = 1
	} else {
		user_stat = 0
	}
	s = fmt.Sprintf("UPDATE users SET user_stat = %d, overdue = %d WHERE user = '%s'", user_stat, overdue, user)
	_, err = lib.db.Exec(s)
	if err != nil {
		return user_stat, err
	}
	return user_stat, nil
}
