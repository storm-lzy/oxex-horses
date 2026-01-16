<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getPosts, deletePost, topPost } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const posts = ref<any[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const fetchPosts = async () => {
  loading.value = true
  try {
    const res = await getPosts({ page: currentPage.value, size: pageSize.value })
    posts.value = res.list || []
    total.value = res.total
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchPosts())

const handleDelete = async (post: any) => {
  await ElMessageBox.confirm(`确定要删除帖子 "${post.title}" 吗？`, '确认')
  await deletePost(post.id)
  ElMessage.success('删除成功')
  fetchPosts()
}

const handleTop = async (post: any) => {
  await topPost(post.id)
  ElMessage.success('置顶成功')
  fetchPosts()
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchPosts()
}
</script>

<template>
  <div class="posts-page">
    <div class="page-header">
      <h2>帖子管理</h2>
    </div>

    <el-table :data="posts" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column label="作者">
        <template #default="{ row }">{{ row.user?.username }}</template>
      </el-table-column>
      <el-table-column label="职业">
        <template #default="{ row }">{{ row.occupation?.name }}</template>
      </el-table-column>
      <el-table-column prop="likes_count" label="点赞" width="80" />
      <el-table-column prop="views_count" label="浏览" width="80" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.status === 2" type="danger">置顶</el-tag>
          <el-tag v-else-if="row.status === 1" type="success">正常</el-tag>
          <el-tag v-else type="info">删除</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button v-if="row.status !== 2" type="warning" size="small" @click="handleTop(row)">置顶</el-button>
          <el-button v-if="row.status !== 0" type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-if="total > pageSize"
      v-model:current-page="currentPage"
      :page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="handlePageChange"
      style="margin-top: 16px"
    />
  </div>
</template>
