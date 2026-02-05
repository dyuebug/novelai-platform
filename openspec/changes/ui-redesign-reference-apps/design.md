# Design Document: UI Redesign - Reference Apps

## Context

### Background
当前NovelAI Platform前端缺少完整的项目管理界面和编辑器布局，用户体验不完整。通过深入分析MaliangAINovalWriter（Flutter）和MuMuAINovel（React）两个成熟应用，我们识别了关键的UI/UX模式和技术实现方案。

### Current State
**前端现状**：
- 技术栈：React 18 + TypeScript + Ant Design 5 + Zustand + TailwindCSS
- 现有页面：Settings、ForeshadowManager
- 现有布局：AppLayout（基础响应式设计）
- 缺失：项目管理界面、编辑器布局、主题系统、快捷键系统

**后端现状**：
- Go Gateway：60+ REST API端点，完善的认证和版本控制
- Python AI Service：AI生成、质量评估、约束检查
- 数据库：PostgreSQL + pgvector
- 通信：gRPC（Gateway ↔ AI Service）

### Constraints
1. **技术栈约束**：必须使用React 18 + Ant Design 5 + Zustand
2. **向后兼容**：不能破坏现有API和数据结构
3. **性能约束**：页面加载时间 < 2秒，API响应时间 < 300ms
4. **时间约束**：3-4周完成（120-160小时）
5. **资源约束**：单人开发，需要优先级排序

### Stakeholders
- **用户**：小说创作者，需要高效的创作工具
- **开发者**：需要清晰的架构和可维护的代码
- **运维**：需要稳定的部署和监控

---

## Goals / Non-Goals

### Goals
1. **完整的项目管理界面**：卡片式项目列表、项目详情布局、项目创建向导
2. **抽屉式编辑器布局**：主编辑区 + 抽屉式AI面板，适合中小屏幕
3. **主题系统**：亮色/暗色模式切换，持久化用户偏好
4. **快捷键系统**：提升编辑效率，支持键盘导航
5. **沉浸模式**：专注写作，隐藏所有UI
6. **高级搜索/筛选**：全局搜索、多条件筛选
7. **后端API扩展**：支持新UI功能的API端点

### Non-Goals
1. **多面板编辑器**：不实现MaliangAI的5面板并排布局（复杂度高，用户选择了抽屉式）
2. **实时协作**：不实现WebSocket协作功能（超出当前范围）
3. **移动端原生应用**：仅实现响应式Web应用
4. **AI模型训练**：不涉及AI模型本身的改进
5. **数据库迁移工具**：不实现自动化迁移工具（手动执行Alembic）

---

## Decisions

### Decision 1: 抽屉式 vs 多面板编辑器

**选择**：抽屉式AI面板（类似MuMuAI）

**理由**：
- 用户明确选择抽屉式布局
- 实现复杂度低（20-30小时 vs 50-70小时）
- 更适合中小屏幕（1080p及以下）
- Ant Design Drawer组件成熟稳定

**替代方案**：
- 多面板编辑器（类似MaliangAI）：功能强大但实现复杂，需要自定义拖拽分隔条、面板管理器
- 标签页式AI面板：切换不直观，无法同时查看编辑区和AI输出

**权衡**：
- ✅ 实现快速，用户体验良好
- ❌ 无法同时显示多个AI面板（但用户可以快速切换标签页）

---

### Decision 2: Zustand vs Redux 状态管理

**选择**：Zustand

**理由**：
- 已有技术栈，无需引入新依赖
- 轻量级（~1KB），性能优秀
- API简洁，学习曲线低
- 支持中间件（persist、devtools）

**替代方案**：
- Redux Toolkit：功能强大但样板代码多，过度设计
- Context API：适合简单场景，但性能不如Zustand

**权衡**：
- ✅ 开发效率高，代码简洁
- ❌ 生态不如Redux丰富（但对本项目足够）

---

### Decision 3: react-beautiful-dnd vs react-dnd 拖拽库

**选择**：react-beautiful-dnd

**理由**：
- 专为列表拖拽设计，API简单
- 内置平滑动画和无障碍支持
- 社区活跃，文档完善
- 适合章节排序、面板重排序场景

**替代方案**：
- react-dnd：更灵活但API复杂，适合复杂拖拽场景（如画布编辑器）
- 自己实现：工作量大，容易出bug

**权衡**：
- ✅ 开箱即用，用户体验好
- ❌ 库体积较大（~50KB），但可接受

---

### Decision 4: localStorage vs 后端API 持久化布局配置

