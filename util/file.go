package util

import (
	"cloudStorage/config"
	"io"
	"mime/multipart"
	"os"
)

func GetFileHash(file multipart.File, filename string) (string, int64, error) {
	//save the file to the disk
	tmpfile, err := os.Create(config.FILE_TMP_PATH + filename)
	defer func() {
		//删除tmp文件
		// os.Remove(config.FILE_TMP_PATH + filename)
	}()

	if err != nil {
		return "", 0, err
	}
	defer tmpfile.Close()

	var filesize int64
	//copy the file data to the new file
	if filesize, err = io.Copy(tmpfile, file); err != nil {
		return "", 0, err
	}

	tmpfile.Seek(0, 0)            //seek to the beginning of the file
	filehash := FileSha1(tmpfile) //calculate the hash

	return filehash, filesize, nil
}

func SaveFile(file multipart.File, newfilepath string) error {
	newfile, err := os.Create(newfilepath)
	if err != nil {
		return err
	}
	defer newfile.Close()

	//copy the file data to the new file
	if _, err = io.Copy(newfile, file); err != nil {
		return err
	}
	return nil
}
