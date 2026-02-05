# Specification: Keyboard Shortcuts

## ADDED Requirements

### Requirement: Editor Shortcuts
系统 SHALL 支持常用的文本编辑快捷键，包括加粗、斜体、撤销、重做等。

#### Scenario: Bold text
- **WHEN** 用户按下 Ctrl+B（Windows/Linux）或 Cmd+B（Mac）
- **THEN** 系统将选中文本设置为加粗格式

#### Scenario: Italic text
- **WHEN** 用户按下 Ctrl+I 或 Cmd+I
- **THEN** 系统将选中文本设置为斜体格式

#### Scenario: Undo
- **WHEN** 用户按下 Ctrl+Z 或 Cmd+Z
- **THEN** 系统撤销最近的编辑操作

#### Scenario: Redo
- **WHEN** 用户按下 Ctrl+Y 或 Cmd+Shift+Z
- **THEN** 系统重做最近撤销的操作

#### Scenario: Save
- **WHEN** 用户按下 Ctrl+S 或 Cmd+S
- **THEN** 系统保存当前章节内容

### Requirement: Navigation Shortcuts
系统 SHALL 支持快速导航快捷键，包括搜索、切换章节、打开设置等。

#### Scenario: Quick search
- **WHEN** 用户按下 Ctrl+K 或 Cmd+K
- **THEN** 系统打开快速搜索对话框

#### Scenario: Next chapter
- **WHEN** 用户按下 Ctrl+Right 或 Cmd+Right
- **THEN** 系统切换到下一章节

#### Scenario: Previous chapter
- **WHEN** 用户按下 Ctrl+Left 或 Cmd+Left
- **THEN** 系统切换到上一章节

#### Scenario: Go to chapter list
- **WHEN** 用户按下 Ctrl+Shift+C
- **THEN** 系统聚焦到章节列表

### Requirement: Panel Control Shortcuts
系统 SHALL 支持面板控制快捷键，包括打开/关闭AI抽屉、侧边栏等。

#### Scenario: Toggle AI drawer
- **WHEN** 用户按下 Ctrl+Shift+A
- **THEN** 系统打开或关闭AI聊天抽屉

#### Scenario: Toggle sidebar
- **WHEN** 用户按下 Ctrl+B
- **THEN** 系统折叠或展开左侧章节导航栏

#### Scenario: Toggle immersive mode
- **WHEN** 用户按下 Ctrl+Shift+F 或 F11
- **THEN** 系统进入或退出沉浸模式

### Requirement: AI Shortcuts
系统 SHALL 支持AI功能快捷键，包括快速生成、重写、润色。

#### Scenario: Quick generate
- **WHEN** 用户按下 Ctrl+Shift+G
- **THEN** 系统打开AI生成对话框

#### Scenario: Rewrite selection
- **WHEN** 用户选中文本并按下 Ctrl+Shift+R
- **THEN** 系统打开AI重写对话框，预填选中文本

#### Scenario: Polish selection
- **WHEN** 用户选中文本并按下 Ctrl+Shift+P
- **THEN** 系统打开AI润色对话框，预填选中文本

### Requirement: Shortcuts Help
系统 SHALL 提供快捷键帮助对话框，显示所有可用快捷键。

#### Scenario: Open shortcuts help
- **WHEN** 用户按下 Ctrl+/ 或 ?
- **THEN** 系统打开快捷键帮助对话框，显示所有快捷键列表

#### Scenario: Search shortcuts
- **WHEN** 用户在快捷键帮助对话框中输入搜索关键词
- **THEN** 系统过滤显示匹配的快捷键

#### Scenario: Close shortcuts help
- **WHEN** 用户按下 Esc 或点击关闭按钮
- **THEN** 系统关闭快捷键帮助对话框

### Requirement: Custom Shortcuts
系统 SHALL 允许用户自定义快捷键（可选功能）。

#### Scenario: Customize shortcut
- **WHEN** 用户在设置中点击某个快捷键并按下新的组合键
- **THEN** 系统更新该快捷键绑定

#### Scenario: Reset shortcuts
- **WHEN** 用户点击"恢复默认快捷键"
- **THEN** 系统重置所有快捷键为默认值

#### Scenario: Conflict detection
- **WHEN** 用户设置的快捷键与现有快捷键冲突
- **THEN** 系统显示警告并要求用户确认或修改

### Requirement: Platform-specific Shortcuts
系统 SHALL 根据操作系统自动使用正确的修饰键（Ctrl/Cmd）。

#### Scenario: Windows shortcuts
- **WHEN** 用户在Windows系统使用快捷键
- **THEN** 系统使用Ctrl作为主修饰键

#### Scenario: Mac shortcuts
- **WHEN** 用户在Mac系统使用快捷键
- **THEN** 系统使用Cmd作为主修饰键

#### Scenario: Display correct keys
- **WHEN** 系统显示快捷键提示
- **THEN** 系统根据操作系统显示正确的修饰键符号（⌘、⌃、⌥、⇧）
