<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const activeMenu = ref('/')

onMounted(() => {
  if (userStore.isLoggedIn) {
    userStore.fetchProfile()
  }
})

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

const handleCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      handleLogout()
      break
  }
}
</script>

<template>
  <el-container class="layout-container">
    <!-- 顶部导航 -->
    <el-header class="header">
      <div class="header-content">
        <div class="logo" @click="router.push('/')">
          <span class="logo-icon">🐴</span>
          <span class="logo-text">牛马之家</span>
        </div>

        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          :ellipsis="false"
          class="nav-menu"
          router
        >
          <el-menu-item index="/">首页</el-menu-item>
          <el-menu-item index="/companies">坑逼公司墙</el-menu-item>
          <el-menu-item index="/messages" v-if="userStore.isLoggedIn">私信</el-menu-item>
        </el-menu>

        <div class="header-right">
          <template v-if="userStore.isLoggedIn">
            <el-button type="primary" @click="router.push('/post/create')">
              <el-icon><Edit /></el-icon>
              发帖
            </el-button>
            <el-dropdown @command="handleCommand">
              <span class="user-info">
                <el-avatar :size="32" :src="userStore.user?.avatar || undefined">
                  {{ (userStore.user?.nickname || userStore.user?.username)?.charAt(0) }}
                </el-avatar>
                <span class="username">{{ userStore.user?.nickname || userStore.user?.username }}</span>
                <span :class="['level-badge', `level-${userStore.user?.level}`]">
                  {{ userStore.levelName }}
                </span>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button @click="router.push('/login')">登录</el-button>
            <el-button type="primary" @click="router.push('/register')">注册</el-button>
          </template>
        </div>
      </div>
    </el-header>

    <!-- 主内容区 -->
    <el-main class="main-content">
      <RouterView />
    </el-main>

    <!-- 底部 -->
    <el-footer class="footer">
      <p>© 2026 牛马之家 - 职场人的避风港</p>
      <p class="slogan">打工人打工魂，打工才是人上人 🔥</p>
    </el-footer>
  </el-container>
</template>

<style scoped>
.layout-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 0;
  height: 64px;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #fff;
}

.logo-icon {
  font-size: 28px;
  margin-right: 8px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
}

.nav-menu {
  background: transparent;
  border: none;
  flex: 1;
  margin-left: 48px;
}

.nav-menu :deep(.el-menu-item) {
  color: rgba(255, 255, 255, 0.85);
  border: none;
}

.nav-menu :deep(.el-menu-item:hover),
.nav-menu :deep(.el-menu-item.is-active) {
  color: #fff;
  background: rgba(255, 255, 255, 0.15);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: #fff;
}

.username {
  font-weight: 500;
}

.main-content {
  flex: 1;
  max-width: 1400px;
  width: 100%;
  margin: 0 auto;
  padding: 24px;
}

.footer {
  background: #303133;
  color: rgba(255, 255, 255, 0.7);
  text-align: center;
  padding: 24px;
  height: auto;
}

.footer .slogan {
  margin-top: 8px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
}
</style>
