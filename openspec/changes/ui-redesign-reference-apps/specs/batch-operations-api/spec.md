# Specification: Batch Operations API

## ADDED Requirements

### Requirement: Batch Update Chapters
系统 SHALL 提供API端点批量更新章节，支持修改状态、排序、标签等字段。

#### Scenario: Batch update status
- **WHEN** 客户端发送 POST /api/v1/projects/:id/chapters/batch-update 带多个章节ID和新状态
- **THEN** 系统批量更新章节状态并返回成功数量

#### Scenario: Batch reorder
- **WHEN** 客户端发送章节ID数组和新排序
- **THEN** 系统批量更新章节的sort_order字段

#### Scenario: Batch add tags
- **WHEN** 客户端发送章节ID数组和标签数组
- **THEN** 系统为所有指定章节添加标签

#### Scenario: Transaction support
- **WHEN** 批量操作中某个章节更新失败
- **THEN** 系统回滚所有更新，返回错误信息

### Requirement: Batch Delete Chapters
系统 SHALL 提供API端点批量删除章节。

#### Scenario: Batch soft delete
- **WHEN** 客户端发送 POST /api/v1/projects/:id/chapters/batch-delete 带章节ID数组
- **THEN** 系统批量软删除章节（设置deleted_at字段）

#### Scenario: Validate ownership
- **WHEN** 用户尝试批量删除章节
- **THEN** 系统验证所有章节属于该用户的项目，否则返回403错误

#### Scenario: Limit batch size
- **WHEN** 用户尝试批量删除超过100个章节
- **THEN** 系统返回400错误，要求分批操作

### Requirement: Batch Status Update
系统 SHALL 提供API端点批量更新章节状态。

#### Scenario: Publish multiple chapters
- **WHEN** 客户端发送 POST /api/v1/projects/:id/chapters/batch-status 带章节ID数组和status="published"
- **THEN** 系统批量将章节状态设置为已发布

#### Scenario: Archive multiple chapters
- **WHEN** 客户端发送章节ID数组和status="archived"
- **THEN** 系统批量将章节状态设置为已归档

#### Scenario: Validate status transition
- **WHEN** 用户尝试批量更新章节状态
- **THEN** 系统验证状态转换是否合法（如草稿→已发布，但不能已归档→草稿）

### Requirement: Batch Operation Response
系统 SHALL 返回详细的批量操作结果。

#### Scenario: Success response
- **WHEN** 批量操作全部成功
- **THEN** 系统返回200，包含成功数量和更新的章节ID列表

#### Scenario: Partial success response
- **WHEN** 批量操作部分成功
- **THEN** 系统返回207（Multi-Status），包含成功和失败的章节ID及错误原因

#### Scenario: Complete failure response
- **WHEN** 批量操作全部失败
- **THEN** 系统返回400或500，包含错误信息

### Requirement: Batch Operation Validation
系统 SHALL 验证批量操作的输入数据。

#### Scenario: Validate chapter IDs
- **WHEN** 客户端发送批量操作请求
- **THEN** 系统验证所有章节ID存在且属于指定项目

#### Scenario: Validate operation data
- **WHEN** 客户端发送批量更新数据
- **THEN** 系统验证数据格式和值范围（如状态值、排序值）

#### Scenario: Validate batch size
- **WHEN** 客户端发送批量操作请求
- **THEN** 系统验证操作数量不超过100条

### Requirement: Batch Operation Performance
系统 SHALL 优化批量操作的性能。

#### Scenario: Bulk database update
- **WHEN** 系统执行批量更新
- **THEN** 系统使用单个SQL语句批量更新，而非逐条更新

#### Scenario: Async processing
- **WHEN** 批量操作数量超过50条
- **THEN** 系统异步处理操作，立即返回任务ID，客户端可查询进度

#### Scenario: Progress tracking
- **WHEN** 客户端查询批量操作进度
- **THEN** 系统返回已处理数量和总数量

### Requirement: Batch Reorder API
系统 SHALL 提供专门的批量重排序API，支持拖拽排序。

#### Scenario: Reorder chapters
- **WHEN** 客户端发送 POST /api/v1/projects/:id/chapters/reorder 带章节ID数组（新顺序）
- **THEN** 系统批量更新章节的sort_order字段，保持顺序一致

#### Scenario: Validate order
- **WHEN** 客户端发送重排序请求
- **THEN** 系统验证章节ID数组包含项目的所有章节，无遗漏

#### Scenario: Atomic reorder
- **WHEN** 系统执行重排序
- **THEN** 系统在事务中完成所有更新，确保原子性

### Requirement: Batch Operation Logging
系统 SHALL 记录批量操作的审计日志。

#### Scenario: Log batch operation
- **WHEN** 用户执行批量操作
- **THEN** 系统记录操作类型、影响的章节ID、操作时间、用户ID

#### Scenario: Log operation result
- **WHEN** 批量操作完成
- **THEN** 系统记录成功数量、失败数量、错误信息

#### Scenario: Query operation history
- **WHEN** 用户查询操作历史
- **THEN** 系统返回该用户的批量操作记录
