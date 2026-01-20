import request from '@/utils/request'

export interface Message {
  id: number
  sender_id: number
  receiver_id: number
  content: string
  is_read: boolean
  created_at: string
  sender?: {
    id: number
    username: string
    level: number
  }
  receiver?: {
    id: number
    username: string
    level: number
  }
}

export interface Conversation {
  user: {
    id: number
    username: string
    level: number
  }
  last_message: Message
  unread_count: number
}

// 获取会话列表
export const getConversations = () => {
  return request.get<{ list: Conversation[] }>('/messages')
}

// 获取与某用户的消息列表
export const getMessagesWith = (userId: number, page = 1, size = 50) => {
  return request.get<{ list: Message[]; total: number }>('/messages', {
    params: { user_id: userId, page, size }
  })
}

// 获取未读消息数
export const getUnreadCount = () => {
  return request.get<{ count: number }>('/messages/unread')
}

// 标记消息已读
export const markAsRead = (senderId: number) => {
  return request.post('/messages/read', { sender_id: senderId })
}
