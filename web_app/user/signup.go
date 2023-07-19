package user

import (
	"bookstore/web_app/snowflake"
)

func SignUp(u *User) error {
	// 1. 判断用户存不存在
	err := CheckUserExist(u.Username)
	if err != nil {
		// 数据库查询出错
		return err
	}

	// 2. 生成 UID
	u.UserID = snowflake.GenID()

	// 3. 保存进数据库
	return InsertUser(u)
}
