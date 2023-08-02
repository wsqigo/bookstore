package post

import (
	"bookstore/web_app/community"
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

// CreatePostHandler 创建帖子
func CreatePostHandler(ctx *gin.Context) {
	// 1. 获取参数及参数的校验
	req := postReq{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Debug("ctx.ShouldBindJSON() err", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		controller.ResponseError(ctx, controller.CodeInvalidParam)
		return
	}

	// 2. 创建帖子
	// 从 ctx 获取当前发请求的用户的 user id
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
		zap.L().Error("GenAndInsertPost failed", zap.Error(err))
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
		zap.L().Error("get post detail with invalid param",
			zap.Error(err))
		controller.ResponseError(ctx, controller.CodeInvalidParam)
		return
	}

	post, err := GetPostByID(pid)
	if err != nil {
		zap.L().Error("GetPostByID failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		controller.ResponseError(ctx, controller.CodeServerBusy)
		return
	}

	// 根据作者 id 查询作者信息
	u, err := user.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("GetUserByID() failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		controller.ResponseError(ctx, controller.CodeServerBusy)
		return
	}

	// 根据社区 id 查询社区详细信息
	communityInfo, err := community.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("GetCommunityDetailByID() failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		controller.ResponseError(ctx, controller.CodeServerBusy)
		return
	}

	// 接口数据拼接
	data := ApiPostDetail{
		AuthorName:    u.Username,
		DBPost:        post,
		CommunityName: communityInfo.CommunityName,
	}
	controller.ResponseSuccess(ctx, data)
}

// 获取帖子列表
func getPostList(page, size int64) ([]ApiPostDetail, error) {
	// 获取数据
	posts, err := GetDBPostList(page, size)
	if err != nil {
		return nil, err
	}

	data := make([]ApiPostDetail, 0, len(posts))

	for _, p := range posts {
		// 根据作者 id 查询作者信息
		u, err := user.GetUserById(p.AuthorID)
		if err != nil {
			zap.L().Error("GetUserById failed",
				zap.Int64("author_id", p.AuthorID),
				zap.Error(err))
			return nil, err
		}

		// 根据社区 id 查询社区详细信息
		communityInfo, err := community.GetCommunityDetailByID(p.CommunityID)
		if err != nil {
			zap.L().Error("GetCommunityDetailByID failed",
				zap.Int64("community_id", p.CommunityID),
				zap.Error(err))
			return nil, err
		}

		postDetail := ApiPostDetail{
			DBPost:        p,
			AuthorName:    u.Username,
			CommunityName: communityInfo.CommunityName,
		}
		data = append(data, postDetail)
	}

	return data, nil
}

func getPageInfo(ctx *gin.Context) (int64, int64) {
	pageStr := ctx.Query("page")
	sizeStr := ctx.Query("size")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}

	return page, size
}

// GetPostListHandler 获取帖子列表的处理函数
func GetPostListHandler(ctx *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(ctx)
	// 获取数据
	data, err := getPostList(page, size)
	if err != nil {
		zap.L().Error("GetPostList() failed", zap.Error(err))
		controller.ResponseError(ctx, controller.CodeServerBusy)
		return
	}

	controller.ResponseSuccess(ctx, data)
}
