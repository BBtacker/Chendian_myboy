import React from 'react'
import { BrowserRouter as Router, Routes, Route, useLocation, Navigate } from 'react-router-dom'
import { AnimatePresence } from 'framer-motion'
import { ConfigProvider } from 'antd'
import antdTheme from './theme/antdTheme'
import './styles/App.css'
import Sidebar from './components/Sidebar'
import Header from './components/Header'
import AskPage from './pages/AskPage'
import DiaryPage from './pages/DiaryPage'
import CreateDiaryPage from './pages/CreateDiaryPage'
import ProfilePage from './pages/ProfilePage'
import SharePage from './pages/SharePage'
import SettingsPage from './pages/SettingsPage'
import DreamPage from './pages/DreamPage'
import LoginPage from './pages/LoginPage'
import GamePage from './pages/GamePage'

// 页面过渡组件
const PageTransition = ({ children }) => {
  const location = useLocation();

  return (
    <AnimatePresence mode="wait">
      <div
        key={location.pathname}
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: -20 }}
        transition={{ duration: 0.3 }}
        className="page-transition"
      >
        {children}
      </div>
    </AnimatePresence>
  );
};

function HomePage() {
  return (
    <div className="content">
      <div className="homepage">
        <div className="hero-section">
          <div className="hero-content">
            <h1 className="hero-title">心影日记</h1>
            <p className="hero-subtitle">记录生活点滴，分享美好回忆</p>
            <div className="hero-cta">
              <a href="/diary/create" className="cta-button">开始记录</a>
              <a href="/share" className="cta-secondary">浏览分享</a>
            </div>
          </div>
        </div>
        
        <div className="features-section">
          <h2 className="section-title">核心功能</h2>
          <div className="features-grid">
            <a href="/diary/create" className="feature-card feature-link">
              <div className="feature-icon">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 19V5M12 5L8 9M12 5L16 9" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                </svg>
              </div>
              <h3 className="feature-title">日记记录</h3>
              <p className="feature-description">轻松记录每日心情，上传照片，留下美好回忆</p>
            </a>
            
            <a href="/share" className="feature-card feature-link">
              <div className="feature-icon">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M18 21V5C18 3.89543 17.1046 3 16 3H8C6.89543 3 6 3.89543 6 5V21M16 3L21 8M16 3L11 8" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                </svg>
              </div>
              <h3 className="feature-title">日记分享</h3>
              <p className="feature-description">与他人分享你的日记，交流生活感悟</p>
            </a>
            
            <a href="/ask" className="feature-card feature-link">
              <div className="feature-icon">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 18C15.3137 18 18 15.3137 18 12C18 8.68629 15.3137 6 12 6C8.68629 6 6 8.68629 6 12C6 15.3137 8.68629 18 12 18Z" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M12 14V16" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M12 8V10" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                </svg>
              </div>
              <h3 className="feature-title">智能聊天</h3>
              <p className="feature-description">与AI助手聊天，分享心情，获得温暖回应</p>
            </a>
            
            <a href="/dream" className="feature-card feature-link">
              <div className="feature-icon">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M12 16V12M12 8" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                </svg>
              </div>
              <h3 className="feature-title">心影绘梦</h3>
              <p className="feature-description">上传照片，根据心情生成卡通风格的梦幻图像</p>
            </a>
            
            <a href="/share" className="feature-card feature-link">
              <div className="feature-icon">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M16 21V19C16 17.9391 15.5786 16.9217 14.8284 16.1716C14.0783 15.4214 13.0609 15 12 15C10.9391 15 9.92172 15.4214 9.17157 16.1716C8.42143 16.9217 8 17.9391 8 19V21" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M21 16V8C21 6.89543 20.1046 6 19 6H5C3.89543 6 3 6.89543 3 8V16C3 17.1046 3.89543 18 5 18H19C20.1046 18 21 17.1046 21 16Z" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M8 11C8 9.89543 8.89543 9 10 9H14C15.1046 9 16 9.89543 16 11V13C16 14.1046 15.1046 15 14 15H10C8.89543 15 8 14.1046 8 13V11Z" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                </svg>
              </div>
              <h3 className="feature-title">日记分享圈</h3>
              <p className="feature-description">浏览他人分享的日记，交流生活感悟</p>
            </a>
            
            <a href="/game" className="feature-card feature-link">
              <div className="feature-icon">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M21 6H3C1.89543 6 1 6.89543 1 8V16C1 17.1046 1.89543 18 3 18H21C22.1046 18 23 17.1046 23 16V8C23 6.89543 22.1046 6 21 6Z" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M7 12H9" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M11 12H17" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                </svg>
              </div>
              <h3 className="feature-title">开心大作战</h3>
              <p className="feature-description">用开心机关枪消灭所有负面情绪，放松心情</p>
            </a>
            <a href="/profile" className="feature-card feature-link">
              <div className="feature-icon">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M20 21V19C20 17.9391 19.5786 16.9217 18.8284 16.1716C18.0783 15.4214 17.0609 15 16 15H8C6.93913 15 5.92172 15.4214 5.17157 16.1716C4.42143 16.9217 4 17.9391 4 19V21" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M12 11C14.2091 11 16 9.20914 16 7C16 4.79086 14.2091 3 12 3C9.79086 3 8 4.79086 8 7C8 9.20914 9.79086 11 12 11Z" stroke="#E11D48" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                </svg>
              </div>
              <h3 className="feature-title">个人主页</h3>
              <p className="feature-description">查看和管理你的个人资料与日记</p>
            </a>
          </div>
        </div>
        

        
        <div className="about-section">
          <h2 className="section-title">关于心影日记</h2>
          <p className="about-description">心影日记是一个专注于记录生活点滴和分享美好回忆的平台。在这里，你可以写下每日的所思所想，上传珍贵的照片，记录自己的情绪变化，并与他人分享你的生活感悟。通过心影日记，你可以更好地了解自己，珍藏美好回忆，与他人建立联系，让生活更加充实和有意义。</p>
        </div>
        
        <div className="cta-section">
          <div className="cta-content">
            <h2 className="cta-title">开始你的心影之旅</h2>
            <p className="cta-description">记录生活，分享感动，让每一个瞬间都成为永恒</p>
            <a href="/diary/create" className="cta-button large">立即开始</a>
          </div>
        </div>
      </div>
    </div>
  );
}

