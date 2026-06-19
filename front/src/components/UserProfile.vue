<template>
  <div class="profile-page">
    <!-- 页头 banner -->
    <div class="profile-banner">
      <div class="banner-bg-circle banner-circle-1"></div>
      <div class="banner-bg-circle banner-circle-2"></div>
      <div class="banner-content">
        <div class="banner-avatar-wrap">
          <el-avatar
            :src="profileForm.avatar || defaultAvatar"
            :size="90"
            shape="circle"
            class="banner-avatar"
          />
          <div class="avatar-status-dot"></div>
        </div>
        <div class="banner-info">
          <h2 class="banner-name">{{ profileForm.name || profileForm.username || '用户' }}</h2>
          <p class="banner-username">@{{ profileForm.username }}</p>
          <div class="banner-tags">
            <span class="banner-tag"><el-icon style="font-size: 11px;"><User /></el-icon> 系统用户</span>
            <span class="banner-tag banner-tag-teal">
              <el-icon style="font-size: 11px;"><Clock /></el-icon>
              注册于 {{ formattedCreateTime }}
            </span>
          </div>
        </div>
        <div class="banner-actions">
          <el-button type="primary" @click="openEditDialog" class="action-btn">
            <el-icon><Edit /></el-icon> 编辑资料
          </el-button>
          <el-button @click="openAvatarDialog" class="action-btn-outline">
            <el-icon><Picture /></el-icon> 更换头像
          </el-button>
        </div>
      </div>
    </div>

    <!-- 个人信息卡片 -->
    <el-card class="info-card">
      <template #header>
        <div class="section-card-header">
          <el-icon style="color: var(--med-primary);"><User /></el-icon>
          <span>个人信息</span>
        </div>
      </template>

      <div class="info-grid">
        <div class="info-item">
          <div class="info-item-label">
            <div class="info-icon-wrap info-icon-blue">
              <el-icon :size="14" style="color: var(--med-info);"><Message /></el-icon>
            </div>
            邮箱
          </div>
          <div class="info-item-value">{{ profileForm.email || '未设置' }}</div>
        </div>

        <div class="info-item">
          <div class="info-item-label">
            <div class="info-icon-wrap info-icon-green">
              <el-icon :size="14" style="color: var(--med-success);"><Phone /></el-icon>
            </div>
            手机号
          </div>
          <div class="info-item-value">{{ profileForm.phone || '未设置' }}</div>
        </div>

        <div class="info-item">
          <div class="info-item-label">
            <div class="info-icon-wrap" :style="{ background: profileForm.gender === 1 ? 'rgba(41,128,185,0.12)' : 'rgba(231,76,60,0.1)' }">
              <el-icon :size="14" :style="{ color: profileForm.gender === 1 ? '#2980B9' : '#E74C3C' }">
                <component :is="profileForm.gender === 1 ? 'Male' : profileForm.gender === 2 ? 'Female' : 'QuestionFilled'" />
              </el-icon>
            </div>
            性别
          </div>
          <div class="info-item-value">{{ genderText }}</div>
        </div>

        <div class="info-item">
          <div class="info-item-label">
            <div class="info-icon-wrap info-icon-orange">
              <el-icon :size="14" style="color: var(--med-warning);"><Calendar /></el-icon>
            </div>
            生日
          </div>
          <div class="info-item-value">{{ profileForm.birthday || '未设置' }}</div>
        </div>

        <div class="info-item info-item-full">
          <div class="info-item-label">
            <div class="info-icon-wrap info-icon-purple">
              <el-icon :size="14" style="color: var(--med-ai-purple);"><Location /></el-icon>
            </div>
            地址
          </div>
          <div class="info-item-value address-val">{{ profileForm.address || '未设置' }}</div>
        </div>
      </div>
    </el-card>

    <!-- 编辑信息弹窗 -->
    <el-dialog
      v-model="editDialogVisible"
      title="修改个人信息"
      width="620px"
      :before-close="handleEditDialogClose"
      class="med-dialog"
      destroy-on-close
    >
      <el-form
        ref="editFormRef"
        :model="editForm"
        :rules="editRules"
        label-width="90px"
        class="edit-form"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="editForm.username" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="editForm.name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="editForm.email" placeholder="请输入邮箱" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="editForm.phone" placeholder="请输入手机号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="性别" prop="gender">
              <el-select v-model="editForm.gender" placeholder="请选择性别" clearable style="width: 100%">
                <el-option label="男" :value="1" />
                <el-option label="女" :value="2" />
                <el-option label="未知" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="生日" prop="birthday">
              <el-date-picker v-model="editForm.birthday" type="date" placeholder="请选择生日" format="YYYY-MM-DD" value-format="YYYY-MM-DD" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="地址" prop="address">
              <el-input v-model="editForm.address" type="textarea" :rows="2" placeholder="请输入地址" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="新密码" prop="newPassword">
              <el-input v-model="editForm.newPassword" type="password" placeholder="不修改请留空" show-password />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="确认密码" prop="confirmPassword" v-if="editForm.newPassword">
              <el-input v-model="editForm.confirmPassword" type="password" placeholder="请再次输入密码" show-password />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="updateProfile" :loading="loading">确认修改</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 修改头像弹窗 -->
    <el-dialog
      v-model="avatarDialogVisible"
      title="更换头像"
      width="460px"
      :before-close="handleAvatarDialogClose"
      class="med-dialog"
      destroy-on-close
    >
      <div class="avatar-edit">
        <el-avatar :src="profileForm.avatar || defaultAvatar" :size="120" shape="circle" class="avatar-preview-img" />
        <el-upload
          class="avatar-upload"
          action="/api/user/avatar"
          :show-file-list="false"
          :on-success="handleAvatarSuccess"
          :on-error="handleAvatarError"
          :before-upload="beforeAvatarUpload"
          :auto-upload="true"
          :headers="uploadHeaders"
          name="avatar"
        >
          <el-button type="primary" size="large">
            <el-icon><Upload /></el-icon> 选择图片
          </el-button>
          <p class="upload-tip">支持 JPG、PNG 格式，文件小于 2MB</p>
        </el-upload>
      </div>
      <template #footer>
        <el-button @click="avatarDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <el-alert v-if="errorMessage" :title="errorMessage" type="error" show-icon :closable="false" class="page-alert" />
    <el-alert v-if="successMessage" :title="successMessage" type="success" show-icon :closable="false" class="page-alert" />
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import request from '../utils/request'
import { User, Edit, Picture, Message, Phone, Male, Female, Calendar, Location, Clock, Upload, QuestionFilled, Close, Check } from '@element-plus/icons-vue'

