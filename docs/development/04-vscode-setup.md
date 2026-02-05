# VSCode 配置指南 - 全栈开发

## 适用场景

- 全栈开发（Go + Python + React）
- 前端 React 开发
- 轻量级后端开发
- 多语言项目协作

---

## 1. 打开项目

### 打开整个项目（推荐）

1. 启动 VSCode
2. 选择 `File` → `Open Folder`
3. 导航到 `novelai-platform` 目录
4. 点击 `选择文件夹`

### 使用工作区

创建 `novelai.code-workspace` 文件：

```json
{
  "folders": [
    { "path": "go-gateway", "name": "Go Gateway" },
    { "path": "python-ai-service", "name": "Python AI Service" },
    { "path": "frontend", "name": "Frontend" }
  ],
  "settings": {
    "files.exclude": {
      "**/__pycache__": true,
      "**/.pytest_cache": true,
      "**/node_modules": true
    }
  }
}
```

然后 `File` → `Open Workspace from File` 打开。

---

## 2. 必装扩展

### Go 开发

| 扩展 | ID | 用途 |
|------|-----|------|
| Go | `golang.go` | Go 语言支持 |

### Python 开发

| 扩展 | ID | 用途 |
|------|-----|------|
| Python | `ms-python.python` | Python 语言支持 |
| Pylance | `ms-python.vscode-pylance` | 智能提示 |
| Black Formatter | `ms-python.black-formatter` | 代码格式化 |

### 前端开发

| 扩展 | ID | 用途 |
|------|-----|------|
| ESLint | `dbaeumer.vscode-eslint` | JS/TS 代码检查 |
| Prettier | `esbenp.prettier-vscode` | 代码格式化 |
| Tailwind CSS IntelliSense | `bradlc.vscode-tailwindcss` | Tailwind 提示 |

### 通用工具

| 扩展 | ID | 用途 |
|------|-----|------|
| GitLens | `eamodio.gitlens` | Git 增强 |
| Docker | `ms-azuretools.vscode-docker` | Docker 支持 |
| Thunder Client | `rangav.vscode-thunder-client` | API 测试 |
| vscode-proto3 | `zxh404.vscode-proto3` | Protobuf 支持 |
| DotENV | `mikestead.dotenv` | .env 文件支持 |

### 一键安装

在终端执行：

```powershell
code --install-extension golang.go
code --install-extension ms-python.python
code --install-extension ms-python.vscode-pylance
code --install-extension ms-python.black-formatter
code --install-extension dbaeumer.vscode-eslint
code --install-extension esbenp.prettier-vscode
code --install-extension bradlc.vscode-tailwindcss
code --install-extension eamodio.gitlens
code --install-extension ms-azuretools.vscode-docker
code --install-extension rangav.vscode-thunder-client
code --install-extension zxh404.vscode-proto3
code --install-extension mikestead.dotenv
```

---

## 3. 配置文件

### 项目设置 (.vscode/settings.json)

在项目根目录创建 `.vscode/settings.json`：

```json
{
  // Go 配置
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": "explicit"
    }
  },

  // Python 配置
  "python.defaultInterpreterPath": "${workspaceFolder}/python-ai-service/.venv/Scripts/python.exe",
  "python.analysis.typeCheckingMode": "basic",
  "[python]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "ms-python.black-formatter"
  },

  // TypeScript/React 配置
  "[typescript]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[typescriptreact]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[javascript]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },

  // 通用配置
  "editor.tabSize": 2,
  "editor.insertSpaces": true,
  "files.eol": "\n",
  "files.trimTrailingWhitespace": true,
  "files.insertFinalNewline": true,

  // 文件关联
  "files.associations": {
    "*.proto": "proto3"
  },

  // 搜索排除
  "search.exclude": {
    "**/node_modules": true,
    "**/.venv": true,
    "**/dist": true,
    "**/__pycache__": true
  }
}
```

---

## 4. 运行/调试配置

### 创建 launch.json

在 `.vscode/launch.json` 中配置：

```json
{
  "version": "0.2.0",
  "configurations": [
    // Go Gateway
    {
      "name": "Go Gateway",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/go-gateway/cmd/server",
      "cwd": "${workspaceFolder}/go-gateway",
      "envFile": "${workspaceFolder}/.env"
    },

    // Python AI Service (FastAPI)
    {
      "name": "Python FastAPI",
      "type": "debugpy",
      "request": "launch",
      "module": "uvicorn",
      "args": ["app.main:app", "--reload", "--host", "0.0.0.0", "--port", "8000"],
      "cwd": "${workspaceFolder}/python-ai-service",
      "envFile": "${workspaceFolder}/.env",
      "jinja": true
    },

    // Python gRPC Server
    {
      "name": "Python gRPC",
      "type": "debugpy",
      "request": "launch",
      "program": "${workspaceFolder}/python-ai-service/app/grpc_server.py",
      "cwd": "${workspaceFolder}/python-ai-service",
      "envFile": "${workspaceFolder}/.env"
    },

    // Frontend (Vite)
    {
      "name": "Frontend Debug",
      "type": "chrome",
      "request": "launch",
      "url": "http://localhost:3000",
      "webRoot": "${workspaceFolder}/frontend/src",
      "sourceMapPathOverrides": {
        "webpack:///./src/*": "${webRoot}/*"
      }
    }
  ],
  "compounds": [
    {
      "name": "Full Stack",
      "configurations": ["Go Gateway", "Python gRPC", "Frontend Debug"]
    }
  ]
}
```

---

## 5. 任务配置

### 创建 tasks.json

在 `.vscode/tasks.json` 中配置：

