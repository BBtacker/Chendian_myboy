import React, { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Button, Card, Typography, Modal } from 'antd';
import { getDoubaoAnswer, chatHistoryManager } from '../utils/doubaoApi';
import { chatApi } from '../utils/api';
import '../styles/AskPage.css';

const { Title, Text, Paragraph } = Typography;

/* ── Avatar components ───────────────────────────────── */
function AIAvatar() {
  return <div className="msg-avatar msg-avatar--ai">🤖</div>;
}
function UserAvatar() {
  return <div className="msg-avatar msg-avatar--user">👤</div>;
}

/* ── Single message row ──────────────────────────────── */
function Message({ message }) {
  const isUser = message.isUser;
  return (
    <motion.div
      className={`message-row ${isUser ? 'message-row--user' : 'message-row--ai'}`}
      initial={{ opacity: 0, y: 10, scale: 0.97 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{ duration: 0.28, ease: 'easeOut' }}
    >
      {!isUser && <AIAvatar />}
      <div className="message-bubble">
        <div className="message-text">{message.content}</div>
        <div className="message-time">{message.timestamp}</div>
      </div>
      {isUser && <UserAvatar />}
    </motion.div>
  );
}

/* ── Typing indicator ────────────────────────────────── */
function TypingIndicator() {
  return (
    <motion.div
      className="message-row message-row--ai"
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0 }}
    >
      <AIAvatar />
      <div className="message-bubble message-bubble--typing">
        <div className="typing-dots">
          {[0, 1, 2].map(i => (
            <motion.span
              key={i}
              className="typing-dot"
              animate={{ y: [0, -7, 0] }}
              transition={{ duration: 0.75, repeat: Infinity, delay: i * 0.15 }}
            />
          ))}
        </div>
      </div>
    </motion.div>
  );
}

/* ── Quick prompts ───────────────────────────────────── */
const QUICK_PROMPTS = [
  '今天心情不太好，帮我开导一下',
  '我想写一篇日记，给我一些建议',
  '怎么保持积极乐观的心态？',
];

