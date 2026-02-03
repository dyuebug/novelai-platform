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
| API Gateway | Go 1.21+, Gin, gRPC, JWT (RS256) |
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

## 快速开始

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
# Go Gateway
cp go-gateway/.env.example go-gateway/.env

# Python AI Service
cp python-ai-service/.env.example python-ai-service/.env

# Frontend
cp frontend/.env.example frontend/.env
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
- API 文档: http://localhost:8080/swagger/index.html
- AI Service: http://localhost:8001

## Docker Compose 部署

```bash
docker-compose up -d
```

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
