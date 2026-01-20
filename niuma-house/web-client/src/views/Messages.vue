<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getConversations, getMessagesWith, markAsRead, type Message, type Conversation } from '@/api/message'

const route = useRoute()
const userStore = useUserStore()

// 会话列表
const conversations = ref<Conversation[]>([])
// 当前选中的用户
const selectedUser = ref<{ id: number; username: string } | null>(null)
// 消息列表
const messages = ref<Message[]>([])
// 新消息输入
const newMessage = ref('')
// WebSocket 连接
const ws = ref<WebSocket | null>(null)
// 消息容器引用
const messagesContainer = ref<HTMLElement | null>(null)
// 加载状态
const loading = ref(false)

// 从路由获取目标用户
onMounted(async () => {
  await fetchConversations()

  // 如果路由带有 userId 参数，直接选中该用户
  const userId = route.query.userId
  const username = route.query.username as string
  if (userId && username) {
    selectUser({ id: Number(userId), username })
  }

  // 建立 WebSocket 连接
  connectWebSocket()
})

onUnmounted(() => {
  if (ws.value) {
    ws.value.close()
  }
})

// 监听路由变化
watch(() => route.query, (newQuery: Record<string, string | (string | null)[] | null | undefined>) => {
  if (newQuery.userId && newQuery.username) {
    selectUser({ id: Number(newQuery.userId), username: newQuery.username as string })
  }
}, { immediate: false })

// 获取会话列表
const fetchConversations = async () => {
  try {
    const res = await getConversations()
    conversations.value = res.list || []
  } catch (error) {
    console.error('获取会话列表失败:', error)
  }
}

// 选择用户开始聊天
const selectUser = async (user: { id: number; username: string }) => {
  selectedUser.value = user
  await fetchMessages()
  await markAsRead(user.id)
  scrollToBottom()
}

// 获取与某用户的消息
const fetchMessages = async () => {
  if (!selectedUser.value) return
  loading.value = true
  try {
    const res = await getMessagesWith(selectedUser.value.id)
    messages.value = res.list || []
  } catch (error) {
    console.error('获取消息失败:', error)
  } finally {
    loading.value = false
  }
}

// 建立 WebSocket 连接
const connectWebSocket = () => {
  const token = localStorage.getItem('token')
  if (!token) return

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${protocol}//${window.location.hostname}:8080/api/ws/chat?token=${token}`
  
  ws.value = new WebSocket(wsUrl)

  ws.value.onopen = () => {
    console.log('WebSocket 连接成功')
  }

  ws.value.onmessage = (event: MessageEvent) => {
    try {
      const data = JSON.parse(event.data)
      // 收到新消息
      if (data.type === 'message') {
        const msg = data.data as Message
        // 如果是当前会话的消息，添加到列表
        if (selectedUser.value && 
            (msg.sender_id === selectedUser.value.id || msg.receiver_id === selectedUser.value.id)) {
          messages.value.push(msg)
          scrollToBottom()
          // 标记已读
          if (msg.sender_id === selectedUser.value.id) {
            markAsRead(selectedUser.value.id)
          }
        }
        // 刷新会话列表
        fetchConversations()
      }
    } catch (error) {
      console.error('解析 WebSocket 消息失败:', error)
    }
  }

  ws.value.onerror = (error: Event) => {
    console.error('WebSocket 错误:', error)
  }

  ws.value.onclose = () => {
    console.log('WebSocket 连接关闭')
    // 5秒后尝试重连
    setTimeout(() => {
      if (userStore.isLoggedIn) {
        connectWebSocket()
      }
    }, 5000)
  }
}