export default {
  name: 'UserProfile',
  components: { User, Edit, Picture, Message, Phone, Male, Female, Calendar, Location, Clock, Upload, QuestionFilled, Close, Check },
  setup() {
    const editFormRef = ref(null)
    const profileForm = ref({ id: '', username: '', createTime: '', avatar: '', name: '', email: '', phone: '', gender: 0, birthday: '', address: '' })
    const editForm = ref({ id: '', username: '', name: '', email: '', phone: '', gender: 0, birthday: '', address: '', newPassword: '', confirmPassword: '' })
    const editDialogVisible = ref(false)
    const avatarDialogVisible = ref(false)
    const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
    const uploadHeaders = computed(() => { const t = localStorage.getItem('token'); return t ? { Authorization: t } : {} })
    const genderText = computed(() => ({ 1: '男', 2: '女' }[profileForm.value.gender] || '未知'))
    const formattedCreateTime = computed(() => {
      try {
        if (!profileForm.value.createTime) return '未知'
        const d = new Date(profileForm.value.createTime)
        if (isNaN(d.getTime())) return profileForm.value.createTime || '未知'
        return d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
      } catch { return '未知' }
    })
    const editRules = {
      email: [{ type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }],
      phone: [{ pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }],
      newPassword: [{ min: 6, max: 20, message: '密码长度应在6到20个字符之间', trigger: 'blur' }],
      confirmPassword: [{ validator: (rule, value, callback) => { if (value && value !== editForm.value.newPassword) callback(new Error('两次输入的密码不一致')); else callback() }, trigger: 'blur' }]
    }
    const loading = ref(false)
    const errorMessage = ref('')
    const successMessage = ref('')

    const openEditDialog = () => {
      editForm.value = { ...profileForm.value, newPassword: '', confirmPassword: '' }
      editDialogVisible.value = true
    }
    const openAvatarDialog = () => { avatarDialogVisible.value = true }
    const handleEditDialogClose = (done) => { done() }
    const handleAvatarDialogClose = (done) => { done() }

    const handleAvatarSuccess = (res) => {
      if (res.code === 1) { profileForm.value.avatar = res.data; successMessage.value = '头像修改成功'; getCurrentUser(); avatarDialogVisible.value = false }
      else errorMessage.value = res.msg || '头像修改失败'
    }
    const handleAvatarError = () => { errorMessage.value = '头像修改失败，请重试' }
    const beforeAvatarUpload = (file) => {
      if (!file.type.startsWith('image/')) { errorMessage.value = '头像必须是图片格式!'; return false }
      if (file.size / 1024 / 1024 >= 2) { errorMessage.value = '头像大小不能超过 2MB!'; return false }
      return true
    }

    const getCurrentUser = async () => {
      try {
        const res = await request.get('/user')
        if (res.code === 1) {
          const d = res.data
          profileForm.value = { id: d.id, username: d.username, createTime: d.createTime, avatar: d.avatar || '', name: d.name || '', email: d.email || '', phone: d.phone || '', gender: d.gender || 0, birthday: d.birthday || '', address: d.address || '' }
        } else errorMessage.value = res.msg || '获取用户信息失败'
      } catch (e) { errorMessage.value = e.message || '获取用户信息失败' }
    }

    const updateProfile = async () => {
      if (!editFormRef.value) return
      await editFormRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true; errorMessage.value = ''; successMessage.value = ''
          try {
            const updateData = { id: editForm.value.id, username: editForm.value.username, name: editForm.value.name, email: editForm.value.email, phone: editForm.value.phone, gender: editForm.value.gender, birthday: editForm.value.birthday, address: editForm.value.address }
            if (editForm.value.newPassword) updateData.password = editForm.value.newPassword
            const res = await request.put('/user', updateData)
            if (res.code === 1) { successMessage.value = '用户信息更新成功'; editDialogVisible.value = false; await getCurrentUser() }
            else errorMessage.value = res.msg || '更新失败'
          } catch (e) { errorMessage.value = e.message || '更新用户信息失败' }
          finally { loading.value = false }
        }
      })
    }

    onMounted(() => { getCurrentUser() })

    return { editFormRef, profileForm, editForm, formattedCreateTime, genderText, editRules, loading, errorMessage, successMessage, defaultAvatar, uploadHeaders, editDialogVisible, avatarDialogVisible, handleAvatarSuccess, handleAvatarError, beforeAvatarUpload, openEditDialog, openAvatarDialog, handleEditDialogClose, handleAvatarDialogClose, updateProfile }
  }
}
</script>

