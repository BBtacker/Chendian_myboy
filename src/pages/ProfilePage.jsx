import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import '../styles/ProfilePage.css'
import { getUserFromStorage, saveUserToStorage } from '../utils/diaryData'
import { Chart as ChartJS, ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, Title } from 'chart.js'
import { Pie, Bar } from 'react-chartjs-2'
import { Modal, Button, message } from 'antd'
import { uploadImageToOSS, deleteImageFromOSS } from '../utils/doubaoApi'
import { diaryApi, userApi } from '../utils/api'

// 注册Chart.js组件
ChartJS.register(
  ArcElement, 
  Tooltip, 
  Legend, 
  CategoryScale, 
  LinearScale, 
  BarElement, 
  Title
)

function ProfilePage() {
  const navigate = useNavigate()
  // 从localStorage获取用户数据
  const [user, setUser] = useState(getUserFromStorage())

  // 从localStorage获取日记数据
  const [diaries, setDiaries] = useState([])

  // 心情统计数据
  const [moodStats, setMoodStats] = useState({})

  // 当组件挂载时加载数据
  useEffect(() => {
    const currentUser = JSON.parse(localStorage.getItem('currentUser'))
    if (currentUser) {
      // 先从后端获取最新的用户信息（包含封面等数据）
      userApi.getUserById(currentUser.id).then(apiUser => {
        if (apiUser) {
          const merged = { ...currentUser, ...apiUser }
          setUser(merged)
          localStorage.setItem('currentUser', JSON.stringify(merged))
          saveUserToStorage(merged)
        }
      }).catch(() => {
        // 后端不可用时使用本地数据
      })
      fetchDiaries(currentUser.id)
    }
  }, [])
  
  // 从后端API获取日记数据
  const fetchDiaries = async (userId) => {
    try {
      const diaries = await diaryApi.getDiariesByUserId(userId)
      
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
      
      setDiaries(diariesWithPhotos)
      calculateMoodStats(diariesWithPhotos)
    } catch (error) {
      console.error('获取日记数据失败:', error)
      message.error('获取日记数据失败，请稍后重试')
    }
  }

  // 计算心情统计数据
  const calculateMoodStats = (diaryList) => {
    const stats = {}
    
    // 初始化心情类型
    const moodTypes = ['开心', '悲伤', '愤怒', '惊讶', '自然']
    moodTypes.forEach(mood => {
      stats[mood] = 0
    })
    
    // 统计各心情的日记数量
    diaryList.forEach(diary => {
      // 使用默认心情'自然'如果没有selectedMood字段
      const mood = diary.selectedMood || '自然'
      if (stats.hasOwnProperty(mood)) {
        stats[mood] = stats[mood] + 1
      } else {
        // 如果心情不在预定义列表中，也默认为'自然'
        stats['自然'] = stats['自然'] + 1
      }
    })
    
    setMoodStats(stats)
  }

  // 处理删除日记
  const handleDeleteDiary = async (diaryId) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这篇日记吗？删除后将无法恢复。',
      okText: '删除',
      okType: 'danger',
      cancelText: '取消',
      async onOk() {
        // 先获取要删除的日记
        const diaryToDelete = diaries.find(d => d.id === diaryId);
        const currentUser = JSON.parse(localStorage.getItem('currentUser'));
        
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
          const updatedDiaries = diaries.filter(d => d.id !== diaryId);
          setDiaries(updatedDiaries);
          calculateMoodStats(updatedDiaries);
          message.success('日记删除成功！');
        } catch (error) {
          console.error('删除日记失败:', error);
          message.error('删除日记失败，请稍后重试');
        }
      },
      onCancel() {
        console.log('取消删除')
      }
    })
  }

  // 编辑模式状态
  const [isEditing, setIsEditing] = useState(false)
  const [editBio, setEditBio] = useState(user.bio)

  // 保存用户数据到后端和本地
  const saveUser = async (updatedUser) => {
    saveUserToStorage(updatedUser)
    setUser(updatedUser)
    // 同步更新localStorage中的currentUser
    if (updatedUser.id) {
      localStorage.setItem('currentUser', JSON.stringify(updatedUser))
    }
    // 保存到后端
    try {
      await userApi.updateUser(updatedUser)
    } catch (e) {
      console.warn('保存用户信息到后端失败:', e)
    }
  }

  // 处理头像上传
  const handleAvatarUpload = async (e) => {
    const file = e.target.files[0]
    if (file) {
      try {
        // 如果有旧头像，先从OSS删除
        if (user.avatar && user.avatar.startsWith('http')) {
          try {
            await deleteImageFromOSS(user.avatar)
          } catch (error) {
            console.error('删除旧头像失败:', error)
          }
        }
        
        // 上传新头像到阿里云OSS
        const imageUrl = await uploadImageToOSS(file)
        const updatedUser = { ...user, avatar: imageUrl }
        await saveUser(updatedUser)
      } catch (error) {
        console.error('上传头像失败:', error)
        message.error(`上传头像失败: ${error.message}`)
      }
    }
  }

  // 处理封面图片上传
  const handleCoverPhotoUpload = async (e) => {
    const file = e.target.files[0]
    if (file) {
      try {
        // 如果有旧封面，先从OSS删除
        if (user.coverPhoto && user.coverPhoto.startsWith('http')) {
          try {
            await deleteImageFromOSS(user.coverPhoto)
          } catch (error) {
            console.error('删除旧封面失败:', error)
          }
        }
        
        // 上传新封面到阿里云OSS
        const imageUrl = await uploadImageToOSS(file)
        const updatedUser = { ...user, coverPhoto: imageUrl }
        await saveUser(updatedUser)
      } catch (error) {
        console.error('上传封面失败:', error)
        message.error(`上传封面失败: ${error.message}`)
      }
    }
  }

  // 处理签名保存
  const handleBioSave = () => {
    const updatedUser = { ...user, bio: editBio }
    saveUserToStorage(updatedUser)
    setUser(updatedUser)
    // 保存到后端
    userApi.updateUser(updatedUser).catch(e => console.warn('保存签名到后端失败:', e))
    setIsEditing(false)
  }

  // 饼图数据配置
  const pieData = {
    labels: Object.keys(moodStats),
    datasets: [
      {
        data: Object.values(moodStats),
        backgroundColor: [
          'rgba(255, 205, 86, 0.8)', // 开心 - 黄色
          'rgba(54, 162, 235, 0.8)', // 悲伤 - 蓝色
          'rgba(255, 99, 132, 0.8)', // 愤怒 - 红色
          'rgba(153, 102, 255, 0.8)', // 惊讶 - 紫色
          'rgba(75, 192, 192, 0.8)' // 自然 - 绿色
        ],
        borderColor: [
          'rgba(255, 205, 86, 1)',
          'rgba(54, 162, 235, 1)',
          'rgba(255, 99, 132, 1)',
          'rgba(153, 102, 255, 1)',
          'rgba(75, 192, 192, 1)'
        ],
        borderWidth: 1
      }
    ]
  }

  // 饼图配置选项
  const pieOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom'
      },
      title: {
        display: true,
        text: '心情分布占比'
      }
    }
  }

  // 条形图数据配置
  const barData = {
    labels: Object.keys(moodStats),
    datasets: [
      {
        label: '日记篇数',
        data: Object.values(moodStats),
        backgroundColor: [
          'rgba(255, 205, 86, 0.8)',
          'rgba(54, 162, 235, 0.8)',
          'rgba(255, 99, 132, 0.8)',
          'rgba(153, 102, 255, 0.8)',
          'rgba(75, 192, 192, 0.8)'
        ],
        borderColor: [
          'rgba(255, 205, 86, 1)',
          'rgba(54, 162, 235, 1)',
          'rgba(255, 99, 132, 1)',
          'rgba(153, 102, 255, 1)',
          'rgba(75, 192, 192, 1)'
        ],
        borderWidth: 1
      }
    ]
  }

  // 条形图配置选项
  const barOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top'
      },
      title: {
        display: true,
        text: '各心情日记篇数'
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          stepSize: 1
        }
      }
    }
  }

  return (
    <div className="profile-page">
      <div className="profile-header">
        <div className={`cover-photo ${!user.coverPhoto || user.coverPhoto === '' ? 'no-cover' : ''}`} style={user.coverPhoto && user.coverPhoto !== '' ? { backgroundImage: `url(${user.coverPhoto})` } : {}}>
          <input 
            type="file" 
            id="cover-upload" 
            accept="image/*" 
            onChange={handleCoverPhotoUpload}
            className="cover-upload-input"
          />
          <label htmlFor="cover-upload" className="cover-upload-btn">
            📷 更换封面
          </label>
          <div className="avatar-section">
            <div className="avatar-container">
              <img src={user.avatar} alt="用户头像" className="avatar" />
              <input 
                type="file" 
                id="avatar-upload" 
                accept="image/*" 
                onChange={handleAvatarUpload}
                className="avatar-upload-input"
              />
              <label htmlFor="avatar-upload" className="avatar-upload-btn">
                📷
              </label>
            </div>
            <div className="user-info">
              {isEditing ? (
                <div className="bio-edit">
                  <textarea 
                    value={editBio} 
                    onChange={(e) => setEditBio(e.target.value)}
                    placeholder="请输入签名..."
                    rows={2}
                  />
                  <div className="bio-edit-actions">
                    <button className="save-btn" onClick={handleBioSave}>保存</button>
                    <button className="cancel-btn" onClick={() => {
                      setEditBio(user.bio)
                      setIsEditing(false)
                    }}>取消</button>
                  </div>
                </div>
              ) : (
                <div className="bio-section">
                  <p className="bio">{user.bio || '这个人很懒，什么都没写'}</p>
                  <button className="edit-bio-btn" onClick={() => setIsEditing(true)}>编辑签名</button>
                </div>
              )}
            </div>
          </div>
        </div>
        <div className="stats-section">
          <div className="stat-item">
            <div className="stat-number">{user.followers}</div>
            <div className="stat-label">粉丝</div>
          </div>
          <div className="stat-divider"></div>
          <div className="stat-item">
            <div className="stat-number">{user.following}</div>
            <div className="stat-label">关注</div>
          </div>
          <div className="stat-divider"></div>
          <div className="stat-item">
            <div className="stat-number">{user.likes}</div>
            <div className="stat-label">获赞</div>
          </div>
        </div>
      </div>

      <div className="profile-content">
        <div className="content-header">
          <h3>我的日记</h3>
          <button className="create-diary-btn" onClick={() => navigate('/diary/create')}>写新日记</button>
        </div>

        {/* 心情统计图表 */}
        <div className="stats-charts">
          <div className="chart-container">
            <h4>心情分布占比</h4>
            <div className="chart-wrapper">
              <Pie data={pieData} options={pieOptions} />
            </div>
          </div>
          <div className="chart-container">
            <h4>各心情日记篇数</h4>
            <div className="chart-wrapper">
              <Bar data={barData} options={barOptions} />
            </div>
          </div>
        </div>

        <div className="diary-grid">
          {diaries.map(diary => (
            <div key={diary.id} className="diary-card">
              {/* 兼容旧格式和新格式的封面图片显示 */}
              {((diary.uploadedPhotos && diary.uploadedPhotos.length > 0) || (diary.photos && diary.photos.length > 0)) && (
                <div className="diary-photo">
                  {/* 优先使用新格式的uploadedPhotos */}
                  {diary.uploadedPhotos && diary.uploadedPhotos.length > 0 ? (
                    <img src={diary.uploadedPhotos[0].url || diary.uploadedPhotos[0].data} alt={diary.title} />
                  ) : (
                    /* 兼容旧格式的photos */
                    <img src={diary.photos[0]} alt={diary.title} />
                  )}
                </div>
              )}
              <div className="diary-info">
                <h4 className="diary-title">{diary.title}</h4>
                <div className="diary-meta">
                  <span className="diary-date">{diary.diaryDate}</span>
                  <span className="diary-mood">
                    <span className="mood-icon">
                      😐
                    </span>
                    自然
                  </span>
                </div>
                <div className="diary-stats">
                  <span className="like-count">❤️ 0</span>
                  <button 
                    className="delete-diary-btn"
                    onClick={(e) => {
                      e.stopPropagation(); // 阻止事件冒泡
                      handleDeleteDiary(diary.id);
                    }}
                  >
                    🗑️ 删除
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export default ProfilePage