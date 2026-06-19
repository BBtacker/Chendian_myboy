<template>
  <div class="chat-page">
    <!-- 左侧历史对话列表 -->
    <div class="conv-sidebar" v-show="showConversationList">
      <div class="conv-sidebar-header">
        <div class="conv-header-title">
          <el-icon style="color: var(--med-primary); font-size: 16px;"><Clock /></el-icon>
          <span>历史对话</span>
        </div>
        <el-button type="primary" size="small" plain @click="startNewConversation" class="new-conv-btn">
          <el-icon><Plus /></el-icon> 新对话
        </el-button>
      </div>

      <div class="conv-list">
        <div
          v-for="conversation in conversationList"
          :key="conversation.id"
          :class="['conv-item', { active: activeConversation && activeConversation.id === conversation.id }]"
          @click="loadConversation(conversation)"
        >
          <div class="conv-item-icon">
            <el-icon style="font-size: 14px;"><ChatLineRound /></el-icon>
          </div>
          <div class="conv-item-content">
            <div class="conv-item-title">{{ conversation.title }}</div>
            <div class="conv-item-time">{{ formatConversationTime(conversation.lastUpdateTime) }}</div>
          </div>
        </div>

        <div v-if="conversationList.length === 0" class="conv-empty">
          <el-icon style="font-size: 32px; color: var(--med-text-muted); margin-bottom: 10px;"><ChatLineRound /></el-icon>
          <div>暂无历史对话</div>
          <div style="font-size: 12px; margin-top: 4px;">点击右上方新建对话</div>
        </div>
      </div>
    </div>

    <!-- 主聊天区域 -->
    <div class="chat-main">
      <!-- 顶部栏 -->
      <div class="chat-topbar">
        <div class="chat-topbar-left">
          <div class="ai-status-dot"></div>
          <div>
            <div class="chat-topbar-title">腺体面容 AI 助手</div>
            <div class="chat-topbar-sub">基于大语言模型 · 医学专业问答</div>
          </div>
        </div>
        <el-button
          v-if="activeConversation && activeConversation.id"
          type="danger"
          plain
          size="small"
          @click="clearConversation"
          class="clear-btn"
        >
          <el-icon><Delete /></el-icon> 清除对话
        </el-button>
      </div>

      <!-- 消息区域 -->
      <div class="messages-area" ref="messagesContainer">
        <!-- 欢迎提示（无消息时） -->
        <div v-if="messageList.length === 0 && !loading" class="welcome-hint">
          <div class="welcome-icon">
            <el-icon :size="40" style="color: var(--med-ai-purple);"><Headset /></el-icon>
          </div>
          <div class="welcome-title">你好，我是 AI 医疗助手</div>
          <div class="welcome-desc">我可以解答关于腺样体面容的医学问题，为您提供专业建议</div>
          <div class="quick-questions">
            <div class="quick-q-item" @click="setQuickQuestion('什么是腺样体面容？')">什么是腺样体面容？</div>
            <div class="quick-q-item" @click="setQuickQuestion('腺样体面容如何治疗？')">腺样体面容如何治疗？</div>
            <div class="quick-q-item" @click="setQuickQuestion('腺样体面容的早期症状有哪些？')">腺样体面容的早期症状？</div>
          </div>
        </div>

        <!-- 消息列表 -->
        <div
          v-for="message in messageList"
          :key="message.id"
          :class="['message-row', message.sender === 0 ? 'user-row' : 'ai-row']"
        >
          <!-- AI消息 -->
          <template v-if="message.sender !== 0">
            <div class="msg-avatar ai-avatar">
              <el-icon :size="18" style="color: #fff;"><Headset /></el-icon>
            </div>
            <div class="msg-bubble ai-bubble">
              <div class="msg-text">{{ message.content }}</div>
              <div class="msg-time">{{ formatTime(message.sendTime) }}</div>
            </div>
          </template>

          <!-- 用户消息 -->
          <template v-else>
            <div class="msg-bubble user-bubble">
              <div class="msg-text">{{ message.content }}</div>
              <div class="msg-time user-time">{{ formatTime(message.sendTime) }}</div>
            </div>
            <div class="msg-avatar user-avatar">
              <el-icon :size="18" style="color: #fff;"><User /></el-icon>
            </div>
          </template>
        </div>

        <!-- 加载中 -->
        <div v-if="loading" class="message-row ai-row">
          <div class="msg-avatar ai-avatar">
            <el-icon :size="18" style="color: #fff;"><Headset /></el-icon>
          </div>
          <div class="msg-bubble ai-bubble typing-bubble">
            <span class="dot"></span>
            <span class="dot"></span>
            <span class="dot"></span>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="input-area">
        <div class="input-hint">
          <el-icon style="font-size: 12px; color: var(--med-text-muted);"><InfoFilled /></el-icon>
          <span>Enter 发送，Shift+Enter 换行</span>
        </div>
        <div class="input-row">
          <el-button
            :type="isListening ? 'danger' : 'default'"
            circle
            @click="toggleSpeechRecognition"
            :disabled="loading"
            class="voice-btn"
          >
            <el-icon><Microphone /></el-icon>
          </el-button>
          <el-input
            v-model="inputMessage"
            type="textarea"
            :rows="2"
            placeholder="请输入您的医学问题..."
            @keydown.enter.exact.prevent="sendMessage"
            :disabled="loading"
            class="chat-input"
            resize="none"
          />
          <el-button
            type="primary"
            :loading="loading"
            :disabled="!inputMessage.trim() || loading"
            @click="sendMessage"
            class="send-btn"
          >
            <el-icon v-if="!loading"><Position /></el-icon>
            <span>发送</span>
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, nextTick } from 'vue'
import request from '../utils/request'
import { Headset, Delete, Clock, Plus, ChatLineRound, Microphone, Position, User, InfoFilled } from '@element-plus/icons-vue'

