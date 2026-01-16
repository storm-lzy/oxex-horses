<script setup lang="ts">
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

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
</script>

<template>
  <div class="profile-container">
    <div class="profile-card">
      <div class="profile-header">
        <el-avatar :size="80" class="avatar">
          {{ userStore.user?.username?.charAt(0) }}
        </el-avatar>
        <div class="user-info">
          <h2>{{ userStore.user?.username }}</h2>
          <span :class="['level-badge', `level-${userStore.user?.level}`]">
            {{ userStore.levelName }}
          </span>
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

.avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  font-size: 32px;
  color: #fff;
}

.user-info h2 {
  font-size: 24px;
  margin-bottom: 8px;
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
