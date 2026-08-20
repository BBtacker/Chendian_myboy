# AI 腺样体面容智能筛查平台

<div align="center">

基于 **go-zero 微服务架构** + **Vue 3** 的儿童腺样体面容 AI 智能筛查系统。

异步诊断流水线 · RAG 检索增强 · 可靠消息投递（Outbox）· 多图综合分析

</div>

## 项目简介

为儿童腺样体面容早期筛查提供高隐私保护的诊断服务。用户上传（或拍照）多张面部照片，系统自动剔除无脸照片、综合多图特征，通过 **EfficientNet-B3 特征提取 → Milvus 向量 RAG 检索 → DeepSeek 大模型推理** 的完整流水线，输出腺样体面容与否的判定、严重程度分级、医学特征描述及相似参考病例，并支持 PDF / Excel 报告导出。

- 后端：**go-zero v1.7.6**（`rest` API 网关 + `zrpc` 微服务），gRPC + etcd 服务注册发现
- 前端：**Vue 3 + Element Plus + Vite**
- 模型：**EfficientNet-B3**（ONNX Runtime，独立 Python 微服务）+ **Milvus**（RAG）+ **DeepSeek**（`deepseek-chat`）

## 功能特性

- 🖼️ **多图面容分析** — 支持批量上传/拍照（最多 9 张），MediaPipe 人脸门控自动剔除无脸照片，跨图特征平均综合判断
- ⚡ **异步诊断流水线** — 上传后异步处理，网关轮询任务状态，AI 推理不阻塞 HTTP 请求
- 🧩 **RAG 检索增强** — EfficientNet-B3 提取 1536 维特征存入 Milvus，按相似度召回历史病例作为大模型参考
- 🔒 **可靠消息投递** — Outbox 本地消息表 + RabbitMQ 死信队列 + 指数退避重试，保证任务不丢、不重复消费
- 🔐 **乐观锁并发安全** — 关键表 version 字段 CAS 更新，保证任务状态与用户资料并发安全
- 🛡️ **Redis Lua 限流** — 滑动窗口原子限流，防恶意刷接口
- 📊 **检测记录与统计** — 历史记录分页/筛选，病情分布与等级统计可视化
- 💬 **AI 医疗对话** — 基于 DeepSeek 的耳鼻喉科智能助手
- 📄 **报告导出** — PDF 检测报告 + Excel 数据导出

## 系统架构

```
┌────────────────────────────────────────────────────────────────────┐
│                    前端 (Vue 3 + Element Plus, :5174)               │
│   登录/注册 │ 面容分析(上传/拍照) │ 历史记录 │ AI对话 │ 统计 │ 导出  │
└────────────────────────────────┬───────────────────────────────────┘
                                 │ HTTP + JWT (Vite 代理 /api → :8080)
┌────────────────────────────────▼───────────────────────────────────┐
│                     API 网关 gateway (:8080, go-zero rest)          │
│   JWT 认证中间件 │ 路由转发 │ CORS │ 静态文件 │ AI对话 │ 分析聚合轮询 │
└──────┬──────────────┬──────────────┬──────────────┬─────────────────┘
       │ gRPC         │ gRPC         │ gRPC         │ gRPC
       ▼              ▼              ▼              ▼
┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐
│ 认证服务    │  │ 上传服务    │  │ 诊断服务    │  │ 报告服务    │
│ auth :8081│  │upload:8082│  │diagnosis  │  │report:8084│
│ 登录/注册  │  │ 图片存储   │  │  :8083     │  │ Excel/PDF │
│ JWT签发    │  │           │  │ 异步诊断流水线│  │ 导出       │
└─────┬─────┘  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘
      │              │         ┌────▼────┐         │
      │              │         │RabbitMQ │  Outbox中继投递
      │              │         └────┬────┘         │
      │              │              │ 消费者        │
      │              │              ▼               │
      └──────┬───────┴──────┬───────┴───────┬───────┘
             │              │               │
        ┌────▼────┐   ┌─────▼─────┐   ┌─────▼─────┐
        │ MySQL   │   │  Redis    │   │  Milvus   │
        │ 7 张表   │   │ 缓存/限流  │   │ 向量 RAG  │
        └─────────┘   └───────────┘   └───────────┘
             │  etcd 服务注册与发现（所有微服务）
             └── 本地磁盘 uploads/（图片存储）
```

**AI 特征提取链路（独立 Python 微服务）**：

```
┌──────────────────────────────────────────────────────────────┐
│  feature-extractor (:8085, FastAPI + ONNX Runtime)            │
│  POST /extract ──► EfficientNet-B3 推理 ──► 1536 维特征向量    │
│                    MediaPipe Face Mesh ──► 人脸检测 + 几何测量  │
└──────────────────────────────────────────────────────────────┘
```

