import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import '../styles/CreateDiaryPage.css'
import { uploadImageToOSS, deleteImageFromOSS } from '../utils/doubaoApi'
import { message } from 'antd'
import { diaryApi } from '../utils/api'

function CreateDiaryPage() {
  const navigate = useNavigate()
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [selectedMood, setSelectedMood] = useState('')
  const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split('T')[0])
  const [uploadedPhotos, setUploadedPhotos] = useState([])
  const [isSubmitting, setIsSubmitting] = useState(false)

  // 心情选项
  const moodOptions = [
    { id: 1, name: '开心', icon: '😊', color: '#f1c40f' },
    { id: 2, name: '悲伤', icon: '😢', color: '#3498db' },
    { id: 3, name: '愤怒', icon: '😠', color: '#e74c3c' },
    { id: 4, name: '惊讶', icon: '😲', color: '#9b59b6' },
    { id: 5, name: '自然', icon: '😐', color: '#2ecc71' }
  ]

  // 处理照片上传
  const handlePhotoUpload = async (e) => {
    const files = Array.from(e.target.files)
    
    for (const file of files) {
      try {
        // 先上传到阿里云OSS
        const imageUrl = await uploadImageToOSS(file)
        
        const photoData = {
          id: Date.now() + Math.random(),
          name: file.name,
          type: file.type,
          url: imageUrl // 保存OSS返回的URL
        }
        setUploadedPhotos(prev => [...prev, photoData])
      } catch (error) {
        console.error('上传图片失败:', error)
        message.error(`上传图片失败: ${error.message}`)
      }
    }
  }

  // 删除照片
  const handleDeletePhoto = async (photoId) => {
    const photoToDelete = uploadedPhotos.find(photo => photo.id === photoId);
    
    // 如果是OSS图片，先从OSS删除
    if (photoToDelete && photoToDelete.url) {
      try {
        await deleteImageFromOSS(photoToDelete.url);
      } catch (error) {
        console.error('从OSS删除图片失败:', error);
      }
    }
    
    // 从本地状态删除
    setUploadedPhotos(uploadedPhotos.filter(photo => photo.id !== photoId));
  }

  // 处理表单提交
  const handleSubmit = (e) => {
    e.preventDefault()
    setIsSubmitting(true)

    // 获取当前用户信息
    const currentUser = JSON.parse(localStorage.getItem('currentUser'))
    if (!currentUser) {
      setIsSubmitting(false)
      message.error('用户未登录，请重新登录')
      navigate('/login')
      return
    }

    // 准备日记数据
    const diaryData = {
      userId: currentUser.id,
      title,
      content,
      diaryDate: selectedDate,
      selectedMood
    }

    // 保存日记到后端API
    setTimeout(async () => {
      try {
        // 创建日记，获取日记ID
        const createdDiary = await diaryApi.createDiary(diaryData)
        
        // 保存图片到diary_photo表
        if (uploadedPhotos.length > 0 && createdDiary.id) {
          for (const photo of uploadedPhotos) {
            await diaryApi.addDiaryPhoto({
              diaryId: createdDiary.id,
              photoUrl: photo.url
            })
          }
        }
        
        setIsSubmitting(false)
        message.success('日记保存成功！')
        // 导航到日记列表页面
        navigate('/diary')
      } catch (error) {
        setIsSubmitting(false)
        console.error('保存日记失败:', error)
        message.error(`保存日记失败: ${error.message}`)
      }
    }, 1000)
  }

  return (
    <div className="create-diary-page">
      <h2>写日记</h2>
      <form onSubmit={handleSubmit} className="diary-form">
        <div className="form-section">
          <h3>基本信息</h3>
          <div className="form-row">
            <div className="form-group">
              <label htmlFor="title">日记标题</label>
              <input
                type="text"
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="请输入日记标题"
                required
              />
            </div>
            <div className="form-group">
              <label htmlFor="date">日期</label>
              <input
                type="date"
                id="date"
                value={selectedDate}
                onChange={(e) => setSelectedDate(e.target.value)}
                required
              />
            </div>
          </div>
        </div>

        <div className="form-section">
          <h3>心情选择</h3>
          <div className="mood-options">
            {moodOptions.map(mood => (
              <div
                key={mood.id}
                className={`mood-option ${selectedMood === mood.name ? 'selected' : ''}`}
                style={{ borderColor: mood.color }}
                onClick={() => setSelectedMood(mood.name)}
              >
                <span className="mood-icon">{mood.icon}</span>
                <span className="mood-name">{mood.name}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="form-section">
          <h3>日记内容</h3>
          <div className="form-group">
            <label htmlFor="content">内容</label>
            <textarea
              id="content"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="请输入日记内容..."
              rows={8}
              required
            ></textarea>
          </div>
        </div>

        <div className="form-section">
          <h3>上传照片</h3>
          <div className="photo-upload-section">
            <div className="photo-upload-area">
              <input
                type="file"
                id="photo-upload"
                accept="image/*"
                multiple
                onChange={handlePhotoUpload}
                className="photo-input"
              />
              <label htmlFor="photo-upload" className="photo-upload-label">
                <div className="upload-icon">📁</div>
                <p>点击或拖拽文件到此处上传</p>
                <p className="upload-hint">支持JPG、PNG格式，最多上传9张</p>
              </label>
            </div>

            {uploadedPhotos.length > 0 && (
              <div className="photo-preview">
                {uploadedPhotos.map(photo => (
                  <div key={photo.id} className="photo-item">
                    <img src={photo.url} alt="预览" />
                    <button
                      type="button"
                      className="delete-photo-btn"
                      onClick={() => handleDeletePhoto(photo.id)}
                    >
                      ×
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        <div className="form-actions">
          <button type="button" className="cancel-btn" onClick={() => navigate('/diary')}>取消</button>
          <button type="submit" className="submit-btn" disabled={isSubmitting}>
            {isSubmitting ? '保存中...' : '保存日记'}
          </button>
        </div>
      </form>
    </div>
  )
}

export default CreateDiaryPage