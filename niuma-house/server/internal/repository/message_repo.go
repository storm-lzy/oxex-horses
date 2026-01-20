package repository

import (
	"niuma-house/internal/model"
	"niuma-house/pkg/database"

	"gorm.io/gorm"
)

// MessageRepository 私信仓储
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository 创建私信仓储
func NewMessageRepository() *MessageRepository {
	return &MessageRepository{db: database.GetDB()}
}

// Create 创建私信
func (r *MessageRepository) Create(message *model.Message) error {
	return r.db.Create(message).Error
}

// FindByID 根据 ID 查找私信
func (r *MessageRepository) FindByID(id uint) (*model.Message, error) {
	var message model.Message
	err := r.db.Preload("Sender").Preload("Receiver").First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// ListConversations 获取用户的所有会话（最新消息）
func (r *MessageRepository) ListConversations(userID uint) ([]model.Message, error) {
	var messages []model.Message

	// 使用原生 SQL 获取与每个用户的最新消息
	err := r.db.Preload("Sender").Preload("Receiver").
		Where("id IN (SELECT MAX(id) FROM messages WHERE sender_id = ? OR receiver_id = ? GROUP BY IF(sender_id = ?, receiver_id, sender_id))", userID, userID, userID).
		Order("created_at DESC").
		Find(&messages).Error

	return messages, err
}

// ListByConversation 获取两个用户之间的消息
func (r *MessageRepository) ListByConversation(userID, otherUserID uint, page, size int) ([]model.Message, int64, error) {
	var messages []model.Message
	var total int64

	query := r.db.Model(&model.Message{}).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, otherUserID, otherUserID, userID)

	query.Count(&total)

	offset := (page - 1) * size
	err := query.Preload("Sender").Preload("Receiver").
		Order("created_at DESC").
		Offset(offset).Limit(size).
		Find(&messages).Error

	return messages, total, err
}

// MarkAsRead 标记消息为已读
func (r *MessageRepository) MarkAsRead(receiverID, senderID uint) error {
	return r.db.Model(&model.Message{}).
		Where("receiver_id = ? AND sender_id = ? AND is_read = false", receiverID, senderID).
		Update("is_read", true).Error
}

// CountUnread 统计未读消息数
func (r *MessageRepository) CountUnread(userID uint) int64 {
	var count int64
	r.db.Model(&model.Message{}).
		Where("receiver_id = ? AND is_read = false", userID).
		Count(&count)
	return count
}

// CountUnreadFrom 统计来自特定用户的未读消息数
func (r *MessageRepository) CountUnreadFrom(receiverID, senderID uint) int64 {
	var count int64
	r.db.Model(&model.Message{}).
		Where("receiver_id = ? AND sender_id = ? AND is_read = false", receiverID, senderID).
		Count(&count)
	return count
}
