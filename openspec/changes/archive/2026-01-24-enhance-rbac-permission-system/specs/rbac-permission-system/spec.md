# Spec: RBAC权限管理系统

**Module**: RBAC Permission System  
**Priority**: P0  
**Status**: 🔄 In Progress  
**Related Proposal**: [enhance-rbac-permission-system](../../proposal.md)

---

## 📋 模块概述

RBAC权限管理系统是DMH平台的核心安全模块，实现基于角色的访问控制(Role-Based Access Control)。系统支持4种用户角色，提供完整的认证、授权和权限管理功能，确保不同角色用户只能访问其权限范围内的功能和数据。

---

## 🎯 功能需求

## ADDED Requirements

### Requirement: 用户认证管理
系统SHALL提供完整的用户认证功能，包括注册、登录、密码管理和会话控制。

#### Scenario: H5用户注册成功
- **WHEN** 用户通过H5页面提供有效的用户名、密码和手机号
- **THEN** 系统SHALL创建新用户账号并分配默认角色(participant)
- **AND** 系统SHALL返回JWT token和用户信息

#### Scenario: 平台管理员后台创建用户
- **WHEN** 平台管理员通过后台管理系统创建用户
- **THEN** 系统SHALL创建用户账号并分配指定角色
- **AND** 平台管理员角色只能通过后台系统创建，不允许H5注册

#### Scenario: 用户登录成功
- **WHEN** 用户提供正确的用户名和密码
- **THEN** 系统SHALL验证用户身份并生成JWT token
- **AND** token SHALL包含用户ID、角色信息和有效期

#### Scenario: 密码安全验证
- **WHEN** 用户设置或修改密码
- **THEN** 系统SHALL使用bcrypt加密存储密码
- **AND** 密码SHALL满足最小安全要求(长度≥6位)

#### Scenario: 会话超时控制
- **WHEN** JWT token超过有效期(24小时)
- **THEN** 系统SHALL拒绝请求并返回401未授权错误
- **AND** 前端SHALL清除本地token并跳转到登录页

### Requirement: 角色权限体系
系统SHALL实现4种用户角色，每种角色具有明确的权限范围和功能边界。

#### Scenario: 平台管理员权限
- **WHEN** 平台管理员(platform_admin)访问任何功能
- **THEN** 系统SHALL允许访问所有功能模块
- **AND** 可以管理所有品牌、活动、用户和系统配置

#### Scenario: 品牌管理员权限
- **WHEN** 品牌管理员(brand_admin)访问品牌相关功能
- **THEN** 系统SHALL只允许访问其管理的品牌数据
- **AND** 可以管理品牌信息(编辑品牌资料、上传品牌logo等)
- **AND** 可以管理品牌素材库(上传、分类、删除素材)
- **AND** 可以创建、编辑、删除、发布本品牌的活动
- **AND** 可以查看本品牌的数据统计和报表

#### Scenario: 活动参与者权限
- **WHEN** 活动参与者(participant)访问系统功能
- **THEN** 系统SHALL只允许访问个人相关功能
- **AND** 可以参与活动、查看个人奖励和申请提现

#### Scenario: 匿名用户权限
- **WHEN** 匿名用户(anonymous)访问系统
- **THEN** 系统SHALL只允许访问公开功能
- **AND** 可以浏览活动列表、查看活动详情和注册账号

### Requirement: API权限控制
系统SHALL在API层面实现细粒度的权限控制，确保每个接口都有适当的权限检查。

#### Scenario: JWT token验证
- **WHEN** 客户端请求需要认证的API接口
- **THEN** 系统SHALL验证Authorization header中的JWT token
- **AND** token无效时SHALL返回401未授权错误

#### Scenario: 权限检查机制
- **WHEN** 用户访问受保护的API接口
- **THEN** 系统SHALL根据URL和HTTP方法确定所需权限
- **AND** 检查用户角色是否具有该权限

#### Scenario: 数据级权限隔离
- **WHEN** 品牌管理员查询活动数据
- **THEN** 系统SHALL只返回其管理品牌的活动数据
- **AND** 不能访问其他品牌的数据

### Requirement: 用户注册权限控制
系统SHALL根据用户角色类型实现不同的注册方式和权限控制。

#### Scenario: H5注册限制角色
- **WHEN** 用户通过H5页面注册
- **THEN** 系统SHALL只允许创建participant角色的用户
- **AND** 不允许通过H5注册创建管理员角色

#### Scenario: 品牌管理员角色分配
- **WHEN** 需要创建品牌管理员用户
- **THEN** 系统SHALL要求平台管理员通过后台管理系统操作
- **AND** 同时在brand_admins表中建立品牌关联关系

#### Scenario: 平台管理员创建限制
- **WHEN** 需要创建平台管理员用户
- **THEN** 系统SHALL只允许现有平台管理员通过后台系统创建
- **AND** 平台管理员角色不能通过任何前端注册方式获得

#### Scenario: 匿名用户转换
- **WHEN** 匿名用户完成H5注册流程
- **THEN** 系统SHALL将其转换为participant角色
- **AND** 获得相应的登录凭据和基础权限

### Requirement: 用户管理功能
系统SHALL提供完整的用户管理功能，支持用户创建、查询、更新、状态管理和密码重置。

#### Scenario: 后台创建用户账号
- **WHEN** 平台管理员通过后台管理系统创建新用户
- **THEN** 系统SHALL验证用户名和手机号的唯一性
- **AND** 创建用户记录并分配指定角色
- **AND** 生成初始密码并通知用户

#### Scenario: 用户状态管理
- **WHEN** 平台管理员变更用户状态
- **THEN** 系统SHALL更新用户状态(active/disabled/locked)
- **AND** 禁用用户SHALL立即失去系统访问权限
- **AND** 记录状态变更日志和操作人

