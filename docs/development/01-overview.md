# NovelAI Platform 开发环境配置指南

## 目录

本指南帮助开发者在 Windows 环境下配置各种 IDE 进行项目开发。

### 文档索引

| 编号 | 文档 | 适用场景 |
|------|------|----------|
| 01 | [概述](./01-overview.md) | 项目结构与 IDE 选择建议 |
| 02 | [GoLand 配置](./02-goland-setup.md) | Go Gateway 服务开发 |
| 03 | [PyCharm 配置](./03-pycharm-setup.md) | Python AI Service 开发 |
| 04 | [VSCode 配置](./04-vscode-setup.md) | 全栈开发（推荐） |
| 05 | [VS2022 配置](./05-vs2022-setup.md) | C++ 扩展开发（可选） |

---

## 项目结构

```
novelai-platform/
├── go-gateway/          # Go 网关服务 (Gin + gRPC)
├── python-ai-service/   # Python AI 服务 (FastAPI + gRPC)
├── frontend/            # React 前端 (Vite + TypeScript)
├── proto/               # Protocol Buffers 定义
├── keys/                # JWT 密钥对
└── docker-compose.yml   # 容器编排配置
```

---

## IDE 选择建议

| 服务 | 推荐 IDE | 备选 |
|------|----------|------|
| Go Gateway | **GoLand** | VSCode + Go 扩展 |
| Python AI Service | **PyCharm** | VSCode + Python 扩展 |
| React Frontend | **VSCode** | WebStorm |
| 全栈开发 | **VSCode** | 多 IDE 组合 |

---

## 环境要求

### 必需软件

| 软件 | 版本要求 | 下载地址 |
|------|----------|----------|
| Go | 1.21+ | https://go.dev/dl/ |
| Python | 3.11+ | https://www.python.org/downloads/ |
| Node.js | 20+ | https://nodejs.org/ |
| PostgreSQL | 15+ | https://www.postgresql.org/download/ |
| Git | 最新版 | https://git-scm.com/download/win |

### 可选软件

| 软件 | 用途 |
|------|------|
| Docker Desktop | 容器化部署 |
| protoc | Protocol Buffers 编译 |
| Make | 构建自动化 |

---

## 快速开始

### 1. 克隆项目

```powershell
git clone <repository-url>
cd novelai-platform
```

### 2. 配置环境变量

```powershell
# 复制环境变量模板
copy .env.example .env

# 编辑 .env 文件，填入必要配置
notepad .env
```

### 3. 生成 JWT 密钥（如未生成）

```powershell
# 创建 keys 目录
mkdir keys

# 生成 RSA 密钥对
openssl genrsa -out keys/private.pem 2048
openssl rsa -in keys/private.pem -pubout -out keys/public.pem
```

### 4. 选择开发模式

- **单服务开发**: 使用对应 IDE（GoLand/PyCharm）
- **全栈开发**: 使用 VSCode 打开整个项目
- **前端独立开发**: 使用 Mock Server

---

## 下一步

根据你的开发需求，选择对应的 IDE 配置文档：

- 开发 Go Gateway → [02-goland-setup.md](./02-goland-setup.md)
- 开发 Python AI Service → [03-pycharm-setup.md](./03-pycharm-setup.md)
- 全栈开发 → [04-vscode-setup.md](./04-vscode-setup.md)
