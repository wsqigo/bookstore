package database

import (
	"database/sql"
	"fmt"
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
	// 非常重要: 确保 QueryRow 之后调用 Scan 方法，否则持有的数据库链接不会被释放
	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}

	fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
}

// 查询多条数据示例
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}

	// 非常重要: 关闭 rows 释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err = rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
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

	fmt.Printf("insert success, the id is %d.\n", theID)
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

func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 3)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}

	n, err := ret.RowsAffected() // 操作影响的函数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}

	fmt.Printf("delete success, affected rows:%d\n", n)
}

// 预处理查询示例
func prepareQueryDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}

	defer stmt.Close()
	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		u := &user{}
		err = rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}

		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}

// 预处理插入示例
func prepareInsertDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}

	defer stmt.Close()
	_, err = stmt.Exec("小王子", 18)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}

	_, err = stmt.Exec("沙河娜扎", 18)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}

	fmt.Println("insert success.")
}

func sqlInjectDemo(name string) {
	sqlStr := fmt.Sprintf("select id, name, age from user where name='%s'", name)
	fmt.Printf("SQL:%s\n", sqlStr)

	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}

	for rows.Next() {
		var u user
		err = rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("user:%#v\n", u)
	}
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

// 事务操作示例
func transactionDemo() {
	tx, err := db.Begin() // 开启事务
	if err != nil {
		fmt.Printf("begin transfer")
		return
	}

	defer func() {
		if err != nil {
			fmt.Println("事务回滚啦...")
			tx.Rollback()
		}
	}()

	sqlStr1 := "update user set age=30 where id = ?"
	ret1, err := tx.Exec(sqlStr1, 2)
	if err != nil {
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}

	affRow1, err := ret1.RowsAffected()
	if err != nil {
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	sqlStr2 := "update user set age=40 where id = ?"
	ret2, err := tx.Exec(sqlStr2, 3)
	if err != nil {
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}

	affRow2, err := ret2.RowsAffected()
	if err != nil {
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	fmt.Println(affRow1, affRow2)
	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("事务提交啦...")
		tx.Commit()
	}

	fmt.Println("exec trans success!")
}
