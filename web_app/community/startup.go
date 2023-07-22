package community

import (
	"bookstore/web_app/controller"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*
---- 跟社区相关的 ----
*/

func GetCommunityList() ([]DBCommunity, error) {
	// 查数据库 查找到所有的 community
	return getCommunityList()
}

func GetCommunityConf(ctx *gin.Context) {
	// 查询到所有的社区 (community_id, community_name)
	// 以列表的形式返回
	data, err := GetCommunityList()
	if err != nil {
		zap.L().Error("GetCommunityList() failed", zap.Error(err))
		controller.ResponseError(ctx, controller.CodeServerBusy) // 不轻易把服务端报错暴露给外面
		return
	}

	controller.ResponseSuccess(ctx, data)
}

// GetCommunityDetail 社区分类详情
func GetCommunityDetail(ctx *gin.Context) {
	// 1. 获取社区id
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		controller.ResponseError(ctx, controller.CodeInvalidParam)
		return
	}

	data, err := GetCommunityDetailByID(id)
	if err != nil {
		zap.L().Error("GetCommunityDetail failed", zap.Error(err))
		controller.ResponseError(ctx, controller.CodeServerBusy)
		return
	}

	controller.ResponseSuccess(ctx, data)
}
