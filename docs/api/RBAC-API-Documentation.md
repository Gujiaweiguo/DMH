# RBAC权限管理系统 API文档

## 概述

本文档描述了DMH数字营销中台RBAC权限管理系统的完整API接口。系统支持四种用户角色：平台管理员、品牌管理员、参与者和匿名用户，提供完整的用户认证、权限控制、品牌管理和安全审计功能。

## 基础信息

- **API版本**: v1.0
- **基础URL**: `https://api.dmh.com/api/v1`
- **认证方式**: JWT Bearer Token
- **数据格式**: JSON
- **字符编码**: UTF-8

## 认证机制

### JWT Token认证

所有需要认证的API都需要在请求头中包含JWT Token：

```http
Authorization: Bearer <your-jwt-token>
```

### Token获取

通过登录接口获取Token：

```http
POST /api/v1/auth/login
```

### Token刷新

Token过期前可通过刷新接口获取新Token：

```http
POST /api/v1/auth/refresh-token
```

## 用户角色和权限

### 角色定义

| 角色 | 代码 | 描述 | 权限范围 |
|------|------|------|----------|
| 平台管理员 | `platform_admin` | 系统超级管理员 | 所有权限 |
| 品牌管理员 | `brand_admin` | 品牌业务管理员 | 分配品牌的管理权限 |
| 参与者 | `participant` | 普通用户 | 基础查看权限 |
| 匿名用户 | `anonymous` | 未登录用户 | 公开接口权限 |

### 权限编码规范

权限编码采用 `资源:操作` 格式：

- `user:manage` - 用户管理权限
- `brand:manage` - 品牌管理权限
- `campaign:manage` - 活动管理权限
- `order:read` - 订单查看权限
- `security:manage` - 安全管理权限

## API接口分类

### 1. 认证授权接口

#### 1.1 用户登录

**接口地址**: `POST /api/v1/auth/login`

**请求参数**:
```json
{
  "username": "string",  // 用户名，3-50字符
  "password": "string"   // 密码，6位以上
}
```

**响应数据**:
```json
{
  "token": "string",     // JWT Token
  "userId": 123,         // 用户ID
  "username": "string",  // 用户名
  "phone": "string",     // 手机号
  "realName": "string",  // 真实姓名
  "roles": ["string"],   // 用户角色列表
  "brandIds": [123]      // 品牌管理员的品牌ID列表（可选）
}
```

**错误响应**:
```json
{
  "code": 400,
  "message": "用户名或密码错误"
}
```

#### 1.2 用户注册

**接口地址**: `POST /api/v1/auth/register`

**请求参数**:
```json
{
  "username": "string",     // 用户名，3-50字符，字母数字下划线
  "password": "string",     // 密码，需符合密码策略
  "phone": "string",        // 手机号，11位数字
  "email": "string",        // 邮箱（可选）
  "realName": "string"      // 真实姓名（可选）
}
```

**响应数据**: 同登录接口

**注意**: 注册接口只能创建`participant`角色用户

#### 1.3 Token刷新

**接口地址**: `POST /api/v1/auth/refresh-token`

**请求参数**:
```json
{
  "token": "string"  // 当前Token
}
```

**响应数据**:
```json
{
  "token": "string"  // 新的Token
}
```

#### 1.4 获取用户信息

**接口地址**: `GET /api/v1/auth/userinfo`

**认证**: 需要Token

**响应数据**:
```json
{
  "id": 123,
  "username": "string",
  "phone": "string",
  "email": "string",
  "realName": "string",
  "avatar": "string",
  "status": "active",
  "roles": ["string"],
  "brandIds": [123],
  "createdAt": "2024-01-01 12:00:00"
}
```

### 2. 用户管理接口

#### 2.1 创建用户（管理员）

**接口地址**: `POST /api/v1/admin/users`

**权限要求**: `user:manage`

**请求参数**:
```json
{
  "username": "string",     // 用户名
  "password": "string",     // 密码
  "phone": "string",        // 手机号
  "email": "string",        // 邮箱（可选）
  "realName": "string",     // 真实姓名（可选）
  "role": "string",         // 角色：platform_admin/brand_admin/participant
  "brandIds": [123]         // 品牌ID列表（品牌管理员必填）
}
```

**响应数据**: 用户信息对象

#### 2.2 获取用户列表

**接口地址**: `GET /api/v1/admin/users`

**权限要求**: `user:manage`

**查询参数**:
- `page`: 页码（默认1）
- `pageSize`: 每页数量（默认20）
- `role`: 角色过滤
- `status`: 状态过滤
- `keyword`: 关键词搜索

