package cache

import (
	"context"
	"fmt"
	"log"
	"sync"

	"niuma-house/pkg/config"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
)

// InitRedis 初始化 Redis 连接
func InitRedis(cfg *config.RedisConfig) *redis.Client {
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		ctx := context.Background()
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}

		log.Println("Redis connected successfully")
	})
	return rdb
}

// GetRedis 获取 Redis 单例
func GetRedis() *redis.Client {
	if rdb == nil {
		log.Fatal("Redis not initialized. Call InitRedis first.")
	}
	return rdb
}
