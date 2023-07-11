package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// 声明一个全局的 rdb 变量
var rdb *redis.Client

func Init() error {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port"),
		),
		Password: viper.GetString("redis.password"), // no password set
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})

	return nil
}

func Close() {
	_ = rdb.Close()
}