/* ── Main Page ───────────────────────────────────────── */
function AskPage() {
  const WELCOME_MESSAGE = {
    id: 1,
    content: '你好！我是心影助手，专注于日记、心情和生活记录。有什么可以帮你的吗？',
    isUser: false,
    timestamp: new Date().toLocaleTimeString(),
  };

  const [messages, setMessages] = useState([WELCOME_MESSAGE]);
  const [inputValue, setInputValue] = useState('');
  const messagesEndRef = useRef(null);
  const [typingMessage, setTypingMessage] = useState(null);
  const [isTyping, setIsTyping] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isLoadingHistory, setIsLoadingHistory] = useState(true);
  const [currentSessionId, setCurrentSessionId] = useState(null);
  const timeoutRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    setIsLoadingHistory(true);
    try {
      const history = chatHistoryManager.getChatHistory();
      if (history && history.length > 0) setMessages(history);

      // 创建后端会话
      const currentUser = localStorage.getItem('currentUser')
      if (currentUser) {
        const user = JSON.parse(currentUser)
        chatApi.createSession(user.id, '心影助手对话').then(session => {
          if (session && session.session) {
            setCurrentSessionId(session.session.id)
          }
        }).catch(() => {}) // 静默失败，不影响前端使用
      }
    } catch (e) {
      console.error('加载聊天历史记录失败:', e);
    } finally {
      setIsLoadingHistory(false);
    }
  }, []);

  useEffect(() => {
    if (!isLoadingHistory) chatHistoryManager.saveChatHistory(messages);
  }, [messages, isLoadingHistory]);

  useEffect(() => { scrollToBottom(); }, [messages, typingMessage]);

  useEffect(() => () => {
    if (timeoutRef.current) clearTimeout(timeoutRef.current);
  }, []);

  const clearChatHistory = () => {
    Modal.confirm({
      title: '清除聊天历史',
      content: '确定要清除所有聊天历史记录吗？此操作不可恢复。',
      okText: '确定',
      okType: 'danger',
      cancelText: '取消',
      onOk() {
        setMessages([WELCOME_MESSAGE]);
        chatHistoryManager.clearChatHistory();
        if (currentSessionId) {
          chatApi.deleteSession(currentSessionId).catch(() => {})
        }
      }
    });
  };

  const getAIAnswer = async () => {
    try {
      setIsLoading(true);
      return await getDoubaoAnswer(messages);
    } catch {
      return '抱歉，我暂时无法回答您的问题，请稍后再试。';
    } finally {
      setIsLoading(false);
    }
  };

  const typeWriterEffect = (message, fullText, index = 0, speed = 28) => {
    if (index < fullText.length) {
      setTypingMessage({ ...message, content: fullText.substring(0, index + 1) });
      if (timeoutRef.current) clearTimeout(timeoutRef.current);
      timeoutRef.current = setTimeout(() => typeWriterEffect(message, fullText, index + 1, speed), speed);
    } else {
      const finalMessage = { ...message, content: fullText };
      setMessages(prev => [...prev, finalMessage]);
      setTypingMessage(null);
      setIsTyping(false);
      timeoutRef.current = null;

      // 保存AI回复到后端
      if (currentSessionId) {
        chatApi.saveMessage({ sessionId: currentSessionId, sender: 'assistant', content: fullText }).catch(() => {})
      }
    }
  };

  const sendMessage = async (text) => {
    const msg = (text || inputValue).trim();
    if (!msg || isLoading) return;
    const userMessage = { id: Date.now(), content: msg, isUser: true, timestamp: new Date().toLocaleTimeString() };
    const updatedMessages = [...messages, userMessage];
    setMessages(updatedMessages);
    setInputValue('');
    setIsTyping(true);

    // 保存用户消息到后端
    if (currentSessionId) {
      chatApi.saveMessage({ sessionId: currentSessionId, sender: 'user', content: msg }).catch(() => {})
    }

    const answer = await getDoubaoAnswer(updatedMessages);
    typeWriterEffect({ id: Date.now() + 1, content: '', isUser: false, timestamp: new Date().toLocaleTimeString() }, answer);
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey && !isLoading) { e.preventDefault(); sendMessage(); }
  };

  const canSend = !!inputValue.trim() && !isLoading && !isLoadingHistory;

  return (
    <div className="ask-page">

      {/* ── Header Banner ── */}
      <motion.div
        className="ask-header"
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <div className="ask-header__orbs">
          <div className="ask-orb ask-orb--1" />
          <div className="ask-orb ask-orb--2" />
        </div>
        <div className="ask-header__body">
          <div className="ask-header__left">
            <motion.div
              className="ask-header__icon"
              animate={{ rotate: [0, -10, 10, 0] }}
              transition={{ duration: 3.5, repeat: Infinity }}
            >🤖</motion.div>
            <div>
              <Title level={2} className="ask-header__title">智能心影助手</Title>
              <Paragraph className="ask-header__subtitle">
                与 AI 聊天，分享心情，获得温暖回应
              </Paragraph>
            </div>
          </div>
          <Button
            className="btn-clear"
            onClick={clearChatHistory}
            disabled={isLoading || isLoadingHistory}
          >
            🗑 清除历史
          </Button>
        </div>
      </motion.div>

      {/* ── Chat Card ── */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.18 }}
      >
        <Card className="chat-card" variant="outlined">

          {/* Messages */}
          <div className="chat-messages">
            {isLoadingHistory ? (
              <div className="chat-loading">
                <div className="chat-loading__spinner" />
                <Text type="secondary">加载聊天历史...</Text>
              </div>
            ) : (
              <>
                {messages.map(msg => <Message key={msg.id} message={msg} />)}
                <AnimatePresence>
                  {isTyping && !typingMessage && <TypingIndicator key="typing" />}
                </AnimatePresence>
                {typingMessage && <Message key={`t-${typingMessage.id}`} message={typingMessage} />}
              </>
            )}
            <div ref={messagesEndRef} />
          </div>

          {/* Quick prompts */}
          {messages.length <= 1 && !isLoadingHistory && (
            <div className="quick-prompts">
              {QUICK_PROMPTS.map((p, i) => (
                <motion.button
                  key={i}
                  className="quick-prompt-btn"
                  onClick={() => sendMessage(p)}
                  whileHover={{ scale: 1.02, y: -1 }}
                  whileTap={{ scale: 0.97 }}
                  initial={{ opacity: 0, y: 8 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.3 + i * 0.08 }}
                >
                  💬 {p}
                </motion.button>
              ))}
            </div>
          )}

          {/* Input Area */}
          <div className="chat-input-area">
            <div className="chat-input-wrapper">
              <textarea
                className="chat-textarea"
                placeholder="请输入您的问题... (Enter 发送，Shift+Enter 换行)"
                value={inputValue}
                onChange={(e) => setInputValue(e.target.value)}
                onKeyPress={handleKeyPress}
                rows={2}
                disabled={isLoading || isLoadingHistory}
              />
              <Button
                className="btn-send"
                type="primary"
                onClick={() => sendMessage()}
                disabled={!canSend}
                loading={isLoading}
              >
                {!isLoading && '➤'}
              </Button>
            </div>
            <Text className="input-hint">Enter 发送 · Shift+Enter 换行</Text>
          </div>

        </Card>
      </motion.div>
    </div>
  );
}

export default AskPage;