export default {
  name: 'ChatAssistant',
  components: { Headset, Delete, Clock, Plus, ChatLineRound, Microphone, Position, User, InfoFilled },
  setup() {
    const conversationList = ref([])
    const activeConversation = ref(null)
    const messageList = ref([])
    const inputMessage = ref('')
    const loading = ref(false)
    const messagesContainer = ref(null)
    const showConversationList = ref(true)
    const isListening = ref(false)
    let recognition = null

    const initSpeechRecognition = () => {
      const SR = window.SpeechRecognition || window.webkitSpeechRecognition
      if (SR) {
        recognition = new SR()
        recognition.lang = 'zh-CN'
        recognition.continuous = false
        recognition.interimResults = false
        recognition.onstart = () => { isListening.value = true }
        recognition.onresult = (e) => { inputMessage.value += e.results[0][0].transcript }
        recognition.onerror = () => { isListening.value = false }
        recognition.onend = () => { isListening.value = false }
      }
    }

    const toggleSpeechRecognition = () => {
      if (!recognition) return
      if (isListening.value) { recognition.stop(); isListening.value = false }
      else { try { recognition.start() } catch (e) { console.error(e) } }
    }

    const formatTime = (time) => {
      if (!time) return ''
      const d = new Date(time)
      return `${String(d.getHours()).padStart(2,'0')}:${String(d.getMinutes()).padStart(2,'0')}`
    }

    const formatConversationTime = (time) => {
      if (!time) return ''
      const d = new Date(time)
      const now = new Date()
      if (d.toDateString() === now.toDateString()) return `${String(d.getHours()).padStart(2,'0')}:${String(d.getMinutes()).padStart(2,'0')}`
      const yesterday = new Date(now)
      yesterday.setDate(yesterday.getDate() - 1)
      if (d.toDateString() === yesterday.toDateString()) return '昨天'
      return `${d.getMonth()+1}-${d.getDate()}`
    }

    const scrollToBottom = () => {
      if (messagesContainer.value) messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }

    const loadConversationList = async () => {
      try {
        const res = await request.get('/conversation/list')
        if (res.code === 1) conversationList.value = res.data || []
      } catch (e) { console.error(e) }
    }

    const startNewConversation = () => {
      activeConversation.value = null
      messageList.value = []
      inputMessage.value = ''
      showConversationList.value = false
    }

    const loadConversation = async (conv) => {
      try {
        const res = await request.get(`/conversation/messages?conversationId=${conv.id}`)
        if (res.code === 1) {
          activeConversation.value = conv
          messageList.value = res.data || []
          showConversationList.value = false
          await nextTick(); scrollToBottom()
        }
      } catch (e) { console.error(e) }
    }

    const createConversation = async (firstMessage) => {
      try {
        const res = await request.post('/conversation/create', `firstMessage=${encodeURIComponent(firstMessage)}`, { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } })
        if (res.code === 1) { activeConversation.value = res.data; await loadConversationList(); return res.data.id }
        return null
      } catch (e) { console.error(e); return null }
    }

    const sendMessage = async () => {
      if (!inputMessage.value.trim() || loading.value) return
      const content = inputMessage.value
      if (!activeConversation.value) {
        const id = await createConversation(content)
        if (!id) { messageList.value.push({ id: Date.now(), sender: 1, content: '创建对话失败，请重试', sendTime: new Date() }); return }
      }
      await sendToAI(content)
    }

    const sendToAI = async (content) => {
      messageList.value.push({ id: Date.now(), sender: 0, content, sendTime: new Date() })
      inputMessage.value = ''
      loading.value = true
      await nextTick(); scrollToBottom()
      try {
        const res = await request.post('/conversation/send', `conversationId=${activeConversation.value.id}&content=${encodeURIComponent(content)}`, { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } })
        if (res.code === 1) {
          messageList.value.push({ id: Date.now()+1, sender: 1, content: res.data, sendTime: new Date() })
          await loadConversationList()
        } else {
          messageList.value.push({ id: Date.now()+1, sender: 1, content: res.msg || '发送失败，请重试', sendTime: new Date() })
        }
      } catch (e) {
        messageList.value.push({ id: Date.now()+1, sender: 1, content: '网络错误，请重试', sendTime: new Date() })
      } finally {
        loading.value = false
        await nextTick(); scrollToBottom()
      }
    }

    const clearConversation = async () => {
      if (activeConversation.value?.id) {
        try { await request.delete(`/conversation/${activeConversation.value.id}`) } catch (e) { console.error(e) }
      }
      messageList.value = []
      activeConversation.value = null
      inputMessage.value = ''
      showConversationList.value = true
    }

    const setQuickQuestion = (q) => { inputMessage.value = q; sendMessage() }

    onMounted(async () => { await loadConversationList(); initSpeechRecognition() })

    return { conversationList, activeConversation, messageList, inputMessage, loading, messagesContainer, showConversationList, isListening, formatTime, formatConversationTime, sendMessage, startNewConversation, loadConversation, clearConversation, toggleSpeechRecognition, setQuickQuestion }
  }
}
</script>

