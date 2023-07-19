package router

import (
	"bookstore/web_app/controller"
	"bookstore/web_app/logger"
	"bookstore/web_app/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		// gin 设置成发布模式
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由
	r.POST("/signup", controller.SignUpHandler)
	// 登录
	r.POST("/login", controller.LoginHandler)
	// auth
	r.POST("/auth", controller.LoginHandler)
	// Token
	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(ctx *gin.Context) {
		// 如果是登录的用户，判断请求头中是否有有效的 JWT
		ctx.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
