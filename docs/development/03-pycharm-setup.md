# PyCharm 配置指南 - Python AI Service 开发

## 适用场景

- Python AI Service 开发
- FastAPI 后端开发
- gRPC 客户端/服务端开发
- AI 模型集成

---

## 1. 打开项目

### 方式一：打开子目录（推荐）

1. 启动 PyCharm
2. 选择 `File` → `Open`
3. 导航到 `novelai-platform/python-ai-service` 目录
4. 点击 `OK`

### 方式二：打开整个项目

1. 启动 PyCharm
2. 选择 `File` → `Open`
3. 导航到 `novelai-platform` 目录
4. 点击 `OK`
5. 右键 `python-ai-service` → `Mark Directory as` → `Sources Root`

---

## 2. 配置 Python 解释器

### 创建虚拟环境

1. 打开 `File` → `Settings` (Ctrl+Alt+S)
2. 导航到 `Project: python-ai-service` → `Python Interpreter`
3. 点击齿轮图标 → `Add...`
4. 选择 `Virtualenv Environment` → `New environment`
5. 配置：

| 字段 | 值 |
|------|-----|
| Location | `python-ai-service\.venv` |
| Base interpreter | Python 3.11+ |

6. 点击 `OK`

### 使用已有虚拟环境

如果已创建虚拟环境：

1. 点击齿轮图标 → `Add...`
2. 选择 `Virtualenv Environment` → `Existing environment`
3. 选择 `python-ai-service\.venv\Scripts\python.exe`

---

## 3. 安装依赖

### 方式一：通过 PyCharm

1. 打开 `requirements.txt`
2. PyCharm 会提示 `Install requirements`
3. 点击安装

### 方式二：通过终端

```powershell
cd python-ai-service
.\.venv\Scripts\activate
pip install -r requirements.txt
```

---

## 4. 配置运行/调试

### FastAPI 服务配置

1. 点击右上角 `Add Configuration...`
2. 点击 `+` → `Python`
3. 配置如下：

| 字段 | 值 |
|------|-----|
| Name | `FastAPI Server` |
| Script path | 选择 `Module name` |
| Module name | `uvicorn` |
| Parameters | `app.main:app --reload --host 0.0.0.0 --port 8000` |
| Working directory | `$ProjectFileDir$` |
| Environment variables | 见下方 |

### gRPC 服务配置

1. 点击 `+` → `Python`
2. 配置如下：

| 字段 | 值 |
|------|-----|
| Name | `gRPC Server` |
| Script path | `app/grpc_server.py` |
| Working directory | `$ProjectFileDir$` |
| Environment variables | 见下方 |

### 环境变量

点击 `Environment variables` 右侧的 `...`，添加：

```
DATABASE_URL=postgresql+asyncpg://postgres:password@localhost:5432/novelai
OPENAI_API_KEY=your_openai_key
ANTHROPIC_API_KEY=your_anthropic_key
GOOGLE_API_KEY=your_google_key
```

或者使用 EnvFile 插件加载 `.env` 文件。

---

## 5. 配置 EnvFile 插件

### 安装插件

1. 打开 `File` → `Settings`
2. 导航到 `Plugins`
3. 搜索 `EnvFile`
4. 点击 `Install`
5. 重启 PyCharm

### 使用 EnvFile

1. 编辑运行配置
2. 切换到 `EnvFile` 标签
3. 勾选 `Enable EnvFile`
4. 点击 `+` 添加 `.env` 文件路径

---

## 6. 代码风格配置

### 配置 Black 格式化

1. 打开 `File` → `Settings`
2. 导航到 `Tools` → `Black`
3. 勾选 `On save`
4. 配置 Black 路径（通常自动检测）

### 配置 isort

1. 导航到 `Editor` → `Code Style` → `Python`
2. 切换到 `Imports` 标签
3. 勾选 `Sort imports`

### 配置 Flake8/Pylint

1. 导航到 `Tools` → `External Tools`
2. 点击 `+` 添加：

| 字段 | 值 |
|------|-----|
| Name | `Flake8` |
| Program | `$PyInterpreterDirectory$/flake8` |
| Arguments | `$FilePath$` |
| Working directory | `$ProjectFileDir$` |

---

## 7. 调试技巧

### 设置断点

- 点击代码行号左侧设置断点
- 右键断点可设置条件断点

### 启动调试

1. 选择运行配置
2. 点击 `Debug` 按钮（虫子图标）
3. 或按 `Shift+F9`

### 调试 FastAPI

调试模式下访问 `http://localhost:8000/docs` 可查看 Swagger UI。

### 调试控制台

- `F8`: 单步跳过
- `F7`: 单步进入
- `Shift+F8`: 单步跳出
- `F9`: 继续执行
- `Alt+F8`: 计算表达式

---

## 8. 数据库工具

PyCharm Professional 内置数据库工具：

1. 打开 `View` → `Tool Windows` → `Database`
2. 点击 `+` → `Data Source` → `PostgreSQL`
3. 配置连接信息
4. 点击 `Test Connection` 验证

---

## 9. 常用插件

| 插件 | 用途 |
|------|------|
| EnvFile | .env 文件支持 |
| Protocol Buffers | .proto 文件支持 |
| .ignore | .gitignore 支持 |
| Rainbow Brackets | 括号高亮 |
| GitToolBox | Git 增强 |

---

## 10. 测试配置

### 配置 pytest

1. 打开 `File` → `Settings`
2. 导航到 `Tools` → `Python Integrated Tools`
3. 设置 `Default test runner` 为 `pytest`

### 创建测试配置

1. 点击 `Add Configuration...`
2. 点击 `+` → `pytest`
3. 配置：

| 字段 | 值 |
|------|-----|
| Name | `All Tests` |
| Target | `Custom` |
| Additional Arguments | `-v` |
| Working directory | `$ProjectFileDir$` |

---

## 常见问题

### Q: pip install 很慢？

配置国内镜像：

```powershell
pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
```

### Q: 找不到模块？

1. 确保虚拟环境已激活
2. 检查 `Sources Root` 设置
3. 右键 `python-ai-service` → `Mark Directory as` → `Sources Root`

### Q: gRPC 相关导入报错？

重新生成 gRPC 代码：

```powershell
cd python-ai-service
python -m grpc_tools.protoc -I../proto --python_out=./app/proto --grpc_python_out=./app/proto ../proto/ai_service.proto
```

### Q: 异步调试问题？

确保使用 `asyncio` 调试模式：

1. 打开 `File` → `Settings`
2. 导航到 `Build, Execution, Deployment` → `Python Debugger`
3. 勾选 `Gevent compatible`
