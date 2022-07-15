package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"cloudStorage/meta"
	"cloudStorage/util"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//return the upload form
		b, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "upload err")
			return
		}
		io.WriteString(w, string(b))
	} else if r.Method == "POST" {
		//handle the upload
		f, fh, err := r.FormFile("file")
		if err != nil {
			io.WriteString(w, "upload err")
			return
		}
		defer f.Close()

		//create a file meta
		fileMeta := meta.FileMeta{
			FileName: fh.Filename,
			Location: "./tmp/" + fh.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		//save the file to the disk
		newfile, err := os.Create(fileMeta.Location)
		if err != nil {
			io.WriteString(w, "create file err")
			return
		}
		defer newfile.Close()

		//copy the file data to the new file
		if fileMeta.FileSize, err = io.Copy(newfile, f); err != nil {
			io.WriteString(w, "failed to save file")
			return
		}

		newfile.Seek(0, 0)                         //seek to the beginning of the file
		fileMeta.FileSha1 = util.FileSha1(newfile) //calculate the sha1
		meta.UpdateFileMeta(fileMeta)              //update the file meta to map

		http.Redirect(w, r, "/file/upload/success", http.StatusFound) //redirect to the success page
	}
}

// UploadSucHandler handles the success page
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload success")
}

// GetFileMetaHandler handles the request of getting file meta
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(data)
}

// GetFileMetasHandler handles the request of getting file metas
func GetFileMetasHandler(w http.ResponseWriter, r *http.Request) {
	fileMetas := meta.GetFileMetas()
	data, err := json.Marshal(fileMetas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(b)
}

func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	op := r.Form.Get("op")
	filehash := r.Form.Get("filehash")
	newfilename := r.Form.Get("filename")
	fmt.Println(op, filehash, newfilename)

	if op != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(filehash)
	oldFileName := curFileMeta.FileName
	curFileMeta.FileName = newfilename
	curFileMeta.Location = "./tmp/" + newfilename
	meta.UpdateFileMeta(curFileMeta)

	//rename the file
	os.Rename("./tmp/"+oldFileName, "./tmp/"+newfilename)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func FileMetaDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")

	//delete the file
	fMeta := meta.GetFileMeta(filehash)

	meta.RemoveFileMeta(filehash)
	os.Remove(fMeta.Location)

	w.WriteHeader(http.StatusOK)
}
