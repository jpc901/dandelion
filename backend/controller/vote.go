package controller

import (
	"bluebell/logic"
	"bluebell/models"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

// PostVoteController 给帖子投票
// @Summary 给帖子投票
// @Description 描述
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param obj query models.ParamVoteData true "帖子投票参数"
// @Security ApiKeyAuth
// @Success 200
// @Router /api/v1/vote [post]
func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	// 获取当前用户id
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	if err = logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
