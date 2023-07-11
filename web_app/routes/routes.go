package routes

import (
	"bookstore/web_app/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.Handle(http.MethodGet, "/", func(ctx *gin.Context) {

	})

	return r
}
