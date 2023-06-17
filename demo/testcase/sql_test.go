package testcase

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// https://liwenzhou.com/posts/Go/mysql/

type user struct {
	id   int
	age  int
	name string
}

// 查询单条数据示例
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"
	var u user
	// 非常重要
	// QueryRow 总是返回非 nil 的值，直到返回值的 Scan 方法被调用时，才会返回被延迟的错误。
	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}

	fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
}

// func

// 查询多条数据示例
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}

	// 非常重要：关闭 rows 释放持有的数据库连接
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err = rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v", err)
			return
		}

		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}

// 插入数据
func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "王五", 18)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}

	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}

	fmt.Printf("insert success, the id id %d.\n", theID)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 3)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}

	n, err := ret.RowsAffected() // 操作影响的函数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}

	fmt.Printf("update success, affected rows:%d\n", n)
}

// 定义一个全局对象 db
// 表示连接的数据库对象（结构体实例），
// 内部维护着一个具有零到多个底层连接的连接池，它可以安全地被多个 goroutine 同时使用
var db *sql.DB

// 定义一个初始化数据库的函数
func initDB() error {
	// DSN: Data Source Name
	dsn := "root:root@tcp(127.0.0.1:3306)/sql_demo?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码正确
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// 尝试与数据库建立连接（校验 dsn 是否正确）
	err = db.Ping()
	if err != nil {
		fmt.Printf("connect to db failed, err:%v\n", err)
		return err
	}

	// 数值需要业务具体情况来确定
	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxOpenConns(200) // 最大连接数

	return nil
}

func TestMysql(t *testing.T) {
	err := initDB() // 调用初始化数据库的函数
	if err != nil {
		fmt.Printf("init db failed, err:%v\n", err)
		return
	}

	fmt.Println("connect to db success")
}
