package gorm

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	username := "root"  // 账号
	password := "root"  // 密码
	host := "127.0.0.1" // 数据库地址，可以是 Ip 或者域名
	port := 3306        // 数据库端口
	dbName := "gorm"    // 数据库名
	timeout := "10s"    // 连接超时，10s

	// root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&timeout=%s",
		username, password, host, port, dbName, timeout)
	var err error
	// 连接 mysql，获得 DB 类型示例，用于后面的数据库读写操作
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）
		// 如果没有这方面的要求，您可以在初始化时禁用它，这样可以获得 60% 的性能提升
		SkipDefaultTransaction: true,
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix:   "f_", // table前缀
		//	SingularTable: true, // 单数表名 Student -> student
		//	NoLowerCase:   true, // 关闭小写转换
		//},
		// 配置要显示的日志等级
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("connect database failed, err: " + err.Error())
	}
}

/*
如果想自定义日志的显示
newLogger := logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // （日志输出的目标，前缀和日志包含的内容
	logger.Config{
		SlowThreshold: time.Second, // 慢 SQL 阈值
		LogLevel: logger.Info, // 日志级别
		IgnoreRecordNotFoundError: true, // 忽略 ErrRecordNotFound 错误
		Colorful: true, // 使用彩色打印
	},
)
db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
	Logger: newLogger,
})
*/
