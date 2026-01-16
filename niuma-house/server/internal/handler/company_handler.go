package handler

import (
	"strconv"

	"niuma-house/internal/middleware"
	"niuma-house/internal/service"
	"niuma-house/pkg/response"

	"github.com/gin-gonic/gin"
)

var companyService = service.NewCompanyService()

// GetCompanies 获取公司列表
func GetCompanies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	companies, total, err := companyService.List(page, size)
	if err != nil {
		response.Fail(c, response.CodeServerError, "获取公司列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  companies,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// SearchCompanies 搜索公司
func SearchCompanies(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	companies, total, err := companyService.Search(c.Request.Context(), keyword, page, size)
	if err != nil {
		response.Fail(c, response.CodeServerError, "搜索失败")
		return
	}

	response.Success(c, gin.H{
		"list":    companies,
		"total":   total,
		"page":    page,
		"size":    size,
		"keyword": keyword,
	})
}

// GetCompany 获取公司详情
func GetCompany(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	company, err := companyService.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.CodeNotFound, "公司不存在")
		return
	}

	response.Success(c, company)
}

// CreateCompany 创建公司
func CreateCompany(c *gin.Context) {
	var req service.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetCurrentUserID(c)
	company, err := companyService.Create(userID, &req)
	if err != nil {
		response.Fail(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, company)
}

// AdminGetCompanies 管理端获取公司列表
func AdminGetCompanies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	companies, total, err := companyService.AdminList(page, size)
	if err != nil {
		response.Fail(c, response.CodeServerError, "获取公司列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  companies,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// AdminDeleteCompany 管理员删除公司
func AdminDeleteCompany(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := companyService.AdminDelete(uint(id)); err != nil {
		response.Fail(c, response.CodeServerError, "删除失败")
		return
	}
	response.Success(c, nil)
}
