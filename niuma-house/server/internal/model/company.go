package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// StringArray 字符串数组类型 (JSON存储)
type StringArray []string

// Scan 实现 sql.Scanner 接口
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal StringArray value")
	}
	return json.Unmarshal(bytes, s)
}

// Value 实现 driver.Valuer 接口
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Company 坑逼公司实体
type Company struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;index" json:"name"`
	City      string         `gorm:"size:50" json:"city"`
	Tags      StringArray    `gorm:"type:json" json:"tags"`       // ["拖欠工资", "暴力裁员"]
	RiskLevel int            `gorm:"default:1" json:"risk_level"` // 1-5 星避雷等级
	Evidence  StringArray    `gorm:"type:json" json:"evidence"`   // 证据图片 MinIO Keys
	Content   string         `gorm:"type:text" json:"content"`    // 详细描述
	CreatorID uint           `gorm:"not null" json:"creator_id"`
	Creator   *User          `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Status    int            `gorm:"default:1;index" json:"status"` // 1: 正常, 0: 删除
	ViewCount int            `gorm:"default:0" json:"view_count"`
	CreatedAt time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Company) TableName() string {
	return "companies"
}

// 预置避雷标签
var DefaultCompanyTags = []string{
	"拖欠工资",
	"暴力裁员",
	"996严重",
	"单休",
	"PUA文化",
	"不交社保",
	"领导傻逼",
	"加班无加班费",
	"画大饼",
	"钱少事多",
}
