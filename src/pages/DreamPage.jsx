import React, { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Button, Row, Col, Flex, Typography, Alert, Spin, Image, Card, List,
} from 'antd';
import '../styles/DreamPage.css';
import { generateImage } from '../utils/doubaoApi';
import { dreamApi } from '../utils/api';

const { Title, Text, Paragraph } = Typography;

/* ── Mood options ─────────────────────────────────────────── */
const moodOptions = [
  { value: 'happy',   emoji: '😊', label: '开心', color: '#f59e0b', gradient: 'linear-gradient(135deg,#f59e0b,#fbbf24)' },
  { value: 'sad',     emoji: '😢', label: '悲伤', color: '#3b82f6', gradient: 'linear-gradient(135deg,#3b82f6,#60a5fa)' },
  { value: 'angry',   emoji: '😠', label: '愤怒', color: '#ef4444', gradient: 'linear-gradient(135deg,#ef4444,#f87171)' },
  { value: 'anxious', emoji: '😰', label: '焦虑', color: '#10b981', gradient: 'linear-gradient(135deg,#10b981,#34d399)' },
  { value: 'tired',   emoji: '😴', label: '疲惫', color: '#8b5cf6', gradient: 'linear-gradient(135deg,#8b5cf6,#a78bfa)' },
];

const imageMapping = {
  happy:   'https://picsum.photos/seed/anime-happy1/600/450',
  sad:     'https://picsum.photos/seed/anime-happy2/600/450',
  angry:   'https://picsum.photos/seed/anime-happy3/600/450',
  anxious: 'https://picsum.photos/seed/anime-happy4/600/450',
  tired:   'https://picsum.photos/seed/anime-happy5/600/450',
  default: 'https://picsum.photos/seed/anime-happy6/600/450',
};

/* ── Sub-components ───────────────────────────────────────── */

function MoodCard({ option, selected, onSelect }) {
  return (
    <motion.button
      className={`mood-card${selected ? ' mood-card--selected' : ''}`}
      style={{ '--mood-color': option.color, '--mood-gradient': option.gradient }}
      onClick={() => onSelect(option.value)}
      whileHover={{ scale: 1.06, y: -2 }}
      whileTap={{ scale: 0.95 }}
    >
      <span className="mood-card__emoji">{option.emoji}</span>
      <span className="mood-card__label">{option.label}</span>
    </motion.button>
  );
}

function EmptyResult() {
  const sparkles = ['✨', '🌟', '💫', '🌸', '⭐'];
  return (
    <motion.div
      className="result-empty"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
    >
      <div className="result-empty__sparkles">
        {sparkles.map((s, i) => (
          <motion.span
            key={i}
            animate={{ y: [0, -10, 0], opacity: [0.4, 1, 0.4] }}
            transition={{ duration: 2 + i * 0.3, repeat: Infinity, delay: i * 0.25 }}
          >
            {s}
          </motion.span>
        ))}
      </div>
      <div className="result-empty__icon">🎨</div>
      <Text className="result-empty__text">
        上传照片并选择心情<br />AI将为你绘制专属动漫形象
      </Text>
    </motion.div>
  );
}

function ProcessingState() {
  return (
    <motion.div
      className="result-processing"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
    >
      <Spin size="large" />
      <div className="processing-dots">
        {[0, 1, 2].map(i => (
          <motion.div
            key={i}
            className="processing-dot"
            animate={{ scale: [1, 1.6, 1], opacity: [0.3, 1, 0.3] }}
            transition={{ duration: 1.2, repeat: Infinity, delay: i * 0.3 }}
          />
        ))}
      </div>
      <Text className="processing-text">AI正在绘制你的动漫心情图...</Text>
      <Text type="secondary" style={{ fontSize: 13 }}>请稍候，施法中 🔮</Text>
    </motion.div>
  );
}

