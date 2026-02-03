# Spec: 章节管理模块

> FR-005, FR-006 | 章节 CRUD、版本管理、AI 生成

---

## 1. 功能需求

### 1.1 创建章节 (FR-005-01)

**输入**:
```json
{
  "project_id": "uuid",
  "title": "string (1-200 chars)",
  "content": "string (optional, max 100000 chars)",
  "chapter_number": "int (optional, auto-increment)"
}
```

**输出**:
```json
{
  "id": "uuid",
  "project_id": "uuid",
  "chapter_number": 1,
  "title": "string",
  "content": "string",
  "status": "draft",
  "word_count": 0,
  "created_at": "ISO 8601"
}
```

**验收标准**:
- [ ] chapter_number 自动递增 (项目内最大值 + 1)
- [ ] 自动计算 word_count
- [ ] 初始状态为 `draft`
- [ ] 更新项目 total_chapters

### 1.2 获取章节列表 (FR-005-02)

**输入**:
```
GET /api/v1/projects/{project_id}/chapters?page=1&page_size=50
```

**输出**:
```json
{
  "items": [
    {
      "id": "uuid",
      "chapter_number": 1,
      "title": "string",
      "status": "draft",
      "word_count": 2500,
      "updated_at": "ISO 8601"
    }
  ],
  "total": 100
}
```

**验收标准**:
- [ ] 默认按 chapter_number 升序
- [ ] 列表不返回 content (节省带宽)
- [ ] 支持状态过滤

### 1.3 获取章节详情 (FR-005-03)

**输入**:
```
GET /api/v1/chapters/{id}
```

**输出**: 完整章节对象 + 元数据

**验收标准**:
- [ ] 返回完整 content
- [ ] 包含追读力指标 (hook_type, hook_strength)
- [ ] 包含最新分析结果摘要

### 1.4 更新章节 (FR-005-04)

**输入**:
```json
{
  "title": "string (optional)",
  "content": "string (optional)",
  "status": "enum (optional)",
  "pov_character": "string (optional)",
  "location": "string (optional)"
}
```

**验收标准**:
- [ ] 内容变更时重新计算 word_count
- [ ] 内容变更时更新项目 total_words
- [ ] 自动创建版本记录

### 1.5 删除章节 (FR-005-05)

**输入**:
```
DELETE /api/v1/chapters/{id}
```

**验收标准**:
- [ ] 物理删除章节及关联数据
- [ ] 更新项目 total_chapters 和 total_words
- [ ] 不自动重排 chapter_number

### 1.6 重排章节 (FR-005-06)

**输入**:
```json
{
  "chapter_ids": ["uuid1", "uuid2", "uuid3"]
}
```

**验收标准**:
- [ ] 按数组顺序更新 chapter_number
- [ ] 原子操作
- [ ] 仅更新同一项目的章节

---

## 2. 版本管理

### 2.1 获取版本历史 (FR-006-01)

**输入**:
```
GET /api/v1/chapters/{id}/versions
```

**输出**:
```json
{
  "items": [
    {
      "id": "uuid",
      "version_number": 3,
      "source": "ai_rewrite",
      "word_count": 2800,
      "created_at": "ISO 8601"
    }
  ]
}
```

**验收标准**:
- [ ] 按 version_number 降序
- [ ] 列表不返回 content

### 2.2 获取版本详情 (FR-006-02)

**输入**:
```
GET /api/v1/chapters/{id}/versions/{version_number}
```

**输出**: 完整版本内容

### 2.3 恢复版本 (FR-006-03)

**输入**:
```
POST /api/v1/chapters/{id}/versions/{version_number}/restore
```

**验收标准**:
- [ ] 将指定版本内容复制到当前章节
- [ ] 创建新版本记录 (source=restore)
- [ ] 不删除历史版本

### 2.4 版本对比 (FR-006-04)

**输入**:
```
GET /api/v1/chapters/{id}/versions/diff?from=1&to=3
```

**输出**:
```json
{
  "from_version": 1,
  "to_version": 3,
  "diff": [
    { "type": "add", "content": "新增内容", "position": 100 },
    { "type": "remove", "content": "删除内容", "position": 50 }
  ],
  "stats": {
    "additions": 500,
    "deletions": 200
  }
}
```

---

## 3. AI 生成

### 3.1 流式生成章节 (FR-005-07)

**输入**:
```
POST /api/v1/chapters/{id}/generate-stream
Content-Type: application/json

{
  "prompt": "string (optional)",
  "use_outline": true,
  "use_context": true,
  "context_chapters": 3,
  "temperature": 0.7,
  "max_tokens": 4096
}
```

**输出**: SSE 流

```
event: content
data: {"text": "第一章开头..."}

event: content
data: {"text": "继续生成..."}

event: tool_call
data: {"tool": "retrieve_context", "args": {...}}

event: tool_result
data: {"result": "检索到的上下文..."}

event: done
data: {"word_count": 2500, "usage": {"input_tokens": 1000, "output_tokens": 2500}}
```

**验收标准**:
- [ ] 支持取消 (AbortController)
- [ ] 生成完成后自动保存
- [ ] 创建版本记录 (source=ai_generate)
- [ ] 记录用量

### 3.2 局部重写 (FR-005-08)

**输入**:
```
POST /api/v1/chapters/{id}/partial-regenerate-stream
Content-Type: application/json

{
  "selection_start": 100,
  "selection_end": 500,
  "instruction": "让这段更有张力"
}
```

**验收标准**:
- [ ] 仅重写选中部分
- [ ] 保持上下文连贯
- [ ] 创建版本记录 (source=ai_rewrite)

### 3.3 润色 (FR-005-09)

**输入**:
```
POST /api/v1/chapters/{id}/polish-stream
Content-Type: application/json

{
  "aspects": ["grammar", "style", "pacing"]
}
```

**验收标准**:
- [ ] 返回修改建议列表
- [ ] 支持逐条应用或全部应用
- [ ] 创建版本记录 (source=ai_polish)

---

## 4. PBT 属性

### 4.1 章节编号唯一性

**属性**: 同一项目内 chapter_number 唯一
**验证**: 并发创建 100 个章节，验证无重复编号

### 4.2 字数统计准确性

**属性**: word_count = 实际中文字数 + 英文单词数
**验证**: 生成各种内容 (纯中文、纯英文、混合、emoji)，验证统计准确

### 4.3 版本不可变性

**属性**: 已创建的版本内容不可修改
**验证**: 创建版本后尝试修改，验证失败

---

## 5. 边界条件

| 场景 | 预期行为 |
|------|----------|
| content 超过 100000 字符 | 400 验证失败 |
| 并发更新同一章节 | 乐观锁冲突返回 409 |
| 生成时项目被删除 | 生成中断，返回 404 |
| 恢复不存在的版本 | 404 |
| 选区超出内容范围 | 400 |

---

*Spec ID: CHAPTER-001*
