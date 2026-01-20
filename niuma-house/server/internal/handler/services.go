package handler

import (
	"sync"

	"niuma-house/internal/service"
)

var (
	userSvc    *service.UserService
	postSvc    *service.PostService
	companySvc *service.CompanyService
	commentSvc *service.CommentService
	messageSvc *service.MessageService

	userOnce    sync.Once
	postOnce    sync.Once
	companyOnce sync.Once
	commentOnce sync.Once
	messageOnce sync.Once
)

// GetUserService 获取用户服务（懒加载）
func GetUserService() *service.UserService {
	userOnce.Do(func() {
		userSvc = service.NewUserService()
	})
	return userSvc
}

// GetPostService 获取帖子服务（懒加载）
func GetPostService() *service.PostService {
	postOnce.Do(func() {
		postSvc = service.NewPostService()
	})
	return postSvc
}

// GetCompanyService 获取公司服务（懒加载）
func GetCompanyService() *service.CompanyService {
	companyOnce.Do(func() {
		companySvc = service.NewCompanyService()
	})
	return companySvc
}

// GetCommentService 获取评论服务（懒加载）
func GetCommentService() *service.CommentService {
	commentOnce.Do(func() {
		commentSvc = service.NewCommentService()
	})
	return commentSvc
}

// GetMessageService 获取消息服务（懒加载）
func GetMessageService() *service.MessageService {
	messageOnce.Do(func() {
		messageSvc = service.NewMessageService()
	})
	return messageSvc
}
