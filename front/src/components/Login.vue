<template>
  <div class="login-page">
    <!-- 左侧品牌区域 -->
    <div class="brand-panel">
      <div class="brand-bg-circle brand-circle-1"></div>
      <div class="brand-bg-circle brand-circle-2"></div>
      <div class="brand-content">
        <div class="brand-logo">
          <div class="brand-icon">
            <div class="icon-circle"></div>
            <div class="icon-wave"></div>
          </div>
        </div>
        <div class="brand-badge">AI · 医疗检测平台</div>
        <h1 class="brand-name">腺体面容识别系统</h1>
        <p class="brand-desc">
          基于先进的人工智能与生物特征识别技术，为临床医生提供精准、快速的腺样体面容智能检测解决方案。
        </p>

        <div class="features">
          <div class="feature-card">
            <div class="feature-icon-wrap" style="background: rgba(39,174,96,0.18);">
              <el-icon style="color: #2ECC71; font-size: 20px;"><Lock /></el-icon>
            </div>
            <div>
              <div class="feature-name">高度安全</div>
              <div class="feature-desc">多重加密保护，防止患者信息泄露</div>
            </div>
          </div>
          <div class="feature-card">
            <div class="feature-icon-wrap" style="background: rgba(50,224,196,0.18);">
              <el-icon style="color: #32E0C4; font-size: 20px;"><Lightning /></el-icon>
            </div>
            <div>
              <div class="feature-name">快速识别</div>
              <div class="feature-desc">毫秒级响应，实时输出检测结果</div>
            </div>
          </div>
          <div class="feature-card">
            <div class="feature-icon-wrap" style="background: rgba(124,77,255,0.18);">
              <el-icon style="color: #7C4DFF; font-size: 20px;"><Check /></el-icon>
            </div>
            <div>
              <div class="feature-name">精准可靠</div>
              <div class="feature-desc">AI算法保障，识别准确率达78%</div>
            </div>
          </div>
        </div>

        <div class="brand-footer">
          <span class="brand-tag">专业</span>
          <span class="brand-tag">精准</span>
          <span class="brand-tag">智能</span>
        </div>
      </div>
    </div>

    <!-- 右侧登录表单 -->
    <div class="form-panel">
      <div class="form-illustration">
        <div class="illustration-bg"></div>
        <div class="illustration-elements">
          <div class="illustration-circle"></div>
          <div class="illustration-wave"></div>
          <div class="illustration-dots"></div>
        </div>
      </div>
      <div class="login-card">
        <div class="login-card-header">
          <div class="login-avatar-badge">
            <el-icon :size="28" style="color: var(--med-primary);"><User /></el-icon>
          </div>
          <h2 class="login-title">欢迎登录</h2>
          <p class="login-subtitle">请输入您的账号信息以继续使用</p>
        </div>

        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          class="login-form"
          @submit.prevent="handleLogin"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名"
              size="large"
              clearable
              class="med-input"
            >
              <template #prefix>
                <el-icon style="color: var(--med-primary);"><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              show-password
              size="large"
              class="med-input"
            >
              <template #prefix>
                <el-icon style="color: var(--med-primary);"><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>

          <div class="form-options">
            <el-checkbox v-model="rememberMe" class="remember-check">记住我</el-checkbox>
            <a href="#" class="forgot-link">忘记密码?</a>
          </div>

          <el-button
            type="primary"
            native-type="submit"
            :loading="loading"
            size="large"
            class="login-btn"
            @click="handleLogin"
          >
            <el-icon v-if="!loading" style="margin-right: 6px;"><Right /></el-icon>
            {{ loading ? '登录中...' : '登录系统' }}
          </el-button>
        </el-form>

        <div class="divider-line">
          <span>其他登录方式</span>
        </div>

        <div class="social-row">
          <button class="social-btn">
            <el-icon><ChatRound /></el-icon>
          </button>
          <button class="social-btn">
            <el-icon><ChatDotSquare /></el-icon>
          </button>
          <button class="social-btn">
            <el-icon><Message /></el-icon>
          </button>
        </div>

        <p class="switch-link">
          还没有账号?
          <a href="#" @click.prevent="goToRegister">立即注册</a>
        </p>

        <el-alert
          v-if="errorMessage"
          :title="errorMessage"
          type="error"
          show-icon
          :closable="false"
          class="login-error"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import request from '../utils/request'
import { UserFilled, Lock, Lightning, Check, User, Right, ChatRound, ChatDotSquare, Message } from '@element-plus/icons-vue'

