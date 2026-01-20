package handler

import (
	"strconv"

	"niuma-house/internal/middleware"
	"niuma-house/internal/repository"
	"niuma-house/internal/service"
	"niuma-house/pkg/response"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	user, err := GetUserService().Register(&req)
	if err != nil {
		response.Fail(c, response.CodeUserExists, err.Error())
		return
	}

	response.Success(c, user)
}

// Login 用户登录
func Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeInvalidParams, "参数错误")
		return
	}

	resp, err := GetUserService().Login(&req)
	if err != nil {
		response.Fail(c, response.CodeWrongPassword, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetProfile 获取用户资料
func GetProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	user, err := GetUserService().GetProfile(userID)
	if err != nil {
		response.Fail(c, response.CodeNotFound, "用户不存在")
		return
	}

	response.Success(c, user)
}

// UpdateProfile 更新用户资料
func UpdateProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req struct {
		OccupationID uint `json:"occupation_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeInvalidParams, "参数错误")
		return
	}

	if err := GetUserService().UpdateProfile(userID, req.OccupationID); err != nil {
		response.Fail(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetOccupations 获取职业列表
func GetOccupations(c *gin.Context) {
	occupations, err := repository.NewOccupationRepository().List()
	if err != nil {
		response.Fail(c, response.CodeServerError, "获取职业列表失败")
		return
	}
	response.Success(c, occupations)
}

// AdminGetUsers 管理端获取用户列表
func AdminGetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	users, total, err := GetUserService().List(page, size)
	if err != nil {
		response.Fail(c, response.CodeServerError, "获取用户列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  users,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// BanUser 封禁用户
func BanUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := GetUserService().Ban(uint(id)); err != nil {
		response.Fail(c, response.CodeServerError, "封禁失败")
		return
	}
	response.Success(c, nil)
}

// UnbanUser 解封用户
func UnbanUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := GetUserService().Unban(uint(id)); err != nil {
		response.Fail(c, response.CodeServerError, "解封失败")
		return
	}
	response.Success(c, nil)
}
