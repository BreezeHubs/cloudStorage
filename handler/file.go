package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"cloudStorage/config"
	dao "cloudStorage/dao"
	"cloudStorage/dao/mq"
	"cloudStorage/store/oss"
	"cloudStorage/util"
)

type FileData struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Hash     string `json:"hash"`
	Location string `json:"location"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}

//文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	//handle the upload
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		util.ErrorResponse(w, "form"+err.Error())
		return
	}
	defer file.Close()

	//暂存至tmp目录获取size和hash，获取后删除
	tmpfilehash, tmpfilesize, err := util.GetFileHash(file, fileHeader.Filename)
	if err != nil {
		util.ErrorResponse(w, "tmp: "+err.Error())
		return
	}

	//保存至正式目录
	newfilepath := config.FILE_STATIC_PATH + tmpfilehash + path.Ext(fileHeader.Filename)
	if err := util.SaveFile(file, newfilepath); err != nil {
		util.ErrorResponse(w, "file: "+err.Error())
		return
	}

	//写入ceph存储
	file.Seek(0, 0)
	// data, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	util.ErrorResponse(w, "ceph: "+err.Error())
	// 	return
	// }
	// bucket := ceph.GetCephBucket("userfile")
	// cephPath := "/ceph/" + tmpfilehash

	// // fmt.Println(cephPath) 事先需要创建bucket，bucket.PutBucket(s3.PublicRead)
	// if err := bucket.Put(cephPath, data, "octet-stream", s3.PublicRead); err != nil {
	// 	util.ErrorResponse(w, "bucket put: "+err.Error())
	// 	return
	// }

	//写入oss
	ossPath := "oss/" + tmpfilehash
	// ossClient, err := oss.Bucket()
	// if err != nil {
	// 	util.ErrorResponse(w, "oss: "+err.Error())
	// 	return
	// }
	// if err := ossClient.PutObject(ossPath, file); err != nil {
	// 	util.ErrorResponse(w, "oss put: "+err.Error())
	// 	return
	// }
	data := mq.TransferData{
		FileHash:      tmpfilehash,
		CurLocation:   newfilepath,
		DestLocation:  ossPath,
		DestStoreType: mq.StoreOSS,
	}
	pubData, err := json.Marshal(data)
	if err != nil {
		util.ErrorResponse(w, "marshal: "+err.Error())
		return
	}
	if err := mq.Publish(
		config.TRANS_EXCHANGE_NAME,
		config.TRANS_OSS_ROUTING_KEY,
		pubData,
	); err != nil {
		util.ErrorResponse(w, "mq.Publish: "+err.Error())
		return
	}

	//create a file meta
	fileData := FileData{
		Name:     fileHeader.Filename,
		Location: ossPath, //cephPath, //newfilepath,
		Hash:     tmpfilehash,
		Size:     tmpfilesize,
	}

	userid, _ := strconv.Atoi(r.Form.Get("userid"))
	//update the file meta to db
	if err = dao.OnFileUploadFinished(
		userid,
		fileData.Hash,
		fileData.Name,
		fileData.Size,
		fileData.Location,
		false,
	); err != nil {
		util.ErrorResponse(w, "db: "+err.Error())
		return
	}

	util.SuccessResponse(w, fileData)
}

//秒传文件
func FastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("hash")
	filename := r.Form.Get("filename")

	fileMeta, err := dao.GetFileMeta(filehash)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	//create a file meta
	fileData := FileData{
		Name:     filename,
		Location: fileMeta.FileAddr.String,
		Hash:     filename,
		Size:     fileMeta.FileSize.Int64,
	}

	userid, _ := strconv.Atoi(r.Form.Get("userid"))
	//update the file meta to db
	if err = dao.OnFileUploadFinished(
		userid,
		fileData.Hash,
		fileData.Name,
		fileData.Size,
		fileData.Location,
		true,
	); err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	util.SuccessResponse(w, fileData)
}

//文件信息
func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["hash"][0]
	tfile, err := dao.GetUserFileMeta(filehash)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	util.SuccessResponse(w, &FileData{
		Name:     tfile.FileName.String,
		Size:     tfile.FileSize.Int64,
		Hash:     tfile.FileHash,
		Location: tfile.FileAddr.String,
		CreateAt: tfile.CreateAt.String,
		UpdateAt: tfile.UpdateAt.String,
	})
}

//文件列表
func ListHandler(w http.ResponseWriter, _ *http.Request) {
	tfiles, err := dao.GetUserFileList()
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}
	var files []*FileData
	for _, tfile := range tfiles {
		files = append(files, &FileData{
			Name:     tfile.FileName.String,
			Size:     tfile.FileSize.Int64,
			Hash:     tfile.FileHash,
			Location: tfile.FileAddr.String,
			CreateAt: tfile.CreateAt.String,
			UpdateAt: tfile.UpdateAt.String,
		})
	}
	util.SuccessResponse(w, files)
}

//文件下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("hash")
	tfile, err := dao.GetUserFileMeta(filehash)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	f, err := os.Open(tfile.FileAddr.String)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	util.DownloadFile(w, tfile.FileName.String, b)
}

//oss文件下载
func DownloadURLHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("hash")
	tfile, err := dao.GetFileMeta(filehash)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	//判断存储在oss还是ceph...

	signedURL, err := oss.DownloadUrl(tfile.FileAddr.String)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}
	util.SuccessResponse(w, map[string]string{"url": signedURL})
}

//文件信息更新
func FileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	op := r.Form.Get("op")
	filehash := r.Form.Get("hash")
	newfilename := r.Form.Get("filename")

	if op != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	tfile, err := dao.GetUserFileMeta(filehash)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}
	// oldFileName := tfile.FileName
	// tfile.FileName = sql.NullString{
	// 	String: newfilename,
	// 	Valid:  true,
	// }
	// tfile.FileAddr = sql.NullString{
	// 	String: config.FILE_STATIC_PATH + newfilename,
	// 	Valid:  true,
	// }

	//rename the file
	// if err := os.Rename(config.FILE_STATIC_PATH+oldFileName.String, config.FILE_STATIC_PATH+newfilename); err != nil {
	// 	util.ErrorResponse(w, err.Error())
	// 	return
	// }
	if err := dao.UpdateUserFile(filehash, newfilename+path.Ext(tfile.FileName.String)); err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	util.SuccessResponse(w, tfile)
}

//文件删除
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")

	//delete the file
	tfile, err := dao.GetUserFileMeta(filehash)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	// if err = os.Remove(tfile.FileAddr.String); err != nil {
	// 	util.ErrorResponse(w, err.Error())
	// 	return
	// }
	if err := dao.DeleteUserFile(tfile.FileHash); err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	util.SuccessResponse(w, nil)
}
