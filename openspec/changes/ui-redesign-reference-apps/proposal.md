# UI重新设计提案 - 参考MaliangAI和MuMuAI

## Why

当前NovelAI Platform前端缺少完整的项目管理界面和编辑器布局，用户无法有效管理多个项目或使用AI辅助创作功能。通过参考MaliangAINovalWriter和MuMuAINovel两个成熟应用的UI/UX模式，我们可以快速构建一个功能完整、用户体验优秀的前端界面，满足小说创作者的核心需求。

## What Changes

**前端新增功能**：
- 实现卡片式项目管理界面（ProjectList、ProjectDetail）
- 实现抽屉式AI面板编辑器（EditorLayout、AIChatPanel）
- 实现项目创建向导流程（ProjectWizardNew）
- 实现主题系统（亮色/暗色模式切换）
- 实现快捷键系统（编辑、导航、面板控制）
- 实现沉浸模式（专注写作）
- 实现高级搜索/筛选功能

**后端新增功能**：
- 实现用户布局配置API（持久化面板状态）
- 扩展项目元数据（封面、标签、统计信息）
- 实现文件上传API（项目封面、角色头像）
- 实现章节批量操作API（拖拽排序、批量更新）
- 实现全局搜索API（跨项目搜索）
- 实现高级筛选API（多条件筛选）

**技术栈约束**：
- 前端：React 18 + TypeScript + Ant Design 5 + Zustand + TailwindCSS
- 后端：Go Gateway + Python AI Service（现有架构）
- 拖拽：react-beautiful-dnd
- 富文本：TipTap（已有依赖）

## Capabilities

### New Capabilities

- `project-management-ui`: 项目管理界面，包括卡片式项目列表、项目详情布局、项目创建向导
- `drawer-editor-layout`: 抽屉式编辑器布局，包括主编辑区、抽屉式AI面板、工具栏
- `theme-system`: 主题系统，支持亮色/暗色模式切换、自定义颜色、主题持久化
- `keyboard-shortcuts`: 快捷键系统，支持编辑、导航、面板控制快捷键
- `immersive-mode`: 沉浸模式，隐藏所有UI仅显示编辑区
- `advanced-search-filter`: 高级搜索/筛选，支持全局搜索、多条件筛选、搜索历史
- `user-layout-config-api`: 用户布局配置API，持久化面板状态和用户偏好
- `project-metadata-extension`: 项目元数据扩展，支持封面、标签、统计信息
- `file-upload-api`: 文件上传API，支持图片上传（封面、头像）
- `batch-operations-api`: 批量操作API，支持章节排序、批量更新、批量删除

### Modified Capabilities

<!-- 无现有能力需要修改 -->

## Impact

**前端影响**：
- 新增页面：ProjectList.tsx、ProjectDetail.tsx、ProjectWizardNew.tsx、EditorLayout.tsx
- 新增组件：AIChatPanel.tsx、ThemeSwitcher.tsx、KeyboardShortcutsHelp.tsx、SearchBar.tsx
- 修改组件：AppLayout.tsx（集成主题系统）
- 新增Store：useEditorLayoutStore.ts、useThemeStore.ts、useSearchStore.ts
- 新增依赖：react-beautiful-dnd

**后端影响**：
- Go Gateway新增模型：UserLayoutConfig、ProjectMetadata、FileUpload
- Go Gateway新增端点：10+ 新API端点
- Go Gateway新增服务：UserLayoutService、FileUploadService、SearchService
- Python AI Service无影响（仅前端调用现有API）

**数据库影响**：
- 新增表：user_layout_configs、files
- 扩展表：projects（新增metadata字段）、chapters（新增tags字段）

**兼容性**：
- 向后兼容：所有新字段有默认值，不影响现有功能
- API版本：保持v1，新端点不影响现有端点

**工作量估算**：
- 第一阶段（项目管理UI）：30-40小时
- 第二阶段（编辑器布局）：20-30小时
- 第三阶段（增强功能）：30-40小时
- 第四阶段（后端API）：40-50小时
- 总计：120-160小时（3-4周）
