# 心影日记 - 智能日记平台架构设计

## 1. 项目概述

心影日记是一个基于React和Spring Boot的智能日记平台，专注于为用户提供安全、温暖的情感记录空间，并通过AI技术提供心灵慰藉和支持。

## 2. 技术架构

### 2.1 前端架构

- **框架**：React 18
- **路由**：React Router
- **UI组件**：Element Plus
- **数据可视化**：ECharts
- **HTTP请求**：Axios
- **样式**：CSS3/Flexbox/Grid
- **构建工具**：Vite

### 2.2 后端架构

- **框架**：Spring Boot 3
- **ORM**：MyBatis
- **数据库**：MySQL 8.0
- **开发语言**：Java 17
- **AI服务**：API
- **云存储**：阿里云OSS

## 3. 项目结构

### 3.1 前端结构

```
src/
├── App.jsx              # 应用主组件，包含路由配置
├── components/          # 通用组件
│   ├── Sidebar.jsx      # 侧边栏导航组件
│   └── Header.jsx       # 头部组件
├── pages/               # 页面组件
│   ├── HomePage.jsx     # 首页
│   ├── CreateDiaryPage.jsx  # 写日记页面
│   ├── DiaryPage.jsx    # 日记列表页面
│   ├── AskPage.jsx      # AI智能聊天页面
│   ├── ProfilePage.jsx  # 个人主页
│   ├── SharePage.jsx    # 日记分享页面
│   └── DreamPage.jsx    # 心影绘梦页面
├── styles/              # CSS样式文件
├── utils/               # 工具函数
│   ├── api.js           # API请求工具
│   ├── diaryData.js     # 日记数据管理
│   └── doubaoApi.js     # API封装
└── assets/              # 静态资源
```

### 3.2 后端结构

```
Light_Heart_Diary_01/
├── src/main/java/com/lhy/face_recognition/
│   ├── config/          # 配置类
│   │   ├── CorsConfig.java
│   │   ├── OssConfig.java
│   │   └── RestTemplateConfig.java
│   ├── controller/      # REST API控制器
│   │   ├── AnalysisController.java
│   │   ├── DiaryController.java
│   │   ├── DoubaoController.java
│   │   └── OssController.java
│   ├── entity/          # 实体类
│   │   ├── AnalysisRecord.java
│   │   ├── Diary.java
│   │   ├── DiaryMood.java
│   │   ├── DiaryPhoto.java
│   │   └── Mood.java
│   ├── mapper/          # MyBatis映射器
│   │   ├── AnalysisRecordMapper.java
│   │   ├── DiaryMapper.java
│   │   ├── DiaryMoodMapper.java
│   │   ├── DiaryPhotoMapper.java
│   │   └── MoodMapper.java
│   ├── service/         # 业务逻辑层
│   │   ├── impl/        # 实现类
│   │   ├── AnalysisRecordService.java
│   │   ├── DiaryService.java
│   │   ├── DoubaoService.java
│   │   ├── ExpressionRecognitionService.java
│   │   ├── MoodService.java
│   │   └── OssService.java
│   └── FaceRecognitionApplication.java  # 应用入口
├── src/main/resources/
│   ├── mapper/          # MyBatis XML映射文件
│   ├── sql/             # SQL脚本
│   └── application.yml  # 应用配置
└── pom.xml              # Maven依赖配置
```

## 4. 核心功能设计

### 4.1 日记记录与管理
- 富文本日记编辑
- 图片上传与管理
- 心情标签分类
- 日记隐私设置
- 日记搜索与筛选

### 4.2 社交分享功能
- 日记分享圈
- 他人日记浏览
- 评论与点赞
- 关注与粉丝功能

### 4.3 AI智能聊天
- 基于豆包API的智能对话
- 心灵慰藉主题优化
- 聊天历史记录
- 个性化回复

