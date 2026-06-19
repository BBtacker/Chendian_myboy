<template>
  <div class="face-analysis-container">
    <div class="header">
      <h2 class="section-header"><el-icon><Camera /></el-icon> AI腺样体面容分析</h2>
      <p class="header-subtitle">上传面部照片或使用摄像头进行智能检测与分析</p>
      <div class="header-features">
        <div class="feature-item">
          <el-icon class="feature-icon"><MagicStick /></el-icon>
          <span>AI智能识别</span>
        </div>
        <div class="feature-item">
          <el-icon class="feature-icon"><TrendCharts /></el-icon>
          <span>精准分析</span>
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
          <p class="card-subtitle">支持批量上传或拍照检测</p>
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
        
        <el-upload
          v-else
          class="upload-area"
          drag
          :auto-upload="false"
          :show-file-list="false"
          :on-change="handleFileChange"
          :multiple="true"
          accept="image/*"
        >
          <div v-if="!selectedFiles.length" class="upload-placeholder">
            <div class="upload-icon-container">
              <el-icon class="el-icon--upload upload-icon">
                <UploadFilled />
              </el-icon>
            </div>
            <div class="el-upload__text">
              将图片拖到此处，或 <em>点击上传</em>
            </div>
            <div class="el-upload__tip">
              <p>支持 JPG、PNG 格式图片，可多选</p>
              <p class="file-limit">最多可上传9张图片</p>
            </div>
            <div class="upload-hints">
              <el-icon class="hint-icon"><InfoFilled /></el-icon>
              <span>建议上传正面免冠照片，光线充足</span>
            </div>
          </div>
          <div v-else class="preview-container">
            <div class="multi-preview">
              <div 
                v-for="(file, index) in selectedFiles" 
                :key="index" 
                class="preview-item"
              >
                <div class="preview-wrapper">
                  <img :src="file.previewUrl" :alt="'预览图片' + (index + 1)" class="preview-image-small" />
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
            </div>
          </div>
        </el-upload>
        
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
            <el-button 
              type="primary" 
              @click="captureImage"
              size="large"
              class="action-button capture-button"
            >
              <el-icon><CameraFilled /></el-icon> 拍摄
            </el-button>
            
            <el-button 
              @click="stopCamera"
              size="large"
              class="action-button secondary"
            >
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
    
    <div v-if="showResults" class="analysis-results-section" ref="resultsSection">
      <transition name="fade-scale">
        <div v-if="showSingleResult && analysisResult" class="single-result-container">
          <div class="result-header-section">
            <div class="result-icon-wrapper" :class="{ 'warning-bg': analysisResult.isGlandFace, 'success-bg': !analysisResult.isGlandFace }">
              <svg class="result-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" :fill="analysisResult.isGlandFace ? '#ef4444' : '#22c55e'"/>
              </svg>
            </div>
            <div class="result-title-section">
              <h2 class="result-main-title">AI分析报告</h2>
              <p class="result-subtitle">智能面容检测 · 专业分析</p>
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
          
          <div class="result-content-grid">
            <div class="result-image-card">
              <div class="image-card-header">
                <el-icon><Picture /></el-icon>
                <span>分析图像</span>
              </div>
              <div class="image-wrapper">
                <el-image 
                  :src="analysisResult.imagePath" 
                  :preview-src-list="[analysisResult.imagePath]"
                  fit="contain"
                  class="result-image-preview"
                />
              </div>
            </div>
            
            <div class="result-details-card">
              <div class="details-card-header">
                <el-icon><DataAnalysis /></el-icon>
                <span>详细数据</span>
              </div>
              
              <div class="detail-item">
                <div class="detail-label">
                  <el-icon><Timer /></el-icon>
                  <span>检测时间</span>
                </div>
                <div class="detail-value">{{ formatDateTime(analysisResult.testTime) }}</div>
              </div>
              
              <div class="detail-item">
                <div class="detail-label">
                  <el-icon><TrendCharts /></el-icon>
                  <span>置信度</span>
                </div>
                <div class="confidence-display">
                  <div class="confidence-bar-wrapper">
                    <div class="confidence-track"></div>
                    <div class="confidence-fill" :style="{ width: (analysisResult.confidence * 100) + '%' }"></div>
                  </div>
                  <span class="confidence-number">{{ (analysisResult.confidence * 100).toFixed(2) }}%</span>
                </div>
              </div>
              
              <div class="detail-item full-width">
                <div class="detail-label">
                  <el-icon><Document /></el-icon>
                  <span>分析描述</span>
                </div>
                <div class="detail-value description">{{ analysisResult.visualizationDescription }}</div>
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
        </div>
      </transition>
      
      <transition name="fade-scale">
        <div v-if="showCompareResults && compareResults.length" class="compare-results-container">
          <div class="compare-header-section">
            <div class="compare-icon-wrapper">
              <svg class="result-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M9 19V6L4 9M9 19L14 9M9 19H15M15 19V5L20 8M15 19L10 9M15 19H9" fill="#0891B2"/>
              </svg>
            </div>
            <div class="compare-title-wrapper">
              <h2 class="compare-main-title">AI对比分析报告</h2>
              <p class="compare-subtitle">多样本智能分析 · 数据对比</p>
            </div>
            <div class="compare-badge">
              <svg class="badge-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z" fill="#22c55e"/>
              </svg>
              {{ compareResults.length }} 样本
            </div>
          </div>
          
          <div class="compare-results-grid">
            <div 
              v-for="(result, index) in compareResults" 
              :key="index" 
              class="compare-result-card"
            >
              <div class="compare-card-header">
                <span class="sample-number">样本 {{ index + 1 }}</span>
                <el-tag :type="result.isGlandFace ? 'danger' : 'success'" size="small">
                  {{ result.isGlandFace ? '腺体面容' : '正常' }}
                </el-tag>
              </div>
              
              <div class="compare-card-content">
                <div class="compare-image-wrapper">
                  <el-image 
                    :src="result.imagePath" 
                    fit="cover"
                    :preview-src-list="[result.imagePath]"
                    preview-teleported
                    class="compare-image"
                  >
                    <template #error>
                      <div class="image-slot">
                        <el-icon><PictureRounded /></el-icon>
                      </div>
                    </template>
                  </el-image>
                </div>
                
                <div class="compare-details">
                  <div class="compare-detail-row">
                    <span class="compare-label"><el-icon><Timer /></el-icon> {{ formatDateTime(result.testTime) }}</span>
                  </div>
                  
                  <div class="compare-detail-row">
                    <span class="compare-label">等级:</span>
                    <el-tag :type="getLevelType(result.level)" size="small">{{ result.level }}</el-tag>
                  </div>
                  
                  <div class="compare-detail-row">
                    <span class="compare-label">置信度:</span>
                    <div class="mini-confidence-bar">
                      <div class="mini-confidence-fill" :style="{ width: (result.confidence * 100) + '%' }"></div>
                    </div>
                    <span class="mini-confidence-value">{{ (result.confidence * 100).toFixed(1) }}%</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <div v-if="compareResults.length > 1" class="compare-summary-card">
            <div class="summary-card-header">
              <el-icon><DataBoard /></el-icon>
              <span>对比分析总结</span>
            </div>
            
            <div class="summary-stats">
              <div class="stat-item">
                <div class="stat-icon-wrapper">
                  <el-icon><Collection /></el-icon>
                </div>
                <div class="stat-number">{{ compareResults.length }}</div>
                <div class="stat-label">总检测数</div>
              </div>
              
              <div class="stat-item">
                <div class="stat-icon-wrapper warning">
                  <el-icon><Warning /></el-icon>
                </div>
                <div class="stat-number danger">{{ compareResults.filter(r => r.isGlandFace).length }}</div>
                <div class="stat-label">腺体面容数</div>
              </div>
              
              <div class="stat-item">
                <div class="stat-icon-wrapper info">
                  <el-icon><PieChart /></el-icon>
                </div>
                <div class="stat-number primary">{{ (compareResults.reduce((sum, r) => sum + r.confidence, 0) / compareResults.length * 100).toFixed(2) }}%</div>
                <div class="stat-label">平均置信度</div>
              </div>
            </div>
          </div>
          
          <div class="compare-actions">
            <el-button type="primary" @click="exportCompareResults" class="export-btn">
              <el-icon><Download /></el-icon> 导出对比报告
            </el-button>
            <el-button @click="clearResults" class="clear-btn">
              <el-icon><Refresh /></el-icon> 重新分析
            </el-button>
          </div>
        </div>
      </transition>
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
import { ref, onMounted, onUnmounted, nextTick } from 'vue';
import request, { downloadRequest, faceAnalysisRequest } from '../utils/request';
import { ElMessage, ElLoading } from 'element-plus';
import {
  Camera,
  Upload,
  UploadFilled,
  Delete,
  CameraFilled,
  Close,
  Search,
  Refresh,
  Picture,
  PictureRounded,
  Timer,
  User,
  DataAnalysis,
  TrendCharts,
  Document,
  DataBoard,
  Collection,
  Warning,
  PieChart,
  Download
} from '@element-plus/icons-vue';

