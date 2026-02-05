## Context

novelai-platform 项目采用微服务架构（Go Gateway + Python AI Service + React Frontend），使用 docker-compose 进行容器编排。当前存在配置方案不统一的问题：

**当前状态**：
- Go Gateway 代码通过 viper 读取 `config.yaml`，配置项包括 `database.host`、`grpc.ai_service_addr`、`auth.jwt_secret`
- docker-compose.yml 提供环境变量 `DATABASE_URL`、`GRPC_AI_SERVICE_ADDR`、`JWT_PRIVATE_KEY_PATH`
- 前端构建时使用 Vite 的 `import.meta.env.VITE_API_BASE_URL`，值在构建时固化到 JS bundle
- 数据库迁移需要手动执行 `docker compose --profile migrate run --rm migrate`
- AI Service 健康检查使用 curl 但基础镜像不包含该工具

**约束条件**：
- 必须保持向后兼容，现有使用 config.yaml 的部署方式不能破坏
- 不能修改数据库 schema 或 API 接口
- 镜像大小增加应控制在合理范围（< 10MB）
- 首次部署时间应在 5 分钟内完成

**利益相关者**：
- 开发者：需要简化本地开发环境搭建
- 运维人员：需要统一的配置管理方式
- 最终用户：不受影响（内部架构变更）

## Goals / Non-Goals

**Goals:**
- 实现 `docker compose up -d` 一键部署，无需手动迁移或配置文件编辑
- 统一配置方案为环境变量优先，符合 12-factor app 原则
- 保持向后兼容，支持 config.yaml 作为 fallback
- 修复所有阻碍容器化部署的技术问题（健康检查、JWT 方案、前端配置）
- 提供完整的 .env.example 和部署文档

**Non-Goals:**
- 不重构现有业务逻辑或 API 接口
- 不改变数据库 schema 或迁移文件内容
- 不实现生产环境部署方案（HTTPS、域名、负载均衡等）
- 不添加配置热重载功能
- 不实现配置加密或密钥管理服务集成

## Decisions

### Decision 1: 环境变量优先级高于 config.yaml

**选择**：修改 Go Gateway 配置加载逻辑，优先读取环境变量，config.yaml 作为 fallback。

**理由**：
- 符合 12-factor app 原则，配置与代码分离
- 容器化部署的标准实践，便于不同环境（dev/staging/prod）使用相同镜像
- 向后兼容，不破坏现有部署方式

**替代方案**：
- ❌ 仅使用 config.yaml：需要为每个环境维护不同配置文件，容器化部署复杂
- ❌ 仅使用环境变量：破坏向后兼容性，现有部署需要迁移
- ❌ 使用配置中心（Consul/etcd）：过度设计，增加部署复杂度

**实现方式**：
```go
// 伪代码示例
func LoadConfig() *Config {
    cfg := &Config{}

    // 1. 先加载 config.yaml（如果存在）
    if fileExists("config.yaml") {
        viper.ReadInConfig()
        viper.Unmarshal(cfg)
    }

    // 2. 环境变量覆盖（优先级更高）
    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.Database = parsePostgresURL(dbURL)
    }
    if grpcAddr := os.Getenv("GRPC_AI_SERVICE_ADDR"); grpcAddr != "" {
        cfg.GRPC.AIServiceAddr = grpcAddr
    }
    if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
        cfg.Auth.JWTSecret = jwtSecret
    }

    return cfg
}
```

### Decision 2: JWT 认证改用 HMAC-SHA256

**选择**：从 RS256（非对称加密）改为 HMAC-SHA256（对称加密），使用环境变量 JWT_SECRET。

**理由**：
- 简化部署：无需生成和管理 RSA 密钥对文件
- 单体应用场景：当前架构中只有 Gateway 生成和验证 token，不需要非对称加密的优势
- 安全性足够：HMAC-SHA256 配合强密钥（32+ 字符）满足安全要求
- 配置统一：与其他环境变量配置方式一致

