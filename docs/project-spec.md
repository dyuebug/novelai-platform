# Spec: 项目管理模块

> FR-003, FR-004 | 项目 CRUD、大纲管理

---

## 1. 功能需求

### 1.1 创建项目 (FR-003-01)

**输入**:
```json
{
  "title": "string (1-200 chars)",
  "description": "string (optional, max 5000 chars)",
  "genre": "enum: 玄幻|都市|科幻|历史|言情|其他",
  "metadata": "object (optional)"
}
```

**输出**:
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "string",
  "description": "string",
  "genre": "string",
  "status": "draft",
  "total_chapters": 0,
  "total_words": 0,
  "created_at": "ISO 8601",
  "updated_at": "ISO 8601"
}
```

**验收标准**:
- [ ] 标题不能为空
- [ ] 自动关联当前用户
- [ ] 初始状态为 `draft`
- [ ] 返回 201 Created

### 1.2 获取项目列表 (FR-003-02)

**输入**:
```
GET /api/v1/projects?page=1&page_size=20&status=writing&sort=-updated_at
```

**参数**:
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码 (>=1) |
| page_size | int | 20 | 每页数量 (1-100) |
| status | string | - | 状态过滤 |
| genre | string | - | 类型过滤 |
| sort | string | -updated_at | 排序字段 |

**输出**:
```json
{
  "items": [...],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "total_pages": 5
}
```

**验收标准**:
- [ ] 仅返回当前用户的项目
- [ ] 不返回已软删除的项目 (is_deleted=true)
- [ ] 支持多字段排序
- [ ] page_size 超过 100 时截断为 100

### 1.3 获取项目详情 (FR-003-03)

**输入**:
```
GET /api/v1/projects/{id}
```

**输出**: 完整项目对象 + 统计信息

**验收标准**:
- [ ] 项目不存在返回 404
- [ ] 非所有者访问返回 403
- [ ] 已删除项目返回 404

### 1.4 更新项目 (FR-003-04)

**输入**:
```json
{
  "title": "string (optional)",
  "description": "string (optional)",
  "genre": "string (optional)",
  "status": "enum (optional)",
  "cover_url": "string (optional)"
}
```

**验收标准**:
- [ ] 仅更新提供的字段
- [ ] 自动更新 `updated_at`
- [ ] 状态转换验证 (draft -> writing -> completed -> archived)

### 1.5 删除项目 (FR-003-05)

**输入**:
```
DELETE /api/v1/projects/{id}
```

**验收标准**:
- [ ] 软删除: 设置 `is_deleted=true`, `deleted_at=NOW()`
- [ ] 返回 204 No Content
- [ ] 关联数据不物理删除

### 1.6 恢复项目 (FR-003-06)

**输入**:
```
POST /api/v1/projects/{id}/restore
```

**验收标准**:
- [ ] 仅对已删除项目有效
- [ ] 清除 `is_deleted` 和 `deleted_at`
- [ ] 返回恢复后的项目

---

## 2. 大纲管理

### 2.1 创建大纲 (FR-004-01)

**输入**:
```json
{
  "project_id": "uuid",
  "title": "string (optional)",
  "content": "string",
  "outline_type": "enum: main|volume|chapter",
  "parent_id": "uuid (optional)",
  "start_chapter": "int (optional)",
  "end_chapter": "int (optional)"
}
```

**验收标准**:
- [ ] main 类型每个项目只能有一个
- [ ] parent_id 必须属于同一项目
- [ ] 自动计算 sort_order

### 2.2 大纲树结构 (FR-004-02)

**输入**:
```
GET /api/v1/projects/{id}/outlines?tree=true
```

**输出**:
```json
{
  "id": "uuid",
  "title": "主线大纲",
  "content": "...",
  "children": [
    {
      "id": "uuid",
      "title": "第一卷",
      "children": [...]
    }
  ]
}
```

**验收标准**:
- [ ] 递归构建树结构
- [ ] 按 sort_order 排序
- [ ] 最大深度 5 层

### 2.3 大纲排序 (FR-004-03)

**输入**:
```json
{
  "outline_ids": ["uuid1", "uuid2", "uuid3"]
}
```

**验收标准**:
- [ ] 按数组顺序更新 sort_order
- [ ] 仅更新同一父级下的大纲
- [ ] 原子操作，失败回滚

---

## 3. PBT 属性

### 3.1 项目隔离

**属性**: 用户 A 无法访问用户 B 的项目
**验证**: 随机生成 user_id + project_id 组合，验证非所有者访问返回 403

### 3.2 软删除可恢复

**属性**: 删除后恢复的项目与删除前一致
**验证**: 创建项目 -> 删除 -> 恢复 -> 比较字段

### 3.3 大纲树完整性

**属性**: 删除父大纲时子大纲级联删除或提升
**验证**: 创建多层大纲 -> 删除中间层 -> 验证树结构

---

## 4. 边界条件

| 场景 | 预期行为 |
|------|----------|
| 标题 201 字符 | 400 验证失败 |
| page_size=0 | 400 验证失败 |
| page_size=101 | 截断为 100 |
| 删除已删除项目 | 404 |
| 恢复未删除项目 | 400 |
| 大纲深度超过 5 层 | 400 |

---

*Spec ID: PROJECT-001*
