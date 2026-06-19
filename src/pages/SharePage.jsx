import React, { useState, useEffect } from 'react'
import '../styles/SharePage.css'
import { getCommentsFromStorage, addComment, deleteComment } from '../utils/diaryData'
import { getPlaceholderImage, getAvatarPlaceholder } from '../utils/placeholders'
import { Modal, message } from 'antd'

// 将 picsum.photos URL 替换为本地占位图
function replacePicsumWithPlaceholder(url) {
  if (!url || !url.includes('picsum.photos')) return url
  const match = url.match(/\/seed\/(\w+)\/(\d+)\/(\d+)/)
  if (match) {
    const [, seed, width, height] = match
    return getPlaceholderImage(parseInt(width), parseInt(height), seed)
  }
  return getPlaceholderImage(400, 300, Math.random().toString())
}

// 模拟日记分享数据
const rawDiaries = [
    {
      id: 1,
      user: {
        id: 1,
        username: '用户123',
        avatar: 'https://picsum.photos/seed/user1/100/100'
      },
      title: '美好的一天',
      content: '今天天气真好，和朋友一起去公园玩，拍了很多照片，心情非常愉快！',
      date: '2025-11-28',
      mood: '开心',
      moodIcon: '😊',
      likes: 25,
      comments: 12,
      shares: 5,
      photos: ['https://picsum.photos/seed/park1/400/300', 'https://picsum.photos/seed/friends1/400/300'],
      createdAt: '2025-11-28 15:30:00'
    },
    {
      id: 2,
      user: {
        id: 2,
        username: '用户456',
        avatar: 'https://picsum.photos/seed/user2/100/100'
      },
      title: '美食之旅',
      content: '今天尝试了新的餐厅，食物非常美味，尤其是甜点，强烈推荐！',
      date: '2025-11-27',
      mood: '开心',
      moodIcon: '😊',
      likes: 45,
      comments: 20,
      shares: 8,
      photos: ['https://picsum.photos/seed/food1/400/500', 'https://picsum.photos/seed/dessert1/400/400'],
      createdAt: '2025-11-27 20:15:00'
    },
    {
      id: 3,
      user: {
        id: 3,
        username: '用户789',
        avatar: 'https://picsum.photos/seed/user3/100/100'
      },
      title: '工作压力大',
      content: '最近工作任务太多，感觉压力很大，需要好好调整一下状态。',
      date: '2025-11-26',
      mood: '悲伤',
      moodIcon: '😢',
      likes: 15,
      comments: 8,
      shares: 3,
      photos: ['https://picsum.photos/seed/work1/400/350', 'https://picsum.photos/seed/stress1/400/350'],
      createdAt: '2025-11-26 18:45:00'
    },
    {
      id: 4,
      user: {
        id: 1,
        username: '用户123',
        avatar: 'https://picsum.photos/seed/user1/100/100'
      },
      title: '旅行日记',
      content: '今天去了海边，看到了美丽的日落，心情非常放松。',
      date: '2025-11-25',
      mood: '惊讶',
      moodIcon: '😲',
      likes: 60,
      comments: 25,
      shares: 15,
      photos: ['https://picsum.photos/seed/beach1/400/450', 'https://picsum.photos/seed/sunset1/400/300', 'https://picsum.photos/seed/ocean1/400/500'],
      createdAt: '2025-11-25 17:30:00'
    },
    {
      id: 5,
      user: {
        id: 4,
        username: '用户101',
        avatar: 'https://picsum.photos/seed/user4/100/100'
      },
      title: '学习新技能',
      content: '今天开始学习摄影，虽然一开始有点困难，但看到自己拍的第一张照片，还是很有成就感的！',
      date: '2025-11-24',
      mood: '自然',
      moodIcon: '😐',
      likes: 30,
      comments: 15,
      shares: 7,
      photos: ['https://picsum.photos/seed/camera1/400/380', 'https://picsum.photos/seed/photography1/400/380'],
      createdAt: '2025-11-24 14:20:00'
    },
    {
      id: 6,
      user: {
        id: 5,
        username: '用户102',
        avatar: 'https://picsum.photos/seed/user5/100/100'
      },
      title: '宠物日常',
      content: '我家的小猫今天学会了新把戏，太可爱了！忍不住拍了很多照片记录下来。',
      date: '2025-11-23',
      mood: '开心',
      moodIcon: '😊',
      likes: 55,
      comments: 30,
      shares: 12,
      photos: ['https://picsum.photos/seed/cat1/400/420', 'https://picsum.photos/seed/kitten1/400/360'],
      createdAt: '2025-11-23 19:45:00'
    },
    {
      id: 7,
      user: {
        id: 6,
        username: '用户103',
        avatar: 'https://picsum.photos/seed/user6/100/100'
      },
      title: '健身打卡',
      content: '坚持健身一个月了，虽然很累，但看到自己的变化，觉得一切都值得！',
      date: '2025-11-22',
      mood: '自然',
      moodIcon: '😐',
      likes: 40,
      comments: 22,
      shares: 9,
      photos: ['https://picsum.photos/seed/gym1/400/480', 'https://picsum.photos/seed/fitness1/400/480'],
      createdAt: '2025-11-22 18:15:00'
    },
    {
      id: 8,
      user: {
        id: 2,
        username: '用户456',
        avatar: 'https://picsum.photos/seed/user2/100/100'
      },
      title: '电影观后感',
      content: '今天看了一部非常感人的电影，哭了好几次。电影真的能带给我们很多思考。',
      date: '2025-11-21',
      mood: '悲伤',
      moodIcon: '😢',
      likes: 35,
      comments: 18,
      shares: 6,
      photos: ['https://picsum.photos/seed/cinema1/400/320', 'https://picsum.photos/seed/movie1/400/320'],
      createdAt: '2025-11-21 21:30:00'
    },
    {
      id: 9,
      user: {
        id: 7,
        username: '用户104',
        avatar: 'https://picsum.photos/seed/user7/100/100'
      },
      title: '花园一角',
      content: '自己种的花终于开了，每天看着它们成长，心情都变得很好。',
      date: '2025-11-20',
      mood: '开心',
      moodIcon: '😊',
      likes: 48,
      comments: 24,
      shares: 10,
      photos: ['https://picsum.photos/seed/garden1/400/460', 'https://picsum.photos/seed/flowers1/400/340'],
      createdAt: '2025-11-20 16:50:00'
    },
    {
      id: 10,
      user: {
        id: 8,
        username: '用户105',
        avatar: 'https://picsum.photos/seed/user8/100/100'
      },
      title: '加班日常',
      content: '又是一个加班的夜晚，希望项目能顺利完成。',
      date: '2025-11-19',
      mood: '生气',
      moodIcon: '😠',
      likes: 22,
      comments: 14,
      shares: 4,
      photos: ['https://picsum.photos/seed/office1/400/390', 'https://picsum.photos/seed/night1/400/390'],
      createdAt: '2025-11-19 22:10:00'
    },
    {
      id: 11,
      user: {
        id: 3,
        username: '用户789',
        avatar: 'https://picsum.photos/seed/user3/100/100'
      },
      title: '朋友聚会',
      content: '好久没见的朋友今天聚在一起，聊了很多往事，感觉时间过得真快。',
      date: '2025-11-18',
      mood: '开心',
      moodIcon: '😊',
      likes: 52,
      comments: 28,
      shares: 13,
      photos: ['https://picsum.photos/seed/party1/400/440', 'https://picsum.photos/seed/friends2/400/370', 'https://picsum.photos/seed/gathering1/400/410'],
      createdAt: '2025-11-18 20:45:00'
    },
    {
      id: 12,
      user: {
        id: 9,
        username: '用户106',
        avatar: 'https://picsum.photos/seed/user9/100/100'
      },
      title: '阅读时光',
      content: '今天读了一本好书，收获满满。阅读真的能让人平静下来。',
      date: '2025-11-17',
      mood: '自然',
      moodIcon: '😐',
      likes: 33,
      comments: 16,
      shares: 8,
      photos: ['https://picsum.photos/seed/book1/400/330', 'https://picsum.photos/seed/library1/400/330'],
      createdAt: '2025-11-17 13:20:00'
    }
  ]

