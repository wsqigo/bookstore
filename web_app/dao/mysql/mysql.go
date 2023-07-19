package mysql

import (
	"bookstore/web_app/conf"
	"fmt"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func Init(cfg *conf.MysqlConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var err error
	// 也可以使用 MustConnect，连接不成功就 panic
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed, err:", zap.Error(err))
		return err
	}

	DB.SetMaxOpenConns(cfg.MaxOpenConns)
	DB.SetMaxIdleConns(cfg.MaxIdleConns)
	return nil
}

func Close() {
	_ = DB.Close()
}