export default {
  name: 'FaceAnalysis',
  components: {
    Camera,
    Upload,
    UploadFilled,
    Delete,
    CameraFilled,
    Close,
    Search,
    Refresh,
    Picture,
    PictureRounded,
    Timer,
    User,
    DataAnalysis,
    TrendCharts,
    Document,
    DataBoard,
    Collection,
    Warning,
    PieChart,
    Download
  },
  setup() {
    const selectedFiles = ref([]);
    const loading = ref(false);
    const analysisResult = ref(null);
    const compareResults = ref([]);
    const error = ref('');
    const isCameraActive = ref(false);
    const videoElement = ref(null);
    const canvasElement = ref(null);
    const isMobile = ref(false);
    const showResults = ref(false);
    const showSingleResult = ref(false);
    const showCompareResults = ref(false);
    const resultsSection = ref(null);
    let stream = null;

    const scrollToResults = () => {
      nextTick(() => {
        if (resultsSection.value) {
          resultsSection.value.scrollIntoView({
            behavior: 'smooth',
            block: 'start'
          });
        }
      });
    };

    const checkIsMobile = () => {
      isMobile.value = window.innerWidth <= 768;
    };

    onMounted(() => {
      checkIsMobile();
      window.addEventListener('resize', checkIsMobile);
    });

    onUnmounted(() => {
      window.removeEventListener('resize', checkIsMobile);
      stopCamera();
    });

    const handleFileChange = (file, fileList) => {
      const newFiles = fileList.map(f => {
        return {
          raw: f.raw,
          name: f.name,
          previewUrl: URL.createObjectURL(f.raw)
        };
      });
      
      selectedFiles.value = newFiles;
      error.value = '';
    };

    const removeFile = (index) => {
      const file = selectedFiles.value[index];
      if (file && file.previewUrl) {
        URL.revokeObjectURL(file.previewUrl);
      }
      selectedFiles.value.splice(index, 1);
    };

    const clearFiles = () => {
      selectedFiles.value.forEach(file => {
        if (file.previewUrl) {
          URL.revokeObjectURL(file.previewUrl);
        }
      });
      selectedFiles.value = [];
      clearResults();
    };

    const clearResults = () => {
      analysisResult.value = null;
      compareResults.value = [];
      showResults.value = false;
      showSingleResult.value = false;
      showCompareResults.value = false;
    };

    const startCamera = async () => {
      try {
        clearFiles();
        isCameraActive.value = true;
        error.value = '';
        
        stream = await navigator.mediaDevices.getUserMedia({ 
          video: { 
            facingMode: 'user',
            width: { ideal: 1280 },
            height: { ideal: 720 }
          }, 
          audio: false 
        });
        
        if (videoElement.value) {
          videoElement.value.srcObject = stream;
        }
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
          const file = new File([blob], 'camera-capture.jpg', { type: 'image/jpeg' });
          selectedFiles.value = [{
            raw: file,
            name: 'camera-capture.jpg',
            previewUrl: URL.createObjectURL(file)
          }];
          stopCamera();
        }
      }, 'image/jpeg', 0.95);
    };

    const stopCamera = () => {
      if (stream) {
        const tracks = stream.getTracks();
        tracks.forEach(track => track.stop());
        stream = null;
      }
      isCameraActive.value = false;
    };

    const formatDateTime = (dateString) => {
      if (!dateString) return '';
      const date = new Date(dateString);
      if (isNaN(date.getTime())) {
        return '';
      }
      return date.toLocaleString('zh-CN');
    };
    
    const getLevelType = (level) => {
      if (level && (level.includes('轻') || level.includes('微'))) {
        return 'success';
      } else if (level && (level.includes('中') || level.includes('等'))) {
        return 'warning';
      } else if (level && level.includes('重')) {
        return 'danger';
      } else {
        return 'info';
      }
    };

    const analyzeFaces = async () => {
      if (!selectedFiles.value.length) {
        error.value = '请先选择图片';
        return;
      }

      loading.value = true;
      error.value = '';
      compareResults.value = [];

      try {
        if (selectedFiles.value.length === 1) {
          await analyzeSingleFace(selectedFiles.value[0].raw);
          return;
        }

        const results = [];
        
        for (const file of selectedFiles.value) {
          try {
            const formData = new FormData();
            formData.append('image', file.raw);
            
            const response = await faceAnalysisRequest.post('/doubao/analyzeFace', formData);
            
            if (response && response.code === 1) {
              results.push(response.data);
            } else {
              console.warn(`图片 ${file.name} 分析失败:`, response?.msg || '未知错误');
              ElMessage.warning(`图片 ${file.name} 分析失败: ${response?.msg || '未知错误'}`);
            }
          } catch (err) {
            console.error(`分析图片 ${file.name} 失败:`, err);
            ElMessage.error(`分析图片 ${file.name} 失败: ${err.message || '未知错误'}`);
          }
        }
        
        if (results.length > 0) {
          compareResults.value = results;
          showResults.value = true;
          showCompareResults.value = true;
          scrollToResults();
        } else {
          error.value = '所有图片分析都失败了';
        }
      } catch (err) {
        console.error('分析请求失败，错误详情:', err);
        if (err.response) {
          if (err.response.status === 401) {
            error.value = '登录已过期，请重新登录';
          } else if (err.response.data) {
            error.value = err.response.data.msg || `分析请求失败，状态码: ${err.response.status}`;
          } else {
            error.value = `分析请求失败，状态码: ${err.response.status}`;
          }
        } else if (err.request) {
          error.value = '网络连接异常，请检查网络设置';
        } else {
          error.value = '分析请求失败: ' + (err.message || '未知错误');
        }
      } finally {
        loading.value = false;
      }
    };

    const analyzeSingleFace = async (file) => {
      try {
        const formData = new FormData();
        formData.append('image', file);
        
        console.log('发送分析请求...');
        const response = await faceAnalysisRequest.post('/doubao/analyzeFace', formData);

        console.log('API完整响应:', response);
        if (response && response.code === 1) {
          console.log('分析成功，数据:', response.data);
          analysisResult.value = response.data;
          error.value = '';
          showResults.value = true;
          showSingleResult.value = true;
          scrollToResults();
        } else {
          const errorMsg = (response && response.msg) ? response.msg : '分析失败';
          console.log('分析失败:', errorMsg);
          error.value = errorMsg;
        }
      } catch (err) {
        console.error('分析请求失败，错误详情:', err);
        if (err.response) {
          if (err.response.status === 401) {
            error.value = '登录已过期，请重新登录';
          } else if (err.response.data) {
            error.value = err.response.data.msg || `分析请求失败，状态码: ${err.response.status}`;
          } else {
            error.value = `分析请求失败，状态码: ${err.response.status}`;
          }
        } else if (err.request) {
          error.value = '网络错误，请检查网络连接';
        } else {
          error.value = `请求错误: ${err.message}`;
        }
      } finally {
        loading.value = false;
      }
    };

    const exportToPDF = () => {
      if (!analysisResult.value) {
        ElMessage.error('没有可导出的分析结果');
        return;
      }

      const printContent = document.createElement('div');
      printContent.style.padding = '20px';
      printContent.style.fontFamily = 'Arial, sans-serif';
      printContent.innerHTML = `
        <div style="text-align: center; margin-bottom: 20px;">
          <h2>腺体面容检测报告</h2>
        </div>
        <div style="margin-bottom: 20px;">
          <div style="display: flex; justify-content: center; gap: 15px; margin-bottom: 20px;">
            <span style="padding: 10px 20px; border-radius: 4px; font-weight: bold; background-color: ${analysisResult.value.isGlandFace ? '#fef0f0' : '#f0f9ff'}; color: ${analysisResult.value.isGlandFace ? '#f56c6c' : '#409eff'}">
              ${analysisResult.value.isGlandFace ? '腺体面容' : '非腺体面容'}
            </span>
            <span style="padding: 10px 20px; border-radius: 4px; font-weight: bold; background-color: ${getLevelType(analysisResult.value.level) === 'primary' ? '#ecf5ff' : getLevelType(analysisResult.value.level) === 'warning' ? '#fdf6ec' : getLevelType(analysisResult.value.level) === 'danger' ? '#fef0f0' : '#f4f4f5'}; color: ${getLevelType(analysisResult.value.level) === 'primary' ? '#409eff' : getLevelType(analysisResult.value.level) === 'warning' ? '#e6a23c' : getLevelType(analysisResult.value.level) === 'danger' ? '#f56c6c' : '#909399'}">
              ${analysisResult.value.level}
            </span>
          </div>
          
          <div style="margin-bottom: 20px;">
            <table style="width: 100%; border-collapse: collapse;">
              <tr>
                <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold;">检测时间</td>
                <td style="border: 1px solid #ddd; padding: 12px;">${formatDateTime(analysisResult.value.createTime)}</td>
              </tr>
              <tr>
                <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold;">置信度</td>
                <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold; color: #409eff;">${(analysisResult.value.confidence * 100).toFixed(2)}%</td>
              </tr>
              <tr>
                <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold;">可视化描述</td>
                <td style="border: 1px solid #ddd; padding: 12px;">${analysisResult.value.visualizationDescription}</td>
              </tr>
            </table>
          </div>
          
          <div style="text-align: center;">
            <h4 style="margin-bottom: 15px;">分析图像</h4>
            <img src="${analysisResult.value.imagePath}" style="max-width: 300px; max-height: 300px; border-radius: 8px; box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);" />
          </div>
        </div>
        <div style="margin-top: 30px; text-align: center; color: #999; font-size: 12px;">
          报告生成时间: ${new Date().toLocaleString('zh-CN')}
        </div>
      `;

      const printWindow = window.open('', '_blank');
      printWindow.document.write(`
        <html>
          <head>
            <title>腺体面容检测报告</title>
            <style>
              @media print {
                body {
                  -webkit-print-color-adjust: exact;
                  print-color-adjust: exact;
                }
              }
            </style>
          </head>
          <body>
            ${printContent.innerHTML}
            <script>
              window.onload = function() {
                window.print();
                window.close();
              }
            <\/script>
          </body>
        </html>
      `);
      printWindow.document.close();
    };

    const exportCompareResults = () => {
      if (!compareResults.value.length) {
        ElMessage.error('没有可导出的分析结果');
        return;
      }

      const printContent = document.createElement('div');
      printContent.style.padding = '20px';
      printContent.style.fontFamily = 'Arial, sans-serif';
      
      let resultsHTML = '';
      compareResults.value.forEach((result, index) => {
        resultsHTML += `
          <div style="margin-bottom: 30px; page-break-inside: avoid;">
            <h3 style="text-align: center; margin-bottom: 20px;">图片 ${index + 1}</h3>
            <div style="display: flex; gap: 20px; flex-wrap: wrap;">
              <div style="flex: 1; min-width: 300px;">
                <table style="width: 100%; border-collapse: collapse; margin-bottom: 20px;">
                  <tr>
                    <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold;">检测时间</td>
                    <td style="border: 1px solid #ddd; padding: 12px;">${formatDateTime(result.createTime)}</td>
                  </tr>
                  <tr>
                    <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold;">腺体面容</td>
                    <td style="border: 1px solid #ddd; padding: 12px;">
                      <span style="padding: 4px 8px; border-radius: 4px; background-color: ${result.isGlandFace ? '#fef0f0' : '#f0f9ff'}; color: ${result.isGlandFace ? '#f56c6c' : '#409eff'}">
                        ${result.isGlandFace ? '是' : '否'}
                      </span>
                    </td>
                  </tr>
                  <tr>
                    <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold;">等级</td>
                    <td style="border: 1px solid #ddd; padding: 12px;">
                      <span style="padding: 4px 8px; border-radius: 4px; background-color: ${getLevelType(result.level) === 'primary' ? '#ecf5ff' : getLevelType(result.level) === 'warning' ? '#fdf6ec' : getLevelType(result.level) === 'danger' ? '#fef0f0' : '#f4f4f5'}; color: ${getLevelType(result.level) === 'primary' ? '#409eff' : getLevelType(result.level) === 'warning' ? '#e6a23c' : getLevelType(result.level) === 'danger' ? '#f56c6c' : '#909399'}">
                        ${result.level}
                      </span>
                    </td>
                  </tr>
                  <tr>
                    <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold;">置信度</td>
                    <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold; color: #409eff;">${(result.confidence * 100).toFixed(2)}%</td>
                  </tr>
                  <tr>
                    <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold;">可视化描述</td>
                    <td style="border: 1px solid #ddd; padding: 12px;">${result.visualizationDescription}</td>
                  </tr>
                </table>
              </div>
              <div style="flex: 1; min-width: 300px; text-align: center;">
                <h4 style="margin-bottom: 15px;">分析图像</h4>
                <img src="${result.imagePath}" style="max-width: 300px; max-height: 300px; border-radius: 8px; box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);" />
              </div>
            </div>
          </div>
        `;
      });

      if (compareResults.value.length > 1) {
        const total = compareResults.value.length;
        const glandFaceCount = compareResults.value.filter(r => r.isGlandFace).length;
        const avgConfidence = compareResults.value.reduce((sum, r) => sum + r.confidence, 0) / total;
        
        resultsHTML += `
          <div style="margin-top: 30px; page-break-inside: avoid;">
            <h3 style="text-align: center; margin-bottom: 20px;">对比总结</h3>
            <div style="display: flex; gap: 20px;">
              <div style="flex: 1; text-align: center; padding: 20px; background-color: #f5f7fa; border-radius: 8px;">
                <div style="font-size: 24px; font-weight: bold; color: #333;">${total}</div>
                <div style="color: #666;">总检测数</div>
              </div>
              <div style="flex: 1; text-align: center; padding: 20px; background-color: #f5f7fa; border-radius: 8px;">
                <div style="font-size: 24px; font-weight: bold; color: #f56c6c;">${glandFaceCount}</div>
                <div style="color: #666;">腺体面容数</div>
              </div>
              <div style="flex: 1; text-align: center; padding: 20px; background-color: #f5f7fa; border-radius: 8px;">
                <div style="font-size: 24px; font-weight: bold; color: #409eff;">${(avgConfidence * 100).toFixed(2)}%</div>
                <div style="color: #666;">平均置信度</div>
              </div>
            </div>
          </div>
        `;
      }

      printContent.innerHTML = `
        <div style="text-align: center; margin-bottom: 30px;">
          <h1>腺体面容对比分析报告</h1>
        </div>
        <div>
          ${resultsHTML}
        </div>
        <div style="margin-top: 40px; text-align: center; color: #999; font-size: 12px;">
          报告生成时间: ${new Date().toLocaleString('zh-CN')}
        </div>
      `;

      const printWindow = window.open('', '_blank');
      printWindow.document.write(`
        <html>
          <head>
            <title>腺体面容对比分析报告</title>
            <style>
              @media print {
                body {
                  -webkit-print-color-adjust: exact;
                  print-color-adjust: exact;
                }
              }
            </style>
          </head>
          <body>
            ${printContent.innerHTML}
            <script>
              window.onload = function() {
                window.print();
                window.close();
              }
            <\/script>
          </body>
        </html>
      `);
      printWindow.document.close();
    };

    return {
      selectedFiles,
      loading,
      analysisResult,
      compareResults,
      error,
      isCameraActive,
      videoElement,
      canvasElement,
      isMobile,
      showResults,
      showSingleResult,
      showCompareResults,
      resultsSection,
      handleFileChange,
      removeFile,
      clearFiles,
      clearResults,
      startCamera,
      captureImage,
      stopCamera,
      formatDateTime,
      getLevelType,
      analyzeFaces,
      exportToPDF,
      exportCompareResults
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

.header h2 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #fff;
}

