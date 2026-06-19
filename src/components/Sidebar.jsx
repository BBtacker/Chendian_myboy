import React from 'react'
import { useLocation } from 'react-router-dom'
import './Sidebar.css'
import logo from '../assets/logo.png'

function Sidebar() {
  // 获取当前路径
  const location = useLocation();
  
  // 导航菜单项数据
  const navItems = [
    { id: 1, name: '首页', icon: '🏠', href: '/' },
    { id: 2, name: '写日记', icon: '📝', href: '/diary/create' },
    { id: 3, name: '日记列表', icon: '📚', href: '/diary' },
    { id: 4, name: '日记分享圈', icon: '🌍', href: '/share' },
    { id: 5, name: '个人主页', icon: '👤', href: '/profile' },
    { id: 6, name: '智能聊天', icon: '💬', href: '/ask' },
    { id: 7, name: '心影绘梦', icon: '🎨', href: '/dream' },
    { id: 9, name: '情绪分拣站', icon: '🎮', href: '/game' },
    { id: 8, name: '系统设置', icon: '⚙️', href: '/settings' }
  ];

  return (
    <div className="sidebar">
      <div className="sidebar-logo">
        <div className="logo-icon">
          <img src={logo} alt="MyDiary" className="logo-img" />
        <h3>MyDiary</h3></div>
        <h3>心影日记</h3>
      </div>
      <nav className="sidebar-nav">
        <ul>
          {navItems.map((item) => (
            <li 
              key={item.id} 
              className={`nav-item ${location.pathname === item.href ? 'active' : ''}`}
            >
              <a href={item.href}>
                <span className="nav-icon">{item.icon}</span>
                <span className="nav-text">{item.name}</span>
              </a>
            </li>
          ))}
        </ul>
      </nav>
      <div className="sidebar-footer">
        <div className="system-info">
          <div className="version">版本: 1.0.0</div>
          <div className="status">
            <span className="status-indicator"></span>
            <span>系统正常</span>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Sidebar
