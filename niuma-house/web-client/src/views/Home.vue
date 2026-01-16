<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getPosts, type Post } from '@/api/post'
import { getOccupations } from '@/api/user'

const router = useRouter()

const posts = ref<Post[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const occupationId = ref<number | undefined>()
const occupations = ref<{ id: number; name: string }[]>([])

const fetchPosts = async () => {
  loading.value = true
  try {
    const res = await getPosts({
      occupation_id: occupationId.value,
      page: currentPage.value,
      size: pageSize.value
    })
    posts.value = res.list || []
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const fetchOccupations = async () => {
  occupations.value = await getOccupations()
}

onMounted(() => {
  fetchOccupations()
  fetchPosts()
})

const handleFilterChange = () => {
  currentPage.value = 1
  fetchPosts()
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchPosts()
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('zh-CN')
}
</script>

<template>
  <div class="home-container">
    <!-- 横幅 -->
    <div class="banner">
      <h1>🐴 欢迎来到牛马之家</h1>
      <p>这里是职场牛马的避风港，分享工作经验，吐槽坑逼公司</p>
    </div>

    <!-- 筛选器 -->
    <div class="filter-bar">
      <el-select
        v-model="occupationId"
        placeholder="全部职业"
        clearable
        @change="handleFilterChange"
      >
        <el-option
          v-for="occ in occupations"
          :key="occ.id"
          :label="occ.name"
          :value="occ.id"
        />
      </el-select>
    </div>

    <!-- 帖子列表 -->
    <div class="post-list" v-loading="loading">
      <el-empty v-if="posts.length === 0" description="暂无帖子，快来发布第一篇吧！" />
      
      <div
        v-for="post in posts"
        :key="post.id"
        class="post-card hover-card"
        @click="router.push(`/post/${post.id}`)"
      >
        <div class="post-header">
          <div class="author-info">
            <el-avatar :size="40">{{ post.user?.username?.charAt(0) }}</el-avatar>
            <div class="author-detail">
              <span class="author-name">{{ post.user?.username }}</span>
              <span :class="['level-badge', `level-${post.user?.level}`]">
                Lv.{{ post.user?.level }}
              </span>
            </div>
          </div>
          <el-tag size="small">{{ post.occupation?.name }}</el-tag>
        </div>

        <h3 class="post-title">
          <el-tag v-if="post.status === 2" type="danger" size="small" class="top-tag">置顶</el-tag>
          {{ post.title }}
        </h3>
        <p class="post-excerpt">{{ post.content.substring(0, 150) }}...</p>

        <div class="post-footer">
          <span><el-icon><View /></el-icon> {{ post.views_count }}</span>
          <span><el-icon><Star /></el-icon> {{ post.likes_count }}</span>
          <span class="post-date">{{ formatDate(post.created_at) }}</span>
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
.home-container {
  max-width: 800px;
  margin: 0 auto;
}

.banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  padding: 48px 32px;
  border-radius: 16px;
  text-align: center;
  margin-bottom: 24px;
}

.banner h1 {
  font-size: 32px;
  margin-bottom: 12px;
}

.banner p {
  font-size: 16px;
  opacity: 0.9;
}

.filter-bar {
  margin-bottom: 24px;
}

.post-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.post-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.author-detail {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.author-name {
  font-weight: 500;
}

.post-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.top-tag {
  flex-shrink: 0;
}

.post-excerpt {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 12px;
}

.post-footer {
  display: flex;
  gap: 24px;
  color: #909399;
  font-size: 14px;
}

.post-footer span {
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
