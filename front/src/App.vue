<template>
  <div id="app">
    <div v-if="!isLoggedIn && currentView === 'login'">
      <Login @login-success="handleLoginSuccess" @switch-to-register="currentView = 'register'" />
    </div>
    <div v-else-if="!isLoggedIn && currentView === 'register'">
      <Register @switch-to-login="currentView = 'login'" />
    </div>
    <div v-else class="main-layout">
      <!-- 侧边栏 -->
      <aside class="sidebar">
        <!-- 系统标题 -->
        <div class="sidebar-brand">
          <div class="brand-icon">
            <el-icon :size="32" color="#32E0C4"><Camera /></el-icon>
          </div>
          <div class="brand-text">
            <div class="brand-title">AI腺样体</div>
            <div class="brand-subtitle">面容检测系统</div>
          </div>
        </div>

        <!-- 分隔线 -->
        <div class="sidebar-divider"></div>

        <!-- 导航菜单 -->
        <nav class="sidebar-nav">
          <div
            v-for="item in menuItems"
            :key="item.key"
            :class="['nav-item', { active: activePage === item.key }]"
            @click="switchToPage(item.key)"
          >
            <div class="nav-indicator"></div>
            <el-icon class="nav-icon"><component :is="item.icon" /></el-icon>
            <span class="nav-label">{{ item.label }}</span>
            <div v-if="activePage === item.key" class="nav-badge"></div>
          </div>
        </nav>

        <!-- 底部用户操作 -->
        <div class="sidebar-footer">
          <div class="sidebar-divider"></div>
          <div class="sidebar-user" v-if="userInfo">
            <el-avatar
              :src="userInfo.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'"
              :size="36"
              class="sidebar-avatar"
            />
            <div class="sidebar-username">{{ userInfo.name || userInfo.username }}</div>
          </div>
          <button class="logout-btn" @click="handleLogout">
            <el-icon><SwitchButton /></el-icon>
            <span>退出登录</span>
          </button>
        </div>
      </aside>

      <!-- 主内容区 -->
      <div class="content-wrapper">
        <!-- 顶部栏 -->
        <header class="topbar">
          <div class="topbar-left">
            <div class="page-breadcrumb">
              <span class="breadcrumb-system">检测系统</span>
              <el-icon class="breadcrumb-sep"><ArrowRight /></el-icon>
              <span class="breadcrumb-page">{{ getPageTitle(activePage) }}</span>
            </div>
            <h1 class="page-title">{{ getPageTitle(activePage) }}</h1>
          </div>

          <div class="topbar-right">
            <div class="topbar-date">{{ currentDate }}</div>
            <div class="topbar-user" v-if="userInfo">
              <el-dropdown @command="handleUserCommand" trigger="click">
                <div class="user-pill">
                  <el-avatar
                    :src="userInfo.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'"
                    :size="38"
                    class="user-pill-avatar"
                  />
                  <div class="user-pill-info">
                    <div class="user-pill-name">{{ userInfo.name || userInfo.username }}</div>
                    <div class="user-pill-role">医生 / 系统用户</div>
                  </div>
                  <el-icon class="user-pill-arrow"><ArrowDown /></el-icon>
                </div>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="profile">
                      <el-icon><User /></el-icon> 个人中心
                    </el-dropdown-item>
                    <el-dropdown-item divided command="logout">
                      <el-icon><SwitchButton /></el-icon> 退出登录
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </header>

        <!-- 页面内容 -->
        <main class="page-main">
          <transition name="slide-fade" mode="out-in">
            <template v-if="shouldKeepAlive">
              <keep-alive include="FaceAnalysis">
                <component :is="currentComponent" :key="activePage" />
              </keep-alive>
            </template>
            <template v-else>
              <component :is="currentComponent" :key="activePage" />
            </template>
          </transition>
        </main>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch, computed } from 'vue'
import {
  Camera,
  Cpu,
  Document,
  User,
  TrendCharts,
  SwitchButton,
  ArrowDown,
  ArrowRight,
  Setting,
  DataAnalysis
} from '@element-plus/icons-vue'
import request from './utils/request'
import FaceAnalysis from './components/FaceAnalysis.vue'
import ChatAssistant from './components/ChatAssistant.vue'
import TestResult from './components/TestResult.vue'
import Statistics from './components/Statistics.vue'
import UserProfile from './components/UserProfile.vue'
import Login from './components/Login.vue'
import Register from './components/Register.vue'

