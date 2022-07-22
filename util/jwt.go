package util

import (
	"fmt"
	"time"

	dao "cloudStorage/dao"
)

func GenToken(username string) string {
	//40位字符：md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	token := MD5([]byte(username + ts + "_tokensalt"))
	return token + ts[:8]
}

func ValidToken(userid int, token string) bool {
	if len(token) != 40 {
		return false
	}
	return token == dao.GetToken(userid)
}