export default {
  name: 'Login',
  components: { UserFilled, Lock, Lightning, Check, User, Right, ChatRound, ChatDotSquare, Message },
  emits: ['login-success', 'switch-to-register'],
  setup(props, { emit }) {
    const loginFormRef = ref()
    const loginForm = reactive({ username: '', password: '' })
    const loginRules = {
      username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
      password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
    }
    const rememberMe = ref(false)
    const loading = ref(false)
    const errorMessage = ref('')

    const handleLogin = async () => {
      if (!loginFormRef.value) return
      await loginFormRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true
          errorMessage.value = ''
          try {
            const response = await request.post('/login', loginForm)
            if (response.code === 1) {
              localStorage.setItem('token', response.data)
              emit('login-success')
              window.location.reload()
            } else {
              errorMessage.value = response.msg || '登录失败'
            }
          } catch (error) {
            errorMessage.value = error.message || '登录请求失败，请检查网络连接'
          } finally {
            loading.value = false
          }
        }
      })
    }

    return { loginFormRef, loginForm, loginRules, rememberMe, loading, errorMessage, handleLogin, goToRegister: () => emit('switch-to-register') }
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  min-height: 100vh;
  background: var(--med-bg);
}

/* ===== 左侧品牌区域 ===== */
.brand-panel {
  flex: 1;
  background: linear-gradient(145deg, var(--med-sidebar) 0%, var(--med-sidebar-mid) 50%, #0D4F5E 100%);
  padding: 60px 52px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.brand-bg-circle {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
}

.brand-circle-1 {
  width: 320px;
  height: 320px;
  top: -100px;
  right: -100px;
  background: radial-gradient(circle, rgba(50, 224, 196, 0.1) 0%, transparent 70%);
}

.brand-circle-2 {
  width: 240px;
  height: 240px;
  bottom: -60px;
  left: -60px;
  background: radial-gradient(circle, rgba(124, 77, 255, 0.1) 0%, transparent 70%);
}

.brand-content {
  position: relative;
  z-index: 1;
  max-width: 420px;
}

.brand-logo {
  width: 88px;
  height: 88px;
  background: rgba(50, 224, 196, 0.12);
  border: 1px solid rgba(50, 224, 196, 0.3);
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
  overflow: hidden;
}

.brand-icon {
  position: relative;
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-circle {
  width: 30px;
  height: 30px;
  background: var(--med-accent);
  border-radius: 50%;
  position: relative;
  z-index: 2;
  box-shadow: 0 0 20px rgba(50, 224, 196, 0.6);
  animation: pulse 2s infinite;
}

.icon-wave {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(50, 224, 196, 0.3) 0%, transparent 70%);
  animation: wave 3s infinite;
}

@keyframes pulse {
  0% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.1); opacity: 0.8; }
  100% { transform: scale(1); opacity: 1; }
}

@keyframes wave {
  0% { transform: scale(1); opacity: 0.6; }
  50% { transform: scale(1.2); opacity: 0.3; }
  100% { transform: scale(1); opacity: 0.6; }
}

.brand-badge {
  display: inline-block;
  background: rgba(50, 224, 196, 0.15);
  border: 1px solid rgba(50, 224, 196, 0.3);
  color: var(--med-accent);
  font-size: 12px;
  font-weight: 600;
  padding: 5px 14px;
  border-radius: 20px;
  letter-spacing: 1px;
  margin-bottom: 18px;
}

.brand-name {
  font-size: 40px;
  font-weight: 800;
  color: #fff;
  line-height: 1.2;
  margin: 0 0 20px 0;
  letter-spacing: 1px;
}

.brand-desc {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.6);
  line-height: 1.8;
  margin: 0 0 36px 0;
}

.features {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 36px;
}

.feature-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 18px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: var(--med-radius-md);
  transition: all 0.25s ease;
}

.feature-card:hover {
  background: rgba(255, 255, 255, 0.09);
  transform: translateX(4px);
}

.feature-icon-wrap {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.feature-name {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 3px;
}

.feature-desc {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  line-height: 1.4;
}

.brand-footer {
  display: flex;
  gap: 10px;
}

.brand-tag {
  background: rgba(255,255,255,0.08);
  border: 1px solid rgba(255,255,255,0.15);
  color: rgba(255,255,255,0.6);
  padding: 5px 14px;
  border-radius: 20px;
  font-size: 12px;
  letter-spacing: 1px;
}

/* ===== 右侧表单区域 ===== */
.form-panel {
  width: 480px;
  min-width: 480px;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  box-shadow: -4px 0 30px rgba(0,0,0,0.06);
  position: relative;
  overflow: hidden;
}

.form-illustration {
  position: absolute;
  top: 0;
  right: 0;
  width: 200px;
  height: 200px;
  pointer-events: none;
}

.illustration-bg {
  position: absolute;
  top: -50px;
  right: -50px;
  width: 150px;
  height: 150px;
  background: radial-gradient(circle, rgba(50, 224, 196, 0.1) 0%, transparent 70%);
  border-radius: 50%;
}

.illustration-elements {
  position: relative;
  width: 100%;
  height: 100%;
}

.illustration-circle {
  position: absolute;
  top: 20px;
  right: 30px;
  width: 40px;
  height: 40px;
  background: var(--med-primary);
  border-radius: 50%;
  opacity: 0.1;
  animation: float 4s infinite ease-in-out;
}

.illustration-wave {
  position: absolute;
  top: 60px;
  right: 20px;
  width: 60px;
  height: 60px;
  border: 2px solid var(--med-accent);
  border-radius: 50%;
  opacity: 0.1;
  animation: wave 3s infinite ease-in-out;
}

.illustration-dots {
  position: absolute;
  top: 40px;
  right: 80px;
  width: 10px;
  height: 10px;
  background: var(--med-ai-purple);
  border-radius: 50%;
  opacity: 0.2;
  animation: pulse 2s infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-10px) rotate(180deg); }
}

