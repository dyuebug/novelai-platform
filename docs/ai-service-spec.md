# Spec: AI 服务模块

> FR-011 ~ FR-015 | LLM 调用、RAG 检索、质量评估

---

## 1. LLM 调用服务 (FR-011)

### 1.1 流式生成 (FR-011-01)

**gRPC 接口**:
```protobuf
rpc GenerateStream(GenerateRequest) returns (stream GenerateResponse);

message GenerateRequest {
    string prompt = 1;
    string model = 2;           // 默认: gpt-4o
    string provider = 3;        // 默认: openai
    float temperature = 4;      // 默认: 0.7, 范围: 0-2
    int32 max_tokens = 5;       // 默认: 4096, 最大: 128000
    string system_prompt = 6;
    string user_id = 7;
    string project_id = 8;
    map<string, string> metadata = 9;
}

message GenerateResponse {
    string content = 1;
    bool done = 2;
    string error = 3;
    UsageInfo usage = 4;  // 仅在 done=true 时填充
}
```

**验收标准**:
- [ ] 支持 OpenAI, Anthropic, Gemini 三个提供商
- [ ] 流式输出延迟 < 100ms
- [ ] 自动记录用量到 usage_records
- [ ] 支持取消请求

### 1.2 提供商配置 (FR-011-02)

| 提供商 | 默认模型 | 支持模型 |
|--------|----------|----------|
| openai | gpt-4o | gpt-4o, gpt-4o-mini, gpt-4-turbo |
| anthropic | claude-sonnet-4-20250514 | claude-sonnet-4-20250514, claude-opus-4-20250514 |
| gemini | gemini-2.5-pro | gemini-2.5-pro, gemini-2.5-flash |

**验收标准**:
- [ ] 不支持的模型返回 400
- [ ] API Key 从用户设置或环境变量获取
- [ ] 支持自定义 base_url (OpenAI 兼容)

---

## 2. RAG 检索服务 (FR-012)

### 2.1 分层检索 (FR-012-01)

**gRPC 接口**:
```protobuf
rpc RetrieveContext(RetrieveRequest) returns (RetrieveResponse);

message RetrieveRequest {
    string project_id = 1;
    string query = 2;
    int32 top_k = 3;            // 默认: 10
    repeated string layers = 4;  // 默认: ["settings", "summaries", "chunks"]
    int32 chapter_range_start = 5;  // 可选: 限制章节范围
    int32 chapter_range_end = 6;
}

message RetrieveResponse {
    map<string, LayerResult> results = 1;
    int32 total_tokens = 2;
}

message LayerResult {
    repeated RetrievedItem items = 1;
}

message RetrievedItem {
    string content = 1;
    float score = 2;
    string source_type = 3;
    string source_id = 4;
    int32 chapter_number = 5;
}
```

**检索层级**:
| 层级 | 内容类型 | 优先级 |
|------|----------|--------|
| settings | 世界设定、角色信息 | 高 |
| summaries | 章节摘要 | 中 |
| chunks | 内容片段 | 低 |
| characters | 角色详情 | 按需 |

**验收标准**:
- [ ] 相似度阈值 > 0.7 才返回
- [ ] 按相似度降序排列
- [ ] 支持章节范围过滤
- [ ] 返回 token 统计

### 2.2 向量索引更新 (FR-012-02)

**触发条件**:
- 章节内容更新
- 角色/地点/组织/设定创建或更新
- 手动触发重建

**验收标准**:
- [ ] 异步处理，不阻塞主流程
- [ ] 增量更新，仅处理变更内容
- [ ] 支持批量重建

---

## 3. 质量评估服务 (FR-013)

### 3.1 追读力分析 (FR-013-01)

**gRPC 接口**:
```protobuf
rpc AnalyzeReadingPower(AnalyzeRequest) returns (AnalyzeResponse);

message AnalyzeRequest {
    string chapter_id = 1;
    string content = 2;
}

message AnalyzeResponse {
    float overall_score = 1;  // 0-10
    HookAnalysis hook = 2;
    CoolpointAnalysis coolpoint = 3;
    MicropayoffAnalysis micropayoff = 4;
}

message HookAnalysis {
    string hook_type = 1;      // crisis, mystery, emotion, choice, desire
    string hook_strength = 2;  // strong, medium, weak
    string hook_content = 3;   // 钩子内容摘录
    float score = 4;
}
```

