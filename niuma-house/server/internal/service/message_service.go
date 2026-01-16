package service

import (
	"niuma-house/internal/model"
	"niuma-house/internal/repository"
)

// MessageService 私信服务
type MessageService struct {
	messageRepo *repository.MessageRepository
}

// NewMessageService 创建私信服务
func NewMessageService() *MessageService {
	return &MessageService{
		messageRepo: repository.NewMessageRepository(),
	}
}

// SendMessage 发送私信
func (s *MessageService) SendMessage(senderID, receiverID uint, content string) (*model.Message, error) {
	message := &model.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		IsRead:     false,
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}

	return message, nil
}

// GetConversations 获取会话列表
func (s *MessageService) GetConversations(userID uint) ([]model.Message, error) {
	return s.messageRepo.ListConversations(userID)
}

// GetMessages 获取与某用户的消息列表
func (s *MessageService) GetMessages(userID, otherUserID uint, page, size int) ([]model.Message, int64, error) {
	return s.messageRepo.ListByConversation(userID, otherUserID, page, size)
}

// MarkAsRead 标记消息已读
func (s *MessageService) MarkAsRead(receiverID, senderID uint) error {
	return s.messageRepo.MarkAsRead(receiverID, senderID)
}

// GetUnreadCount 获取未读消息数
func (s *MessageService) GetUnreadCount(userID uint) int64 {
	return s.messageRepo.CountUnread(userID)
}
