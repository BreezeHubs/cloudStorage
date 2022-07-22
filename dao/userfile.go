package dao

import (
	mydb "cloudStorage/dao/mysql"
	"database/sql"

	"github.com/pkg/errors"
)

func OnUserFileUploadFinished(userid int, filehash string, filename string, filesize int64) error {
	if rows, err := mydb.Exec(
		mydb.DBWriteConn(),
		"insert ignore into tbl_user_file(`user_id`,`file_hash`,`file_name`,`file_size`,`status`) values(?,?,?,?,1)",
		userid, filehash, filename, filesize,
	); err != nil {
		return errors.Wrap(err, "OnUserFileUploadFinished mydb.Exec")
	} else if rows <= 0 {
		return errors.New("OnUserFileUploadFinished failed")
	}
	return nil
}

type TableUserFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
	CreateAt sql.NullString
	UpdateAt sql.NullString
}

func GetUserFileMeta(filehash string) (*TableUserFile, error) {
	tfile := TableUserFile{}
	if err := mydb.Read(
		mydb.DBReadConn(),
		`select 
			tbl_user_file.file_hash,tbl_user_file.file_name,tbl_user_file.file_size,tbl_file.file_addr,tbl_user_file.create_at,tbl_user_file.update_at
		from tbl_user_file 
		LEFT JOIN tbl_file 
			ON tbl_file.file_hash = tbl_user_file.file_hash 
		where tbl_user_file.file_hash=? and tbl_user_file.status=1 
		limit 1`,
		[]any{filehash},
		&tfile.FileHash, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr, &tfile.CreateAt, &tfile.UpdateAt,
	); err != nil {
		if err == sql.ErrNoRows {
			// not found
			return nil, nil
		}
		return nil, errors.Wrap(err, "GetFileMeta s.QueryRow")
	}
	return &tfile, nil
}

func GetUserFileList() ([]*TableUserFile, error) {
	tfiles := make([]*TableUserFile, 0)
	if err := mydb.ReadAny(
		mydb.DBReadConn(),
		`select 
		tbl_user_file.file_hash,tbl_user_file.file_name,tbl_user_file.file_size,tbl_file.file_addr,tbl_user_file.create_at,tbl_user_file.update_at
		from tbl_user_file 
		LEFT JOIN tbl_file 
			ON tbl_file.file_hash = tbl_user_file.file_hash 
		where tbl_user_file.status=1`,
		nil,
		func(rows *sql.Rows) error {
			tfile := TableUserFile{}
			if err := rows.Scan(&tfile.FileHash, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr, &tfile.CreateAt, &tfile.UpdateAt); err != nil {
				return errors.Wrap(err, "GetFileList mydb.ReadAny() rows.Scan")
			}
			tfiles = append(tfiles, &tfile)
			return nil
		},
	); err != nil {
		return nil, errors.Wrap(err, "GetFileList mydb.ReadAny()")
	}
	return tfiles, nil
}

func DeleteUserFile(filehash string) error {
	if _, err := mydb.Exec(
		mydb.DBWriteConn(),
		"update tbl_user_file set status=0 where file_hash=?",
		filehash,
	); err != nil {
		return errors.Wrap(err, "DeleteFile s.Exec")
	}
	return nil
}

func UpdateUserFile(filehash string, fileName string) error {
	if _, err := mydb.Exec(
		mydb.DBWriteConn(),
		"update tbl_user_file set file_name=? where file_hash=?",
		fileName, filehash,
	); err != nil {
		return errors.Wrap(err, "DeleteFile s.Exec")
	}
	return nil
}
