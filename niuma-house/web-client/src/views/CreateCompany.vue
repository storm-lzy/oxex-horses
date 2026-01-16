<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { createCompany } from '@/api/company'
import { ElMessage } from 'element-plus'

const router = useRouter()

const form = ref({
  name: '',
  city: '',
  tags: [] as string[],
  risk_level: 3,
  content: '',
  evidence: [] as string[]
})
const loading = ref(false)
const newTag = ref('')

const defaultTags = [
  '拖欠工资', '暴力裁员', '996严重', '单休', 'PUA文化',
  '不交社保', '领导傻逼', '加班无加班费', '画大饼', '钱少事多'
]

const addTag = () => {
  if (newTag.value && !form.value.tags.includes(newTag.value)) {
    form.value.tags.push(newTag.value)
    newTag.value = ''
  }
}

const removeTag = (tag: string) => {
  form.value.tags = form.value.tags.filter(t => t !== tag)
}

const selectTag = (tag: string) => {
  if (!form.value.tags.includes(tag)) {
    form.value.tags.push(tag)
  }
}

const handleSubmit = async () => {
  if (!form.value.name) {
    ElMessage.warning('请填写公司名称')
    return
  }
  if (form.value.tags.length === 0) {
    ElMessage.warning('请至少选择一个标签')
    return
  }

  loading.value = true
  try {
    const company = await createCompany(form.value)
    ElMessage.success('添加成功！')
    router.push(`/company/${company.id}`)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="create-company">
    <el-button @click="router.back()" class="back-btn">
      <el-icon><ArrowLeft /></el-icon> 返回
    </el-button>

    <div class="form-card">
      <h2>🚨 添加避雷公司</h2>

      <el-form @submit.prevent="handleSubmit" class="company-form" label-width="100px">
        <el-form-item label="公司名称" required>
          <el-input v-model="form.name" placeholder="请输入公司名称" maxlength="100" />
        </el-form-item>

        <el-form-item label="所在城市">
          <el-input v-model="form.city" placeholder="如：北京、上海" maxlength="50" />
        </el-form-item>

        <el-form-item label="避雷等级" required>
          <el-rate v-model="form.risk_level" :max="5" show-text :texts="['一般', '较差', '很差', '超级坑', '地狱级']" />
        </el-form-item>

        <el-form-item label="避雷标签" required>
          <div class="tag-section">
            <div class="selected-tags">
              <el-tag
                v-for="tag in form.tags"
                :key="tag"
                type="danger"
                closable
                @close="removeTag(tag)"
              >
                {{ tag }}
              </el-tag>
            </div>
            <div class="add-tag">
              <el-input v-model="newTag" placeholder="自定义标签" size="small" style="width: 150px" />
              <el-button size="small" @click="addTag">添加</el-button>
            </div>
            <div class="preset-tags">
              <span>快速选择：</span>
              <el-tag
                v-for="tag in defaultTags"
                :key="tag"
                :type="form.tags.includes(tag) ? 'danger' : 'info'"
                size="small"
                class="preset-tag"
                @click="selectTag(tag)"
              >
                {{ tag }}
              </el-tag>
            </div>
          </div>
        </el-form-item>

        <el-form-item label="详细描述">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="6"
            placeholder="描述一下这家公司的具体问题..."
          />
        </el-form-item>

        <el-form-item>
          <el-button type="danger" size="large" :loading="loading" @click="handleSubmit">
            提交避雷
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.create-company {
  max-width: 700px;
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

.tag-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.selected-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  min-height: 32px;
}

.add-tag {
  display: flex;
  gap: 8px;
}

.preset-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.preset-tags span {
  color: #909399;
  font-size: 14px;
}

.preset-tag {
  cursor: pointer;
}
</style>
