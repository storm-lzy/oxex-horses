package repository

import (
	"niuma-house/internal/model"
	"niuma-house/pkg/database"

	"gorm.io/gorm"
)

// PostRepository 帖子仓储
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository 创建帖子仓储
func NewPostRepository() *PostRepository {
	return &PostRepository{db: database.GetDB()}
}

// Create 创建帖子
func (r *PostRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

// FindByID 根据 ID 查找帖子
func (r *PostRepository) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("User").Preload("Occupation").
		Where("status > 0").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// Update 更新帖子
func (r *PostRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

// Delete 删除帖子（软删除）
func (r *PostRepository) Delete(id uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).
		Update("status", 0).Error
}

// List 帖子列表
func (r *PostRepository) List(occupationID uint, page, size int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{}).Where("status > 0")
	if occupationID > 0 {
		query = query.Where("occupation_id = ?", occupationID)
	}

	query.Count(&total)

	offset := (page - 1) * size
	err := query.Preload("User").Preload("Occupation").
		Order("status DESC, created_at DESC"). // 置顶优先
		Offset(offset).Limit(size).
		Find(&posts).Error

	return posts, total, err
}

// IncrementLikes 增加点赞数
func (r *PostRepository) IncrementLikes(postID uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", postID).
		UpdateColumn("likes_count", gorm.Expr("likes_count + 1")).Error
}

// DecrementLikes 减少点赞数
func (r *PostRepository) DecrementLikes(postID uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", postID).
		UpdateColumn("likes_count", gorm.Expr("likes_count - 1")).Error
}

// IncrementViews 增加浏览数
func (r *PostRepository) IncrementViews(postID uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", postID).
		UpdateColumn("views_count", gorm.Expr("views_count + 1")).Error
}

// SetTop 置顶帖子
func (r *PostRepository) SetTop(postID uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", postID).
		Update("status", 2).Error
}

// AdminList 管理端列表（含已删除）
func (r *PostRepository) AdminList(page, size int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	r.db.Model(&model.Post{}).Count(&total)

	offset := (page - 1) * size
	err := r.db.Preload("User").Preload("Occupation").
		Order("created_at DESC").
		Offset(offset).Limit(size).
		Find(&posts).Error

	return posts, total, err
}

// LikeRepository 点赞仓储
type LikeRepository struct {
	db *gorm.DB
}

// NewLikeRepository 创建点赞仓储
func NewLikeRepository() *LikeRepository {
	return &LikeRepository{db: database.GetDB()}
}

// Like 点赞
func (r *LikeRepository) Like(postID, userID uint) error {
	like := model.PostLike{PostID: postID, UserID: userID}
	return r.db.Create(&like).Error
}

// Unlike 取消点赞
func (r *LikeRepository) Unlike(postID, userID uint) error {
	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&model.PostLike{}).Error
}

// IsLiked 是否已点赞
func (r *LikeRepository) IsLiked(postID, userID uint) bool {
	var count int64
	r.db.Model(&model.PostLike{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count)
	return count > 0
}

// FavoriteRepository 收藏仓储
type FavoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏仓储
func NewFavoriteRepository() *FavoriteRepository {
	return &FavoriteRepository{db: database.GetDB()}
}

// Favorite 收藏
func (r *FavoriteRepository) Favorite(postID, userID uint) error {
	fav := model.PostFavorite{PostID: postID, UserID: userID}
	return r.db.Create(&fav).Error
}

// Unfavorite 取消收藏
func (r *FavoriteRepository) Unfavorite(postID, userID uint) error {
	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&model.PostFavorite{}).Error
}

// IsFavorited 是否已收藏
func (r *FavoriteRepository) IsFavorited(postID, userID uint) bool {
	var count int64
	r.db.Model(&model.PostFavorite{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count)
	return count > 0
}
