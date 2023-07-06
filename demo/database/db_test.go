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

	//queryRowDemo()
	//queryMultiRowDemo()
	//insertRowDemo()
	//queryMultiRowDemo()
	fmt.Println("connect to db success")

	//prepareInsertDemo()
	//prepareQueryDemo()

	//sqlInjectDemo("qimi")
	////sqlInjectDemo("qimi' or 1=1#")
	//sqlInjectDemo("qimi' union select * from user #")
	//sqlInjectDemo("qimi' and (select count(*) from user) < 10 #")

	transactionDemo()
}
