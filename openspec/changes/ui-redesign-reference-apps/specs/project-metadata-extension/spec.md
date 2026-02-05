# Specification: Project Metadata Extension

## ADDED Requirements

### Requirement: Cover Image Support
系统 SHALL 扩展项目模型，支持封面图片URL字段。

#### Scenario: Store cover URL
- **WHEN** 用户上传项目封面图片
- **THEN** 系统将图片URL存储到项目的metadata.cover_image_url字段

#### Scenario: Return cover URL
- **WHEN** 客户端请求项目详情
- **THEN** 系统在响应中包含封面图片URL

#### Scenario: Default cover
- **WHEN** 项目没有封面图片
- **THEN** 系统返回null或默认封面URL

### Requirement: Tags Support
系统 SHALL 扩展项目模型，支持标签数组字段。

#### Scenario: Store tags
- **WHEN** 用户为项目添加标签（玄幻、都市、科幻等）
- **THEN** 系统将标签数组存储到项目的metadata.tags字段

#### Scenario: Filter by tags
- **WHEN** 用户按标签筛选项目
- **THEN** 系统返回包含指定标签的所有项目

#### Scenario: Tag validation
- **WHEN** 用户添加标签
- **THEN** 系统验证标签长度（1-20字符）和数量（最多10个）

### Requirement: Statistics Support
系统 SHALL 扩展项目模型，支持统计信息字段。

#### Scenario: Calculate total chapters
- **WHEN** 系统返回项目详情
- **THEN** 系统计算并返回项目的总章节数

#### Scenario: Calculate total words
- **WHEN** 系统返回项目详情
- **THEN** 系统计算并返回项目的总字数

#### Scenario: Calculate average chapter length
- **WHEN** 系统返回项目详情
- **THEN** 系统计算并返回平均章节字数

#### Scenario: Calculate completion percent
- **WHEN** 系统返回项目详情
- **THEN** 系统根据已完成章节数计算完成百分比

#### Scenario: Track last update time
- **WHEN** 用户修改项目或章节
- **THEN** 系统更新项目的last_update_time字段

### Requirement: Custom Fields Support
系统 SHALL 支持自定义元数据字段，允许用户存储额外信息。

#### Scenario: Store custom fields
- **WHEN** 用户添加自定义字段（目标字数、预计完成日期等）
- **THEN** 系统将自定义字段存储到metadata.custom_fields对象

#### Scenario: Retrieve custom fields
- **WHEN** 客户端请求项目详情
- **THEN** 系统在响应中包含所有自定义字段

#### Scenario: Validate custom fields
- **WHEN** 用户添加自定义字段
- **THEN** 系统验证字段名（不能与保留字段冲突）和值类型

### Requirement: Update Metadata API
系统 SHALL 提供API端点更新项目元数据。

#### Scenario: Update cover image
- **WHEN** 客户端发送 PUT /api/v1/projects/:id/metadata 带cover_image_url
- **THEN** 系统更新项目封面URL并返回200

#### Scenario: Update tags
- **WHEN** 客户端发送 PUT /api/v1/projects/:id/metadata 带tags数组
- **THEN** 系统更新项目标签并返回200

#### Scenario: Update custom fields
- **WHEN** 客户端发送 PUT /api/v1/projects/:id/metadata 带custom_fields对象
- **THEN** 系统合并更新自定义字段并返回200

#### Scenario: Partial update
- **WHEN** 客户端仅发送部分元数据字段
- **THEN** 系统仅更新提供的字段，保留其他字段不变

### Requirement: Statistics API
系统 SHALL 提供API端点获取项目统计信息。

#### Scenario: Get project statistics
- **WHEN** 客户端发送 GET /api/v1/projects/:id/statistics
- **THEN** 系统返回项目统计信息（总章节数、总字数、平均章节长度、完成百分比、最后更新时间）

#### Scenario: Real-time calculation
- **WHEN** 系统计算统计信息
- **THEN** 系统实时查询数据库，确保数据准确

#### Scenario: Cache statistics
- **WHEN** 统计信息计算成本高
- **THEN** 系统缓存统计结果5分钟，减少数据库负载

### Requirement: Backward Compatibility
系统 SHALL 确保元数据扩展向后兼容现有项目。

#### Scenario: Existing projects
- **WHEN** 系统查询旧项目（无元数据字段）
- **THEN** 系统返回默认元数据值（空标签、null封面、计算统计）

#### Scenario: Migration
- **WHEN** 系统启动时检测到旧数据
- **THEN** 系统自动为现有项目添加metadata字段（默认值）

#### Scenario: API compatibility
- **WHEN** 客户端请求项目列表或详情
- **THEN** 系统始终返回metadata字段，即使为空对象