<style scoped>
.chat-page {
  display: flex;
  height: calc(100vh - 66px);
  background: var(--med-bg);
  overflow: hidden;
}

/* ===== 历史对话侧边栏 ===== */
.conv-sidebar {
  width: 260px;
  min-width: 260px;
  background: #fff;
  border-right: 1px solid var(--med-border-light);
  display: flex;
  flex-direction: column;
  box-shadow: var(--med-shadow-sm);
}

.conv-sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 16px;
  border-bottom: 1px solid var(--med-border-light);
}

.conv-header-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 700;
  color: var(--med-text);
}

.new-conv-btn {
  border-radius: var(--med-radius-sm) !important;
  border-color: var(--med-primary) !important;
  color: var(--med-primary) !important;
}

.conv-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.conv-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 12px;
  border-radius: var(--med-radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 4px;
  border: 1px solid transparent;
}

.conv-item:hover {
  background: var(--med-bg);
  border-color: var(--med-border);
}

.conv-item.active {
  background: var(--med-primary-bg);
  border-color: rgba(13, 115, 119, 0.25);
}

.conv-item-icon {
  width: 32px; height: 32px;
  background: var(--med-bg); border-radius: 8px;
  display: flex; align-items: center; justify-content: center;
  color: var(--med-primary); flex-shrink: 0;
}

.conv-item.active .conv-item-icon {
  background: var(--med-primary-bg);
}

.conv-item-title {
  font-size: 13px; font-weight: 500; color: var(--med-text);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}

.conv-item-time { font-size: 11px; color: var(--med-text-muted); margin-top: 2px; }

.conv-empty {
  display: flex; flex-direction: column; align-items: center;
  justify-content: center; padding: 60px 20px;
  color: var(--med-text-muted); font-size: 13px; text-align: center;
}

/* ===== 主聊天区域 ===== */
.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

.chat-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 24px;
  background: #fff;
  border-bottom: 1px solid var(--med-border-light);
  box-shadow: var(--med-shadow-sm);
}

.chat-topbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.ai-status-dot {
  width: 10px; height: 10px;
  background: var(--med-success);
  border-radius: 50%;
  box-shadow: 0 0 0 3px rgba(39, 174, 96, 0.2);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { box-shadow: 0 0 0 3px rgba(39, 174, 96, 0.2); }
  50% { box-shadow: 0 0 0 6px rgba(39, 174, 96, 0.08); }
}

.chat-topbar-title { font-size: 16px; font-weight: 700; color: var(--med-text); line-height: 1.2; }
.chat-topbar-sub { font-size: 12px; color: var(--med-text-muted); }

.clear-btn {
  border-radius: var(--med-radius-sm) !important;
  font-size: 13px !important;
}

/* ===== 消息区域 ===== */
.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: #F6F8FB;
}

/* 欢迎提示 */
.welcome-hint {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 40px;
  text-align: center;
  flex: 1;
}

.welcome-icon {
  width: 80px; height: 80px;
  background: var(--med-ai-purple-light);
  border: 1px solid rgba(124, 77, 255, 0.2);
  border-radius: 24px;
  display: flex; align-items: center; justify-content: center;
  margin-bottom: 20px;
}

