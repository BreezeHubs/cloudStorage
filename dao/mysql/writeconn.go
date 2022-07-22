package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var writedb *sql.DB

func init() {
	writedb, _ = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/tbl_file?charset=utf8mb4")
	writedb.SetMaxOpenConns(1000)
	err := writedb.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func DBWriteConn() *sql.DB {
	return writedb
}
