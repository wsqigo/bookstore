package database

import (
	"fmt"
	"testing"
)

func TestMysql(t *testing.T) {
	err := initDB() // 调用初始化数据库的函数
	if err != nil {
		fmt.Printf("init db failed, err:%v\n", err)
		return
	}

	fmt.Println("connect to db success")
}
