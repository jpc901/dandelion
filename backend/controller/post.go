package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建帖子
// @Summary 创建帖子
// @Description 创建新帖子，存入数据库并在redis中记录该帖子的分数和所处社区
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT_AToken"
// @Param obj body models.Post false "参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseCreatePost
// @Router /api/v1/post [post]
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
// @Summary 通过post id获取post详情
// @Description 通过post id获取post内容以及所所在社区和作者名
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int64 true "帖子id"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostDetail
// @Router /api/v1/post/{id} [get]
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
// @Summary 概况
// @Description 描述
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param page path string false "页码"
// @Param size path string false "页面大小"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /api/v1/posts [post]
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
// @Summary 获取帖子分页数据
// @Description 根据社区id（可以为空）、页码、数量返回分页数据
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /api/v2/posts [get]
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
