package handler

import (
	"strconv"

	"niuma-house/internal/middleware"
	"niuma-house/internal/service"
	"niuma-house/pkg/response"

	"github.com/gin-gonic/gin"
)

var postService = service.NewPostService()

// GetPosts 获取帖子列表
func GetPosts(c *gin.Context) {
	occupationID, _ := strconv.ParseUint(c.Query("occupation_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	posts, total, err := postService.List(uint(occupationID), page, size)
	if err != nil {
		response.Fail(c, response.CodeServerError, "获取帖子列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  posts,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetPost 获取帖子详情
func GetPost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetCurrentUserID(c)

	post, isLiked, isFavorited, err := postService.GetByID(uint(id), userID)
	if err != nil {
		response.Fail(c, response.CodeNotFound, "帖子不存在")
		return
	}

	response.Success(c, gin.H{
		"post":         post,
		"is_liked":     isLiked,
		"is_favorited": isFavorited,
	})
}

// CreatePost 创建帖子
func CreatePost(c *gin.Context) {
	var req service.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetCurrentUserID(c)
	post, err := postService.Create(userID, &req)
	if err != nil {
		response.Fail(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, post)
}

// UpdatePost 更新帖子
func UpdatePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetCurrentUserID(c)

	var req service.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeInvalidParams, "参数错误")
		return
	}

	if err := postService.Update(uint(id), userID, &req); err != nil {
		response.Fail(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeletePost 删除帖子
func DeletePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetCurrentUserID(c)

	if err := postService.Delete(uint(id), userID); err != nil {
		response.Fail(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// LikePost 点赞帖子
func LikePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetCurrentUserID(c)

	if err := postService.Like(uint(id), userID); err != nil {
		response.Fail(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// UnlikePost 取消点赞
func UnlikePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetCurrentUserID(c)

	if err := postService.Unlike(uint(id), userID); err != nil {
		response.Fail(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// FavoritePost 收藏帖子
func FavoritePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetCurrentUserID(c)

	if err := postService.Favorite(uint(id), userID); err != nil {
		response.Fail(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// UnfavoritePost 取消收藏
func UnfavoritePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetCurrentUserID(c)

	if err := postService.Unfavorite(uint(id), userID); err != nil {
		response.Fail(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// AdminGetPosts 管理端获取帖子列表
func AdminGetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	posts, total, err := postService.AdminList(page, size)
	if err != nil {
		response.Fail(c, response.CodeServerError, "获取帖子列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  posts,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// AdminDeletePost 管理员删除帖子
func AdminDeletePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := postService.AdminDelete(uint(id)); err != nil {
		response.Fail(c, response.CodeServerError, "删除失败")
		return
	}
	response.Success(c, nil)
}

// TopPost 置顶帖子
func TopPost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := postService.SetTop(uint(id)); err != nil {
		response.Fail(c, response.CodeServerError, "置顶失败")
		return
	}
	response.Success(c, nil)
}
