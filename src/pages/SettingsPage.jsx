import React, { useState, useEffect } from 'react'
import '../styles/SettingsPage.css'

function SettingsPage() {
  // 主题颜色选项
  const themeColors = [
    { name: '默认', color: '#6366f1' },
    { name: '红色', color: '#ef4444' },
    { name: '绿色', color: '#22c55e' },
    { name: '蓝色', color: '#3b82f6' },
    { name: '紫色', color: '#a855f7' },
    { name: '橙色', color: '#f97316' }
  ]

  // 字体大小选项
  const fontSizeOptions = [
    { name: '小', size: '14px' },
    { name: '中', size: '16px' },
    { name: '大', size: '18px' },
    { name: '特大', size: '20px' }
  ]

  // 当前设置
  const [currentTheme, setCurrentTheme] = useState(localStorage.getItem('themeColor') || '#6366f1')
  const [currentFontSize, setCurrentFontSize] = useState(localStorage.getItem('fontSize') || '16px')
  const [notifications, setNotifications] = useState(JSON.parse(localStorage.getItem('notifications')) || true)
  const [darkMode, setDarkMode] = useState(JSON.parse(localStorage.getItem('darkMode')) || false)

  // 应用主题颜色
  useEffect(() => {
    document.documentElement.style.setProperty('--primary-color', currentTheme)
    localStorage.setItem('themeColor', currentTheme)
  }, [currentTheme])

  // 应用字体大小
  useEffect(() => {
    document.documentElement.style.setProperty('--font-size', currentFontSize)
    localStorage.setItem('fontSize', currentFontSize)
  }, [currentFontSize])

  // 应用深色模式
  useEffect(() => {
    if (darkMode) {
      document.documentElement.classList.add('dark-mode')
    } else {
      document.documentElement.classList.remove('dark-mode')
    }
    localStorage.setItem('darkMode', JSON.stringify(darkMode))
  }, [darkMode])

  // 应用通知设置
  useEffect(() => {
    localStorage.setItem('notifications', JSON.stringify(notifications))
  }, [notifications])

  return (
    <div className="settings-page">
      <div className="settings-header">
        <h2>系统设置</h2>
        <p>自定义你的心影日记体验</p>
      </div>

      <div className="settings-content">
        {/* 主题设置 */}
        <div className="settings-section">
          <h3 className="section-title">主题设置</h3>
          
          <div className="setting-item">
            <label className="setting-label">主题颜色</label>
            <div className="theme-colors">
              {themeColors.map((theme, index) => (
                <div 
                  key={index}
                  className={`color-option ${currentTheme === theme.color ? 'active' : ''}`}
                  style={{ backgroundColor: theme.color }}
                  onClick={() => setCurrentTheme(theme.color)}
                  title={theme.name}
                >
                  {currentTheme === theme.color && <span className="checkmark">✓</span>}
                </div>
              ))}
            </div>
          </div>

          <div className="setting-item">
            <label className="setting-label">深色模式</label>
            <div className="toggle-switch">
              <input 
                type="checkbox" 
                id="dark-mode-toggle"
                checked={darkMode}
                onChange={(e) => setDarkMode(e.target.checked)}
              />
              <label htmlFor="dark-mode-toggle" className="toggle-label"></label>
            </div>
          </div>
        </div>

        {/* 字体设置 */}
        <div className="settings-section">
          <h3 className="section-title">字体设置</h3>
          
          <div className="setting-item">
            <label className="setting-label">字体大小</label>
            <div className="font-size-options">
              {fontSizeOptions.map((option, index) => (
                <button 
                  key={index}
                  className={`font-size-btn ${currentFontSize === option.size ? 'active' : ''}`}
                  style={{ fontSize: option.size }}
                  onClick={() => setCurrentFontSize(option.size)}
                >
                  {option.name}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* 通知设置 */}
        <div className="settings-section">
          <h3 className="section-title">通知设置</h3>
          
          <div className="setting-item">
            <label className="setting-label">接收通知</label>
            <div className="toggle-switch">
              <input 
                type="checkbox" 
                id="notifications-toggle"
                checked={notifications}
                onChange={(e) => setNotifications(e.target.checked)}
              />
              <label htmlFor="notifications-toggle" className="toggle-label"></label>
            </div>
          </div>
        </div>

        {/* 关于 */}
        <div className="settings-section">
          <h3 className="section-title">关于</h3>
          
          <div className="about-info">
            <p className="app-name">心影日记</p>
            <p className="app-version">版本 1.0.0</p>
            <p className="app-description">记录生活点滴，分享美好回忆</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default SettingsPage