**响应数据**:
```json
{
  "total": 100,
  "users": [
    {
      "id": 123,
      "username": "string",
      "phone": "string",
      "email": "string",
      "realName": "string",
      "status": "active",
      "roles": ["string"],
      "brandIds": [123],
      "createdAt": "2024-01-01 12:00:00"
    }
  ]
}
```

#### 2.3 更新用户信息

**接口地址**: `PUT /api/v1/admin/users/{id}`

**权限要求**: `user:manage`

**请求参数**:
```json
{
  "realName": "string",     // 真实姓名（可选）
  "email": "string",        // 邮箱（可选）
  "role": "string",         // 角色（可选）
  "status": "string",       // 状态：active/disabled/locked（可选）
  "brandIds": [123]         // 品牌ID列表（可选）
}
```

#### 2.4 重置用户密码

**接口地址**: `POST /api/v1/admin/users/{id}/reset-password`

**权限要求**: `user:manage`

**请求参数**:
```json
{
  "newPassword": "string"  // 新密码
}
```

#### 2.5 删除用户

**接口地址**: `DELETE /api/v1/admin/users/{id}`

**权限要求**: `user:manage`

### 3. 品牌管理接口

#### 3.1 创建品牌

**接口地址**: `POST /api/v1/brands`

**权限要求**: `brand:manage`

**请求参数**:
```json
{
  "name": "string",         // 品牌名称
  "logo": "string",         // 品牌Logo URL（可选）
  "description": "string"   // 品牌描述（可选）
}
```

**响应数据**:
```json
{
  "id": 123,
  "name": "string",
  "logo": "string",
  "description": "string",
  "status": "active",
  "createdAt": "2024-01-01 12:00:00",
  "updatedAt": "2024-01-01 12:00:00"
}
```

#### 3.2 获取品牌列表

**接口地址**: `GET /api/v1/brands`

**权限要求**: `brand:manage`

**响应数据**:
```json
{
  "total": 50,
  "brands": [
    {
      "id": 123,
      "name": "string",
      "logo": "string",
      "description": "string",
      "status": "active",
      "createdAt": "2024-01-01 12:00:00",
      "updatedAt": "2024-01-01 12:00:00"
    }
  ]
}
```

**注意**: 品牌管理员只能看到分配给自己的品牌

#### 3.3 获取品牌详情

**接口地址**: `GET /api/v1/brands/{id}`

**权限要求**: `brand:manage` + 品牌权限

#### 3.4 更新品牌信息

**接口地址**: `PUT /api/v1/brands/{id}`

**权限要求**: `brand:manage` + 品牌权限

**请求参数**:
```json
{
  "name": "string",         // 品牌名称（可选）
  "logo": "string",         // 品牌Logo URL（可选）
  "description": "string",  // 品牌描述（可选）
  "status": "string"        // 状态：active/disabled（可选）
}
```

#### 3.5 获取品牌统计

**接口地址**: `GET /api/v1/brands/{id}/stats`

**权限要求**: `brand:manage` + 品牌权限

**响应数据**:
```json
{
  "brandId": 123,
  "totalCampaigns": 50,
  "activeCampaigns": 10,
  "totalOrders": 1000,
  "totalRevenue": 50000.00,
  "totalRewards": 5000.00,
  "participantCount": 500,
  "conversionRate": 0.15,
  "lastUpdated": "2024-01-01 12:00:00"
}
```

### 4. 品牌素材管理接口

#### 4.1 创建品牌素材

**接口地址**: `POST /api/v1/brands/{id}/assets`

**权限要求**: `brand:manage` + 品牌权限

**请求参数**:
```json
{
  "name": "string",         // 素材名称
  "type": "string",         // 素材类型：image/video/document
  "category": "string",     // 素材分类（可选）
  "tags": "string",         // 素材标签（可选）
  "fileUrl": "string",      // 文件URL
  "fileSize": 1024,         // 文件大小（字节）
  "description": "string"   // 素材描述（可选）
}
```

#### 4.2 获取品牌素材列表

**接口地址**: `GET /api/v1/brands/{id}/assets`

**权限要求**: `brand:manage` + 品牌权限

**响应数据**:
```json
{
  "total": 20,
  "assets": [
    {
      "id": 123,
      "brandId": 456,
      "name": "string",
      "type": "image",
      "category": "logo",
      "tags": "official,primary",
      "fileUrl": "string",
      "fileSize": 1024,
      "description": "string",
      "createdAt": "2024-01-01 12:00:00",
      "updatedAt": "2024-01-01 12:00:00"
    }
  ]
}
```