<style scoped>
.profile-page {
  padding: 24px;
  min-height: 100%;
  background: var(--med-bg);
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Banner */
.profile-banner {
  padding: 28px 32px;
  background: linear-gradient(135deg, var(--med-sidebar) 0%, var(--med-sidebar-mid) 55%, #0C4A56 100%);
  border-radius: var(--med-radius-xl);
  box-shadow: var(--med-shadow-lg);
  position: relative;
  overflow: hidden;
}

.banner-bg-circle { position: absolute; border-radius: 50%; pointer-events: none; }
.banner-circle-1 { width: 280px; height: 280px; top: -80px; right: -60px; background: radial-gradient(circle, rgba(50,224,196,0.1) 0%, transparent 70%); }
.banner-circle-2 { width: 180px; height: 180px; bottom: -50px; left: 200px; background: radial-gradient(circle, rgba(124,77,255,0.08) 0%, transparent 70%); }

.banner-content {
  position: relative; z-index: 1;
  display: flex; align-items: center; gap: 24px;
}

.banner-avatar-wrap { position: relative; flex-shrink: 0; }
.banner-avatar { border: 3px solid rgba(50, 224, 196, 0.5); box-shadow: 0 6px 20px rgba(0,0,0,0.3); }

.avatar-status-dot {
  position: absolute; bottom: 4px; right: 4px;
  width: 14px; height: 14px;
  background: var(--med-success);
  border-radius: 50%; border: 2px solid rgba(10,35,66,0.8);
}

.banner-info { flex: 1; }
.banner-name { font-size: 24px; font-weight: 800; color: #fff; margin: 0 0 4px 0; }
.banner-username { font-size: 14px; color: rgba(255,255,255,0.5); margin: 0 0 12px 0; }

.banner-tags { display: flex; gap: 8px; flex-wrap: wrap; }
.banner-tag {
  display: inline-flex; align-items: center; gap: 4px;
  background: rgba(255,255,255,0.08); border: 1px solid rgba(255,255,255,0.15);
  color: rgba(255,255,255,0.65); font-size: 12px; padding: 4px 10px; border-radius: 20px;
}
.banner-tag-teal { border-color: rgba(50,224,196,0.3); color: rgba(50,224,196,0.9); background: rgba(50,224,196,0.08); }

.banner-actions { display: flex; gap: 10px; flex-shrink: 0; }

.action-btn {
  background: linear-gradient(135deg, var(--med-primary) 0%, var(--med-primary-light) 100%) !important;
  border: none !important;
  padding: 10px 18px !important;
  border-radius: var(--med-radius-md) !important;
  font-weight: 600 !important;
  box-shadow: 0 4px 12px rgba(13,115,119,0.4) !important;
}

.action-btn-outline {
  background: rgba(255,255,255,0.08) !important;
  border: 1px solid rgba(255,255,255,0.25) !important;
  color: rgba(255,255,255,0.85) !important;
  padding: 10px 18px !important;
  border-radius: var(--med-radius-md) !important;
}

.action-btn-outline:hover {
  background: rgba(255,255,255,0.15) !important;
  border-color: rgba(255,255,255,0.4) !important;
}

/* 信息卡片 */
.info-card :deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid var(--med-border-light);
}

.section-card-header {
  display: flex; align-items: center; gap: 8px;
  font-size: 15px; font-weight: 700; color: var(--med-text);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 4px 0;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 18px;
  background: var(--med-bg);
  border-radius: var(--med-radius-md);
  border: 1px solid var(--med-border-light);
  transition: all 0.2s ease;
}

.info-item:hover {
  border-color: var(--med-border);
  background: #fff;
  box-shadow: var(--med-shadow-sm);
}

.info-item-full { grid-column: span 2; }

.info-item-label {
  display: flex; align-items: center; gap: 10px;
  font-size: 13px; font-weight: 600; color: var(--med-text-secondary);
  min-width: 90px; flex-shrink: 0;
}

.info-icon-wrap {
  width: 28px; height: 28px; border-radius: 7px;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}

.info-icon-blue { background: var(--med-info-light); }
.info-icon-green { background: var(--med-success-light); }
.info-icon-orange { background: var(--med-warning-light); }
.info-icon-purple { background: var(--med-ai-purple-light); }

.info-item-value { font-size: 14px; color: var(--med-text); font-weight: 500; }
.address-val { white-space: pre-wrap; line-height: 1.5; }

/* 弹窗 */
:deep(.med-dialog .el-dialog__header) {
  background: linear-gradient(135deg, var(--med-primary) 0%, var(--med-primary-light) 100%);
  padding: 18px 24px;
  margin: 0;
}

:deep(.med-dialog .el-dialog__title) { color: #fff; font-size: 16px; font-weight: 700; }
:deep(.med-dialog .el-dialog__headerbtn .el-dialog__close) { color: rgba(255,255,255,0.8); }
:deep(.med-dialog .el-dialog__headerbtn:hover .el-dialog__close) { color: #fff; }

.edit-form { padding: 20px 4px; }
:deep(.edit-form .el-form-item__label) { font-weight: 600; color: var(--med-text-secondary); font-size: 13px; }

.dialog-footer { display: flex; justify-content: flex-end; gap: 10px; }

/* 头像编辑 */
.avatar-edit { display: flex; flex-direction: column; align-items: center; padding: 30px 20px; gap: 20px; }

.avatar-preview-img { border: 3px solid var(--med-border); box-shadow: var(--med-shadow-md); }

.avatar-upload { display: flex; flex-direction: column; align-items: center; gap: 8px; }
.upload-tip { font-size: 12px; color: var(--med-text-muted); margin: 0; }

/* 提示弹窗 */
.page-alert {
  position: fixed; top: 20px; left: 50%;
  transform: translateX(-50%);
  width: calc(100% - 40px); max-width: 500px;
  border-radius: var(--med-radius-md) !important; z-index: 9999;
  box-shadow: var(--med-shadow-lg);
}

@media (max-width: 768px) {
  .profile-page { padding: 16px; gap: 14px; }
  .banner-content { flex-direction: column; text-align: center; }
  .banner-info { display: flex; flex-direction: column; align-items: center; }
  .banner-tags { justify-content: center; }
  .banner-actions { justify-content: center; }
  .info-grid { grid-template-columns: 1fr; }
  .info-item-full { grid-column: span 1; }
}
</style>