```json
{
  "version": "2.0.0",
  "tasks": [
    // Go Gateway
    {
      "label": "Go: Build Gateway",
      "type": "shell",
      "command": "go build -o bin/server.exe ./cmd/server",
      "options": { "cwd": "${workspaceFolder}/go-gateway" },
      "group": "build"
    },
    {
      "label": "Go: Run Tests",
      "type": "shell",
      "command": "go test ./...",
      "options": { "cwd": "${workspaceFolder}/go-gateway" },
      "group": "test"
    },
    {
      "label": "Go: Migrate Up",
      "type": "shell",
      "command": "migrate -path migrations -database \"postgres://postgres:password@localhost:5432/novelai?sslmode=disable\" up",
      "options": { "cwd": "${workspaceFolder}/go-gateway" }
    },

    // Python AI Service
    {
      "label": "Python: Install Dependencies",
      "type": "shell",
      "command": ".venv\\Scripts\\pip install -r requirements.txt",
      "options": { "cwd": "${workspaceFolder}/python-ai-service" }
    },
    {
      "label": "Python: Run Tests",
      "type": "shell",
      "command": ".venv\\Scripts\\pytest -v",
      "options": { "cwd": "${workspaceFolder}/python-ai-service" },
      "group": "test"
    },

    // Frontend
    {
      "label": "Frontend: Install",
      "type": "shell",
      "command": "npm install",
      "options": { "cwd": "${workspaceFolder}/frontend" }
    },
    {
      "label": "Frontend: Dev",
      "type": "shell",
      "command": "npm run dev",
      "options": { "cwd": "${workspaceFolder}/frontend" },
      "isBackground": true,
      "problemMatcher": []
    },
    {
      "label": "Frontend: Build",
      "type": "shell",
      "command": "npm run build",
      "options": { "cwd": "${workspaceFolder}/frontend" },
      "group": "build"
    },
    {
      "label": "Frontend: Mock Server",
      "type": "shell",
      "command": "node mock-server.js",
      "options": { "cwd": "${workspaceFolder}/frontend" },
      "isBackground": true,
      "problemMatcher": []
    },

    // Docker
    {
      "label": "Docker: Up",
      "type": "shell",
      "command": "docker-compose up -d",
      "options": { "cwd": "${workspaceFolder}" }
    },
    {
      "label": "Docker: Down",
      "type": "shell",
      "command": "docker-compose down",
      "options": { "cwd": "${workspaceFolder}" }
    },
    {
      "label": "Docker: Logs",
      "type": "shell",
      "command": "docker-compose logs -f",
      "options": { "cwd": "${workspaceFolder}" }
    }
  ]
}
```

### 运行任务

- `Ctrl+Shift+P` → `Tasks: Run Task`
- 或 `Ctrl+Shift+B` 运行构建任务

---

## 6. 终端配置

### 多终端布局

1. 打开终端 (`Ctrl+``)
2. 点击 `+` 创建新终端
3. 右键终端标签 → `Rename` 命名

建议布局：
- Terminal 1: `go-gateway` (Go 服务)
- Terminal 2: `python-ai-service` (Python 服务)
- Terminal 3: `frontend` (前端开发)

### PowerShell 配置

在 `settings.json` 中添加：

```json
{
  "terminal.integrated.defaultProfile.windows": "PowerShell",
  "terminal.integrated.profiles.windows": {
    "PowerShell": {
      "source": "PowerShell",
      "icon": "terminal-powershell"
    }
  }
}
```

---

## 7. 调试技巧

### 多服务调试

1. 使用 `Full Stack` 复合配置同时启动多个服务
2. 或分别启动各服务调试

### 断点调试

- 点击行号左侧设置断点
- 条件断点：右键断点 → `Edit Breakpoint`
- 日志点：右键 → `Add Logpoint`

### 调试控制

- `F5`: 启动/继续
- `F10`: 单步跳过
- `F11`: 单步进入
- `Shift+F11`: 单步跳出
- `Shift+F5`: 停止

---

## 8. Git 集成

### 内置 Git

- 左侧 `Source Control` 面板 (`Ctrl+Shift+G`)
- 查看更改、暂存、提交

### GitLens 功能

- 行内 blame 信息
- 文件历史
- 分支比较

---

## 9. 前端开发专项

### 启动开发服务器

```powershell
cd frontend
npm run dev
```

### 使用 Mock Server

```powershell
cd frontend
node mock-server.js
```

### React DevTools

1. 安装 Chrome 扩展 `React Developer Tools`
2. 在 Chrome DevTools 中使用 `Components` 和 `Profiler` 面板

### Tailwind CSS

- 安装 `Tailwind CSS IntelliSense` 扩展
- 悬停类名查看样式
- 自动补全类名

---

## 10. 常见问题

### Q: Go 扩展提示 gopls 错误？

```powershell
go install golang.org/x/tools/gopls@latest
```

### Q: Python 解释器找不到？

1. `Ctrl+Shift+P` → `Python: Select Interpreter`
2. 选择 `python-ai-service\.venv\Scripts\python.exe`

### Q: ESLint 不工作？

确保 `frontend` 目录下有 `.eslintrc` 配置文件。

### Q: 终端中文乱码？

在 `settings.json` 中添加：

```json
{
  "terminal.integrated.defaultProfile.windows": "PowerShell",
  "terminal.integrated.env.windows": {
    "PYTHONIOENCODING": "utf-8"
  }
}
```

### Q: 保存时格式化不生效？

检查：
1. 对应语言的格式化器已安装
2. `editor.formatOnSave` 已启用
3. 文件类型关联正确
