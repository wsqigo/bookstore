package controller

import (
	"bookstore/web_app/dao/mysql"
	"bookstore/web_app/model"
	"bookstore/web_app/snowflake"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SignUp(param *SignUpReq) error {
	// 1. 判断用户存不存在
	err := mysql.CheckUserExist(param.Username)
	if err != nil {
		// 数据库查询出错
		return err
	}

	// 2. 生成 UID
	userID := snowflake.GenID()
	u := &model.User{
		UserID:   userID,
		Username: param.Username,
		Password: param.Password,
		Email:    param.Email,
	}
	// 3. 保存进数据库
	return mysql.InsertUser(u)
}

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
			ctx.JSON(http.StatusOK, gin.H{
				"msg": err,
			})
			return
		}

		// validator.ValidateErrors 类型错误则使用翻译器
		// 并使用removeTopStruct函数去除字段名中的结构体名称标识
		ctx.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)),
		})
		return
	}
	// 2. 业务处理
	err = SignUp(req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 3. 返回响应
	ctx.JSON(http.StatusOK, "ok")
}
