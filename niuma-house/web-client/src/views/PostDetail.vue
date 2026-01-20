<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getPost, likePost, unlikePost, favoritePost, unfavoritePost, getComments, createComment, type Post } from '@/api/post'

import { useUserStore } from '@/stores/user'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const post = ref<Post | null>(null)
const isLiked = ref(false)
const isFavorited = ref(false)
const comments = ref<any[]>([])
const newComment = ref('')
const loading = ref(false)

const postId = computed(() => Number(route.params.id))

const fetchPost = async () => {
  loading.value = true
  try {
    const res = await getPost(postId.value)
    post.value = res.post
    isLiked.value = res.is_liked
    isFavorited.value = res.is_favorited
  } finally {
    loading.value = false
  }
}

const fetchComments = async () => {
  const res = await getComments(postId.value)
  comments.value = res.list || []
}

onMounted(() => {
  fetchPost()
  fetchComments()
})

const handleLike = async () => {
  if (isLiked.value) {
    await unlikePost(postId.value)
    isLiked.value = false
    if (post.value) post.value.likes_count--
  } else {
    await likePost(postId.value)
    isLiked.value = true
    if (post.value) post.value.likes_count++
  }
}

const handleFavorite = async () => {
  if (isFavorited.value) {
    await unfavoritePost(postId.value)
    isFavorited.value = false
    ElMessage.success('已取消收藏')
  } else {
    await favoritePost(postId.value)
    isFavorited.value = true
    ElMessage.success('收藏成功')
  }
}

const submitComment = async () => {
  if (!newComment.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }
  await createComment(postId.value, { content: newComment.value })
  newComment.value = ''
  ElMessage.success('评论成功')
  fetchComments()
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}
</script>

<template>
  <div class="post-detail" v-loading="loading">
    <el-button @click="router.back()" class="back-btn">
      <el-icon><ArrowLeft /></el-icon> 返回
    </el-button>

    <template v-if="post">
      <div class="post-card">
        <div class="post-header">
          <div class="author-info">
            <el-avatar :size="48" :src="post.user?.avatar || undefined">
              {{ (post.user?.nickname || post.user?.username)?.charAt(0) }}
            </el-avatar>
            <div class="author-detail">
              <span class="author-name">{{ post.user?.nickname || post.user?.username }}</span>
              <span :class="['level-badge', `level-${post.user?.level}`]">
                Lv.{{ post.user?.level }}
              </span>
              <span class="post-meta">{{ post.occupation?.name }} · {{ formatDate(post.created_at) }}</span>
            </div>
            <el-button
              v-if="userStore.isLoggedIn && post.user?.id !== userStore.user?.id"
              type="primary"
              size="small"
              @click="router.push({ path: '/messages', query: { userId: post.user?.id, username: post.user?.username } })"
            >
              <el-icon><ChatDotRound /></el-icon>
              私信
            </el-button>
          </div>
        </div>

        <h1 class="post-title">{{ post.title }}</h1>

        <div class="post-content">
          <MdPreview :modelValue="post.content" />
        </div>

        <div class="post-actions">
          <el-button :type="isLiked ? 'primary' : 'default'" @click="handleLike">
            <el-icon><Star /></el-icon>
            {{ isLiked ? '已点赞' : '点赞' }} {{ post.likes_count }}
          </el-button>
          <el-button :type="isFavorited ? 'warning' : 'default'" @click="handleFavorite">
            <el-icon><Collection /></el-icon>
            {{ isFavorited ? '已收藏' : '收藏' }}
          </el-button>
          <span class="view-count"><el-icon><View /></el-icon> {{ post.views_count }} 阅读</span>
        </div>
      </div>

      <!-- 评论区 -->
      <div class="comments-section">
        <h3>评论 ({{ comments.length }})</h3>

        <div class="comment-input" v-if="userStore.isLoggedIn">
          <el-input
            v-model="newComment"
            type="textarea"
            :rows="3"
            placeholder="写下你的评论..."
          />
          <el-button type="primary" @click="submitComment" class="submit-btn">发表评论</el-button>
        </div>

        <div class="comment-list">
          <div v-for="comment in comments" :key="comment.id" class="comment-item">
            <el-avatar :size="36">{{ comment.user?.username?.charAt(0) }}</el-avatar>
            <div class="comment-content">
              <div class="comment-header">
                <span class="comment-author">{{ comment.user?.username }}</span>
                <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
              </div>
              <p class="comment-text">{{ comment.content }}</p>
            </div>
          </div>
          <el-empty v-if="comments.length === 0" description="暂无评论" />
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.post-detail {
  max-width: 900px;
  margin: 0 auto;
}

.back-btn {
  margin-bottom: 16px;
}

.post-card {
  background: #fff;
  border-radius: 12px;
  padding: 32px;
  margin-bottom: 24px;
}

.post-header {
  margin-bottom: 24px;
}

.author-info {
  display: flex;
  gap: 12px;
}

.author-detail {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.author-name {
  font-weight: 600;
  font-size: 16px;
}

.post-meta {
  color: #909399;
  font-size: 14px;
}

.post-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 24px;
  line-height: 1.4;
}

.post-content {
  margin-bottom: 24px;
}

.post-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-top: 16px;
  border-top: 1px solid #ebeef5;
}

.view-count {
  color: #909399;
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: auto;
}

.comments-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
}

.comments-section h3 {
  margin-bottom: 16px;
}

.comment-input {
  margin-bottom: 24px;
}

.submit-btn {
  margin-top: 12px;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.comment-item {
  display: flex;
  gap: 12px;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.comment-author {
  font-weight: 500;
}

.comment-time {
  color: #909399;
  font-size: 12px;
}

.comment-text {
  color: #303133;
  line-height: 1.6;
}
</style>