export default {
  name: 'App',
  components: {
    FaceAnalysis, ChatAssistant, TestResult, Statistics, UserProfile,
    Login, Register,
    Camera, Cpu, Document, User, TrendCharts, SwitchButton, ArrowDown, ArrowRight, Setting, DataAnalysis
  },
  setup() {
    const isLoggedIn = ref(false)
    const currentView = ref('login')
    const activePage = ref('face-analysis')
    const userInfo = ref(null)
    const shouldKeepAlive = ref(true)

    const menuItems = [
      { key: 'face-analysis', label: '面容分析', icon: 'Camera' },
      { key: 'chat-assistant', label: 'AI助手', icon: 'Cpu' },
      { key: 'test-result', label: '检测记录', icon: 'Document' },
      { key: 'statistics', label: '统计分析', icon: 'TrendCharts' },
      { key: 'user-profile', label: '个人中心', icon: 'User' }
    ]

    const currentDate = computed(() => {
      const now = new Date()
      return now.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'short' })
    })

    const currentComponent = computed(() => {
      const componentMap = {
        'face-analysis': FaceAnalysis,
        'chat-assistant': ChatAssistant,
        'test-result': TestResult,
        'statistics': Statistics,
        'user-profile': UserProfile
      }
      return componentMap[activePage.value] || FaceAnalysis
    })

    onMounted(() => {
      const token = localStorage.getItem('token')
      isLoggedIn.value = !!token
      if (token) fetchUserInfo()
    })

    watch(isLoggedIn, (newVal) => {
      if (newVal) fetchUserInfo()
      else userInfo.value = null
    })

    const fetchUserInfo = async () => {
      try {
        const response = await request.get('/user')
        if (response.code === 1) userInfo.value = response.data
      } catch (error) {
        console.error('获取用户信息失败:', error)
      }
    }

    const handleLoginSuccess = () => { isLoggedIn.value = true }

    const handleLogout = () => {
      localStorage.removeItem('token')
      isLoggedIn.value = false
      activePage.value = 'face-analysis'
      shouldKeepAlive.value = false
      setTimeout(() => {
        shouldKeepAlive.value = true
      }, 0)
    }

    const switchToPage = (page) => { activePage.value = page }

    const handleUserCommand = (command) => {
      if (command === 'profile') activePage.value = 'user-profile'
      else if (command === 'logout') handleLogout()
    }

    const getPageTitle = (page) => {
      const titles = {
        'face-analysis': '面容分析',
        'chat-assistant': 'AI智能助手',
        'test-result': '检测记录',
        'statistics': '统计分析',
        'user-profile': '个人中心'
      }
      return titles[page] || 'AI腺样体面容检测系统'
    }

    return {
      isLoggedIn, currentView, activePage, userInfo, menuItems, currentDate, currentComponent, shouldKeepAlive,
      handleLoginSuccess, handleLogout, switchToPage, handleUserCommand, getPageTitle
    }
  }
}
</script>

<style>
/* ===== CSS变量定义 ===== */
:root {
  /* 颜色变量 */
  --med-primary: #1976D2;
  --med-primary-light: #2196F3;
  --med-primary-bg: #e3f2fd;
  --med-accent: #32E0C4;
  --med-ai-purple: #7C4DFF;
  --med-text: #333333;
  --med-text-secondary: #666666;
  --med-text-muted: #999999;
  --med-bg: #f5f7fa;
  --med-bg-section: #fafbfc;
  --med-border: #e0e0e0;
  --med-border-light: #ebeef5;
  --med-sidebar: #1a2035;
  --med-sidebar-mid: #1e2742;
  --med-sidebar-light: #232c4a;
  --med-sidebar-text: rgba(255, 255, 255, 0.7);
  --med-sidebar-active: rgba(50, 224, 196, 0.15);
  --med-success: #27ae60;
  --med-success-light: #e8f5e8;
  --med-danger: #e74c3c;
  --med-danger-light: #fce8e8;
  --med-warning: #f39c12;
  --med-warning-light: #fef5e7;
  --med-info: #2980b9;
  --med-info-light: #e8f4f8;
  
  /* 阴影变量 */
  --med-shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.1);
  --med-shadow-md: 0 4px 16px rgba(0, 0, 0, 0.12);
  
  /* 圆角变量 */
  --med-radius-sm: 4px;
  --med-radius-md: 8px;
  --med-radius-lg: 12px;
  --med-radius-xl: 16px;
}