/* ── Main Page ────────────────────────────────────────────── */
function DreamPage() {
  const [selectedImage, setSelectedImage] = useState(null);
  const [processedImage, setProcessedImage] = useState(null);
  const [isProcessing, setIsProcessing] = useState(false);
  const [mood, setMood] = useState('');
  const [error, setError] = useState('');
  const [dreamHistory, setDreamHistory] = useState([]);
  const fileInputRef = useRef(null);

  // 加载绘梦历史记录
  useEffect(() => {
    const loadDreamHistory = () => {
      try {
        const history = localStorage.getItem('dreamHistory');
        if (history) {
          setDreamHistory(JSON.parse(history));
        }
      } catch (e) {
        console.error('加载绘梦历史记录失败:', e);
      }
    };
    loadDreamHistory();
  }, []);

  // 保存绘梦历史记录
  useEffect(() => {
    try {
      localStorage.setItem('dreamHistory', JSON.stringify(dreamHistory));
    } catch (e) {
      console.error('保存绘梦历史记录失败:', e);
    }
  }, [dreamHistory]);

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = (ev) => {
      setSelectedImage(ev.target.result);
      setProcessedImage(null);
      setError('');
    };
    reader.readAsDataURL(file);
  };

  const handleProcess = async () => {
    if (!selectedImage) { setError('请先选择一张照片'); return; }
    if (!mood) { setError('请选择当前心情'); return; }
    setIsProcessing(true);
    setError('');
    try {
      // 提取base64数据，去掉前缀
      const base64Data = selectedImage.split(',')[1];
      // 调用生成图片API
      const generatedImage = await generateImage(base64Data, mood);
      // 检查返回的数据类型
      let imageUrl;
      if (generatedImage && generatedImage.startsWith('http')) {
        // 如果是URL，直接使用
        console.log('使用URL类型的图片:', generatedImage);
        imageUrl = generatedImage;
      } else {
        // 如果是base64，构建完整的图片URL
        console.log('使用base64类型的图片');
        imageUrl = `data:image/png;base64,${generatedImage}`;
      }
      setProcessedImage(imageUrl);
      
      // 添加到历史记录
      const newDream = {
        id: Date.now(),
        imageUrl,
        mood,
        timestamp: new Date().toISOString(),
        originalImage: selectedImage
      };
      setDreamHistory(prev => [newDream, ...prev].slice(0, 50)); // 保留最近50条记录

      // 保存绘梦记录到后端
      try {
        const currentUser = localStorage.getItem('currentUser')
        if (currentUser) {
          const user = JSON.parse(currentUser)
          await dreamApi.saveDreamImage({
            userId: user.id,
            prompt: mood,
            mood,
            imageUrl,
            originalImage: selectedImage
          })
        }
      } catch (e) {
        console.warn('保存绘梦记录到后端失败:', e)
      }
    } catch (error) {
      console.error('生成图片失败:', error);
      setError(`生成图片失败: ${error.message}`);
    } finally {
      setIsProcessing(false);
    }
  };

  const handleReset = () => {
    setSelectedImage(null);
    setProcessedImage(null);
    setMood('');
    setError('');
    if (fileInputRef.current) fileInputRef.current.value = '';
  };

  const handleDownload = () => {
    if (!processedImage) return;
    
    // 创建临时a标签
    const link = document.createElement('a');
    link.href = processedImage;
    
    // 设置文件名
    const timestamp = new Date().toISOString().slice(0, 19).replace(/[-:]/g, '');
    link.download = `心影绘梦_${mood}_${timestamp}.png`;
    
    // 触发下载
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const canProcess = !!selectedImage && !!mood && !isProcessing;

  return (
    <div className="dream-page">

      {/* ── Header Banner ── */}
      <motion.div
        className="dream-header"
        initial={{ opacity: 0, y: -24 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
      >
        <div className="dream-header__orbs">
          <div className="dream-orb dream-orb--1" />
          <div className="dream-orb dream-orb--2" />
          <div className="dream-orb dream-orb--3" />
        </div>
        <div className="dream-header__content">
          <motion.div
            className="dream-header__badge"
            animate={{ rotate: [0, 8, -8, 0] }}
            transition={{ duration: 4, repeat: Infinity }}
          >
            ✨
          </motion.div>
          <Title level={2} className="dream-header__title">心影绘梦</Title>
          <Paragraph className="dream-header__subtitle">
            上传一张照片，让 AI 为你创作专属的动漫风格心情图
          </Paragraph>
        </div>
      </motion.div>

      {/* ── Two-column Content ── */}
      <Row gutter={[24, 24]} className="dream-content" align="stretch">

        {/* Left Panel */}
        <Col xs={24} lg={12}>
          <motion.div
            style={{ height: '100%' }}
            initial={{ opacity: 0, x: -24 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
          >
            <Card className="dream-card" variant="outlined">

              {/* Upload Section */}
              <div className="dream-section">
                <Flex align="center" gap={8} className="dream-section__header">
                  <span className="section-icon">📸</span>
                  <Title level={4} className="section-title">上传照片</Title>
                </Flex>

                <div
                  className="upload-zone"
                  onClick={() => fileInputRef.current?.click()}
                >
                  <input
                    type="file"
                    accept="image/*"
                    ref={fileInputRef}
                    onChange={handleImageChange}
                    style={{ display: 'none' }}
                  />
                  <AnimatePresence mode="wait">
                    {selectedImage ? (
                      <motion.div
                        key="preview"
                        className="upload-zone__preview"
                        initial={{ opacity: 0, scale: 0.92 }}
                        animate={{ opacity: 1, scale: 1 }}
                        exit={{ opacity: 0, scale: 0.92 }}
                        transition={{ duration: 0.3 }}
                      >
                        <img src={selectedImage} alt="上传预览" className="preview-image" />
                        <div className="preview-overlay">
                          <span>🔄 点击更换照片</span>
                        </div>
                      </motion.div>
                    ) : (
                      <motion.div
                        key="placeholder"
                        className="upload-zone__placeholder"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                      >
                        <motion.div
                          className="upload-icon"
                          animate={{ y: [0, -8, 0] }}
                          transition={{ duration: 2.2, repeat: Infinity, ease: 'easeInOut' }}
                        >
                          ☁️
                        </motion.div>
                        <Text strong className="upload-zone__title">点击上传照片</Text>
                        <Text type="secondary" className="upload-zone__hint">
                          支持 JPG、PNG 格式 · 不超过 10MB
                        </Text>
                      </motion.div>
                    )}
                  </AnimatePresence>
                </div>
              </div>

              {/* Mood Section */}
              <div className="dream-section">
                <Flex align="center" gap={8} className="dream-section__header">
                  <span className="section-icon">🌈</span>
                  <Title level={4} className="section-title">当前心情</Title>
                </Flex>
                <div className="mood-grid">
                  {moodOptions.map(option => (
                    <MoodCard
                      key={option.value}
                      option={option}
                      selected={mood === option.value}
                      onSelect={setMood}
                    />
                  ))}
                </div>
              </div>

              {/* Error Alert */}
              <AnimatePresence>
                {error && (
                  <motion.div
                    initial={{ opacity: 0, height: 0, marginBottom: 0 }}
                    animate={{ opacity: 1, height: 'auto', marginBottom: 16 }}
                    exit={{ opacity: 0, height: 0, marginBottom: 0 }}
                    style={{ overflow: 'hidden' }}
                  >
                    <Alert
                      message={error}
                      type="error"
                      showIcon
                      closable
                      onClose={() => setError('')}
                    />
                  </motion.div>
                )}
              </AnimatePresence>

              {/* Action Buttons */}
              <Flex gap={12}>
                <Button
                  className="btn-generate"
                  type="primary"
                  size="large"
                  block
                  disabled={!canProcess}
                  onClick={handleProcess}
                  loading={isProcessing}
                >
                  {isProcessing ? '绘制中...' : '✨ 生成绘梦'}
                </Button>
                <Button
                  className="btn-reset"
                  size="large"
                  onClick={handleReset}
                  disabled={isProcessing}
                >
                  重置
                </Button>
              </Flex>

            </Card>
          </motion.div>
        </Col>

        {/* Right Panel — Result */}
        <Col xs={24} lg={12}>
          <motion.div
            style={{ height: '100%' }}
            initial={{ opacity: 0, x: 24 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5, delay: 0.3 }}
          >
            <Card className="dream-card result-card" variant="outlined">
              <Flex align="center" gap={8} className="dream-section__header">
                <span className="section-icon">🎨</span>
                <Title level={4} className="section-title">绘梦结果</Title>
              </Flex>

              <div className="result-display">
                <AnimatePresence mode="wait">
                  {isProcessing ? (
                    <ProcessingState key="processing" />
                  ) : processedImage ? (
                    <motion.div
                      key="result"
                      className="result-image-wrap"
                      initial={{ opacity: 0, scale: 0.82 }}
                      animate={{ opacity: 1, scale: 1 }}
                      exit={{ opacity: 0, scale: 0.9 }}
                      transition={{ duration: 0.55, type: 'spring', bounce: 0.32 }}
                    >
                      <div className="result-image-frame">
                        <Image
                          src={processedImage}
                          alt="动漫绘梦结果"
                          style={{ width: '100%', display: 'block', maxHeight: 320, objectFit: 'cover' }}
                          preview={{ mask: '预览大图' }}
                        />
                      </div>

                      <motion.div
                        className="result-success-badge"
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.45 }}
                      >
                        <span>🎉</span>
                        <Text className="success-text">你的动漫心情图已生成！</Text>
                        <span>🎉</span>
                      </motion.div>

                      <Button className="btn-download" size="small" onClick={handleDownload}>
                        💾 保存图片
                      </Button>
                    </motion.div>
                  ) : (
                    <EmptyResult key="empty" />
                  )}
                </AnimatePresence>
              </div>
            </Card>
          </motion.div>
        </Col>

      </Row>

      {/* ── 绘梦历史记录 ── */}
      <motion.div
        className="dream-history"
        initial={{ opacity: 0, y: 24 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.4 }}
      >
        <Card className="dream-card" variant="outlined">
          <Flex align="center" gap={8} className="dream-section__header">
            <span className="section-icon">📚</span>
            <Title level={4} className="section-title">绘梦历史</Title>
          </Flex>
          
          <div className="history-list-container">
            {dreamHistory.length === 0 ? (
              <div className="empty-history">
                <div className="empty-icon">🎨</div>
                <Text>还没有绘梦记录</Text>
                <Text type="secondary">生成绘梦后会自动保存到历史记录中</Text>
              </div>
            ) : (
              <List
                dataSource={dreamHistory}
                renderItem={item => (
                  <List.Item key={item.id} className="history-item">
                    <List.Item.Meta
                      avatar={
                        <Image
                          src={item.imageUrl}
                          width={80}
                          height={60}
                          style={{ objectFit: 'cover', borderRadius: 4 }}
                        />
                      }
                      title={
                        <div className="history-item-title">
                          <span className="mood-tag">
                            {moodOptions.find(m => m.value === item.mood)?.emoji} {moodOptions.find(m => m.value === item.mood)?.label}
                          </span>
                          <span className="history-time">
                            {new Date(item.timestamp).toLocaleString()}
                          </span>
                        </div>
                      }
                      description={
                        <div className="history-actions">
                          <Button 
                            size="small" 
                            onClick={() => setProcessedImage(item.imageUrl)}
                          >
                            查看
                          </Button>
                          <Button 
                            size="small" 
                            onClick={() => {
                              const link = document.createElement('a');
                              link.href = item.imageUrl;
                              link.download = `心影绘梦_${item.mood}_${new Date(item.timestamp).toISOString().slice(0, 19).replace(/[-:]/g, '')}.png`;
                              document.body.appendChild(link);
                              link.click();
                              document.body.removeChild(link);
                            }}
                          >
                            下载
                          </Button>
                        </div>
                      }
                    />
                  </List.Item>
                )}
              />
            )}
          </div>
        </Card>
      </motion.div>
    </div>
  )
}

export default DreamPage;
