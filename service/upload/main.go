package main

import (
	"cloudStorage/handler"
	"cloudStorage/util"
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	NewRouter("/file/upload", "POST", HttpInterceptor(handler.UploadHandler))          //文件上传
	NewRouter("/file/fastupload", "POST", HttpInterceptor(handler.FastUploadHandler))  //秒传文件
	NewRouter("/file/list", "GET", HttpInterceptor(handler.ListHandler))               //文件列表
	NewRouter("/file", "GET", HttpInterceptor(handler.GetFileHandler))                 //文件信息
	NewRouter("/file/download", "GET", HttpInterceptor(handler.DownloadHandler))       //文件下载
	NewRouter("/file/downloadurl", "GET", HttpInterceptor(handler.DownloadURLHandler)) //文件下载地址
	NewRouter("/file/update", "POST", HttpInterceptor(handler.FileUpdateHandler))      //文件信息更新
	NewRouter("/file/delete", "GET", HttpInterceptor(handler.FileDeleteHandler))       //文件删除

	NewRouter("/user/signup", "POST", handler.SignupHandler)            //用户注册
	NewRouter("/user/signin", "POST", handler.SigninHandler)            //用户登录
	NewRouter("/user", "GET", HttpInterceptor(handler.UserInfoHandler)) //用户信息

	//分块上传接口
	NewRouter("/file/mpupload/init", "GET", HttpInterceptor(handler.InitialMultipartUploadHandler)) //初始化分块上传
	NewRouter("/file/mpupload/uppart", "POST", HttpInterceptor(handler.UploadPartHandler))          //上传文件分块
	NewRouter("/file/mpupload/complete", "POST", HttpInterceptor(handler.CompleteUploadHander))     //通知上传合并

	//静态文件服务器
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("failed to start server, err:" + err.Error())
	}
}

func NewRouter(url string, method string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			fmt.Println("OK " + url + " " + "[" + method + "]")
			handler(w, r)
		} else {
			fmt.Println("404 " + url + " " + "[" + method + "]")
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func HttpInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			userid, _ := strconv.Atoi(r.Form.Get("userid"))
			token := r.Form.Get("token")

			//校验token
			if b := util.ValidToken(userid, token); !b {
				util.ErrorResponse(w, "Invalid token")
				return
			}

			r.Header.Set("userid", r.Form.Get("userid"))
			h(w, r)
		},
	)
}
