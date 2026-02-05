# NovelAI Platform 文档中心

## 文档分类

### 开发环境配置

| 编号 | 文档 | 说明 |
|------|------|------|
| 01 | [开发环境概述](./development/01-overview.md) | 项目结构、IDE 选择、环境要求 |
| 02 | [GoLand 配置](./development/02-goland-setup.md) | Go Gateway 服务开发 |
| 03 | [PyCharm 配置](./development/03-pycharm-setup.md) | Python AI Service 开发 |
| 04 | [VSCode 配置](./development/04-vscode-setup.md) | 全栈开发（推荐） |
| 05 | [VS2022 配置](./development/05-vs2022-setup.md) | C++ 扩展开发 |

---

## 快速导航

### 按角色选择

| 角色 | 推荐文档 |
|------|----------|
| 后端开发（Go） | [02-goland-setup.md](./development/02-goland-setup.md) |
| 后端开发（Python） | [03-pycharm-setup.md](./development/03-pycharm-setup.md) |
| 前端开发 | [04-vscode-setup.md](./development/04-vscode-setup.md) |
| 全栈开发 | [04-vscode-setup.md](./development/04-vscode-setup.md) |

### 按任务选择

| 任务 | 推荐文档 |
|------|----------|
| 首次配置开发环境 | [01-overview.md](./development/01-overview.md) |
| 调试 Go 服务 | [02-goland-setup.md](./development/02-goland-setup.md) |
| 调试 Python 服务 | [03-pycharm-setup.md](./development/03-pycharm-setup.md) |
| 前端独立开发 | [04-vscode-setup.md](./development/04-vscode-setup.md) |
| 性能分析 | [05-vs2022-setup.md](./development/05-vs2022-setup.md) |

---

## 项目结构

```
novelai-platform/
├── go-gateway/          # Go 网关服务
├── python-ai-service/   # Python AI 服务
├── frontend/            # React 前端
├── proto/               # Protocol Buffers
├── keys/                # JWT 密钥
├── docs/                # 文档目录
│   ├── README.md        # 本文件
│   └── development/     # 开发环境配置
└── docker-compose.yml   # 容器编排
```

---

## 贡献指南

### 文档命名规范

- 使用两位数字编号：`01-`, `02-`, ...
- 使用小写字母和连字符：`setup-guide.md`
- 按主题分类存放

### 文档格式

- 使用 Markdown 格式
- 包含目录（如文档较长）
- 提供代码示例
- 列出常见问题
