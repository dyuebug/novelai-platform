# Implementation Tasks: UI Redesign - Reference Apps

## 1. 项目管理UI - 基础组件

- [x] 1.1 创建ProjectCard组件（卡片式项目展示，包含封面、标题、描述、标签、统计信息）
- [x] 1.2 创建ProjectGrid组件（响应式网格布局，支持3-4列/2列/1列）
- [x] 1.3 创建ProjectFilters组件（搜索框、标签筛选、状态筛选、日期筛选）
- [x] 1.4 创建EmptyState组件（空状态提示和创建按钮）
- [x] 1.5 创建ProjectActions组件（编辑、删除、查看等快速操作）

## 2. 项目管理UI - 页面实现

- [x] 2.1 创建ProjectList页面（集成ProjectGrid、ProjectFilters、分页）
- [x] 2.2 实现项目搜索功能（实时搜索、防抖处理）
- [x] 2.3 实现项目筛选功能（多条件组合筛选）
- [x] 2.4 实现项目排序功能（按创建时间、更新时间、标题排序）
- [x] 2.5 实现批量选择模式（复选框、批量操作工具栏）
- [x] 2.6 实现批量删除功能（确认对话框、API调用）
- [x] 2.7 实现批量导出功能（生成压缩包）
- [x] 2.8 实现批量标签管理（添加/移除标签）

## 3. 项目详情布局

- [x] 3.1 创建ProjectDetail页面（侧边栏 + 主内容区布局）
- [x] 3.2 实现侧边栏导航（章节、角色、地点、组织、世界设定、伏笔菜单）
- [x] 3.3 实现侧边栏折叠功能（220px ↔ 80px）
- [x] 3.4 实现移动端Drawer导航（替代固定侧边栏）
- [x] 3.5 实现嵌套路由（子页面渲染）
- [x] 3.6 实现面包屑导航（显示当前位置）

## 4. 项目创建向导

- [x] 4.1 创建ProjectWizardNew组件（多步表单）
- [x] 4.2 实现第一步：基本信息表单（标题、描述、类型、标签）
- [x] 4.3 实现第二步：AI生成设定（提示词输入、模型选择）
- [x] 4.4 实现SSE流式生成UI（进度显示、逐字渲染）
- [x] 4.5 实现生成中止功能（停止按钮、中断SSE连接）
- [x] 4.6 实现向导导航（上一步、下一步、取消）
- [x] 4.7 实现表单验证（必填字段、格式验证）
- [x] 4.8 实现项目保存（API调用、成功跳转）
- [x] 4.9 实现恢复未完成向导（localStorage缓存、提示恢复）

## 5. 编辑器布局 - 基础组件

- [x] 5.1 创建EditorLayout组件（主编辑区 + 工具栏 + 侧边栏）
- [x] 5.2 创建EditorToolbar组件（保存、撤销、重做、AI助手、设置按钮）
- [x] 5.3 创建ChapterSidebar组件（章节列表、当前章节高亮）
- [x] 5.4 实现章节切换功能（点击加载章节内容）
- [x] 5.5 实现章节拖拽排序（react-beautiful-dnd集成）
- [x] 5.6 实现侧边栏折叠功能（隐藏章节列表）

## 6. 富文本编辑器

- [x] 6.1 集成TipTap编辑器（基础配置）
- [x] 6.2 实现编辑器工具栏（加粗、斜体、标题、列表等）
- [x] 6.3 实现字数统计（实时更新）
- [x] 6.4 实现自动保存（3秒防抖、API调用）
- [x] 6.5 实现保存状态提示（保存中、已保存、保存失败）
- [x] 6.6 实现撤销/重做功能（Ctrl+Z、Ctrl+Y）
- [x] 6.7 实现编辑器快捷键（加粗、斜体、保存等）

## 7. AI聊天抽屉

