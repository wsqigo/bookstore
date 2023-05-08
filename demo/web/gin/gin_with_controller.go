package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
	Engine实现了
*/

type UserController struct {
}

func (u *UserController) GetUser(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello, world")
}
