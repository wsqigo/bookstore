package redis

import (
	"bookstore/web_app/conf"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// 声明一个全局的 rdb 变量
var rdb *redis.Client

func Init(cfg *conf.RedisConfig) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	return nil
}

func Close() {
	_ = rdb.Close()
}