.header-subtitle {
  margin: 0;
  font-size: 14px;
  color: rgba(255,255,255,0.65);
}

.header-features {
  display: flex;
  gap: 20px;
  margin-top: 16px;
  flex-wrap: wrap;
}

.header-features .feature-item {
  display: flex;
  align-items: center;
  gap: 6px;
  background: rgba(255,255,255,0.1);
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 13px;
  color: rgba(255,255,255,0.85);
  border: 1px solid rgba(255,255,255,0.15);
}

.feature-icon { font-size: 16px !important; color: var(--med-accent); }

.upload-card {
  border-radius: var(--med-radius-lg) !important;
  box-shadow: var(--med-shadow-md) !important;
  border: 1px solid var(--med-border-light) !important;
  overflow: hidden;
}

.upload-card :deep(.el-card__header) {
  padding: 16px 22px;
  background: var(--med-bg-section);
  border-bottom: 1px solid var(--med-border-light);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-header h3 {
  margin: 0;
  color: var(--med-text);
  font-size: 16px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-subtitle { font-size: 12px; color: var(--med-text-muted); margin: 0; }

.upload-section {
  padding: 24px;
}

.camera-preview {
  position: relative;
  width: 100%;
  height: 380px;
  background: #0A0F1A;
  border-radius: var(--med-radius-md);
  overflow: hidden;
  margin-bottom: 20px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.25);
  border: 1px solid rgba(50,224,196,0.2);
}

.video-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.camera-overlay {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.25);
}

