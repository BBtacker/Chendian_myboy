import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import '../styles/LoginPage.css'
import { authApi } from '../utils/api'

function LoginPage() {
  const navigate = useNavigate()
  
  // 状态管理
  const [isLogin, setIsLogin] = useState(true) // true: 登录，false: 注册
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  // 表单验证
  const validateForm = () => {
    setError('')
    
    if (!username.trim()) {
      setError('请输入用户名')
      return false
    }
    
    if (!password) {
      setError('请输入密码')
      return false
    }
    
    if (!isLogin && password !== confirmPassword) {
      setError('两次输入的密码不一致')
      return false
    }
    
    if (password.length < 6) {
      setError('密码长度不能少于6位')
      return false
    }
    
    return true
  }

  // 登录/注册请求
  const handleSubmit = async (e) => {
    e.preventDefault()
    
    if (!validateForm()) {
      return
    }
    
    setLoading(true)
    
    try {
      // 发送API请求
      const data = isLogin 
        ? await authApi.login(username, password)
        : await authApi.register(username, password)
      
      if (data.success) {
        // 保存用户信息到localStorage
        localStorage.setItem('currentUser', JSON.stringify(data.user))
        localStorage.setItem('isLoggedIn', 'true')
        
        // 跳转到首页
        navigate('/')
      } else {
        setError(data.message)
      }
    } catch (err) {
      setError(isLogin ? '登录失败，请检查用户名和密码' : '注册失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-page">
      {/* 装饰性背景元素 */}
      <div className="login-background">
        <div className="background-particles">
          {[...Array(20)].map((_, i) => (
            <div 
              key={i} 
              className="particle"
              style={{
                left: `${Math.random() * 100}%`,
                top: `${Math.random() * 100}%`,
                animationDelay: `${Math.random() * 5}s`,
                animationDuration: `${Math.random() * 10 + 10}s`,
                opacity: Math.random() * 0.5 + 0.1
              }}
            ></div>
          ))}
        </div>
      </div>
      
      <div className="login-container">
        <div className="login-header">
          <div className="app-icon">📝</div>
          <h1 className="app-title">MyDiary</h1>
          <p className="app-subtitle">记录生活点滴，分享美好回忆</p>
          <div className="app-slogan">珍藏每一个心动瞬间，让回忆永远鲜活</div>
        </div>
        
        <div className="login-form-container">
          <div className="form-tabs">
            <button 
              className={`tab-btn ${isLogin ? 'active' : ''}`}
              onClick={() => setIsLogin(true)}
            >
              登录
            </button>
            <button 
              className={`tab-btn ${!isLogin ? 'active' : ''}`}
              onClick={() => setIsLogin(false)}
            >
              注册
            </button>
          </div>
          
          <form className="login-form" onSubmit={handleSubmit}>
            <div className="error-container">
              {error && <div className="error-message">{error}</div>}
            </div>
            
            <div className="form-group">
              <label htmlFor="username" className="form-label">用户名</label>
              <input 
                type="text" 
                id="username" 
                className="form-input"
                placeholder="请输入用户名"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                disabled={loading}
              />
            </div>
            
            <div className="form-group">
              <label htmlFor="password" className="form-label">密码</label>
              <input 
                type="password" 
                id="password" 
                className="form-input"
                placeholder="请输入密码"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={loading}
              />
            </div>
            
            {!isLogin && (
              <div className="form-group">
                <label htmlFor="confirmPassword" className="form-label">确认密码</label>
                <input 
                  type="password" 
                  id="confirmPassword" 
                  className="form-input"
                  placeholder="请再次输入密码"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  disabled={loading}
                />
              </div>
            )}
            
            <button 
              type="submit" 
              className="submit-btn"
              disabled={loading}
            >
              {loading ? (
                <span className="loading-spinner"></span>
              ) : (
                isLogin ? '登录' : '注册'
              )}
            </button>
          </form>
          
          <div className="login-footer">
            <p>登录即表示同意我们的<a href="#" className="footer-link">服务条款</a>和<a href="#" className="footer-link">隐私政策</a></p>
            <div className="login-tips">
              💡 提示：使用真实邮箱注册可方便找回密码
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default LoginPage