### 4.4 数据可视化
- 心情占比饼状图
- 心情数量条形图
- 日记统计分析

### 4.5 个人中心
- 个人信息管理
- 主页背景自定义
- 头像与签名设置

### 4.6 心影绘梦
- 上传照片生成卡通风格图像
- 根据心情生成不同风格的图像

## 5. 数据库设计

### 5.1 核心表结构

#### 5.1.1 diary（日记表）
- id: 日记ID
- user_id: 用户ID
- title: 日记标题
- content: 日记内容
- selected_date: 日记日期
- created_at: 创建时间
- updated_at: 更新时间
- privacy: 隐私设置

#### 5.1.2 diary_mood（日记心情表）
- id: 记录ID
- diary_id: 日记ID
- mood_id: 心情ID
- intensity: 心情强度
- created_at: 创建时间

#### 5.1.3 diary_photo（日记照片表）
- id: 照片ID
- diary_id: 日记ID
- photo_url: 照片URL
- description: 照片描述
- created_at: 创建时间

#### 5.1.4 mood（心情表）
- id: 心情ID
- mood_name: 心情名称
- mood_color: 心情颜色
- created_at: 创建时间

#### 5.1.5 analysis_record（分析记录表）
- id: 记录ID
- user_id: 用户ID
- image_url: 图片URL
- main_emotion: 主要表情
- emotion_details: 表情详情
- created_at: 创建时间

## 6. API接口设计

### 6.1 日记相关接口
- POST /api/diary: 创建日记
- GET /api/diary: 获取日记列表
- GET /api/diary/{id}: 获取日记详情
- PUT /api/diary/{id}: 更新日记
- DELETE /api/diary/{id}: 删除日记

### 6.2 心情相关接口
- GET /api/mood: 获取心情列表
- POST /api/diary/{id}/mood: 添加日记心情

### 6.3 照片相关接口
- POST /api/oss/upload: 上传照片
- GET /api/diary/{id}/photos: 获取日记照片

### 6.4 AI聊天接口
- POST /api/doubao/chat: 智能聊天
- GET /api/doubao/history: 获取聊天历史

### 6.5 分析相关接口
- POST /api/analysis: 表情分析
- GET /api/analysis: 获取分析记录

## 7. 部署架构

### 7.1 开发环境
- 前端：Vite开发服务器（http://localhost:5173）
- 后端：Spring Boot内置Tomcat（http://localhost:8080）
- 数据库：本地MySQL

### 7.2 生产环境
- 前端：Nginx静态资源服务器
- 后端：Tomcat应用服务器
- 数据库：MySQL集群
- 存储：阿里云OSS
- 缓存：Redis

## 8. 技术亮点

### 8.1 AI智能聊天
- 集成豆包API，提供智能对话功能
- 针对心灵慰藉场景优化，提供温暖、专业的回复
- 支持多轮对话和上下文理解

### 8.2 数据可视化
- 使用ECharts实现精美的心情统计图表
- 实时更新，直观展示用户情感变化
- 响应式设计，适配不同设备

### 8.3 云存储服务
- 集成阿里云OSS，安全可靠的图片存储
- 支持大规模图片上传和快速访问
- 自动备份，防止数据丢失

### 8.4 模块化设计
- 前后端分离架构，便于扩展和维护
- 清晰的代码结构，易于团队协作
- 丰富的API接口，支持二次开发

## 9. 未来规划

1. 增强AI情感分析能力，实现日记内容的智能分析
2. 开发移动端应用，支持跨平台使用
3. 引入专业心理咨询师对接服务
4. 增加社交功能，打造情感社区
5. 优化推荐算法，提供个性化内容推荐
6. 增强数据安全和隐私保护
7. 支持多语言版本

## 10. 项目愿景

心影日记致力于成为用户最信赖的情感记录和心灵陪伴平台，通过技术的力量，让每个人都能找到情感的出口，获得心灵的慰藉和成长。