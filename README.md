# NovelAI Platform

AI 驱动的智能小说创作平台，提供从世界观构建、角色设计、大纲生成到章节创作的全流程智能辅助。

## 架构概览

```
novelai-platform/
├── go-gateway/           # Go API Gateway (Gin + gRPC)
├── python-ai-service/    # Python AI Service (FastAPI + gRPC)
├── frontend/             # React 18 Frontend (Vite + TypeScript)
├── scripts/              # 工具脚本
└── docs/                 # API 规范文档
```

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | React 18, TypeScript, Vite, Ant Design, TailwindCSS, Zustand |
| API Gateway | Go 1.21+, Gin, gRPC, JWT (HMAC-SHA256) |
| AI 服务 | Python 3.11+, FastAPI, gRPC, OpenAI/Anthropic/Gemini |
| 数据库 | PostgreSQL 15+ (pgvector), Redis |
| 部署 | Docker, Docker Compose |

## 核心功能

### Phase 1: 基础架构
- [x] 用户认证 (注册/登录/JWT)
- [x] OAuth2 登录 (GitHub/Google/LinuxDo)
- [x] 密码重置 (邮件验证)
- [x] gRPC 通信 (支持 mTLS)

### Phase 2: 核心功能
- [x] 项目管理 (CRUD + 软删除)
- [x] 章节管理 (CRUD + 版本控制)
- [x] AI 生成 (流式输出)
- [x] 世界观管理 (角色/地点/组织/设定)
- [x] RAG 检索 (pgvector 向量存储)

### Phase 3: 高级功能
- [x] 追读力分析 (钩子/爽点/微兑现)
- [x] 一致性检查 (角色/时间线/世界观)
- [x] 多 Agent 审查 (6 Agent 并行)
- [x] 宪法约束 (硬约束/软约束 + 豁免)
- [x] 伏笔管理 (CRUD + 提醒 + 时间线)

## 快速开始（Docker 部署 - 推荐）

### 环境要求

- Docker 20.10+
- Docker Compose 2.0+

### 1. 克隆项目

```bash
git clone https://github.com/your-org/novelai-platform.git
cd novelai-platform
```

### 2. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，配置以下必需项：
# 1. POSTGRES_PASSWORD - 数据库密码
# 2. JWT_SECRET - JWT 密钥（生成命令：openssl rand -base64 32）
# 3. 至少一个 AI Provider API Key (OPENAI_API_KEY/ANTHROPIC_API_KEY/GEMINI_API_KEY)
```

**生成 JWT 密钥：**

```bash
# 生成强随机密钥（推荐）
openssl rand -base64 32

# 或使用其他方法生成至少 32 字符的随机字符串
```

### 3. 启动服务

```bash
# 一键启动所有服务（包括自动数据库迁移）
docker compose up -d

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f
```

### 4. 验证部署

**健康检查：**

```bash
# 检查 Gateway
curl http://localhost:8080/healthz

# 检查 AI Service
curl http://localhost:8001/health

# 检查 Frontend
curl http://localhost/
```

**访问应用：**

- 前端: http://localhost
- API Gateway: http://localhost:8080
- AI Service: http://localhost:8001

### 5. 停止服务

```bash
# 停止所有服务
docker compose down

# 停止并删除数据卷（清空数据库）
docker compose down -v
```

---

## 本地开发环境

### 环境要求

- Go 1.21+
- Python 3.11+
- Node.js 20+
- PostgreSQL 15+ (with pgvector extension)
- Redis 7+

### 1. 克隆项目

```bash
git clone https://github.com/your-org/novelai-platform.git
cd novelai-platform
```

### 2. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件配置必需项
```

### 3. 启动数据库

```bash
# PostgreSQL with pgvector
docker run -d --name novelai-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=novelai \
  -p 5432:5432 \
  pgvector/pgvector:pg15

# Redis
docker run -d --name novelai-redis \
  -p 6379:6379 \
  redis:7-alpine
```

### 4. 执行数据库迁移

```bash
cd go-gateway
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/novelai?sslmode=disable" up
```

### 5. 启动服务

