package controller

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//--- 跟社区相关 ---

// CommunityHandler 获取所有社区信息
func CommunityHandler(c *gin.Context) {
	// 查询所有的社区 community_id community_name 以列表的形式返回

	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错暴露给外面
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区id
	idStr := c.Param("id")                     // 获取URL参数
	id, err := strconv.ParseInt(idStr, 10, 64) // 转换
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 查询社区分类详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错暴露给外面
	}
	ResponseSuccess(c, data)
}
