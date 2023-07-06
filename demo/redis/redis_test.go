package redis

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// docker cli
// docker run --name redis -p 6379:6379 -d redis
// docker run -it --network host --rm redis redis-cli

func TestRedisDemo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 执行命令获取结果
	val, err := rdb.Get(ctx, "hello").Result()
	fmt.Println(val, err)

	// 先获取到命令对象
	cmder := rdb.Get(ctx, "hello")
	fmt.Println(cmder.Val()) // 获取值
	fmt.Println(cmder.Err()) // 获取错误

	// 直接执行命令获取错误
	err = rdb.Set(ctx, "counter", 10, time.Hour).Err()

	// 执行执行命令获取值
	value := rdb.Get(ctx, "counter").Val()
	fmt.Println(value)
}

func TestHashRedis(t *testing.T) {
	hgetDemo()
}

func TestDoDemo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 直接执行命令获取错误
	err := rdb.Do(ctx, "set", "go", "best", "EX", 3600).Err()
	fmt.Println(err)

	// 执行命令获取结果
	val, err := rdb.Do(ctx, "get", "go").Result()
	fmt.Println(val, err)
}
