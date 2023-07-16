package mysql

import (
	"bookstore/web_app/model"
	"crypto/md5"
	"encoding/hex"
	"errors"
)

// 把每一步数据库操作封装成函数

const secret = "wsqigo"

// CheckUserExist 检查指定 name 的用户是否存在
func CheckUserExist(name string) error {
	sqlStr := `select count(user_id) from user where username = ?`

	var count int
	err := db.Get(&count, sqlStr, name)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("username already exist. name = " + name)
	}

	return nil
}

// InsertUser 插入一条新的用户记录
func InsertUser(user *model.User) error {
	// 对密码进行加密
	password := encryptPassword(user.Password)
	// 执行 SQL 语句入库
	sqlStr := `insert into user(user_id, username, password, email) values(?,?,?,?)`
	_, err := db.Exec(sqlStr, user.UserID, user.Username, password, user.Email)

	return err
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))

	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