// 发送消息
const sendMessage = () => {
  if (!newMessage.value.trim()) {
    ElMessage.warning('请输入消息内容')
    return
  }

  if (!selectedUser.value) {
    ElMessage.warning('请先选择一个聊天对象')
    return
  }

  if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
    ElMessage.error('连接已断开，正在重连...')
    connectWebSocket()
    return
  }

  // 发送消息
  const messageData = {
    type: 'message',
    receiver_id: selectedUser.value.id,
    content: newMessage.value.trim()
  }
  ws.value.send(JSON.stringify(messageData))

  // 本地立即显示消息
  messages.value.push({
    id: Date.now(),
    sender_id: userStore.user!.id,
    receiver_id: selectedUser.value.id,
    content: newMessage.value.trim(),
    is_read: false,
    created_at: new Date().toISOString()
  })

  newMessage.value = ''
  scrollToBottom()
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// 格式化时间
const formatTime = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const isToday = date.toDateString() === now.toDateString()
  
  if (isToday) {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

// 判断是否是自己发的消息
const isOwnMessage = (msg: Message) => {
  return msg.sender_id === userStore.user?.id
}
</script>

<template>
  <div class="messages-container">
    <!-- 会话列表 -->
    <div class="messages-sidebar">
      <h3>会话列表</h3>
      <div class="conversation-list" v-if="conversations.length > 0">
        <div
          v-for="conv in conversations"
          :key="conv.user.id"
          :class="['conversation-item', { active: selectedUser?.id === conv.user.id }]"
          @click="selectUser(conv.user)"
        >
          <el-avatar :size="40">{{ conv.user.username.charAt(0) }}</el-avatar>
          <div class="conversation-info">
            <div class="conversation-header">
              <span class="conversation-name">{{ conv.user.username }}</span>
              <span v-if="conv.unread_count > 0" class="unread-badge">{{ conv.unread_count }}</span>
            </div>
            <p class="last-message">{{ conv.last_message?.content }}</p>
          </div>
        </div>
      </div>
      <el-empty v-else description="暂无会话" />
    </div>

    <!-- 聊天区域 -->
    <div class="messages-main">
      <template v-if="selectedUser">
        <div class="chat-header">
          <h3>与 {{ selectedUser.username }} 的对话</h3>
        </div>
        <div class="chat-messages" ref="messagesContainer" v-loading="loading">
          <template v-if="messages.length > 0">
            <div
              v-for="msg in messages"
              :key="msg.id"
              :class="['message-item', { own: isOwnMessage(msg) }]"
            >
              <el-avatar :size="32" v-if="!isOwnMessage(msg)">
                {{ selectedUser.username.charAt(0) }}
              </el-avatar>
              <div class="message-bubble">
                <p class="message-content">{{ msg.content }}</p>
                <span class="message-time">{{ formatTime(msg.created_at) }}</span>
              </div>
              <el-avatar :size="32" v-if="isOwnMessage(msg)">
                {{ userStore.user?.username?.charAt(0) }}
              </el-avatar>
            </div>
          </template>
          <el-empty v-else description="开始聊天吧" />
        </div>
        <div class="chat-input">
          <el-input
            v-model="newMessage"
            placeholder="输入消息..."
            @keyup.enter="sendMessage"
          />
          <el-button type="primary" @click="sendMessage">发送</el-button>
        </div>
      </template>
      <template v-else>
        <div class="no-chat">
          <el-empty description="选择一个会话开始聊天，或从帖子中点击作者的私信按钮发起对话" />
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.messages-container {
  display: flex;
  height: calc(100vh - 200px);
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.messages-sidebar {
  width: 300px;
  border-right: 1px solid #ebeef5;
  display: flex;
  flex-direction: column;
}

.messages-sidebar h3 {
  padding: 16px;
  margin: 0;
  border-bottom: 1px solid #ebeef5;
}

.conversation-list {
  flex: 1;
  overflow-y: auto;
}

.conversation-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.conversation-item:hover {
  background: #f5f7fa;
}

.conversation-item.active {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
}

.conversation-info {
  flex: 1;
  min-width: 0;
}

.conversation-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}

.conversation-name {
  font-weight: 500;
}

.unread-badge {
  background: #f56c6c;
  color: #fff;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
}

.last-message {
  margin: 0;
  font-size: 13px;
  color: #909399;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.messages-main {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.chat-header {
  padding: 16px;
  border-bottom: 1px solid #ebeef5;
}

.chat-header h3 {
  margin: 0;
}

.chat-messages {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  max-width: 70%;
}

.message-item.own {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.message-bubble {
  background: #f5f7fa;
  padding: 10px 14px;
  border-radius: 12px;
  max-width: 100%;
}

.message-item.own .message-bubble {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}

.message-content {
  margin: 0;
  word-wrap: break-word;
  line-height: 1.5;
}

.message-time {
  display: block;
  font-size: 11px;
  color: #909399;
  margin-top: 4px;
  text-align: right;
}

.message-item.own .message-time {
  color: rgba(255, 255, 255, 0.7);
}

.chat-input {
  display: flex;
  gap: 12px;
  padding: 16px;
  border-top: 1px solid #ebeef5;
}

.chat-input .el-input {
  flex: 1;
}

.no-chat {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
