# AI腺样体面容检测系统

基于 Spring Boot + Vue 3 的 AI 面部识别诊断系统，利用豆包大模型进行腺样体面容分析与检测。

## 功能特性

- **面容分析** - 上传面部照片，AI自动识别是否为腺样体面容并评估程度
- **批量检测** - 支持多张图片批量分析对比
- **检测记录** - 查看和管理历史检测记录
- **统计分析** - 检测数据可视化图表展示
- **AI智能助手** - 基于豆包大模型的对话助手
- **报告导出** - 支持 PDF 报告导出和 Excel 数据导出

## 技术栈

| 后端 | 前端 |
|------|------|
| Spring Boot 3.5 | Vue 3 |
| MyBatis | Vite 7 |
| MySQL 8.0+ | Element Plus |
| PageHelper | ECharts |
| JWT认证 | Axios |
| 阿里云OSS | |
| 豆包大模型API | |

## 快速开始

### 前置条件

- JDK 17+
- Node.js 20.19+
- MySQL 8.0+
- Maven 3.6+

### 1. 克隆项目

```bash
git clone https://github.com/BBtacker/Chendian_myboy.git
cd faceTest
```

### 2. 数据库配置

创建数据库：

```sql
CREATE DATABASE face DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### ⚠️ 3. 配置API密钥（重要）

项目中的敏感配置通过 `application.yml` 管理，该文件已被 `.gitignore` 忽略，**不会提交到GitHub**。

```bash
# 复制配置模板
cp src/main/resources/application-example.yml src/main/resources/application.yml
```

然后编辑 `application.yml`，填入您的真实配置：

| 配置项 | 说明 | 获取方式 |
|--------|------|----------|
| `spring.datasource.password` | 数据库密码 | 您的MySQL密码 |
| `aliyun.oss.accessKeyId` | 阿里云OSS AccessKey | [阿里云RAM控制台](https://ram.console.aliyun.com/) |
| `aliyun.oss.accessKeySecret` | 阿里云OSS AccessKey Secret | 同上 |
| `doubao.api.key` | 豆包大模型API密钥 | [火山引擎控制台](https://console.volcengine.com/) |

### 4. 启动后端

```bash
mvn spring-boot:run
```

后端服务运行在 `http://localhost:8080`

### 5. 启动前端

```bash
cd front
npm install
npm run dev
```

前端服务运行在 `http://localhost:5174`

## 项目结构

```
faceTest/
├── front/                    # 前端项目 (Vue 3)
│   ├── src/
│   │   ├── components/       # Vue组件
│   │   │   ├── FaceAnalysis.vue   # 面容分析
│   │   │   ├── TestResult.vue     # 检测记录
│   │   │   ├── Statistics.vue     # 统计分析
│   │   │   ├── ChatAssistant.vue  # AI助手
│   │   │   ├── Login.vue          # 登录
│   │   │   └── Register.vue       # 注册
│   │   ├── utils/request.js       # Axios请求封装
│   │   └── App.vue                # 根组件
│   └── package.json
├── src/                      # 后端项目 (Spring Boot)
│   ├── main/
│   │   ├── java/.../controller/  # 控制器
│   │   ├── java/.../service/     # 服务层
│   │   ├── java/.../mapper/      # 数据访问层
│   │   ├── java/.../pojo/        # 实体类
│   │   ├── java/.../utils/       # 工具类
│   │   └── resources/
│   │       ├── mapper/           # MyBatis XML映射
│   │       ├── application.yml   # 应用配置（已gitignore）
│   │       └── application-example.yml  # 配置模板
│   └── pom.xml
├── .gitignore
├── LICENSE
└── README.md
```

## API接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/login` | POST | 用户登录 |
| `/register` | POST | 用户注册 |
| `/doubao/analyzeFace` | POST | 面容分析（上传图片） |
| `/testResult/list` | GET | 获取检测记录列表 |
| `/testResult` | PUT | 更新检测记录 |
| `/testResult/{id}` | DELETE | 删除检测记录 |
| `/testResult/export` | GET | 导出Excel |
| `/testResult/exportPdf/{id}` | GET | 导出PDF |
| `/statistics/data` | GET | 获取统计数据 |
| `/conversation/*` | * | AI对话相关接口 |

## 安全说明

- **API密钥保护**：所有敏感配置（数据库密码、阿里云OSS密钥、豆包API密钥）通过 `application.yml` 管理，该文件已加入 `.gitignore`，不会泄露
- **用户认证**：使用 JWT Token 进行身份验证
- **文件上传限制**：单文件最大 10MB

## 开源协议

本项目基于 [MIT License](LICENSE) 开源。

## 作者

- long67 (BBtacker)