function SharePage() {
  // 用本地占位图替换 picsum.photos（在国内网络环境下可正常显示）
  const [diaries, setDiaries] = useState(() =>
    rawDiaries.map(diary => ({
      ...diary,
      user: { ...diary.user, avatar: replacePicsumWithPlaceholder(diary.user.avatar) },
      photos: diary.photos.map(p => replacePicsumWithPlaceholder(p)),
    }))
  )

  // 从localStorage获取评论数据
  const [comments, setComments] = useState([])

  // 初始化评论数据
  useEffect(() => {
    const storedComments = getCommentsFromStorage()
    // 如果没有评论数据，添加一些初始评论
    if (storedComments.length === 0) {
      const initialComments = [
        // 日记1的评论
        {
          id: 1,
          diaryId: 1,
          user: {
            id: 2,
            username: '用户456',
            avatar: 'https://picsum.photos/seed/user2/100/100'
          },
          content: '看起来真的很开心呢！',
          createdAt: '2025-11-28 16:00:00'
        },
        {
          id: 2,
          diaryId: 1,
          user: {
            id: 3,
            username: '用户789',
            avatar: 'https://picsum.photos/seed/user3/100/100'
          },
          content: '风景真美，下次一起去吧！',
          createdAt: '2025-11-28 16:30:00'
        },
        {
          id: 3,
          diaryId: 1,
          user: {
            id: 4,
            username: '用户101',
            avatar: 'https://picsum.photos/seed/user4/100/100'
          },
          content: '好羡慕你们的生活！',
          createdAt: '2025-11-28 17:00:00'
        },
        // 日记2的评论
        {
          id: 4,
          diaryId: 2,
          user: {
            id: 1,
            username: '用户123',
            avatar: 'https://picsum.photos/seed/user1/100/100'
          },
          content: '这家餐厅看起来不错，下次我也要去试试！',
          createdAt: '2025-11-27 21:00:00'
        },
        {
          id: 5,
          diaryId: 2,
          user: {
            id: 3,
            username: '用户789',
            avatar: 'https://picsum.photos/seed/user3/100/100'
          },
          content: '甜点看起来好好吃！',
          createdAt: '2025-11-27 21:30:00'
        },
        // 日记3的评论
        {
          id: 6,
          diaryId: 3,
          user: {
            id: 1,
            username: '用户123',
            avatar: 'https://picsum.photos/seed/user1/100/100'
          },
          content: '工作压力大的时候要注意休息哦！',
          createdAt: '2025-11-26 19:00:00'
        },
        {
          id: 7,
          diaryId: 3,
          user: {
            id: 2,
            username: '用户456',
            avatar: 'https://picsum.photos/seed/user2/100/100'
          },
          content: '可以试试冥想，有助于缓解压力！',
          createdAt: '2025-11-26 19:30:00'
        },
        // 日记4的评论
        {
          id: 8,
          diaryId: 4,
          user: {
            id: 5,
            username: '用户102',
            avatar: 'https://picsum.photos/seed/user5/100/100'
          },
          content: '海边的日落真的太美了！',
          createdAt: '2025-11-25 18:00:00'
        },
        {
          id: 9,
          diaryId: 4,
          user: {
            id: 6,
            username: '用户103',
            avatar: 'https://picsum.photos/seed/user6/100/100'
          },
          content: '我也想去海边放松一下！',
          createdAt: '2025-11-25 18:30:00'
        },
        {
          id: 10,
          diaryId: 4,
          user: {
            id: 7,
            username: '用户104',
            avatar: 'https://picsum.photos/seed/user7/100/100'
          },
          content: '照片拍得真好！',
          createdAt: '2025-11-25 19:00:00'
        }
      ]
      // 保存初始评论到localStorage
      initialComments.forEach(comment => {
        addComment(comment)
      })
      setComments(initialComments)
    } else {
      setComments(storedComments)
    }
  }, [])

  // 当前查看的日记
  const [currentDiary, setCurrentDiary] = useState(null)
  // 评论输入内容
  const [commentInput, setCommentInput] = useState('')
  // 点赞状态
  const [likedDiaries, setLikedDiaries] = useState(new Set())

  // 计算每个日记的实际评论数量
  const calculateCommentCounts = () => {
    const commentCounts = {};
    diaries.forEach(diary => {
      const count = comments.filter(comment => comment.diaryId === diary.id).length;
      commentCounts[diary.id] = count;
    });
    return commentCounts;
  };

  // 更新日记的评论数量
  const updateDiaryCommentCounts = () => {
    const commentCounts = calculateCommentCounts();
    setDiaries(prevDiaries => 
      prevDiaries.map(diary => ({
        ...diary,
        comments: commentCounts[diary.id] || 0
      }))
    );
  };

  // 组件挂载时更新评论数量
  useEffect(() => {
    updateDiaryCommentCounts();
  }, [comments]);

  // 打开日记详情
  const openDiaryDetail = (diary) => {
    setCurrentDiary(diary)
  }

  // 关闭日记详情
  const closeDiaryDetail = () => {
    setCurrentDiary(null)
    setCommentInput('')
  }

  // 切换点赞状态
  const toggleLike = (diaryId) => {
    const newLikedDiaries = new Set(likedDiaries)
    if (newLikedDiaries.has(diaryId)) {
      newLikedDiaries.delete(diaryId)
    } else {
      newLikedDiaries.add(diaryId)
    }
    setLikedDiaries(newLikedDiaries)
  }

  // 发送评论
  const sendComment = () => {
    if (!commentInput.trim() || !currentDiary) return

    const newComment = {
      diaryId: currentDiary.id,
      user: {
        id: 1,
        username: '用户123',
        avatar: 'https://picsum.photos/seed/user1/100/100'
      },
      content: commentInput.trim(),
      createdAt: new Date().toLocaleString()
    }

    // 使用真实的addComment函数添加评论
    const addedComment = addComment(newComment)
    setComments([...comments, addedComment])
    setCommentInput('')
    
    // 更新当前日记的评论数
    setCurrentDiary(prev => ({
      ...prev,
      comments: prev.comments + 1
    }))
  }

  // 处理删除评论
  const handleDeleteComment = (commentId) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这条评论吗？删除后将无法恢复。',
      okText: '删除',
      okType: 'danger',
      cancelText: '取消',
      onOk() {
        // 使用真实的deleteComment函数删除评论
        const updatedComments = deleteComment(commentId)
        setComments(updatedComments)
        
        // 更新当前日记的评论数
        if (currentDiary) {
          const currentDiaryComments = updatedComments.filter(comment => comment.diaryId === currentDiary.id)
          setCurrentDiary(prev => ({
            ...prev,
            comments: currentDiaryComments.length
          }))
        }
        
        message.success('评论删除成功！')
      },
      onCancel() {
        console.log('取消删除')
      }
    })
  }

  return (
    <div className="share-page">
      <div className="page-header">
        <h2>日记分享圈</h2>
        <p>发现他人的精彩生活，分享你的故事</p>
      </div>

      {!currentDiary ? (
        <div className="diary-waterfall">
          {diaries.map(diary => (
            <div key={diary.id} className="diary-card" onClick={() => openDiaryDetail(diary)}>
              <div className="diary-header">
                <div className="user-avatar">
                  <img src={diary.user.avatar} alt={diary.user.username} />
                </div>
                <div className="user-info">
                  <div className="username">{diary.user.username}</div>
                  <div className="diary-date">{diary.date}</div>
                </div>
              </div>
              <div className="diary-content">
                <h3 className="diary-title">{diary.title}</h3>
                <p className="diary-text">{diary.content}</p>
              </div>
              <div className="diary-photos">
                {diary.photos.map((photo, index) => (
                  <div key={index} className="diary-photo" style={{ height: `${Math.floor(Math.random() * 100) + 200}px` }}>
                    <img src={photo} alt={`日记照片 ${index + 1}`} />
                  </div>
                ))}
              </div>
              <div className="diary-stats">
                <div className="stat-item">
                  <span className="stat-icon">
                    {diary.moodIcon || 
                     (diary.mood === '开心' && '😊') ||
                     (diary.mood === '悲伤' && '😢') ||
                     (diary.mood === '愤怒' && '😠') ||
                     (diary.mood === '惊讶' && '😲') ||
                     (diary.mood === '自然' && '😐') ||
                     (diary.mood === '生气' && '😠') ||
                     '😊'}
                  </span>
                  <span className="stat-text">{diary.mood}</span>
                </div>
                <div className="stat-divider"></div>
                <div className="stat-item">
                  <span className="stat-icon">❤️</span>
                  <span className="stat-text">{likedDiaries.has(diary.id) ? diary.likes + 1 : diary.likes}</span>
                </div>
                <div className="stat-item">
                  <span className="stat-icon">💬</span>
                  <span className="stat-text">{diary.comments}</span>
                </div>
                <div className="stat-item">
                  <span className="stat-icon">🔄</span>
                  <span className="stat-text">{diary.shares}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="diary-detail-overlay" onClick={closeDiaryDetail}>
          <div className="diary-detail" onClick={(e) => e.stopPropagation()}>
            <div className="diary-detail-header">
              <button className="close-btn" onClick={closeDiaryDetail}>×</button>
              <div className="user-info">
                <div className="user-avatar">
                  <img src={currentDiary.user.avatar} alt={currentDiary.user.username} />
                </div>
                <div className="user-name-date">
                  <div className="username">{currentDiary.user.username}</div>
                  <div className="diary-date">{currentDiary.date}</div>
                </div>
              </div>
            </div>
            <div className="diary-detail-content">
              <h2 className="diary-title">{currentDiary.title}</h2>
              <div className="diary-mood">
                <span className="mood-icon">
                  {currentDiary.moodIcon || 
                   (currentDiary.mood === '开心' && '😊') ||
                   (currentDiary.mood === '悲伤' && '😢') ||
                   (currentDiary.mood === '愤怒' && '😠') ||
                   (currentDiary.mood === '惊讶' && '😲') ||
                   (currentDiary.mood === '自然' && '😐') ||
                   (currentDiary.mood === '生气' && '😠') ||
                   '😊'}
                </span>
                <span className="mood-text">{currentDiary.mood}</span>
              </div>
              <p className="diary-text">{currentDiary.content}</p>
              <div className="diary-detail-photos">
                {currentDiary.photos.map((photo, index) => (
                  <div key={index} className="diary-detail-photo">
                    <img src={photo} alt={`日记照片 ${index + 1}`} />
                  </div>
                ))}
              </div>
            </div>
            <div className="diary-detail-actions">
              <button 
                className={`action-btn like-btn ${likedDiaries.has(currentDiary.id) ? 'liked' : ''}`}
                onClick={() => toggleLike(currentDiary.id)}
              >
                <span className="action-icon">❤️</span>
                <span className="action-text">{likedDiaries.has(currentDiary.id) ? currentDiary.likes + 1 : currentDiary.likes}</span>
              </button>
              <button className="action-btn comment-btn">
                <span className="action-icon">💬</span>
                <span className="action-text">{currentDiary.comments}</span>
              </button>
              <button className="action-btn share-btn">
                <span className="action-icon">🔄</span>
                <span className="action-text">{currentDiary.shares}</span>
              </button>
            </div>
            <div className="diary-comments">
              <h3 className="comments-title">评论 ({currentDiary.comments})</h3>
              <div className="comments-list">
                {comments.filter(comment => comment.diaryId === currentDiary.id).map(comment => (
                  <div key={comment.id} className="comment-item">
                    <div className="comment-avatar">
                      <img src={comment.user.avatar} alt={comment.user.username} />
                    </div>
                    <div className="comment-content">
                      <div className="comment-header">
                        <span className="comment-username">{comment.user.username}</span>
                        <span className="comment-time">{comment.createdAt}</span>
                        {/* 只有当前用户（id为1）的评论才显示删除按钮 */}
                        {comment.user.id === 1 && (
                          <button 
                            className="delete-comment-btn"
                            onClick={() => handleDeleteComment(comment.id)}
                          >
                            删除
                          </button>
                        )}
                      </div>
                      <p className="comment-text">{comment.content}</p>
                    </div>
                  </div>
                ))}
              </div>
              <div className="comment-input-section">
                <textarea 
                  value={commentInput} 
                  onChange={(e) => setCommentInput(e.target.value)}
                  placeholder="写下你的评论..."
                  rows={2}
                />
                <button 
                  className="send-comment-btn" 
                  onClick={sendComment}
                  disabled={!commentInput.trim()}
                >
                  发送
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default SharePage