package dao

import (
	mydb "cloudStorage/dao/mysql"
	"database/sql"

	"github.com/pkg/errors"
)

//用户注册
func UserSignUp(username string, password string) error {
	s, err := mydb.DBWriteConn().Prepare("insert ignore into tbl_user(`user_name`,`user_pwd`) values(?,?)")
	if err != nil {
		return errors.Wrap(err, "UserSignUp mydb.DBConn().Prepare")
	}
	defer s.Close()

	res, err := s.Exec(username, password)
	if err != nil {
		return errors.Wrap(err, "UserSignUp s.Exec")
	}

	if rows, err := res.RowsAffected(); err != nil {
		return errors.Wrap(err, "UserSignUp res.RowsAffected")
	} else if rows <= 0 {
		return errors.New("UserSignUp failed")
	}
	return nil
}

//用户登录
func UserSignIn(username string, password string) (*UserData, error) {
	s, err := mydb.DBReadConn().Prepare("select id,user_name,email,phone,create_at,status from tbl_user where user_name=? and user_pwd=? limit 1")
	if err != nil {
		return nil, errors.Wrap(err, "UserSignIn mydb.DBConn().Prepare")
	}
	defer s.Close()

	var user UserData
	err = s.QueryRow(username, password).Scan(&user.Id, &user.UserName, &user.Email, &user.Phone, &user.SignupAt, &user.Status)
	if err != nil {
		return nil, errors.Wrap(err, "账户或密码错误")
	}
	if user.Status == 0 {
		return nil, errors.New("该用户不存在")
	}
	if user.Status == 2 {
		return nil, errors.New("该用户被禁用，请联系管理员处理")
	}
	return &user, nil
}

//用户信息
type UserData struct {
	Id           int    `json:"user_id"`
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	SignupAt     string `json:"signup_at"`
	LastActiveAt string `json:"last_active_at"`
	Status       int    `json:"status"`
}

func GetUserInfo(userid int) (*UserData, error) {
	s, err := mydb.DBReadConn().Prepare("select id,user_name,email,phone,create_at,status from tbl_user where id=? limit 1")
	if err != nil {
		return nil, errors.Wrap(err, "GetUserInfo mydb.DBConn().Prepare")
	}
	defer s.Close()

	var user UserData
	err = s.QueryRow(userid).Scan(&user.Id, &user.UserName, &user.Email, &user.Phone, &user.SignupAt, &user.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "GetUserInfo s.QueryRow")
	}
	return &user, nil
}

func UpdateToken(userid int, token string) bool {
	s, err := mydb.DBWriteConn().Prepare("replace into tbl_user_token(`user_id`,`user_token`) values(?,?)")
	if err != nil {
		return false
	}
	defer s.Close()

	_, err = s.Exec(userid, token)
	if err != nil {
		return false
	}
	return true
}

func GetToken(userid int) string {
	s, err := mydb.DBReadConn().Prepare("select user_token from tbl_user_token where user_id=? limit 1")
	if err != nil {
		return ""
	}
	defer s.Close()

	var token string
	err = s.QueryRow(userid).Scan(&token)
	if err != nil {
		return ""
	}
	return token
}
