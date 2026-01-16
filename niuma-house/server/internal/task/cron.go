package task

import (
	"log"

	"niuma-house/internal/model"
	"niuma-house/pkg/database"

	"github.com/robfig/cron/v3"
)

var cronScheduler *cron.Cron

// StartCronJobs 启动定时任务
func StartCronJobs() {
	cronScheduler = cron.New()

	// 每天凌晨 2:00 执行
	cronScheduler.AddFunc("0 2 * * *", dailyTask)

	// 每小时执行一次 - 清理过期数据
	cronScheduler.AddFunc("0 * * * *", hourlyCleanup)

	cronScheduler.Start()
	log.Println("Cron jobs started")
}

// StopCronJobs 停止定时任务
func StopCronJobs() {
	if cronScheduler != nil {
		cronScheduler.Stop()
	}
}

// dailyTask 每日任务
func dailyTask() {
	log.Println("Running daily task...")

	db := database.GetDB()

	// 重新计算所有用户等级（校准）
	var users []model.User
	db.Find(&users)

	for _, user := range users {
		newLevel := model.CalculateLevel(user.Exp)
		if newLevel != user.Level {
			db.Model(&user).Update("level", newLevel)
			log.Printf("User %d level calibrated: %d -> %d", user.ID, user.Level, newLevel)
		}
	}

	log.Println("Daily task completed")
}

// hourlyCleanup 每小时清理
func hourlyCleanup() {
	log.Println("Running hourly cleanup...")

	// 可以添加清理逻辑，如清理过期会话等

	log.Println("Hourly cleanup completed")
}
