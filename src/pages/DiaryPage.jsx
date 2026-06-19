import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import '../styles/DiaryPage.css'
import { uploadImageToOSS, deleteImageFromOSS } from '../utils/doubaoApi'
import { message } from 'antd'
import { diaryApi } from '../utils/api'
import { filterDiaries } from '../utils/diaryData'

function DiaryPage() {
  const navigate = useNavigate()
  // 从localStorage获取日记数据
  const [allDiaries, setAllDiaries] = useState([])
  // 显示的日记列表（可能经过筛选）
  const [displayedDiaries, setDisplayedDiaries] = useState([])
  // 筛选条件
  const [filters, setFilters] = useState({
    startDate: '',
    endDate: '',
    mood: ''
  })

  // 从localStorage获取当前用户信息
  const [currentUser, setCurrentUser] = useState(null)

  // 当组件挂载时加载日记数据
  useEffect(() => {
    // 获取当前用户信息
    const user = JSON.parse(localStorage.getItem('currentUser'))
    if (user) {
      setCurrentUser(user)
      // 从后端API获取日记数据
      fetchDiaries(user.id)
    } else {
      // 如果用户未登录，跳转到登录页面
      navigate('/login')
    }
  }, [])

  // 从后端API获取日记数据
  const fetchDiaries = async (userId) => {
    try {
      const diaries = await diaryApi.getDiariesByUserId(userId)
      
      // 确保diaries是数组再处理
      if (!Array.isArray(diaries)) {
        console.warn('日记数据格式异常:', diaries)
        setAllDiaries([])
        setDisplayedDiaries([])
        return
      }
      
      // 为每个日记获取照片信息
      const diariesWithPhotos = await Promise.all(
        diaries.map(async (diary) => {
          try {
            const photos = await diaryApi.getPhotosByDiaryId(diary.id)
            return {
              ...diary,
              uploadedPhotos: photos.map(photo => ({
                id: photo.id,
                url: photo.photoUrl
              }))
            }
          } catch (error) {
            console.error(`获取日记 ${diary.id} 的照片失败:`, error)
            return diary
          }
        })
      )
      
      setAllDiaries(diariesWithPhotos)
      setDisplayedDiaries(diariesWithPhotos)
    } catch (error) {
      console.error('获取日记数据失败:', error)
      message.error('获取日记数据失败，请稍后重试')
    }
  }

  // 处理筛选条件变化
  const handleFilterChange = (e) => {
    const { name, value } = e.target
    setFilters(prev => ({ ...prev, [name]: value }))
  }

  // 处理筛选
  const handleFilter = () => {
    const filtered = filterDiaries(allDiaries, filters)
    setDisplayedDiaries(filtered)
  }

  // 处理重置
  const handleReset = () => {
    setFilters({
      startDate: '',
      endDate: '',
      mood: ''
    })
    setDisplayedDiaries(allDiaries)
  }

  // 处理删除日记
  const handleDeleteDiary = async (diaryId) => {
    if (window.confirm('确定要删除这篇日记吗？')) {
      // 先获取要删除的日记
      const diaryToDelete = allDiaries.find(d => d.id === diaryId);
      
      // 删除关联的OSS图片
      if (diaryToDelete && diaryToDelete.uploadedPhotos) {
        for (const photo of diaryToDelete.uploadedPhotos) {
          if (photo.url) {
            try {
              await deleteImageFromOSS(photo.url);
            } catch (error) {
              console.error('从OSS删除图片失败:', error);
            }
          }
        }
      }
      
      // 删除日记
      try {
        await diaryApi.deleteDiary(diaryId, currentUser.id);
        // 更新本地状态
        const updatedDiaries = allDiaries.filter(d => d.id !== diaryId);
        setAllDiaries(updatedDiaries);
        setDisplayedDiaries(updatedDiaries);
        message.success('日记删除成功！');
      } catch (error) {
        console.error('删除日记失败:', error);
        message.error('删除日记失败，请稍后重试');
      }
    }
  }

  // 查看和编辑状态
  const [currentDiary, setCurrentDiary] = useState(null);
  const [isViewModalVisible, setIsViewModalVisible] = useState(false);
  const [isEditModalVisible, setIsEditModalVisible] = useState(false);
  
  // 编辑表单状态
  const [editForm, setEditForm] = useState({
    title: '',
    content: '',
    diaryDate: new Date().toISOString().split('T')[0],
    selectedMood: '',
    uploadedPhotos: []
  });

  // 处理查看日记详情
  const handleViewDiary = (diaryId) => {
    const diary = allDiaries.find(d => d.id === diaryId);
    if (diary) {
      setCurrentDiary(diary);
      setIsViewModalVisible(true);
    }
  };

  // 处理编辑日记
  const handleEditDiary = (diary) => {
    setCurrentDiary(diary);
    setEditForm({
      title: diary.title,
      content: diary.content,
      diaryDate: diary.diaryDate,
      selectedMood: diary.selectedMood || '',
      uploadedPhotos: diary.uploadedPhotos || []
    });
    setIsEditModalVisible(true);
  };

  // 处理编辑表单输入变化
  const handleEditInputChange = (field, value) => {
    setEditForm(prev => ({
      ...prev,
      [field]: value
    }));
  };

  // 处理编辑模态框中的图片上传
  const handleEditPhotoUpload = async (e) => {
    const files = Array.from(e.target.files);
    
    for (const file of files) {
      try {
        // 先上传到阿里云OSS
        const imageUrl = await uploadImageToOSS(file);
        
        const photoData = {
          id: Date.now() + Math.random(),
          name: file.name,
          type: file.type,
          url: imageUrl // 保存OSS返回的URL
        };
        setEditForm(prev => ({
          ...prev,
          uploadedPhotos: [...prev.uploadedPhotos, photoData]
        }));
      } catch (error) {
        console.error('上传图片失败:', error);
        message.error(`上传图片失败: ${error.message}`);
      }
    }
  };

  // 删除编辑模态框中的照片
  const handleDeleteEditPhoto = async (photoId) => {
    const photoToDelete = editForm.uploadedPhotos.find(photo => photo.id === photoId);
    
    // 如果是OSS图片，先从OSS删除
    if (photoToDelete && photoToDelete.url) {
      try {
        await deleteImageFromOSS(photoToDelete.url);
      } catch (error) {
        console.error('从OSS删除图片失败:', error);
      }
    }
    
    // 从本地状态删除
    setEditForm(prev => ({
      ...prev,
      uploadedPhotos: prev.uploadedPhotos.filter(photo => photo.id !== photoId)
    }));
  };

  // 保存编辑后的日记
  const saveEditedDiary = async () => {
    if (!editForm.title.trim() || !editForm.content.trim()) {
      message.warning('请填写所有必填项');
      return;
    }

    // 获取原始日记的图片
    const originalPhotos = currentDiary.uploadedPhotos || [];
    const newPhotos = editForm.uploadedPhotos || [];
    
    // 找出被删除的图片
    const deletedPhotos = originalPhotos.filter(photo => 
      !newPhotos.some(newPhoto => newPhoto.id === photo.id)
    );
    
    // 从OSS删除被删除的图片
    for (const photo of deletedPhotos) {
      if (photo.url) {
        try {
          await deleteImageFromOSS(photo.url);
        } catch (error) {
          console.error('从OSS删除图片失败:', error);
        }
      }
    }

    // 准备更新的日记数据
    const updatedDiary = {
      id: currentDiary.id,
      userId: currentUser.id,
      title: editForm.title,
      content: editForm.content,
      diaryDate: editForm.diaryDate,
      selectedMood: editForm.selectedMood,
      uploadedPhotos: editForm.uploadedPhotos,
      updatedAt: new Date().toISOString()
    };

    // 调用后端API更新日记
    try {
      const result = await diaryApi.updateDiary(updatedDiary);
      if (result) {
        // 更新本地状态
        const updatedDiaries = allDiaries.map(diary => {
          if (diary.id === currentDiary.id) {
            return result;
          }
          return diary;
        });
        setAllDiaries(updatedDiaries);
        setDisplayedDiaries(updatedDiaries);
        
        // 关闭模态框
        setIsEditModalVisible(false);
        message.success('日记更新成功！');
      } else {
        message.error('更新日记失败，请稍后重试');
      }
    } catch (error) {
      console.error('更新日记失败:', error);
      message.error('更新日记失败，请稍后重试');
    }
  };

  return (
    <div className="diary-page">
      <div className="page-header">
        <h2>日记列表</h2>
        <button className="create-btn" onClick={() => navigate('/diary/create')}>
          <span className="btn-icon">+</span>
          写新日记
        </button>
      </div>

      <div className="diary-filters">
        <div className="filter-group">
          <label htmlFor="startDate">日期范围</label>
          <input 
            type="date" 
            id="startDate" 
            name="startDate"
            value={filters.startDate}
            onChange={handleFilterChange}
          />
          <span className="filter-separator">至</span>
          <input 
            type="date" 
            id="endDate"
            name="endDate"
            value={filters.endDate}
            onChange={handleFilterChange}
          />
        </div>
        <div className="filter-group">
          <label htmlFor="mood">心情</label>
          <select 
            id="mood" 
            name="mood"
            value={filters.mood}
            onChange={handleFilterChange}
          >
            <option value="">全部</option>
            <option value="开心">开心</option>
            <option value="悲伤">悲伤</option>
            <option value="愤怒">愤怒</option>
            <option value="惊讶">惊讶</option>
            <option value="自然">自然</option>
          </select>
        </div>
        <div className="filter-group">
          <button className="filter-btn" onClick={handleFilter}>筛选</button>
          <button className="reset-btn" onClick={handleReset}>重置</button>
        </div>
      </div>

      <div className="diary-list">
        {displayedDiaries.length === 0 ? (
          <div className="empty-state">
            <div className="empty-icon">📝</div>
            <h3>还没有日记</h3>
            <p>点击右上角的"写新日记"按钮开始记录吧</p>
          </div>
        ) : (
          displayedDiaries.map(diary => (
            <div key={diary.id} className="diary-card">
              {/* 兼容旧格式和新格式的图片显示，放在前面 */}
              {((diary.uploadedPhotos && diary.uploadedPhotos.length > 0) || (diary.photos && diary.photos.length > 0)) && (
                <div className="diary-photos">
                  {/* 优先使用新格式的uploadedPhotos */}
                  {(diary.uploadedPhotos || []).map((photo, index) => (
                    <div key={photo.id || index} className="diary-photo-item">
                      <img src={photo.url || photo.data} alt={`日记照片 ${index + 1}`} />
                    </div>
                  ))}
                  {/* 兼容旧格式的photos */}
                  {(diary.photos || []).map((photo, index) => (
                    <div key={index} className="diary-photo-item">
                      <img src={photo} alt={`日记照片 ${index + 1}`} />
                    </div>
                  ))}
                </div>
              )}
              <div className="diary-header">
                <div className="diary-meta">
                  <h3 className="diary-title">{diary.title}</h3>
                  <div className="diary-info">
                    <span className="diary-date">{diary.diaryDate}</span>
                    <span className="diary-mood">
                      <span className="mood-icon">
                        {diary.selectedMood === '开心' && '😊'}
                        {diary.selectedMood === '悲伤' && '😢'}
                        {diary.selectedMood === '愤怒' && '😠'}
                        {diary.selectedMood === '惊讶' && '�'}
                        {(diary.selectedMood === '自然' || !diary.selectedMood) && '�'}
                      </span>
                      {diary.selectedMood || '自然'}
                    </span>
                  </div>
                </div>
                <div className="diary-actions">
                  <button className="action-btn view-btn" onClick={() => handleViewDiary(diary.id)}>
                    查看
                  </button>
                  <button className="action-btn edit-btn" onClick={() => handleEditDiary(diary)}>
                    编辑
                  </button>
                  <button className="action-btn delete-btn" onClick={() => handleDeleteDiary(diary.id)}>
                    删除
                  </button>
                </div>
              </div>
              <div className="diary-content">
                <p>{diary.content}</p>
              </div>
              <div className="diary-footer">
                <span className="diary-created-at">
                  创建于 {diary.createdAt}
                </span>
              </div>
            </div>
          ))
        )}
      </div>

      {/* 查看日记模态框 */}
      {isViewModalVisible && currentDiary && (
        <div className="modal-overlay">
          <div className="modal-content view-modal">
            <div className="modal-header">
              <h3>查看日记</h3>
              <button className="modal-close-btn" onClick={() => setIsViewModalVisible(false)}>×</button>
            </div>
            <div className="modal-body">
              <div className="modal-section">
                <h4>{currentDiary.title}</h4>
                <div className="diary-meta">
                  <span className="diary-date">{currentDiary.diaryDate}</span>
                  <span className="diary-mood">
                    <span className="mood-icon">
                      {currentDiary.selectedMood === '开心' && '😊'}
                      {currentDiary.selectedMood === '悲伤' && '😢'}
                      {currentDiary.selectedMood === '愤怒' && '�'}
                      {currentDiary.selectedMood === '惊讶' && '😲'}
                      {(currentDiary.selectedMood === '自然' || !currentDiary.selectedMood) && '�😐'}
                    </span>
                    {currentDiary.selectedMood || '自然'}
                  </span>
                </div>
              </div>
              <div className="modal-section">
                <p>{currentDiary.content}</p>
              </div>
              {((currentDiary.uploadedPhotos && currentDiary.uploadedPhotos.length > 0) || (currentDiary.photos && currentDiary.photos.length > 0)) && (
                <div className="modal-section">
                  <h5>照片</h5>
                  <div className="modal-photos">
                    {/* 优先使用新格式的uploadedPhotos */}
                    {(currentDiary.uploadedPhotos || []).map((photo, index) => (
                      <div key={photo.id || index} className="modal-photo-item">
                        <img src={photo.url || photo.data} alt={`日记照片 ${index + 1}`} />
                      </div>
                    ))}
                    {/* 兼容旧格式的photos */}
                    {(currentDiary.photos || []).map((photo, index) => (
                      <div key={index} className="modal-photo-item">
                        <img src={photo} alt={`日记照片 ${index + 1}`} />
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
            <div className="modal-footer">
              <button className="modal-btn close-btn" onClick={() => setIsViewModalVisible(false)}>
                关闭
              </button>
            </div>
          </div>
        </div>
      )}

      {/* 编辑日记模态框 */}
      {isEditModalVisible && currentDiary && (
        <div className="modal-overlay">
          <div className="modal-content edit-modal">
            <div className="modal-header">
              <h3>编辑日记</h3>
              <button className="modal-close-btn" onClick={() => setIsEditModalVisible(false)}>×</button>
            </div>
            <div className="modal-body">
              <div className="form-group">
                <label htmlFor="edit-title">日记标题</label>
                <input
                  type="text"
                  id="edit-title"
                  value={editForm.title}
                  onChange={(e) => handleEditInputChange('title', e.target.value)}
                  placeholder="请输入日记标题"
                />
              </div>
              <div className="form-row">
                <div className="form-group">
                  <label htmlFor="edit-date">日期</label>
                  <input
                    type="date"
                    id="edit-date"
                    value={editForm.diaryDate}
                    onChange={(e) => handleEditInputChange('diaryDate', e.target.value)}
                  />
                </div>
                <div className="form-group">
                  <label>心情</label>
                  <div className="mood-options">
                    {[
                      { name: '开心', icon: '😊', color: '#f1c40f' },
                      { name: '悲伤', icon: '😢', color: '#3498db' },
                      { name: '愤怒', icon: '😠', color: '#e74c3c' },
                      { name: '惊讶', icon: '😲', color: '#9b59b6' },
                      { name: '自然', icon: '😐', color: '#2ecc71' }
                    ].map(mood => (
                      <div
                        key={mood.name}
                        className={`mood-option ${editForm.selectedMood === mood.name ? 'selected' : ''}`}
                        style={{ borderColor: mood.color }}
                        onClick={() => handleEditInputChange('selectedMood', mood.name)}
                      >
                        <span className="mood-icon">{mood.icon}</span>
                        <span className="mood-name">{mood.name}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
              <div className="form-group">
                <label htmlFor="edit-content">内容</label>
                <textarea
                  id="edit-content"
                  value={editForm.content}
                  onChange={(e) => handleEditInputChange('content', e.target.value)}
                  placeholder="请输入日记内容..."
                  rows={6}
                ></textarea>
              </div>
              <div className="form-group">
                <label>图片管理</label>
                <div className="photo-upload-section">
                  <div className="photo-upload-area">
                    <input
                      type="file"
                      accept="image/*"
                      multiple
                      onChange={handleEditPhotoUpload}
                      className="photo-input"
                      id="edit-photo-upload"
                    />
                    <label className="photo-upload-label" htmlFor="edit-photo-upload">
                      <div className="upload-icon">📁</div>
                      <p>点击或拖拽文件到此处上传</p>
                      <p className="upload-hint">支持JPG、PNG格式，最多上传9张</p>
                    </label>
                  </div>

                  {editForm.uploadedPhotos && editForm.uploadedPhotos.length > 0 && (
                    <div className="photo-preview">
                      {editForm.uploadedPhotos.map(photo => (
                        <div key={photo.id} className="photo-item">
                          <img src={photo.url || photo.data} alt="预览" />
                          <button
                            type="button"
                            className="delete-photo-btn"
                            onClick={() => handleDeleteEditPhoto(photo.id)}
                          >
                            ×
                          </button>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </div>
            <div className="modal-footer">
              <button className="modal-btn cancel-btn" onClick={() => setIsEditModalVisible(false)}>
                取消
              </button>
              <button className="modal-btn save-btn" onClick={saveEditedDiary}>
                保存
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default DiaryPage