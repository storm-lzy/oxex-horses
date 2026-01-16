<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getOccupations } from '@/api/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  username: '',
  password: '',
  confirmPassword: '',
  occupation_id: null as number | null
})
const loading = ref(false)
const occupations = ref<{ id: number; name: string }[]>([])

onMounted(async () => {
  occupations.value = await getOccupations()
})

const handleRegister = async () => {
  if (!form.value.username || !form.value.password || !form.value.occupation_id) {
    ElMessage.warning('请填写所有必填项')
    return
  }
  if (form.value.password !== form.value.confirmPassword) {
    ElMessage.warning('两次密码输入不一致')
    return
  }
  if (form.value.password.length < 6) {
    ElMessage.warning('密码长度不能少于6位')
    return
  }

  loading.value = true
  try {
    await userStore.registerAction({
      username: form.value.username,
      password: form.value.password,
      occupation_id: form.value.occupation_id
    })
    ElMessage.success('注册成功！请登录')
    router.push('/login')
  } catch (error) {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <span class="logo-icon">🐴</span>
        <h1>加入牛马之家</h1>
        <p>成为光荣的职场牛马！</p>
      </div>

      <el-form @submit.prevent="handleRegister" class="register-form">
        <el-form-item>
          <el-input
            v-model="form.username"
            placeholder="用户名 (3-20个字符)"
            size="large"
            prefix-icon="User"
          />
        </el-form-item>
        <el-form-item>
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码 (至少6位)"
            size="large"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="确认密码"
            size="large"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-select
            v-model="form.occupation_id"
            placeholder="选择你的职业"
            size="large"
            style="width: 100%"
          >
            <el-option
              v-for="occ in occupations"
              :key="occ.id"
              :label="occ.name"
              :value="occ.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleRegister"
            class="register-btn"
          >
            注册
          </el-button>
        </el-form-item>
      </el-form>

      <div class="register-footer">
        <span>已有账号？</span>
        <router-link to="/login">立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-card {
  background: #fff;
  border-radius: 16px;
  padding: 48px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo-icon {
  font-size: 48px;
}

.register-header h1 {
  font-size: 28px;
  margin: 16px 0 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.register-header p {
  color: #909399;
}

.register-form {
  margin-bottom: 24px;
}

.register-btn {
  width: 100%;
}

.register-footer {
  text-align: center;
  color: #909399;
}

.register-footer a {
  color: #667eea;
  font-weight: 500;
}
</style>
