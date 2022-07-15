package main

import (
	"cloudStorage/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)            //处理文件上传
	http.HandleFunc("/file/upload/success", handler.UploadSucHandler) //处理文件上传成功
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)         //处理获取文件元信息
	http.HandleFunc("/file/metas", handler.GetFileMetasHandler)       //处理获取批量文件元信息
	http.HandleFunc("/file/download", handler.DownloadHandler)        //处理文件下载
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)    //处理文件元信息更新
	http.HandleFunc("/file/delete", handler.FileMetaDeleteHandler)    //处理文件元信息删除

	//静态文件服务器
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("failed to start server, err:" + err.Error())
	}
}
