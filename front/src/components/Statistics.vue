<template>
  <div class="statistics-page">
    <!-- 页头 -->
    <div class="page-header">
      <div class="page-header-left">
        <h2 class="page-header-title">
          <el-icon style="color: var(--med-accent);"><PieChart /></el-icon>
          检测统计分析
        </h2>
        <p class="page-header-sub">全面了解检测数据趋势与分布情况</p>
      </div>
    </div>

    <!-- 统计概览卡片 -->
    <div class="overview-grid">
      <div class="stat-card stat-card-blue">
        <div class="stat-card-icon">
          <el-icon :size="26" style="color: #2980B9;"><Document /></el-icon>
        </div>
        <div class="stat-card-body">
          <div class="stat-card-value">{{ totalTests }}</div>
          <div class="stat-card-label">总检测次数</div>
        </div>
        <div class="stat-card-wave"></div>
      </div>

      <div class="stat-card stat-card-red">
        <div class="stat-card-icon">
          <el-icon :size="26" style="color: #E74C3C;"><Warning /></el-icon>
        </div>
        <div class="stat-card-body">
          <div class="stat-card-value">{{ glandFaceCount }}</div>
          <div class="stat-card-label">腺体面容次数</div>
        </div>
        <div class="stat-card-wave"></div>
      </div>

      <div class="stat-card stat-card-green">
        <div class="stat-card-icon">
          <el-icon :size="26" style="color: #27AE60;"><Check /></el-icon>
        </div>
        <div class="stat-card-body">
          <div class="stat-card-value">{{ nonGlandFaceCount }}</div>
          <div class="stat-card-label">非腺体面容次数</div>
        </div>
        <div class="stat-card-wave"></div>
      </div>

      <div class="stat-card stat-card-orange">
        <div class="stat-card-icon">
          <el-icon :size="26" style="color: #F39C12;"><Star /></el-icon>
        </div>
        <div class="stat-card-body">
          <div class="stat-card-value">{{ accuracyRate }}%</div>
          <div class="stat-card-label">准确率</div>
        </div>
        <div class="stat-card-wave"></div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="charts-grid">
      <el-card class="chart-card">
        <template #header>
          <div class="chart-card-header trend-header">
            <el-icon style="color: var(--med-info);"><TrendCharts /></el-icon>
            <span>检测趋势</span>
            <div class="chart-legend">
              <span class="legend-item"><span class="legend-dot" style="background: #2980B9;"></span>总检测</span>
              <span class="legend-item"><span class="legend-dot" style="background: #E74C3C;"></span>腺体面容</span>
            </div>
          </div>
        </template>
        <div ref="trendChart" class="chart-container"></div>
      </el-card>

      <el-card class="chart-card">
        <template #header>
          <div class="chart-card-header level-header">
            <el-icon style="color: var(--med-ai-purple);"><DataAnalysis /></el-icon>
            <span>等级分布</span>
          </div>
        </template>
        <div ref="levelChart" class="chart-container"></div>
      </el-card>
    </div>

    <!-- 筛选条件 -->
    <el-card class="filter-card">
      <template #header>
        <div class="section-card-header">
          <el-icon style="color: var(--med-primary);"><Filter /></el-icon>
          <span>筛选条件</span>
        </div>
      </template>
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="filterForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 280px;"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadStatistics">
            <el-icon><Search /></el-icon> 查询
          </el-button>
          <el-button @click="resetFilter">
            <el-icon><Refresh /></el-icon> 重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 详细统计表格 -->
    <el-card class="table-card">
      <template #header>
        <div class="section-card-header">
          <el-icon style="color: var(--med-secondary);"><List /></el-icon>
          <span>每日详细统计</span>
        </div>
      </template>
      <el-table
        :data="statisticsDetail"
        border
        stripe
        style="width: 100%"
        :header-cell-style="{ background: '#EBF4F8', color: 'var(--med-text)', fontWeight: '600', fontSize: '14px' }"
      >
        <el-table-column prop="date" label="日期" align="center" min-width="120" />
        <el-table-column prop="total" label="总检测" align="center" min-width="100">
          <template #default="scope">
            <el-tag style="background: var(--med-info-light); border-color: rgba(41,128,185,0.25); color: var(--med-info);">
              {{ scope.row.total }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="glandFace" label="腺体面容" align="center" min-width="110">
          <template #default="scope">
            <el-tag type="danger">{{ scope.row.glandFace }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="nonGlandFace" label="非腺体面容" align="center" min-width="120">
          <template #default="scope">
            <el-tag type="success">{{ scope.row.nonGlandFace }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="light" label="轻度" align="center" min-width="90">
          <template #default="scope">
            <el-tag style="background: rgba(39,174,96,0.1); border-color: rgba(39,174,96,0.3); color: #27AE60;">{{ scope.row.light }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="moderate" label="中度" align="center" min-width="90">
          <template #default="scope">
            <el-tag type="warning">{{ scope.row.moderate }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="severe" label="重度" align="center" min-width="90">
          <template #default="scope">
            <el-tag type="danger">{{ scope.row.severe }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script>
import { ref, nextTick, onMounted } from 'vue'
import * as echarts from 'echarts'
import request from '../utils/request'
import { PieChart, Document, Warning, Check, Star, TrendCharts, DataAnalysis, Filter, Search, Refresh, List } from '@element-plus/icons-vue'

export default {
  name: 'Statistics',
  components: { PieChart, Document, Warning, Check, Star, TrendCharts, DataAnalysis, Filter, Search, Refresh, List },
  setup() {
    const totalTests = ref(0)
    const glandFaceCount = ref(0)
    const nonGlandFaceCount = ref(0)
    const accuracyRate = ref(0)
    const statisticsDetail = ref([])
    const trendChart = ref(null)
    const levelChart = ref(null)
    let trendChartInstance = null
    let levelChartInstance = null
    const filterForm = ref({ dateRange: [] })

    const loadStatistics = async () => {
      try {
        const params = {}
        if (filterForm.value.dateRange?.length === 2) {
          params.startDate = filterForm.value.dateRange[0]
          params.endDate = filterForm.value.dateRange[1]
        }
        const res = await request.get('/statistics/overview', { params })
        if (res.code === 1) {
          const d = res.data
          totalTests.value = d.totalTests || 0
          glandFaceCount.value = d.glandFaceCount || 0
          nonGlandFaceCount.value = d.nonGlandFaceCount || 0
          accuracyRate.value = d.accuracyRate ? (d.accuracyRate * 100).toFixed(2) : 0
          await loadDetailStatistics()
          await nextTick()
          renderTrendChart(d.trendData || [])
          renderLevelChart(d.levelData || [])
        }
      } catch (e) { console.error(e) }
    }

    const loadDetailStatistics = async () => {
      try {
        const params = {}
        if (filterForm.value.dateRange?.length === 2) {
          params.startDate = filterForm.value.dateRange[0]
          params.endDate = filterForm.value.dateRange[1]
        }
        const res = await request.get('/statistics/detail', { params })
        if (res.code === 1) statisticsDetail.value = res.data || []
      } catch (e) { console.error(e) }
    }

    const renderTrendChart = (data) => {
      if (!trendChart.value) return
      if (!trendChartInstance) trendChartInstance = echarts.init(trendChart.value)
      trendChartInstance.setOption({
        tooltip: { trigger: 'axis', backgroundColor: '#fff', borderColor: '#DEE5EE', textStyle: { color: '#2C3E50' } },
        legend: { data: ['总检测次数', '腺体面容次数'], bottom: 0, textStyle: { color: '#6B7E8E' } },
        grid: { top: 20, left: 40, right: 20, bottom: 40 },
        xAxis: { type: 'category', data: data.map(i => i.date), axisLine: { lineStyle: { color: '#DEE5EE' } }, axisLabel: { color: '#6B7E8E', fontSize: 12 } },
        yAxis: { type: 'value', axisLine: { show: false }, splitLine: { lineStyle: { color: '#EDF2F7' } }, axisLabel: { color: '#6B7E8E' } },
        series: [
          { name: '总检测次数', type: 'line', data: data.map(i => i.total), smooth: true, lineStyle: { color: '#2980B9', width: 3 }, itemStyle: { color: '#2980B9' }, areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(41,128,185,0.15)' }, { offset: 1, color: 'rgba(41,128,185,0)' }] } } },
          { name: '腺体面容次数', type: 'line', data: data.map(i => i.glandFace), smooth: true, lineStyle: { color: '#E74C3C', width: 3 }, itemStyle: { color: '#E74C3C' }, areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(231,76,60,0.12)' }, { offset: 1, color: 'rgba(231,76,60,0)' }] } } }
        ]
      })
    }

    const renderLevelChart = (data) => {
      if (!levelChart.value) return
      if (!levelChartInstance) levelChartInstance = echarts.init(levelChart.value)
      const colorMap = { '轻度': '#22c55e', '中期': '#eab308', '重度': '#ef4444', '正常': '#3b82f6', '非腺体面容': '#06b6d4', '非腺体': '#8b5cf6' }
      levelChartInstance.setOption({
        tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)', backgroundColor: '#fff', borderColor: '#DEE5EE', textStyle: { color: '#2C3E50' } },
        legend: { bottom: 0, textStyle: { color: '#6B7E8E' } },
        series: [{
          type: 'pie',
          radius: ['40%', '68%'],
          center: ['50%', '45%'],
          data: data.map(i => ({ name: i.level, value: i.count, itemStyle: { color: colorMap[i.level] || '#7C4DFF' } })),
          emphasis: { itemStyle: { shadowBlur: 12, shadowColor: 'rgba(0,0,0,0.15)' } },
          label: { color: '#2C3E50', fontSize: 12 },
          labelLine: { lineStyle: { color: '#DEE5EE' } }
        }]
      })
    }

    const resetFilter = () => { filterForm.value.dateRange = []; loadStatistics() }
    const handleResize = () => { trendChartInstance?.resize(); levelChartInstance?.resize() }

    onMounted(() => { loadStatistics(); window.addEventListener('resize', handleResize) })

    return { totalTests, glandFaceCount, nonGlandFaceCount, accuracyRate, statisticsDetail, trendChart, levelChart, filterForm, loadStatistics, resetFilter }
  }
}
</script>

<style scoped>
.statistics-page {
  padding: 24px;
  min-height: 100%;
  background: var(--med-bg);
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 页头 */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  background: linear-gradient(135deg, var(--med-sidebar) 0%, var(--med-sidebar-mid) 60%, #0C4A56 100%);
  border-radius: var(--med-radius-lg);
  box-shadow: var(--med-shadow-md);
}

.page-header-title {
  display: flex; align-items: center; gap: 10px;
  font-size: 22px; font-weight: 700; color: #fff; margin: 0 0 6px 0;
}

.page-header-sub { font-size: 13px; color: rgba(255,255,255,0.55); margin: 0; }

/* 概览卡片网格 */
.overview-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: #fff;
  border-radius: var(--med-radius-lg);
  box-shadow: var(--med-shadow-sm);
  position: relative;
  overflow: hidden;
  border: 1px solid var(--med-border-light);
  transition: all 0.25s ease;
}

.stat-card:hover { transform: translateY(-3px); box-shadow: var(--med-shadow-md); }

.stat-card-icon {
  width: 54px; height: 54px;
  border-radius: var(--med-radius-md);
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}

.stat-card-blue .stat-card-icon { background: var(--med-info-light); }
.stat-card-red .stat-card-icon { background: var(--med-danger-light); }
.stat-card-green .stat-card-icon { background: var(--med-success-light); }
.stat-card-orange .stat-card-icon { background: var(--med-warning-light); }

.stat-card-value { font-size: 28px; font-weight: 800; color: var(--med-text); line-height: 1; margin-bottom: 6px; }
.stat-card-label { font-size: 13px; color: var(--med-text-muted); }

.stat-card-wave {
  position: absolute; right: -10px; bottom: -10px;
  width: 70px; height: 70px; border-radius: 50%;
  opacity: 0.06;
}

.stat-card-blue .stat-card-wave { background: #2980B9; }
.stat-card-red .stat-card-wave { background: #E74C3C; }
.stat-card-green .stat-card-wave { background: #27AE60; }
.stat-card-orange .stat-card-wave { background: #F39C12; }

/* 图表网格 */
.charts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.chart-card :deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid var(--med-border-light);
}

.chart-card-header {
  display: flex; align-items: center; gap: 8px;
  font-size: 15px; font-weight: 700; color: var(--med-text);
}

.chart-legend {
  display: flex; gap: 16px; margin-left: auto;
}

.legend-item {
  display: flex; align-items: center; gap: 6px;
  font-size: 12px; color: var(--med-text-secondary);
}

.legend-dot { width: 10px; height: 10px; border-radius: 50%; display: inline-block; }

.chart-container { width: 100%; height: 320px; }

/* 筛选和表格卡片 */
.filter-card :deep(.el-card__header),
.table-card :deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid var(--med-border-light);
}

.section-card-header {
  display: flex; align-items: center; gap: 8px;
  font-size: 15px; font-weight: 700; color: var(--med-text);
}

.filter-form {
  display: flex; flex-wrap: wrap; gap: 16px; align-items: center;
  padding: 4px 0;
}

:deep(.el-table) { border-radius: var(--med-radius-md) !important; overflow: hidden; }
:deep(.el-table .el-table__cell) { padding: 12px 0; }
:deep(.el-table__row:hover > td) { background: var(--med-primary-bg) !important; }

@media (max-width: 1100px) {
  .overview-grid { grid-template-columns: repeat(2, 1fr); }
  .charts-grid { grid-template-columns: 1fr; }
}

@media (max-width: 640px) {
  .statistics-page { padding: 16px; gap: 14px; }
  .overview-grid { grid-template-columns: 1fr 1fr; gap: 12px; }
  .stat-card { padding: 16px; gap: 12px; }
  .stat-card-value { font-size: 22px; }
  .chart-container { height: 250px; }
}
</style>