- [x] 7.1 创建AIChatPanel组件（Ant Design Drawer）
- [x] 7.2 实现抽屉打开/关闭（从右侧滑出）
- [x] 7.3 实现抽屉宽度调整（拖拽边缘，280px-600px）
- [x] 7.4 实现多标签页（生成、重写、润色）
- [x] 7.5 实现AI提供商选择（OpenAI、Anthropic、Gemini）
- [x] 7.6 实现AI模型选择（动态模型列表）
- [x] 7.7 实现提示词输入（多行文本框）
- [x] 7.8 实现SSE流式生成（逐字显示）
- [x] 7.9 实现生成中止功能（停止按钮）
- [x] 7.10 实现生成结果操作（插入、复制、重新生成）

## 8. 主题系统

- [x] 8.1 创建useThemeStore（Zustand store）
- [x] 8.2 创建ThemeProvider组件（Ant Design ConfigProvider + CSS Variables）
- [x] 8.3 定义亮色主题配置（颜色、字体、边框等）
- [x] 8.4 定义暗色主题配置（颜色、字体、边框等）
- [x] 8.5 创建ThemeSwitcher组件（亮色/暗色/跟随系统选项）
- [x] 8.6 实现主题切换功能（更新ConfigProvider和CSS Variables）
- [x] 8.7 实现主题过渡动画（200ms平滑过渡）
- [x] 8.8 实现主题持久化（localStorage保存）
- [x] 8.9 实现跟随系统主题（监听系统主题变化）
- [x] 8.10 集成主题到AppLayout（全局应用）

## 9. 快捷键系统

- [x] 9.1 创建useKeyboardShortcuts Hook（监听键盘事件）
- [x] 9.2 实现编辑快捷键（Ctrl+B、Ctrl+I、Ctrl+S等）
- [x] 9.3 实现导航快捷键（Ctrl+K、Ctrl+Left、Ctrl+Right等）
- [x] 9.4 实现面板控制快捷键（Ctrl+Shift+A、Ctrl+B等）
- [x] 9.5 实现AI快捷键（Ctrl+Shift+G、Ctrl+Shift+R等）
- [x] 9.6 创建KeyboardShortcutsHelp组件（快捷键列表对话框）
- [x] 9.7 实现快捷键帮助打开（Ctrl+/或?）
- [x] 9.8 实现快捷键搜索（过滤快捷键列表）
- [x] 9.9 实现平台检测（Windows用Ctrl，Mac用Cmd）
- [x] 9.10 实现快捷键冲突检测（避免与浏览器快捷键冲突）

## 10. 沉浸模式

- [x] 10.1 创建useImmersiveMode Hook（管理沉浸模式状态）
- [x] 10.2 实现进入沉浸模式（隐藏工具栏、侧边栏、AI抽屉）
- [x] 10.3 实现退出沉浸模式（恢复所有UI元素）
- [x] 10.4 实现沉浸模式快捷键（Ctrl+Shift+F、F11、Esc）
- [x] 10.5 实现沉浸模式UI（半透明字数统计、退出按钮）
- [x] 10.6 实现字数统计自动隐藏（5秒无输入后淡出）
- [x] 10.7 实现浏览器全屏支持（Fullscreen API）
- [x] 10.8 实现沉浸模式持久化（localStorage保存偏好）
- [x] 10.9 实现自动进入沉浸模式（可选设置）
- [x] 10.10 实现通知抑制（沉浸模式下不显示通知）

## 11. 全局搜索

- [x] 11.1 创建GlobalSearch组件（搜索对话框）
- [x] 11.2 实现搜索框（输入框、搜索图标）
- [x] 11.3 实现实时搜索（300ms防抖、API调用）
- [x] 11.4 实现搜索结果显示（按类型分组）
- [x] 11.5 实现搜索结果高亮（匹配关键词高亮）
- [x] 11.6 实现搜索结果导航（点击跳转到详情页）
- [x] 11.7 实现搜索历史（保存最近10条搜索）
- [x] 11.8 实现搜索建议（自动完成、相关搜索）
- [x] 11.9 实现搜索快捷键（Ctrl+K打开搜索）
- [x] 11.10 实现搜索性能优化（虚拟滚动、分页）

## 12. 高级筛选

