# GoLand 配置指南 - Go Gateway 开发

## 适用场景

- Go Gateway 服务开发
- gRPC 服务端开发
- 数据库迁移管理

---

## 1. 打开项目

### 方式一：打开子目录（推荐）

1. 启动 GoLand
2. 选择 `File` → `Open`
3. 导航到 `novelai-platform/go-gateway` 目录
4. 点击 `OK`

### 方式二：打开整个项目

1. 启动 GoLand
2. 选择 `File` → `Open`
3. 导航到 `novelai-platform` 目录
4. 点击 `OK`
5. 右键 `go-gateway` → `Mark Directory as` → `Go Module`

---

## 2. 配置 Go SDK

1. 打开 `File` → `Settings` (Ctrl+Alt+S)
2. 导航到 `Go` → `GOROOT`
3. 点击 `+` 添加 Go SDK
4. 选择 Go 安装目录（如 `C:\Go` 或 `C:\Program Files\Go`）
5. 确保版本 ≥ 1.21

---

## 3. 配置 Go Modules

1. 打开 `File` → `Settings`
2. 导航到 `Go` → `Go Modules`
3. 勾选 `Enable Go modules integration`
4. 设置 `GOPROXY`:
   ```
   https://goproxy.cn,direct
   ```
5. 点击 `Apply`

---

## 4. 下载依赖

在 GoLand 终端中执行：

```powershell
cd go-gateway
go mod download
go mod tidy
```

或者右键 `go.mod` 文件 → `Go` → `Download Modules`

---

## 5. 配置运行/调试

### 创建运行配置

1. 点击右上角 `Add Configuration...`
2. 点击 `+` → `Go Build`
3. 配置如下：

| 字段 | 值 |
|------|-----|
| Name | `Go Gateway` |
| Run kind | `Package` |
| Package path | `novelai-platform/go-gateway/cmd/server` |
| Working directory | `$ProjectFileDir$` |
| Environment | 见下方 |

### 环境变量配置

点击 `Environment` 右侧的 `...`，添加：

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=novelai
JWT_PRIVATE_KEY_PATH=../keys/private.pem
JWT_PUBLIC_KEY_PATH=../keys/public.pem
GRPC_AI_SERVICE_ADDR=localhost:50051
```

或者勾选 `Use .env file`，选择项目根目录的 `.env` 文件。

---

## 6. 数据库迁移

### 安装 golang-migrate

```powershell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 创建迁移配置

1. 点击 `Add Configuration...`
2. 点击 `+` → `Shell Script`
3. 配置如下：

| 字段 | 值 |
|------|-----|
| Name | `DB Migrate Up` |
| Script text | `migrate -path migrations -database "postgres://postgres:password@localhost:5432/novelai?sslmode=disable" up` |
| Working directory | `$ProjectFileDir$/go-gateway` |

---

## 7. 代码风格配置

### 导入 Go 代码风格

1. 打开 `File` → `Settings`
2. 导航到 `Editor` → `Code Style` → `Go`
3. 配置：
   - `Tabs and Indents`: 使用 Tab，Tab size = 4
   - `Imports`: 勾选 `Group imports`，`Sort imports`

### 启用 gofmt/goimports

1. 打开 `File` → `Settings`
2. 导航到 `Tools` → `File Watchers`
3. 点击 `+` → `go fmt` 或 `goimports`

---

## 8. 调试技巧

### 设置断点

- 点击代码行号左侧设置断点
- 右键断点可设置条件断点

### 启动调试

1. 点击运行配置旁的 `Debug` 按钮（虫子图标）
2. 或按 `Shift+F9`

### 调试控制台

- `F8`: 单步跳过
- `F7`: 单步进入
- `Shift+F8`: 单步跳出
- `F9`: 继续执行

---

## 9. 常用插件

| 插件 | 用途 |
|------|------|
| Protocol Buffers | .proto 文件语法高亮 |
| Database Tools | 数据库管理（内置） |
| EnvFile | .env 文件支持 |
| GitToolBox | Git 增强 |

### 安装插件

1. 打开 `File` → `Settings`
2. 导航到 `Plugins`
3. 搜索并安装

---

## 10. 数据库工具

GoLand 内置数据库工具：

1. 打开 `View` → `Tool Windows` → `Database`
2. 点击 `+` → `Data Source` → `PostgreSQL`
3. 配置连接：

| 字段 | 值 |
|------|-----|
| Host | localhost |
| Port | 5432 |
| User | postgres |
| Password | your_password |
| Database | novelai |

4. 点击 `Test Connection` 验证
5. 点击 `OK`

---

## 常见问题

### Q: go mod download 很慢？

设置 Go 代理：

```powershell
go env -w GOPROXY=https://goproxy.cn,direct
```

### Q: 找不到依赖包？

1. 右键 `go.mod` → `Go` → `Sync Dependencies`
2. 或在终端执行 `go mod tidy`

### Q: 调试时无法连接数据库？

确保：
1. PostgreSQL 服务已启动
2. 环境变量配置正确
3. 数据库已创建

```powershell
# 检查 PostgreSQL 服务
Get-Service postgresql*

# 创建数据库
psql -U postgres -c "CREATE DATABASE novelai;"
```
