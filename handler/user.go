package handler

import (
	"cloudStorage/util"
	"net/http"
	"strconv"

	dao "cloudStorage/dao"
)

const (
	PWD_SALT = "*#890"
)

//用户注册
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) < 3 || len(password) < 5 {
		util.ErrorResponse(w, "Invalid username or password")
		return
	}

	enc_pwd := util.Sha1([]byte(password + PWD_SALT))
	if err := dao.UserSignUp(username, enc_pwd); err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}
	util.SuccessResponse(w, nil)
}

//用户登录
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	//校验用户名和密码
	enc_pwd := util.Sha1([]byte(password + PWD_SALT))
	user, err := dao.UserSignIn(username, enc_pwd)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	//生成token，并存入数据库
	token := util.GenToken(username)
	// fmt.Println(user)
	b := dao.UpdateToken(user.Id, token)
	if !b {
		util.ErrorResponse(w, err.Error())
		return
	}

	//登录成功，重定向到首页
	util.SuccessResponse(w, map[string]interface{}{
		"location": "/static/view/home.html",
		"token":    token,
		"userid":   user.Id,
	})
}

//用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userid, _ := strconv.Atoi(r.Form.Get("userid"))

	//获取用户信息
	user, err := dao.GetUserInfo(userid)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	util.SuccessResponse(w, user)
}
