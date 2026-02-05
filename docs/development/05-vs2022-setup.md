# Visual Studio 2022 配置指南

## 适用场景

- C/C++ 原生扩展开发
- .NET 后端服务开发（如需扩展）
- CMake 项目构建
- 性能分析与调试

> **注意**: NovelAI Platform 主要使用 Go、Python、TypeScript 开发。VS2022 主要用于：
> - 开发高性能 C++ 扩展模块
> - Python C 扩展开发
> - 需要 Windows 原生 API 的场景

---

## 1. 安装工作负载

### 必需组件

启动 Visual Studio Installer，安装以下工作负载：

| 工作负载 | 用途 |
|----------|------|
| 使用 C++ 的桌面开发 | C++ 开发基础 |
| Python 开发 | Python 支持（可选） |
| Node.js 开发 | 前端开发（可选） |

### 单独组件

在 `单独组件` 标签中确保安装：

- CMake tools for Windows
- Git for Windows
- Windows 10/11 SDK

---

## 2. 打开项目

### 方式一：打开文件夹

1. 启动 VS2022
2. 选择 `打开本地文件夹`
3. 导航到 `novelai-platform` 目录
4. 点击 `选择文件夹`

### 方式二：CMake 项目

如果有 CMake 项目：

1. `文件` → `打开` → `CMake...`
2. 选择 `CMakeLists.txt`

---

## 3. 配置 CMake（C++ 扩展开发）

### 创建 CMakeLists.txt

如需开发 C++ 扩展，在项目中创建：

```cmake
cmake_minimum_required(VERSION 3.20)
project(novelai_extensions)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Python 扩展示例
find_package(Python3 COMPONENTS Interpreter Development REQUIRED)

# 添加扩展模块
add_library(fast_tokenizer SHARED
    src/tokenizer.cpp
)

target_include_directories(fast_tokenizer PRIVATE
    ${Python3_INCLUDE_DIRS}
)

target_link_libraries(fast_tokenizer PRIVATE
    ${Python3_LIBRARIES}
)
```

### CMake 设置

1. `项目` → `CMake 设置`
2. 配置：

| 设置 | 值 |
|------|-----|
| 配置类型 | Debug / Release |
| 工具集 | msvc_x64_x64 |
| CMake 生成器 | Ninja |

---

## 4. Python 开发支持

### 配置 Python 环境

1. `工具` → `选项` → `Python` → `环境`
2. 点击 `+ 添加环境`
3. 选择 `现有环境`
4. 路径：`python-ai-service\.venv\Scripts\python.exe`

### 打开 Python 项目

1. `文件` → `打开` → `文件夹`
2. 选择 `python-ai-service`
3. VS2022 会自动识别 Python 项目

---

## 5. 调试配置

### launch.vs.json

在 `.vs` 目录创建 `launch.vs.json`：

```json
{
  "version": "0.2.1",
  "configurations": [
    {
      "type": "python",
      "name": "Python: FastAPI",
      "request": "launch",
      "module": "uvicorn",
      "args": ["app.main:app", "--reload"],
      "cwd": "${workspaceRoot}/python-ai-service",
      "env": {
        "PYTHONPATH": "${workspaceRoot}/python-ai-service"
      }
    },
    {
      "type": "default",
      "name": "Node: Frontend",
      "project": "frontend/package.json",
      "projectTarget": "dev"
    }
  ]
}
```

### 调试 C++ 扩展

1. 设置断点
2. `调试` → `开始调试` (F5)
3. 选择目标可执行文件

---

## 6. 终端配置

### 开发者 PowerShell

1. `视图` → `终端`
2. 下拉选择 `Developer PowerShell`

### 配置多终端

点击终端窗口的 `+` 创建多个终端实例。

---

## 7. Git 集成

### 内置 Git 支持

1. `视图` → `Git 更改`
2. 查看更改、暂存、提交

### 分支管理

1. 点击状态栏的分支名
2. 创建/切换分支

---

## 8. 性能分析

### CPU 分析

1. `调试` → `性能探查器`
2. 选择 `CPU 使用率`
3. 点击 `开始`

### 内存分析

1. `调试` → `性能探查器`
2. 选择 `.NET 对象分配` 或 `内存使用率`

---

## 9. 扩展推荐

| 扩展 | 用途 |
|------|------|
| GitHub Copilot | AI 代码补全 |
| CodeMaid | 代码清理 |
| Productivity Power Tools | 效率工具集 |
| Markdown Editor | Markdown 支持 |

### 安装扩展

1. `扩展` → `管理扩展`
2. 搜索并安装

---

## 10. 与其他 IDE 协作

### 推荐工作流

由于 NovelAI Platform 主要使用 Go/Python/TypeScript：

| 任务 | 推荐工具 |
|------|----------|
| Go Gateway 开发 | GoLand 或 VSCode |
| Python AI Service | PyCharm 或 VSCode |
| React Frontend | VSCode |
| C++ 扩展开发 | **VS2022** |
| 性能分析 | **VS2022** |

### 项目文件共存

VS2022 的 `.vs` 目录不会影响其他 IDE：

```
novelai-platform/
├── .vs/              # VS2022 配置（已在 .gitignore）
├── .vscode/          # VSCode 配置
├── .idea/            # JetBrains IDE 配置
└── ...
```

---

## 常见问题

### Q: CMake 配置失败？

1. 确保安装了 CMake 工具
2. 检查 CMakeLists.txt 语法
3. `项目` → `删除缓存并重新配置`

### Q: Python 环境识别不到？

1. 确保虚拟环境已创建
2. 手动添加环境路径
3. 重启 VS2022

### Q: 编译 C++ 扩展报错？

检查：
1. Windows SDK 已安装
2. Python 开发头文件可用
3. 编译器版本兼容

### Q: 调试时找不到符号？

1. 确保使用 Debug 配置编译
2. 检查 PDB 文件生成
3. `工具` → `选项` → `调试` → `符号` 配置符号服务器
