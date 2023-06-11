package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取当前用户id
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	// 3. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 4. 返回响应
	ResponseSuccess(c, nil)

}

// GetPostDetailHandler 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	//1.获取参数（取出url中id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64) // 10进制，64位
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//fmt.Println(pid)
	//2.根据id取出帖子数据（查询数据库）
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	////3.返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 或缺帖子列表接口
func GetPostListHandler(c *gin.Context) {

	page, size := getPageInfo(c)
	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// GetPostListHandler2 根据时间或者分数获取帖子列表， 升级版
// 1. 获取参数
// 2. 去redis查询id列表
// 3. 根据id去数据库查询帖子详情
func GetPostListHandler2(c *gin.Context) {

	op, err := strconv.ParseInt(c.Query("community_id"), 10, 64)
	if op > 0 && err == nil { // 如果query包含community_id则走另外一个handler
		GetCommunityPostListHandler(c)
		return
	}

	p := &models.ParamPostList{
		Page:        1,
		Size:        10,
		Order:       models.OrderTime,
		CommunityID: 0,
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取数据
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func GetCommunityPostListHandler(c *gin.Context) {
	p := &models.ParamPostList{
		Page:        1,
		Size:        10,
		Order:       models.OrderTime,
		CommunityID: 0,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler with invalid params",
			zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	zap.L().Debug("debug param info", zap.Any("param", p))
	data, err := logic.GetCommunityPostListV2(p)
	if err != nil {
		zap.L().Error("logic.GetCommunityListV2() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
