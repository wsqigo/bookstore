package shutdown

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func TestShutdown(t *testing.T) {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		time.Sleep(5 * time.Second)
		ctx.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 开启一个 goroutine 启动服务
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("start server failed, err: %v", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个 5 秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的 Ctrl+C 就是触发系统 SIGINT 信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能捕获，所以不需要添加它
	// signal.Notify 把收到的 syscall.SIGINT 或 syscall.SIGTERM 信号转发给 quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞

	<-quit // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("Shutdown Server ...")

	// 创建一个 5s 超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 5s 内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过 5s 就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}

func TestRestart(t *testing.T) {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "hello gin!")
	})

	// 默认 endless 服务器会监听下列信号:
	// syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM 和 syscall.SIGSTP
	// 接收到 SIGHUP 信号将触发 `fork/restart` 实现优雅重启（kill -1 pid 会发送 SIGHUP 信号）
	// 接收到 syscall.SIGINT 或 syscall.SIGTERM 信号将触发优雅关机
	// 接收到 SIGUSR2 信号将触发 HammerTime
	// SIGUSR1 和 SIGSTP 被用来触发一些用户自定义的 hook 函数
	err := endless.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("listen: %s\n", err)
	}

	log.Println("Server exiting")
}