**选择**：localStorage + 后端API（双重持久化）

**理由**：
- localStorage：即时响应，无需网络请求
- 后端API：跨设备同步，数据安全
- 双重持久化：localStorage优先，后台异步同步到后端

**替代方案**：
- 仅localStorage：无法跨设备同步
- 仅后端API：每次加载需要网络请求，延迟高

**权衡**：
- ✅ 最佳用户体验，支持跨设备
- ❌ 实现复杂度略高，需要同步逻辑

---

### Decision 5: Ant Design ConfigProvider vs CSS Variables 主题系统

**选择**：Ant Design ConfigProvider + CSS Variables

**理由**：
- ConfigProvider：控制所有Ant Design组件的主题
- CSS Variables：控制自定义组件和全局样式
- 结合使用：确保主题一致性

**替代方案**：
- 仅CSS Variables：无法控制Ant Design组件
- 仅ConfigProvider：无法控制自定义样式

**权衡**：
- ✅ 完整的主题控制，一致性好
- ❌ 需要维护两套主题配置

---

### Decision 6: TipTap vs Slate 富文本编辑器

**选择**：TipTap

**理由**：
- 已有依赖（package.json中已包含）
- 基于ProseMirror，性能优秀
- 扩展性强，支持自定义节点和标记
- 文档完善，社区活跃

**替代方案**：
- Slate：更灵活但API复杂，学习曲线陡峭
- Quill：简单但扩展性差

**权衡**：
- ✅ 功能强大，满足需求
- ❌ 学习曲线中等

---

### Decision 7: SSE vs WebSocket 流式生成

**选择**：SSE（Server-Sent Events）

