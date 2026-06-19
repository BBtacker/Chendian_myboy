<template>
  <div class="register-page">
    <!-- 左侧品牌区域 -->
    <div class="brand-panel">
      <div class="brand-bg-circle brand-circle-1"></div>
      <div class="brand-bg-circle brand-circle-2"></div>
      <div class="brand-content">
        <div class="brand-logo">
          <el-icon :size="52" color="#32E0C4"><UserFilled /></el-icon>
        </div>
        <div class="brand-badge">AI · 医疗检测平台</div>
        <h1 class="brand-name">加入我们<br/>开始检测</h1>
        <p class="brand-desc">
          注册账号后，您将获得完整的腺样体面容AI检测能力，支持多人管理、历史记录查询和统计分析。
        </p>
        <div class="steps">
          <div class="step-item">
            <div class="step-num">1</div>
            <div class="step-text">填写基本信息</div>
          </div>
          <div class="step-arrow">→</div>
          <div class="step-item">
            <div class="step-num">2</div>
            <div class="step-text">完成注册</div>
          </div>
          <div class="step-arrow">→</div>
          <div class="step-item">
            <div class="step-num">3</div>
            <div class="step-text">开始检测</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧注册表单 -->
    <div class="form-panel">
      <div class="register-card">
        <div class="register-card-header">
          <div class="register-avatar-badge">
            <el-icon :size="28" style="color: var(--med-primary);"><UserFilled /></el-icon>
          </div>
          <h2 class="register-title">创建账号</h2>
          <p class="register-subtitle">请填写以下信息完成注册</p>
        </div>

        <el-form
          ref="registerFormRef"
          :model="registerForm"
          :rules="registerRules"
          class="register-form"
          @submit.prevent="handleRegister"
          label-position="top"
        >
          <div class="form-grid">
            <el-form-item label="用户名" prop="username">
              <el-input
                v-model="registerForm.username"
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

            <el-form-item label="姓名" prop="name">
              <el-input
                v-model="registerForm.name"
                placeholder="请输入真实姓名"
                size="large"
                clearable
                class="med-input"
              >
                <template #prefix>
                  <el-icon style="color: var(--med-primary);"><UserFilled /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="邮箱" prop="email">
              <el-input
                v-model="registerForm.email"
                placeholder="请输入邮箱地址"
                size="large"
                clearable
                class="med-input"
              >
                <template #prefix>
                  <el-icon style="color: var(--med-primary);"><Message /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="手机号" prop="phone">
              <el-input
                v-model="registerForm.phone"
                placeholder="请输入手机号"
                size="large"
                clearable
                class="med-input"
              >
                <template #prefix>
                  <el-icon style="color: var(--med-primary);"><Phone /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="密码" prop="password">
              <el-input
                v-model="registerForm.password"
                type="password"
                placeholder="请设置密码（6-20位）"
                show-password
                size="large"
                class="med-input"
              >
                <template #prefix>
                  <el-icon style="color: var(--med-primary);"><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input
                v-model="registerForm.confirmPassword"
                type="password"
                placeholder="请再次输入密码"
                show-password
                size="large"
                class="med-input"
              >
                <template #prefix>
                  <el-icon style="color: var(--med-primary);"><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </div>

          <el-button
            type="primary"
            :loading="loading"
            size="large"
            class="register-btn"
            @click="handleRegister"
          >
            {{ loading ? '注册中...' : '立即注册' }}
          </el-button>
        </el-form>

        <p class="switch-link">
          已有账号?
          <a href="#" @click.prevent="goToLogin">立即登录</a>
        </p>

        <el-alert v-if="errorMessage" :title="errorMessage" type="error" show-icon :closable="false" class="form-alert" />
        <el-alert v-if="successMessage" :title="successMessage" type="success" show-icon :closable="false" class="form-alert" />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import request from '../utils/request'
import { UserFilled, Lock, User, Message, Phone } from '@element-plus/icons-vue'

export default {
  name: 'Register',
  components: { UserFilled, Lock, User, Message, Phone },
  emits: ['switch-to-login'],
  setup(props, { emit }) {
    const registerFormRef = ref()
    const registerForm = reactive({ username: '', name: '', email: '', phone: '', password: '', confirmPassword: '' })
    const registerRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '用户名长度应在3-20个字符之间', trigger: 'blur' }
      ],
      name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
      email: [{ type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }],
      phone: [{ pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, max: 20, message: '密码长度应在6-20个字符之间', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请确认密码', trigger: 'blur' },
        {
          validator: (rule, value, callback) => {
            if (value !== registerForm.password) callback(new Error('两次输入的密码不一致'))
            else callback()
          },
          trigger: 'blur'
        }
      ]
    }
    const loading = ref(false)
    const errorMessage = ref('')
    const successMessage = ref('')

    const handleRegister = async () => {
      if (!registerFormRef.value) return
      await registerFormRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true
          errorMessage.value = ''
          successMessage.value = ''
          try {
            const userData = { username: registerForm.username, name: registerForm.name, email: registerForm.email, phone: registerForm.phone, password: registerForm.password }
            const response = await request.post('/user/register', userData)
            if (response.code === 1) {
              successMessage.value = '注册成功，即将跳转到登录页...'
              Object.keys(registerForm).forEach(k => registerForm[k] = '')
              setTimeout(() => emit('switch-to-login'), 2500)
            } else {
              errorMessage.value = response.msg || '注册失败'
            }
          } catch (error) {
            errorMessage.value = error.message || '注册请求失败'
          } finally {
            loading.value = false
          }
        }
      })
    }

    return { registerFormRef, registerForm, registerRules, loading, errorMessage, successMessage, handleRegister, goToLogin: () => emit('switch-to-login') }
  }
}
</script>

