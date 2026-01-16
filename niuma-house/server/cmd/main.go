package main

import (
	"flag"
	"fmt"
	"log"

	"niuma-house/internal/model"
	"niuma-house/internal/mq"
	"niuma-house/internal/router"
	"niuma-house/internal/task"
	"niuma-house/pkg/cache"
	"niuma-house/pkg/config"
	"niuma-house/pkg/database"
	"niuma-house/pkg/jwt"
	"niuma-house/pkg/queue"
	"niuma-house/pkg/storage"
)

func main() {
	// 命令行参数
	configPath := flag.String("config", "../config/config.yaml", "config file path")
	flag.Parse()

	// 加载配置
	cfg := config.LoadConfig(*configPath)
	log.Printf("Config loaded: server port=%d, mode=%s", cfg.Server.Port, cfg.Server.Mode)

	// 初始化 JWT
	jwt.Init(&cfg.JWT)

	// 初始化 MySQL
	db := database.InitMySQL(&cfg.MySQL)

	// 自动迁移
	if err := model.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化 Redis
	cache.InitRedis(&cfg.Redis)

	// 初始化 MinIO
	storage.InitMinIO(&cfg.MinIO)

	// 初始化 RabbitMQ
	queue.InitRabbitMQ(&cfg.RabbitMQ)
	defer queue.Close()

	// 启动 MQ 消费者
	go mq.StartExpConsumer()

	// 启动定时任务
	task.StartCronJobs()

	// 启动 HTTP 服务器
	r := router.SetupRouter(cfg)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