**理由**：
- 后端已实现SSE接口（/api/wizard-stream/*）
- 单向通信足够（服务器→客户端）
- 实现简单，无需WebSocket服务器
- 自动重连机制

**替代方案**：
- WebSocket：双向通信，但本场景不需要
- 轮询：延迟高，资源浪费

**权衡**：
- ✅ 实现简单，性能好
- ❌ 不支持双向通信（但不需要）

---

### Decision 8: 文件存储：本地文件系统 vs 云存储

**选择**：本地文件系统（初期），支持云存储扩展

**理由**：
- 初期部署简单，无需配置云服务
- 适合Docker本地部署场景
- 预留云存储接口，后续可扩展

**替代方案**：
- 仅云存储：增加部署复杂度和成本
- 数据库BLOB：性能差，不适合大文件

**权衡**：
- ✅ 部署简单，成本低
- ❌ 无法跨服务器共享文件（但可通过云存储扩展）

---

### Decision 9: 图片处理：服务端 vs 客户端

**选择**：服务端处理（Go Gateway）

**理由**：
- 统一处理逻辑，确保一致性
- 减少客户端负担
- 支持缩略图生成、尺寸优化

**替代方案**：
- 客户端处理：增加客户端复杂度，不同浏览器兼容性问题

**权衡**：
- ✅ 一致性好，用户体验好
- ❌ 增加服务器负载（但可接受）

---

### Decision 10: 实现顺序：先项目管理，后编辑器

**选择**：Phase 1 项目管理 → Phase 2 编辑器 → Phase 3 增强功能 → Phase 4 后端API

**理由**：
- 用户明确选择此顺序
- 项目管理是入口，优先级高
- 编辑器依赖项目管理的数据结构
- 增强功能可以逐步添加

**替代方案**：
- 先编辑器：用户无法管理项目
- 并行开发：集成复杂度高

**权衡**：
- ✅ 逐步交付，风险可控
- ❌ 用户需要等待完整功能

---

## Architecture

### Frontend Architecture

```
frontend/
├── src/
│   ├── pages/                      # 页面组件
│   │   ├── ProjectList.tsx         # 项目列表（卡片式）
│   │   ├── ProjectDetail.tsx       # 项目详情布局
│   │   ├── ProjectWizardNew.tsx    # 项目创建向导
│   │   ├── EditorLayout.tsx        # 编辑器布局
│   │   └── Settings.tsx            # 设置页面（已存在）
│   ├── components/                 # 通用组件
│   │   ├── editor/
│   │   │   ├── AIChatPanel.tsx     # AI聊天抽屉
│   │   │   ├── RichTextEditor.tsx  # 富文本编辑器
│   │   │   ├── EditorToolbar.tsx   # 编辑器工具栏
│   │   │   └── ChapterSidebar.tsx  # 章节导航栏
│   │   ├── project/
│   │   │   ├── ProjectCard.tsx     # 项目卡片
│   │   │   ├── ProjectGrid.tsx     # 项目网格
│   │   │   └── ProjectFilters.tsx  # 项目筛选器
│   │   ├── theme/
│   │   │   ├── ThemeSwitcher.tsx   # 主题切换器
│   │   │   └── ThemeProvider.tsx   # 主题提供者
│   │   └── search/
│   │       ├── GlobalSearch.tsx    # 全局搜索
│   │       └── SearchBar.tsx       # 搜索栏
│   ├── store/                      # Zustand状态管理
│   │   ├── useEditorLayoutStore.ts # 编辑器布局状态
│   │   ├── useThemeStore.ts        # 主题状态
│   │   ├── useSearchStore.ts       # 搜索状态
│   │   └── useProjectStore.ts      # 项目状态
│   ├── hooks/                      # 自定义Hooks
│   │   ├── useKeyboardShortcuts.ts # 快捷键Hook
│   │   ├── useImmersiveMode.ts     # 沉浸模式Hook
│   │   └── useSSE.ts               # SSE流式Hook
│   └── services/                   # API服务
│       ├── api.ts                  # API客户端（已存在）
│       ├── projectService.ts       # 项目API
│       ├── layoutService.ts        # 布局配置API
│       └── uploadService.ts        # 文件上传API
```

### Backend Architecture

```
go-gateway/
├── internal/
│   ├── model/
│   │   ├── user_layout.go          # 用户布局配置模型
│   │   ├── file.go                 # 文件元数据模型
│   │   └── project.go              # 项目模型（扩展metadata）
│   ├── repository/
│   │   ├── user_layout_repository.go
│   │   ├── file_repository.go
│   │   └── project_repository.go   # 扩展方法
│   ├── service/
│   │   ├── user_layout_service.go
│   │   ├── file_upload_service.go
│   │   ├── search_service.go
│   │   └── batch_operation_service.go
│   ├── handler/
│   │   ├── user_layout_handler.go
│   │   ├── file_upload_handler.go
│   │   ├── search_handler.go
│   │   └── batch_operation_handler.go
│   └── router/
│       └── router.go               # 新增路由
```

### Data Flow

```
用户操作 → React组件 → Zustand Store → API Service → Go Gateway → 数据库/AI Service
                ↓                                              ↓
           localStorage持久化                            后端持久化
```

---

## Risks / Trade-offs

### Risk 1: 性能问题 - 大量项目加载慢
**风险**：用户有100+项目时，项目列表加载缓慢

**缓解措施**：
- 实现虚拟滚动（react-window）
- 分页加载（每页20个项目）
- 图片懒加载
- 缓存项目列表（5分钟）

---

### Risk 2: 浏览器兼容性 - 旧浏览器不支持
**风险**：IE11、旧版Safari不支持某些特性

**缓解措施**：
- 明确支持的浏览器版本（Chrome 90+, Firefox 88+, Safari 14+, Edge 90+）
- 使用Babel polyfill
- 优雅降级（如不支持全屏API时提示用户）

---

### Risk 3: 数据迁移 - 现有项目缺少元数据
**风险**：现有项目没有metadata字段，导致显示异常

**缓解措施**：
- 数据库迁移脚本自动添加metadata字段（默认值）
- 前端兼容null/undefined metadata
- 后端API始终返回metadata（即使为空对象）

---

### Risk 4: 文件上传安全 - 恶意文件上传
**风险**：用户上传恶意文件（病毒、脚本）

**缓解措施**：
- 严格验证文件类型（MIME type + 文件头）
- 限制文件大小（图片5MB）
- 文件扫描（ClamAV或云服务）
- 存储隔离（上传文件不在Web根目录）

---

### Risk 5: 快捷键冲突 - 与浏览器快捷键冲突
**风险**：自定义快捷键与浏览器快捷键冲突

**缓解措施**：
- 避免使用常见浏览器快捷键（Ctrl+T, Ctrl+W等）
- 提供快捷键自定义功能
- 显示快捷键帮助（Ctrl+/）

---

### Risk 6: SSE连接中断 - 网络不稳定
**风险**：SSE连接中断导致生成失败

**缓解措施**：
- 自动重连机制（最多3次）
- 显示连接状态
- 提供手动重试按钮
- 保存部分生成结果

---

### Risk 7: 状态同步 - localStorage与后端不一致
**风险**：localStorage和后端API的布局配置不一致

**缓解措施**：
- localStorage优先（即时响应）
- 后台异步同步到后端（debounce 3秒）
- 定期从后端拉取最新配置（每次登录）
- 冲突时以后端为准

---

### Risk 8: 批量操作失败 - 部分成功部分失败
**风险**：批量操作中部分章节更新失败

**缓解措施**：
- 使用数据库事务（全部成功或全部回滚）
- 返回详细的错误信息（哪些成功、哪些失败）
- 提供重试机制

---

## Migration Plan

### Phase 1: 项目管理UI（第1-2周）

**步骤**：
1. 创建ProjectList、ProjectCard、ProjectGrid组件
2. 创建ProjectDetail布局和路由
3. 创建ProjectWizardNew向导
4. 集成现有项目API
5. 测试响应式布局

**验收标准**：
- 用户可以查看项目列表（卡片式）
- 用户可以创建新项目（向导流程）
- 用户可以进入项目详情页面
- 移动端布局正常

---

### Phase 2: 编辑器布局（第2-3周）

**步骤**：
1. 创建EditorLayout、RichTextEditor组件
2. 创建AIChatPanel抽屉组件
3. 集成TipTap富文本编辑器
4. 实现SSE流式生成UI
5. 创建EditorToolbar和ChapterSidebar
6. 测试编辑器功能

**验收标准**：
- 用户可以编辑章节内容
- 用户可以打开AI抽屉并生成内容
- 用户可以切换章节
- 自动保存功能正常

---

### Phase 3: 增强功能（第3周）

**步骤**：
1. 实现主题系统（ThemeProvider、ThemeSwitcher）
2. 实现快捷键系统（useKeyboardShortcuts Hook）
3. 实现沉浸模式（useImmersiveMode Hook）
4. 实现全局搜索（GlobalSearch组件）
5. 测试所有增强功能

**验收标准**：
- 用户可以切换亮色/暗色主题
- 用户可以使用快捷键操作
- 用户可以进入沉浸模式
- 用户可以全局搜索内容

---

### Phase 4: 后端API（第4周）

**步骤**：
1. 实现用户布局配置API
2. 实现文件上传API
3. 扩展项目元数据API
4. 实现批量操作API
5. 实现搜索API
6. 数据库迁移
7. API测试

**验收标准**：
- 所有新API端点正常工作
- 数据库迁移成功
- API性能满足要求（<300ms）
- 单元测试覆盖率>80%

---

### Rollback Strategy

**如果Phase 1失败**：
- 回滚前端代码到主分支
- 无后端变更，无需回滚

**如果Phase 4失败**：
- 回滚数据库迁移（Alembic downgrade）
- 回滚Go Gateway代码
- 前端降级使用现有API

---

## Open Questions

### Q1: 是否需要实现项目模板功能？
**状态**：待决定
**影响**：如果需要，增加10-15小时工作量
**建议**：Phase 5实现（当前范围外）

### Q2: 文件上传是否需要支持视频/音频？
**状态**：待决定
**影响**：如果需要，增加文件处理逻辑和存储空间
**建议**：初期仅支持图片，后续扩展

### Q3: 是否需要实现协作功能（多用户编辑）？
**状态**：明确不需要（Non-Goals）
**影响**：无
**建议**：未来版本考虑

### Q4: 主题系统是否需要支持自定义颜色？
**状态**：待决定
**影响**：如果需要，增加5-8小时工作量
**建议**：Phase 3可选实现

### Q5: 搜索功能是否需要支持全文搜索（Elasticsearch）？
**状态**：待决定
**影响**：如果需要，增加部署复杂度和成本
**建议**：初期使用PostgreSQL全文搜索，后续扩展

---

## Success Metrics

### 用户体验指标
- 页面加载时间 < 2秒
- API响应时间 < 300ms
- 搜索响应时间 < 500ms
- 主题切换延迟 < 200ms

### 功能完整性指标
- 10个能力规范100%实现
- 单元测试覆盖率 > 80%
- E2E测试覆盖核心流程

### 代码质量指标
- ESLint无错误
- TypeScript无类型错误
- 代码审查通过

---

## Timeline

| 阶段 | 时间 | 工作量 | 交付物 |
|------|------|--------|--------|
| Phase 1 | 第1-2周 | 30-40小时 | 项目管理UI |
| Phase 2 | 第2-3周 | 20-30小时 | 编辑器布局 |
| Phase 3 | 第3周 | 30-40小时 | 增强功能 |
| Phase 4 | 第4周 | 40-50小时 | 后端API |
| **总计** | **3-4周** | **120-160小时** | **完整UI重新设计** |
