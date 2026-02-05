# Specification: Drawer Editor Layout

## ADDED Requirements

### Requirement: Main Editor Area
系统 SHALL 提供主编辑区域，占据屏幕中央，支持富文本编辑、工具栏、字数统计。

#### Scenario: Display editor
- **WHEN** 用户打开章节编辑页面
- **THEN** 系统显示富文本编辑器，包含工具栏和编辑区

#### Scenario: Toolbar actions
- **WHEN** 用户点击工具栏按钮（加粗、斜体、标题等）
- **THEN** 系统应用对应的文本格式

#### Scenario: Word count display
- **WHEN** 用户输入文本
- **THEN** 系统实时更新字数统计显示

#### Scenario: Auto-save
- **WHEN** 用户停止输入超过3秒
- **THEN** 系统自动保存章节内容

### Requirement: AI Chat Drawer
系统 SHALL 提供抽屉式AI聊天面板，从右侧滑出，支持多标签页（生成、重写、润色）。

#### Scenario: Open AI drawer
- **WHEN** 用户点击工具栏的"AI助手"按钮
- **THEN** 系统从右侧滑出AI聊天抽屉，默认宽度480px

#### Scenario: Switch AI tabs
- **WHEN** 用户点击抽屉内的标签页（生成、重写、润色）
- **THEN** 系统切换到对应的AI功能界面

#### Scenario: Close AI drawer
- **WHEN** 用户点击抽屉的关闭按钮或点击遮罩层
- **THEN** 系统关闭AI抽屉，恢复全屏编辑区

#### Scenario: Resize drawer
- **WHEN** 用户拖拽抽屉边缘
- **THEN** 系统调整抽屉宽度（280px-600px范围内）

### Requirement: AI Generation Interface
系统 SHALL 在AI抽屉内提供生成界面，支持选择模型、提供商、输入提示词、流式生成。

#### Scenario: Select AI provider
- **WHEN** 用户在AI抽屉中选择提供商（OpenAI、Anthropic、Gemini）
- **THEN** 系统更新可用模型列表

#### Scenario: Select AI model
- **WHEN** 用户选择模型（GPT-4o、Claude 3 Opus等）
- **THEN** 系统记录用户选择并用于后续生成

#### Scenario: Input prompt
- **WHEN** 用户在提示词输入框输入内容
- **THEN** 系统启用"生成"按钮

#### Scenario: Stream generation
- **WHEN** 用户点击"生成"按钮
- **THEN** 系统通过SSE流式接收并逐字显示AI生成的内容

#### Scenario: Abort generation
- **WHEN** 用户在生成过程中点击"停止"按钮
- **THEN** 系统中止SSE连接并停止生成

#### Scenario: Insert generated content
- **WHEN** 用户点击生成结果的"插入"按钮
- **THEN** 系统将生成内容插入到编辑器光标位置

### Requirement: Editor Toolbar
系统 SHALL 提供编辑器顶部工具栏，包含保存、撤销、重做、AI助手、设置等按钮。

#### Scenario: Save button
- **WHEN** 用户点击"保存"按钮
- **THEN** 系统立即保存章节内容并显示成功提示

#### Scenario: Undo and redo
- **WHEN** 用户点击"撤销"或"重做"按钮
- **THEN** 系统撤销或重做最近的编辑操作

#### Scenario: AI assistant button
- **WHEN** 用户点击"AI助手"按钮
- **THEN** 系统打开AI聊天抽屉

#### Scenario: Settings button
- **WHEN** 用户点击"设置"按钮
- **THEN** 系统打开编辑器设置对话框

### Requirement: Chapter Navigation Sidebar
系统 SHALL 提供左侧章节导航栏，显示项目的所有章节，支持快速切换。

#### Scenario: Display chapter list
- **WHEN** 用户在编辑器页面
- **THEN** 系统在左侧显示章节列表，当前章节高亮

#### Scenario: Switch chapter
- **WHEN** 用户点击章节列表中的某个章节
- **THEN** 系统加载该章节内容到编辑器

#### Scenario: Collapse sidebar
- **WHEN** 用户点击侧边栏折叠按钮
- **THEN** 系统隐藏章节列表，仅显示折叠按钮

#### Scenario: Drag to reorder chapters
- **WHEN** 用户拖拽章节列表项
- **THEN** 系统重新排序章节并保存新顺序

### Requirement: Responsive Editor Layout
系统 SHALL 根据屏幕尺寸自动调整编辑器布局，移动端使用全屏编辑器。

#### Scenario: Desktop layout
- **WHEN** 用户在桌面端（>1024px）访问编辑器
- **THEN** 系统显示左侧章节导航 + 中央编辑区 + 右侧AI抽屉（按需）

#### Scenario: Tablet layout
- **WHEN** 用户在平板端（768px-1024px）访问编辑器
- **THEN** 系统隐藏章节导航，通过顶部按钮切换

#### Scenario: Mobile layout
- **WHEN** 用户在移动端（<768px）访问编辑器
- **THEN** 系统显示全屏编辑器，AI功能通过底部工具栏访问

### Requirement: Loading States
系统 SHALL 在加载章节内容或AI生成时显示加载状态。

#### Scenario: Loading chapter
- **WHEN** 用户切换章节
- **THEN** 系统在编辑区显示加载骨架屏

#### Scenario: AI generating
- **WHEN** AI正在生成内容
- **THEN** 系统在AI抽屉中显示流式输出动画

#### Scenario: Saving indicator
- **WHEN** 系统正在保存章节
- **THEN** 系统在工具栏显示"保存中..."提示