#### Scenario: 用户密码重置
- **WHEN** 平台管理员重置用户密码
- **THEN** 系统SHALL生成新的临时密码
- **AND** 强制用户下次登录时修改密码
- **AND** 记录密码重置操作日志

#### Scenario: 用户角色分配
- **WHEN** 平台管理员为用户分配角色
- **THEN** 系统SHALL更新用户角色关联关系
- **AND** 新角色权限SHALL立即生效
- **AND** 记录角色变更日志

#### Scenario: 品牌管理员分配
- **WHEN** 平台管理员指定用户为品牌管理员
- **THEN** 系统SHALL在brand_admins表中创建关联记录
- **AND** 用户SHALL获得该品牌的管理权限
- **AND** 可以同时管理多个品牌

### Requirement: 权限缓存优化
系统SHALL实现权限信息缓存机制，提高权限检查的性能和响应速度。

#### Scenario: 权限信息缓存
- **WHEN** 系统首次查询用户权限信息
- **THEN** 系统SHALL将权限信息缓存到内存中
- **AND** 后续权限检查SHALL优先使用缓存数据

#### Scenario: 缓存失效更新
- **WHEN** 用户角色或权限发生变更
- **THEN** 系统SHALL立即清除相关缓存
- **AND** 下次权限检查SHALL重新查询数据库

### Requirement: 品牌管理员关系管理
系统SHALL为平台管理员提供品牌管理员与品牌关系的完整管理功能，支持绑定、解绑和变更操作。

#### Scenario: 绑定品牌管理员
- **WHEN** 平台管理员为用户绑定品牌管理权限
- **THEN** 系统SHALL在brand_admins表中创建关联记录
- **AND** 用户SHALL立即获得该品牌的管理权限
- **AND** 记录绑定操作日志

#### Scenario: 解绑品牌管理员
- **WHEN** 平台管理员解除用户的品牌管理权限
- **THEN** 系统SHALL删除brand_admins表中的关联记录
- **AND** 用户SHALL立即失去该品牌的管理权限
- **AND** 记录解绑操作日志

#### Scenario: 变更品牌管理员权限
- **WHEN** 平台管理员调整品牌管理员的品牌范围
- **THEN** 系统SHALL更新brand_admins表中的关联记录
- **AND** 新的品牌权限SHALL立即生效
- **AND** 记录权限变更日志

#### Scenario: 多品牌管理支持
- **WHEN** 品牌管理员被分配多个品牌
- **THEN** 系统SHALL支持一个用户管理多个品牌
- **AND** 在数据查询时SHALL正确过滤各品牌数据
- **AND** 权限检查SHALL验证用户对特定品牌的访问权限

#### Scenario: 品牌管理员权限查询
- **WHEN** 查询用户的品牌管理权限
- **THEN** 系统SHALL返回用户管理的所有品牌列表
- **AND** 包含品牌基本信息和权限范围
- **AND** 支持按品牌状态过滤
系统SHALL为品牌管理员提供完整的品牌管理功能，包括品牌信息、素材库、活动和数据管理。

#### Scenario: 品牌信息管理
- **WHEN** 品牌管理员编辑品牌信息
- **THEN** 系统SHALL允许修改品牌名称、描述、logo等基本信息
- **AND** 只能修改其管理的品牌信息

#### Scenario: 品牌素材库管理
- **WHEN** 品牌管理员管理素材库
- **THEN** 系统SHALL允许上传、分类、编辑、删除品牌素材
- **AND** 素材包括图片、视频、文档等多种类型
- **AND** 只能管理本品牌的素材资源

#### Scenario: 品牌活动管理
- **WHEN** 品牌管理员管理活动
- **THEN** 系统SHALL允许创建、编辑、删除、发布本品牌的活动
- **AND** 可以配置活动的动态表单和奖励规则
- **AND** 不能访问其他品牌的活动

#### Scenario: 品牌数据查看
- **WHEN** 品牌管理员查看数据统计
- **THEN** 系统SHALL提供本品牌的完整数据报表
- **AND** 包括活动参与数、订单统计、奖励发放、用户分析等
- **AND** 不能查看其他品牌的数据
系统SHALL实现提现申请和审核的权限控制，确保只有授权用户可以进行相关操作。

#### Scenario: 提现申请权限
- **WHEN** 用户申请提现
- **THEN** 系统SHALL检查用户是否为已认证用户
- **AND** 验证用户余额是否足够

#### Scenario: 提现审核权限
- **WHEN** 用户尝试审核提现申请
- **THEN** 系统SHALL验证用户是否为平台管理员
- **AND** 只有平台管理员可以批准或拒绝提现

#### Scenario: 提现状态更新
- **WHEN** 平台管理员审核提现申请
- **THEN** 系统SHALL使用数据库事务确保数据一致性
- **AND** 记录审核人和审核时间

### Requirement: 品牌管理功能
系统SHALL实现提现申请和审核的权限控制，确保只有授权用户可以进行相关操作。

#### Scenario: 提现申请权限
- **WHEN** 用户申请提现
- **THEN** 系统SHALL检查用户是否为已认证用户
- **AND** 验证用户余额是否足够

#### Scenario: 提现审核权限
- **WHEN** 用户尝试审核提现申请
- **THEN** 系统SHALL验证用户是否为平台管理员
- **AND** 只有平台管理员可以批准或拒绝提现

#### Scenario: 提现状态更新
- **WHEN** 平台管理员审核提现申请
- **THEN** 系统SHALL使用数据库事务确保数据一致性
- **AND** 记录审核人和审核时间

### Requirement: 菜单权限管理
系统SHALL提供完整的菜单权限管理功能，支持菜单结构管理和角色菜单权限配置。

#### Scenario: 菜单结构管理
- **WHEN** 平台管理员管理菜单结构
- **THEN** 系统SHALL支持菜单的增加、删除、修改和排序
- **AND** 支持多级菜单结构和菜单分组
- **AND** 区分后台管理菜单和H5用户菜单

