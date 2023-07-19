package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code": 10001, // 程序中的错误码
	"msg" xxx, 	   // 提示信息
	"data" {} 	   // 数据
}
*/

type ResponseData struct {
	Code ResCode `json:"code"`
	Msg  any     `json:"msg"`
	Data any     `json:"data"`
}

func ResponseError(ctx *gin.Context, code ResCode) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  code.Msg(),
		"data": nil,
	})
}

func ResponseSuccess(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

func ResponseErrorWithMsg(ctx *gin.Context, code ResCode, msg any) {
	ctx.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