```bash
# Terminal 1: Go Gateway
cd go-gateway
go run cmd/gateway/main.go

# Terminal 2: Python AI Service
cd python-ai-service
pip install -r requirements.txt
python -m uvicorn app.main:app --port 8001

# Terminal 3: Frontend
cd frontend
npm install
npm run dev
```

### 6. 访问应用

- 前端: http://localhost:5173
- API Gateway: http://localhost:8080
- AI Service: http://localhost:8001

---

## 常见问题排查

### 端口冲突

如果遇到端口占用错误，检查以下端口是否被占用：

```bash
# Windows
netstat -ano | findstr "5432 6379 8080 8001 50051 80"

# Linux/Mac
lsof -i :5432 -i :6379 -i :8080 -i :8001 -i :50051 -i :80
```

解决方案：
1. 停止占用端口的进程
2. 或修改 docker-compose.yml 中的端口映射

### 数据库连接失败

**症状**：Gateway 启动失败，日志显示 "Database connection timeout"

**排查步骤**：

```bash
# 1. 检查 PostgreSQL 容器状态
docker compose ps postgres

# 2. 检查 PostgreSQL 日志
docker compose logs postgres

# 3. 手动测试连接
docker exec -it novelai-postgres psql -U postgres -d novelai
```

**解决方案**：
- 确认 POSTGRES_PASSWORD 配置正确
- 等待 PostgreSQL 健康检查通过（约 10-30 秒）
- 检查防火墙设置

### JWT 认证失败

**症状**：登录后 API 请求返回 401 Unauthorized

**排查步骤**：

```bash
# 检查 JWT_SECRET 是否配置
docker compose exec gateway env | grep JWT_SECRET
```

**解决方案**：
- 确认 .env 文件中 JWT_SECRET 已配置且长度 >= 32 字符
- 重启 Gateway 服务：`docker compose restart gateway`
- 清除浏览器缓存和 localStorage

### AI Service 健康检查失败

**症状**：ai-service 容器一直显示 unhealthy

**排查步骤**：

```bash
# 1. 检查容器日志
docker compose logs ai-service

# 2. 手动测试健康检查
docker compose exec ai-service curl -f http://localhost:8001/health
```

**解决方案**：
- 确认至少配置了一个 AI Provider API Key
- 检查 Python 依赖是否正确安装
- 重新构建镜像：`docker compose build ai-service`

### 前端无法连接后端

**症状**：前端页面加载正常，但 API 请求失败

**排查步骤**：

```bash
# 1. 检查 Gateway 是否正常
curl http://localhost:8080/healthz

# 2. 检查前端 Nginx 配置
docker compose exec frontend cat /etc/nginx/conf.d/default.conf
```

**解决方案**：
- 确认 VITE_API_BASE_URL 配置正确（Docker 环境应为 `http://gateway:8080`）
- 重启前端服务：`docker compose restart frontend`
- 检查浏览器控制台网络请求

### 数据库迁移失败

**症状**：init-db 容器退出，状态码非 0

**排查步骤**：

```bash
# 查看迁移日志
docker compose logs init-db
```

**解决方案**：
- 确认 PostgreSQL 已完全启动（健康检查通过）
- 检查迁移文件语法是否正确
- 手动执行迁移排查问题：
  ```bash
  docker compose run --rm init-db
  ```

---

## Docker Compose 部署

**已集成到快速开始部分，请参阅上方"快速开始（Docker 部署 - 推荐）"章节。**

---

## 开发指南

### 生成 mTLS 证书 (可选)

```bash
./scripts/generate-certs.sh
```

### 编译 Protobuf

```bash
# Go
cd go-gateway
protoc --go_out=. --go-grpc_out=. proto/ai_service.proto

# Python
cd python-ai-service
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. app/proto/ai_service.proto
```

### 运行测试

```bash
# Go
cd go-gateway && go test ./...

# Python
cd python-ai-service && pytest

# Frontend
cd frontend && npm test
```

## API 规范

详细 API 规范请参阅 `docs/` 目录：

- [认证服务](docs/auth-spec.md)
- [项目管理](docs/project-spec.md)
- [章节管理](docs/chapter-spec.md)
- [世界观管理](docs/world-spec.md)
- [AI 服务](docs/ai-service-spec.md)
- [伏笔管理](docs/foreshadow-spec.md)

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request。
