package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户实体
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Nickname     string         `gorm:"size:50" json:"nickname"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	OccupationID uint           `gorm:"not null" json:"occupation_id"`
	Occupation   *Occupation    `gorm:"foreignKey:OccupationID" json:"occupation,omitempty"`
	Level        int            `gorm:"default:1" json:"level"`
	Exp          int            `gorm:"default:0" json:"exp"`
	Role         string         `gorm:"size:20;default:'user'" json:"role"` // user, admin, super_admin
	Status       int            `gorm:"default:1" json:"status"`            // 1: 正常, 0: 封禁
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子 - 密码加密
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// LevelName 等级名称
func (u *User) LevelName() string {
	return GetLevelName(u.Level)
}

// 等级阈值配置
var LevelThresholds = map[int]int{
	1: 0,     // 普通牛马
	2: 100,   // 内卷牛马
	3: 500,   // 精英牛马
	4: 2000,  // 天选牛马
	5: 10000, // 核动力牛马
}

// LevelNames 等级名称
var LevelNames = map[int]string{
	1: "普通牛马",
	2: "内卷牛马",
	3: "精英牛马",
	4: "天选牛马",
	5: "核动力牛马",
}

// GetLevelName 获取等级名称
func GetLevelName(level int) string {
	if name, ok := LevelNames[level]; ok {
		return name
	}
	return "普通牛马"
}

// CalculateLevel 根据经验值计算等级
func CalculateLevel(exp int) int {
	level := 1
	for l := 5; l >= 1; l-- {
		if exp >= LevelThresholds[l] {
			level = l
			break
		}
	}
	return level
}
