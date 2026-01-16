<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

// 简化版私信页面
const selectedUser = ref<null | { id: number; username: string }>(null)
const messages = ref<any[]>([])
const newMessage = ref('')

const sendMessage = () => {
  if (!newMessage.value.trim()) {
    ElMessage.warning('请输入消息内容')
    return
  }
  // TODO: 实现 WebSocket 发送
  ElMessage.info('私信功能需要 WebSocket 连接，请确保后端已启动')
  newMessage.value = ''
}
</script>

<template>
  <div class="messages-container">
    <div class="messages-sidebar">
      <h3>会话列表</h3>
      <el-empty description="暂无会话" />
    </div>

    <div class="messages-main">
      <template v-if="selectedUser">
        <div class="chat-header">
          <h3>与 {{ selectedUser.username }} 的对话</h3>
        </div>
        <div class="chat-messages">
          <div v-for="msg in messages" :key="msg.id" class="message-item">
            {{ msg.content }}
          </div>
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
          <el-empty description="选择一个会话开始聊天" />
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
}

.messages-sidebar {
  width: 280px;
  border-right: 1px solid #ebeef5;
  padding: 16px;
}

.messages-sidebar h3 {
  margin-bottom: 16px;
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

.chat-messages {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
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
