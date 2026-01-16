package service

import (
	"errors"

	"niuma-house/internal/model"
	"niuma-house/internal/mq"
	"niuma-house/internal/repository"
)

// PostService 帖子服务
type PostService struct {
	postRepo *repository.PostRepository
	likeRepo *repository.LikeRepository
	favRepo  *repository.FavoriteRepository
}

// NewPostService 创建帖子服务
func NewPostService() *PostService {
	return &PostService{
		postRepo: repository.NewPostRepository(),
		likeRepo: repository.NewLikeRepository(),
		favRepo:  repository.NewFavoriteRepository(),
	}
}

// CreatePostRequest 创建帖子请求
type CreatePostRequest struct {
	Title        string `json:"title" binding:"required,max=200"`
	Content      string `json:"content" binding:"required"`
	OccupationID uint   `json:"occupation_id" binding:"required"`
}

// UpdatePostRequest 更新帖子请求
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"max=200"`
	Content string `json:"content"`
}

// Create 创建帖子
func (s *PostService) Create(userID uint, req *CreatePostRequest) (*model.Post, error) {
	post := &model.Post{
		UserID:       userID,
		OccupationID: req.OccupationID,
		Title:        req.Title,
		Content:      req.Content,
		Status:       1,
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	// 发送经验值消息
	mq.PublishExpMessage(userID, mq.ActionPost, 5)

	return post, nil
}

// GetByID 获取帖子详情
func (s *PostService) GetByID(id uint, userID uint) (*model.Post, bool, bool, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, false, false, err
	}

	// 增加浏览量
	s.postRepo.IncrementViews(id)

	// 检查是否点赞/收藏
	isLiked := s.likeRepo.IsLiked(id, userID)
	isFavorited := s.favRepo.IsFavorited(id, userID)

	return post, isLiked, isFavorited, nil
}

// Update 更新帖子
func (s *PostService) Update(id, userID uint, req *UpdatePostRequest) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return errors.New("无权编辑他人帖子")
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	return s.postRepo.Update(post)
}

// Delete 删除帖子
func (s *PostService) Delete(id, userID uint) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return errors.New("无权删除他人帖子")
	}

	return s.postRepo.Delete(id)
}

// List 帖子列表
func (s *PostService) List(occupationID uint, page, size int) ([]model.Post, int64, error) {
	return s.postRepo.List(occupationID, page, size)
}

// Like 点赞
func (s *PostService) Like(postID, userID uint) error {
	if s.likeRepo.IsLiked(postID, userID) {
		return errors.New("已点赞过该帖子")
	}

	if err := s.likeRepo.Like(postID, userID); err != nil {
		return err
	}

	s.postRepo.IncrementLikes(postID)

	// 获取帖子作者，给作者加经验
	post, err := s.postRepo.FindByID(postID)
	if err == nil && post.UserID != userID {
		mq.PublishExpMessage(post.UserID, mq.ActionLiked, 2)
	}

	return nil
}

// Unlike 取消点赞
func (s *PostService) Unlike(postID, userID uint) error {
	if !s.likeRepo.IsLiked(postID, userID) {
		return errors.New("未点赞过该帖子")
	}

	if err := s.likeRepo.Unlike(postID, userID); err != nil {
		return err
	}

	return s.postRepo.DecrementLikes(postID)
}

// Favorite 收藏
func (s *PostService) Favorite(postID, userID uint) error {
	if s.favRepo.IsFavorited(postID, userID) {
		return errors.New("已收藏过该帖子")
	}
	return s.favRepo.Favorite(postID, userID)
}

// Unfavorite 取消收藏
func (s *PostService) Unfavorite(postID, userID uint) error {
	if !s.favRepo.IsFavorited(postID, userID) {
		return errors.New("未收藏过该帖子")
	}
	return s.favRepo.Unfavorite(postID, userID)
}

// AdminList 管理端列表
func (s *PostService) AdminList(page, size int) ([]model.Post, int64, error) {
	return s.postRepo.AdminList(page, size)
}

// AdminDelete 管理员删除
func (s *PostService) AdminDelete(postID uint) error {
	return s.postRepo.Delete(postID)
}

// SetTop 置顶
func (s *PostService) SetTop(postID uint) error {
	return s.postRepo.SetTop(postID)
}