<style scoped>
.register-page {
  display: flex;
  min-height: 100vh;
  background: var(--med-bg);
}

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

.brand-bg-circle { position: absolute; border-radius: 50%; pointer-events: none; }
.brand-circle-1 { width: 320px; height: 320px; top: -100px; right: -100px; background: radial-gradient(circle, rgba(50,224,196,0.1) 0%, transparent 70%); }
.brand-circle-2 { width: 240px; height: 240px; bottom: -60px; left: -60px; background: radial-gradient(circle, rgba(124,77,255,0.1) 0%, transparent 70%); }

.brand-content { position: relative; z-index: 1; max-width: 420px; }

.brand-logo {
  width: 88px; height: 88px;
  background: rgba(50, 224, 196, 0.12);
  border: 1px solid rgba(50, 224, 196, 0.3);
  border-radius: 24px;
  display: flex; align-items: center; justify-content: center;
  margin-bottom: 24px;
}

.brand-badge {
  display: inline-block;
  background: rgba(50,224,196,0.15); border: 1px solid rgba(50,224,196,0.3);
  color: var(--med-accent); font-size: 12px; font-weight: 600;
  padding: 5px 14px; border-radius: 20px; letter-spacing: 1px; margin-bottom: 18px;
}

.brand-name { font-size: 40px; font-weight: 800; color: #fff; line-height: 1.2; margin: 0 0 20px 0; }
.brand-desc { font-size: 15px; color: rgba(255,255,255,0.6); line-height: 1.8; margin: 0 0 40px 0; }

.steps {
  display: flex; align-items: center; gap: 12px;
  padding: 20px 24px; background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.08); border-radius: var(--med-radius-md);
}

.step-item { display: flex; flex-direction: column; align-items: center; gap: 6px; }
.step-num {
  width: 32px; height: 32px;
  background: rgba(50, 224, 196, 0.2); border: 1px solid rgba(50,224,196,0.4);
  border-radius: 50%; display: flex; align-items: center; justify-content: center;
  color: var(--med-accent); font-size: 14px; font-weight: 700;
}
.step-text { font-size: 12px; color: rgba(255,255,255,0.6); white-space: nowrap; }
.step-arrow { color: rgba(50,224,196,0.5); font-size: 18px; }

/* 右侧表单 */
.form-panel {
  width: 520px; min-width: 520px; background: #fff;
  display: flex; align-items: center; justify-content: center;
  padding: 40px; box-shadow: -4px 0 30px rgba(0,0,0,0.06);
  overflow-y: auto;
}

.register-card { width: 100%; max-width: 420px; }

.register-card-header { text-align: center; margin-bottom: 28px; }

.register-avatar-badge {
  width: 64px; height: 64px;
  background: var(--med-primary-bg); border: 2px solid rgba(13,115,119,0.2);
  border-radius: 20px; display: flex; align-items: center; justify-content: center;
  margin: 0 auto 14px;
}

.register-title { font-size: 24px; font-weight: 700; color: var(--med-text); margin: 0 0 6px 0; }
.register-subtitle { font-size: 13px; color: var(--med-text-muted); margin: 0; }

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 16px;
}

.form-grid :deep(.el-form-item) { margin-bottom: 16px; }
.form-grid :deep(.el-form-item__label) { font-weight: 600; color: var(--med-text); font-size: 13px; padding-bottom: 6px; }

:deep(.med-input .el-input__wrapper) {
  border-radius: var(--med-radius-sm) !important;
  border: 1.5px solid var(--med-border) !important;
  box-shadow: none !important;
  height: 44px;
  transition: all 0.25s ease;
}
:deep(.med-input .el-input__wrapper:hover) { border-color: var(--med-primary-light) !important; }
:deep(.med-input .el-input__wrapper.is-focus) {
  border-color: var(--med-primary) !important;
  box-shadow: 0 0 0 3px var(--med-primary-bg) !important;
}

.register-btn {
  width: 100%; height: 48px !important;
  font-size: 16px !important; font-weight: 600 !important;
  border-radius: var(--med-radius-sm) !important;
  background: linear-gradient(135deg, var(--med-primary) 0%, var(--med-primary-light) 100%) !important;
  border: none !important; letter-spacing: 1px;
  box-shadow: 0 4px 16px rgba(13,115,119,0.3) !important;
  margin-top: 8px; margin-bottom: 20px;
}

.register-btn:hover {
  transform: translateY(-2px) !important;
  box-shadow: 0 8px 24px rgba(13,115,119,0.4) !important;
}

.switch-link { text-align: center; font-size: 14px; color: var(--med-text-muted); margin: 0; }
.switch-link a { color: var(--med-primary); font-weight: 600; text-decoration: none; margin-left: 4px; }
.switch-link a:hover { text-decoration: underline; }

.form-alert { margin-top: 16px; border-radius: var(--med-radius-sm) !important; }

@media (max-width: 900px) {
  .brand-panel { display: none; }
  .form-panel { width: 100%; min-width: 0; }
  .form-grid { grid-template-columns: 1fr; }
}
</style>
