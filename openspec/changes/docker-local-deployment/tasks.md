## 1. Go Gateway 环境变量配置支持

- [x] 1.1 修改 `internal/config/config.go` 添加环境变量读取逻辑
- [x] 1.2 实现 `parsePostgresURL()` 函数解析 DATABASE_URL
- [x] 1.3 添加环境变量优先级：DATABASE_URL > config.yaml database.*
- [x] 1.4 添加环境变量优先级：GRPC_AI_SERVICE_ADDR > config.yaml grpc.ai_service_addr
- [x] 1.5 添加环境变量优先级：JWT_SECRET > config.yaml auth.jwt_secret
- [x] 1.6 添加配置验证：启动时检查必需配置项（数据库、JWT密钥）
- [x] 1.7 添加配置日志：启动时输出配置来源（env var 或 config.yaml）

## 2. JWT 认证方案改为 HMAC

- [x] 2.1 修改 `internal/auth/jwt.go` 移除 RS256 相关代码
- [x] 2.2 实现 HMAC-SHA256 token 生成函数
- [x] 2.3 实现 HMAC-SHA256 token 验证函数
- [x] 2.4 添加 JWT_SECRET 长度验证（最小 32 字符，不足时警告）
- [x] 2.5 更新 JWT 中间件使用新的 HMAC 验证逻辑
- [x] 2.6 移除 `internal/config/config.go` 中的 jwt_private_key_path 和 jwt_public_key_path 配置项
- [x] 2.7 更新单元测试使用 HMAC 方案

## 3. AI Service 健康检查修复

- [x] 3.1 修改 `python-ai-service/Dockerfile` 添加 curl 安装
- [x] 3.2 优化 Dockerfile：在同一 RUN 层清理 apt 缓存
- [x] 3.3 验证 /health 端点返回正确的 JSON 格式
- [x] 3.4 修改健康检查日志级别为 DEBUG（避免日志刷屏）
- [x] 3.5 测试健康检查在服务启动期间返回 503

## 4. 前端 API 地址运行时注入

- [x] 4.1 创建 `frontend/nginx.conf.template` 包含 ${VITE_API_BASE_URL} 占位符
- [x] 4.2 创建 `frontend/docker-entrypoint.sh` 脚本
- [x] 4.3 在 entrypoint 中实现 envsubst 替换逻辑
- [x] 4.4 添加 VITE_API_BASE_URL 格式验证（必须是 http/https URL）
- [x] 4.5 添加 Nginx 配置语法验证：`nginx -t`
- [x] 4.6 设置默认值：VITE_API_BASE_URL=${VITE_API_BASE_URL:-http://localhost:8080}
- [x] 4.7 修改 `frontend/Dockerfile` 使用多阶段构建
- [x] 4.8 在 Dockerfile 中 COPY entrypoint 脚本并设置可执行权限
- [x] 4.9 设置 ENTRYPOINT 为 docker-entrypoint.sh

## 5. 独立数据库初始化容器

- [x] 5.1 在 `docker-compose.yml` 添加 init-db 服务定义
- [x] 5.2 配置 init-db 使用 migrate/migrate 镜像
- [x] 5.3 挂载 `./go-gateway/migrations` 到 init-db 容器
- [x] 5.4 配置 init-db 的 command 参数（-path, -database, up）
- [x] 5.5 设置 init-db depends_on postgres (condition: service_healthy)
- [x] 5.6 设置 init-db restart: on-failure
- [x] 5.7 修改 gateway 服务 depends_on init-db (condition: service_completed_successfully)
- [x] 5.8 修改 ai-service 服务 depends_on init-db (condition: service_completed_successfully)

## 6. Docker Compose 配置更新

- [x] 6.1 从 gateway 服务移除 JWT_PRIVATE_KEY_PATH 和 JWT_PUBLIC_KEY_PATH 环境变量
- [x] 6.2 添加 gateway 服务环境变量：JWT_SECRET=${JWT_SECRET}
- [x] 6.3 从 gateway 服务移除 keys 目录挂载
- [x] 6.4 更新 ai-service 健康检查配置（interval: 30s, timeout: 10s, retries: 3）
- [x] 6.5 更新 frontend 服务添加 VITE_API_BASE_URL 环境变量
- [x] 6.6 移除 migrate profile（不再需要手动迁移）

## 7. 环境变量配置文件

- [x] 7.1 更新 `.env.example` 添加 JWT_SECRET 配置项
- [x] 7.2 在 .env.example 中添加 JWT_SECRET 生成命令注释
- [x] 7.3 从 .env.example 移除 JWT_PRIVATE_KEY_PATH 和 JWT_PUBLIC_KEY_PATH
- [x] 7.4 添加所有环境变量的详细注释说明
- [x] 7.5 验证 .env.example 包含所有必需配置项
- [x] 7.6 确认 .env 文件在 .gitignore 中

## 8. 清理和文档

- [x] 8.1 删除 `keys/` 目录（如果存在）
- [x] 8.2 更新 README.md 部署步骤（移除手动迁移说明）
- [x] 8.3 添加 JWT_SECRET 生成步骤到部署文档
- [x] 8.4 添加配置验证步骤到部署文档
- [x] 8.5 添加健康检查验证命令到部署文档
- [x] 8.6 添加常见问题排查指南（端口冲突、权限问题等）
- [x] 8.7 更新 docker-compose.yml 中的注释

## 9. 测试和验证

- [x] 9.1 全新环境测试：删除所有容器和卷，执行 docker compose up -d
- [x] 9.2 验证 postgres 健康检查通过
- [x] 9.3 验证 redis 健康检查通过
- [x] 9.4 验证 init-db 自动执行迁移并成功退出
- [x] 9.5 验证 ai-service 健康检查通过（curl 可用）
- [x] 9.6 验证 gateway 健康检查通过：curl http://localhost:8080/healthz
- [x] 9.7 验证 frontend 可访问：curl http://localhost/
- [x] 9.8 测试用户注册功能（验证 JWT HMAC 生成）
- [x] 9.9 测试用户登录功能（验证 JWT HMAC 验证）
- [x] 9.10 测试 AI 生成功能（验证 gRPC 连接）
- [x] 9.11 测试配置 fallback：移除环境变量，使用 config.yaml 启动
- [x] 9.12 测试幂等性：重复执行 docker compose up -d 验证无错误

## 10. 回归测试

- [x] 10.1 验证所有现有 API 端点正常工作
- [x] 10.2 验证数据库连接池正常
- [x] 10.3 验证 Redis 缓存功能正常
- [x] 10.4 验证文件上传功能正常
- [x] 10.5 验证 WebSocket 连接正常（如果有）
- [x] 10.6 检查日志无异常错误
- [x] 10.7 验证容器资源使用在合理范围（内存、CPU）
