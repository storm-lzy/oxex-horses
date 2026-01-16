<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { getDashboardStats } from '@/api/admin'
import * as echarts from 'echarts'

const stats = ref<any>(null)
const loading = ref(false)

const pieChartRef = ref<HTMLElement | null>(null)
const lineChartRef = ref<HTMLElement | null>(null)

onMounted(async () => {
  loading.value = true
  try {
    stats.value = await getDashboardStats()
    await nextTick()
    initCharts()
  } finally {
    loading.value = false
  }
})

const initCharts = () => {
  // 职业分布饼图
  if (pieChartRef.value && stats.value?.occupation_stats) {
    const pieChart = echarts.init(pieChartRef.value)
    pieChart.setOption({
      title: { text: '职业分布', left: 'center' },
      tooltip: { trigger: 'item' },
      legend: { bottom: '5%', left: 'center' },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        data: stats.value.occupation_stats,
        emphasis: { itemStyle: { shadowBlur: 10 } }
      }]
    })
  }

  // 每日新增折线图
  if (lineChartRef.value && stats.value?.daily_users) {
    const lineChart = echarts.init(lineChartRef.value)
    lineChart.setOption({
      title: { text: '每日新增趋势 (近7天)', left: 'center' },
      tooltip: { trigger: 'axis' },
      legend: { top: 30 },
      xAxis: {
        type: 'category',
        data: stats.value.daily_users.map((d: any) => d.date)
      },
      yAxis: { type: 'value' },
      series: [
        { name: '新增用户', type: 'line', data: stats.value.daily_users.map((d: any) => d.count) },
        { name: '新增帖子', type: 'line', data: stats.value.daily_posts?.map((d: any) => d.count) || [] },
        { name: '新增公司', type: 'line', data: stats.value.daily_companies?.map((d: any) => d.count) || [] }
      ]
    })
  }
}
</script>

<template>
  <div class="dashboard" v-loading="loading">
    <el-row :gutter="24" class="stat-row">
      <el-col :span="6">
        <div class="stat-card">
          <div class="value">{{ stats?.totals?.users || 0 }}</div>
          <div class="label">总用户数</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="value" style="color: #67c23a">{{ stats?.totals?.posts || 0 }}</div>
          <div class="label">总帖子数</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="value" style="color: #f56c6c">{{ stats?.totals?.companies || 0 }}</div>
          <div class="label">避雷公司数</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="value" style="color: #e6a23c">{{ stats?.totals?.comments || 0 }}</div>
          <div class="label">总评论数</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="24">
      <el-col :span="12">
        <div class="chart-card">
          <div ref="pieChartRef" class="chart"></div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="chart-card">
          <div ref="lineChartRef" class="chart"></div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.stat-row {
  margin-bottom: 24px;
}
.chart-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
}
.chart {
  height: 400px;
}
</style>
