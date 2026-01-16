<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCompanies, searchCompanies, type Company } from '@/api/company'

const router = useRouter()

const companies = ref<Company[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const keyword = ref('')

const fetchCompanies = async () => {
  loading.value = true
  try {
    const res = keyword.value
      ? await searchCompanies({ keyword: keyword.value, page: currentPage.value, size: pageSize.value })
      : await getCompanies({ page: currentPage.value, size: pageSize.value })
    companies.value = res.list || []
    total.value = res.total
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCompanies()
})

const handleSearch = () => {
  currentPage.value = 1
  fetchCompanies()
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchCompanies()
}

const getRiskStars = (level: number) => {
  return '⚠️'.repeat(level)
}
</script>

<template>
  <div class="companies-container">
    <!-- 横幅 -->
    <div class="banner danger-banner">
      <h1>🚨 坑逼公司墙</h1>
      <p>揭露职场黑幕，让牛马不再被坑！</p>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-input
        v-model="keyword"
        placeholder="搜索公司名称..."
        size="large"
        @keyup.enter="handleSearch"
      >
        <template #append>
          <el-button @click="handleSearch">
            <el-icon><Search /></el-icon>
          </el-button>
        </template>
      </el-input>
      <el-button type="danger" size="large" @click="router.push('/company/create')">
        <el-icon><Plus /></el-icon>
        添加避雷公司
      </el-button>
    </div>

    <!-- 公司列表 -->
    <div class="company-list" v-loading="loading">
      <el-empty v-if="companies.length === 0" description="暂无数据" />

      <div
        v-for="company in companies"
        :key="company.id"
        class="company-card hover-card"
        @click="router.push(`/company/${company.id}`)"
      >
        <div class="company-header">
          <h3 class="company-name">{{ company.name }}</h3>
          <span class="risk-stars">{{ getRiskStars(company.risk_level) }}</span>
        </div>

        <div class="company-city" v-if="company.city">
          <el-icon><Location /></el-icon>
          {{ company.city }}
        </div>

        <div class="tag-list" v-if="company.tags?.length">
          <el-tag
            v-for="tag in company.tags.slice(0, 5)"
            :key="tag"
            type="danger"
            size="small"
          >
            {{ tag }}
          </el-tag>
        </div>

        <div class="company-footer">
          <span><el-icon><View /></el-icon> {{ company.view_count }}</span>
          <span class="post-date">{{ new Date(company.created_at).toLocaleDateString('zh-CN') }}</span>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <el-pagination
      v-if="total > pageSize"
      v-model:current-page="currentPage"
      :page-size="pageSize"
      :total="total"
      layout="prev, pager, next"
      @current-change="handlePageChange"
      class="pagination"
    />
  </div>
</template>

<style scoped>
.companies-container {
  max-width: 800px;
  margin: 0 auto;
}

.danger-banner {
  background: linear-gradient(135deg, #f56c6c 0%, #e6a23c 100%);
  color: #fff;
  padding: 48px 32px;
  border-radius: 16px;
  text-align: center;
  margin-bottom: 24px;
}

.danger-banner h1 {
  font-size: 32px;
  margin-bottom: 12px;
}

.search-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.search-bar .el-input {
  flex: 1;
}

.company-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.company-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  border-left: 4px solid #f56c6c;
}

.company-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.company-name {
  font-size: 18px;
  font-weight: 600;
}

.risk-stars {
  font-size: 16px;
}

.company-city {
  color: #909399;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 12px;
}

.tag-list {
  margin-bottom: 12px;
}

.company-footer {
  display: flex;
  gap: 24px;
  color: #909399;
  font-size: 14px;
}

.company-footer span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.post-date {
  margin-left: auto;
}

.pagination {
  margin-top: 24px;
  justify-content: center;
}
</style>