.camera-instructions {
  text-align: center;
  color: white;
}

.camera-instructions .capture-icon {
  font-size: 44px;
  margin-bottom: 12px;
  color: var(--med-accent);
}

.camera-instructions p {
  font-size: 16px;
  margin: 0 0 4px 0;
}

.instruction-detail { font-size: 12px; color: rgba(255,255,255,0.65); }

.upload-area {
  margin-bottom: 24px;
}

:deep(.el-upload-dragger) {
  border: 2px dashed var(--med-border) !important;
  border-radius: var(--med-radius-md) !important;
  background: var(--med-bg) !important;
  transition: all 0.25s ease !important;
}

:deep(.el-upload-dragger:hover) {
  border-color: var(--med-primary) !important;
  background: var(--med-primary-bg) !important;
}

.upload-placeholder {
  padding: 40px 0;
}

.upload-icon-container { margin-bottom: 16px; }

.upload-icon {
  font-size: 52px !important;
  color: var(--med-primary) !important;
}

.el-upload__text {
  font-size: 15px;
  color: var(--med-text-secondary);
  margin-bottom: 8px;
}

.el-upload__text em { color: var(--med-primary); font-style: normal; font-weight: 600; }

.el-upload__tip {
  color: var(--med-text-muted);
  font-size: 13px;
  margin-top: 10px;
  line-height: 1.6;
}

