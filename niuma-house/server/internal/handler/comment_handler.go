package handler

import (
	"strconv"

	"niuma-house/internal/middleware"
	"niuma-house/internal/service"
	"niuma-house/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetComments 获取评论列表
func GetComments(c *gin.Context) {
	postID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	comments, total, err := GetCommentService().List(uint(postID), page, size)
	if err != nil {
		response.Fail(c, response.CodeServerError, "获取评论列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  comments,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	postID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetCurrentUserID(c)

	var req service.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeInvalidParams, "参数错误")
		return
	}

	comment, err := GetCommentService().Create(uint(postID), userID, &req)
	if err != nil {
		response.Fail(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, comment)
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetCurrentUserID(c)

	if err := GetCommentService().Delete(uint(id), userID); err != nil {
		response.Fail(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
