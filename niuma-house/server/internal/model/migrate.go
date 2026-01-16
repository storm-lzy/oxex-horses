package model

import (
	"log"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移所有模型
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&User{},
		&Occupation{},
		&Post{},
		&PostLike{},
		&PostFavorite{},
		&Company{},
		&Comment{},
		&Message{},
	)
	if err != nil {
		return err
	}

	// 初始化预置数据
	initOccupations(db)
	initAdminUser(db)

	log.Println("Database migration completed")
	return nil
}

// initOccupations 初始化职业分类
func initOccupations(db *gorm.DB) {
	for _, occ := range DefaultOccupations {
		var existing Occupation
		result := db.Where("id = ?", occ.ID).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			db.Create(&occ)
		}
	}
	log.Println("Default occupations initialized")
}

// initAdminUser 初始化管理员账号
func initAdminUser(db *gorm.DB) {
	var admin User
	result := db.Where("username = ?", "admin").First(&admin)
	if result.Error == gorm.ErrRecordNotFound {
		admin = User{
			Username:     "admin",
			Password:     "admin123", // 会在 BeforeCreate 中加密
			OccupationID: 1,
			Level:        5,
			Exp:          99999,
			Role:         "super_admin",
			Status:       1,
		}
		db.Create(&admin)
		log.Println("Admin user created: admin / admin123")
	}
}
