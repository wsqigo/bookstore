package controller

import (
	"bookstore/web_app/user"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpReq 定义请求的参数结构体
type SignUpReq struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Email      string `json:"email" bind:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

func SignUpHandler(ctx *gin.Context) {
	// 1. 获取参数和参数校验
	req := &SignUpReq{}
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		// 请求参数有误，直击返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断 err 是不是 validator.ValidationErrors 错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}

		// validator.ValidateErrors 类型错误则使用翻译器
		// 并使用removeTopStruct函数去除字段名中的结构体名称标识
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2. 业务处理
	err = user.SignUp(&user.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		zap.L().Error("sign up failed", zap.Error(err))
		if errors.Is(err, user.ErrorUserExist) {
			ResponseError(ctx, CodeUserExist)
			return
		}
		ResponseErrorWithMsg(ctx, CodeServerBusy, "")
		return
	}

	// 3. 返回响应
	ResponseSuccess(ctx, CodeSuccess)
}

// 登录请求参数
type loginReq struct {
	UserName string `json:"user_name" binding:"require"`
	Password string `json:"password" binding:"password"`
}

func LoginHandler(ctx *gin.Context) {
	// 1. 获取请求参数及参数校验
	l := &loginReq{}
	err := ctx.ShouldBindJSON(l)
	if err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("login with invalid param",
			zap.String("username", l.UserName), zap.Error(err))
		// 判断 err 是不是 validator.ValidationErrors 类型
		errs, ok := err.(*validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			// 错误翻译一下
			"msg": removeTopStruct(errs.Translate(trans)),
		})
		return
	}

	// 2. 业务逻辑处理
	token, err := user.Login(user.User{
		Username: l.UserName,
		Password: l.Password,
	})
	if err != nil {
		zap.L().Error("login failed", zap.Error(err))
		if errors.Is(err, user.ErrorUserNotExist) {
			ResponseError(ctx, CodeUserExist)
			return
		}
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	// 3. 返回响应
	ResponseSuccess(ctx, token)
}
