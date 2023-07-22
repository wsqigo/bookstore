package post

import (
	"bookstore/web_app/controller"
	"bookstore/web_app/user"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type postReq struct {
	AuthorID    int64  `json:"author_id"`
	CommunityID int64  `json:"community_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
}

func CreatePost(ctx *gin.Context) {
	// 1. 获取参数及参数的校验
	req := postReq{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		controller.ResponseError(ctx, controller.CodeInvalidParam)
		return
	}
	// 2. 创建帖子
	// 获取 user id
	userID, err := user.GetCurrentUserID(ctx)
	if err != nil {
		controller.ResponseError(ctx, controller.CodeNeedLogin)
		return
	}

	p := DBPost{
		AuthorID:    userID,
		CommunityID: req.CommunityID,
		Title:       req.Title,
		Content:     req.Content,
	}

	err = GenAndInsertPost(p)
	if err != nil {
		zap.L().Error("CreatePost failed", zap.Error(err))
		controller.ResponseError(ctx, controller.CodeServerBusy)
		return
	}

	// 3. 返回响应
	controller.ResponseSuccess(ctx, nil)
}

// GetPostDetailHandler 获取帖子详情的处理函数
func GetPostDetailHandler(ctx *gin.Context) {
	pidStr := ctx.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		controller.ResponseError(ctx, controller.CodeInvalidParam)
		return
	}

	data, err := GetPostByID(pid)
	if err != nil {
		zap.L().Error("GetPostByID failed", zap.Error(err))
		controller.ResponseError(ctx, controller.CodeServerBusy)
		return
	}

	// 根据 id 取出帖子数据（查数据库）
	controller.ResponseSuccess(ctx, data)

}