#### Scenario: 页面操作权限配置
- **WHEN** 配置页面操作权限
- **THEN** 系统SHALL支持增删改查、导出、转发等操作权限
- **AND** 每个菜单项可配置多种操作权限
- **AND** 支持按钮级别的权限控制

#### Scenario: 角色菜单权限分配
- **WHEN** 为角色分配菜单权限
- **THEN** 系统SHALL支持选择性分配菜单访问权限
- **AND** 支持为每个菜单配置具体的操作权限
- **AND** 权限变更SHALL立即生效

#### Scenario: 用户菜单权限查询
- **WHEN** 用户登录系统
- **THEN** 系统SHALL根据用户角色返回可访问的菜单列表
- **AND** 包含每个菜单的操作权限信息
- **AND** 前端根据权限动态显示菜单和按钮

#### Scenario: 菜单权限继承
- **WHEN** 配置多级菜单权限
- **THEN** 系统SHALL支持权限继承机制
- **AND** 子菜单可继承父菜单的权限设置
- **AND** 支持覆盖继承的权限配置

### Requirement: 安全审计日志
系统SHALL记录所有重要的安全相关操作，提供完整的审计追踪能力。

#### Scenario: 用户操作日志
- **WHEN** 用户执行重要操作(登录、权限变更、数据修改)
- **THEN** 系统SHALL记录操作日志
- **AND** 日志SHALL包含用户ID、操作类型、时间戳和IP地址

#### Scenario: 权限变更日志
- **WHEN** 管理员修改用户角色或权限
- **THEN** 系统SHALL记录权限变更日志
- **AND** 日志SHALL包含变更前后的权限状态

#### Scenario: 安全事件监控
- **WHEN** 检测到异常登录或权限滥用
- **THEN** 系统SHALL记录安全事件
- **AND** 可选择性地触发安全告警

---

## 🔌 API 接口定义

### 1. 用户认证接口

#### H5用户注册（仅限participant角色）
```
POST /api/v1/auth/register
Content-Type: application/json

Request Body:
{
  "username": "testuser",
  "password": "123456",
  "phone": "13800138000",
  "email": "test@example.com"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "phone": "13800138000",
    "email": "test@example.com",
    "roles": ["participant"],
    "status": "active"
  }
}

Note: H5注册只能创建participant角色用户
```

#### 后台管理员创建用户（支持所有角色）
```
POST /api/v1/admin/users
Authorization: Bearer <admin_token>
Content-Type: application/json

Request Body:
{
  "username": "brandmanager",
  "password": "123456",
  "phone": "13800138001",
  "email": "brand@example.com",
  "roles": ["brand_admin"],
  "brandIds": [1, 2]
}

Response:
{
  "id": 2,
  "username": "brandmanager",
  "phone": "13800138001",
  "email": "brand@example.com",
  "roles": ["brand_admin"],
  "brandIds": [1, 2],
  "status": "active",
  "createdAt": "2025-01-02T10:00:00Z"
}

Note: 只有平台管理员可以通过此接口创建用户
```

#### 用户登录
```
POST /api/v1/auth/login
Content-Type: application/json

Request Body:
{
  "username": "admin",
  "password": "123456"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin",
    "roles": ["platform_admin"],
    "brandIds": [],
    "permissions": ["*"]
  }
}
```

#### 获取用户信息
```
GET /api/v1/auth/userinfo
Authorization: Bearer <token>

Response:
{
  "id": 1,
  "username": "admin",
  "phone": "13800138000",
  "email": "admin@example.com",
  "roles": ["platform_admin"],
  "brandIds": [],
  "permissions": ["*"],
  "status": "active"
}
```

### 2. 用户管理接口

#### 后台创建用户（管理员专用）
```
POST /api/v1/admin/users
Authorization: Bearer <admin_token>
Content-Type: application/json

Request Body:
{
  "username": "newuser",
  "password": "123456",
  "phone": "13800138001",
  "email": "newuser@example.com",
  "roles": ["brand_admin"],
  "brandIds": [1, 2]
}

Response:
{
  "id": 2,
  "username": "newuser",
  "phone": "13800138001",
  "email": "newuser@example.com",
  "roles": ["brand_admin"],
  "brandIds": [1, 2],
  "status": "active",
  "createdAt": "2025-01-02T10:00:00Z"
}

Note: 此接口只能由平台管理员调用，支持创建所有角色用户
```

#### 获取用户列表
```
GET /api/v1/users?page=1&pageSize=20&role=brand_admin
Authorization: Bearer <admin_token>

Response:
{
  "total": 50,
  "users": [
    {
      "id": 2,
      "username": "branduser",
      "phone": "13800138001",
      "roles": ["brand_admin"],
      "brandIds": [1],
      "status": "active",
      "createdAt": "2025-01-02T10:00:00Z"
    }
  ]
}
```

#### 更新用户状态
```
PUT /api/v1/users/:id/status
Authorization: Bearer <admin_token>
Content-Type: application/json

Request Body:
{
  "status": "disabled",
  "reason": "违规操作"
}

Response:
{
  "id": 2,
  "username": "testuser",
  "status": "disabled",
  "updatedAt": "2025-01-02T10:00:00Z",
  "updatedBy": 1,
  "reason": "违规操作"
}
```

#### 重置用户密码
```
POST /api/v1/users/:id/reset-password
Authorization: Bearer <admin_token>
Content-Type: application/json

Request Body:
{
  "forceChange": true,
  "notifyUser": true
}

Response:
{
  "id": 2,
  "username": "testuser",
  "temporaryPassword": "Temp123456",
  "forceChange": true,
  "resetAt": "2025-01-02T10:00:00Z",
  "resetBy": 1
}
```

