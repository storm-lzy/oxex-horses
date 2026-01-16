<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUsers, banUser, unbanUser } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const users = ref<any[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await getUsers({ page: currentPage.value, size: pageSize.value })
    users.value = res.list || []
    total.value = res.total
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchUsers())

const handleBan = async (user: any) => {
  await ElMessageBox.confirm(`确定要封禁用户 "${user.username}" 吗？`, '确认')
  await banUser(user.id)
  ElMessage.success('封禁成功')
  fetchUsers()
}

const handleUnban = async (user: any) => {
  await unbanUser(user.id)
  ElMessage.success('解封成功')
  fetchUsers()
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchUsers()
}

const levelNames: Record<number, string> = {
  1: '普通牛马', 2: '内卷牛马', 3: '精英牛马', 4: '天选牛马', 5: '核动力牛马'
}
</script>

<template>
  <div class="users-page">
    <div class="page-header">
      <h2>用户管理</h2>
    </div>

    <el-table :data="users" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column label="职业">
        <template #default="{ row }">{{ row.occupation?.name }}</template>
      </el-table-column>
      <el-table-column label="等级">
        <template #default="{ row }">
          Lv.{{ row.level }} {{ levelNames[row.level] || '' }}
        </template>
      </el-table-column>
      <el-table-column prop="exp" label="经验值" />
      <el-table-column label="状态">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '正常' : '封禁' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button v-if="row.status === 1" type="danger" size="small" @click="handleBan(row)">封禁</el-button>
          <el-button v-else type="success" size="small" @click="handleUnban(row)">解封</el-button>
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