// 受保护的路由组件
const ProtectedRoute = ({ children }) => {
  const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true';
  return isLoggedIn ? children : <Navigate to="/login" replace />;
};

// 公共路由组件（不显示侧边栏和头部）
const PublicRoute = ({ children }) => {
  return <div className="public-page">{children}</div>;
};

function App() {
  return (
    <ConfigProvider theme={antdTheme}>
    <Router>
      <Routes>
        {/* 公共路由 - 不显示侧边栏和头部 */}
        <Route 
          path="/login" 
          element={
            <PublicRoute>
              <LoginPage />
            </PublicRoute>
          } 
        />
        
        {/* 受保护的路由 - 显示侧边栏和头部 */}
        <Route 
          path="/*" 
          element={
            <ProtectedRoute>
              <div className="app-container">
                <Sidebar />
                <div className="main-content">
                  <Header />
                  <PageTransition>
                    <Routes>
                      <Route path="/" element={<HomePage />} />
                      <Route path="/ask" element={<AskPage />} />
                      <Route path="/diary" element={<DiaryPage />} />
                      <Route path="/diary/create" element={<CreateDiaryPage />} />
                      <Route path="/profile" element={<ProfilePage />} />
                      <Route path="/share" element={<SharePage />} />
                      <Route path="/dream" element={<DreamPage />} />
                      <Route path="/settings" element={<SettingsPage />} />
                      <Route path="/game" element={<GamePage />} />
                    </Routes>
                  </PageTransition>
                </div>
              </div>
            </ProtectedRoute>
          } 
        />
      </Routes>
    </Router>
    </ConfigProvider>
  )
}

export default App
