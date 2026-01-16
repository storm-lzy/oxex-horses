<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { createPost } from '@/api/post'
import { getOccupations } from '@/api/user'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { ElMessage } from 'element-plus'

const router = useRouter()

const form = ref({
  title: '',
  content: '',
  occupation_id: null as number | null
})
const loading = ref(false)
const occupations = ref<{ id: number; name: string }[]>([])

onMounted(async () => {
  occupations.value = await getOccupations()
})

const handleSubmit = async () => {
  if (!form.value.title || !form.value.content || !form.value.occupation_id) {
    ElMessage.warning('请填写所有必填项')
    return
  }

  loading.value = true
  try {
    const post = await createPost({
      title: form.value.title,
      content: form.value.content,
      occupation_id: form.value.occupation_id
    })
    ElMessage.success('发布成功！')
    router.push(`/post/${post.id}`)
  } catch (error) {
    // 错误已处理
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="create-post">
    <el-button @click="router.back()" class="back-btn">
      <el-icon><ArrowLeft /></el-icon> 返回
    </el-button>

    <div class="form-card">
      <h2>发布帖子</h2>

      <el-form @submit.prevent="handleSubmit" class="post-form">
        <el-form-item label="标题" required>
          <el-input
            v-model="form.title"
            placeholder="请输入帖子标题"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="职业分类" required>
          <el-select v-model="form.occupation_id" placeholder="选择职业分类">
            <el-option
              v-for="occ in occupations"
              :key="occ.id"
              :label="occ.name"
              :value="occ.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="内容" required>
          <MdEditor
            v-model="form.content"
            :preview="false"
            style="height: 400px"
            placeholder="支持 Markdown 格式..."
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" size="large" :loading="loading" @click="handleSubmit">
            发布
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.create-post {
  max-width: 900px;
  margin: 0 auto;
}

.back-btn {
  margin-bottom: 16px;
}

.form-card {
  background: #fff;
  border-radius: 12px;
  padding: 32px;
}

.form-card h2 {
  margin-bottom: 24px;
}

.post-form :deep(.el-form-item__label) {
  font-weight: 500;
}
</style>