## 技术栈

### 后端微服务

| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.22 | 微服务开发语言 |
| go-zero | v1.7.6 | 微服务框架（`zrpc` gRPC + `rest` HTTP） |
| GORM | v1.25.12 | ORM，事务 / 乐观锁（version） |
| gRPC + Protobuf | v1.66.2 / v1.34.2 | 服务间通信，etcd 服务发现 |
| etcd | v3.5.15 | 服务注册与发现 |
| RabbitMQ | 3.13 | 异步诊断队列 + 死信队列（DLX） |
| Milvus | v2.4.0 | 向量数据库（RAG 检索） |
| Redis | 7 | 缓存 + Lua 滑动窗口限流 |
| MySQL | 8.0 | 用户/任务/结果/Outbox 等 7 张表 |
| DeepSeek API | `deepseek-chat` | 诊断推理 + AI 对话 |

### AI 与算法

| 技术 | 用途 |
|------|------|
| EfficientNet-B3 | 面部视觉特征提取（1536 维，ONNX Runtime 推理，独立 Python 微服务 :8085） |
| MediaPipe Face Mesh | 人脸检测门控 + 几何测量（口裂比、面宽高比、鼻唇角等） |
| Milvus RAG | 特征向量相似检索 Top-5 历史病例，作为大模型「参考答案」 |
| DeepSeek | 综合诊断推理（几何测量为首要硬指标 + B3 特征描述），输出结构化 JSON |

### 前端

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.5 | 前端框架 |
| Vite | 7 | 构建工具（`/api` 代理 → 网关 :8080） |
| Element Plus | 2.11 | UI 组件库 |
| ECharts | 6.0 | 数据可视化 |

## 核心设计

### 1. 异步诊断流水线

```
[前端] 上传图片
  └─ [网关] 逐张调上传服务存图 → 调诊断服务 SubmitDiagnosis（同步，毫秒级返回 task_id）
        └─ [诊断服务] 同一事务内双写 diagnosis_task + outbox_message（Outbox 模式）
              └─ [Outbox 中继] 每 2s 扫描投递到 RabbitMQ
                    └─ [消费者] 幂等检查 → 乐观锁抢占
                          ├─ ① EfficientNet-B3 特征提取（1536 维）＋ MediaPipe 人脸门控
                          ├─ ② 跨图特征平均（L2 归一化）+ 几何测量平均
                          ├─ ③ Milvus 相似病例检索（Top-5）
                          ├─ ④ DeepSeek 综合推理（几何测量 + 特征描述 + 参考病例）
                          ├─ ⑤ 结果写入 MySQL，特征向量写入 Milvus
                          └─ ⑥ 任务状态置为完成
  └─ [网关] 每 2s 轮询 GetTaskStatus（最长 90s）→ 完成后返回诊断结果给前端
```

### 2. Outbox 本地消息表

「数据库写 + 发消息」无法原子完成，会导致任务创建成功但消息丢失。方案：在**同一个 MySQL 事务**内写入任务记录和 `outbox_message` 记录，事务提交后由中继器扫描投递 RabbitMQ，失败按指数退避重试，保证**最终一致、绝不丢消息**。

### 3. 幂等 + 死信队列 + 指数退避

- **幂等键**：`task_no` 唯一索引，消费者处理前按任务状态判断是否重复
- **死信队列**：处理失败的消息 Nack 进入 DLX，由死信消费者兜底
- **指数退避**：重试间隔 `1 << retryCount` 秒（1s、2s、4s…），超 `max_retry`（默认 5）标记失败

### 4. 乐观锁

`user`、`diagnosis_task`、`diagnosis_result` 均含 `version` 字段，更新使用
`UPDATE ... SET version = version + 1 WHERE id = ? AND version = ?`，`RowsAffected == 0` 表示并发冲突，由业务层决定重试或放弃。

### 5. Redis Lua 滑动窗口限流

限流的「读窗口计数 → 判断 → 写回」三步必须原子，使用 Lua 脚本在 Redis 单线程内一次性完成，并对 Redis 故障做 fail-open（放行），避免限流组件成为新的单点。

### 6. 多图人脸门控

MediaPipe 检测不到人脸的图片自动剔除（记录到 `skipped_images` 并向前端说明），仅有效人脸图片参与跨图特征平均与综合判断；MediaPipe 不可用时宽松视为有人脸，避免误杀全部照片。

## 项目结构

