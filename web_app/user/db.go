package user

import (
	"bookstore/web_app/code"
	"bookstore/web_app/dao/mysql"
	"bookstore/web_app/util"
)

// 把每一步数据库操作封装成函数

var (
	db = mysql.GetDBConn()
)

// CheckUserExist 检查指定 name 的用户是否存在
func CheckUserExist(name string) error {
	sqlStr := `select count(user_id) from user where username = ?`

	var count int
	err := db.Get(&count, sqlStr, name)
	if err != nil {
		return err
	}

	if count > 0 {
		return code.ErrorUserExist
	}

	return nil
}

func GetUserByName(name string) (*User, error) {
	user := &User{}
	sqlStr := `select user_id, username, password, email from user where username = ?`
	err := db.Get(user, sqlStr, name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByID(id int64) (*User, error) {
	user := &User{}
	sqlStr := `select user_id, username, password, email from user where user_id = ?`
	err := db.Get(user, sqlStr, id)
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
