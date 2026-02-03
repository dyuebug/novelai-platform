# Go Gateway

Go 1.21+ / Gin / Viper / Zap / RateLimit / Gobreaker / gRPC / SSE

## 运行

```bash
go mod tidy
go run ./cmd/gateway
```

## 配置

复制示例配置并按需修改：

```bash
cp ./configs/config.yaml.example ./configs/config.yaml
```

## 路由

- POST /api/v1/auth/login
- POST /api/v1/auth/register
- GET /api/v1/projects
- POST /api/v1/projects
- GET /api/v1/projects/:id/chapters
- GET /api/v1/stream/generate (SSE)
- GET /healthz

## 架构

```
go-gateway/
├── cmd/gateway/          # 入口
├── internal/
│   ├── app/              # 应用初始化
│   ├── config/           # 配置管理
│   ├── router/           # 路由注册
│   ├── middleware/       # 中间件
│   ├── handler/          # HTTP 处理器
│   ├── service/          # 业务服务
│   └── grpcclient/       # gRPC 客户端
├── pkg/
│   ├── logger/           # 日志
│   └── response/         # 响应封装
├── api/proto/            # Proto 定义
└── configs/              # 配置文件
```

## 说明

- gRPC 客户端使用 JSON 编解码占位（便于无生成器运行），生产建议使用 protobuf 生成代码。
- SSE 接口会转发 AI 流式响应。
- JWT 校验为 HMAC 最小实现，OAuth2/RBAC 建议后续单独模块化实现。
