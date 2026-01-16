package service

import (
	"errors"

	"niuma-house/internal/model"
	"niuma-house/internal/mq"
	"niuma-house/internal/repository"
)

// CommentService 评论服务
type CommentService struct {
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository
}

// NewCommentService 创建评论服务
func NewCommentService() *CommentService {
	return &CommentService{
		commentRepo: repository.NewCommentRepository(),
		postRepo:    repository.NewPostRepository(),
	}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}

// Create 创建评论
func (s *CommentService) Create(postID, userID uint, req *CreateCommentRequest) (*model.Comment, error) {
	// 检查帖子是否存在
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, errors.New("帖子不存在")
	}

	comment := &model.Comment{
		PostID:   postID,
		UserID:   userID,
		Content:  req.Content,
		ParentID: req.ParentID,
		Status:   1,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// 给帖子作者加经验（自己评论自己的帖子不加）
	if post.UserID != userID {
		mq.PublishExpMessage(post.UserID, mq.ActionCommented, 1)
	}

	return comment, nil
}

// List 评论列表
func (s *CommentService) List(postID uint, page, size int) ([]model.Comment, int64, error) {
	return s.commentRepo.ListByPostID(postID, page, size)
}

// Delete 删除评论
func (s *CommentService) Delete(commentID, userID uint) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return errors.New("无权删除他人评论")
	}

	return s.commentRepo.Delete(commentID)
}
