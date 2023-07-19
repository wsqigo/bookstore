package user

import (
	"bookstore/web_app/dao/mysql"
	"bookstore/web_app/util"
	"errors"
)

// 把每一步数据库操作封装成函数

var (
	ErrorUserExist = errors.New("username already exist")
)

var db = mysql.DB

// CheckUserExist 检查指定 name 的用户是否存在
func CheckUserExist(name string) error {
	sqlStr := `select count(user_id) from user where username = ?`

	var count int
	err := db.Get(&count, sqlStr, name)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrorUserExist
	}

	return nil
}

func GetUserByName(name string) (*User, error) {
	user := &User{}
	sqlStr := `select user_id, username, password, email from user wher username = ?`
	err := db.Get(user, sqlStr, name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// InsertUser 插入一条新的用户记录
func InsertUser(user *User) error {
	// 对密码进行加密
	password := util.EncryptPassword(user.Password)
	// 执行 SQL 语句入库
	sqlStr := `insert into user(user_id, username, password, email) values(?,?,?,?)`
	_, err := db.Exec(sqlStr, user.UserID, user.Username, password, user.Email)

	return err
}
