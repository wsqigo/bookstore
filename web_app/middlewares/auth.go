package middlewares

import (
	"bookstore/web_app/controller"
	"bookstore/web_app/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

// JWTAuthMiddleware 基于 JWT 的认证中间件
func JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// 客户端携带 Token 有三种方式
		// 1. 放在请求头 2. 放在请求体 3. 放在URI
		// 这里假设 Token 放在 Header 的 Authorization 中，并使用 Bearer 开头
		// Authorization: Bearer xxxxxx.xxxxx
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			controller.ResponseError(ctx, controller.CodeNeedLogin)
			ctx.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(ctx, controller.CodeInvalidToken)
			ctx.Abort()
			return
		}

		// parts[1] 是获取到的 tokenString，我们使用之前定义好的解析 jWT 函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(ctx, controller.CodeInvalidToken)
			ctx.Abort()
			return
		}
		// 将当前请求的 userID 信息保存到请求的上下文 ctx 上
		ctx.Set(CtxUserIDKey, mc)
		// 后续的处理函数可以通过 ctx.Get(CtxUserIDKey) 来获取当前请求的用户信息
		ctx.Next()
	}
}
