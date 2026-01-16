package repository

import (
	"niuma-house/internal/model"
	"niuma-house/pkg/database"

	"gorm.io/gorm"
)

// CommentRepository 评论仓储
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓储
func NewCommentRepository() *CommentRepository {
	return &CommentRepository{db: database.GetDB()}
}

// Create 创建评论
func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// FindByID 根据 ID 查找评论
func (r *CommentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// Delete 删除评论
func (r *CommentRepository) Delete(id uint) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).
		Update("status", 0).Error
}

// ListByPostID 根据帖子 ID 获取评论列表
func (r *CommentRepository) ListByPostID(postID uint, page, size int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	r.db.Model(&model.Comment{}).
		Where("post_id = ? AND status = 1", postID).
		Count(&total)

	offset := (page - 1) * size
	err := r.db.Preload("User").
		Where("post_id = ? AND status = 1", postID).
		Order("created_at ASC").
		Offset(offset).Limit(size).
		Find(&comments).Error

	return comments, total, err
}

// CountByPostID 统计帖子评论数
func (r *CommentRepository) CountByPostID(postID uint) int64 {
	var count int64
	r.db.Model(&model.Comment{}).
		Where("post_id = ? AND status = 1", postID).
		Count(&count)
	return count
}
