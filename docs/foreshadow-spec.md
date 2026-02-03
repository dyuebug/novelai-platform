# Spec: 伏笔管理模块

> FR-016 | 伏笔埋设、追踪、回收

---

## 1. 功能需求

### 1.1 创建伏笔 (FR-016-01)

**输入**:
```json
{
  "project_id": "uuid",
  "title": "string (1-200 chars)",
  "description": "string",
  "foreshadow_type": "enum: plot|character|world|item|mystery",
  "importance": "enum: critical|high|medium|low",
  "planted_chapter_id": "uuid (optional)",
  "planted_content": "string (optional)",
  "expected_resolution_chapter": "int (optional)",
  "related_character_ids": ["uuid"],
  "notes": "string (optional)"
}
```

**输出**:
```json
{
  "id": "uuid",
  "project_id": "uuid",
  "title": "string",
  "status": "planned",
  "created_at": "ISO 8601"
}
```

**验收标准**:
- [ ] 初始状态为 `planned`
- [ ] 关联角色必须属于同一项目
- [ ] 自动计算预计回收章节 (如未提供)

### 1.2 埋设伏笔 (FR-016-02)

**输入**:
```json
{
  "planted_chapter_id": "uuid",
  "planted_content": "string"
}
```

**验收标准**:
- [ ] 状态从 `planned` 变为 `planted`
- [ ] 记录埋设章节号
- [ ] 支持多次埋设 (暗示)

### 1.3 回收伏笔 (FR-016-03)

**输入**:
```json
{
  "resolved_chapter_id": "uuid",
  "resolved_content": "string"
}
```

**验收标准**:
- [ ] 状态从 `planted` 或 `hinted` 变为 `resolved`
- [ ] 记录回收章节号
- [ ] 计算埋设到回收的章节跨度

### 1.4 伏笔列表 (FR-016-04)

**输入**:
```
GET /api/v1/projects/{id}/foreshadows?status=planted&importance=critical
```

**输出**:
```json
{
  "items": [
    {
      "id": "uuid",
      "title": "神秘老人的身份",
      "status": "planted",
      "importance": "critical",
      "planted_chapter_number": 5,
      "expected_resolution_chapter": 50,
      "overdue": false
    }
  ],
  "stats": {
    "total": 20,
    "planned": 5,
    "planted": 10,
    "resolved": 4,
    "abandoned": 1,
    "overdue": 2
  }
}
```

**验收标准**:
- [ ] 支持状态和重要性过滤
- [ ] 返回统计信息
- [ ] 标记逾期未回收的伏笔

### 1.5 伏笔提醒 (FR-016-05)

**触发条件**:
- 当前章节号 >= expected_resolution_chapter - 5
- 伏笔状态为 `planted` 或 `hinted`

**输出**:
```json
{
  "reminders": [
    {
      "foreshadow_id": "uuid",
      "title": "神秘老人的身份",
      "planted_chapter": 5,
      "expected_resolution": 50,
      "current_chapter": 45,
      "urgency": "high"
    }
  ]
}
```

**验收标准**:
- [ ] 在章节编辑页面显示提醒
- [ ] 支持关闭提醒
- [ ] 逾期伏笔高亮显示

---

## 2. 状态流转

```
planned -> planted -> hinted -> resolved
    |         |          |
    v         v          v
abandoned  abandoned  abandoned
```

| 状态 | 说明 |
|------|------|
| planned | 计划中，尚未埋设 |
| planted | 已埋设，等待回收 |
| hinted | 已暗示，强化伏笔 |
| resolved | 已回收 |
| abandoned | 已放弃 |

---

## 3. PBT 属性

### 3.1 状态流转合法性

**属性**: 状态只能按定义的路径流转
**验证**: 尝试非法状态转换，验证失败

### 3.2 章节关联一致性

**属性**: planted_chapter_number <= resolved_chapter_number
**验证**: 创建各种章节组合，验证约束

### 3.3 逾期检测准确性

**属性**: 当前章节 > expected_resolution_chapter 且未回收时标记为逾期
**验证**: 模拟章节进度，验证逾期标记

---

## 4. 边界条件

| 场景 | 预期行为 |
|------|----------|
| 回收未埋设的伏笔 | 400 |
| 埋设已回收的伏笔 | 400 |
| 删除已埋设的伏笔 | 状态变为 abandoned |
| 关联不存在的章节 | 400 |
| expected_resolution < planted_chapter | 400 |

---

*Spec ID: FORESHADOW-001*
