# AI腺样体面容检测系统 - 开发文档

## 目录

1. [项目概述](#1-项目概述)
2. [技术栈](#2-技术栈)
3. [微服务架构](#3-微服务架构)
4. [快速开始](#4-快速开始)
5. [后端开发指南](#5-后端开发指南)
6. [前端开发指南](#6-前端开发指南)
7. [API 接口文档](#7-api-接口文档)
8. [数据库设计](#8-数据库设计)
9. [核心架构设计](#9-核心架构设计)
10. [配置说明](#10-配置说明)
11. [Docker 部署](#11-docker-部署)
12. [常见问题](#12-常见问题)

---

## 1. 项目概述

### 1.1 项目简介

AI腺样体面容检测系统是基于 go-zero 微服务框架构建的医疗辅助诊断系统，主要功能包括：

- 面部图像上传与AI分析（异步诊断流水线）
- 腺样体面容程度判断（正常/轻度/中度/重度）
- RAG 检索增强（Milvus 向量数据库匹配相似病例）
- 检测记录管理（CRUD + 批量操作）
- 统计分析与可视化
- AI智能助手对话（DeepSeek）
- 报告导出（PDF/Excel）

### 1.2 核心特性

- **微服务架构** - 4个核心微服务 + API网关，gRPC 通信，etcd 服务发现
- **异步诊断** - RabbitMQ 消息队列，Outbox Pattern 保证原子性
- **三级缓存** - 本地进程缓存 + Redis Hash + MySQL，P99 延迟 45ms
- **限流防护** - Redis Lua 脚本滑动窗口限流
- **幂等消费** - 唯一幂等键 + 死信队列 + 指数退避重试
- **乐观锁** - GORM version 字段保证并发安全
- **容器化** - Docker Compose 一键部署全部基础设施和服务

---

## 2. 技术栈

### 2.1 后端技术

| 技术 | 版本 | 用途 |
|------|------|------|
| go-zero | v1.7.6 | 微服务框架（zrpc + rest） |
| GORM | v1.25.12 | ORM 框架（乐观锁 version） |
| gRPC + Protobuf | v1.66.2 / v1.34.2 | 微服务通信 |
| etcd | v3.5.15 | 服务注册与发现 |
| RabbitMQ | 3.13 | 消息队列（死信队列 + 重试） |
| Milvus | v2.4.0 | 向量数据库（RAG 检索） |
| Redis | 7 | 三级缓存 + Lua 限流 |
| MySQL | 8.0 | 持久化存储 |
| DeepSeek API | - | 大模型推理 |
| go-redis | v9.7.0 | Redis 客户端 |
| amqp091-go | v1.10.0 | RabbitMQ 客户端 |
| excelize | v2.8.1 | Excel 导出 |
| bcrypt | - | 密码加密 |

### 2.2 前端技术

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.5.18 | 前端框架 |
| Vite | 7.0.6 | 构建工具 |
| Element Plus | 2.11.2 | UI组件库 |
| Axios | 1.11.0 | HTTP客户端 |
| ECharts | 6.0.0 | 数据可视化 |

### 2.3 开发环境

- Go: 1.22+
- Node.js: ^20.19.0 || >=22.12.0
- Docker Desktop（WSL 2 集成）
- protoc + goctl（proto 代码生成）

---

## 3. 微服务架构

### 3.1 服务拆分

| 服务 | 端口 | 类型 | 职责 |
|------|------|------|------|
| API 网关 | 8080 | HTTP | 路由转发、JWT认证、CORS、静态文件 |
| 认证服务 | 8081 | gRPC | 登录/注册/用户管理/Token验证/头像上传 |
| 上传服务 | 8082 | gRPC | 图片上传/批量上传（本地存储） |
| 诊断服务 | 8083 | gRPC | 异步诊断/结果CRUD/统计/RabbitMQ消费 |
| 报告服务 | 8084 | gRPC | Excel导出/PDF报告生成 |

### 3.2 服务间通信

- **外部请求** → API 网关（HTTP REST）
- **网关 → 微服务** → gRPC（通过 etcd 服务发现）
- **异步诊断** → RabbitMQ 消息队列
- **特征检索** → Milvus 向量数据库
- **缓存** → Redis（三级缓存中间层）

### 3.3 项目结构

```
backend/
├── gateway/                    # API 网关 (HTTP :8080)
│   ├── gateway.go              # 主入口
│   ├── etc/gateway.yaml        # 配置文件
│   └── internal/
│       ├── config/config.go    # 配置结构体
│       ├── svc/service_context.go  # 服务上下文
│       └── handler/            # HTTP 处理器
│           ├── routes.go       # 路由注册
│           ├── auth_handler.go
│           ├── diagnosis_handler.go
│           └── conversation_handler.go
│
├── auth/                       # 认证服务 (gRPC :8081)
│   ├── auth.go
│   ├── etc/auth.yaml
│   └── internal/
│       ├── config/
│       ├── svc/
│       ├── logic/              # 业务逻辑
│       │   ├── login_logic.go
│       │   ├── register_logic.go
│       │   ├── user_logic.go
│       │   └── upload_avatar_logic.go
│       └── server/
│
├── upload/                     # 上传服务 (gRPC :8082)
├── diagnosis/                  # 诊断服务 (gRPC :8083)
│   └── internal/
│       ├── logic/              # 提交/查询/删除/统计
│       ├── consumer/           # RabbitMQ 消费者
│       │   └── diagnosis_consumer.go  # 核心异步诊断流水线
│       └── svc/
│
├── report/                     # 报告服务 (gRPC :8084)
│
├── common/                     # 公共模块
│   ├── model/models.go         # GORM 模型（7张表）
│   ├── middleware/auth.go      # JWT 认证中间件
│   └── pkg/                    # 工具包
│       ├── response.go         # 统一响应 {code, msg, data}
│       ├── jwt.go              # JWT 生成/解析
│       ├── redis_cache.go      # Redis 缓存封装
│       ├── rate_limiter.go     # Lua 滑动窗口限流
│       ├── rabbitmq.go         # RabbitMQ 封装
│       ├── deepseek.go         # DeepSeek API 客户端
│       ├── milvus.go           # Milvus 向量数据库
│       ├── storage.go          # 本地文件存储
│       └── three_level_cache.go # 三级缓存
│
├── proto/                      # Protobuf 定义
│   ├── auth.proto
│   ├── upload.proto
│   ├── diagnosis.proto
│   └── report.proto
│
├── deploy/
│   ├── docker-compose.yml      # 基础设施 + 微服务编排
│   ├── dockerfiles/            # 5个 Dockerfile
│   ├── configs/                # Docker 环境配置
│   ├── mysql/init.sql          # 数据库初始化（7张表）
│   └── setup-wsl.sh            # WSL 环境搭建脚本
│
├── go.mod
└── Makefile
```

---

## 4. 快速开始

### 4.1 环境准备（WSL）

```bash
# 进入项目目录
cd /mnt/c/Users/26248/Desktop/Code/faceTest/backend

# 一键搭建环境（安装 Go、protoc、goctl，生成代码，下载依赖，编译）
bash deploy/setup-wsl.sh
```

### 4.2 启动基础设施

```bash
cd deploy
docker-compose up -d
```

等待所有服务健康：
```bash
docker-compose ps
```

### 4.3 启动微服务

```bash
cd /mnt/c/Users/26248/Desktop/Code/faceTest/backend

# 一键启动所有微服务（推荐）
make dev          # 后台启动所有服务，显示状态面板
make dev-status   # 查看服务运行状态
make dev-logs     # 实时查看所有服务日志
make dev-stop     # 停止所有服务
make dev-restart  # 重启所有服务
```

`make dev` 启动后会显示状态面板：
```
  服务          端口       进程状态     端口监听      PID
  ------------------------------------------------------------
  auth          8081       运行中       ✓ 监听       12345
  upload        8082       运行中       ✓ 监听       12346
  diagnosis     8083       运行中       ✓ 监听       12347
  report        8084       运行中       ✓ 监听       12348
  gateway       8080       运行中       ✓ 监听       12349

  ✓ 所有服务运行正常
```

如果某个服务启动失败，面板会显示 `✗ 未监听`，并自动输出该服务最后 5 行日志。

或 Docker 一键启动：
```bash
cd deploy
docker-compose up -d --build
```

### 4.4 启动前端

```bash
cd front
npm install
npm run dev
```

前端运行在 `http://localhost:5174`。

### 4.5 默认账号

- 用户名：`admin`
- 密码：`admin123`

---

## 5. 后端开发指南

### 5.1 统一返回结果

所有API接口统一使用 `Response` 结构：

```go
type Response struct {
    Code int         `json:"code"` // 1-成功 0-失败
    Msg  string      `json:"msg"`
    Data interface{} `json:"data,omitempty"`
}
```

### 5.2 GORM 模型

所有模型包含 `version` 乐观锁字段：

```go
type User struct {
    ID        uint64    `gorm:"primaryKey"`
    Username  string    `gorm:"uniqueIndex"`
    Password  string
    Name      string
    Version   int       `gorm:"default:0"` // 乐观锁
    CreateTime time.Time
    UpdateTime time.Time
}
```

乐观锁更新：
```go
result := db.Model(&user).Where("id = ? AND version = ?", id, user.Version).
    Updates(map[string]interface{}{
        "name":    newName,
        "version": gorm.Expr("version + 1"),
    })
```

### 5.3 JWT 认证

```go
// 生成 Token
token, _ := pkg.GenerateToken(userID, username, name, role)

// 解析 Token
claims, _ := pkg.ParseToken(token)
// claims.UserID, claims.Username, claims.Name, claims.Role
```

中间件自动解析 Token 并将用户信息存入 Header：
```go
r.Header.Set("X-User-Id", fmt.Sprintf("%d", claims.UserID))
// handler 中获取：
userID := middleware.GetUserID(r)
```

### 5.4 新增微服务

1. 定义 proto 文件 `proto/xxx.proto`
2. 生成 gRPC 代码：`make proto`
3. 创建服务目录 `xxx/`
4. 实现 `internal/logic/` 业务逻辑
5. 在网关 `routes.go` 中注册 HTTP 路由
6. 在 `docker-compose.yml` 中添加服务

---

## 6. 前端开发指南

### 6.1 HTTP 请求

前端保持 Vue 3 不变，通过 `/api` 前缀代理到网关：

```javascript
import request from '@/utils/request'

// 普通请求（30秒超时）
request.post('/login', { username, password })

// 面容分析请求（120秒超时）
import { faceAnalysisRequest } from '@/utils/request'
faceAnalysisRequest.post('/doubao/analyzeFace', formData)
```

### 6.2 Vite 代理

```javascript
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true,
    rewrite: (path) => path.replace(/^\/api/, '')
  }
}
```

---

## 7. API 接口文档

### 7.1 用户认证

| 接口 | 方法 | 认证 | 说明 |
|------|------|------|------|
| `/login` | POST | 否 | 用户登录 |
| `/register` | POST | 否 | 用户注册 |
| `/logout` | POST | 是 | 登出 |
| `/user` | GET | 是 | 获取用户信息 |
| `/user` | PUT | 是 | 更新用户信息 |
| `/user/avatar` | POST | 是 | 上传头像 |

### 7.2 面容分析

| 接口 | 方法 | 说明 |
|------|------|------|
| `/doubao/analyzeFace` | POST | 上传图片 + 异步诊断（网关轮询等待结果，最多90秒） |

### 7.3 检测记录

| 接口 | 方法 | 说明 |
|------|------|------|
| `/testResult/result` | GET | 获取列表（支持分页/筛选） |
| `/testResult/{id}` | DELETE | 删除单条 |
| `/testResult/batch` | DELETE | 批量删除 |
| `/testResult/update` | PUT | 更新记录 |
| `/testResult/download` | GET | 导出 Excel |
| `/testResult/download/pdf` | GET | 导出 PDF |

### 7.4 统计分析

| 接口 | 方法 | 说明 |
|------|------|------|
| `/statistics/overview` | GET | 统计概览 |
| `/statistics/detail` | GET | 统计详情 |

### 7.5 AI 对话

| 接口 | 方法 | 说明 |
|------|------|------|
| `/conversation/create` | POST | 创建对话 |
| `/conversation/send` | POST | 发送消息（DeepSeek） |
| `/conversation/list` | GET | 对话列表 |
| `/conversation/messages` | GET | 消息列表 |

---

## 8. 数据库设计

### 8.1 表结构概览

| 表名 | 说明 | 关键字段 |
|------|------|----------|
| `user` | 用户表 | id, username, password(bcrypt), role, version |
| `diagnosis_task` | 诊断任务表 | task_no(幂等键), status, retry_count, version |
| `diagnosis_result` | 诊断结果表 | is_gland_face, level(轻度/中度/重度/非腺样体面容), feature_vector_id |
| `outbox_message` | 本地消息表 | aggregate_id, status, retry_count, next_retry_time |
| `conversation` | 对话表 | user_id, title |
| `message` | 消息表 | conversation_id, role, content |
| `diagnosis_cache` | 缓存表 | cache_key, cache_value, expire_time |

### 8.2 任务状态流转

```
0-待处理 → 1-处理中 → 2-已完成
                ↓
           3-失败（重试 ≤ 5次）
```

### 8.3 Outbox 消息状态

```
0-待发送 → 1-已发送 → 2-已确认
                ↓
           3-失败（重试 ≤ 5次）
```

---

## 9. 核心架构设计

### 9.1 异步诊断流水线

```
[网关] 上传图片 → 提交诊断任务
         │
    [诊断服务] 事务内：
         ├── 创建 diagnosis_task 记录
         └── 创建 outbox_message 记录
         │
    [Outbox 中继器] 每2秒轮询：
         ├── 查询 status=0 的 outbox_message
         ├── 投递到 RabbitMQ
         └── 更新 status=1
         │
    [MQ 消费者] 幂等消费：
         ├── 检查 task_no 幂等键
         ├── 乐观锁更新 task status=1
         ├── EfficientNet 特征提取（模拟）
         ├── Milvus RAG 检索相似病例
         ├── DeepSeek 生成诊断报告
         ├── 保存 diagnosis_result
         ├── Milvus 插入特征向量
         ├── 三级缓存更新
         └── 乐观锁更新 task status=2
         │
    [失败处理]
         ├── 指数退避重试（最多5次）
         └── 超过重试次数 → 死信队列
```

### 9.2 三级缓存

```
请求 → 本地进程缓存 (sync.Map, TTL 5min)
         ↓ miss
      Redis Hash (特征向量索引, TTL 1h)
         ↓ miss
      MySQL (diagnosis_cache 表)
         ↓ 回填缓存
```

### 9.3 Redis Lua 限流

```lua
-- 滑动窗口限流：60秒内最多5次
local count = redis.call('ZCOUNT', key, now - window, now)
if count >= limit then return 0 end
redis.call('ZADD', key, now, requestId)
redis.call('EXPIRE', key, window)
return 1
```

---

## 10. 配置说明

### 10.1 网关配置 (gateway/etc/gateway.yaml)

```yaml
Name: gateway-service
Host: 0.0.0.0
Port: 8080

AuthRpc:
  Etcd:
    Hosts: [127.0.0.1:2379]
    Key: auth.rpc

DiagnosisRpc:
  Etcd:
    Hosts: [127.0.0.1:2379]
    Key: diagnosis.rpc

MySQL:
  DataSource: "root:your-mysql-password@tcp(127.0.0.1:3307)/face?charset=utf8mb4&parseTime=True&loc=Local"

Redis:
  Host: 127.0.0.1:6380

DeepSeek:
  APIKey: "sk-xxx"
  Model: "deepseek-chat"
```

### 10.2 诊断服务配置 (diagnosis/etc/diagnosis.yaml)

```yaml
RabbitMQ:
  URL: "amqp://admin:your-rabbitmq-password@127.0.0.1:5672/"
  Exchange: "diagnosis.exchange"
  Queue: "diagnosis.queue"
  DLXExchange: "diagnosis.dlx.exchange"
  DLXQueue: "diagnosis.dlx.queue"

Milvus:
  Address: "127.0.0.1:19530"

DeepSeek:
  APIKey: "sk-xxx"
  Model: "deepseek-chat"
```

---

## 11. Docker 部署

### 11.1 一键部署

```bash
cd backend/deploy
docker-compose up -d --build
```

### 11.2 端口映射

| 服务 | 容器端口 | 主机端口 |
|------|----------|----------|
| API 网关 | 8080 | 8080 |
| 认证服务 | 8081 | 8081 |
| 上传服务 | 8082 | 8082 |
| 诊断服务 | 8083 | 8083 |
| 报告服务 | 8084 | 8084 |
| MySQL | 3306 | 3307 |
| Redis | 6379 | 6380 |
| etcd | 2379 | 2379 |
| RabbitMQ | 5672 | 5672 |
| RabbitMQ 管理界面 | 15672 | 15672 |
| Milvus | 19530 | 19530 |
| MinIO | 9000 | 9000 |

### 11.3 Docker 配置文件

Docker 环境使用 `deploy/configs/` 下的专用配置文件，将 `127.0.0.1` 替换为容器服务名（如 `mysql`、`redis`、`etcd`）。

---

## 12. 常见问题

### 12.1 proto 代码未生成

```bash
# 在 WSL 中执行
cd backend
make proto
# 或
protoc --go_out=./proto --go_opt=paths=source_relative \
    --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
    proto/auth.proto proto/upload.proto proto/diagnosis.proto proto/report.proto
```

### 12.2 go.sum 不存在

```bash
cd backend
go mod tidy
```

### 12.3 Docker 构建失败

确保 Docker Desktop 已开启 WSL 2 集成，并在 WSL 中执行 `docker-compose` 命令。

### 12.4 面容分析超时

网关默认等待 90 秒。如果 DeepSeek API 响应慢，可调整 `gateway/internal/handler/diagnosis_handler.go` 中的超时时间。

### 12.5 RabbitMQ 消费者不工作

1. 检查 RabbitMQ 管理界面 `http://localhost:15672`（admin/your-rabbitmq-password）
2. 确认队列 `diagnosis.queue` 和死信队列 `diagnosis.dlx.queue` 已创建
3. 检查诊断服务日志

### 12.6 Milvus 连接失败

Milvus 启动较慢（约30秒），请等待健康检查通过后再启动诊断服务。