/* ===== 全局重置 ===== */
html, body, #app {
  height: 100%;
  margin: 0;
  padding: 0;
  font-family: 'PingFang SC', 'Microsoft YaHei', 'Helvetica Neue', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--med-text);
  background-color: var(--med-bg);
}

* { box-sizing: border-box; }

/* ===== 主布局 ===== */
.main-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
  background: var(--med-bg);
}

/* ===== 侧边栏 ===== */
.sidebar {
  width: 220px;
  min-width: 220px;
  height: 100vh;
  background: linear-gradient(180deg, var(--med-sidebar) 0%, var(--med-sidebar-mid) 60%, var(--med-sidebar-light) 100%);
  display: flex;
  flex-direction: column;
  box-shadow: 4px 0 20px rgba(0, 0, 0, 0.2);
  position: relative;
  z-index: 10;
  overflow: hidden;
}

.sidebar::before {
  content: '';
  position: absolute;
  top: -80px;
  right: -80px;
  width: 200px;
  height: 200px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(50, 224, 196, 0.08) 0%, transparent 70%);
  pointer-events: none;
}

/* 品牌区域 */
.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 20px 20px;
}

.brand-icon {
  width: 48px;
  height: 48px;
  background: rgba(50, 224, 196, 0.12);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(50, 224, 196, 0.25);
  flex-shrink: 0;
}

.brand-text {
  flex: 1;
  overflow: hidden;
}

.brand-title {
  font-size: 16px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 1px;
  line-height: 1.3;
}

.brand-subtitle {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 0.5px;
  margin-top: 2px;
}

/* 分隔线 */
.sidebar-divider {
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.1), transparent);
  margin: 0 16px;
}

/* 导航 */
.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 11px 14px;
  border-radius: var(--med-radius-md);
  cursor: pointer;
  position: relative;
  transition: all 0.25s ease;
  color: var(--med-sidebar-text);
  user-select: none;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.9);
}

.nav-item.active {
  background: var(--med-sidebar-active);
  color: var(--med-accent);
  box-shadow: inset 0 0 0 1px rgba(50, 224, 196, 0.2);
}

.nav-indicator {
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%) scaleY(0);
  width: 3px;
  height: 22px;
  background: var(--med-accent);
  border-radius: 0 3px 3px 0;
  transition: transform 0.25s ease;
}

.nav-item.active .nav-indicator {
  transform: translateY(-50%) scaleY(1);
}

.nav-icon {
  font-size: 18px !important;
  flex-shrink: 0;
}

.nav-label {
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.3px;
}

.nav-badge {
  width: 6px;
  height: 6px;
  background: var(--med-accent);
  border-radius: 50%;
  margin-left: auto;
}