#### 分配用户角色
```
POST /api/v1/users/:id/roles
Authorization: Bearer <admin_token>
Content-Type: application/json

Request Body:
{
  "roleIds": [2, 3],
  "brandIds": [1, 2]
}

Response:
{
  "id": 2,
  "username": "testuser",
  "roles": [
    {
      "id": 2,
      "code": "brand_admin",
      "name": "品牌管理员"
    }
  ],
  "brandIds": [1, 2],
  "updatedAt": "2025-01-02T10:00:00Z"
}
```

#### 更新品牌信息
```
PUT /api/v1/brands/:id
Authorization: Bearer <brand_admin_token>
Content-Type: application/json

Request Body:
{
  "name": "更新的品牌名称",
  "description": "品牌描述",
  "logo": "https://example.com/logo.png",
  "website": "https://brand.example.com"
}

Response:
{
  "id": 1,
  "name": "更新的品牌名称",
  "description": "品牌描述",
  "logo": "https://example.com/logo.png",
  "website": "https://brand.example.com",
  "status": "active",
  "updatedAt": "2025-01-02T10:00:00Z"
}
```

#### 获取品牌素材列表
```
GET /api/v1/brands/:id/materials?page=1&pageSize=20&type=image
Authorization: Bearer <brand_admin_token>

Response:
{
  "total": 100,
  "materials": [
    {
      "id": 1,
      "name": "活动海报.jpg",
      "type": "image",
      "url": "https://example.com/materials/poster.jpg",
      "size": 1024000,
      "createdAt": "2025-01-02T10:00:00Z"
    }
  ]
}
```

#### 上传品牌素材
```
POST /api/v1/brands/:id/materials
Authorization: Bearer <brand_admin_token>
Content-Type: multipart/form-data

Request Body:
- file: [binary file data]
- name: "活动海报"
- category: "posters"

Response:
{
  "id": 2,
  "name": "活动海报",
  "type": "image",
  "url": "https://example.com/materials/poster2.jpg",
  "size": 2048000,
  "category": "posters",
  "createdAt": "2025-01-02T10:00:00Z"
}
```

#### 绑定品牌管理员
```
POST /api/v1/brands/:id/admins
Authorization: Bearer <admin_token>
Content-Type: application/json

Request Body:
{
  "userIds": [2, 3],
  "permissions": ["full"] // full: 完整权限, limited: 受限权限
}

Response:
{
  "brandId": 1,
  "brandName": "品牌A",
  "admins": [
    {
      "userId": 2,
      "username": "brand_manager1",
      "permissions": ["full"],
      "assignedAt": "2025-01-02T10:00:00Z",
      "assignedBy": 1
    }
  ]
}
```

#### 解绑品牌管理员
```
DELETE /api/v1/brands/:id/admins/:userId
Authorization: Bearer <admin_token>

Response:
{
  "brandId": 1,
  "userId": 2,
  "message": "品牌管理员权限已解除",
  "removedAt": "2025-01-02T10:00:00Z",
  "removedBy": 1
}
```

#### 获取品牌管理员列表
```
GET /api/v1/brands/:id/admins
Authorization: Bearer <admin_token>

Response:
{
  "brandId": 1,
  "brandName": "品牌A",
  "admins": [
    {
      "userId": 2,
      "username": "brand_manager1",
      "realName": "张三",
      "phone": "13800138001",
      "permissions": ["full"],
      "assignedAt": "2025-01-02T10:00:00Z",
      "status": "active"
    }
  ]
}
```

#### 获取用户管理的品牌列表
```
GET /api/v1/users/:id/brands
Authorization: Bearer <admin_token>

Response:
{
  "userId": 2,
  "username": "brand_manager1",
  "brands": [
    {
      "brandId": 1,
      "brandName": "品牌A",
      "permissions": ["full"],
      "assignedAt": "2025-01-02T10:00:00Z",
      "status": "active"
    },
    {
      "brandId": 2,
      "brandName": "品牌B", 
      "permissions": ["limited"],
      "assignedAt": "2025-01-02T11:00:00Z",
      "status": "active"
    }
  ]
}
```

#### 批量管理品牌管理员
```
PUT /api/v1/users/:id/brands
Authorization: Bearer <admin_token>
Content-Type: application/json

Request Body:
{
  "brandIds": [1, 2, 3],
  "action": "assign", // assign: 分配, remove: 移除, replace: 替换
  "permissions": ["full"]
}

Response:
{
  "userId": 2,
  "username": "brand_manager1",
  "updatedBrands": [
    {
      "brandId": 1,
      "brandName": "品牌A",
      "action": "assigned",
      "permissions": ["full"]
    }
  ],
  "updatedAt": "2025-01-02T10:00:00Z"
}
```
```
GET /api/v1/brands/:id/statistics?startDate=2025-01-01&endDate=2025-01-31
Authorization: Bearer <brand_admin_token>

Response:
{
  "brandId": 1,
  "period": {
    "startDate": "2025-01-01",
    "endDate": "2025-01-31"
  },
  "campaigns": {
    "total": 10,
    "active": 5,
    "completed": 5
  },
  "participants": {
    "total": 1000,
    "new": 200
  },
  "orders": {
    "total": 800,
    "totalAmount": 80000.00,
    "avgAmount": 100.00
  },
  "rewards": {
    "totalPaid": 8000.00,
    "totalParticipants": 500
  }
}
```

#### 获取品牌数据统计

#### 获取角色列表
```
GET /api/v1/roles
Authorization: Bearer <admin_token>

Response:
{
  "roles": [
    {
      "id": 1,
      "name": "平台管理员",
      "code": "platform_admin",
      "description": "系统最高权限",
      "permissions": ["*"]
    },
    {
      "id": 2,
      "name": "品牌管理员", 
      "code": "brand_admin",
      "description": "品牌级权限",
      "permissions": ["campaign:create", "campaign:read", "campaign:update", "order:read"]
    }
  ]
}
```

### 5. 菜单权限管理接口

