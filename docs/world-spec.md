# Spec: 世界观管理模块

> FR-007 ~ FR-010 | 角色、地点、组织、世界设定

---

## 1. 角色管理 (FR-007)

### 1.1 创建角色 (FR-007-01)

**输入**:
```json
{
  "project_id": "uuid",
  "name": "string (1-100 chars)",
  "aliases": ["string"],
  "gender": "string (optional)",
  "age": "string (optional)",
  "race": "string (optional)",
  "identity": "string (optional)",
  "appearance": "string (optional, max 5000)",
  "personality": "string (optional, max 5000)",
  "background": "string (optional, max 10000)",
  "motivation": "string (optional, max 5000)",
  "abilities": ["string"],
  "character_type": "enum: protagonist|antagonist|supporting|minor",
  "is_protagonist": "boolean"
}
```

**验收标准**:
- [ ] 名称在项目内唯一
- [ ] 自动生成向量嵌入 (异步)
- [ ] is_protagonist=true 时 character_type 自动设为 protagonist

### 1.2 角色关系 (FR-007-02)

**创建关系**:
```json
{
  "from_character_id": "uuid",
  "to_character_id": "uuid",
  "relationship_type": "enum: family|friend|enemy|lover|master|subordinate",
  "relationship_name": "string (e.g., 父子, 师徒)",
  "description": "string (optional)",
  "intensity": "int (1-10)"
}
```

**验收标准**:
- [ ] 不能创建自关联
- [ ] 双向关系自动创建反向记录
- [ ] 关系强度默认 5

### 1.3 角色经历 (FR-007-03)

**创建经历**:
```json
{
  "character_id": "uuid",
  "career_type": "enum: occupation|title|achievement|event|transformation",
  "title": "string",
  "description": "string (optional)",
  "start_chapter": "int (optional)",
  "end_chapter": "int (optional)",
  "related_organization_id": "uuid (optional)",
  "abilities_gained": ["string"],
  "abilities_lost": ["string"]
}
```

**验收标准**:
- [ ] start_chapter <= end_chapter (如果都提供)
- [ ] 按 sort_order 排序

---

## 2. 地点管理 (FR-008)

### 2.1 创建地点 (FR-008-01)

**输入**:
```json
{
  "project_id": "uuid",
  "name": "string (1-200 chars)",
  "location_type": "enum: city|country|building|natural|other",
  "description": "string (optional, max 10000)",
  "atmosphere": "string (optional)",
  "parent_id": "uuid (optional)",
  "climate": "string (optional)",
  "notable_characters": ["uuid"]
}
```

**验收标准**:
- [ ] 支持层级结构 (最大 5 层)
- [ ] parent_id 必须属于同一项目
- [ ] 自动生成向量嵌入

### 2.2 地点树结构 (FR-008-02)

**输入**:
```
GET /api/v1/projects/{id}/locations?tree=true
```

**验收标准**:
- [ ] 递归构建树结构
- [ ] 包含子地点数量统计

---

## 3. 组织管理 (FR-009)

### 3.1 创建组织 (FR-009-01)

**输入**:
```json
{
  "project_id": "uuid",
  "name": "string (1-200 chars)",
  "org_type": "enum: sect|kingdom|company|family|gang",
  "description": "string (optional)",
  "ideology": "string (optional)",
  "scale": "enum: small|medium|large|massive",
  "influence_level": "int (1-10)",
  "parent_id": "uuid (optional)",
  "leader_id": "uuid (optional)",
  "headquarters_id": "uuid (optional)"
}
```

**验收标准**:
- [ ] leader_id 必须是同项目角色
- [ ] headquarters_id 必须是同项目地点
- [ ] 支持层级结构

### 3.2 组织成员 (FR-009-02)

**添加成员**:
```json
{
  "organization_id": "uuid",
  "character_id": "uuid",
  "position": "string",
  "rank": "int",
  "joined_chapter": "int (optional)",
  "permissions": ["string"]
}
```

**验收标准**:
- [ ] 同一角色可加入多个组织
- [ ] 同一组织内角色唯一
- [ ] 支持成员状态变更 (active -> inactive -> expelled)

---

## 4. 世界设定 (FR-010)

### 4.1 创建设定 (FR-010-01)

**输入**:
```json
{
  "project_id": "uuid",
  "name": "string (1-200 chars)",
  "setting_type": "enum: magic_system|technology|social|geography|history|culture|economy|religion",
  "description": "string",
  "rules": [
    { "name": "string", "description": "string" }
  ],
  "limitations": [
    { "name": "string", "description": "string" }
  ],
  "parent_id": "uuid (optional)",
  "related_characters": ["uuid"],
  "related_locations": ["uuid"],
  "related_organizations": ["uuid"]
}
```

**验收标准**:
- [ ] 自动生成向量嵌入
- [ ] 支持层级结构 (如: 魔法体系 -> 火系魔法 -> 火球术)
- [ ] 关联实体必须属于同一项目

---

## 5. PBT 属性

### 5.1 角色名称唯一性

**属性**: 同一项目内角色名称唯一
**验证**: 并发创建同名角色，仅一个成功

### 5.2 关系对称性

**属性**: 创建 A->B 关系时自动创建 B->A 反向关系
**验证**: 创建关系后查询双方，验证关系存在

### 5.3 层级深度限制

**属性**: 地点/组织/设定层级不超过 5 层
**验证**: 尝试创建第 6 层，验证失败

---

## 6. 边界条件

| 场景 | 预期行为 |
|------|----------|
| 角色名称 101 字符 | 400 验证失败 |
| 创建自关联关系 | 400 |
| 关联不存在的角色 | 400 |
| 关联其他项目的角色 | 400 |
| 删除有成员的组织 | 级联删除成员关系 |
| 删除有子地点的地点 | 子地点提升到上级 |

---

*Spec ID: WORLD-001*
