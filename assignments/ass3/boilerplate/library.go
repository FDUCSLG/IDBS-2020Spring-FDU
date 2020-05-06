package main

import (
	"fmt"

	// mysql connector
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	User     = "user"
	Password = "password"
	DBName   = "ass3"
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

// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error {
	return nil
}

// AddBook add a book into the library
func (lib *Library) AddBook(title, auther, ISBN string) error {
	return nil
}

// etc...

func main() {
	fmt.Println("Welcome to the Library Management System!")
}