.file-limit { font-size: 12px; color: var(--med-text-muted); }

.upload-hints {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 14px;
  padding: 8px 16px;
  background: var(--med-info-light);
  border-radius: 20px;
  display: inline-flex;
  font-size: 12px;
  color: var(--med-info);
}

.hint-icon { font-size: 14px !important; }

.preview-container { padding: 8px 0; }

.multi-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  justify-content: center;
}

.preview-item {
  position: relative;
  width: 110px;
  text-align: center;
}

.preview-image-small {
  width: 100%;
  height: 110px;
  object-fit: cover;
  border-radius: var(--med-radius-md);
  border: 2px solid var(--med-border-light);
  box-shadow: var(--med-shadow-sm);
}

.file-name {
  margin: 10px 0 5px 0;
  font-size: 12px;
  color: #666;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.remove-btn {
  position: absolute;
  top: -8px;
  right: -8px;
  background: var(--med-danger) !important;
  border: none !important;
  box-shadow: 0 2px 6px rgba(231, 76, 60, 0.4) !important;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
  flex-wrap: wrap;
  padding-top: 4px;
}

.action-button {
  min-width: 130px;
  padding: 12px 22px !important;
  font-size: 14px !important;
  border-radius: var(--med-radius-md) !important;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-weight: 600 !important;
  transition: all 0.25s ease !important;
}

