<template>
  <div class="face-analysis-container">
    <div class="header">
      <h2 class="section-header"><el-icon><Camera /></el-icon> AI腺样体面容分析</h2>
      <p class="header-subtitle">上传多张面部照片，系统会自动剔除无脸照片、综合多图特征生成更精准的诊断</p>
      <div class="header-features">
        <div class="feature-item">
          <el-icon class="feature-icon"><MagicStick /></el-icon>
          <span>多图综合</span>
        </div>
        <div class="feature-item">
          <el-icon class="feature-icon"><TrendCharts /></el-icon>
          <span>几何测量为主</span>
        </div>
        <div class="feature-item">
          <el-icon class="feature-icon"><Document /></el-icon>
          <span>专业报告</span>
        </div>
      </div>
    </div>

    <el-card class="upload-card">
      <template #header>
        <div class="card-header">
          <h3><el-icon><Upload /></el-icon> 图片上传</h3>
          <p class="card-subtitle">支持批量上传或拍照检测（最多9张）</p>
        </div>
      </template>

      <div class="upload-section">
        <div v-if="isCameraActive" class="camera-preview">
          <div class="camera-frame">
            <video ref="videoElement" autoplay playsinline class="video-preview"></video>
            <canvas ref="canvasElement" style="display: none;"></canvas>
            <div class="camera-overlay">
              <div class="camera-guideline"></div>
              <div class="camera-instructions">
                <el-icon class="capture-icon"><Camera /></el-icon>
                <p>请将面部对准摄像头</p>
                <p class="instruction-detail">确保光线充足，面部清晰可见</p>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="dropzone" @click="triggerFileInput" @dragover.prevent @dragenter.prevent @drop.prevent="onDrop">
          <input ref="fileInput" type="file" accept="image/*" multiple hidden @change="onFileInputChange" />
          <div v-if="!selectedFiles.length" class="upload-placeholder">
            <div class="upload-icon-container">
              <el-icon class="el-icon--upload upload-icon"><UploadFilled /></el-icon>
            </div>
            <div class="el-upload__text">将图片拖到此处，或 <em>点击上传</em></div>
            <div class="el-upload__tip">
              <p>支持 JPG、PNG 格式，可多选（最多9张）</p>
              <p class="file-limit">建议上传正面免冠、光线充足的照片</p>
            </div>
          </div>
          <div v-else class="multi-preview">
            <div
              v-for="(file, index) in selectedFiles"
              :key="file.id"
              class="preview-item"
            >
              <div class="preview-wrapper">
                <img :src="file.previewUrl" :alt="file.name" class="preview-image-small" />
                <div class="preview-overlay">
                  <el-button
                    type="danger"
                    size="small"
                    @click.stop="removeFile(index)"
                    circle
                    class="remove-btn"
                  >
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
              </div>
              <p class="file-name">{{ file.name }}</p>
            </div>
            <div class="preview-item add-tile" @click.stop="triggerFileInput">
              <el-icon class="add-icon"><Plus /></el-icon>
              <span>继续添加</span>
            </div>
          </div>
        </div>

        <div class="action-buttons">
          <el-button
            v-if="!isCameraActive"
            type="primary"
            @click="startCamera"
            size="large"
            class="action-button"
          >
            <el-icon><Camera /></el-icon> 拍照检测
          </el-button>

          <template v-else>
            <el-button type="primary" @click="captureImage" size="large" class="action-button capture-button">
              <el-icon><CameraFilled /></el-icon> 拍摄
            </el-button>
            <el-button @click="stopCamera" size="large" class="action-button secondary">
              <el-icon><Close /></el-icon> 取消
            </el-button>
          </template>

          <el-button
            v-if="selectedFiles.length && !isCameraActive"
            type="primary"
            :disabled="loading"
            :loading="loading"
            @click="analyzeFaces"
            size="large"
            class="action-button analyze-button"
          >
            <el-icon><Search /></el-icon> {{ loading ? '分析中...' : '开始分析' }}
          </el-button>

          <el-button
            v-if="selectedFiles.length && !isCameraActive"
            @click="clearFiles"
            size="large"
            class="action-button secondary"
          >
            <el-icon><Refresh /></el-icon> 重新选择
          </el-button>
        </div>
      </div>
    </el-card>

    <div v-if="showResults && analysisResult" class="analysis-results-section" ref="resultsSection">
      <div class="single-result-container">
        <div class="result-header-section">
          <div class="result-icon-wrapper" :class="{ 'warning-bg': analysisResult.isGlandFace, 'success-bg': !analysisResult.isGlandFace }">
            <svg class="result-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" :fill="analysisResult.isGlandFace ? '#ef4444' : '#22c55e'"/>
            </svg>
          </div>
          <div class="result-title-section">
            <h2 class="result-main-title">AI综合分析报告</h2>
            <p class="result-subtitle">多图综合 · 专业分析</p>
            <div class="result-tags">
              <div class="tag-badge" :class="{ 'danger': analysisResult.isGlandFace, 'success': !analysisResult.isGlandFace }">
                <svg class="tag-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z" :fill="analysisResult.isGlandFace ? '#ef4444' : '#22c55e'"/>
                </svg>
                {{ analysisResult.isGlandFace ? '腺体面容' : '非腺体面容' }}
              </div>
              <div class="tag-badge info">
                <svg class="tag-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M13 9L18 14L16.59 15.41L13 11.83V19H11V11.83L7.41 15.41L6 14L11 9V3H13V9Z" fill="#0891B2"/>
                </svg>
                {{ analysisResult.level }}
              </div>
            </div>
          </div>
          <div class="result-status-indicator" :class="{ 'warning': analysisResult.isGlandFace, 'success': !analysisResult.isGlandFace }">
            <div class="pulse-ring"></div>
            <div class="status-dot"></div>
          </div>
        </div>

        <div class="result-summary-meta">
          本次综合 {{ analysisResult.imagePaths ? analysisResult.imagePaths.length : 1 }} 张照片
          <template v-if="analysisResult.skippedImages && analysisResult.skippedImages.length">
            ，已自动剔除 {{ analysisResult.skippedImages.length }} 张未检测到人脸的照片
          </template>
        </div>

        <div class="result-content-grid">
          <div class="result-image-card">
            <div class="image-card-header">
              <el-icon><Picture /></el-icon>
              <span>分析图像（{{ analysisResult.imagePaths ? analysisResult.imagePaths.length : 1 }}张）</span>
            </div>
            <div class="image-wrapper gallery">
              <div
                v-for="(img, idx) in displayImages"
                :key="idx"
                class="gallery-item"
              >
                <el-image :src="img" :preview-src-list="displayImages" :initial-index="idx" fit="cover" class="gallery-image" />
              </div>
            </div>
            <div v-if="analysisResult.skippedImages && analysisResult.skippedImages.length" class="skipped-warning">
              <el-alert
                type="warning"
                :closable="false"
                show-icon
                :title="`以下 ${analysisResult.skippedImages.length} 张照片未检测到人脸，已自动排除出综合判断：`"
                :description="skippedNames"
              />
            </div>
          </div>

          <div class="result-details-card">
            <div class="details-card-header">
              <el-icon><DataAnalysis /></el-icon>
              <span>详细数据</span>
            </div>

            <div class="detail-item">
              <div class="detail-label"><el-icon><Timer /></el-icon><span>检测时间</span></div>
              <div class="detail-value">{{ formatDateTime(analysisResult.testTime) }}</div>
            </div>


            <div class="detail-item full-width">
              <div class="detail-label"><el-icon><Document /></el-icon><span>分析描述</span></div>
              <div class="detail-value description">{{ analysisResult.visualizationDescription }}</div>
            </div>

            <div v-if="analysisResult.referenceCases" class="detail-item full-width">
              <div class="detail-label"><el-icon><Collection /></el-icon><span>相似参考病例</span></div>
              <div class="detail-value reference">{{ analysisResult.referenceCases }}</div>
            </div>
          </div>
        </div>

        <div class="action-buttons-result">
          <el-button type="primary" @click="exportToPDF" class="export-btn">
            <el-icon><Download /></el-icon> 导出PDF报告
          </el-button>
          <el-button @click="clearResults" class="clear-btn">
            <el-icon><Refresh /></el-icon> 重新分析
          </el-button>
        </div>
      </div>
    </div>

    <el-alert
      v-if="error"
      :title="error"
      type="error"
      show-icon
      closable
      @close="error = ''"
      class="error-alert"
    />
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue';
import request, { downloadRequest, faceAnalysisRequest } from '../utils/request';
import { ElMessage } from 'element-plus';
import {
  Camera,
  Upload,
  UploadFilled,
  Delete,
  Plus,
  CameraFilled,
  Close,
  Search,
  Refresh,
  Picture,
  Timer,
  DataAnalysis,
  TrendCharts,
  Document,
  Collection,
  Download
} from '@element-plus/icons-vue';

