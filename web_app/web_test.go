package web_app

import (
	"bookstore/web_app/conf"
	"bookstore/web_app/dao/mysql"
	"bookstore/web_app/dao/redis"
	"bookstore/web_app/logger"
	"bookstore/web_app/routes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/spf13/viper"
)

// Go Web 开发较通用的脚手架模版
func TestWeb(t *testing.T) {
	// 1. 加载配置
	if err := conf.Init(); err != nil {
		fmt.Println("init conf failed, err:", err)
		return
	}
	// 2. 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Println("init logger failed, err:", err)
		return
	}
	// zap 底层有缓冲。在任何情况下执行 defer logger.Sync() 是一个很好的习惯
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")
	// 3. 初始化 MySQL 连接
	if err := mysql.Init(); err != nil {
		fmt.Println("init mysql failed, err:", err)
		return
	}
	defer mysql.Close()
	// 4. 初始化 Redis 连接
	if err := redis.Init(); err != nil {
		fmt.Println("init redis failed, err:", err)
		return
	}
	defer redis.Close()
	// 5. 注册路由
	r := routes.Setup()
	// 6. 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个 goroutine 启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个 5s 的超时
	// 创建一个接收信号的通道
	quit := make(chan os.Signal, 1)
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的 Ctrl+C 就是触发系统 SIGINT 信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能捕获，所以不需要添加它
	// signal.Notify 把收到的 syscall.SIGINT 或 syscall.SIGTERM 信号转发给 quit
	<-quit
	zap.L().Info("Shutdown Server ...")

	// 创建一个 5s 超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 5s 内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过 5s 就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
