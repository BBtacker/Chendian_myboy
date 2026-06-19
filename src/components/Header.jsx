import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import './Header.css'
import { diaryApi } from '../utils/api'

function Header() {
  const navigate = useNavigate()
  const [showLogoutMenu, setShowLogoutMenu] = useState(false)
  const [currentUser, setCurrentUser] = useState(null)
  const [todayDiariesCount, setTodayDiariesCount] = useState(0)
  const [totalDiariesCount, setTotalDiariesCount] = useState(0)

  // 从localStorage获取用户信息和日记统计
  useEffect(() => {
    // 获取用户信息
    const user = localStorage.getItem('currentUser')
    if (user) {
      setCurrentUser(JSON.parse(user))
    }

    // 计算日记统计数据
    const calculateDiaryStats = async () => {
      const user = localStorage.getItem('currentUser')
      if (user) {
        const currentUser = JSON.parse(user)
        const userId = currentUser.id
        
        if (userId) {
          try {
            const diaries = await diaryApi.getDiariesByUserId(userId)
            const today = new Date().toISOString().split('T')[0]
            
            // 确保diaries是数组再操作
            if (Array.isArray(diaries)) {
              const todayCount = diaries.filter(diary => {
                return diary.diaryDate && diary.diaryDate.startsWith(today)
              }).length
              setTodayDiariesCount(todayCount)
              setTotalDiariesCount(diaries.length)
            } else {
              console.warn('日记数据格式异常:', diaries)
              setTodayDiariesCount(0)
              setTotalDiariesCount(0)
            }
          } catch (error) {
            console.error('获取日记统计数据失败:', error)
          }
        }
      }
    }

    // 立即执行一次
    const fetchStats = async () => {
      await calculateDiaryStats()
    }
    fetchStats()
    
    // 监听localStorage变化，实时更新统计数据
    const handleStorageChange = () => {
      calculateDiaryStats()
    }

    window.addEventListener('storage', handleStorageChange)
    return () => window.removeEventListener('storage', handleStorageChange)
  }, [])

  // 处理退出登录
  const handleLogout = () => {
    localStorage.removeItem('currentUser')
    localStorage.removeItem('isLoggedIn')
    setCurrentUser(null)
    setShowLogoutMenu(false)
    navigate('/login')
  }

  // 切换退出菜单显示
  const toggleLogoutMenu = () => {
    setShowLogoutMenu(!showLogoutMenu)
  }

  return (
    <div className="header">
      <div className="header-left">
        <div className="header-title">
          <h1>心影日记</h1>
          <p className="subtitle">记录生活点滴，珍藏美好回忆</p>
        </div>
      </div>
      
      <div className="header-center">
        <div className="search-box">
          <input 
            type="text" 
            placeholder="搜索日记、心情或照片..." 
            className="search-input"
          />
          <button className="search-btn">🔍</button>
        </div>
      </div>
      
      <div className="header-right">
        <div className="header-stats">
          <div className="stat-item">
            <div className="stat-number">{todayDiariesCount}</div>
            <div className="stat-label">今日日记</div>
          </div>
          <div className="stat-divider"></div>
          <div className="stat-item">
            <div className="stat-number">{totalDiariesCount}</div>
            <div className="stat-label">日记总数</div>
          </div>
        </div>
        
        <div className="header-user" onClick={toggleLogoutMenu}>
          <div className="user-info">
            <div className="user-name">{currentUser?.username || '用户'}</div>
            <div className="user-role">普通用户</div>
          </div>
          <div className="user-avatar">
            {currentUser?.avatar ? (
              <img src={currentUser.avatar} alt="用户头像" className="avatar-img" />
            ) : (
              <span>{(currentUser?.username || 'US').substring(0, 2).toUpperCase()}</span>
            )}
          </div>
          
          {/* 退出菜单 */}
          {showLogoutMenu && (
            <div className="logout-menu">
              <div className="menu-item" onClick={() => navigate('/profile')}>
                <span className="menu-icon">👤</span>
                <span className="menu-text">个人主页</span>
              </div>
              <div className="menu-item" onClick={handleLogout}>
                <span className="menu-icon">🚪</span>
                <span className="menu-text">退出登录</span>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default Header
