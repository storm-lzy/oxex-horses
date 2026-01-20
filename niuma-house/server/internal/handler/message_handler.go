package handler

import (
	"strconv"

	"niuma-house/internal/middleware"
	"niuma-house/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetMessages 获取消息列表
func GetMessages(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	otherUserID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "50"))

	if otherUserID == 0 {
		// 获取会话列表
		conversations, err := GetMessageService().GetConversations(userID)
		if err != nil {
			response.Fail(c, response.CodeServerError, "获取会话列表失败")
			return
		}

		// 转换为前端期望的格式
		type ConversationItem struct {
			User struct {
				ID       uint   `json:"id"`
				Username string `json:"username"`
				Level    int    `json:"level"`
			} `json:"user"`
			LastMessage  interface{} `json:"last_message"`
			UnreadCount  int64       `json:"unread_count"`
		}

		var list []ConversationItem
		for _, msg := range conversations {
			var item ConversationItem

			// 确定对方用户
			if msg.SenderID == userID {
				item.User.ID = msg.ReceiverID
				if msg.Receiver != nil {
					item.User.Username = msg.Receiver.Username
					item.User.Level = msg.Receiver.Level
				}
			} else {
				item.User.ID = msg.SenderID
				if msg.Sender != nil {
					item.User.Username = msg.Sender.Username
					item.User.Level = msg.Sender.Level
				}
			}

			item.LastMessage = msg
			// 统计与该用户的未读消息数
			item.UnreadCount = GetMessageService().GetUnreadCountFrom(userID, item.User.ID)

			list = append(list, item)
		}

		response.Success(c, gin.H{"list": list})
		return
	}

	// 获取与某用户的消息
	messages, total, err := GetMessageService().GetMessages(userID, uint(otherUserID), page, size)
	if err != nil {
		response.Fail(c, response.CodeServerError, "获取消息列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  messages,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetUnreadCount 获取未读消息数
func GetUnreadCount(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	count := GetMessageService().GetUnreadCount(userID)
	response.Success(c, gin.H{"count": count})
}

// MarkAsRead 标记消息已读
func MarkAsRead(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req struct {
		SenderID uint `json:"sender_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeInvalidParams, "参数错误")
		return
	}

	if err := GetMessageService().MarkAsRead(userID, req.SenderID); err != nil {
		response.Fail(c, response.CodeServerError, "操作失败")
		return
	}

	response.Success(c, nil)
}