**替代方案**：
- ❌ 保持 RS256：需要在容器中挂载密钥文件，增加部署复杂度和安全风险
- ❌ 使用 EdDSA：更现代但 Go JWT 库支持不够成熟
- ❌ 使用外部认证服务（OAuth2）：过度设计，增加依赖

**权衡**：
- ✅ 优势：部署简单，配置统一，性能更好（HMAC 比 RSA 快）
- ⚠️ 劣势：密钥泄露风险更高（对称密钥），不支持多服务验证场景
- 🔄 迁移影响：所有现有 token 失效，用户需要重新登录

**安全措施**：
- JWT_SECRET 最小长度 32 字符（启动时验证并警告）
- 使用 `openssl rand -base64 32` 生成强随机密钥
- .env 文件添加到 .gitignore，避免提交到版本控制

### Decision 3: 前端 API 地址运行时注入

**选择**：使用 entrypoint 脚本 + envsubst 在容器启动时替换 Nginx 配置中的 API 地址。

**理由**：
- Vite 构建时环境变量会固化到 JS bundle，无法在运行时修改
- Nginx 配置支持运行时替换，是容器化前端的标准实践
- 支持同一镜像部署到不同环境（localhost、staging、production）

**替代方案**：
- ❌ 构建时指定 API 地址：每个环境需要单独构建镜像，违反容器化原则
- ❌ 使用 window.env.js 动态加载：需要额外 HTTP 请求，增加首屏加载时间
- ❌ 使用相对路径 /api + Nginx 反向代理：需要修改前端代码中所有 API 调用

**实现方式**：
```dockerfile
# frontend/Dockerfile
FROM node:18 AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf.template /etc/nginx/templates/default.conf.template
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
```

```bash
# docker-entrypoint.sh
#!/bin/sh
export VITE_API_BASE_URL=${VITE_API_BASE_URL:-http://localhost:8080}
envsubst '${VITE_API_BASE_URL}' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf
exec "$@"
```

### Decision 4: 独立初始化容器自动执行迁移

**选择**：添加 init-db 服务，使用 depends_on 确保在应用服务前执行迁移。

**理由**：
- 自动化：无需手动执行 `--profile migrate`
- 幂等性：重复执行不会出错，golang-migrate 自动跟踪已应用的迁移
- 依赖管理：通过 depends_on 确保执行顺序
- 故障隔离：迁移失败时应用服务不会启动

**替代方案**：
- ❌ Gateway entrypoint 中执行迁移：并发启动时可能冲突，增加 Gateway 启动时间
- ❌ 保持手动执行：不符合一键部署目标
- ❌ 使用 Kubernetes Init Container：过度设计，docker-compose 场景不适用

**实现方式**：
```yaml
# docker-compose.yml
services:
  init-db:
    image: migrate/migrate
    volumes:
      - ./go-gateway/migrations:/migrations:ro
    command: [
      "-path", "/migrations",
      "-database", "postgres://postgres:${POSTGRES_PASSWORD}@postgres:5432/novelai?sslmode=disable",
      "up"
    ]
    depends_on:
      postgres:
        condition: service_healthy
    restart: on-failure

  gateway:
    depends_on:
      init-db:
        condition: service_completed_successfully
```

### Decision 5: AI Service 健康检查安装 curl

**选择**：在 python-ai-service/Dockerfile 中添加 `apt-get install curl`。

**理由**：
- 简单直接：一行 RUN 命令解决问题
- 镜像增加可控：约 5MB，可接受
- 标准实践：curl 是容器健康检查的常用工具

**替代方案**：
- ❌ 使用 wget：python:3.11-slim 也不包含 wget
- ❌ 使用 Python urllib：健康检查命令过长，不够优雅
- ❌ 改用 TCP 端口检查：无法验证应用层是否正常

**优化措施**：
```dockerfile
RUN apt-get update && \
    apt-get install -y --no-install-recommends curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
```

## Risks / Trade-offs