.welcome-title { font-size: 20px; font-weight: 700; color: var(--med-text); margin-bottom: 8px; }
.welcome-desc { font-size: 14px; color: var(--med-text-muted); margin-bottom: 28px; line-height: 1.6; }

.quick-questions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: center;
  max-width: 500px;
}

.quick-q-item {
  padding: 8px 16px;
  background: #fff;
  border: 1px solid var(--med-border);
  border-radius: 20px;
  font-size: 13px;
  color: var(--med-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.quick-q-item:hover {
  border-color: var(--med-primary);
  color: var(--med-primary);
  background: var(--med-primary-bg);
}

/* 消息行 */
.message-row {
  display: flex;
  align-items: flex-end;
  gap: 10px;
  max-width: 75%;
}

.ai-row {
  align-self: flex-start;
}

.user-row {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.msg-avatar {
  width: 36px; height: 36px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}

.ai-avatar {
  background: linear-gradient(135deg, var(--med-ai-purple) 0%, #9C6FFF 100%);
  box-shadow: 0 3px 10px rgba(124, 77, 255, 0.35);
}

.user-avatar {
  background: linear-gradient(135deg, var(--med-primary) 0%, var(--med-primary-light) 100%);
  box-shadow: 0 3px 10px rgba(13, 115, 119, 0.35);
}

.msg-bubble {
  padding: 12px 16px;
  border-radius: 18px;
  max-width: 100%;
  word-break: break-word;
}

.ai-bubble {
  background: #fff;
  box-shadow: var(--med-shadow-sm);
  border: 1px solid var(--med-border-light);
  border-bottom-left-radius: 6px;
}

.user-bubble {
  background: linear-gradient(135deg, var(--med-primary) 0%, var(--med-primary-light) 100%);
  border-bottom-right-radius: 6px;
  box-shadow: 0 4px 16px rgba(13, 115, 119, 0.25);
}

.msg-text {
  font-size: 14px;
  line-height: 1.7;
  white-space: pre-wrap;
  color: var(--med-text);
}

.user-bubble .msg-text { color: #fff; }

.msg-time {
  font-size: 11px;
  color: var(--med-text-muted);
  margin-top: 6px;
  text-align: right;
}

.user-time { color: rgba(255,255,255,0.65); }

/* 打字动画 */
.typing-bubble {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 14px 18px;
}

.dot {
  width: 8px; height: 8px;
  background: var(--med-primary-light);
  border-radius: 50%;
  animation: bounce 1.4s infinite ease-in-out;
}

.dot:nth-child(2) { animation-delay: 0.2s; }
.dot:nth-child(3) { animation-delay: 0.4s; }

@keyframes bounce {
  0%, 60%, 100% { transform: translateY(0); }
  30% { transform: translateY(-6px); }
}

/* ===== 输入区域 ===== */
.input-area {
  padding: 16px 24px 20px;
  background: #fff;
  border-top: 1px solid var(--med-border-light);
}

.input-hint {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: var(--med-text-muted);
  margin-bottom: 10px;
}

.input-row {
  display: flex;
  align-items: flex-end;
  gap: 10px;
}

.voice-btn {
  width: 44px !important;
  height: 44px !important;
  border-radius: 50% !important;
  flex-shrink: 0;
  border-color: var(--med-border) !important;
  color: var(--med-text-secondary) !important;
  margin-bottom: 4px;
}

:deep(.chat-input .el-textarea__inner) {
  border-radius: var(--med-radius-md) !important;
  border: 1.5px solid var(--med-border) !important;
  resize: none;
  font-size: 14px;
  padding: 12px 16px;
  line-height: 1.6;
  transition: all 0.25s ease;
}

:deep(.chat-input .el-textarea__inner:focus) {
  border-color: var(--med-primary) !important;
  box-shadow: 0 0 0 3px var(--med-primary-bg) !important;
}

.send-btn {
  height: 52px !important;
  padding: 0 20px !important;
  border-radius: var(--med-radius-md) !important;
  background: linear-gradient(135deg, var(--med-primary) 0%, var(--med-primary-light) 100%) !important;
  border: none !important;
  font-size: 14px !important;
  font-weight: 600 !important;
  display: flex; align-items: center; gap: 6px;
  box-shadow: 0 4px 12px rgba(13, 115, 119, 0.3) !important;
  flex-shrink: 0;
  margin-bottom: 4px;
}

.send-btn:hover:not(:disabled) {
  transform: translateY(-1px) !important;
  box-shadow: 0 6px 18px rgba(13, 115, 119, 0.4) !important;
}

@media (max-width: 768px) {
  .conv-sidebar { width: 200px; min-width: 200px; }
  .message-row { max-width: 90%; }
}
</style>