- [x] 12.1 创建AdvancedFilter组件（筛选面板）
- [x] 12.2 实现标签筛选（多选标签）
- [x] 12.3 实现状态筛选（草稿、进行中、已完成）
- [x] 12.4 实现日期范围筛选（日期选择器）
- [x] 12.5 实现字数范围筛选（最小值、最大值输入）
- [x] 12.6 实现筛选条件组合（AND逻辑）
- [x] 12.7 实现筛选结果显示（实时更新）
- [x] 12.8 实现筛选条件重置（清除所有筛选）
- [x] 12.9 实现筛选条件持久化（localStorage保存）

## 13. 后端API - 用户布局配置

- [x] 13.1 创建UserLayoutConfig模型（Go struct + GORM标签）
- [x] 13.2 创建数据库迁移文件（user_layout_configs表）
- [x] 13.3 创建UserLayoutRepository（CRUD方法）
- [x] 13.4 创建UserLayoutService（业务逻辑）
- [x] 13.5 创建UserLayoutHandler（HTTP处理器）
- [x] 13.6 实现GET /api/v1/user/layout-config/:layoutType（获取配置）
- [x] 13.7 实现POST /api/v1/user/layout-config/:layoutType（保存配置）
- [x] 13.8 实现DELETE /api/v1/user/layout-config/:layoutType（删除配置）
- [x] 13.9 添加路由到router.go
- [x] 13.10 添加单元测试（覆盖率>80%）

## 14. 后端API - 文件上传

- [x] 14.1 创建File模型（Go struct + GORM标签）
- [x] 14.2 创建数据库迁移文件（files表）
- [x] 14.3 创建FileRepository（CRUD方法）
- [x] 14.4 创建FileUploadService（上传逻辑、图片处理）
- [x] 14.5 创建FileUploadHandler（HTTP处理器）
- [x] 14.6 实现POST /api/v1/upload（上传文件）
- [x] 14.7 实现GET /api/v1/files/:id（获取文件元数据）
- [x] 14.8 实现DELETE /api/v1/files/:id（删除文件）
- [x] 14.9 实现文件类型验证（MIME type + 文件头）
- [x] 14.10 实现文件大小限制（图片5MB）
- [x] 14.11 实现缩略图生成（200x200）
- [x] 14.12 实现图片优化（压缩质量80%）
- [x] 14.13 添加路由到router.go
- [ ] 14.14 添加单元测试（覆盖率>80%）

## 15. 后端API - 项目元数据扩展

- [x] 15.1 扩展Project模型（添加metadata字段）
- [x] 15.2 创建数据库迁移文件（添加metadata列）
- [x] 15.3 更新ProjectRepository（支持metadata查询）
- [x] 15.4 更新ProjectService（计算统计信息）
- [x] 15.5 实现PUT /api/v1/projects/:id/metadata（更新元数据）
- [x] 15.6 实现GET /api/v1/projects/:id/statistics（获取统计信息）
- [ ] 15.7 实现统计信息缓存（Redis，5分钟）
- [x] 15.8 更新现有项目API（返回metadata字段）
- [x] 15.9 添加向后兼容处理（旧项目默认metadata）
- [ ] 15.10 添加单元测试（覆盖率>80%）

## 16. 后端API - 批量操作

- [x] 16.1 创建BatchOperationService（批量操作逻辑）
- [x] 16.2 创建BatchOperationHandler（HTTP处理器）
- [x] 16.3 实现POST /api/v1/projects/:id/chapters/batch-update（批量更新）
- [x] 16.4 实现POST /api/v1/projects/:id/chapters/batch-delete（批量删除）
- [x] 16.5 实现POST /api/v1/projects/:id/chapters/batch-status（批量状态更新）
- [x] 16.6 实现POST /api/v1/projects/:id/chapters/reorder（重排序）
- [x] 16.7 实现事务支持（全部成功或全部回滚）
- [x] 16.8 实现批量操作验证（章节ID、数据格式）
- [x] 16.9 实现批量操作限制（最多100条）
- [x] 16.10 实现批量操作日志（审计记录）
- [x] 16.11 添加路由到router.go
- [ ] 16.12 添加单元测试（覆盖率>80%）

## 17. 后端API - 搜索功能