/* 侧边栏底部 */
.sidebar-footer {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.sidebar-user {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 10px;
}

.sidebar-avatar {
  border: 2px solid rgba(50, 224, 196, 0.4);
  flex-shrink: 0;
}

.sidebar-username {
  font-size: 13px;
  color: rgba(255,255,255,0.7);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 10px 14px;
  border: 1px solid rgba(231, 76, 60, 0.35);
  background: rgba(231, 76, 60, 0.08);
  color: #FC8680;
  border-radius: var(--med-radius-md);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.25s ease;
}

.logout-btn:hover {
  background: rgba(231, 76, 60, 0.18);
  border-color: rgba(231, 76, 60, 0.6);
  color: #FF6B6B;
}

/* ===== 内容区域 ===== */
.content-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

/* ===== 顶部栏 ===== */
.topbar {
  height: 66px;
  min-height: 66px;
  background: #fff;
  border-bottom: 1px solid var(--med-border-light);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 28px;
  box-shadow: var(--med-shadow-sm);
  z-index: 5;
}

.topbar-left {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.page-breadcrumb {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 2px;
}

.breadcrumb-system {
  font-size: 11px;
  color: var(--med-text-muted);
}

.breadcrumb-sep {
  font-size: 10px !important;
  color: var(--med-text-muted);
}

.breadcrumb-page {
  font-size: 11px;
  color: var(--med-primary);
  font-weight: 500;
}

.page-title {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: var(--med-text);
  line-height: 1.2;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.topbar-date {
  font-size: 13px;
  color: var(--med-text-secondary);
  background: var(--med-bg);
  padding: 6px 12px;
  border-radius: 20px;
  border: 1px solid var(--med-border);
}

.user-pill {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 14px 6px 6px;
  border-radius: 30px;
  border: 1px solid var(--med-border);
  background: var(--med-bg);
  cursor: pointer;
  transition: all 0.25s ease;
}

.user-pill:hover {
  border-color: var(--med-primary-light);
  background: var(--med-primary-bg);
  box-shadow: var(--med-shadow-sm);
}

.user-pill-avatar {
  border: 2px solid var(--med-primary-light);
}

.user-pill-info {
  display: flex;
  flex-direction: column;
}

.user-pill-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--med-text);
  line-height: 1.2;
}

.user-pill-role {
  font-size: 11px;
  color: var(--med-text-muted);
}

.user-pill-arrow {
  font-size: 12px !important;
  color: var(--med-text-muted);
  margin-left: 2px;
}

/* ===== 页面主体 ===== */
.page-main {
  flex: 1;
  overflow-y: auto;
  background: var(--med-bg);
}

/* ===== Element Plus 全局覆盖 ===== */
.el-container, .el-aside, .el-header, .el-main {
  margin: 0; padding: 0;
}

/* 按钮样式优化 */
.el-button {
  border-radius: var(--med-radius-sm) !important;
  font-weight: 500 !important;
  transition: all 0.25s ease !important;
}

.el-button--primary {
  background: var(--med-primary) !important;
  border-color: var(--med-primary) !important;
}

.el-button--primary:hover {
  background: var(--med-primary-light) !important;
  border-color: var(--med-primary-light) !important;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(13, 115, 119, 0.3) !important;
}

/* 卡片优化 */
.el-card {
  border-radius: var(--med-radius-lg) !important;
  border: 1px solid var(--med-border-light) !important;
  box-shadow: var(--med-shadow-md) !important;
}

/* Tag优化 */
.el-tag--success {
  background: var(--med-success-light) !important;
  border-color: rgba(39, 174, 96, 0.25) !important;
  color: var(--med-success) !important;
}

.el-tag--danger {
  background: var(--med-danger-light) !important;
  border-color: rgba(231, 76, 60, 0.25) !important;
  color: var(--med-danger) !important;
}

.el-tag--warning {
  background: var(--med-warning-light) !important;
  border-color: rgba(243, 156, 18, 0.25) !important;
  color: var(--med-warning) !important;
}

.el-tag--info {
  background: var(--med-info-light) !important;
  border-color: rgba(41, 128, 185, 0.25) !important;
  color: var(--med-info) !important;
}

/* 分页优化 */
.el-pagination.is-background .el-pager li.is-active {
  background-color: var(--med-primary) !important;
}

/* 表单输入框 */
.el-input__wrapper:focus-within,
.el-textarea__inner:focus {
  box-shadow: 0 0 0 1px var(--med-primary-light) inset !important;
}

/* ===== 页面切换过渡特效 ===== */
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.slide-fade-enter-from {
  opacity: 0;
  transform: translateX(40px) scale(0.95);
}

.slide-fade-leave-to {
  opacity: 0;
  transform: translateX(-40px) scale(0.95);
}

.page-main {
  position: relative;
  overflow-y: auto;
}

/* 登录注册页面过渡 */
.login-register-fade-enter-active,
.login-register-fade-leave-active {
  transition: all 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.login-register-fade-enter-from,
.login-register-fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* 减少动画偏好 */
@media (prefers-reduced-motion: reduce) {
  .slide-fade-enter-active,
  .slide-fade-leave-active,
  .login-register-fade-enter-active,
  .login-register-fade-leave-active {
    transition: none;
  }
  
  .slide-fade-enter-from,
  .slide-fade-leave-to,
  .login-register-fade-enter-from,
  .login-register-fade-leave-to {
    opacity: 1;
    transform: none;
  }
}
</style>
