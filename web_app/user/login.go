package user

import (
	"bookstore/web_app/pkg/jwt"
	"bookstore/web_app/util"
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrorUserNotLogin = errors.New("user is not login")
	ErrorUserNotExist = errors.New("user is not exist")
)

func Login(user User) (string, error) {
	// 根据用户名查询用户信息
	dbUser, err := GetUserByName(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrorUserNotExist
		}
		return "", err
	}

	// 判断密码是否正确
	password := util.EncryptPassword(user.Password)
	if password != dbUser.Password {
		return "", errors.New("password is not correct")
	}

	// 生成 JWT Token
	return jwt.GenToken(dbUser.UserID, dbUser.Username)
}

// GetCurrentUser 获取当前登录的用户 ID
func GetCurrentUser(ctx *gin.Context) (int64, error) {
	uid, ok := ctx.Get("userID")
	if !ok {
		return 0, ErrorUserExist
	}

	userID, ok := uid.(int64)
	if !ok {
		return 0, ErrorUserNotLogin
	}

	return userID, nil
}
