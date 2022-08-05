package dao

import (
	mydb "cloudStorage/dao/mysql"
	"database/sql"

	"github.com/pkg/errors"
)

func OnFileUploadFinished(userid int, filehash string, filename string, filesize int64, fileaddr string, isFastUpload bool) error {
	db := mydb.DBWriteConn()
	txdb, err := db.Begin()
	if err != nil {
		return err
	}

	if !isFastUpload {
		if _, err := mydb.Exec(
			db,
			"insert ignore into tbl_file(`file_hash`,`file_name`,`file_size`,`file_addr`,`status`) values(?,?,?,?,1)",
			filehash, filename, filesize, fileaddr,
		); err != nil {
			txdb.Rollback()
			return errors.Wrap(err, "OnFileUploadFinished mydb.Exec")
		}
		//  else if rows <= 0 {
		// 	txdb.Rollback()
		// 	return errors.New("OnFileUploadFinished failed")
		// }
	}

	if rows, err := mydb.Exec(
		mydb.DBWriteConn(),
		"insert ignore into tbl_user_file(`user_id`,`file_hash`,`file_name`,`file_size`,`status`) values(?,?,?,?,1)",
		userid, filehash, filename, filesize,
	); err != nil {
		return errors.Wrap(err, "OnUserFileUploadFinished mydb.Exec")
	} else if rows <= 0 {
		return errors.New("OnUserFileUploadFinished failed")
	}

	txdb.Commit()
	return nil
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
	CreateAt sql.NullString
	UpdateAt sql.NullString
}

func GetFileMeta(filehash string) (*TableFile, error) {
	tfile := TableFile{}
	if err := mydb.Read(
		mydb.DBReadConn(),
		"select file_hash,file_name,file_size,file_addr,create_at,update_at from tbl_file where file_hash=? and status=1 limit 1",
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

func GetFileList() ([]*TableFile, error) {
	tfiles := make([]*TableFile, 0)
	if err := mydb.ReadAny(
		mydb.DBReadConn(),
		"select file_hash,file_name,file_size,file_addr,create_at,update_at from tbl_file where status=1",
		nil,
		func(rows *sql.Rows) error {
			tfile := TableFile{}
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

func DeleteFile(filehash string) error {
	if _, err := mydb.Exec(
		mydb.DBWriteConn(),
		"update tbl_file set status=0 where file_hash=?",
		filehash,
	); err != nil {
		return errors.Wrap(err, "DeleteFile s.Exec")
	}
	return nil
}

func UpdateFile(filehash string, f *TableFile) error {
	if _, err := mydb.Exec(
		mydb.DBWriteConn(),
		"update tbl_file set file_hash=?,file_name=?,file_size=?,file_addr=? where file_hash=?",
		f.FileHash, f.FileName, f.FileSize, f.FileAddr, filehash,
	); err != nil {
		return errors.Wrap(err, "UpdateFile s.Exec")
	}
	return nil
}

func UpdateFileLocation(filehash string, fileaddr string) error {
	if _, err := mydb.Exec(
		mydb.DBWriteConn(),
		"update tbl_file set file_addr=? where file_hash=?",
		fileaddr, filehash,
	); err != nil {
		return errors.Wrap(err, "UpdateFileLocation s.Exec")
	}
	return nil
}
