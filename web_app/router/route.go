package router

import (
	"bookstore/web_app/community"
	"bookstore/web_app/controller"
	"bookstore/web_app/logger"
	"bookstore/web_app/middlewares"
	"bookstore/web_app/post"
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

	v1 := r.Group("/api/v1")
	// 注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	// 登录
	v1.POST("/login", controller.LoginHandler)

	// 应用 JWT 认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())

	{
		v1.GET("/community", community.GetCommunityConf)
		v1.GET("/community/:id", community.GetCommunityDetail)
		v1.POST("/post", post.CreatePost)
		v1.POST("/post/:id", post.GetPostDetailHandler)
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