@keyframes wave {
  0%, 100% { transform: scale(1) rotate(0deg); opacity: 0.1; }
  50% { transform: scale(1.2) rotate(180deg); opacity: 0.15; }
}

@keyframes pulse {
  0%, 100% { transform: scale(1); opacity: 0.2; }
  50% { transform: scale(1.2); opacity: 0.3; }
}

.login-card {
  width: 100%;
  max-width: 380px;
}

.login-card-header {
  text-align: center;
  margin-bottom: 36px;
}

.login-avatar-badge {
  width: 64px;
  height: 64px;
  background: var(--med-primary-bg);
  border: 2px solid rgba(13, 115, 119, 0.2);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.login-title {
  font-size: 26px;
  font-weight: 700;
  color: var(--med-text);
  margin: 0 0 8px 0;
}

.login-subtitle {
  font-size: 14px;
  color: var(--med-text-muted);
  margin: 0;
}

.login-form {
  margin-bottom: 4px;
}

:deep(.med-input .el-input__wrapper) {
  border-radius: var(--med-radius-sm) !important;
  border: 1.5px solid var(--med-border) !important;
  box-shadow: none !important;
  height: 48px;
  transition: all 0.25s ease;
}

:deep(.med-input .el-input__wrapper:hover) {
  border-color: var(--med-primary-light) !important;
}

:deep(.med-input .el-input__wrapper.is-focus) {
  border-color: var(--med-primary) !important;
  box-shadow: 0 0 0 3px var(--med-primary-bg) !important;
}

:deep(.med-input .el-input__inner) {
  font-size: 15px;
  color: var(--med-text);
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: -4px 0 20px;
}

:deep(.remember-check .el-checkbox__label) {
  font-size: 13px;
  color: var(--med-text-secondary);
}

:deep(.remember-check .el-checkbox__inner) {
  border-color: var(--med-border);
}

:deep(.remember-check .el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: var(--med-primary);
  border-color: var(--med-primary);
}

.forgot-link {
  font-size: 13px;
  color: var(--med-primary);
  text-decoration: none;
  font-weight: 500;
}

.forgot-link:hover { text-decoration: underline; }

.login-btn {
  width: 100%;
  height: 48px !important;
  font-size: 16px !important;
  font-weight: 600 !important;
  border-radius: var(--med-radius-sm) !important;
  background: linear-gradient(135deg, var(--med-primary) 0%, var(--med-primary-light) 100%) !important;
  border: none !important;
  letter-spacing: 1px;
  box-shadow: 0 4px 16px rgba(13, 115, 119, 0.3) !important;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-btn:hover {
  transform: translateY(-2px) !important;
  box-shadow: 0 8px 24px rgba(13, 115, 119, 0.4) !important;
}

.divider-line {
  position: relative;
  text-align: center;
  margin-bottom: 20px;
}

.divider-line::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: var(--med-border);
}

.divider-line span {
  position: relative;
  background: #fff;
  padding: 0 16px;
  font-size: 13px;
  color: var(--med-text-muted);
}

.social-row {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.social-btn {
  flex: 1;
  height: 44px;
  border: 1.5px solid var(--med-border);
  background: var(--med-bg);
  border-radius: var(--med-radius-sm);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--med-text-secondary);
  font-size: 18px;
  transition: all 0.25s ease;
}

.social-btn:hover {
  border-color: var(--med-primary);
  color: var(--med-primary);
  background: var(--med-primary-bg);
}

.switch-link {
  text-align: center;
  font-size: 14px;
  color: var(--med-text-muted);
  margin: 0;
}

.switch-link a {
  color: var(--med-primary);
  font-weight: 600;
  text-decoration: none;
  margin-left: 4px;
}

.switch-link a:hover { text-decoration: underline; }

.login-error {
  margin-top: 16px;
  border-radius: var(--med-radius-sm) !important;
}

@media (max-width: 900px) {
  .brand-panel { display: none; }
  .form-panel { width: 100%; min-width: 0; }
}
</style>
