package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var readdb *sql.DB

func init() {
	readdb, _ = sql.Open("mysql", "root:12345678@tcp(localhost:3307)/tbl_file?charset=utf8mb4")
	readdb.SetMaxOpenConns(1000)
	err := readdb.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func DBReadConn() *sql.DB {
	return readdb
}
