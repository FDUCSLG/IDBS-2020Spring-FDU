package fdulib

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Library struct {
	db *sqlx.DB
	//Optimally, each client should have a separate MySQL account e.g.'HGX' or 'H3'
	User string
	Password string
	DBName string
}

type Book struct {
	Book_id          string
	Book_stat        int
	Title            string
	Edition          int
	Author           sql.NullString
	ISBN             sql.NullString
	Added_date       string
	Added_by_admin   string
	Expire_days      int
	Removed_date     sql.NullString
	Removed_by_admin sql.NullString
	Removed_msg      sql.NullString
}

type UserHistory struct {
	Record_id     string
	Lent_date     string
	Expire_date   string
	Returned_date sql.NullString
	Book_id       string
	Title         string
	Author        sql.NullString
}