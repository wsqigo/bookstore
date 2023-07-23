package post

import (
	community2 "bookstore/web_app/community"
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

	// 根据 id 取出帖子数据（查数据库）
	post, err := GetPostByID(pid)
	if err != nil {
		zap.L().Error("GetPostByID failed", zap.Int64("pid", pid), zap.Error(err))
		controller.ResponseError(ctx, controller.CodeServerBusy)
		return
	}

	// 根据作者 id 查询作者信息
	u, err := user.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("GetUserById failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		controller.ResponseError(ctx, controller.CodeServerBusy)
		return
	}

	// 根据社区id查询社区详细信息
	community, err := community2.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("GetCommunityDetailByID failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		controller.ResponseError(ctx, controller.CodeServerBusy)
		return
	}

	res := ApiPostDetail{
		AuthorName:  u.Username,
		DBPost:      post,
		DBCommunity: community,
	}

	controller.ResponseSuccess(ctx, res)
}

func GetPostDetailList() ([]ApiPostDetail, error) {
	posts, err := GetPostList()
	if err != nil {
		return nil, err
	}

	data := make([]ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据作者id查询作者信息
		u, err := user.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("GetUserByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			return nil, err
		}

		// 根据社区id查询社区详细信息
		community, err := community2.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("GetCommunityDetailByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			return nil, err
		}

		data = append(data, ApiPostDetail{
			AuthorName:  u.Username,
			DBPost:      post,
			DBCommunity: community,
		})
	}

	return data, nil
}

// GetPostListHandler 获取帖子列表的处理函数
func GetPostListHandler(ctx *gin.Context) {

}