.action-button:hover { transform: translateY(-2px) !important; }

.analyze-button {
  background: linear-gradient(135deg, var(--med-primary) 0%, var(--med-primary-light) 100%) !important;
  border: none !important;
  box-shadow: 0 4px 14px rgba(13,115,119,0.3) !important;
}

.capture-button {
  background: linear-gradient(135deg, var(--med-ai-purple) 0%, #9C6FFF 100%) !important;
  border: none !important;
  box-shadow: 0 4px 14px rgba(124,77,255,0.3) !important;
}

.action-button.secondary {
  background: var(--med-bg) !important;
  border: 1px solid var(--med-border) !important;
  color: var(--med-text-secondary) !important;
}

.action-button.secondary:hover {
  background: #fff !important;
  border-color: var(--med-primary) !important;
  color: var(--med-primary) !important;
}

.analysis-results-section {
  margin-top: 20px;
}

.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.fade-scale-enter-from,
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(20px);
}

.single-result-container,
.compare-results-container {
  animation: slideUp 0.5s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.result-header-section,
.compare-header-section {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 24px;
  padding: 32px;
  background: linear-gradient(135deg, #ffffff 0%, #F0FDFA 100%);
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(8, 145, 178, 0.08);
  border: 1px solid rgba(34, 211, 238, 0.2);
  position: relative;
  overflow: hidden;
}

.result-header-section::before,
.compare-header-section::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(34, 211, 238, 0.08) 0%, transparent 60%);
  pointer-events: none;
}

.result-icon-wrapper,
.compare-icon-wrapper {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, rgba(34, 211, 238, 0.15) 0%, rgba(34, 197, 94, 0.1) 100%);
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid rgba(34, 211, 238, 0.3);
  flex-shrink: 0;
  position: relative;
  z-index: 1;
  transition: all 0.3s ease;
}

.result-icon-wrapper.success-bg {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.15) 0%, rgba(34, 197, 94, 0.08) 100%);
  border-color: rgba(34, 197, 94, 0.35);
}

.result-icon-wrapper.warning-bg {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.15) 0%, rgba(239, 68, 68, 0.08) 100%);
  border-color: rgba(239, 68, 68, 0.35);
}

