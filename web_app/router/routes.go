package router

import (
	"bookstore/web_app/controller"
	"bookstore/web_app/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由
	r.POST("/signup", controller.SignUpHandler)

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
