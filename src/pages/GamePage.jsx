import React, { useState, useEffect, useRef, useCallback } from 'react'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler } from 'chart.js'
import { Line } from 'react-chartjs-2'
import '../styles/GamePage.css'
import { Card, Progress, Button, Statistic, Row, Col } from 'antd'
import { gameApi } from '../utils/api'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const EMOTIONS = {
  negative: [
    { name: '悲伤', color: '#3498db', emoji: '😢' },
    { name: '愤怒', color: '#e74c3c', emoji: '😠' },
    { name: '焦虑', color: '#9b59b6', emoji: '😰' },
    { name: '疲惫', color: '#95a5a6', emoji: '😴' }
  ],
  positive: [
    { name: '开心', color: '#f1c40f', emoji: '😄' },
    { name: '平静', color: '#2ecc71', emoji: '😌' },
    { name: '幸福', color: '#e91e63', emoji: '🥰' },
    { name: '感动', color: '#00bcd4', emoji: '🥲' }
  ],
  bomb: { name: '纠结', color: '#1a1a1a', emoji: '💣' }
}

function GamePage() {
  const [gameState, setGameState] = useState('start')
  const [isPaused, setIsPaused] = useState(false)
  const [score, setScore] = useState(0)
  const [highScore, setHighScore] = useState(0)
  const [combo, setCombo] = useState(0)
  const [energy, setEnergy] = useState(0)
  const [timeLeft, setTimeLeft] = useState(60)
  const [gripperX, setGripperX] = useState(50)
  const [gripperY, setGripperY] = useState(50)
  const [packages, setPackages] = useState([])
  const [heldPackage, setHeldPackage] = useState(null)
  const [shake, setShake] = useState(false)
  const [glowLeft, setGlowLeft] = useState(false)
  const [glowRight, setGlowRight] = useState(false)
  const [showStars, setShowStars] = useState(false)
  const [soundEnabled, setSoundEnabled] = useState(true)
  const [gameStats, setGameStats] = useState({ negative: 0, positive: 0 })
  const [gameHistory, setGameHistory] = useState([])
  const [todayNegative, setTodayNegative] = useState(0)
  const [todayPositive, setTodayPositive] = useState(0)
  const [droppingPackage, setDroppingPackage] = useState(null)
  const [dropProgress, setDropProgress] = useState(0)

  const gameLoopRef = useRef(null)
  const timerRef = useRef(null)
  const packageIdRef = useRef(0)
  const audioContextRef = useRef(null)
  const slowModeRef = useRef(false)

  useEffect(() => {
    const savedHistory = localStorage.getItem('gameHistory')
    const savedToday = localStorage.getItem('todayStats')
    const savedDate = localStorage.getItem('todayDate')
    const savedGameState = localStorage.getItem('gameState')
    const today = new Date().toDateString()
    
    if (savedHistory) {
      setGameHistory(JSON.parse(savedHistory))
    }
    
    if (savedToday && savedDate === today) {
      const stats = JSON.parse(savedToday)
      setTodayNegative(stats.negative || 0)
      setTodayPositive(stats.positive || 0)
    } else {
      setTodayNegative(0)
      setTodayPositive(0)
    }
    
    if (savedGameState) {
      const gameStateData = JSON.parse(savedGameState)
      setScore(gameStateData.score || 0)
      setHighScore(gameStateData.highScore || 0)
      setCombo(gameStateData.combo || 0)
      setEnergy(gameStateData.energy || 0)
      setTimeLeft(gameStateData.timeLeft || 60)
      setGripperX(gameStateData.gripperX || 50)
      setGripperY(gameStateData.gripperY || 50)
      setGameStats(gameStateData.gameStats || { negative: 0, positive: 0 })
      if (gameStateData.gameState === 'playing') {
        setGameState('playing')
        
      }
    }
  }, [])

  const saveGameHistory = useCallback((finalScore) => {
    const newHistory = [...gameHistory, {
      score: finalScore,
      time: new Date().toLocaleTimeString()
    }].slice(-10)
    setGameHistory(newHistory)
    localStorage.setItem('gameHistory', JSON.stringify(newHistory))
  }, [gameHistory])

  const saveTodayStats = useCallback((negative, positive) => {
    const newNegative = todayNegative + negative
    const newPositive = todayPositive + positive
    setTodayNegative(newNegative)
    setTodayPositive(newPositive)
    localStorage.setItem('todayStats', JSON.stringify({ negative: newNegative, positive: newPositive }))
    localStorage.setItem('todayDate', new Date().toDateString())
  }, [todayNegative, todayPositive])

  // 保存游戏状态到localStorage
  useEffect(() => {
    if (gameState === 'playing') {
      const gameStateData = {
        gameState,
        score,
        highScore,
        combo,
        energy,
        timeLeft,
        gripperX,
        gripperY,
        gameStats
      }
      localStorage.setItem('gameState', JSON.stringify(gameStateData))
    }
  }, [gameState, score, highScore, combo, energy, timeLeft, gripperX, gripperY, gameStats])

  const initAudio = useCallback(() => {
    if (!audioContextRef.current) {
      audioContextRef.current = new (window.AudioContext || window.webkitAudioContext)()
    }
  }, [])

  const playSound = useCallback((type) => {
    if (!soundEnabled || !audioContextRef.current) return
    
    const ctx = audioContextRef.current
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    
    osc.connect(gain)
    gain.connect(ctx.destination)
    
    switch(type) {
      case 'correct':
        osc.frequency.setValueAtTime(523, ctx.currentTime)
        osc.frequency.setValueAtTime(659, ctx.currentTime + 0.1)
        gain.gain.setValueAtTime(0.2, ctx.currentTime)
        gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.3)
        osc.start(ctx.currentTime)
        osc.stop(ctx.currentTime + 0.3)
        break
      case 'wrong':
        osc.type = 'sawtooth'
        osc.frequency.setValueAtTime(200, ctx.currentTime)
        gain.gain.setValueAtTime(0.15, ctx.currentTime)
        gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.2)
        osc.start(ctx.currentTime)
        osc.stop(ctx.currentTime + 0.2)
        break
      case 'pickup':
        osc.frequency.setValueAtTime(440, ctx.currentTime)
        gain.gain.setValueAtTime(0.1, ctx.currentTime)
        gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.1)
        osc.start(ctx.currentTime)
        osc.stop(ctx.currentTime + 0.1)
        break
      case 'heal':
        const notes = [262, 330, 392, 523]
        notes.forEach((freq, i) => {
          const o = ctx.createOscillator()
          const g = ctx.createGain()
          o.connect(g)
          g.connect(ctx.destination)
          o.frequency.setValueAtTime(freq, ctx.currentTime + i * 0.15)
          g.gain.setValueAtTime(0.15, ctx.currentTime + i * 0.15)
          g.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + i * 0.15 + 0.2)
          o.start(ctx.currentTime + i * 0.15)
          o.stop(ctx.currentTime + i * 0.15 + 0.2)
        })
        break
      case 'drop':
        osc.frequency.setValueAtTime(300, ctx.currentTime)
        osc.frequency.exponentialRampToValueAtTime(100, ctx.currentTime + 0.3)
        gain.gain.setValueAtTime(0.15, ctx.currentTime)
        gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.3)
        osc.start(ctx.currentTime)
        osc.stop(ctx.currentTime + 0.3)
        break
    }
  }, [soundEnabled])

  const spawnPackage = useCallback(() => {
    const isBomb = Math.random() < 0.1
    let emotion
    let type
    
    if (isBomb) {
      emotion = EMOTIONS.bomb
      type = 'bomb'
    } else {
      const isPositive = Math.random() < 0.5
      type = isPositive ? 'positive' : 'negative'
      const emotions = EMOTIONS[type]
      emotion = emotions[Math.floor(Math.random() * emotions.length)]
    }
    
    return {
      id: packageIdRef.current++,
      x: 20 + Math.random() * 60,
      y: 0,
      emotion,
      type,
      falling: true
    }
  }, [])

  const triggerShake = () => {
    setShake(true)
    setTimeout(() => setShake(false), 300)
  }

  const triggerGlow = (side) => {
    if (side === 'left') {
      setGlowLeft(true)
      setTimeout(() => setGlowLeft(false), 500)
    } else {
      setGlowRight(true)
      setTimeout(() => setGlowRight(false), 500)
    }
  }

  const triggerHealing = () => {
    setShowStars(true)
    playSound('heal')
    slowModeRef.current = true
    setTimeout(() => {
      setShowStars(false)
      slowModeRef.current = false
    }, 3000)
  }

  const checkDropZone = useCallback((pkg, gripperPos) => {
    if (!pkg) return null
    
    const gripperLeft = gripperPos - 10
    const gripperRight = gripperPos + 10
    
    if (gripperRight < 30) {
      return 'left'
    } else if (gripperLeft > 70) {
      return 'right'
    }
    return null
  }, [])

  const startGame = () => {
    initAudio()
    setGameState('playing')
    setScore(0)
    setCombo(0)
    setEnergy(0)
    setTimeLeft(60)
    setGripperX(50)
    setGripperY(50)
    setPackages([])
    setHeldPackage(null)
    setGameStats({ negative: 0, positive: 0 })
    setDroppingPackage(null)
    setIsPaused(false)
    slowModeRef.current = false
    packageIdRef.current = 0
    
    // 保存初始游戏状态
    const gameStateData = {
      gameState: 'playing',
      score: 0,
      highScore,
      combo: 0,
      energy: 0,
      timeLeft: 60,
      gripperX: 50,
      gripperY: 50,
      gameStats: { negative: 0, positive: 0 }
    }
    localStorage.setItem('gameState', JSON.stringify(gameStateData))
  }

  const endGame = () => {
    setGameState('ended')
    const finalHighScore = Math.max(highScore, score)
    setHighScore(finalHighScore)
    saveGameHistory(score)
    saveTodayStats(gameStats.negative, gameStats.positive)
    if (gameLoopRef.current) {
      clearInterval(gameLoopRef.current)
    }
    if (timerRef.current) {
      clearInterval(timerRef.current)
    }
    
    // 保存游戏得分到后端
    try {
      const currentUser = localStorage.getItem('currentUser')
      if (currentUser) {
        const user = JSON.parse(currentUser)
        gameApi.saveScore({
          userId: user.id,
          score,
          highScore: finalHighScore,
          positiveCount: gameStats.positive,
          negativeCount: gameStats.negative,
          combo,
          energy,
          gameDate: new Date().toISOString().split('T')[0]
        })
      }
    } catch (e) {
      console.warn('保存游戏得分到后端失败:', e)
    }
    
    // 游戏结束后清除保存的状态
    localStorage.removeItem('gameState')
  }

  const exitGame = () => {
    setGameState('start')
    if (gameLoopRef.current) {
      clearInterval(gameLoopRef.current)
    }
    if (timerRef.current) {
      clearInterval(timerRef.current)
    }
    
    // 退出游戏后清除保存的状态
    localStorage.removeItem('gameState')
  }

  const generateDiaryDraft = () => {
    const { negative, positive } = gameStats
    let draft = ''
    if (positive > negative) {
      draft = `今天你消解了 ${negative} 份负面情绪，收集了 ${positive} 份正面能量，是能量满满的一天！`
    } else if (negative > positive) {
      draft = `今天你消解了 ${negative} 份负面情绪，收集了 ${positive} 份正面能量，虽然有些挑战，但你做得很棒！`
    } else {
      draft = `今天你消解了 ${negative} 份负面情绪，收集了 ${positive} 份正面能量，是平衡的一天！`
    }
    return draft
  }

  useEffect(() => {
    const handleKeyDown = (e) => {
      if (gameState !== 'playing') return
      
      switch(e.key.toLowerCase()) {
        case 'p':
          setIsPaused(!isPaused)
          break
        case 'a':
        case 'arrowleft':
          if (!isPaused) {
            setGripperX(prev => Math.max(10, prev - 5))
          }
          break
        case 'd':
        case 'arrowright':
          if (!isPaused) {
            setGripperX(prev => Math.min(90, prev + 5))
          }
          break
        case 'w':
        case 'arrowup':
          if (!isPaused) {
            setGripperY(prev => Math.max(20, prev - 5))
          }
          break
        case 's':
        case 'arrowdown':
          if (!isPaused) {
            setGripperY(prev => Math.min(80, prev + 5))
          }
          break
        case ' ':
          if (!isPaused) {
            e.preventDefault()
            if (!heldPackage) {
              const nearbyPackage = packages.find(pkg => 
                pkg.falling && 
                Math.abs(pkg.x - gripperX) < 15 && 
                Math.abs(pkg.y - gripperY) < 15
              )
              if (nearbyPackage) {
                setHeldPackage(nearbyPackage)
                setPackages(prev => prev.filter(p => p.id !== nearbyPackage.id))
                playSound('pickup')
              }
            } else {
              const dropZone = checkDropZone(heldPackage, gripperX)
              if (dropZone) {
                setDroppingPackage({
                  ...heldPackage,
                  x: gripperX,
                  startY: gripperY,
                  targetY: 85,
                  startTime: Date.now(),
                  duration: 300,
                  targetZone: dropZone
                })
                setHeldPackage(null)
                setDropProgress(0)
                playSound('drop')
              }
            }
          }
          break
      }
    }
    
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [gameState, isPaused, gripperX, gripperY, heldPackage, packages, combo, checkDropZone, playSound])

  useEffect(() => {
    if (!droppingPackage) return
    
    const animateDrop = () => {
      if (isPaused) {
        requestAnimationFrame(animateDrop)
        return
      }
      
      const elapsed = Date.now() - droppingPackage.startTime
      const progress = Math.min(elapsed / droppingPackage.duration, 1)
      setDropProgress(progress)
      
      if (progress >= 1) {
        let correct = false
        if (droppingPackage.type === 'bomb') {
          correct = droppingPackage.targetZone === 'left'
        } else {
          correct = (droppingPackage.type === 'negative' && droppingPackage.targetZone === 'left') ||
                    (droppingPackage.type === 'positive' && droppingPackage.targetZone === 'right')
        }
        
        if (correct) {
          playSound('correct')
          triggerGlow(droppingPackage.targetZone)
          setCombo(prev => prev + 1)
          const comboBonus = Math.floor(combo / 5) * 5
          setScore(prev => prev + 10 + comboBonus)
          
          if (droppingPackage.type === 'negative' || droppingPackage.type === 'bomb') {
            setGameStats(prev => ({ ...prev, negative: prev.negative + 1 }))
          } else {
            setGameStats(prev => ({ ...prev, positive: prev.positive + 1 }))
            setEnergy(prev => {
              const newEnergy = prev + 1
              if (newEnergy >= 10) {
                triggerHealing()
                return 0
              }
              return newEnergy
            })
          }
          
          if ((combo + 1) % 5 === 0) {
            slowModeRef.current = true
            setTimeout(() => { slowModeRef.current = false }, 5000)
          }
        } else {
          playSound('wrong')
          triggerShake()
          setCombo(0)
          setScore(prev => Math.max(0, prev - 5))
        }
        
        setDroppingPackage(null)
        setDropProgress(0)
        return
      }
      
      requestAnimationFrame(animateDrop)
    }
    
    requestAnimationFrame(animateDrop)
  }, [droppingPackage, isPaused, combo, playSound])

  useEffect(() => {
    if (gameState !== 'playing') return
    
    const timer = setInterval(() => {
      if (!isPaused) {
        setTimeLeft(prev => {
          if (prev <= 1) {
            endGame()
            return 0
          }
          return prev - 1
        })
      }
    }, 1000)
    
    const gameLoop = setInterval(() => {
      if (!isPaused) {
        if (Math.random() < 0.06) {
          setPackages(prev => [...prev, spawnPackage()])
        }
        
        const fallSpeed = slowModeRef.current ? 0.8 : 1.5
        setPackages(prev => prev.map(pkg => ({
          ...pkg,
          y: pkg.y + fallSpeed
        })).filter(pkg => pkg.y < 110))
      }
    }, 100)
    
    timerRef.current = timer
    gameLoopRef.current = gameLoop
    return () => {
      clearInterval(timer)
      clearInterval(gameLoop)
    }
  }, [gameState, isPaused, spawnPackage])

  const chartData = {
    labels: gameHistory.map((_, i) => `第${i + 1}局`),
    datasets: [
      {
        label: '得分',
        data: gameHistory.map(h => h.score),
        borderColor: '#667eea',
        backgroundColor: 'rgba(102, 126, 234, 0.1)',
        borderWidth: 3,
        fill: true,
        tension: 0.4,
        pointRadius: 6,
        pointHoverRadius: 8
      }
    ]
  }

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: {
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        padding: 12,
        titleFont: { size: 14 },
        bodyFont: { size: 13 }
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        grid: { color: 'rgba(0, 0, 0, 0.05)' }
      },
      x: {
        grid: { display: false }
      }
    }
  }

  return (
    <div className={`emotion-sorting-game ${shake ? 'shake' : ''}`}>
      <div className="game-container">
        <div className="game-header">
          <Card className="header-card">
            <div className="header-content">
              <h1>🎮 情绪分拣站</h1>
              <div className="header-buttons">
                <Button 
                  type="text" 
                  icon={soundEnabled ? '🔊' : '🔇'}
                  onClick={() => setSoundEnabled(!soundEnabled)}
                />
                {gameState === 'playing' && (
                  <>
                    <Button type="default" onClick={() => setIsPaused(!isPaused)}>
                      {isPaused ? '继续游戏' : '暂停游戏'}
                    </Button>
                    <Button type="default" danger onClick={exitGame}>
                      退出游戏
                    </Button>
                  </>
                )}
              </div>
            </div>
          </Card>
        </div>

        {gameState === 'start' && (
          <div className="game-start-section">
            <Card className="start-card">
              <div className="start-content">
                <h2>🎯 情绪分拣站</h2>
                <Button type="primary" size="large" onClick={startGame} className="start-btn">
                  开始游戏 🚀
                </Button>
              </div>
            </Card>
          </div>
        )}

        {gameState === 'playing' && (
          <>
            <div className="game-stats">
              <Card className="stat-card">
                <div className="stat-row">
                  <span className="stat-label">⏱️ 时间</span>
                  <span className="stat-value">{timeLeft}s</span>
                </div>
              </Card>
              <Card className="stat-card">
                <div className="stat-row">
                  <span className="stat-label">⭐ 分数</span>
                  <span className="stat-value">{score}</span>
                </div>
              </Card>
              <Card className="stat-card">
                <div className="stat-row">
                  <span className="stat-label">🔥 连击</span>
                  <span className="stat-value">{combo}</span>
                </div>
              </Card>
              <Card className="stat-card energy-card">
                <div className="energy-bar">
                  <span className="stat-label">⚡ 能量</span>
                  <Progress percent={energy * 10} showInfo={false} strokeColor="#e91e63" />
                </div>
              </Card>
            </div>

            <div className="game-area">
              {showStars && [...Array(20)].map((_, i) => (
                <div 
                  key={i} 
                  className="star"
                  style={{
                    left: `${Math.random() * 100}%`,
                    animationDelay: `${Math.random() * 2}s`
                  }}
                >
                  ✨
                </div>
              ))}

              {packages.map(pkg => (
                <div
                  key={pkg.id}
                  className="package"
                  style={{
                    left: `${pkg.x}%`,
                    top: `${pkg.y}%`,
                    background: `linear-gradient(135deg, ${pkg.emotion.color}, ${pkg.emotion.color}dd)`
                  }}
                >
                  <span className="package-emoji">{pkg.emotion.emoji}</span>
                  <span className="package-name">{pkg.emotion.name}</span>
                </div>
              ))}

              {droppingPackage && (
                <div
                  key="dropping"
                  className="package dropping"
                  style={{
                    left: `${droppingPackage.x}%`,
                    top: `${droppingPackage.startY + (droppingPackage.targetY - droppingPackage.startY) * dropProgress}%`,
                    background: `linear-gradient(135deg, ${droppingPackage.emotion.color}, ${droppingPackage.emotion.color}dd)`
                  }}
                >
                  <span className="package-emoji">{droppingPackage.emotion.emoji}</span>
                  <span className="package-name">{droppingPackage.emotion.name}</span>
                </div>
              )}

              <div 
                className="gripper"
                style={{ left: `${gripperX}%`, top: `${gripperY}%` }}
              >
                <svg width="60" height="80" viewBox="0 0 60 80">
                  <rect x="25" y="0" width="10" height="40" fill="#888" rx="2"/>
                  <rect x="5" y="35" width="20" height="8" fill="#666" rx="2"/>
                  <rect x="35" y="35" width="20" height="8" fill="#666" rx="2"/>
                  <rect x="10" y="30" width="8" height="25" fill="#666" rx="2"/>
                  <rect x="42" y="30" width="8" height="25" fill="#666" rx="2"/>
                </svg>
                {heldPackage && (
                  <div 
                    className="held-package"
                    style={{
                      background: `linear-gradient(135deg, ${heldPackage.emotion.color}, ${heldPackage.emotion.color}dd)`
                    }}
                  >
                    <span className="package-emoji">{heldPackage.emotion.emoji}</span>
                    <span className="package-name">{heldPackage.emotion.name}</span>
                  </div>
                )}
              </div>

              <div className={`sorting-box left ${glowLeft ? 'glow' : ''}`}>
                <div className="box-icon">✕</div>
                <div className="box-label">消解箱</div>
                <div className="box-hint">负面情绪</div>
              </div>
              
              <div className={`sorting-box right ${glowRight ? 'glow' : ''}`}>
                <div className="box-icon">⭐</div>
                <div className="box-label">能量箱</div>
                <div className="box-hint">正面情绪</div>
              </div>
              
              {isPaused && (
                <div className="pause-overlay">
                  <div className="pause-content">
                    <h2>⏸️ 游戏暂停</h2>
                    <p>按 P 键或点击继续游戏按钮恢复</p>
                  </div>
                </div>
              )}
            </div>
            
            <div className="game-intro">
              <Card className="intro-card">
                <h3>🎮 情绪分拣站</h3>
                <p className="intro-text">
                  把"处理情绪"具象成"分拣快递"！快速分拣情绪包裹，负面情绪送去消解，正面情绪存入能量库！
                </p>
                <div className="controls-guide">
                  <h4>操作指南：</h4>
                  <ul>
                    <li>W/S 或 上下方向键 - 上下移动</li>
                    <li>A/D 或 左右方向键 - 左右移动</li>
                    <li>空格键 - 抓取 / 释放包裹</li>
                    <li>P键 - 暂停/继续游戏</li>
                  </ul>
                </div>
                <div className="emotion-guide">
                  <h4>情绪分类：</h4>
                  <div className="emotion-types">
                    <div className="emotion-type">
                      <span className="emoji">😢😠😰😴</span>
                      <span>负面情绪 → 消解箱</span>
                    </div>
                    <div className="emotion-type">
                      <span className="emoji">😄😌🥰🥲</span>
                      <span>正面情绪 → 能量箱</span>
                    </div>
                    <div className="emotion-type">
                      <span className="emoji">💣</span>
                      <span>纠结炸弹 → 快速送消解箱！</span>
                    </div>
                  </div>
                </div>
              </Card>
            </div>

          </>
        )}

        {gameState === 'ended' && (
          <div className="overlay">
            <Card className="end-card">
              <h2>💫 游戏结束</h2>
              <div className="final-score">
                <p>最终得分</p>
                <p className="score-number">{score}</p>
                {score >= highScore && score > 0 && (
                  <p className="new-record">🎉 新纪录！</p>
                )}
              </div>
              
              <div className="diary-draft">
                <h3>📝 今日情绪总结</h3>
                <p>{generateDiaryDraft()}</p>
              </div>
              
              <div className="end-buttons">
                <Button type="default" size="large" onClick={exitGame}>
                  返回首页 🏠
                </Button>
                <Button type="primary" size="large" onClick={startGame}>
                  再来一局 🔄
                </Button>
              </div>
            </Card>
          </div>
        )}
      </div>
    </div>
  )
}

export default GamePage
