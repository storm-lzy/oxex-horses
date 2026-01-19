package database

import (
	"fmt"
	"log"
	"sync"
	"time"

	"niuma-house/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// InitMySQL 初始化 MySQL 连接
func InitMySQL(cfg *config.MySQLConfig) *gorm.DB {
	once.Do(func() {
		log.Printf("Connecting to MySQL: host=%s, port=%d, database=%s, username=%s",
			cfg.Host, cfg.Port, cfg.Database, cfg.Username)

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
			cfg.Charset,
		)

		// 重试连接，最多尝试 30 次，每次间隔 2 秒
		var err error
		for i := 0; i < 30; i++ {
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
			if err == nil {
				break
			}
			log.Printf("Failed to connect to MySQL (attempt %d/30): %v", i+1, err)
			time.Sleep(2 * time.Second)
		}

		if err != nil {
			log.Fatalf("Failed to connect to MySQL after 30 attempts: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get underlying sql.DB: %v", err)
		}

		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Hour)

		log.Println("MySQL connected successfully")
	})
	return db
}

// GetDB 获取数据库单例
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database not initialized. Call InitMySQL first.")
	}
	return db
}