### Risk 1: JWT 密钥泄露风险增加
**风险**：HMAC 使用对称密钥，泄露后攻击者可伪造任意 token。

**缓解措施**：
- 强制最小密钥长度 32 字符
- .env 文件添加到 .gitignore
- 部署文档中强调密钥管理重要性
- 考虑未来添加密钥轮换机制

### Risk 2: 所有用户需要重新登录
**风险**：从 RS256 切换到 HMAC 后，所有现有 token 失效。

**缓解措施**：
- 在部署文档中明确说明
- 提供用户通知模板
- 考虑在 API 响应中添加友好错误提示

### Risk 3: 环境变量配置错误难以调试
**风险**：环境变量拼写错误或格式错误时，错误信息可能不够明确。

**缓解措施**：
- 启动时验证所有必需环境变量
- 提供详细的错误日志，包括期望格式示例
- 在 .env.example 中添加详细注释

### Risk 4: 数据库迁移失败导致服务无法启动
**风险**：init-db 失败时，应用服务被阻塞。

**缓解措施**：
- init-db 使用 restart: on-failure，自动重试
- 提供清晰的失败日志和排查指南
- 支持跳过 init-db 的手动启动模式（用于调试）

### Risk 5: 前端 entrypoint 脚本执行失败
**风险**：envsubst 失败或 Nginx 配置语法错误导致容器无法启动。

**缓解措施**：
- entrypoint 脚本中添加配置验证：`nginx -t`
- 提供默认值：`${VITE_API_BASE_URL:-http://localhost:8080}`
- 在 CI 中测试镜像构建和启动

## Migration Plan

### 部署步骤

**准备阶段**：
1. 备份现有数据库（如果有生产数据）
2. 生成 JWT_SECRET：`openssl rand -base64 32`
3. 复制 .env.example 为 .env，填写所有必需配置

**部署阶段**：
1. 拉取最新代码：`git pull origin main`
2. 构建镜像：`docker compose build`
3. 启动服务：`docker compose up -d`
4. 验证健康检查：
   ```bash
   curl http://localhost:8080/healthz  # Gateway
   curl http://localhost:8001/health   # AI Service
   curl http://localhost/              # Frontend
   ```

**验证阶段**：
1. 检查所有容器状态：`docker compose ps`
2. 查看日志确认无错误：`docker compose logs`
3. 测试用户注册/登录流程
4. 测试 AI 生成功能

### 回滚策略

**如果部署失败**：
1. 停止所有服务：`docker compose down`
2. 恢复到上一个工作版本：`git checkout <previous-commit>`
3. 使用旧的部署方式（手动迁移 + config.yaml）
4. 恢复数据库备份（如果数据损坏）

**如果部分功能异常**：
1. 保持服务运行，查看日志定位问题
2. 修复配置错误（.env 文件）
3. 重启受影响的服务：`docker compose restart <service>`

### 兼容性保证

- 数据库 schema 不变，迁移文件不修改
- API 接口不变，客户端无需更新
- 支持 config.yaml fallback，现有部署方式仍可用
- 环境变量命名遵循现有约定

## Open Questions

1. **是否需要支持配置热重载？**
   - 当前方案需要重启容器才能应用配置变更
   - 可考虑使用 SIGHUP 信号触发配置重载
   - 优先级：低（首次部署后配置变更频率低）

2. **是否需要添加配置验证 CLI 工具？**
   - 提供 `docker compose run gateway validate-config` 命令
   - 在启动前验证所有配置项格式和连通性
   - 优先级：中（可提升用户体验）

3. **是否需要支持多环境配置文件？**
   - 例如 .env.dev、.env.staging、.env.prod
   - 通过 `--env-file` 参数切换
   - 优先级：低（当前目标是本地开发环境）

4. **JWT 密钥轮换机制？**
   - 支持同时验证新旧密钥，平滑过渡
   - 需要修改 JWT 验证逻辑
   - 优先级：低（可作为后续增强）
