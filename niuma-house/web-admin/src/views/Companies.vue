<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getCompanies, deleteCompany } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const companies = ref<any[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const fetchCompanies = async () => {
  loading.value = true
  try {
    const res = await getCompanies({ page: currentPage.value, size: pageSize.value })
    companies.value = res.list || []
    total.value = res.total
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchCompanies())

const handleDelete = async (company: any) => {
  await ElMessageBox.confirm(`确定要删除公司 "${company.name}" 吗？`, '确认')
  await deleteCompany(company.id)
  ElMessage.success('删除成功')
  fetchCompanies()
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchCompanies()
}
</script>

<template>
  <div class="companies-page">
    <div class="page-header">
      <h2>公司管理</h2>
    </div>

    <el-table :data="companies" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="公司名称" min-width="150" />
      <el-table-column prop="city" label="城市" width="100" />
      <el-table-column label="避雷等级" width="120">
        <template #default="{ row }">
          <span style="color: #f56c6c">{{ '⚠️'.repeat(row.risk_level) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="标签" min-width="200">
        <template #default="{ row }">
          <el-tag v-for="tag in (row.tags || []).slice(0, 3)" :key="tag" type="danger" size="small" style="margin-right: 4px">
            {{ tag }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="view_count" label="浏览" width="80" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '正常' : '删除' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
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
