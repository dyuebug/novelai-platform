# Specification: Immersive Mode

## ADDED Requirements

### Requirement: Enter Immersive Mode
系统 SHALL 支持进入沉浸模式，隐藏所有UI元素，仅显示编辑区域。

#### Scenario: Enter via shortcut
- **WHEN** 用户按下 Ctrl+Shift+F 或 F11
- **THEN** 系统隐藏顶部工具栏、左侧导航栏、右侧AI抽屉，编辑区占据全屏

#### Scenario: Enter via button
- **WHEN** 用户点击工具栏的"沉浸模式"按钮
- **THEN** 系统进入沉浸模式

#### Scenario: Smooth transition
- **WHEN** 系统进入沉浸模式
- **THEN** 系统在300ms内平滑过渡，UI元素淡出

### Requirement: Exit Immersive Mode
系统 SHALL 支持退出沉浸模式，恢复所有UI元素。

#### Scenario: Exit via shortcut
- **WHEN** 用户在沉浸模式下按下 Ctrl+Shift+F 或 F11 或 Esc
- **THEN** 系统退出沉浸模式，恢复所有UI元素

#### Scenario: Exit via hover
- **WHEN** 用户在沉浸模式下将鼠标移动到屏幕顶部
- **THEN** 系统显示半透明的退出按钮

#### Scenario: Exit button click
- **WHEN** 用户点击退出按钮
- **THEN** 系统退出沉浸模式

### Requirement: Immersive Mode UI
系统 SHALL 在沉浸模式下提供最小化的UI，包括字数统计、退出按钮。

#### Scenario: Display word count
- **WHEN** 用户在沉浸模式下编辑
- **THEN** 系统在屏幕底部中央显示半透明的字数统计

#### Scenario: Hide word count
- **WHEN** 用户停止输入超过5秒
- **THEN** 系统淡出字数统计显示

#### Scenario: Show on typing
- **WHEN** 用户在沉浸模式下开始输入
- **THEN** 系统淡入字数统计显示

### Requirement: Auto-save in Immersive Mode
系统 SHALL 在沉浸模式下继续自动保存章节内容。

#### Scenario: Auto-save
- **WHEN** 用户在沉浸模式下停止输入超过3秒
- **THEN** 系统自动保存章节内容

#### Scenario: Save indicator
- **WHEN** 系统在沉浸模式下保存内容
- **THEN** 系统在屏幕右下角显示短暂的"已保存"提示

### Requirement: Immersive Mode Persistence
系统 SHALL 记住用户的沉浸模式偏好设置。

#### Scenario: Remember preference
- **WHEN** 用户进入沉浸模式
- **THEN** 系统将偏好保存到localStorage

#### Scenario: Auto-enter on load
- **WHEN** 用户设置了"始终以沉浸模式打开"并加载章节
- **THEN** 系统自动进入沉浸模式

#### Scenario: Disable auto-enter
- **WHEN** 用户在设置中关闭"始终以沉浸模式打开"
- **THEN** 系统不再自动进入沉浸模式

### Requirement: Fullscreen Support
系统 SHALL 支持浏览器全屏API，提供真正的全屏体验。

#### Scenario: Request fullscreen
- **WHEN** 用户在沉浸模式下按下 F11
- **THEN** 系统请求浏览器全屏，隐藏浏览器UI

#### Scenario: Exit fullscreen
- **WHEN** 用户在全屏模式下按下 Esc
- **THEN** 系统退出浏览器全屏，保持沉浸模式

#### Scenario: Fullscreen denied
- **WHEN** 浏览器拒绝全屏请求
- **THEN** 系统显示提示信息，继续保持沉浸模式（非全屏）

### Requirement: Distraction-free Writing
系统 SHALL 在沉浸模式下禁用所有通知和弹窗。

#### Scenario: Suppress notifications
- **WHEN** 用户在沉浸模式下
- **THEN** 系统不显示任何通知提示

#### Scenario: Queue notifications
- **WHEN** 系统在沉浸模式下收到通知
- **THEN** 系统将通知加入队列，退出沉浸模式后显示

#### Scenario: Disable popups
- **WHEN** 用户在沉浸模式下
- **THEN** 系统不显示任何对话框或弹窗（除非用户主动触发）
