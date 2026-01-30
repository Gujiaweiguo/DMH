# 📡 DMH API 文档

## 目录

- [API概述](#api概述)
- [认证授权](#认证授权)
- [用户管理](#用户管理)
- [活动管理](#活动管理)
- [订单管理](#订单管理)
- [奖励管理](#奖励管理)
- [分销商管理](#分销商管理)
- [会员管理](#会员管理)
- [品牌管理](#品牌管理)
- [权限管理](#权限管理)
- [安全管理](#安全管理)

---

## API概述

### 基础信息

- **Base URL**: `https://api.dmh.com/api/v1`
- **协议**: HTTPS
- **数据格式**: JSON
- **字符编码**: UTF-8

### 通用请求头

```http
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
Accept: application/json
```

### 通用响应格式

**成功响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    // 业务数据
  }
}
```

**错误响应**:
```json
{
  "code": 40001,
  "message": "错误描述",
  "data": null
}
```

### 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 40001 | 业务错误（如重复报名） |
| 500 | 服务器内部错误 |

---

## 认证授权

### 用户注册

**接口**: `POST /auth/register`

**请求参数**:
```json
{
  "username": "string",
  "password": "string",
  "phone": "string",
  "email": "string (可选)",
  "realName": "string (可选)"
}
```

**响应**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "userId": 1,
  "username": "testuser",
  "phone": "13800138000",
  "roles": ["participant"]
}
```

### 用户登录

**接口**: `POST /auth/login`

**请求参数**:
```json
{
  "username": "string",
  "password": "string"
}
```

**响应**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "userId": 1,
  "username": "testuser",
  "phone": "13800138000",
  "realName": "测试用户",
  "roles": ["platform_admin"],
  "brandIds": [1, 2]
}
```

### 刷新Token

**接口**: `POST /auth/refresh-token`

**请求参数**:
```json
{
  "token": "string"
}
```

**响应**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 获取当前用户信息

**接口**: `GET /auth/userinfo`

**响应**:
```json
{
  "id": 1,
  "username": "testuser",
  "phone": "13800138000",
  "email": "test@example.com",
  "realName": "测试用户",
  "avatar": "https://...",
  "status": "active",
  "roles": ["platform_admin"],
  "brandIds": [1, 2],
  "createdAt": "2025-01-01T00:00:00Z"
}
```

---

## 用户管理

### 创建用户（管理员）

**接口**: `POST /admin/users`

**权限**: platform_admin

**请求参数**:
```json
{
  "username": "string",
  "password": "string",
  "phone": "string",
  "email": "string (可选)",
  "realName": "string (可选)",
  "role": "platform_admin | brand_admin | participant",
  "brandIds": [1, 2]
}
```

### 获取用户列表

**接口**: `GET /admin/users`

**查询参数**:
- `page`: 页码（默认1）
- `pageSize`: 每页数量（默认20）
- `role`: 角色筛选
- `status`: 状态筛选
- `keyword`: 关键词搜索

**响应**:
```json
{
  "total": 100,
  "users": [
    {
      "id": 1,
      "username": "testuser",
      "phone": "13800138000",
      "realName": "测试用户",
      "roles": ["platform_admin"],
      "status": "active",
      "createdAt": "2025-01-01T00:00:00Z"
    }
  ]
}
```

---

## 活动管理

### 创建活动

**接口**: `POST /campaigns`

**权限**: brand_admin, platform_admin


**请求参数**:
```json
{
  "brandId": 1,
  "name": "新年促销活动",
  "description": "新年大促，推荐有礼",
  "formFields": [
    {
      "type": "text",
      "name": "name",
      "label": "姓名",
      "required": true,
      "placeholder": "请输入姓名"
    },
    {
      "type": "phone",
      "name": "phone",
      "label": "手机号",
      "required": true
    },
    {
      "type": "select",
      "name": "course",
      "label": "意向课程",
      "required": true,
      "options": ["前端开发", "后端开发", "全栈开发"]
    }
  ],
  "rewardRule": 10.00,
  "startTime": "2025-01-01T00:00:00Z",
  "endTime": "2025-12-31T23:59:59Z"
}
```

**响应**:
```json
{
  "id": 1,
  "brandId": 1,
  "brandName": "测试品牌",
  "name": "新年促销活动",
  "description": "新年大促，推荐有礼",
  "formFields": [...],
  "rewardRule": 10.00,
  "startTime": "2025-01-01T00:00:00Z",
  "endTime": "2025-12-31T23:59:59Z",
  "status": "active",
  "createdAt": "2025-01-01T10:00:00Z"
}
```

### 获取活动列表

**接口**: `GET /campaigns`

**查询参数**:
- `page`: 页码
- `pageSize`: 每页数量
- `status`: 状态筛选（active/paused/ended）
- `keyword`: 关键词搜索

**响应**:
```json
{
  "total": 100,
  "campaigns": [
    {
      "id": 1,
      "brandId": 1,
      "brandName": "测试品牌",
      "name": "新年促销活动",
      "status": "active",
      "startTime": "2025-01-01T00:00:00Z",
      "endTime": "2025-12-31T23:59:59Z",
      "createdAt": "2025-01-01T10:00:00Z"
    }
  ]
}
```

### 获取活动详情

**接口**: `GET /campaigns/:id`

**响应**: 同创建活动响应

### 保存页面配置

**接口**: `POST /campaign/page-config/:id`

**请求参数**:
```json
{
  "components": [
    {
      "type": "title",
      "data": {
        "title": "活动标题",
        "subtitle": "副标题"
      }
    },
    {
      "type": "poster",
      "data": {
        "imageUrl": "https://..."
      }
    }
  ],
  "theme": {
    "primaryColor": "#667eea",
    "backgroundColor": "#ffffff"
  }
}
```

---

## 订单管理

### 创建订单

**接口**: `POST /orders`

**请求参数**:
```json
{
  "campaignId": 1,
  "phone": "13800138000",
  "formData": {
    "name": "张三",
    "phone": "13800138000",
    "course": "前端开发"
  },
  "referrerId": 100
}
```

**响应**:
```json
{
  "id": 12345,
  "campaignId": 1,
  "phone": "13800138000",
  "formData": {...},
  "referrerId": 100,
  "status": "pending",
  "amount": 99.00,
  "createdAt": "2025-01-01T10:00:00Z"
}
```

**错误响应**:
```json
{
  "code": 40001,
  "message": "您已参与过该活动"
}
```

### 获取订单详情

**接口**: `GET /orders/:id`

**响应**:
```json
{
  "id": 12345,
  "campaignId": 1,
  "campaignName": "新年促销活动",
  "phone": "13800138000",
  "formData": {...},
  "referrerId": 100,
  "status": "paid",
  "amount": 99.00,
  "payStatus": "paid",
  "tradeNo": "4200001234567890",
  "syncStatus": "synced",
  "createdAt": "2025-01-01T10:00:00Z",
  "paidAt": "2025-01-01T10:05:00Z"
}
```

### 支付回调

**接口**: `POST /orders/payment/callback`

**请求参数**:
```json
{
  "orderId": 12345,
  "payStatus": "paid",
  "amount": 99.00,
  "tradeNo": "4200001234567890",
  "signature": "..."
}
```

**响应**:
```json
{
  "code": "SUCCESS",
  "message": "OK"
}
```

---

## 奖励管理

### 查询用户余额

**接口**: `GET /rewards/balance/:userId`

**响应**:
```json
{
  "userId": 100,
  "balance": 156.50,
  "totalReward": 200.00,
  "updatedAt": "2025-01-01T15:30:00Z"
}
```

### 查询奖励列表

**接口**: `GET /rewards/:userId`

**查询参数**:
- `page`: 页码
- `pageSize`: 每页数量

**响应**:
```json
{
  "total": 10,
  "rewards": [
    {
      "id": 1,
      "userId": 100,
      "orderId": 12345,
      "campaignId": 1,
      "campaignName": "新年促销活动",
      "amount": 10.00,
      "status": "settled",
      "settledAt": "2025-01-01T10:05:02Z",
      "createdAt": "2025-01-01T10:05:00Z"
    }
  ]
}
```

---

## 分销商管理

### 申请成为分销商

**接口**: `POST /distributor/apply`

**请求参数**:
```json
{
  "brandId": 1,
  "reason": "我想成为分销商"
}
```

**响应**:
```json
{
  "id": 1,
  "userId": 100,
  "brandId": 1,
  "status": "pending",
  "reason": "我想成为分销商",
  "createdAt": "2025-01-01T10:00:00Z"
}
```

### 生成推广链接

**接口**: `POST /distributor/link/generate`

**请求参数**:
```json
{
  "campaignId": 1
}
```

**响应**:
```json
{
  "linkId": 1,
  "link": "https://h5.dmh.com/campaign/1?code=ABC123",
  "linkCode": "ABC123",
  "qrcodeUrl": "https://...",
  "campaignId": 1
}
```

### 查询分销统计

**接口**: `GET /distributor/statistics/:brandId`

**响应**:
```json
{
  "distributorId": 1,
  "totalOrders": 100,
  "totalEarnings": 1000.00,
  "todayEarnings": 50.00,
  "monthEarnings": 500.00,
  "subordinatesCount": 10,
  "clickCount": 500,
  "conversionRate": 0.20
}
```

### 查询下级分销商

**接口**: `GET /distributor/subordinates/:brandId`

**响应**:
```json
{
  "total": 10,
  "subordinates": [
    {
      "id": 2,
      "userId": 101,
      "username": "user101",
      "level": 2,
      "totalOrders": 20,
      "totalEarnings": 200.00,
      "createdAt": "2025-01-01T10:00:00Z"
    }
  ]
}
```

### 审批分销商申请（品牌管理员）

**接口**: `POST /brands/:brandId/distributor/approve/:id`

**权限**: brand_admin, platform_admin

**请求参数**:
```json
{
  "action": "approve",
  "level": 1,
  "reason": "审批通过"
}
```

---

## 会员管理

### 获取会员列表

**接口**: `GET /members`

**查询参数**:
- `page`: 页码
- `pageSize`: 每页数量
- `brandId`: 品牌筛选
- `status`: 状态筛选
- `keyword`: 关键词搜索

**响应**:
```json
{
  "total": 1000,
  "members": [
    {
      "id": 1,
      "unionid": "oxxxxxx",
      "phone": "13800138000",
      "nickname": "张三",
      "avatar": "https://...",
      "status": "active",
      "totalOrders": 5,
      "totalPayment": 500.00,
      "totalReward": 50.00,
      "createdAt": "2025-01-01T00:00:00Z"
    }
  ]
}
```

### 获取会员详情

**接口**: `GET /members/:id`

**响应**:
```json
{
  "id": 1,
  "unionid": "oxxxxxx",
  "phone": "13800138000",
  "nickname": "张三",
  "avatar": "https://...",
  "gender": "male",
  "status": "active",
  "firstSource": "wechat",
  "totalOrders": 5,
  "totalPayment": 500.00,
  "totalReward": 50.00,
  "participatedCampaigns": 3,
  "brands": [
    {
      "brandId": 1,
      "brandName": "测试品牌",
      "firstCampaignId": 1
    }
  ],
  "tags": ["VIP", "活跃用户"],
  "createdAt": "2025-01-01T00:00:00Z"
}
```

### 合并会员

**接口**: `POST /members/merge`

**权限**: platform_admin

**请求参数**:
```json
{
  "sourceMemberId": 2,
  "targetMemberId": 1,
  "reason": "重复会员合并"
}
```

### 导出会员

**接口**: `POST /members/export`

**权限**: brand_admin, platform_admin

**请求参数**:
```json
{
  "brandId": 1,
  "startDate": "2025-01-01",
  "endDate": "2025-01-31",
  "status": "active"
}
```

**响应**:
```json
{
  "requestId": "export_123",
  "status": "pending",
  "message": "导出请求已提交，等待审批"
}
```

---

## 品牌管理

### 创建品牌

**接口**: `POST /brands`

**权限**: platform_admin

**请求参数**:
```json
{
  "name": "测试品牌",
  "logo": "https://...",
  "description": "品牌描述"
}
```

### 获取品牌列表

**接口**: `GET /brands`

**响应**:
```json
{
  "total": 10,
  "brands": [
    {
      "id": 1,
      "name": "测试品牌",
      "logo": "https://...",
      "description": "品牌描述",
      "status": "active",
      "createdAt": "2025-01-01T00:00:00Z"
    }
  ]
}
```

### 获取品牌统计

**接口**: `GET /brands/:id/stats`

**响应**:
```json
{
  "brandId": 1,
  "totalCampaigns": 10,
  "activeCampaigns": 5,
  "totalOrders": 1000,
  "totalRevenue": 100000.00,
  "totalRewards": 10000.00,
  "participantCount": 500,
  "conversionRate": 0.50,
  "lastUpdated": "2025-01-01T00:00:00Z"
}
```

---

## 权限管理

### 获取角色列表

**接口**: `GET /roles`

**响应**:
```json
[
  {
    "id": 1,
    "name": "平台管理员",
    "code": "platform_admin",
    "description": "平台最高权限",
    "permissions": ["user:create", "user:update", "user:delete"],
    "createdAt": "2025-01-01T00:00:00Z"
  }
]
```

### 获取权限列表

**接口**: `GET /permissions`

**响应**:
```json
[
  {
    "id": 1,
    "name": "创建用户",
    "code": "user:create",
    "resource": "user",
    "action": "create",
    "description": "创建用户权限"
  }
]
```

### 配置角色权限

**接口**: `POST /roles/permissions`

**权限**: platform_admin

**请求参数**:
```json
{
  "roleId": 1,
  "permissionIds": [1, 2, 3, 4, 5]
}
```

### 获取用户菜单

**接口**: `GET /users/menus`

**查询参数**:
- `platform`: admin | h5

**响应**:
```json
{
  "userId": 1,
  "platform": "admin",
  "menus": [
    {
      "id": 1,
      "name": "用户管理",
      "code": "user_management",
      "path": "/users",
      "icon": "users",
      "sort": 1,
      "type": "menu",
      "children": [
        {
          "id": 2,
          "name": "用户列表",
          "code": "user_list",
          "path": "/users/list",
          "type": "menu"
        }
      ]
    }
  ]
}
```

---

## 安全管理

### 获取密码策略

**接口**: `GET /security/password-policy`

**响应**:
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
  "sessionTimeout": 1440,
  "maxConcurrentSessions": 3
}
```

### 获取审计日志

**接口**: `GET /security/audit-logs`

**查询参数**:
- `page`: 页码
- `pageSize`: 每页数量
- `userId`: 用户筛选
- `action`: 操作筛选
- `startDate`: 开始日期
- `endDate`: 结束日期

**响应**:
```json
{
  "total": 1000,
  "logs": [
    {
      "id": 1,
      "userId": 1,
      "username": "admin",
      "action": "user:create",
      "resource": "user",
      "resourceId": "123",
      "details": "创建用户 testuser",
      "clientIp": "192.168.1.1",
      "userAgent": "Mozilla/5.0...",
      "status": "success",
      "createdAt": "2025-01-01T10:00:00Z"
    }
  ]
}
```

### 获取用户会话

**接口**: `GET /security/sessions`

**响应**:
```json
{
  "total": 5,
  "sessions": [
    {
      "id": "session_123",
      "userId": 1,
      "clientIp": "192.168.1.1",
      "userAgent": "Mozilla/5.0...",
      "loginAt": "2025-01-01T10:00:00Z",
      "lastActiveAt": "2025-01-01T11:00:00Z",
      "expiresAt": "2025-01-02T10:00:00Z",
      "status": "active"
    }
  ]
}
```

### 强制下线用户

**接口**: `POST /security/force-logout/:userId`

**权限**: platform_admin

**请求参数**:
```json
{
  "reason": "安全原因强制下线"
}
```

---

## 提现管理

### 申请提现

**接口**: `POST /withdrawals`

**请求参数**:
```json
{
  "amount": 100.00,
  "bankName": "中国银行",
  "bankAccount": "6222xxxxxxxx1234",
  "accountName": "张三"
}
```

### 获取提现列表

**接口**: `GET /withdrawals`

**查询参数**:
- `page`: 页码
- `pageSize`: 每页数量
- `status`: 状态筛选（pending/approved/rejected）

**响应**:
```json
{
  "total": 10,
  "withdrawals": [
    {
      "id": 1,
      "userId": 100,
      "username": "testuser",
      "amount": 100.00,
      "bankName": "中国银行",
      "bankAccount": "6222xxxxxxxx1234",
      "accountName": "张三",
      "status": "pending",
      "createdAt": "2025-01-01T10:00:00Z"
    }
  ]
}
```

### 审批提现

**接口**: `POST /withdrawals/:id/approve`

**权限**: platform_admin

**请求参数**:
```json
{
  "status": "approved",
  "remark": "审批通过"
}
```

---

## 数据同步

### 查询同步状态

**接口**: `GET /sync/status/:orderId`

**响应**:
```json
{
  "orderId": 12345,
  "syncStatus": "synced",
  "syncTime": "2025-01-01T10:05:30Z",
  "attempts": 1,
  "errorMsg": null
}
```

### 手动重试同步

**接口**: `POST /sync/retry/:orderId`

**权限**: platform_admin

**响应**:
```json
{
  "orderId": 12345,
  "message": "同步任务已加入队列",
  "taskId": "task_abc123"
}
```

### 同步统计

**接口**: `GET /sync/statistics`

**查询参数**:
- `startDate`: 开始日期
- `endDate`: 结束日期

**响应**:
```json
{
  "totalSyncs": 1000,
  "successSyncs": 980,
  "failedSyncs": 20,
  "successRate": 0.98,
  "avgTime": "1.2s"
}
```

### 健康检查

**接口**: `GET /sync/health`

**响应**:
```json
{
  "status": "healthy",
  "database": {
    "connected": true,
    "type": "oracle",
    "host": "external-db.example.com",
    "latency": "15ms"
  },
  "queue": {
    "pending": 5,
    "processing": 2
  }
}
```

---

## 附录

### 状态枚举

**订单状态 (order.status)**:
- `pending` - 待支付
- `paid` - 已支付
- `cancelled` - 已取消
- `refunded` - 已退款

**支付状态 (order.pay_status)**:
- `unpaid` - 未支付
- `paid` - 已支付
- `refunded` - 已退款

**同步状态 (order.sync_status)**:
- `pending` - 待同步
- `syncing` - 同步中
- `synced` - 已同步
- `failed` - 同步失败

**活动状态 (campaign.status)**:
- `active` - 进行中
- `paused` - 已暂停
- `ended` - 已结束

**用户状态 (user.status)**:
- `active` - 正常
- `disabled` - 已禁用
- `locked` - 已锁定

**分销商状态 (distributor.status)**:
- `pending` - 待审批
- `active` - 正常
- `suspended` - 已暂停
- `rejected` - 已拒绝

---

## 相关文档

- [README.md](./README.md) - 项目介绍
- [SETUP.md](./SETUP.md) - 环境搭建指南
- [ARCHITECTURE.md](./ARCHITECTURE.md) - 系统架构
- [DEVELOPMENT.md](./DEVELOPMENT.md) - 开发指南
- [后端API定义](./backend/api/dmh.api) - go-zero API定义文件

---

**文档版本**: v1.0  
**最后更新**: 2025-01-21  
**维护者**: DMH Team