.result-icon-wrapper:hover,
.compare-icon-wrapper:hover {
  transform: scale(1.05) rotate(5deg);
}

.result-icon {
  width: 44px;
  height: 44px;
}

.result-title-section,
.compare-title-wrapper {
  flex: 1;
  position: relative;
  z-index: 1;
}

.result-main-title,
.compare-main-title {
  margin: 0 0 6px 0;
  font-size: 28px;
  font-weight: 800;
  color: #134E4A;
  background: linear-gradient(135deg, #134E4A 0%, #0891B2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.result-subtitle,
.compare-subtitle {
  margin: 0 0 16px 0;
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
}

.result-tags {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.tag-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-radius: 28px;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.3s ease;
  cursor: default;
}

.tag-badge:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.tag-badge.success {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.15) 0%, rgba(34, 197, 94, 0.08) 100%);
  color: #166534;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.tag-badge.danger {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.15) 0%, rgba(239, 68, 68, 0.08) 100%);
  color: #991b1b;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.tag-badge.info {
  background: linear-gradient(135deg, rgba(8, 145, 178, 0.15) 0%, rgba(8, 145, 178, 0.08) 100%);
  color: #155e75;
  border: 1px solid rgba(8, 145, 178, 0.3);
}

.tag-icon {
  width: 18px;
  height: 18px;
}

.result-status-indicator {
  position: relative;
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pulse-ring {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  animation: pulse-ring 2s cubic-bezier(0.455, 0.03, 0.515, 0.955) infinite;
}

.result-status-indicator.success .pulse-ring {
  border: 3px solid rgba(34, 197, 94, 0.3);
}

.result-status-indicator.warning .pulse-ring {
  border: 3px solid rgba(239, 68, 68, 0.3);
}

@keyframes pulse-ring {
  0% {
    transform: scale(0.8);
    opacity: 1;
  }
  100% {
    transform: scale(1.5);
    opacity: 0;
  }
}

.status-dot {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  position: relative;
  z-index: 1;
}

.result-status-indicator.success .status-dot {
  background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
  box-shadow: 0 0 20px rgba(34, 197, 94, 0.4);
}

.result-status-indicator.warning .status-dot {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  box-shadow: 0 0 20px rgba(239, 68, 68, 0.4);
}

.compare-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.15) 0%, rgba(34, 197, 94, 0.08) 100%);
  border-radius: 28px;
  font-size: 15px;
  font-weight: 700;
  color: #166534;
  border: 1px solid rgba(34, 197, 94, 0.3);
  position: relative;
  z-index: 1;
}

.badge-icon {
  width: 20px;
  height: 20px;
}

.result-content-grid {
  display: grid;
  grid-template-columns: 1fr 1.2fr;
  gap: 24px;
}

.result-image-card,
.result-details-card {
  background: #fff;
  border-radius: var(--med-radius-lg);
  box-shadow: var(--med-shadow-md);
  border: 1px solid var(--med-border-light);
  overflow: hidden;
}

.image-card-header,
.details-card-header,
.summary-card-header {
  padding: 16px 20px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-bottom: 1px solid var(--med-border-light);
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--med-text);
  font-size: 15px;
}

.image-wrapper {
  padding: 24px;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
}

.result-image-preview {
  max-width: 100%;
  max-height: 400px;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.detail-item {
  padding: 16px 20px;
  border-bottom: 1px solid var(--med-border-light);
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-item.full-width {
  grid-column: 1 / -1;
}

.detail-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  color: var(--med-text-secondary);
  font-size: 14px;
}

.detail-value {
  color: var(--med-text);
  font-size: 15px;
  line-height: 1.6;
}

.detail-value.description {
  color: var(--med-text-secondary);
  padding: 12px;
  background: var(--med-bg);
  border-radius: 8px;
  border-left: 3px solid var(--med-primary);
}

.confidence-display {
  display: flex;
  align-items: center;
  gap: 16px;
}

.confidence-bar-wrapper {
  flex: 1;
  height: 12px;
  background: var(--med-bg);
  border-radius: 6px;
  overflow: hidden;
  position: relative;
}

.confidence-track {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--med-bg);
}

.confidence-fill {
  position: absolute;
  top: 0;
  left: 0;
  bottom: 0;
  background: linear-gradient(90deg, var(--med-primary) 0%, var(--med-primary-light) 100%);
  border-radius: 6px;
  transition: width 1s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.confidence-number {
  font-size: 20px;
  font-weight: 700;
  color: var(--med-primary);
  min-width: 70px;
  text-align: right;
}

.action-buttons-result,
.compare-actions {
  padding: 20px;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.export-btn,
.clear-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  border-radius: var(--med-radius-md);
  font-weight: 600;
  transition: all 0.25s ease;
}

.export-btn:hover,
.clear-btn:hover {
  transform: translateY(-2px);
}

.compare-results-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.compare-result-card {
  background: #fff;
  border-radius: var(--med-radius-lg);
  box-shadow: var(--med-shadow-md);
  border: 1px solid var(--med-border-light);
  overflow: hidden;
  transition: all 0.3s ease;
}

.compare-result-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.15);
}