### 5. 角色权限管理接口

#### 5.1 获取角色列表

**接口地址**: `GET /api/v1/roles`

**权限要求**: `role:read`

**响应数据**:
```json
[
  {
    "id": 1,
    "name": "平台管理员",
    "code": "platform_admin",
    "description": "系统超级管理员",
    "permissions": ["user:manage", "brand:manage"],
    "createdAt": "2024-01-01 12:00:00"
  }
]
```

#### 5.2 获取权限列表

**接口地址**: `GET /api/v1/permissions`

**权限要求**: `role:read`

**响应数据**:
```json
[
  {
    "id": 1,
    "name": "用户管理",
    "code": "user:manage",
    "resource": "user",
    "action": "manage",
    "description": "用户管理权限"
  }
]
```

#### 5.3 配置角色权限

**接口地址**: `POST /api/v1/roles/permissions`

**权限要求**: `role:config`

**请求参数**:
```json
{
  "roleId": 1,
  "permissionIds": [1, 2, 3]
}
```

#### 5.4 获取用户权限

**接口地址**: `GET /api/v1/users/{id}/permissions`

**权限要求**: `role:read`

**响应数据**:
```json
{
  "userId": 123,
  "roles": ["brand_admin"],
  "permissions": ["brand:manage", "campaign:manage"],
  "brandIds": [1, 2]
}
```

### 6. 菜单权限管理接口

#### 6.1 创建菜单

**接口地址**: `POST /api/v1/menus`

**权限要求**: `menu:create`

**请求参数**:
```json
{
  "name": "string",         // 菜单名称
  "code": "string",         // 菜单编码
  "path": "string",         // 菜单路径
  "icon": "string",         // 菜单图标（可选）
  "parentId": 123,          // 父菜单ID（可选）
  "sort": 1,                // 排序
  "type": "menu",           // 类型：menu/button
  "platform": "admin"       // 平台：admin/h5
}
```

#### 6.2 获取菜单列表

**接口地址**: `GET /api/v1/menus`

**权限要求**: `menu:read`

**响应数据**:
```json
[
  {
    "id": 1,
    "name": "用户管理",
    "code": "user-management",
    "path": "/users",
    "icon": "user",
    "parentId": null,
    "sort": 1,
    "type": "menu",
    "platform": "admin",
    "status": "active",
    "children": [],
    "createdAt": "2024-01-01 12:00:00"
  }
]
```

#### 6.3 获取用户菜单

**接口地址**: `GET /api/v1/users/menus`

**权限要求**: 已登录

**查询参数**:
- `platform`: 平台类型（admin/h5）

**响应数据**:
```json
{
  "userId": 123,
  "platform": "admin",
  "menus": [
    {
      "id": 1,
      "name": "用户管理",
      "code": "user-management",
      "path": "/users",
      "icon": "user",
      "children": []
    }
  ]
}
```

### 7. 安全管理接口

#### 7.1 获取密码策略

**接口地址**: `GET /api/v1/security/password-policy`

**权限要求**: `security:manage`

**响应数据**:
```json
{
  "id": 1,
  "minLength": 8,
  "requireUppercase": true,
  "requireLowercase": true,
  "requireNumbers": true,
  "requireSpecialChars": true,
  "maxAge": 90,
  "historyCount": 5,
  "maxLoginAttempts": 5,
  "lockoutDuration": 30,
  "sessionTimeout": 480,
  "maxConcurrentSessions": 3,
  "createdAt": "2024-01-01 12:00:00",
  "updatedAt": "2024-01-01 12:00:00"
}
```

#### 7.2 更新密码策略

**接口地址**: `PUT /api/v1/security/password-policy`

**权限要求**: `security:manage`

**请求参数**: 同获取密码策略响应数据（可选字段）

#### 7.3 检查密码强度

**接口地址**: `POST /api/v1/security/check-password-strength`

**权限要求**: 已登录

**请求参数**:
```json
{
  "newPassword": "string"  // 要检查的密码
}
```

**响应数据**:
```json
{
  "score": 85,             // 强度评分（0-100）
  "level": "强",           // 强度等级：很弱/弱/中等/强
  "message": "密码强度良好" // 建议信息
}
```

#### 7.4 获取审计日志

**接口地址**: `GET /api/v1/security/audit-logs`

**权限要求**: `security:manage`

**查询参数**:
- `page`: 页码
- `pageSize`: 每页数量
- `username`: 用户名过滤
- `action`: 操作类型过滤
- `startTime`: 开始时间
- `endTime`: 结束时间

