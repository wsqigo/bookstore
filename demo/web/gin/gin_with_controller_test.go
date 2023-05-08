package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestUserController_GetUser(t *testing.T) {
	g := gin.Default()
	ctrl := &UserController{}

	g.GET("/user", ctrl.GetUser)
	g.POST("/user", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello %s", "world")
	})

	g.GET("/static", func(ctx *gin.Context) {
		// 读文件
		// 写文件
	})

	_ = g.Run(":8082")
	//http.ListenAndServe(":8083", g)  engine 本身可以作为一个 Handler 传递到 http 包，用于启动服务器
}

func TestSetRouter(t *testing.T) {
	// Disable Console Color
	// gin.DisableConsoleColor()
	// engine 1. 实现了路由树功能，提供了注册和匹配路由的功能
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(ctx *gin.Context) {
		user := ctx.Params.ByName("name")
		ctx.JSON(http.StatusOK, gin.H{"user": user})
	})

	// 2. 本身作为一个Handler传递到http包，用于启动服务器
	//http.ListenAndServe(":8080", r)
	r.Run(":8080")
}