- [x] 17.1 创建SearchService（搜索逻辑）
- [x] 17.2 创建SearchHandler（HTTP处理器）
- [x] 17.3 实现GET /api/v1/search（全局搜索）
- [x] 17.4 实现GET /api/v1/projects/:id/search（项目内搜索）
- [x] 17.5 实现POST /api/v1/projects/:id/advanced-filter（高级筛选）
- [x] 17.6 实现PostgreSQL全文搜索（tsvector + tsquery）
- [x] 17.7 实现搜索结果高亮（返回匹配上下文）
- [x] 17.8 实现搜索分页（每页20条）
- [x] 17.9 实现搜索性能优化（索引、缓存）
- [x] 17.10 添加路由到router.go
- [ ] 17.11 添加单元测试（覆盖率>80%）

## 18. 前端状态管理

- [x] 18.1 创建useEditorLayoutStore（编辑器布局状态）
- [x] 18.2 创建useThemeStore（主题状态）
- [x] 18.3 创建useSearchStore（搜索状态）
- [x] 18.4 创建useProjectStore（项目状态）
- [x] 18.5 实现状态持久化中间件（localStorage）
- [x] 18.6 实现状态同步到后端（debounce 3秒）
- [x] 18.7 实现状态冲突解决（后端优先）

## 19. 前端API服务

- [x] 19.1 创建layoutService（布局配置API）
- [x] 19.2 创建uploadService（文件上传API）
- [x] 19.3 创建projectService（项目API，扩展metadata）
- [x] 19.4 创建searchService（搜索API）
- [x] 19.5 创建batchService（批量操作API）
- [x] 19.6 实现API错误处理（统一错误提示）
- [x] 19.7 实现API请求拦截器（添加JWT令牌）
- [x] 19.8 实现API响应拦截器（处理401、403等）

## 20. 响应式设计

- [ ] 20.1 实现移动端布局（<768px）
- [ ] 20.2 实现平板端布局（768px-1024px）
- [ ] 20.3 实现桌面端布局（>1024px）
- [ ] 20.4 实现触摸手势支持（移动端）
- [ ] 20.5 测试各种屏幕尺寸（Chrome DevTools）
- [ ] 20.6 测试各种设备（iPhone、iPad、Android）

## 21. 性能优化

- [x] 21.1 实现代码分割（React.lazy + Suspense）
- [ ] 21.2 实现虚拟滚动（react-window，项目列表）
- [ ] 21.3 实现图片懒加载（Intersection Observer）
- [ ] 21.4 实现组件记忆化（React.memo、useMemo、useCallback）
- [ ] 21.5 实现API响应缓存（React Query或SWR）
- [ ] 21.6 优化打包体积（Tree shaking、代码压缩）
- [ ] 21.7 优化首屏加载时间（<2秒）

## 22. 测试

- [ ] 22.1 编写ProjectCard组件测试
- [ ] 22.2 编写ProjectList页面测试
- [ ] 22.3 编写EditorLayout组件测试
- [ ] 22.4 编写AIChatPanel组件测试
- [ ] 22.5 编写ThemeProvider测试
- [ ] 22.6 编写useKeyboardShortcuts Hook测试
- [ ] 22.7 编写API服务测试
- [ ] 22.8 编写E2E测试（Playwright，核心流程）
- [ ] 22.9 确保单元测试覆盖率>80%
- [ ] 22.10 确保E2E测试覆盖核心流程

## 23. 文档

- [ ] 23.1 更新README.md（新功能说明）
- [ ] 23.2 编写用户指南（如何使用新UI）
- [ ] 23.3 编写开发者文档（组件API、状态管理）
- [ ] 23.4 编写API文档（新增端点、请求/响应格式）
- [ ] 23.5 编写部署文档（数据库迁移步骤）

## 24. 部署准备

- [ ] 24.1 执行数据库迁移（Alembic upgrade head）
- [ ] 24.2 构建前端生产版本（npm run build）
- [ ] 24.3 构建后端二进制文件（go build）
- [ ] 24.4 更新Docker镜像（docker compose build）
- [ ] 24.5 测试Docker部署（docker compose up -d）
- [ ] 24.6 验证所有服务健康（健康检查）
- [ ] 24.7 验证所有功能正常（手动测试）
- [ ] 24.8 准备回滚方案（数据库备份、代码备份）
