package sqlx_demo

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// https://www.liwenzhou.com/posts/Go/sqlx/
var db *sqlx.DB

func initDB() error {
	dsn := "root:root@tcp(127.0.0.1:3306)/sql_demo?charset=utf8mb4&parseTime=True"
	var err error
	// 也可以使用 MustConnect 连接不成功就 panic
	//db = sqlx.MustConnect("mysql", dsn)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("connect DB failed, err:", err)
		return err
	}

	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(10)
	return nil
}

type user struct {
	ID   int    `db:"id"`
	Age  int    `db:"age"`
	Name string `db:"name"`
}

func (u *user) Value() (driver.Value, error) {
	return []any{u.Name, u.Age}, nil
}

// 查询单条数据示例
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"
	var u user
	err := db.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}

	fmt.Printf("id:%d name:%s age:%d\n", u.ID, u.Name, u.Age)
}

func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	var users []user
	err := db.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Println("query failed, err:", err)
		return
	}

	fmt.Printf("%#v\n", users)
}

// 插入数据
func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "沙河小王子", 19)
	if err != nil {
		fmt.Println("insert failed, err:", err)
		return
	}

	theID, err := ret.LastInsertId()
	if err != nil {
		fmt.Println("get lastInsertId failed, err", err)
		return
	}

	fmt.Printf("insert success, the id is %d.\n", theID)
}

// 更新数据
func updateDemo() {
	sqlStr := "update user set age = ? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 6)
	if err != nil {
		fmt.Println("update failed, err:", err)
		return
	}

	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("get RowsAffected failed, err:", err)
		return
	}

	fmt.Println("update success, affected rows:", n)
}

// 删除数据
func deleteDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 6)
	if err != nil {
		fmt.Println("delete failed, err:", err)
		return
	}

	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Println("get RowsAffected failed, err:", err)
		return
	}

	fmt.Println("delete success, affected rows:", n)
}

// NameExec
func insertUserDemo() {
	sqlStr := "insert into user (name, age) values (:name, :age)"

	ret, err := db.NamedExec(sqlStr, user{
		ID:   6,
		Name: "wsqigo",
		Age:  28,
	})

	if err != nil {
		fmt.Println("insert failed, err:", err)
		return
	}

	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("get RowsAffected failed, err:", err)
		return
	}

	fmt.Println("insert success, affected rows:", n)
}

func namedQuery() {
	sqlStr := "select * from user where name=:name"
	// 使用 map 做命名查询
	rows, err := db.NamedQuery(sqlStr, map[string]any{
		"name": "七米",
	})
	if err != nil {
		fmt.Println("db.NamedQuery failed, err:", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var u user
		err = rows.StructScan(&u)
		if err != nil {
			fmt.Println("scan failed, err:", err)
			return
		}
		fmt.Printf("%#v\n", u)
	}

	u := user{
		Name: "七米",
	}
	// 使用结构体命名查询，根据结构体字段的 db tag 进行映射
	rows, err = db.NamedQuery(sqlStr, u)
	if err != nil {
		fmt.Println("db.NamedQuery failed, err:", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var u user
		err = rows.StructScan(&u)
		if err != nil {
			fmt.Println("scan failed, err:", err)
			return
		}
		fmt.Printf("%#v\n", u)
	}
}

func transactionDemo() error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println("begin trans failed, err:", err)
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw panic after Rollback
		}

		if err != nil {
			fmt.Println("rollback")
			tx.Rollback()
		} else {
			fmt.Println("commit")
			tx.Commit()
		}
	}()

	sqlStr1 := "update user set age=20 where id=?"
	ret1, err := tx.Exec(sqlStr1, 1)
	if err != nil {
		return err
	}

	n, err := ret1.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}

	sqlStr2 := "Update user set age=50 where i=?"
	ret2, err := tx.Exec(sqlStr2, 5)
	if err != nil {
		return err
	}
	n, err = ret2.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr2 failed")
	}
	return err
}

// BatchInsertUsers 自行构造批量插入的语句
func BatchInsertUsers(users []*user) error {
	// 存放 (?, ?) 的 slice
	valueStrings := make([]string, 0, len(users))
	// 存放 values 的 slice
	valueArgs := make([]any, 0, 2*len(users))
	// 遍历 users 准备相关数据
	for _, u := range users {
		// 此处占位符要与插入值的个数对应
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, u.Name, u.Age)
	}

	// 自行拼接要执行的具体语句
	stmt := fmt.Sprintf("insert into user (name, age) values %s",
		strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	return err
}

// BatchInsertUsers2 使用 sqlx.In 帮我们拼接语句和参数，注意传入的参数是 []any
func BatchInsertUsers2(users []any) error {
	query, args, _ := sqlx.In(
		"insert into user (name, age) VALUES (?), (?), (?)",
		users...) // 如果 arg 实现了 driver.Valuer, sql.In 会通过调用 Value() 来展开它
	fmt.Println(query) // 查看生成的 querystring
	fmt.Println(args)
	_, err := db.Exec(query, args...)
	return err
}

// QueryByIDs 根据给定 ID 查询
func QueryByIDs(ids []int) ([]user, error) {
	// 动态填充 id
	query, args, err := sqlx.In("select name, age from user where id in(?)", ids)
	if err != nil {
		return nil, err
	}

	fmt.Println("sql:", query)
	// sqlx.In 返回带 `?` bindvar 的查询语句，我们使用 Rebind() 重新绑定它
	query = db.Rebind(query)

	var users []user
	err = db.Select(&users, query, args...)
	return users, err
}

// QueryAndOrderByIDs 按照指定 id 查询并维护顺序
func QueryAndOrderByIDs(ids []int) ([]user, error) {
	// 动态填充 id
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, fmt.Sprint(id))
	}

	query, args, err := sqlx.In("select name, age from user where id in(?) order by find_in_set(id, ?)",
		ids, strings.Join(strIDs, ","))
	if err != nil {
		return nil, err
	}

	// sqlx.In 返回带 `?` bindvar 的查询语句，我们使用 Rebind() 重新绑定它
	query = db.Rebind(query)

	var users []user
	err = db.Select(&users, query, args...)
	return users, err
}
