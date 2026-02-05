# Specification: Project Management UI

## ADDED Requirements

### Requirement: Project List Display
系统 SHALL 以卡片式网格布局展示用户的所有项目，每个项目卡片包含封面图片、标题、描述、标签、统计信息和操作按钮。

#### Scenario: Display project cards
- **WHEN** 用户访问项目列表页面
- **THEN** 系统显示所有项目的卡片，每个卡片包含项目封面、标题、描述、标签、字数统计、章节数、最后更新时间

#### Scenario: Empty project list
- **WHEN** 用户没有任何项目
- **THEN** 系统显示空状态提示和"创建新项目"按钮

#### Scenario: Responsive grid layout
- **WHEN** 用户在不同屏幕尺寸下查看项目列表
- **THEN** 系统自动调整卡片网格列数（桌面端3-4列，平板端2列，移动端1列）

### Requirement: Project Card Actions
系统 SHALL 在每个项目卡片上提供快速操作按钮，包括编辑、删除、查看详情。

#### Scenario: Edit project
- **WHEN** 用户点击项目卡片的"编辑"按钮
- **THEN** 系统打开项目编辑对话框，允许修改项目信息

#### Scenario: Delete project
- **WHEN** 用户点击项目卡片的"删除"按钮
- **THEN** 系统显示确认对话框，确认后执行软删除

#### Scenario: View project details
- **WHEN** 用户点击项目卡片或"查看"按钮
- **THEN** 系统导航到项目详情页面

### Requirement: Project Search and Filter
系统 SHALL 提供项目搜索和筛选功能，支持按标题、标签、状态、日期范围筛选。

#### Scenario: Search by title
- **WHEN** 用户在搜索框输入关键词
- **THEN** 系统实时过滤显示标题包含关键词的项目

#### Scenario: Filter by tags
- **WHEN** 用户选择一个或多个标签
- **THEN** 系统仅显示包含所选标签的项目

#### Scenario: Filter by status
- **WHEN** 用户选择项目状态（草稿、进行中、已完成）
- **THEN** 系统仅显示对应状态的项目

### Requirement: Project Detail Layout
系统 SHALL 提供项目详情页面，包含侧边栏导航和主内容区，支持查看章节列表、角色管理、世界设定等子页面。

#### Scenario: Display project sidebar
- **WHEN** 用户进入项目详情页面
- **THEN** 系统显示左侧导航栏，包含章节、角色、地点、组织、世界设定、伏笔等菜单项

#### Scenario: Navigate to sub-pages
- **WHEN** 用户点击侧边栏菜单项
- **THEN** 系统在主内容区显示对应的子页面内容

#### Scenario: Collapsible sidebar
- **WHEN** 用户点击侧边栏折叠按钮
- **THEN** 系统将侧边栏宽度从220px缩小到80px，仅显示图标

#### Scenario: Mobile drawer navigation
- **WHEN** 用户在移动端访问项目详情页面
- **THEN** 系统使用抽屉式导航替代固定侧边栏

### Requirement: Project Creation Wizard
系统 SHALL 提供项目创建向导，引导用户通过多步表单创建新项目，支持AI辅助生成项目设定。

#### Scenario: Start wizard
- **WHEN** 用户点击"创建新项目"按钮
- **THEN** 系统打开项目创建向导，显示第一步表单（基本信息）

#### Scenario: Fill basic information
- **WHEN** 用户填写项目标题、描述、类型并点击"下一步"
- **THEN** 系统验证输入并进入第二步（AI生成设定）

#### Scenario: AI generate settings
- **WHEN** 用户在第二步点击"AI生成"按钮
- **THEN** 系统通过SSE流式生成项目世界观、角色、大纲等设定

#### Scenario: Save and create project
- **WHEN** 用户完成所有步骤并点击"创建项目"
- **THEN** 系统保存项目数据并导航到项目详情页面

#### Scenario: Resume incomplete wizard
- **WHEN** 用户关闭向导后再次打开
- **THEN** 系统提示是否恢复上次未完成的项目创建

### Requirement: Project Import and Export
系统 SHALL 支持项目导入和导出功能，支持JSON、Markdown、EPUB格式。

#### Scenario: Export project
- **WHEN** 用户点击项目的"导出"按钮并选择格式
- **THEN** 系统生成导出文件并触发下载

#### Scenario: Import project
- **WHEN** 用户点击"导入项目"按钮并选择文件
- **THEN** 系统解析文件并创建新项目

#### Scenario: Export progress
- **WHEN** 导出大型项目时
- **THEN** 系统显示导出进度条

### Requirement: Batch Operations
系统 SHALL 支持批量操作，允许用户同时选择多个项目进行删除、导出、标签管理。

#### Scenario: Select multiple projects
- **WHEN** 用户点击项目卡片的复选框
- **THEN** 系统进入批量选择模式，显示批量操作工具栏

#### Scenario: Batch delete
- **WHEN** 用户选择多个项目并点击"批量删除"
- **THEN** 系统显示确认对话框，确认后删除所有选中项目

#### Scenario: Batch export
- **WHEN** 用户选择多个项目并点击"批量导出"
- **THEN** 系统生成包含所有项目的压缩包并下载

#### Scenario: Batch tag management
- **WHEN** 用户选择多个项目并点击"管理标签"
- **THEN** 系统显示标签编辑对话框，允许批量添加或移除标签
