# Spec: 认证与授权模块

> FR-001, FR-002 | 用户认证、OAuth2、RBAC

---

## 1. 功能需求

### 1.1 用户注册 (FR-001-01)

**输入**:
```json
{
  "username": "string (3-50 chars, alphanumeric + underscore)",
  "email": "string (RFC 5322)",
  "password": "string (8-128 chars, 1 upper, 1 lower, 1 digit)"
}
```

**输出**:
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "role": "user",
  "created_at": "ISO 8601"
}
```

**验收标准**:
- [ ] 用户名唯一性校验，冲突返回 409
- [ ] 邮箱唯一性校验，冲突返回 409
- [ ] 密码使用 bcrypt (cost=12) 哈希存储
- [ ] 注册成功返回 201，不返回密码字段

### 1.2 用户登录 (FR-001-02)

**输入**:
```json
{
  "email": "string",
  "password": "string"
}
```

**输出**:
```json
{
  "access_token": "JWT (RS256, TTL=15min)",
  "refresh_token": "opaque (TTL=7d)",
  "token_type": "Bearer",
  "expires_in": 900
}
```

**验收标准**:
- [ ] 邮箱不存在返回 401 (不泄露用户存在性)
- [ ] 密码错误返回 401
- [ ] 连续 5 次失败后锁定账户 30 分钟
- [ ] 登录成功更新 `last_login_at`

### 1.3 OAuth2 登录 (FR-001-03)

**支持提供商**: LinuxDo, GitHub, Google

**流程**:
1. `GET /api/v1/oauth/{provider}/authorize`
   - 生成 state (32 bytes random)
   - 重定向到提供商授权页

2. `GET /api/v1/oauth/{provider}/callback?code=xxx&state=xxx`
   - 验证 state
   - 交换 code 获取 access_token
   - 获取用户信息
   - 创建或关联本地用户
   - 返回 JWT

**验收标准**:
- [ ] 必须使用 PKCE (S256)
- [ ] state 验证失败返回 400
- [ ] 回调 URL 必须在白名单内
- [ ] 新用户自动创建账户
- [ ] 已有用户关联 OAuth 账户

### 1.4 密码重置 (FR-001-04)

**请求重置**:
```
POST /api/v1/auth/password-reset/request
Body: { "email": "user@example.com" }
Response: 202 Accepted (无论邮箱是否存在)
```

**执行重置**:
```
POST /api/v1/auth/password-reset/verify
Body: { "token": "xxx", "new_password": "xxx" }
Response: 200 OK | 400 Invalid Token | 410 Token Expired
```

**验收标准**:
- [ ] Token 有效期 1 小时
- [ ] Token 使用后立即失效
- [ ] 同一邮箱 5 分钟内只能请求一次
- [ ] 新密码不能与最近 3 次相同

### 1.5 Token 刷新 (FR-001-05)

**输入**:
```json
{
  "refresh_token": "string"
}
```

**输出**: 同登录响应

**验收标准**:
- [ ] Refresh Token 无效返回 401
- [ ] Refresh Token 过期返回 401
- [ ] 刷新后旧 Refresh Token 失效 (轮换)

---

## 2. RBAC 权限模型

### 2.1 角色定义

| 角色 | 继承 | 说明 |
|------|------|------|
| `user` | - | 基础用户 |
| `vip` | user | 高级用户 |
| `admin` | vip | 管理员 |

### 2.2 权限矩阵

| 权限 | user | vip | admin |
|------|------|-----|-------|
| project:read (own) | ✓ | ✓ | ✓ |
| project:write (own) | ✓ | ✓ | ✓ |
| project:delete (own) | ✓ | ✓ | ✓ |
| project:read (all) | - | - | ✓ |
| ai:generate | ✓ | ✓ | ✓ |
| ai:generate:advanced | - | ✓ | ✓ |
| ai:analyze | ✓ | ✓ | ✓ |
| admin:users | - | - | ✓ |
| admin:system | - | - | ✓ |

### 2.3 资源级授权

```go
// 中间件示例
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims := c.MustGet("claims").(*Claims)
        resourceID := c.Param("id")

        // 检查资源所有权
        if strings.HasSuffix(permission, "(own)") {
            if !isOwner(claims.UserID, resourceID) {
                c.AbortWithStatus(403)
                return
            }
        }

        // 检查角色权限
        if !hasPermission(claims.Roles, permission) {
            c.AbortWithStatus(403)
            return
        }

        c.Next()
    }
}
```

---

## 3. PBT 属性

### 3.1 密码哈希不可逆

**属性**: 给定任意密码，无法从哈希值还原原始密码
**验证**: 生成 1000 个随机密码，验证哈希后无法通过暴力匹配还原

### 3.2 Token 唯一性

**属性**: 任意两个 Token 不相同
**验证**: 生成 10000 个 Token，验证无重复

### 3.3 权限继承正确性

**属性**: 高级角色拥有低级角色的所有权限
**验证**: 遍历权限矩阵，验证继承关系

---

## 4. 边界条件

| 场景 | 预期行为 |
|------|----------|
| 用户名 2 字符 | 400 验证失败 |
| 用户名 51 字符 | 400 验证失败 |
| 密码无大写字母 | 400 验证失败 |
| 邮箱格式错误 | 400 验证失败 |
| 并发注册相同用户名 | 仅一个成功，另一个 409 |
| Token 过期 1 秒 | 401 |
| 账户锁定期间登录 | 423 Locked |

---

*Spec ID: AUTH-001*