.compare-card-header {
  padding: 14px 16px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-bottom: 1px solid var(--med-border-light);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sample-number {
  font-weight: 700;
  color: var(--med-text);
  font-size: 14px;
}

.compare-card-content {
  padding: 16px;
}

.compare-image-wrapper {
  margin-bottom: 12px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: var(--med-shadow-sm);
}

.compare-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.compare-details {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.compare-detail-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.compare-label {
  color: var(--med-text-secondary);
  display: flex;
  align-items: center;
  gap: 4px;
}

.mini-confidence-bar {
  flex: 1;
  height: 6px;
  background: var(--med-bg);
  border-radius: 3px;
  overflow: hidden;
}

.mini-confidence-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--med-primary) 0%, var(--med-primary-light) 100%);
  border-radius: 3px;
  transition: width 1s ease;
}

.mini-confidence-value {
  font-weight: 700;
  color: var(--med-primary);
  font-size: 13px;
  min-width: 50px;
  text-align: right;
}

.compare-summary-card {
  background: #fff;
  border-radius: var(--med-radius-lg);
  box-shadow: var(--med-shadow-md);
  border: 1px solid var(--med-border-light);
  overflow: hidden;
  margin-bottom: 24px;
}

.summary-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  padding: 24px;
}

.stat-item {
  text-align: center;
  padding: 24px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 16px;
  border: 1px solid var(--med-border-light);
  transition: all 0.3s ease;
}

.stat-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
}

.stat-icon-wrapper {
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, var(--med-primary-bg) 0%, rgba(50, 224, 196, 0.15) 100%);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
  border: 2px solid rgba(50, 224, 196, 0.25);
}

.stat-icon-wrapper.warning {
  background: linear-gradient(135deg, var(--med-danger-light) 0%, rgba(231, 76, 60, 0.1) 100%);
  border-color: rgba(231, 76, 60, 0.25);
}

.stat-icon-wrapper.info {
  background: linear-gradient(135deg, var(--med-info-light) 0%, rgba(41, 128, 185, 0.1) 100%);
  border-color: rgba(41, 128, 185, 0.25);
}

.stat-icon-wrapper .el-icon {
  font-size: 28px;
  color: var(--med-primary);
}

.stat-icon-wrapper.warning .el-icon {
  color: var(--med-danger);
}

.stat-icon-wrapper.info .el-icon {
  color: var(--med-info);
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  color: var(--med-text);
  margin-bottom: 4px;
}

.stat-number.danger {
  color: var(--med-danger);
}

.stat-number.primary {
  color: var(--med-primary);
}

.stat-label {
  font-size: 14px;
  color: var(--med-text-secondary);
  font-weight: 500;
}

.image-slot {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 200px;
  background: linear-gradient(135deg, #f5f7fa, #e4e7f1);
  color: #909399;
  font-size: 32px;
  border-radius: 12px;
}

.error-alert {
  margin-top: 20px;
  border-radius: 12px;
}

@media (max-width: 1024px) {
  .result-content-grid {
    grid-template-columns: 1fr;
  }
  
  .summary-stats {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .face-analysis-container {
    padding: 15px;
  }
  
  .header {
    padding: 20px;
    border-radius: 16px;
    margin-bottom: 20px;
  }
  
  .header h2 {
    font-size: 24px;
  }
  
  .header-subtitle {
    font-size: 14px;
  }
  
  .upload-card {
    margin: 10px;
    border-radius: 16px;
  }
  
  .upload-section {
    padding: 20px;
  }
  
  .card-header h3 {
    font-size: 18px;
  }
  
  .camera-preview {
    height: 300px;
  }
  
  .action-buttons {
    gap: 15px;
  }
  
  .action-button {
    min-width: 120px;
    padding: 12px 20px;
    font-size: 15px;
  }
  
  .multi-preview {
    gap: 15px;
  }
  
  .preview-item {
    width: 100px;
  }
  
  .preview-image-small {
    height: 100px;
  }
  
  .result-header-section,
  .compare-header-section {
    flex-direction: column;
    text-align: center;
    padding: 20px;
  }
  
  .result-main-title,
  .compare-main-title {
    font-size: 20px;
  }
  
  .compare-results-grid {
    grid-template-columns: 1fr;
  }
  
  .action-buttons-result,
  .compare-actions {
    flex-direction: column;
  }
  
  .export-btn,
  .clear-btn {
    width: 100%;
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .face-analysis-container {
    padding: 10px;
  }
  
  .header {
    padding: 15px;
    border-radius: 12px;
  }
  
  .header h2 {
    font-size: 20px;
    gap: 8px;
  }
  
  .header-subtitle {
    font-size: 12px;
  }
  
  .upload-section {
    padding: 15px;
  }
  
  .camera-preview {
    height: 250px;
  }
  
  .upload-placeholder {
    padding: 30px 0;
  }
  
  .upload-icon {
    font-size: 40px;
  }
  
  .action-buttons {
    flex-direction: column;
    gap: 12px;
  }
  
  .action-button {
    width: 100%;
    justify-content: center;
  }
  
  .multi-preview {
    gap: 10px;
  }
  
  .preview-item {
    width: 80px;
  }
  
  .preview-image-small {
    height: 80px;
  }
  
  .file-name {
    font-size: 10px;
  }
  
  .stat-item {
    padding: 20px 15px;
  }
  
  .stat-number {
    font-size: 24px;
  }
  
  .stat-label {
    font-size: 13px;
  }
}
</style>