```
faceTest/
├── backend/                          # go-zero 微服务后端
│   ├── gateway/                      # API 网关（HTTP :8080，go-zero rest）
│   ├── auth/                         # 认证服务（gRPC :8081）
│   ├── upload/                       # 图片上传服务（gRPC :8082）
│   ├── diagnosis/                    # 诊断服务（gRPC :8083，异步流水线核心）
│   │   └── internal/consumer/        #   RabbitMQ 消费者 + Outbox 中继 + 死信消费
│   ├── report/                       # 报告服务（gRPC :8084，Excel/PDF）
│   ├── feature_extractor/            # Python 特征提取服务（FastAPI :8085）
│   │   └── models/                   #   efficientnet_b3.onnx（见 README 获取）
│   ├── common/                       # 公共模块（model / middleware / pkg）
│   ├── proto/                        # Protobuf 定义（auth/upload/diagnosis/report）
│   ├── deploy/                       # Docker Compose + Dockerfile + 配置模板
│   ├── go.mod
│   └── Makefile                      # make dev 一键启动
└── front/                            # Vue 3 前端
    └── src/components/               # FaceAnalysis / TestResult / Statistics 等
```

## 快速开始

### 前置条件

- WSL 2（Ubuntu 22.04+）或 Linux/macOS
- Docker Desktop（WSL 集成）
- Go 1.22+、Node.js 20+、protoc + goctl

### 1. 配置模板

配置文件含密钥，不入库；首次使用先复制模板并填写自己的密钥：

```bash
cd backend

# 各服务配置（本地开发）
cp auth/etc/auth.yaml.example            auth/etc/auth.yaml
cp upload/etc/upload.yaml.example        upload/etc/upload.yaml
cp diagnosis/etc/diagnosis.yaml.example  diagnosis/etc/diagnosis.yaml
cp gateway/etc/gateway.yaml.example      gateway/etc/gateway.yaml
cp report/etc/report.yaml.example        report/etc/report.yaml

# Docker 部署配置
cp deploy/docker-compose.yml.example     deploy/docker-compose.yml
cp deploy/configs/*-docker.yaml.example  deploy/configs/
```

> ⚠️ `JWT.Secret`（auth 与 gateway）必须保持一致；DeepSeek `APIKey`、MySQL、RabbitMQ 密码需替换为你的真实值。

### 2. 启动基础设施（Docker）

```bash
cd backend/deploy
docker-compose up -d
```

| 服务 | 端口 | 说明 |
|------|------|------|
| MySQL | 3307 | 主数据库（库 `face`） |
| Redis | 6380 | 缓存 + 限流 |
| etcd | 2379 | 服务注册与发现 |
| RabbitMQ | 5672 / 15672 | 消息队列（管理界面 15672） |
| Milvus | 19530 | 向量数据库 |

### 3. 启动微服务

```bash
cd backend
make dev          # 后台启动 auth/upload/diagnosis/report/gateway
make dev-status   # 查看服务运行状态
make dev-logs     # 实时查看日志
make dev-stop     # 停止所有服务
```

或 Docker 一键构建全部微服务：

```bash
cd backend/deploy
docker-compose up -d --build
```

### 4. 特征模型

`feature_extractor/models/efficientnet_b3.onnx` 体积较大未入库，请按 `backend/feature_extractor/README.md` 说明放置模型后再启动特征服务（模型缺失时系统自动回退 Mock，不影响主流程运行）。

### 5. 启动前端

```bash
cd front
npm install
npm run dev
```

前端运行在 `http://localhost:5174`，API 自动代理到网关 `http://localhost:8080`。

## API 概览

所有接口经网关 `http://localhost:8080` 访问，前端通过 `/api` 前缀代理。

| 分类 | 接口 | 方法 | 说明 |
|------|------|------|------|
| 认证 | `/login` `/register` | POST | 登录 / 注册 |
| 用户 | `/user` | GET/PUT | 获取 / 更新用户信息 |
| 面容分析 | `/doubao/analyzeFace` | POST | 单图分析（异步，网关轮询 90s） |
| 面容分析 | `/doubao/analyzeMulti` | POST | 多图综合分析（人脸门控 + 特征平均） |
| 记录 | `/testResult/result` | GET | 历史记录（分页 / 筛选） |
| 记录 | `/testResult/{id}` | DELETE | 删除记录 |
| 导出 | `/testResult/download` | GET | 导出 Excel |
| 导出 | `/testResult/download/pdf` | GET | 导出 PDF 报告 |
| 统计 | `/statistics/overview` `/detail` | GET | 统计概览 / 详情 |
| 对话 | `/conversation/*` | — | AI 医疗助手 |

> 注：`/doubao/*` 为历史路径命名，实际处理走 go-zero 微服务流水线（EfficientNet-B3 → Milvus → DeepSeek），与豆包无关。

## 开源协议

本项目基于 [MIT License](LICENSE) 开源。
