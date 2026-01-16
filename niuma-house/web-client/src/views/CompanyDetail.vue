<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCompany, type Company } from '@/api/company'

const route = useRoute()
const router = useRouter()

const company = ref<Company | null>(null)
const loading = ref(false)

const companyId = computed(() => Number(route.params.id))

const fetchCompany = async () => {
  loading.value = true
  try {
    company.value = await getCompany(companyId.value)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCompany()
})

const getRiskStars = (level: number) => {
  return '⚠️'.repeat(level)
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}
</script>

<template>
  <div class="company-detail" v-loading="loading">
    <el-button @click="router.back()" class="back-btn">
      <el-icon><ArrowLeft /></el-icon> 返回
    </el-button>

    <template v-if="company">
      <div class="company-card">
        <div class="company-header">
          <h1>{{ company.name }}</h1>
          <span class="risk-stars">{{ getRiskStars(company.risk_level) }}</span>
        </div>

        <div class="company-meta">
          <span v-if="company.city"><el-icon><Location /></el-icon> {{ company.city }}</span>
          <span><el-icon><View /></el-icon> {{ company.view_count }} 次浏览</span>
          <span>{{ formatDate(company.created_at) }}</span>
        </div>

        <div class="tag-list" v-if="company.tags?.length">
          <el-tag
            v-for="tag in company.tags"
            :key="tag"
            type="danger"
            size="default"
          >
            {{ tag }}
          </el-tag>
        </div>

        <div class="company-content" v-if="company.content">
          <h3>详细描述</h3>
          <p>{{ company.content }}</p>
        </div>

        <div class="evidence-section" v-if="company.evidence?.length">
          <h3>证据截图</h3>
          <el-image
            v-for="(img, idx) in company.evidence"
            :key="idx"
            :src="img"
            :preview-src-list="company.evidence"
            fit="cover"
            class="evidence-img"
          />
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.company-detail {
  max-width: 900px;
  margin: 0 auto;
}

.back-btn {
  margin-bottom: 16px;
}

.company-card {
  background: #fff;
  border-radius: 12px;
  padding: 32px;
  border-left: 4px solid #f56c6c;
}

.company-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.company-header h1 {
  font-size: 28px;
  font-weight: 700;
}

.risk-stars {
  font-size: 24px;
}

.company-meta {
  display: flex;
  gap: 24px;
  color: #909399;
  margin-bottom: 16px;
}

.company-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.tag-list {
  margin-bottom: 24px;
}

.company-content {
  margin-bottom: 24px;
}

.company-content h3 {
  margin-bottom: 12px;
}

.company-content p {
  line-height: 1.8;
  color: #303133;
}

.evidence-section h3 {
  margin-bottom: 12px;
}

.evidence-img {
  width: 200px;
  height: 200px;
  border-radius: 8px;
  margin-right: 12px;
  margin-bottom: 12px;
}
</style>