**响应数据**:
```json
{
  "total": 1000,
  "logs": [
    {
      "id": 1,
      "userId": 123,
      "username": "admin",
      "action": "create_user",
      "resource": "user",
      "resourceId": "456",
      "details": "创建用户详情",
      "clientIp": "192.168.1.100",
      "userAgent": "Mozilla/5.0",
      "status": "success",
      "errorMsg": "",
      "createdAt": "2024-01-01 12:00:00"
    }
  ]
}
```

#### 7.5 获取用户会话

**接口地址**: `GET /api/v1/security/sessions`

**权限要求**: `security:manage`

**响应数据**:
```json
{
  "total": 50,
  "sessions": [
    {
      "id": "session-id-123",
      "userId": 123,
      "clientIp": "192.168.1.100",
      "userAgent": "Mozilla/5.0",
      "loginAt": "2024-01-01 12:00:00",
      "lastActiveAt": "2024-01-01 13:00:00",
      "expiresAt": "2024-01-01 20:00:00",
      "status": "active",
      "createdAt": "2024-01-01 12:00:00"
    }
  ]
}
```

#### 7.6 撤销会话

**接口地址**: `DELETE /api/v1/security/sessions/{sessionId}`

**权限要求**: `security:manage`

#### 7.7 强制用户下线

**接口地址**: `POST /api/v1/security/force-logout/{userId}`

**权限要求**: `security:manage`

**请求参数**:
```json
{
  "reason": "管理员操作"  // 下线原因（可选）
}
```

## 错误码说明

### HTTP状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或Token无效 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 422 | 数据验证失败 |
| 500 | 服务器内部错误 |

### 业务错误码

| 错误码 | 说明 |
|--------|------|
| 10001 | 用户名或密码错误 |
| 10002 | 用户不存在 |
| 10003 | 用户已被禁用 |
| 10004 | 用户已被锁定 |
| 10005 | 密码已过期 |
| 10006 | 登录失败次数过多 |
| 20001 | 权限不足 |
| 20002 | 品牌权限不足 |
| 20003 | 数据访问被拒绝 |
| 30001 | 参数验证失败 |
| 30002 | 数据格式错误 |
| 30003 | 必填参数缺失 |
| 40001 | 资源不存在 |
| 40002 | 资源已存在 |
| 40003 | 资源状态错误 |
| 50001 | 系统内部错误 |
| 50002 | 数据库操作失败 |
| 50003 | 外部服务调用失败 |

### 错误响应格式

```json
{
  "code": 10001,
  "message": "用户名或密码错误",
  "details": "详细错误信息（可选）",
  "timestamp": "2024-01-01T12:00:00Z",
  "path": "/api/v1/auth/login"
}
```

## 数据隔离说明

### 品牌数据隔离

品牌管理员只能访问分配给自己的品牌数据：

1. **品牌信息**: 只能查看和管理分配的品牌
2. **品牌素材**: 只能管理分配品牌的素材
3. **活动数据**: 只能管理分配品牌的活动
4. **订单数据**: 只能查看分配品牌的订单

### 权限检查流程

1. **Token验证**: 验证JWT Token有效性
2. **用户状态检查**: 检查用户是否被禁用或锁定
3. **权限验证**: 检查用户是否有接口所需权限
4. **数据权限验证**: 检查用户是否有访问特定数据的权限
5. **品牌权限验证**: 对于品牌相关接口，检查品牌权限

## 最佳实践

### 1. 认证和授权

- 始终使用HTTPS传输敏感数据
- Token应存储在安全的地方（如HttpOnly Cookie）
- 定期刷新Token，避免长期有效Token
- 实现Token黑名单机制

### 2. 权限设计

- 遵循最小权限原则
- 使用角色-权限模型，避免直接分配权限
- 定期审查和清理无用权限
- 实现权限继承和覆盖机制

### 3. 数据安全

- 敏感数据加密存储
- 实现数据访问审计
- 定期备份重要数据
- 实现数据脱敏机制

### 4. 性能优化

- 使用权限缓存减少数据库查询
- 实现分页查询避免大数据量返回
- 使用索引优化数据库查询
- 实现接口限流防止滥用

## 版本更新说明

### v1.0 (当前版本)

- 完整的RBAC权限管理功能
- 品牌数据隔离机制
- 安全审计和会话管理
- 密码策略和强度检查
- 菜单权限动态控制

### 后续版本规划

- API版本控制机制
- 更细粒度的权限控制
- 多租户支持
- 国际化支持
- 更多安全特性