export default {
  name: 'FaceAnalysis',
  components: {
    Camera,
    Upload,
    UploadFilled,
    Delete,
    Plus,
    CameraFilled,
    Close,
    Search,
    Refresh,
    Picture,
    Timer,
    DataAnalysis,
    TrendCharts,
    Document,
    Collection,
    Download
  },
  setup() {
    const selectedFiles = ref([]);
    const loading = ref(false);
    const analysisResult = ref(null);
    const error = ref('');
    const isCameraActive = ref(false);
    const videoElement = ref(null);
    const canvasElement = ref(null);
    const fileInput = ref(null);
    const resultsSection = ref(null);
    let stream = null;
    let fileSeq = 0;

    const scrollToResults = () => {
      nextTick(() => {
        if (resultsSection.value) {
          resultsSection.value.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
      });
    };

    // ---- 自管理文件列表（彻底修复 el-upload 内部列表导致的“删除又复活”） ----
    const dedupKey = (f) => `${f.name}_${f.size}_${f.lastModified}`;

    const addFiles = (fileList) => {
      const existing = new Set(selectedFiles.value.map(f => f.key));
      for (const f of fileList) {
        const key = dedupKey(f);
        if (existing.has(key)) continue; // 去重，避免重复添加同一文件
        if (selectedFiles.value.length >= 9) {
          ElMessage.warning('最多只能上传 9 张照片');
          break;
        }
        existing.add(key);
        selectedFiles.value.push({
          id: ++fileSeq,
          key,
          raw: f,
          name: f.name,
          previewUrl: URL.createObjectURL(f)
        });
      }
      error.value = '';
    };

    const triggerFileInput = () => {
      if (fileInput.value) fileInput.value.click();
    };

    const onFileInputChange = (e) => {
      if (e.target.files && e.target.files.length) {
        addFiles(Array.from(e.target.files));
      }
      // 清空 input，保证同一文件可再次选择
      e.target.value = '';
    };

    const onDrop = (e) => {
      if (e.dataTransfer && e.dataTransfer.files && e.dataTransfer.files.length) {
        addFiles(Array.from(e.dataTransfer.files));
      }
    };

    const removeFile = (index) => {
      const file = selectedFiles.value[index];
      if (file && file.previewUrl) URL.revokeObjectURL(file.previewUrl);
      selectedFiles.value.splice(index, 1);
    };

    const clearFiles = () => {
      selectedFiles.value.forEach(f => { if (f.previewUrl) URL.revokeObjectURL(f.previewUrl); });
      selectedFiles.value = [];
      clearResults();
    };

    const clearResults = () => {
      analysisResult.value = null;
      showResults.value = false;
    };

    // ---- 摄像头 ----
    const startCamera = async () => {
      try {
        clearFiles();
        isCameraActive.value = true;
        error.value = '';
        stream = await navigator.mediaDevices.getUserMedia({
          video: { facingMode: 'user', width: { ideal: 1280 }, height: { ideal: 720 } },
          audio: false
        });
        if (videoElement.value) videoElement.value.srcObject = stream;
      } catch (err) {
        console.error('无法访问摄像头:', err);
        error.value = '无法访问摄像头，请检查权限设置';
        isCameraActive.value = false;
      }
    };

    const captureImage = () => {
      if (!videoElement.value || !canvasElement.value) return;
      const video = videoElement.value;
      const canvas = canvasElement.value;
      const context = canvas.getContext('2d');
      canvas.width = video.videoWidth;
      canvas.height = video.videoHeight;
      context.drawImage(video, 0, 0, canvas.width, canvas.height);
      canvas.toBlob((blob) => {
        if (blob) {
          const file = new File([blob], `camera-${Date.now()}.jpg`, { type: 'image/jpeg' });
          addFiles([file]);
          stopCamera();
        }
      }, 'image/jpeg', 0.95);
    };

    const stopCamera = () => {
      if (stream) {
        stream.getTracks().forEach(track => track.stop());
        stream = null;
      }
      isCameraActive.value = false;
    };

    // ---- 分析 ----
    const showResults = ref(false);

    const analyzeFaces = async () => {
      if (!selectedFiles.value.length) {
        error.value = '请先选择图片';
        return;
      }
      loading.value = true;
      error.value = '';
      try {
        const formData = new FormData();
        selectedFiles.value.forEach(f => formData.append('images', f.raw));

        const response = await faceAnalysisRequest.post('/doubao/analyzeMulti', formData);
        if (response && response.code === 1) {
          analysisResult.value = response.data;
          showResults.value = true;
          scrollToResults();
        } else {
          error.value = (response && response.msg) ? response.msg : '分析失败';
        }
      } catch (err) {
        console.error('分析请求失败:', err);
        if (err.response) {
          if (err.response.status === 401) error.value = '登录已过期，请重新登录';
          else if (err.response.data) error.value = err.response.data.msg || `分析请求失败，状态码: ${err.response.status}`;
          else error.value = `分析请求失败，状态码: ${err.response.status}`;
        } else if (err.request) {
          error.value = '网络连接异常，请检查网络设置';
        } else {
          error.value = '分析请求失败: ' + (err.message || '未知错误');
        }
      } finally {
        loading.value = false;
      }
    };

    // ---- 展示辅助 ----
    const displayImages = ref([]);
    const skippedNames = ref('');

    // 当分析结果变化时，同步展示用的图片列表与被剔除文件提示
    const syncDisplay = () => {
      const r = analysisResult.value;
      if (!r) { displayImages.value = []; skippedNames.value = ''; return; }
      displayImages.value = (r.imagePaths && r.imagePaths.length) ? r.imagePaths : (r.imagePath ? [r.imagePath] : []);
      if (r.skippedImages && r.skippedImages.length) {
        skippedNames.value = r.skippedImages.map(p => {
          const idx = Math.max(p.lastIndexOf('/'), p.lastIndexOf('\\'));
          return idx >= 0 ? p.slice(idx + 1) : p;
        }).join('、');
      } else {
        skippedNames.value = '';
      }
    };

    // 用 watch 监听 analysisResult 变化
    watch(analysisResult, syncDisplay);

    const formatDateTime = (dateString) => {
      if (!dateString) return '';
      const date = new Date(dateString);
      if (isNaN(date.getTime())) return '';
      return date.toLocaleString('zh-CN');
    };

    const getLevelType = (level) => {
      if (level && (level.includes('轻') || level.includes('微'))) return 'success';
      if (level && (level.includes('中') || level.includes('等'))) return 'warning';
      if (level && level.includes('重')) return 'danger';
      return 'info';
    };

    const exportToPDF = () => {
      const r = analysisResult.value;
      if (!r) { ElMessage.error('没有可导出的分析结果'); return; }

      const imgs = (r.imagePaths && r.imagePaths.length) ? r.imagePaths : (r.imagePath ? [r.imagePath] : []);
      const imgHtml = imgs.map(src =>
        `<div style="flex:1; min-width:240px; text-align:center;"><img src="${src}" style="max-width:100%; max-height:280px; border-radius:8px; box-shadow:0 2px 12px rgba(0,0,0,0.1);" /></div>`
      ).join('');

      const printContent = document.createElement('div');
      printContent.style.padding = '20px';
      printContent.style.fontFamily = 'Arial, sans-serif';
      printContent.innerHTML = `
        <div style="text-align:center; margin-bottom:20px;">
          <h2>腺体面容综合检测报告</h2>
        </div>
        <div style="display:flex; gap:15px; flex-wrap:wrap; justify-content:center; margin-bottom:20px;">
          <span style="padding:10px 20px; border-radius:4px; font-weight:bold; background-color:${r.isGlandFace ? '#fef0f0' : '#f0f9ff'}; color:${r.isGlandFace ? '#f56c6c' : '#409eff'}">
            ${r.isGlandFace ? '腺体面容' : '非腺体面容'}
          </span>
          <span style="padding:10px 20px; border-radius:4px; font-weight:bold; background-color:#ecf5ff; color:#409eff">
            ${r.level}
          </span>
        </div>
        <div style="margin-bottom:20px;">
          <table style="width:100%; border-collapse:collapse;">
            <tr><td style="border:1px solid #ddd; padding:12px; font-weight:bold;">检测时间</td><td style="border:1px solid #ddd; padding:12px;">${formatDateTime(r.testTime)}</td></tr>
            <tr><td style="border:1px solid #ddd; padding:12px; font-weight:bold;">分析描述</td><td style="border:1px solid #ddd; padding:12px;">${r.visualizationDescription}</td></tr>
          </table>
        </div>
        <div style="text-align:center;">
          <h4 style="margin-bottom:15px;">分析图像（${imgs.length}张）</h4>
          <div style="display:flex; gap:20px; flex-wrap:wrap; justify-content:center;">${imgHtml}</div>
        </div>
        <div style="margin-top:30px; text-align:center; color:#999; font-size:12px;">报告生成时间: ${new Date().toLocaleString('zh-CN')}</div>
      `;

      const printWindow = window.open('', '_blank');
      printWindow.document.write(`
        <html><head><title>腺体面容综合检测报告</title>
        <style>@media print { body { -webkit-print-color-adjust: exact; print-color-adjust: exact; } }</style>
        </head><body>${printContent.innerHTML}
        <script>window.onload=function(){window.print();window.close();}<\/script>
        </body></html>
      `);
      printWindow.document.close();
    };

    onMounted(() => {
      syncDisplay();
    });

    onUnmounted(() => {
      stopCamera();
      selectedFiles.value.forEach(f => { if (f.previewUrl) URL.revokeObjectURL(f.previewUrl); });
    });

    return {
      selectedFiles,
      loading,
      analysisResult,
      showResults,
      error,
      isCameraActive,
      videoElement,
      canvasElement,
      fileInput,
      resultsSection,
      triggerFileInput,
      onFileInputChange,
      onDrop,
      removeFile,
      clearFiles,
      clearResults,
      startCamera,
      captureImage,
      stopCamera,
      analyzeFaces,
      displayImages,
      skippedNames,
      formatDateTime,
      getLevelType,
      exportToPDF
    };
  }
};
</script>

<style scoped>
.face-analysis-container {
  padding: 24px;
  background: var(--med-bg);
  min-height: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.header {
  padding: 22px 28px;
  background: linear-gradient(135deg, var(--med-sidebar) 0%, var(--med-sidebar-mid) 60%, #0C4A56 100%);
  border-radius: var(--med-radius-lg);
  box-shadow: var(--med-shadow-md);
  color: white;
}
.header h2 { margin: 0 0 8px 0; font-size: 24px; font-weight: 700; display: flex; align-items: center; gap: 10px; color: #fff; }
.header-subtitle { margin: 0; font-size: 14px; color: rgba(255,255,255,0.65); }
.header-features { display: flex; gap: 20px; margin-top: 16px; flex-wrap: wrap; }
.header-features .feature-item { display: flex; align-items: center; gap: 6px; background: rgba(255,255,255,0.1); padding: 6px 14px; border-radius: 20px; font-size: 13px; color: rgba(255,255,255,0.85); border: 1px solid rgba(255,255,255,0.15); }
.feature-icon { font-size: 16px !important; color: var(--med-accent); }

.upload-card { border-radius: var(--med-radius-lg) !important; box-shadow: var(--med-shadow-md) !important; border: 1px solid var(--med-border-light) !important; overflow: hidden; }
.upload-card :deep(.el-card__header) { padding: 16px 22px; background: var(--med-bg-section); border-bottom: 1px solid var(--med-border-light); }
.card-header { display: flex; align-items: center; justify-content: space-between; }
.card-header h3 { margin: 0; color: var(--med-text); font-size: 16px; font-weight: 700; display: flex; align-items: center; gap: 8px; }
.card-subtitle { font-size: 12px; color: var(--med-text-muted); margin: 0; }
.upload-section { padding: 24px; }

.dropzone { border: 2px dashed var(--med-border) !important; border-radius: var(--med-radius-md); background: var(--med-bg); transition: all 0.25s ease; padding: 40px 0; cursor: pointer; }
.dropzone:hover { border-color: var(--med-primary) !important; background: var(--med-primary-bg) !important; }
.upload-placeholder { padding: 20px 0; text-align: center; }
.upload-icon-container { margin-bottom: 16px; }
.upload-icon { font-size: 52px !important; color: var(--med-primary) !important; }
.el-upload__text { font-size: 15px; color: var(--med-text-secondary); margin-bottom: 8px; }
.el-upload__text em { color: var(--med-primary); font-style: normal; font-weight: 600; }
.el-upload__tip { color: var(--med-text-muted); font-size: 13px; margin-top: 10px; line-height: 1.6; }
.file-limit { font-size: 12px; color: var(--med-text-muted); }

.multi-preview { display: flex; flex-wrap: wrap; gap: 16px; justify-content: center; }
.preview-item { position: relative; width: 110px; text-align: center; }
.preview-image-small { width: 100%; height: 110px; object-fit: cover; border-radius: var(--med-radius-md); border: 2px solid var(--med-border-light); box-shadow: var(--med-shadow-sm); }
.file-name { margin: 10px 0 5px 0; font-size: 12px; color: #666; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.preview-overlay { position: absolute; top: -8px; right: -8px; }
.remove-btn { background: var(--med-danger) !important; border: none !important; box-shadow: 0 2px 6px rgba(231, 76, 60, 0.4) !important; }
.add-tile { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 6px; height: 110px; border: 2px dashed var(--med-border); border-radius: var(--med-radius-md); color: var(--med-text-muted); cursor: pointer; transition: all 0.2s ease; }
.add-tile:hover { border-color: var(--med-primary); color: var(--med-primary); }
.add-icon { font-size: 28px !important; }

.action-buttons { display: flex; justify-content: center; gap: 16px; flex-wrap: wrap; padding-top: 20px; }
.action-button { min-width: 130px; padding: 12px 22px !important; font-size: 14px !important; border-radius: var(--med-radius-md) !important; display: flex; align-items: center; justify-content: center; gap: 8px; font-weight: 600 !important; transition: all 0.25s ease !important; }
.action-button:hover { transform: translateY(-2px) !important; }
.analyze-button { background: linear-gradient(135deg, var(--med-primary) 0%, var(--med-primary-light) 100%) !important; border: none !important; box-shadow: 0 4px 14px rgba(13,115,119,0.3) !important; }
.capture-button { background: linear-gradient(135deg, var(--med-ai-purple) 0%, #9C6FFF 100%) !important; border: none !important; box-shadow: 0 4px 14px rgba(124,77,255,0.3) !important; }
.action-button.secondary { background: var(--med-bg) !important; border: 1px solid var(--med-border) !important; color: var(--med-text-secondary) !important; }
.action-button.secondary:hover { background: #fff !important; border-color: var(--med-primary) !important; color: var(--med-primary) !important; }

.camera-preview { position: relative; width: 100%; height: 380px; background: #0A0F1A; border-radius: var(--med-radius-md); overflow: hidden; margin-bottom: 20px; box-shadow: 0 8px 24px rgba(0,0,0,0.25); border: 1px solid rgba(50,224,196,0.2); }
.video-preview { width: 100%; height: 100%; object-fit: cover; }
.camera-overlay { position: absolute; top: 0; left: 0; right: 0; bottom: 0; display: flex; align-items: center; justify-content: center; background: rgba(0,0,0,0.25); }
.camera-instructions { text-align: center; color: white; }
.camera-instructions .capture-icon { font-size: 44px; margin-bottom: 12px; color: var(--med-accent); }
.camera-instructions p { font-size: 16px; margin: 0 0 4px 0; }
.instruction-detail { font-size: 12px; color: rgba(255,255,255,0.65); }

.analysis-results-section { margin-top: 20px; }
.single-result-container { animation: slideUp 0.5s ease-out; }
@keyframes slideUp { from { opacity: 0; transform: translateY(30px); } to { opacity: 1; transform: translateY(0); } }

.result-header-section { display: flex; align-items: center; gap: 24px; margin-bottom: 16px; padding: 28px 32px; background: linear-gradient(135deg, #ffffff 0%, #F0FDFA 100%); border-radius: 20px; box-shadow: 0 8px 32px rgba(8,145,178,0.08); border: 1px solid rgba(34,211,238,0.2); position: relative; overflow: hidden; }
.result-icon-wrapper { width: 80px; height: 80px; background: linear-gradient(135deg, rgba(34,211,238,0.15) 0%, rgba(34,197,94,0.1) 100%); border-radius: 24px; display: flex; align-items: center; justify-content: center; border: 2px solid rgba(34,211,238,0.3); flex-shrink: 0; }
.result-icon-wrapper.success-bg { background: linear-gradient(135deg, rgba(34,197,94,0.15) 0%, rgba(34,197,94,0.08) 100%); border-color: rgba(34,197,94,0.35); }
.result-icon-wrapper.warning-bg { background: linear-gradient(135deg, rgba(239,68,68,0.15) 0%, rgba(239,68,68,0.08) 100%); border-color: rgba(239,68,68,0.35); }
.result-icon { width: 44px; height: 44px; }
.result-title-section { flex: 1; }
.result-main-title { margin: 0 0 6px 0; font-size: 26px; font-weight: 800; color: #134E4A; background: linear-gradient(135deg, #134E4A 0%, #0891B2 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }
.result-subtitle { margin: 0 0 14px 0; font-size: 14px; color: #64748b; font-weight: 500; }
.result-tags { display: flex; gap: 12px; flex-wrap: wrap; }
.tag-badge { display: flex; align-items: center; gap: 8px; padding: 10px 20px; border-radius: 28px; font-size: 14px; font-weight: 600; }
.tag-badge.success { background: linear-gradient(135deg, rgba(34,197,94,0.15) 0%, rgba(34,197,94,0.08) 100%); color: #166534; border: 1px solid rgba(34,197,94,0.3); }
.tag-badge.danger { background: linear-gradient(135deg, rgba(239,68,68,0.15) 0%, rgba(239,68,68,0.08) 100%); color: #991b1b; border: 1px solid rgba(239,68,68,0.3); }
.tag-badge.info { background: linear-gradient(135deg, rgba(8,145,178,0.15) 0%, rgba(8,145,178,0.08) 100%); color: #155e75; border: 1px solid rgba(8,145,178,0.3); }
.tag-icon { width: 18px; height: 18px; }
.result-status-indicator { position: relative; width: 60px; height: 60px; display: flex; align-items: center; justify-content: center; }
.pulse-ring { position: absolute; width: 100%; height: 100%; border-radius: 50%; animation: pulse-ring 2s cubic-bezier(0.455,0.03,0.515,0.955) infinite; }
.result-status-indicator.success .pulse-ring { border: 3px solid rgba(34,197,94,0.3); }
.result-status-indicator.warning .pulse-ring { border: 3px solid rgba(239,68,68,0.3); }
@keyframes pulse-ring { 0% { transform: scale(0.8); opacity: 1; } 100% { transform: scale(1.5); opacity: 0; } }
.status-dot { width: 24px; height: 24px; border-radius: 50%; position: relative; z-index: 1; }
.result-status-indicator.success .status-dot { background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%); box-shadow: 0 0 20px rgba(34,197,94,0.4); }
.result-status-indicator.warning .status-dot { background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); box-shadow: 0 0 20px rgba(239,68,68,0.4); }

.result-summary-meta { margin: 0 0 16px 4px; font-size: 13px; color: var(--med-text-muted); }

.result-content-grid { display: grid; grid-template-columns: 1fr 1.2fr; gap: 24px; }
.result-image-card, .result-details-card { background: #fff; border-radius: var(--med-radius-lg); box-shadow: var(--med-shadow-md); border: 1px solid var(--med-border-light); overflow: hidden; }
.image-card-header, .details-card-header { padding: 16px 20px; background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%); border-bottom: 1px solid var(--med-border-light); display: flex; align-items: center; gap: 8px; font-weight: 600; color: var(--med-text); font-size: 15px; }
.image-wrapper { padding: 24px; display: flex; justify-content: center; align-items: center; min-height: 300px; background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%); }
.gallery { display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap: 12px; width: 100%; align-items: start; }
.gallery-item { border-radius: 12px; overflow: hidden; box-shadow: var(--med-shadow-sm); }
.gallery-image { width: 100%; height: 160px; }
.skipped-warning { padding: 0 16px 16px; }

.detail-item { padding: 16px 20px; border-bottom: 1px solid var(--med-border-light); display: flex; flex-direction: column; gap: 8px; }
.detail-item.full-width { grid-column: 1 / -1; }
.detail-label { display: flex; align-items: center; gap: 6px; font-weight: 600; color: var(--med-text-secondary); font-size: 14px; }
.detail-value { color: var(--med-text); font-size: 15px; line-height: 1.6; }
.detail-value.description { color: var(--med-text-secondary); padding: 12px; background: var(--med-bg); border-radius: 8px; border-left: 3px solid var(--med-primary); }
.detail-value.reference { color: var(--med-text-secondary); padding: 12px; background: var(--med-bg); border-radius: 8px; white-space: pre-wrap; }


.action-buttons-result { padding: 20px; display: flex; gap: 12px; justify-content: flex-end; }
.export-btn, .clear-btn { display: flex; align-items: center; gap: 6px; padding: 10px 20px; border-radius: var(--med-radius-md); font-weight: 600; transition: all 0.25s ease; }
.export-btn:hover, .clear-btn:hover { transform: translateY(-2px); }

.error-alert { margin-top: 20px; border-radius: 12px; }

@media (max-width: 1024px) {
  .result-content-grid { grid-template-columns: 1fr; }
}
@media (max-width: 768px) {
  .face-analysis-container { padding: 15px; }
  .header { padding: 20px; border-radius: 16px; }
  .header h2 { font-size: 22px; }
  .upload-card { border-radius: 16px; }
  .upload-section { padding: 20px; }
  .action-buttons { gap: 15px; }
  .action-button { min-width: 120px; }
  .multi-preview { gap: 15px; }
  .preview-item { width: 100px; }
  .preview-image-small { height: 100px; }
  .result-header-section { flex-direction: column; text-align: center; padding: 20px; }
  .result-main-title { font-size: 20px; }
  .action-buttons-result { flex-direction: column; }
  .export-btn, .clear-btn { width: 100%; justify-content: center; }
}
</style>
