package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// 声明一个全局的 rdb 变量
var rdb *redis.Client

// 初始化连接
func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 20, // 连接池大小
	})

	s, err := rdb.Ping(context.Background()).Result()
	fmt.Println(s) // Pong
	if err != nil {
		panic(err)
	}
}

func redisExample() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 直接执行命令获取错误
	err := rdb.Set(ctx, "score", 100, 0).Err()
	if err != nil {
		fmt.Printf("set score failed, err:%v\n", err)
		return
	}

	// 执行命令获取结果
	val, err := rdb.Get(ctx, "score").Result()
	if err != nil {
		fmt.Printf("get score failed, err:%v\n", err)
		return
	}
	fmt.Println("score", val)

	// 先获取到命令对象
	cmder := rdb.Get(ctx, "score")
	fmt.Println(cmder.Val()) // 获取值
	fmt.Println(cmder.Err()) // 获取错误

	// 直接执行命令获取值
	value := rdb.Get(ctx, "score").Val()
	fmt.Println(value)

	// redis.Nil 判断
	value2, err := rdb.Get(ctx, "messi").Result()
	if err == redis.Nil {
		fmt.Println("name does not exist")
		return
	}

	if err != nil {
		fmt.Println("get name failed, err:", err)
		return
	}

	fmt.Println("name", value2)
}

func hgetDemo() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	val, err := rdb.HGetAll(ctx, "user:1:info").Result()
	if err != nil {
		// redis.Nil
		// 其他错误
		fmt.Println("hgetall failed, err:", err)
		return
	}
	fmt.Println(val)

	v2 := rdb.HMGet(ctx, "user:1:info", "name", "age").Val()
	fmt.Println(v2)

	v3 := rdb.HGet(ctx, "user:1:info", "age").Val()
	fmt.Println(v3)
}

// zsetDemo 操作zset示例
func zsetDemo() {
	// key
	zsetKey := "language_rank"
	// value
	languages := []*redis.Z{
		{Score: 90.0, Member: "Golang"},
		{Score: 98.0, Member: "Java"},
		{Score: 95.0, Member: "Python"},
		{Score: 97.0, Member: "JavaScript"},
		{Score: 99.0, Member: "C/C++"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// ZADD
	err := rdb.ZAdd(ctx, zsetKey, languages...).Err()
	if err != nil {
		fmt.Printf("zadd failed, err:%v\n", err)
		return
	}
	fmt.Println("zadd success")

	// 把 Golang 的分数加 10
	newScore, err := rdb.ZIncrBy(ctx, zsetKey, 10, "Golang").Result()
	if err != nil {
		fmt.Println("incr golang score failed:", err)
		return
	}

	fmt.Println("Golang score is ", newScore)

	// 取分数最高的 3 个
	ret := rdb.ZRevRangeWithScores(ctx, zsetKey, 0, 2).Val()
	for _, z := range ret {
		fmt.Println(z.Score, z.Member)
	}

	// 取 95~100 分的
	op := &redis.ZRangeBy{
		Min: "95",
		Max: "100",
	}
	ret, err = rdb.ZRevRangeByScoreWithScores(ctx, zsetKey, op).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err:%v\n", err)
		return
	}

	for _, z := range ret {
		fmt.Println(z.Score, z.Member)
	}
}

// scanKeysDemo1 按前缀查找所有 key 示例
func scanKeysDemo1() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 按前缀扫描key
	iter := rdb.Scan(ctx, 0, "prefix*", 0).Iterator()
	for iter.Next(ctx) {
		fmt.Println("keys", iter.Val())
	}

	if err := iter.Err(); err != nil {
		panic(err)
	}
}

// delKeysByMatch 按match格式扫描所有 key 并删除
func delKeysByMatch(match string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	iter := rdb.Scan(ctx, 0, match, 0).Iterator()
	for iter.Next(ctx) {
		err := rdb.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}
}
