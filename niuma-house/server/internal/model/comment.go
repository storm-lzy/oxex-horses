package model

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论实体
type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PostID    uint           `gorm:"not null;index" json:"post_id"`
	Post      *Post          `gorm:"foreignKey:PostID" json:"post,omitempty"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	ParentID  *uint          `gorm:"index" json:"parent_id,omitempty"` // 回复的评论ID
	Status    int            `gorm:"default:1" json:"status"`          // 1: 正常, 0: 删除
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Comment) TableName() string {
	return "comments"
}