**评估维度**:
| 维度 | 说明 | 权重 |
|------|------|------|
| 钩子 (Hook) | 章节开头吸引力 | 30% |
| 爽点 (Coolpoint) | 情绪高潮分布 | 40% |
| 微兑现 (Micropayoff) | 小目标达成感 | 30% |

**验收标准**:
- [ ] 分析结果存入 chapter_analyses
- [ ] 支持批量分析
- [ ] 返回具体改进建议

### 3.2 一致性检查 (FR-013-02)

**检查项**:
| 检查项 | 说明 |
|--------|------|
| 角色一致性 | 性格、能力、外貌是否前后一致 |
| 时间线一致性 | 事件顺序是否合理 |
| 世界观一致性 | 是否违反设定规则 |
| 称呼一致性 | 角色称呼是否统一 |

**输出**:
```json
{
  "consistency_score": 8.5,
  "issues": [
    {
      "type": "character",
      "severity": "warning",
      "description": "张三在第5章是黑发，第8章变成金发",
      "chapter_refs": [5, 8],
      "suggestion": "统一角色外貌描述"
    }
  ]
}
```

**验收标准**:
- [ ] 自动关联相关章节
- [ ] 提供修复建议
- [ ] 支持忽略特定问题

### 3.3 多 Agent 审查 (FR-013-03)

**审查维度**:
| Agent | 职责 |
|-------|------|
| 情节 Agent | 剧情逻辑、节奏 |
| 角色 Agent | 人物塑造、对话 |
| 文笔 Agent | 语言风格、表达 |
| 设定 Agent | 世界观一致性 |
| 读者 Agent | 阅读体验、吸引力 |
| 编辑 Agent | 综合评估、修改建议 |

**验收标准**:
- [ ] 6 个 Agent 并行执行
- [ ] 汇总各 Agent 意见
- [ ] 生成综合评分和建议

---

## 4. 宪法约束服务 (FR-014)

### 4.1 约束检查 (FR-014-01)

**输入**:
```json
{
  "project_id": "uuid",
  "chapter_id": "uuid",
  "content": "string"
}
```

**输出**:
```json
{
  "violations": [
    {
      "constraint_type": "hard",
      "constraint_id": "pov_consistency",
      "description": "视角从第一人称切换到第三人称",
      "location": { "start": 100, "end": 200 },
      "severity": "error"
    }
  ],
  "warnings": [...],
  "passed": false
}
```

**验收标准**:
- [ ] 硬约束违反阻止保存
- [ ] 软约束违反仅警告
- [ ] 支持约束豁免 (override_contracts)

### 4.2 约束豁免 (FR-014-02)

**创建豁免**:
```json
{
  "project_id": "uuid",
  "chapter_id": "uuid",
  "constraint_id": "pov_consistency",
  "rationale_type": "artistic",
  "rationale_text": "此处切换视角是为了制造悬念",
  "payback_plan": "在第10章回归主视角",
  "due_chapter": 10
}
```

**验收标准**:
- [ ] 豁免后该约束不再报错
- [ ] 到期未偿还标记为 overdue
- [ ] 支持豁免审计

---

## 5. PBT 属性

### 5.1 RAG 检索相关性

**属性**: 检索结果与查询语义相关
**验证**: 生成 100 个查询，验证 top-1 结果相关性 > 0.8

### 5.2 评分一致性

**属性**: 相同内容多次评估结果一致 (±0.5)
**验证**: 同一章节评估 10 次，验证标准差 < 0.5

### 5.3 约束检查完整性

**属性**: 所有定义的约束都被检查
**验证**: 构造违反每个约束的内容，验证都被检测到

---

## 6. 边界条件

| 场景 | 预期行为 |
|------|----------|
| 空内容分析 | 返回默认评分 |
| 超长内容 (>100k) | 分段处理 |
| 向量库为空 | 返回空结果，不报错 |
| AI 服务超时 | 重试 3 次后返回 503 |
| 无效模型名 | 400 |
| 配额不足 | 402 |

---

*Spec ID: AI-001*