#### 获取菜单列表
```
GET /api/v1/menus?platform=admin&roleCode=brand_admin
Authorization: Bearer <admin_token>

Response:
{
  "menus": [
    {
      "id": 1,
      "name": "活动管理",
      "path": "/campaigns",
      "icon": "campaign",
      "parentId": null,
      "sort": 1,
      "platform": "admin",
      "permissions": ["read", "create", "update", "delete"],
      "children": [
        {
          "id": 2,
          "name": "活动列表",
          "path": "/campaigns/list",
          "icon": "list",
          "parentId": 1,
          "sort": 1,
          "permissions": ["read", "create", "update", "delete"]
        }
      ]
    }
  ]
}
```

#### 创建菜单项
```
POST /api/v1/menus
Authorization: Bearer <admin_token>
Content-Type: application/json

Request Body:
{
  "name": "品牌管理",
  "path": "/brands",
  "icon": "brand",
  "parentId": null,
  "sort": 2,
  "platform": "admin",
  "description": "品牌信息和素材管理"
}

Response:
{
  "id": 3,
  "name": "品牌管理",
  "path": "/brands",
  "icon": "brand",
  "parentId": null,
  "sort": 2,
  "platform": "admin",
  "description": "品牌信息和素材管理",
  "createdAt": "2025-01-02T10:00:00Z"
}
```

#### 配置角色菜单权限
```
POST /api/v1/roles/:roleId/menu-permissions
Authorization: Bearer <admin_token>
Content-Type: application/json

Request Body:
{
  "menuPermissions": [
    {
      "menuId": 1,
      "permissions": ["read", "create", "update"]
    },
    {
      "menuId": 2,
      "permissions": ["read", "update", "delete", "export"]
    }
  ]
}

Response:
{
  "roleId": 2,
  "roleName": "品牌管理员",
  "menuPermissions": [
    {
      "menuId": 1,
      "menuName": "活动管理",
      "permissions": ["read", "create", "update"],
      "updatedAt": "2025-01-02T10:00:00Z"
    }
  ]
}
```

#### 获取用户菜单权限
```
GET /api/v1/auth/user-menus?platform=admin
Authorization: Bearer <token>

Response:
{
  "userId": 2,
  "platform": "admin",
  "menus": [
    {
      "id": 1,
      "name": "活动管理",
      "path": "/campaigns",
      "icon": "campaign",
      "permissions": ["read", "create", "update"],
      "children": [
        {
          "id": 2,
          "name": "活动列表",
          "path": "/campaigns/list",
          "permissions": ["read", "create", "update", "export"]
        }
      ]
    }
  ]
}
```

#### 更新菜单项
```
PUT /api/v1/menus/:id
Authorization: Bearer <admin_token>
Content-Type: application/json

Request Body:
{
  "name": "活动管理中心",
  "path": "/campaigns",
  "icon": "campaign-new",
  "sort": 1,
  "description": "营销活动管理中心"
}

Response:
{
  "id": 1,
  "name": "活动管理中心",
  "path": "/campaigns",
  "icon": "campaign-new",
  "sort": 1,
  "description": "营销活动管理中心",
  "updatedAt": "2025-01-02T10:00:00Z"
}
```

#### 删除菜单项
```
DELETE /api/v1/menus/:id
Authorization: Bearer <admin_token>

Response:
{
  "id": 3,
  "message": "菜单项已删除",
  "deletedAt": "2025-01-02T10:00:00Z"
}
```
    }
  ]
}

Response:
{
  "message": "角色菜单权限配置成功",
  "roleId": 2,
  "updatedPermissions": 15
}
```

#### 获取用户菜单权限
```
GET /api/v1/auth/user-menus?platform=admin
Authorization: Bearer <user_token>

Response:
{
  "menus": [
    {
      "id": 1,
      "name": "活动管理",
      "path": "/campaigns",
      "icon": "campaign",
      "permissions": ["read", "create", "update"],
      "children": [
        {
          "id": 2,
          "name": "活动列表",
          "path": "/campaigns/list",
          "icon": "list",
          "permissions": ["read", "create", "update", "export"]
        }
      ]
    }
  ]
}
```

#### 权限验证接口
```
POST /api/v1/auth/verify-permission
Authorization: Bearer <token>
Content-Type: application/json

Request Body:
{
  "resource": "campaigns",
  "action": "create"
}

