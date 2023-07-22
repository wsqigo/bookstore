package mysql

import (
	"bookstore/web_app/conf"
	"fmt"
	"sync"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db   *sqlx.DB
	once sync.Once
)

func GetDBConn() *sqlx.DB {
	once.Do(func() {
		Init()
	})

	return db
}

func Init() {
	conf.Init()
	cfg := conf.Conf.MysqlConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var err error
	// 也可以使用 MustConnect，连接不成功就 panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect db failed, err:", zap.Error(err))
		panic("init mysql failed, err: " + err.Error())
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
}

func Close() {
	_ = db.Close()
}
