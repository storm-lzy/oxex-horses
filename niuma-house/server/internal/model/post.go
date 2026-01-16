package model

import (
	"time"

	"gorm.io/gorm"
)

// Post 博客帖子实体
type Post struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	User         *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OccupationID uint           `gorm:"not null;index" json:"occupation_id"`
	Occupation   *Occupation    `gorm:"foreignKey:OccupationID" json:"occupation,omitempty"`
	Title        string         `gorm:"size:200;not null" json:"title"`
	Content      string         `gorm:"type:text;not null" json:"content"`
	LikesCount   int            `gorm:"default:0" json:"likes_count"`
	ViewsCount   int            `gorm:"default:0" json:"views_count"`
	Status       int            `gorm:"default:1;index" json:"status"` // 1: 正常, 0: 删除, 2: 置顶
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Post) TableName() string {
	return "posts"
}

// PostLike 帖子点赞记录
type PostLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null;uniqueIndex:idx_post_user" json:"post_id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_post_user" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (PostLike) TableName() string {
	return "post_likes"
}

// PostFavorite 帖子收藏记录
type PostFavorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null;uniqueIndex:idx_favorite_post_user" json:"post_id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_favorite_post_user" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (PostFavorite) TableName() string {
	return "post_favorites"
}
