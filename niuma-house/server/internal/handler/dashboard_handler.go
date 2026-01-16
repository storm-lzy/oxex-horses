package handler

import (
	"niuma-house/internal/model"
	"niuma-house/pkg/database"
	"niuma-house/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetDashboardStats 获取统计数据
func GetDashboardStats(c *gin.Context) {
	db := database.GetDB()

	// 职业分布
	var occupationStats []struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}
	db.Table("users").
		Select("occupations.name as name, COUNT(*) as value").
		Joins("LEFT JOIN occupations ON users.occupation_id = occupations.id").
		Where("users.deleted_at IS NULL").
		Group("occupations.id, occupations.name").
		Scan(&occupationStats)

	// 每日新增用户 (最近7天)
	var dailyUsers []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	db.Table("users").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyUsers)

	// 每日新增帖子 (最近7天)
	var dailyPosts []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	db.Table("posts").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyPosts)

	// 每日新增公司 (最近7天)
	var dailyCompanies []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	db.Table("companies").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyCompanies)

	// 总计统计
	var totalUsers, totalPosts, totalCompanies, totalComments int64
	db.Model(&model.User{}).Count(&totalUsers)
	db.Model(&model.Post{}).Where("status > 0").Count(&totalPosts)
	db.Model(&model.Company{}).Where("status > 0").Count(&totalCompanies)
	db.Model(&model.Comment{}).Where("status = 1").Count(&totalComments)

	response.Success(c, gin.H{
		"occupation_stats": occupationStats,
		"daily_users":      dailyUsers,
		"daily_posts":      dailyPosts,
		"daily_companies":  dailyCompanies,
		"totals": gin.H{
			"users":     totalUsers,
			"posts":     totalPosts,
			"companies": totalCompanies,
			"comments":  totalComments,
		},
	})
}
