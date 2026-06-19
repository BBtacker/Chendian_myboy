# AI腺样体面容检测系统 - 开发文档

## 目录

1. [项目概述](#1-项目概述)
2. [技术栈](#2-技术栈)
3. [项目结构](#3-项目结构)
4. [快速开始](#4-快速开始)
5. [后端开发指南](#5-后端开发指南)
6. [前端开发指南](#6-前端开发指南)
7. [API 接口文档](#7-api-接口文档)
8. [数据库设计](#8-数据库设计)
9. [配置说明](#9-配置说明)
10. [常见问题](#10-常见问题)

---

## 1. 项目概述

### 1.1 项目简介

AI腺样体面容检测系统是一个基于AI的医疗辅助诊断系统，主要功能包括：

- 面部图像上传与AI分析
- 腺样体面容程度判断（正常/轻度/中期/重度）
- 检测记录管理
- 统计分析与可视化
- AI智能助手对话
- 个人中心管理

### 1.2 核心特性

- 基于豆包大模型的图像识别
- 阿里云OSS图片存储
- JWT用户认证
- 响应式Web界面
- PDF报告导出
- Excel数据导出

---

## 2. 技术栈

### 2.1 后端技术

| 技术 | 版本 | 用途 |
|------|------|------|
| Spring Boot | 3.5.5 | 后端框架 |
| MyBatis | 3.0.3 | ORM框架 |
| MySQL | 8.0+ | 数据库 |
| PageHelper | 1.4.6 | 分页插件 |
| JWT | 0.9.1 | 用户认证 |
| 阿里云OSS SDK | 3.17.4 | 图片存储 |
| Apache POI | 5.2.3 | Excel导出 |
| iText 7 | 7.2.5 | PDF生成 |
| Lombok | - | 简化代码 |

### 2.2 前端技术

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.5.18 | 前端框架 |
| Vite | 7.0.6 | 构建工具 |
| Element Plus | 2.11.2 | UI组件库 |
| Axios | 1.11.0 | HTTP客户端 |
| ECharts | 6.0.0 | 数据可视化 |

### 2.3 开发环境

- Java: JDK 17
- Node.js: ^20.19.0 || &gt;=22.12.0
- Maven: 3.6+
- MySQL: 8.0+

---

## 3. 项目结构

```
faceTest/
├── front/                          # 前端项目
│   ├── public/                     # 静态资源
│   ├── src/
│   │   ├── assets/                 # 资源文件
│   │   ├── components/             # Vue组件
│   │   │   ├── FaceAnalysis.vue    # 面容分析页面
│   │   │   ├── TestResult.vue      # 检测记录页面
│   │   │   ├── Statistics.vue      # 统计分析页面
│   │   │   ├── ChatAssistant.vue   # AI助手页面
│   │   │   ├── UserProfile.vue     # 个人中心
│   │   │   ├── Login.vue           # 登录页
│   │   │   └── Register.vue        # 注册页
│   │   ├── utils/
│   │   │   └── request.js          # Axios配置
│   │   ├── App.vue                 # 根组件
│   │   └── main.js                 # 入口文件
│   ├── vite.config.js              # Vite配置
│   └── package.json                # 依赖配置
│
├── src/                            # 后端项目
│   ├── main/
│   │   ├── java/com/long67/facetest/
│   │   │   ├── controller/         # 控制器层
│   │   │   │   ├── DoubaoController.java
│   │   │   │   ├── TestResultController.java
│   │   │   │   ├── StatisticsController.java
│   │   │   │   ├── ConversationController.java
│   │   │   │   ├── UserController.java
│   │   │   │   └── LoginController.java
│   │   │   ├── service/            # 服务层
│   │   │   │   ├── Impl/           # 服务实现
│   │   │   │   ├── DoubaoService.java
│   │   │   │   ├── TestResultService.java
│   │   │   │   ├── StatisticsService.java
│   │   │   │   ├── ConversationService.java
│   │   │   │   └── UserService.java
│   │   │   ├── mapper/             # 数据访问层
│   │   │   ├── pojo/               # 实体类
│   │   │   │   ├── testResult.java
│   │   │   │   ├── User.java
│   │   │   │   ├── Conversation.java
│   │   │   │   ├── Message.java
│   │   │   │   ├── Result.java
│   │   │   │   └── PageResult.java
│   │   │   ├── utils/              # 工具类
│   │   │   │   ├── AliyunOSSOperator.java
│   │   │   │   ├── JwtUtils.java
│   │   │   │   └── UserThreadLocal.java
│   │   │   ├── config/             # 配置类
│   │   │   │   ├── WebConfig.java
│   │   │   │   └── LoginInterceptor.java
│   │   │   └── FaceTestApplication.java
│   │   └── resources/
│   │       ├── mapper/              # MyBatis XML映射
│   │       └── application.yml      # 应用配置
│   └── test/                        # 测试代码
│
├── pom.xml                          # Maven配置
└── UPDATE_USER_TABLE.sql            # 数据库脚本
```

---

## 4. 快速开始

### 4.1 环境准备

1. **安装JDK 17**
2. **安装Node.js** (v20.19.0或更高)
3. **安装MySQL 8.0+**
4. **安装Maven 3.6+**

### 4.2 数据库配置

1. 创建数据库：
```sql
CREATE DATABASE face DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 执行数据库脚本（需根据实际表结构补充）

### 4.3 后端启动

1. 修改配置文件 `src/main/resources/application.yml`：
```yaml
spring:
  datasource:
    url: jdbc:mysql://localhost:3306/face
    username: root
    password: your_password
```

2. 配置阿里云OSS和豆包API密钥

3. 启动后端：
```bash
mvn spring-boot:run
```

后端服务将在 `http://localhost:8080` 启动

### 4.4 前端启动

1. 安装依赖：
```bash
cd front
npm install
```

2. 启动开发服务器：
```bash
npm run dev
```

前端服务将在 `http://localhost:5174` 启动

### 4.5 构建生产版本

**后端构建：**
```bash
mvn clean package
```

**前端构建：**
```bash
cd front
npm run build
```

---

## 5. 后端开发指南

### 5.1 统一返回结果

所有API接口统一使用 `Result` 类封装返回结果：

```java
public class Result&lt;T&gt; {
    private Integer code;  // 1成功，0失败
    private String msg;     // 提示信息
    private T data;         // 返回数据
    
    public static &lt;T&gt; Result&lt;T&gt; success(T data) {...}
    public static Result&lt;?&gt; error(String msg) {...}
}
```

### 5.2 分页查询

使用 PageHelper 实现分页：

```java
PageHelper.startPage(page, pageSize);
List&lt;testResult&gt; list = testResultMapper.list(query);
Page&lt;testResult&gt; p = (Page&lt;testResult&gt;) list;
return new PageResult&lt;&gt;(p.getTotal(), p.getResult());
```

### 5.3 用户认证

使用 JWT + 拦截器实现认证：

1. 登录成功后生成 Token
2. 前端请求时在 Header 中携带 `Authorization`
3. `LoginInterceptor` 拦截器验证 Token
4. 通过 `UserThreadLocal` 获取当前用户信息

### 5.4 文件上传

支持最大 10MB 的图片上传，文件存储在阿里云 OSS：

```java
// 上传文件
String url = aliyunOSSOperator.upload(file);
```

---

## 6. 前端开发指南

### 6.1 组件结构

主要页面组件：

| 组件 | 功能 |
|------|------|
| FaceAnalysis | 面容分析（上传图片、AI分析、结果展示） |
| TestResult | 检测记录管理（查看、编辑、删除、导出） |
| Statistics | 统计分析（图表展示） |
| ChatAssistant | AI智能助手 |
| UserProfile | 个人中心 |

### 6.2 HTTP 请求

使用封装好的 `request.js`：

```javascript
import request from '@/utils/request'

// 普通请求（30秒超时）
request.post('/doubao/analyzeFace', formData)

// 面容分析请求（120秒超时）
import { faceAnalysisRequest } from '@/utils/request'
faceAnalysisRequest.post('/doubao/analyzeFace', formData)
```

### 6.3 路由与状态管理

本项目使用简单的状态切换而非 Vue Router，通过 `activePage` 控制页面显示：

```javascript
const activePage = ref('face-analysis')  // 当前页面
```

### 6.4 组件缓存

面容分析组件使用 `&lt;keep-alive&gt;` 缓存，切换页面后状态保留：

```vue
&lt;keep-alive include="FaceAnalysis"&gt;
  &lt;component :is="currentComponent" :key="activePage" /&gt;
&lt;/keep-alive&gt;
```

---

## 7. API 接口文档

### 7.1 用户认证

#### 登录
```
POST /login
Content-Type: application/json

Request:
{
  "username": "string",
  "password": "string"
}

Response:
{
  "code": 1,
  "msg": "登录成功",
  "data": "jwt_token"
}
```

#### 注册
```
POST /register
Content-Type: application/json

Request:
{
  "username": "string",
  "password": "string",
  "name": "string"
}

Response:
{
  "code": 1,
  "msg": "注册成功"
}
```

### 7.2 面容分析

#### 上传并分析图片
```
POST /doubao/analyzeFace
Content-Type: multipart/form-data
Authorization: &lt;token&gt;

Request:
image: &lt;file&gt;

Response:
{
  "code": 1,
  "data": {
    "id": 1,
    "imagePath": "https://...",
    "isGlandFace": true,
    "level": "中期",
    "confidence": 85.5,
    "visualizationDescription": "...",
    "createTime": "2024-01-01 12:00:00"
  }
}
```

### 7.3 检测记录

#### 获取检测记录列表
```
GET /testResult/list?page=1&amp;pageSize=10&amp;level=&amp;startTime=&amp;endTime=
Authorization: &lt;token&gt;

Response:
{
  "code": 1,
  "data": {
    "total": 100,
    "records": [...]
  }
}
```

#### 更新检测记录
```
PUT /testResult
Content-Type: application/json
Authorization: &lt;token&gt;

Request:
{
  "id": 1,
  "level": "轻度",
  "visualizationDescription": "..."
}

Response:
{
  "code": 1,
  "msg": "更新成功"
}
```

#### 删除检测记录
```
DELETE /testResult/{id}
Authorization: &lt;token&gt;

Response:
{
  "code": 1,
  "msg": "删除成功"
}
```

#### 导出Excel
```
GET /testResult/export?level=&amp;startTime=&amp;endTime=
Authorization: &lt;token&gt;

Response: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
```

#### 导出PDF
```
GET /testResult/exportPdf/{id}
Authorization: &lt;token&gt;

Response: application/pdf
```

### 7.4 统计分析

#### 获取统计数据
```
GET /statistics/data
Authorization: &lt;token&gt;

Response:
{
  "code": 1,
  "data": {
    "totalCount": 100,
    "levelDistribution": {...},
    "trendData": [...]
  }
}
```

---

## 8. 数据库设计

### 8.1 user 表（用户表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| username | VARCHAR(50) | 用户名（唯一） |
| password | VARCHAR(100) | 密码（加密） |
| name | VARCHAR(50) | 姓名 |
| avatar | VARCHAR(255) | 头像URL |
| create_time | DATETIME | 创建时间 |

### 8.2 test_result 表（检测记录表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| user_id | INT | 用户ID |
| image_path | VARCHAR(255) | 图片URL |
| is_gland_face | BOOLEAN | 是否为腺样体面容 |
| level | VARCHAR(20) | 程度（正常/轻度/中期/重度） |
| confidence | DOUBLE | 置信度 |
| visualization_description | TEXT | 分析描述 |
| test_time | DATETIME | 检测时间 |

### 8.3 conversation 表（对话表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| user_id | INT | 用户ID |
| title | VARCHAR(100) | 对话标题 |
| create_time | DATETIME | 创建时间 |

### 8.4 message 表（消息表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| conversation_id | INT | 对话ID |
| role | VARCHAR(20) | 角色（user/assistant） |
| content | TEXT | 消息内容 |
| create_time | DATETIME | 创建时间 |

---

## 9. 配置说明

### 9.1 后端配置 (application.yml)

```yaml
spring:
  datasource:
    url: jdbc:mysql://localhost:3306/face
    username: root
    password: your_password
  servlet:
    multipart:
      max-file-size: 10MB      # 单文件大小限制
      max-request-size: 10MB    # 请求大小限制

aliyun:
  oss:
    endpoint: https://oss-cn-beijing.aliyuncs.com
    bucketName: your-bucket-name
    accessKeyId: your-access-key-id
    accessKeySecret: your-access-key-secret

doubao:
  api:
    key: your-doubao-api-key
  model:
    name: doubao-seed-1-8-251228
```

### 9.2 前端配置 (vite.config.js)

```javascript
server: {
  port: 5174,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
      rewrite: (path) =&gt; path.replace(/^\/api/, '')
    }
  }
}
```

---

## 10. 常见问题

### 10.1 面容分析超时

**问题**: 分析请求30秒后超时

**解决方案**: 使用 `faceAnalysisRequest` 替代 `request`，超时时间为120秒

### 10.2 切换页面后数据丢失

**问题**: 面容分析页面切换后数据重置

**解决方案**: 使用 `&lt;keep-alive include="FaceAnalysis"&gt;` 缓存组件状态

### 10.3 阿里云OSS上传失败

**问题**: 图片上传失败

**解决方案**: 
1. 检查 `application.yml` 中的 OSS 配置
2. 确认 Bucket 访问权限设置为公共读或私有（私有需要签名URL）
3. 检查网络连接

### 10.4 数据库连接失败

**问题**: 后端启动时报数据库连接错误

**解决方案**:
1. 确认 MySQL 服务已启动
2. 检查 `application.yml` 中的数据库连接配置
3. 确认数据库用户名和密码正确

### 10.5 端口被占用

**问题**: 
- 后端：Port 8080 was already in use
- 前端：Port 5174 is in use

**解决方案**:
- 修改 `application.yml` 中的 `server.port`
- 修改 `vite.config.js` 中的 `server.port`

---

## 附录

### 开发规范

1. **后端命名规范**
   - Controller: `XxxController`
   - Service: `XxxService` / `XxxServiceImpl`
   - Mapper: `XxxMapper`
   - POJO: 驼峰命名，数据库字段下划线转驼峰

2. **前端命名规范**
   - 组件文件: 大驼峰 `FaceAnalysis.vue`
   - 组件名: `name: 'FaceAnalysis'`
   - 变量/函数: 小驼峰 `handleLoginSuccess`

3. **Git 提交规范**
   - feat: 新功能
   - fix: 修复bug
   - docs: 文档更新
   - style: 格式调整
   - refactor: 重构

---


