<template>
  <div class="test-result-page">
    <!-- 页头 -->
    <div class="page-header">
      <div class="page-header-left">
        <h2 class="page-header-title">
          <el-icon style="color: var(--med-accent);"><Document /></el-icon>
          检测记录
        </h2>
        <p class="page-header-sub">查看并管理所有历史面容检测记录</p>
      </div>
      <div class="header-actions">
        <el-button
          type="danger"
          :disabled="selectedResults.length === 0"
          @click="deleteSelectedResults"
          v-if="resultList.length > 0"
          plain
        >
          <el-icon><Delete /></el-icon> 批量删除 ({{ selectedResults.length }})
        </el-button>
        <el-button type="primary" @click="downloadResults" :loading="downloading">
          <el-icon><Download /></el-icon> 下载 Excel
        </el-button>
        <el-button @click="downloadResultsAsPDF" :loading="downloadingPDF" style="border-color: var(--med-success); color: var(--med-success);">
          <el-icon><Document /></el-icon> 下载 PDF
        </el-button>
      </div>
    </div>

    <!-- 筛选卡片 -->
    <el-card class="filter-card">
      <template #header>
        <div class="section-card-header">
          <el-icon style="color: var(--med-primary);"><Filter /></el-icon>
          <span>筛选条件</span>
        </div>
      </template>
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="开始日期">
          <el-date-picker
            v-model="filterForm.startDate"
            type="datetime"
            placeholder="选择开始日期"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 210px;"
          />
        </el-form-item>
        <el-form-item label="结束日期">
          <el-date-picker
            v-model="filterForm.endDate"
            type="datetime"
            placeholder="选择结束日期"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 210px;"
          />
        </el-form-item>
        <el-form-item label="腺体面容">
          <el-select v-model="filterForm.isGlandFace" clearable placeholder="请选择" style="width: 110px;">
            <el-option label="是" :value="true" />
            <el-option label="否" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item label="等级">
          <el-select v-model="filterForm.level" clearable placeholder="请选择等级" style="width: 130px;">
            <el-option label="轻度" value="轻度" />
            <el-option label="中期" value="中期" />
            <el-option label="重度" value="重度" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="searchResults">
            <el-icon><Search /></el-icon> 查询
          </el-button>
          <el-button @click="resetFilter">
            <el-icon><Refresh /></el-icon> 重置
          </el-button>
          <el-button :loading="downloading" @click="downloadResults" style="border-color: var(--med-info); color: var(--med-info);">
            <el-icon><Download /></el-icon> 按条件下载
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 记录表格 -->
    <el-card class="table-card">
      <template #header>
        <div class="section-card-header">
          <el-icon style="color: var(--med-secondary);"><List /></el-icon>
          <span>检测记录列表</span>
          <div class="table-total">共 <strong>{{ total }}</strong> 条记录</div>
        </div>
      </template>

      <el-table
        :data="resultList"
        v-loading="loading"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        border
        stripe
        highlight-current-row
        :header-cell-style="{ background: '#EBF4F8', color: 'var(--med-text)', fontWeight: '600', fontSize: '14px' }"
        class="med-table"
      >
        <el-table-column type="selection" min-width="50" align="center" />
        <el-table-column label="序号" min-width="60" align="center">
          <template #default="scope">
            <span class="seq-num">{{ (currentPage - 1) * pageSize + scope.$index + 1 }}</span>
          </template>
        </el-table-column>

        <el-table-column label="检测图片" min-width="110" align="center">
          <template #default="scope">
            <el-image
              :src="getImageUrl(scope.row.imagePath)"
              fit="cover"
              class="record-img"
              :preview-src-list="[getImageUrl(scope.row.imagePath)]"
              preview-teleported
            >
              <template #error>
                <div class="img-placeholder"><el-icon><Picture /></el-icon></div>
              </template>
            </el-image>
          </template>
        </el-table-column>

        <el-table-column prop="isGlandFace" label="腺体面容" min-width="100" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.isGlandFace ? 'danger' : 'success'" round>
              {{ scope.row.isGlandFace ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="level" label="严重程度" min-width="100" align="center">
          <template #default="scope">
            <el-tag :type="getLevelType(scope.row.level)" round>{{ scope.row.level }}</el-tag>
          </template>
        </el-table-column>


        <el-table-column prop="createTime" label="检测时间" min-width="170" align="center" sortable>
          <template #default="scope">
            <span class="time-cell">{{ scope.row.createTime || scope.row.testTime }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" align="center" min-width="140" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" @click="openEditDialog(scope.row)" circle plain>
              <el-icon><Edit /></el-icon>
            </el-button>
            <el-button size="small" type="danger" @click="deleteResult(scope.row.id)" circle plain>
              <el-icon><Delete /></el-icon>
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 分页 -->
    <div class="pagination-wrap">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="currentPage"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        background
      />
    </div>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑检测结果"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="100px">
        <el-form-item label="腺体面容" prop="isGlandFace">
          <el-radio-group v-model="editForm.isGlandFace">
            <el-radio :label="true">是</el-radio>
            <el-radio :label="false">否</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="严重程度" prop="level">
          <el-select v-model="editForm.level" placeholder="请选择严重程度" style="width: 100%;">
            <el-option label="轻度" value="轻度" />
            <el-option label="中度" value="中度" />
            <el-option label="重度" value="重度" />
            <el-option label="非腺样体面容" value="非腺样体面容" />
          </el-select>
        </el-form-item>
        <el-form-item label="可视化描述" prop="visualizationDescription">
          <el-input 
            v-model="editForm.visualizationDescription" 
            type="textarea" 
            :rows="4"
            placeholder="请输入可视化描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEdit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import request, { downloadRequest } from '../utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, Delete, Download, Filter, Search, Refresh, List, Picture, Edit } from '@element-plus/icons-vue'

export default {
  name: 'TestResult',
  components: { Document, Delete, Download, Filter, Search, Refresh, List, Picture, Edit },
  setup() {
    const resultList = ref([])
    const currentPage = ref(1)
    const pageSize = ref(10)
    const total = ref(0)
    const selectedResults = ref([])
    const downloading = ref(false)
    const downloadingPDF = ref(false)
    const loading = ref(false)
    const filterForm = ref({ startDate: '', endDate: '', isGlandFace: null, level: '' })
    const editDialogVisible = ref(false)
    const editFormRef = ref(null)
    const submitting = ref(false)
    const editForm = ref({
      id: null,
      isGlandFace: false,
      level: '',
      visualizationDescription: ''
    })
    const editRules = {
      isGlandFace: [{ required: true, message: '请选择是否为腺体面容', trigger: 'change' }],
      level: [{ required: true, message: '请选择严重程度', trigger: 'change' }]
    }

    const getResultList = async () => {
      loading.value = true
      try {
        const params = { page: currentPage.value, pageSize: pageSize.value }
        if (filterForm.value.startDate) params.startDate = filterForm.value.startDate
        if (filterForm.value.endDate) params.endDate = filterForm.value.endDate
        if (filterForm.value.isGlandFace !== null) params.isGlandFace = filterForm.value.isGlandFace
        if (filterForm.value.level) params.level = filterForm.value.level
        const res = await request.get('/testResult/result', { params })
        if (res.code === 1) { resultList.value = res.data.records || []; total.value = res.data.total || 0 }
        else ElMessage.error(res.msg || '获取检测记录失败')
      } catch (e) { console.error(e); ElMessage.error('获取检测记录失败') }
      finally { loading.value = false }
    }

    const handleSelectionChange = (sel) => { selectedResults.value = sel }

    const deleteSelectedResults = () => {
      if (!selectedResults.value.length) { ElMessage.warning('请先选择要删除的记录'); return }
      ElMessageBox.confirm(`确定要删除选中的 ${selectedResults.value.length} 条检测记录吗？`, '批量删除确认', { confirmButtonText: '确定删除', cancelButtonText: '取消', type: 'warning' })
        .then(async () => {
          try {
            const ids = selectedResults.value.map(i => i.id)
            const res = await request.delete('/testResult/batch', { data: ids })
            if (res?.code === 1) { ElMessage.success('批量删除成功'); selectedResults.value = []; await getResultList() }
            else ElMessage.error(res.msg || '批量删除失败')
          } catch (e) { ElMessage.error('批量删除失败') }
        }).catch(() => {})
    }

    const deleteResult = (id) => {
      ElMessageBox.confirm('确定要删除这条检测记录吗？', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
        .then(async () => {
          try {
            const res = await request.delete(`/testResult/${id}`)
            if (res?.code === 1) { ElMessage.success('删除成功'); await getResultList() }
            else ElMessage.error(res.msg || '删除失败')
          } catch (e) { ElMessage.error('删除失败') }
        }).catch(() => {})
    }

    const downloadResults = async () => {
      downloading.value = true
      try {
        const params = {}
        if (filterForm.value.startDate) params.startDate = filterForm.value.startDate
        if (filterForm.value.endDate) params.endDate = filterForm.value.endDate
        if (filterForm.value.isGlandFace !== null) params.isGlandFace = filterForm.value.isGlandFace
        if (filterForm.value.level) params.level = filterForm.value.level
        const res = await downloadRequest.get('/testResult/download', { params })
        const blob = new Blob([res.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a'); link.href = url; link.download = `检测记录_${Date.now()}.xlsx`; link.click()
        window.URL.revokeObjectURL(url); ElMessage.success('下载成功')
      } catch (e) { ElMessage.error('下载失败') }
      finally { downloading.value = false }
    }

    const downloadResultsAsPDF = async () => {
      downloadingPDF.value = true
      try {
        const params = {}
        if (filterForm.value.startDate) params.startDate = filterForm.value.startDate
        if (filterForm.value.endDate) params.endDate = filterForm.value.endDate
        const res = await downloadRequest.get('/testResult/download/pdf', { params })
        const blob = new Blob([res.data], { type: 'text/html' })
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a'); link.href = url; link.download = `检测记录_${Date.now()}.pdf.html`; link.click()
        window.URL.revokeObjectURL(url); ElMessage.success('PDF下载成功')
      } catch (e) { ElMessage.error('PDF下载失败') }
      finally { downloadingPDF.value = false }
    }

    const searchResults = () => { currentPage.value = 1; getResultList() }
    const resetFilter = () => { filterForm.value = { startDate: '', endDate: '', isGlandFace: null, level: '' }; currentPage.value = 1; getResultList() }
    const handleSizeChange = (val) => { pageSize.value = val; currentPage.value = 1; getResultList() }
    const handleCurrentChange = (val) => { currentPage.value = val; getResultList() }

    const getLevelType = (level) => {
      if (level?.includes('轻')) return 'success'
      if (level?.includes('中')) return 'warning'
      if (level?.includes('重')) return 'danger'
      return 'info'
    }

    const getImageUrl = (path) => {
      if (!path) return ''
      if (path.startsWith('http')) return path
      return 'https://java-web-ai388.oss-cn-beijing.aliyuncs.com/' + path
    }

    const openEditDialog = (row) => {
      editForm.value = {
        id: row.id,
        isGlandFace: row.isGlandFace,
        level: row.level,
        visualizationDescription: row.visualizationDescription || ''
      }
      editDialogVisible.value = true
    }

    const submitEdit = async () => {
      if (!editFormRef.value) return
      await editFormRef.value.validate(async (valid) => {
        if (valid) {
          submitting.value = true
          try {
            const res = await request.put('/testResult/update', editForm.value)
            if (res.code === 1) {
              ElMessage.success('更新成功')
              editDialogVisible.value = false
              await getResultList()
            } else {
              ElMessage.error(res.msg || '更新失败')
            }
          } catch (e) {
            ElMessage.error('更新失败')
          } finally {
            submitting.value = false
          }
        }
      })
    }

    onMounted(() => { getResultList() })

    return { resultList, loading, downloading, downloadingPDF, currentPage, pageSize, total, filterForm, selectedResults, editDialogVisible, editFormRef, submitting, editForm, editRules, getResultList, handleSelectionChange, deleteSelectedResults, deleteResult, downloadResults, downloadResultsAsPDF, searchResults, resetFilter, handleSizeChange, handleCurrentChange, getLevelType, getImageUrl, openEditDialog, submitEdit }
  }
}
</script>

<style scoped>
.test-result-page {
  padding: 24px;
  min-height: 100%;
  background: var(--med-bg);
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 18px;
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
  flex-wrap: wrap;
  gap: 16px;
}

.page-header-title {
  display: flex; align-items: center; gap: 10px;
  font-size: 22px; font-weight: 700; color: #fff; margin: 0 0 6px 0;
}

.page-header-sub { font-size: 13px; color: rgba(255,255,255,0.55); margin: 0; }

.header-actions { display: flex; gap: 10px; flex-wrap: wrap; align-items: center; }

/* 卡片公共 */
.filter-card :deep(.el-card__header),
.table-card :deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid var(--med-border-light);
}

.section-card-header {
  display: flex; align-items: center; gap: 8px;
  font-size: 15px; font-weight: 700; color: var(--med-text);
}

.table-total {
  margin-left: auto;
  font-size: 13px;
  color: var(--med-text-secondary);
  font-weight: 400;
}

.table-total strong { color: var(--med-primary); font-size: 15px; }

/* 筛选表单 */
.filter-form { display: flex; flex-wrap: wrap; gap: 12px 16px; align-items: center; padding: 4px 0; }

/* 表格 */
:deep(.med-table) { border-radius: var(--med-radius-md) !important; overflow: hidden; }
:deep(.med-table .el-table__cell) { padding: 10px 0; }
:deep(.med-table .el-table__row:hover > td) { background: var(--med-primary-bg) !important; }

.seq-num { font-weight: 700; color: var(--med-primary); font-size: 14px; }

.record-img {
  width: 80px; height: 80px;
  border-radius: var(--med-radius-md);
  border: 2px solid var(--med-border-light);
  object-fit: cover;
  transition: transform 0.25s ease;
}
.record-img:hover { transform: scale(1.08); }

.img-placeholder {
  width: 80px; height: 80px;
  background: var(--med-bg); border-radius: var(--med-radius-md);
  display: flex; align-items: center; justify-content: center;
  color: var(--med-text-muted); font-size: 24px;
}

.time-cell { font-size: 13px; color: var(--med-text-secondary); }
.na-text { font-size: 13px; color: var(--med-text-muted); }

/* 分页 */
.pagination-wrap {
  display: flex; justify-content: center;
  padding: 20px;
  background: #fff;
  border-radius: var(--med-radius-lg);
  box-shadow: var(--med-shadow-sm);
  border: 1px solid var(--med-border-light);
}

:deep(.el-pagination.is-background .el-pager li.is-active) {
  background: var(--med-primary) !important;
}

@media (max-width: 768px) {
  .test-result-page { padding: 16px; gap: 14px; }
  .page-header { flex-direction: column; align-items: flex-start; }
  .header-actions { width: 100%; }
  .filter-form { flex-direction: column; align-items: stretch; }
}
</style>
