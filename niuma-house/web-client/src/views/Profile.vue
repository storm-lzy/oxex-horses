<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { updateProfile, getAvatarUploadUrl } from '@/api/user'

const userStore = useUserStore()

// 编辑状态
const isEditing = ref(false)
const editForm = ref({
  nickname: '',
  avatar: ''
})
const uploading = ref(false)
const saving = ref(false)

// 等级经验阈值
const levelThresholds = [
  { level: 1, name: '普通牛马', exp: 0 },
  { level: 2, name: '内卷牛马', exp: 100 },
  { level: 3, name: '精英牛马', exp: 500 },
  { level: 4, name: '天选牛马', exp: 2000 },
  { level: 5, name: '核动力牛马', exp: 10000 }
]

const getNextLevelExp = () => {
  const currentLevel = userStore.user?.level || 1
  if (currentLevel >= 5) return userStore.user?.exp || 0
  return levelThresholds[currentLevel]?.exp || 100
}

const getProgress = () => {
  const exp = userStore.user?.exp || 0
  const currentLevel = userStore.user?.level || 1
  const currentThreshold = levelThresholds[currentLevel - 1]?.exp || 0
  const nextThreshold = getNextLevelExp()
  if (currentLevel >= 5) return 100
  return Math.min(((exp - currentThreshold) / (nextThreshold - currentThreshold)) * 100, 100)
}

// 显示名称（优先昵称，否则用户名）
const displayName = computed(() => userStore.user?.nickname || userStore.user?.username || '')

// 头像显示
const avatarUrl = computed(() => userStore.user?.avatar || '')

// 开始编辑
const startEdit = () => {
  editForm.value = {
    nickname: userStore.user?.nickname || '',
    avatar: userStore.user?.avatar || ''
  }
  isEditing.value = true
}

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false
}

// 头像上传
const handleAvatarChange = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  // 检查文件大小 (最大 5MB)
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.error('图片大小不能超过 5MB')
    return
  }

  uploading.value = true
  try {
    // 获取预签名 URL
    const res = await getAvatarUploadUrl(file.name)
    
    // 上传到 MinIO
    await fetch(res.upload_url, {
      method: 'PUT',
      body: file,
      headers: {
        'Content-Type': file.type
      }
    })

    // 更新表单
    editForm.value.avatar = res.object_key
    ElMessage.success('头像上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('头像上传失败')
  } finally {
    uploading.value = false
  }
}

// 保存修改
const saveProfile = async () => {
  saving.value = true
  try {
    await updateProfile({
      nickname: editForm.value.nickname,
      avatar: editForm.value.avatar
    })
    await userStore.fetchProfile()
    ElMessage.success('保存成功')
    isEditing.value = false
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="profile-container">
    <div class="profile-card">
      <div class="profile-header">
        <!-- 头像区域 -->
        <div class="avatar-wrapper" :class="{ editing: isEditing }">
          <el-avatar :size="80" class="avatar" :src="avatarUrl || undefined">
            {{ displayName.charAt(0) }}
          </el-avatar>
          <div v-if="isEditing" class="avatar-overlay" @click="() => ($refs.avatarInput as HTMLInputElement)?.click()">
            <el-icon v-if="!uploading"><Camera /></el-icon>
            <el-icon v-else class="is-loading"><Loading /></el-icon>
          </div>
          <input
            ref="avatarInput"
            type="file"
            accept="image/jpeg,image/png,image/gif,image/webp"
            style="display: none"
            @change="handleAvatarChange"
          />
        </div>

        <div class="user-info">
          <!-- 编辑模式 -->
          <template v-if="isEditing">
            <el-input
              v-model="editForm.nickname"
              placeholder="输入昵称"
              maxlength="20"
              show-word-limit
              class="nickname-input"
            />
            <span class="username-hint">用户名: {{ userStore.user?.username }}</span>
          </template>
          <!-- 展示模式 -->
          <template v-else>
            <h2>{{ displayName }}</h2>
            <span :class="['level-badge', `level-${userStore.user?.level}`]">
              {{ userStore.levelName }}
            </span>
          </template>
        </div>

        <!-- 编辑按钮 -->
        <div class="action-buttons">
          <template v-if="isEditing">
            <el-button @click="cancelEdit" :disabled="saving">取消</el-button>
            <el-button type="primary" @click="saveProfile" :loading="saving">保存</el-button>
          </template>
          <template v-else>
            <el-button type="primary" @click="startEdit">
              <el-icon><Edit /></el-icon>
              编辑资料
            </el-button>
          </template>
        </div>
      </div>

      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-value">{{ userStore.user?.exp }}</div>
          <div class="stat-label">经验值</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">Lv.{{ userStore.user?.level }}</div>
          <div class="stat-label">等级</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ userStore.user?.occupation?.name }}</div>
          <div class="stat-label">职业</div>
        </div>
      </div>

      <div class="exp-progress">
        <div class="progress-header">
          <span>升级进度</span>
          <span>{{ userStore.user?.exp }} / {{ getNextLevelExp() }}</span>
        </div>
        <el-progress
          :percentage="getProgress()"
          :stroke-width="12"
          :show-text="false"
          status="success"
        />
      </div>

      <div class="level-table">
        <h3>等级说明</h3>
        <table>
          <thead>
            <tr>
              <th>等级</th>
              <th>称号</th>
              <th>所需经验</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="level in levelThresholds" :key="level.level" :class="{ active: userStore.user?.level === level.level }">
              <td>Lv.{{ level.level }}</td>
              <td>
                <span :class="['level-badge', `level-${level.level}`]">{{ level.name }}</span>
              </td>
              <td>{{ level.exp }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-container {
  max-width: 600px;
  margin: 0 auto;
}

.profile-card {
  background: #fff;
  border-radius: 16px;
  padding: 32px;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 32px;
}

.avatar-wrapper {
  position: relative;
}

.avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  font-size: 32px;
  color: #fff;
}

.avatar-wrapper.editing {
  cursor: pointer;
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 24px;
  opacity: 0;
  transition: opacity 0.2s;
}

.avatar-wrapper.editing:hover .avatar-overlay {
  opacity: 1;
}

.avatar-wrapper.editing .avatar-overlay.is-loading {
  opacity: 1;
}

.user-info {
  flex: 1;
}

.user-info h2 {
  font-size: 24px;
  margin-bottom: 8px;
}

.nickname-input {
  margin-bottom: 8px;
}

.username-hint {
  color: #909399;
  font-size: 12px;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 32px;
}

.stat-item {
  text-align: center;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 12px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #667eea;
}

.stat-label {
  color: #909399;
  font-size: 14px;
  margin-top: 4px;
}

.exp-progress {
  margin-bottom: 32px;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
  color: #606266;
}

.level-table h3 {
  margin-bottom: 16px;
}

.level-table table {
  width: 100%;
  border-collapse: collapse;
}

.level-table th, .level-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ebeef5;
}

.level-table tr.active {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
}
</style>
