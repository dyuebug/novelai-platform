# Specification: Advanced Search and Filter

## ADDED Requirements

### Requirement: Global Search
系统 SHALL 提供全局搜索功能，支持跨项目搜索章节、角色、地点、世界设定等内容。

#### Scenario: Open global search
- **WHEN** 用户按下 Ctrl+K 或点击顶部搜索图标
- **THEN** 系统打开全局搜索对话框

#### Scenario: Search across projects
- **WHEN** 用户在搜索框输入关键词
- **THEN** 系统实时搜索所有项目的章节、角色、地点、世界设定，显示匹配结果

#### Scenario: Search result categories
- **WHEN** 系统显示搜索结果
- **THEN** 系统按类型分组显示（章节、角色、地点、世界设定），每个结果显示标题、摘要、所属项目

#### Scenario: Navigate to result
- **WHEN** 用户点击搜索结果
- **THEN** 系统导航到对应的详情页面并高亮匹配内容

### Requirement: Project-level Search
系统 SHALL 在项目详情页面提供项目内搜索功能。

#### Scenario: Search within project
- **WHEN** 用户在项目详情页面使用搜索框
- **THEN** 系统仅搜索当前项目的内容

#### Scenario: Filter by entity type
- **WHEN** 用户选择实体类型过滤器（章节、角色、地点等）
- **THEN** 系统仅显示该类型的搜索结果

#### Scenario: Search in chapter content
- **WHEN** 用户搜索章节内容
- **THEN** 系统显示包含关键词的章节，并显示匹配段落的上下文

### Requirement: Advanced Filter
系统 SHALL 提供高级筛选功能，支持多条件组合筛选。

#### Scenario: Filter by tags
- **WHEN** 用户选择一个或多个标签
- **THEN** 系统仅显示包含所有选中标签的项目或章节

#### Scenario: Filter by status
- **WHEN** 用户选择状态（草稿、进行中、已完成）
- **THEN** 系统仅显示对应状态的项目或章节

#### Scenario: Filter by date range
- **WHEN** 用户选择日期范围（创建日期或更新日期）
- **THEN** 系统仅显示在该日期范围内的项目或章节

#### Scenario: Filter by word count
- **WHEN** 用户设置字数范围（最小值、最大值）
- **THEN** 系统仅显示字数在该范围内的章节

#### Scenario: Combine filters
- **WHEN** 用户同时应用多个筛选条件
- **THEN** 系统显示满足所有条件的结果

### Requirement: Search History
系统 SHALL 记录用户的搜索历史，支持快速重复搜索。

#### Scenario: Save search history
- **WHEN** 用户执行搜索
- **THEN** 系统将搜索关键词保存到历史记录

#### Scenario: Display search history
- **WHEN** 用户打开搜索框
- **THEN** 系统显示最近10条搜索历史

#### Scenario: Repeat search
- **WHEN** 用户点击历史记录中的某个搜索
- **THEN** 系统重新执行该搜索

#### Scenario: Clear search history
- **WHEN** 用户点击"清除历史"按钮
- **THEN** 系统删除所有搜索历史记录

### Requirement: Search Suggestions
系统 SHALL 提供搜索建议，帮助用户快速找到内容。

#### Scenario: Auto-complete
- **WHEN** 用户在搜索框输入关键词
- **THEN** 系统显示自动完成建议（基于历史搜索和常见搜索）

#### Scenario: Related searches
- **WHEN** 用户执行搜索后
- **THEN** 系统在结果底部显示相关搜索建议

#### Scenario: Popular searches
- **WHEN** 用户打开搜索框且未输入内容
- **THEN** 系统显示热门搜索关键词

### Requirement: Search Performance
系统 SHALL 优化搜索性能，确保快速响应。

#### Scenario: Instant search
- **WHEN** 用户输入搜索关键词
- **THEN** 系统在300ms内返回搜索结果

#### Scenario: Debounce input
- **WHEN** 用户快速输入多个字符
- **THEN** 系统等待用户停止输入300ms后再执行搜索

#### Scenario: Pagination
- **WHEN** 搜索结果超过50条
- **THEN** 系统分页显示结果，每页20条

### Requirement: Search Highlighting
系统 SHALL 在搜索结果中高亮显示匹配的关键词。

#### Scenario: Highlight in results
- **WHEN** 系统显示搜索结果
- **THEN** 系统在标题和摘要中高亮显示匹配的关键词

#### Scenario: Highlight in content
- **WHEN** 用户导航到搜索结果的详情页面
- **THEN** 系统在内容中高亮显示匹配的关键词

#### Scenario: Scroll to match
- **WHEN** 用户打开搜索结果的详情页面
- **THEN** 系统自动滚动到第一个匹配位置
