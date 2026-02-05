## Why

novelai-platform 项目当前存在多个配置不匹配问题,导致无法直接使用 docker-compose 成功部署。Go Gateway 代码读取 config.yaml 但 docker-compose 提供环境变量、JWT 方案不一致(代码 HMAC vs 部署 RS256)、健康检查缺少依赖、前端 API 地址构建时固化等问题,阻碍了全新环境的快速部署。需要统一配置方案并修复这些冲突,实现一键部署能力。

## What Changes

- **修改 Go Gateway 配置加载逻辑**:优先读取环境变量(DATABASE_URL、GRPC_AI_SERVICE_ADDR、JWT_SECRET),向后兼容 config.yaml
- **统一 JWT 认证方案为 HMAC**:删除 RS256 密钥文件要求,使用环境变量 JWT_SECRET
- **修复 AI Service 健康检查**:在 Dockerfile 中安装 curl
- **实现前端 API 地址运行时注入**:使用 envsubst 在容器启动时替换 Nginx 配置
- **添加独立数据库初始化容器**:自动执行迁移,无需手动 --profile migrate
- **完善 .env.example 和部署文档**:包含所有必需配置项和部署步骤

## Capabilities

### New Capabilities
- `env-based-config`:Go Gateway 支持环境变量配置,优先级高于 config.yaml
- `runtime-frontend-config`:前端容器启动时动态注入 API 地址
- `auto-db-migration`:独立初始化容器自动执行数据库迁移
- `health-check-fix`:AI Service 容器包含健康检查所需依赖

### Modified Capabilities
- `jwt-authentication`:从 RS256 密钥文件方案改为 HMAC 环境变量方案

## Impact

**代码变更**:
- `go-gateway/internal/config/`:配置加载逻辑,增加环境变量读取
- `go-gateway/internal/auth/`:JWT 验证逻辑,移除 RS256 实现
- `python-ai-service/Dockerfile`:添加 curl 安装
- `frontend/Dockerfile`:添加 entrypoint 脚本实现运行时配置注入
- `docker-compose.yml`:添加 init-db 服务,移除 keys 挂载

**配置文件**:
- `.env.example`:添加 JWT_SECRET,移除密钥文件路径
- `keys/`:目录可删除

**部署流程**:
- 简化为 `docker compose up -d` 一键启动,无需手动迁移
- 首次部署需配置 .env 文件(数据库密码、JWT 密钥、AI API Keys)

**向后兼容性**:
- Go Gateway 仍支持 config.yaml,但环境变量优先级更高
- 现有使用 config.yaml 的部署方式不受影响