Response:
{
  "allowed": true,
  "reason": "User has platform_admin role"
}
```

---

## 💾 数据存储

### 用户相关表结构

#### users 表
```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码(bcrypt加密)',
    phone VARCHAR(20) UNIQUE COMMENT '手机号',
    email VARCHAR(100) UNIQUE COMMENT '邮箱',
    avatar VARCHAR(255) COMMENT '头像URL',
    real_name VARCHAR(50) COMMENT '真实姓名',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态:active,disabled,locked',
    force_password_change BOOLEAN DEFAULT FALSE COMMENT '强制修改密码',
    password_reset_at DATETIME COMMENT '密码重置时间',
    password_reset_by BIGINT COMMENT '密码重置操作人',
    last_login_at DATETIME COMMENT '最后登录时间',
    last_login_ip VARCHAR(45) COMMENT '最后登录IP',
    login_attempts INT DEFAULT 0 COMMENT '登录失败次数',
    locked_until DATETIME COMMENT '锁定到期时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_phone (phone),
    INDEX idx_email (email),
    INDEX idx_status (status),
    FOREIGN KEY (password_reset_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### user_status_logs 表
```sql
CREATE TABLE user_status_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    old_status VARCHAR(20) COMMENT '原状态',
    new_status VARCHAR(20) NOT NULL COMMENT '新状态',
    reason TEXT COMMENT '变更原因',
    operated_by BIGINT NOT NULL COMMENT '操作人ID',
    operated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (operated_by) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_operated_by (operated_by),
    INDEX idx_operated_at (operated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### roles 表
```sql
CREATE TABLE roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL COMMENT '角色名称',
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '角色编码',
    description TEXT COMMENT '角色描述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### user_roles 表
```sql
CREATE TABLE user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_role (user_id, role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### permissions 表
```sql
CREATE TABLE permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '权限名称',
    code VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码',
    resource VARCHAR(50) NOT NULL COMMENT '资源类型',
    action VARCHAR(50) NOT NULL COMMENT '操作类型',
    description TEXT COMMENT '权限描述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_code (code),
    INDEX idx_resource_action (resource, action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### role_permissions 表
```sql
CREATE TABLE role_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL COMMENT '角色ID',
    permission_id BIGINT NOT NULL COMMENT '权限ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    UNIQUE KEY uk_role_permission (role_id, permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### brand_admins 表
```sql
CREATE TABLE brand_admins (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    brand_id BIGINT NOT NULL COMMENT '品牌ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    permissions JSON COMMENT '权限配置',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态:active,disabled',
    assigned_by BIGINT COMMENT '分配人ID',
    assigned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_by) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE KEY uk_brand_user (brand_id, user_id),
    INDEX idx_user_id (user_id),
    INDEX idx_assigned_by (assigned_by),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### brand_admin_logs 表
```sql
CREATE TABLE brand_admin_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    brand_id BIGINT NOT NULL COMMENT '品牌ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    action VARCHAR(50) NOT NULL COMMENT '操作类型:assign,remove,update',
    old_permissions JSON COMMENT '原权限配置',
    new_permissions JSON COMMENT '新权限配置',
    reason TEXT COMMENT '操作原因',
    operated_by BIGINT NOT NULL COMMENT '操作人ID',
    operated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (operated_by) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_brand_id (brand_id),
    INDEX idx_user_id (user_id),
    INDEX idx_operated_by (operated_by),
    INDEX idx_operated_at (operated_at),
    INDEX idx_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 权限编码规范

权限编码采用 `资源:操作` 格式：

| 权限编码 | 权限名称 | 资源 | 操作 | 说明 |
|---------|---------|------|------|------|
| campaign:create | 创建活动 | campaign | create | 创建营销活动 |
| campaign:read | 查看活动 | campaign | read | 查看活动信息 |
| campaign:update | 编辑活动 | campaign | update | 编辑活动配置 |
| campaign:delete | 删除活动 | campaign | delete | 删除活动 |
| campaign:export | 导出活动 | campaign | export | 导出活动数据 |
| order:read | 查看订单 | order | read | 查看订单信息 |
| order:update | 更新订单 | order | update | 更新订单状态 |
| order:export | 导出订单 | order | export | 导出订单数据 |
| user:create | 创建用户 | user | create | 创建用户账号 |
| user:read | 查看用户 | user | read | 查看用户信息 |
| user:update | 更新用户 | user | update | 更新用户信息 |
| user:delete | 删除用户 | user | delete | 删除用户账号 |
| user:status | 管理用户状态 | user | status | 启用/禁用/锁定用户 |
| user:reset-password | 重置密码 | user | reset-password | 重置用户密码 |
| brand:create | 创建品牌 | brand | create | 创建品牌 |
| brand:read | 查看品牌 | brand | read | 查看品牌信息 |
| brand:update | 更新品牌 | brand | update | 更新品牌信息 |
| brand:delete | 删除品牌 | brand | delete | 删除品牌 |
| brand:assign-admin | 分配品牌管理员 | brand | assign-admin | 为品牌分配管理员 |
| brand:remove-admin | 移除品牌管理员 | brand | remove-admin | 移除品牌管理员 |
| brand:manage-admin | 管理品牌管理员 | brand | manage-admin | 管理品牌管理员关系 |
| material:create | 上传素材 | material | create | 上传品牌素材 |
| material:read | 查看素材 | material | read | 查看素材库 |
| material:update | 编辑素材 | material | update | 编辑素材信息 |
| material:delete | 删除素材 | material | delete | 删除素材文件 |
| menu:create | 创建菜单 | menu | create | 创建菜单项 |
| menu:read | 查看菜单 | menu | read | 查看菜单结构 |
| menu:update | 更新菜单 | menu | update | 更新菜单信息 |
| menu:delete | 删除菜单 | menu | delete | 删除菜单项 |
| role:create | 创建角色 | role | create | 创建用户角色 |
| role:read | 查看角色 | role | read | 查看角色信息 |
| role:update | 更新角色 | role | update | 更新角色权限 |
| role:delete | 删除角色 | role | delete | 删除角色 |
| statistics:read | 查看统计 | statistics | read | 查看数据统计 |
| statistics:export | 导出统计 | statistics | export | 导出统计数据 |
| menu:create | 创建菜单 | menu | create | 创建菜单项 |
| menu:read | 查看菜单 | menu | read | 查看菜单结构 |
| menu:update | 更新菜单 | menu | update | 更新菜单信息 |
| menu:delete | 删除菜单 | menu | delete | 删除菜单项 |
| menu:assign | 分配菜单权限 | menu | assign | 为角色分配菜单权限 |
| withdrawal:create | 申请提现 | withdrawal | create | 申请提现 |
| withdrawal:read | 查看提现 | withdrawal | read | 查看提现记录 |
| withdrawal:approve | 审核提现 | withdrawal | approve | 审核提现申请 |

#### menus 表
```sql
CREATE TABLE menus (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '菜单名称',
    path VARCHAR(200) COMMENT '菜单路径',
    icon VARCHAR(100) COMMENT '菜单图标',
    parent_id BIGINT COMMENT '父菜单ID',
    sort INT NOT NULL DEFAULT 0 COMMENT '排序',
    platform VARCHAR(20) NOT NULL DEFAULT 'admin' COMMENT '平台类型:admin,h5',
    description TEXT COMMENT '菜单描述',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态:active,disabled',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES menus(id) ON DELETE CASCADE,
    INDEX idx_parent_id (parent_id),
    INDEX idx_platform (platform),
    INDEX idx_sort (sort),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### menu_permissions 表
```sql
CREATE TABLE menu_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    menu_id BIGINT NOT NULL COMMENT '菜单ID',
    permission_type VARCHAR(50) NOT NULL COMMENT '权限类型:read,create,update,delete,export',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE,
    UNIQUE KEY uk_menu_permission (menu_id, permission_type),
    INDEX idx_menu_id (menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### role_menu_permissions 表
```sql
CREATE TABLE role_menu_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL COMMENT '角色ID',
    menu_id BIGINT NOT NULL COMMENT '菜单ID',
    permission_type VARCHAR(50) NOT NULL COMMENT '权限类型:read,create,update,delete,export',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE,
    UNIQUE KEY uk_role_menu_permission (role_id, menu_id, permission_type),
    INDEX idx_role_id (role_id),
    INDEX idx_menu_id (menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### brand_materials 表
```sql
CREATE TABLE brand_materials (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    brand_id BIGINT NOT NULL COMMENT '品牌ID',
    name VARCHAR(200) NOT NULL COMMENT '素材名称',
    original_name VARCHAR(255) COMMENT '原始文件名',
    type VARCHAR(50) NOT NULL COMMENT '素材类型:image,video,document',
    category VARCHAR(100) COMMENT '素材分类',
    url VARCHAR(500) NOT NULL COMMENT '素材URL',
    file_size BIGINT COMMENT '文件大小(字节)',
    mime_type VARCHAR(100) COMMENT 'MIME类型',
    width INT COMMENT '图片宽度',
    height INT COMMENT '图片高度',
    duration INT COMMENT '视频时长(秒)',
    created_by BIGINT COMMENT '创建者ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL COMMENT '软删除',
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_brand_id (brand_id),
    INDEX idx_type (type),
    INDEX idx_category (category),
    INDEX idx_created_by (created_by),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### menus 表
```sql
CREATE TABLE menus (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '菜单名称',
    path VARCHAR(200) COMMENT '菜单路径',
    icon VARCHAR(100) COMMENT '菜单图标',
    parent_id BIGINT COMMENT '父菜单ID',
    sort INT NOT NULL DEFAULT 0 COMMENT '排序',
    platform VARCHAR(20) NOT NULL DEFAULT 'admin' COMMENT '平台:admin,h5',
    description TEXT COMMENT '菜单描述',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态:active,disabled',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES menus(id) ON DELETE CASCADE,
    INDEX idx_parent_id (parent_id),
    INDEX idx_platform (platform),
    INDEX idx_sort (sort),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### role_menu_permissions 表
```sql
CREATE TABLE role_menu_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL COMMENT '角色ID',
    menu_id BIGINT NOT NULL COMMENT '菜单ID',
    permissions JSON NOT NULL COMMENT '权限列表:["read","create","update","delete","export","share"]',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE,
    UNIQUE KEY uk_role_menu (role_id, menu_id),
    INDEX idx_role_id (role_id),
    INDEX idx_menu_id (menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### AuthGuard 路由守卫
```typescript
// components/AuthGuard.tsx
import { Navigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

interface AuthGuardProps {
  children: React.ReactNode;
  requiredRoles?: string[];
  requiredPermissions?: string[];
}

export const AuthGuard: React.FC<AuthGuardProps> = ({
  children,
  requiredRoles = [],
  requiredPermissions = []
}) => {
  const { user, isAuthenticated, hasRole, hasPermission } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (requiredRoles.length > 0 && !requiredRoles.some(role => hasRole(role))) {
    return <Navigate to="/403" replace />;
  }

  if (requiredPermissions.length > 0 && !requiredPermissions.some(perm => hasPermission(perm))) {
    return <Navigate to="/403" replace />;
  }

  return <>{children}</>;
};
```

#### PermissionCheck 权限检查组件
```typescript
// components/PermissionCheck.tsx
import { useAuth } from '../hooks/useAuth';

interface PermissionCheckProps {
  children: React.ReactNode;
  roles?: string[];
  permissions?: string[];
  fallback?: React.ReactNode;
}

export const PermissionCheck: React.FC<PermissionCheckProps> = ({
  children,
  roles = [],
  permissions = [],
  fallback = null
}) => {
  const { hasRole, hasPermission } = useAuth();

  const hasRequiredRole = roles.length === 0 || roles.some(role => hasRole(role));
  const hasRequiredPermission = permissions.length === 0 || permissions.some(perm => hasPermission(perm));

  if (hasRequiredRole && hasRequiredPermission) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
};
```

### 认证相关Hook

#### useAuth Hook
```typescript
// hooks/useAuth.ts
import { useState, useEffect, createContext, useContext } from 'react';
import { authApi } from '../services/authApi';

interface User {
  id: number;
  username: string;
  roles: string[];
  permissions: string[];
  brandIds: number[];
}

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  hasRole: (role: string) => boolean;
  hasPermission: (permission: string) => boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const login = async (username: string, password: string) => {
    const response = await authApi.login(username, password);
    const { token, user: userData } = response.data;
    
    localStorage.setItem('dmh_token', token);
    localStorage.setItem('dmh_user', JSON.stringify(userData));
    
    setUser(userData);
    setIsAuthenticated(true);
  };

  const logout = () => {
    localStorage.removeItem('dmh_token');
    localStorage.removeItem('dmh_user');
    setUser(null);
    setIsAuthenticated(false);
  };

  const hasRole = (role: string): boolean => {
    return user?.roles.includes(role) || user?.roles.includes('platform_admin') || false;
  };

  const hasPermission = (permission: string): boolean => {
    if (user?.roles.includes('platform_admin')) return true;
    return user?.permissions.includes(permission) || false;
  };

  useEffect(() => {
    const token = localStorage.getItem('dmh_token');
    const userData = localStorage.getItem('dmh_user');
    
    if (token && userData) {
      setUser(JSON.parse(userData));
      setIsAuthenticated(true);
    }
  }, []);

  return (
    <AuthContext.Provider value={{
      user,
      isAuthenticated,
      login,
      logout,
      hasRole,
      hasPermission
    }}>
      {children}
    </AuthContext.Provider>
  );
};
```

---

## 🔧 后端实现

### JWT认证中间件

```go
// internal/middleware/authmiddleware.go
package middleware

import (
    "context"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
    secret string
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
    return &AuthMiddleware{secret: secret}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 获取Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }

        // 验证Bearer token格式
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
            return
        }

        // 解析JWT token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(m.secret), nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // 提取用户信息
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }

        // 将用户信息添加到context
        ctx := context.WithValue(r.Context(), "userId", claims["userId"])
        ctx = context.WithValue(ctx, "username", claims["username"])
        ctx = context.WithValue(ctx, "roles", claims["roles"])

        next.ServeHTTP(w, r.WithContext(ctx))
    }
}
```

### 权限检查中间件

```go
// internal/middleware/permissionmiddleware.go
package middleware

import (
    "context"
    "net/http"
    "strings"
)

type PermissionMiddleware struct {
    permissionService PermissionService
}

func NewPermissionMiddleware(permissionService PermissionService) *PermissionMiddleware {
    return &PermissionMiddleware{
        permissionService: permissionService,
    }
}

func (m *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 从context获取用户信息
        userId := r.Context().Value("userId")
        roles := r.Context().Value("roles")

        if userId == nil || roles == nil {
            http.Error(w, "User information not found", http.StatusUnauthorized)
            return
        }

        // 平台管理员拥有所有权限
        userRoles := roles.([]interface{})
        for _, role := range userRoles {
            if role.(string) == "platform_admin" {
                next.ServeHTTP(w, r)
                return
            }
        }

        // 提取资源和操作
        resource, action := m.extractResourceAction(r)
        
        // 检查权限
        hasPermission, err := m.permissionService.CheckPermission(userId, resource, action)
        if err != nil {
            http.Error(w, "Permission check failed", http.StatusInternalServerError)
            return
        }

        if !hasPermission {
            http.Error(w, "Insufficient permissions", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    }
}

func (m *PermissionMiddleware) extractResourceAction(r *http.Request) (string, string) {
    path := strings.TrimPrefix(r.URL.Path, "/api/v1/")
    parts := strings.Split(path, "/")
    
    resource := parts[0] // campaigns, orders, users等
    
    var action string
    switch r.Method {
    case "GET":
        action = "read"
    case "POST":
        action = "create"
    case "PUT", "PATCH":
        action = "update"
    case "DELETE":
        action = "delete"
    default:
        action = "read"
    }

    return resource, action
}
```

---

## ✅ 验收标准

### 功能验收
- [ ] 用户可以成功注册和登录
- [ ] JWT token正确生成和验证
- [ ] 4种用户角色权限正确区分
- [ ] API接口权限检查有效
- [ ] 数据级权限隔离正常工作
- [ ] 提现审核权限控制正确
- [ ] 前端权限控制组件正常工作
- [ ] 用户管理功能完整可用

### 安全验收
- [ ] 密码使用bcrypt加密存储
- [ ] JWT token包含必要的安全信息
- [ ] 权限检查无法绕过
- [ ] 数据访问严格按角色隔离
- [ ] 敏感操作有适当的权限控制
- [ ] 会话超时机制正常工作

### 性能验收
- [ ] 权限检查响应时间 < 50ms
- [ ] 支持1000+并发用户
- [ ] 权限缓存机制有效
- [ ] 数据库查询优化良好

---

## 🧪 测试用例

### 认证测试
1. 用户注册 - 正常流程
2. 用户注册 - 用户名重复
3. 用户注册 - 手机号重复
4. 用户登录 - 正确凭据
5. 用户登录 - 错误密码
6. 用户登录 - 用户不存在
7. JWT token验证 - 有效token
8. JWT token验证 - 过期token
9. JWT token验证 - 无效token

### 权限测试
1. 平台管理员 - 访问所有功能
2. 品牌管理员 - 只能访问本品牌数据
3. 活动参与者 - 只能访问个人功能
4. 匿名用户 - 只能访问公开功能
5. API权限检查 - 有权限访问
6. API权限检查 - 无权限访问
7. 数据隔离 - 品牌数据隔离
8. 提现审核 - 只有管理员可审核

### 安全测试
1. 密码加密 - bcrypt验证
2. SQL注入防护测试
3. XSS攻击防护测试
4. CSRF攻击防护测试
5. 权限提升攻击测试
6. 会话劫持防护测试

---

## 📝 开发清单

### 后端开发
- [ ] 完善JWT认证中间件
- [ ] 实现权限检查中间件
- [ ] 创建用户管理API
- [ ] 实现角色权限管理API
- [ ] 优化权限缓存机制
- [ ] 编写权限相关单元测试
- [ ] 实现操作审计日志
- [ ] 完善安全配置

### 前端开发
- [ ] 创建权限控制组件
- [ ] 实现认证相关Hook
- [ ] 优化登录注册页面
- [ ] 创建用户管理页面
- [ ] 实现角色权限配置页面
- [ ] 完善路由权限控制
- [ ] 优化错误处理和用户体验
- [ ] 编写前端权限测试

### 测试和文档
- [ ] 编写完整的测试用例
- [ ] 进行安全测试
- [ ] 性能测试和优化
- [ ] 更新API文档
- [ ] 编写部署指南
- [ ] 创建用户使用手册

---

## 🔗 相关文档
- [RBAC权限系统实现说明](../../../RBAC-PERMISSION-SYSTEM.md)
- [JWT认证机制](../../../backend/api/internal/middleware/authmiddleware.go)
- [用户管理API](../../../backend/api/internal/handler/auth/)
- [前端权限控制](../../../frontend-admin/components/)