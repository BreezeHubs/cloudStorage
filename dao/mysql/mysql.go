package mysql

import (
	"database/sql"
)

func Exec(db *sql.DB, sql string, args ...any) (int64, error) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func Read(db *sql.DB, sql string, args []any, put ...any) error {
	s, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer s.Close()
	return s.QueryRow(args...).Scan(put...)
}

func ReadAny(db *sql.DB, sql string, args []any, f func(rows *sql.Rows) error) error {
	s, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer s.Close()

	rows, err := s.Query(args...)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err := f(rows); err != nil {
			return err
		}
	}
	return nil
}
