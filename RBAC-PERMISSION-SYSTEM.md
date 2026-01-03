# DMH 权限管理系统 (RBAC) 实现说明

## 系统概览

本系统实现了基于角色的访问控制(RBAC - Role-Based Access Control)，支持4种用户角色，每种角色拥有不同的权限范围。

## 1. 角色定义

### 1.1 平台管理员 (platform_admin)
- **权限范围**: 系统最高权限
- **功能权限**:
  - 管理所有品牌和活动
  - 管理所有用户和角色
  - 审核提现申请
  - 查看所有统计数据
  - 系统配置管理

### 1.2 品牌管理员 (brand_admin)
- **权限范围**: 品牌级权限
- **功能权限**:
  - 管理本品牌的活动(创建、编辑、删除、发布)
  - 查看本品牌的订单和数据
  - 查看本品牌的用户和统计
  - 管理本品牌的素材库

### 1.3 活动参与者 (participant)
- **权限范围**: 用户级权限
- **功能权限**:
  - 参加活动(报名、支付)
  - 分享活动获取推广链接
  - 查看个人余额和奖励记录
  - 申请提现
  - 查看个人参与的活动历史

### 1.4 匿名用户 (anonymous)
- **权限范围**: 最小权限
- **功能权限**:
  - 浏览活动列表
  - 查看活动详情
  - 注册账号

## 2. 数据库设计

### 2.1 核心表结构

```sql
-- 用户表
users (id, username, password, phone, email, avatar, real_name, status, created_at, updated_at)

-- 角色表
roles (id, name, code, description, created_at, updated_at)

-- 用户角色关联表
user_roles (id, user_id, role_id, created_at)

-- 权限表
permissions (id, name, code, resource, action, description, created_at)

-- 角色权限关联表
role_permissions (id, role_id, permission_id, created_at)

-- 品牌表
brands (id, name, logo, description, status, created_at, updated_at)

-- 品牌管理员关联表
brand_admins (id, brand_id, user_id, created_at)

-- 提现申请表
withdrawals (id, user_id, amount, bank_name, bank_account, account_name, 
            status, remark, approved_by, approved_at, created_at, updated_at)
```

### 2.2 权限编码规范

权限code采用 `资源:操作` 格式：
- `campaign:create` - 创建活动
- `campaign:read` - 查看活动
- `campaign:update` - 编辑活动
- `campaign:delete` - 删除活动
- `order:read` - 查看订单
- `withdrawal:create` - 申请提现
- `withdrawal:approve` - 审核提现
- `user:create` - 创建用户
- `brand:create` - 创建品牌
- `stats:read` - 查看统计

## 3. API 接口设计

### 3.1 认证相关接口 (无需token)

```
POST /api/v1/auth/register  - 用户注册
POST /api/v1/auth/login     - 用户登录
```

### 3.2 用户信息接口 (需要token)

```
GET /api/v1/auth/userinfo   - 获取当前用户信息
```

### 3.3 品牌管理接口 (需要token + 权限)

```
POST   /api/v1/brands       - 创建品牌 (仅平台管理员)
GET    /api/v1/brands       - 获取品牌列表
GET    /api/v1/brands/:id   - 获取品牌详情
```

### 3.4 活动管理接口 (需要token + 权限)

```
POST   /api/v1/campaigns    - 创建活动 (平台管理员 + 品牌管理员)
GET    /api/v1/campaigns    - 获取活动列表
GET    /api/v1/campaigns/:id - 获取活动详情
```

### 3.5 提现管理接口 (需要token + 权限)

```
POST   /api/v1/withdrawals              - 申请提现 (所有用户)
GET    /api/v1/withdrawals              - 获取提现列表
GET    /api/v1/withdrawals/:id          - 获取提现详情
POST   /api/v1/withdrawals/:id/approve  - 审核提现 (仅平台管理员)
```

## 4. 后端实现

### 4.1 核心文件结构

```
backend/
├── model/
│   ├── user.go                    # 用户、角色、权限模型
│   └── campaign.go                # 活动模型
├── api/
│   ├── dmh.api                    # API定义
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go          # 配置(包含JWT)
│   │   ├── middleware/
│   │   │   └── authmiddleware.go  # 权限检查中间件
│   │   ├── logic/
│   │   │   ├── auth/
│   │   │   │   ├── loginlogic.go      # 登录逻辑
│   │   │   │   ├── registerlogic.go   # 注册逻辑
│   │   │   │   └── getuserinfologic.go # 获取用户信息
│   │   │   ├── brand/
│   │   │   │   └── ...                # 品牌管理逻辑
│   │   │   ├── campaign/
│   │   │   │   └── ...                # 活动管理逻辑
│   │   │   └── withdrawal/
│   │   │       ├── applywithdrawallogic.go   # 申请提现
│   │   │       └── approvewithdrawallogic.go # 审核提现
│   │   └── handler/
│   │       ├── auth/              # 认证handler
│   │       ├── brand/             # 品牌handler
│   │       ├── campaign/          # 活动handler
│   │       └── withdrawal/        # 提现handler
├── common/
│   └── utils/
│       └── crypto.go              # 密码加密工具
└── scripts/
    └── init.sql                   # 数据库初始化脚本
```

### 4.2 JWT Token结构

```json
{
  "exp": 1735689600,        // 过期时间
  "iat": 1735603200,        // 签发时间
  "userId": 1,              // 用户ID
  "username": "admin",      // 用户名
  "roles": ["platform_admin"] // 角色列表
}
```

### 4.3 登录流程

```
1. 用户提交 username + password
2. 查询数据库验证用户存在
3. 使用 bcrypt 验证密码
4. 检查用户状态 (active/disabled)
5. 查询用户角色列表
6. 如果是品牌管理员，查询管理的品牌ID列表
7. 生成JWT Token (包含userId、roles等信息)
8. 返回 token + 用户信息
```

### 4.4 注册流程

```
1. 校验用户名和手机号唯一性
2. 使用 bcrypt 加密密码
3. 开启数据库事务:
   3.1 创建用户记录
   3.2 分配默认角色 (participant)
   3.3 初始化用户余额记录
4. 生成JWT Token
5. 返回 token + 用户信息
```

### 4.5 权限检查流程

```
1. 请求到达 -> JWT中间件验证token
2. 从token中提取 userId 和 roles
3. 进入 AuthMiddleware:
   3.1 从URL提取 resource (如 campaigns、orders)
   3.2 从HTTP method映射 action (GET->read, POST->create等)
   3.3 检查用户是否有该权限:
       - platform_admin: 直接放行(拥有所有权限)
       - 其他角色: 查询数据库检查是否有对应权限
4. 权限通过 -> 执行业务逻辑
5. 权限不通过 -> 返回403 Forbidden
```

### 4.6 提现审核流程

```
1. 用户申请提现:
   - 检查余额是否足够
   - 创建提现申请记录(status=pending)
   - 不立即扣除余额

2. 平台管理员审核:
   - 同意(approved):
     * 开启事务
     * 扣除用户余额
     * 更新提现记录状态
     * 记录审核人和时间
   
   - 拒绝(rejected):
     * 更新提现记录状态
     * 填写拒绝原因
     * 不扣除余额
```

## 5. 前端实现

### 5.1 核心文件

```
frontend-admin/
├── types.ts                      # TypeScript类型定义
├── services/
│   └── auth.ts                   # 认证API服务
├── views/
│   ├── LoginView.tsx             # 登录页面
│   ├── RegisterView.tsx          # 注册页面
│   ├── WithdrawalManageView.tsx  # 提现管理页面
│   └── BrandManageView.tsx       # 品牌管理页面
├── components/
│   ├── AuthGuard.tsx             # 权限守卫组件
│   └── PermissionCheck.tsx       # 权限检查组件
└── utils/
    └── auth.ts                   # 认证工具(token存储)
```

### 5.2 Token存储

```typescript
// localStorage 存储
localStorage.setItem('dmh_token', token);
localStorage.setItem('dmh_user', JSON.stringify(userInfo));

// 请求时自动携带
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('dmh_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

### 5.3 路由守卫

```typescript
// 检查是否登录
const AuthGuard = ({ children }) => {
  const token = localStorage.getItem('dmh_token');
  if (!token) {
    return <Navigate to="/login" />;
  }
  return children;
};

// 检查角色权限
const RoleGuard = ({ roles, children }) => {
  const user = JSON.parse(localStorage.getItem('dmh_user') || '{}');
  const hasRole = roles.some(role => user.roles?.includes(role));
  if (!hasRole) {
    return <Navigate to="/403" />;
  }
  return children;
};
```

### 5.4 根据角色显示菜单

```typescript
const menuItems = [
  {
    label: '活动管理',
    path: '/campaigns',
    roles: ['platform_admin', 'brand_admin'], // 只有这些角色能看到
  },
  {
    label: '品牌管理',
    path: '/brands',
    roles: ['platform_admin'], // 只有平台管理员能看到
  },
  {
    label: '提现审核',
    path: '/withdrawals',
    roles: ['platform_admin'],
  },
  {
    label: '我的奖励',
    path: '/rewards',
    roles: ['participant'], // 普通用户能看到
  },
];

// 根据当前用户角色过滤菜单
const visibleMenus = menuItems.filter(item => 
  item.roles.some(role => currentUser.roles.includes(role))
);
```

## 6. 测试账号

系统初始化时会创建以下测试账号（密码都是 `123456`）：

| 用户名 | 角色 | 说明 |
|-------|------|------|
| admin | platform_admin | 平台管理员，拥有所有权限 |
| brand_manager | brand_admin | 品牌管理员，管理"品牌A" |
| user001 | participant | 普通用户，可参与活动 |

## 7. 安全措施

### 7.1 密码安全
- 使用 bcrypt 加密存储，不存储明文密码
- bcrypt cost = 10，足够安全且性能可接受

### 7.2 Token安全
- JWT签名防止篡改
- Token有效期24小时（可配置）
- 敏感操作需要重新验证

### 7.3 API安全
- 所有需要权限的接口都需要JWT token
- 权限检查在中间件层统一处理
- SQL使用参数化查询，防止SQL注入

### 7.4 提现安全
- 提现需要审核，不能自动到账
- 使用数据库事务保证一致性
- 记录审核人和审核时间，可追溯

## 8. 扩展建议

### 8.1 短期优化
1. 添加验证码登录(手机验证码)
2. 添加找回密码功能
3. 实现更细粒度的权限控制（数据级权限）
4. 添加操作日志记录

### 8.2 中期优化
1. 支持多品牌管理员(一个用户管理多个品牌)
2. 添加权限缓存(Redis)，减少数据库查询
3. 实现Token刷新机制
4. 添加二次认证(2FA)

### 8.3 长期优化
1. 实现动态权限配置系统
2. 支持自定义角色和权限
3. 添加审计日志和合规报告
4. 实现单点登录(SSO)

## 9. 常见问题

### Q1: 如何添加新的权限？
A: 在 `init.sql` 中的 `permissions` 表插入新权限记录，然后在 `role_permissions` 表中关联到对应角色。

### Q2: 如何给用户分配品牌？
A: 在 `brand_admins` 表中插入记录，关联 `brand_id` 和 `user_id`。

### Q3: Token过期后怎么办？
A: 前端收到401响应后，清除本地token并跳转到登录页。后续可实现refresh token机制。

### Q4: 如何实现数据隔离(品牌管理员只看本品牌数据)？
A: 在业务逻辑中根据用户角色和brandIds过滤数据：
```go
// 如果是品牌管理员，只查询其管理的品牌的活动
if isBrandAdmin {
    query = db.Where("brand_id IN ?", userBrandIds)
}
```

### Q5: 提现审核后如何实际打款？
A: 当前版本只是修改状态，实际打款需要接入支付服务商的API（如微信支付企业付款到零钱）。

## 10. 数据库迁移

如果数据库已存在，运行以下SQL添加权限相关表：

```bash
# 备份现有数据库
mysqldump -u root -p dmh > dmh_backup.sql

# 执行新的 init.sql
mysql -u root -p dmh < backend/scripts/init.sql
```

注意：新的 `init.sql` 会添加 `brand_id` 字段到 `campaigns` 表，需要更新现有数据。

## 11. API调用示例

### 11.1 登录

```bash
curl -X POST http://localhost:8888/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }'
```

### 11.2 获取活动列表（需要token）

```bash
curl -X GET http://localhost:8888/api/v1/campaigns \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 11.3 申请提现

```bash
curl -X POST http://localhost:8888/api/v1/withdrawals \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.00,
    "bankName": "中国银行",
    "bankAccount": "6222000000000000",
    "accountName": "张三"
  }'
```

### 11.4 审核提现（仅平台管理员）

```bash
curl -X POST http://localhost:8888/api/v1/withdrawals/1/approve \
  -H "Authorization: Bearer ADMIN_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved",
    "remark": "审核通过"
  }'
```

## 12. 总结

本权限管理系统实现了：
✅ 完整的RBAC权限模型
✅ 4种用户角色(平台管理员、品牌管理员、活动参与者、匿名用户)
✅ JWT认证机制
✅ 密码bcrypt加密
✅ 提现申请和审核流程
✅ 品牌多租户支持
✅ 数据级权限隔离

系统具有良好的可扩展性，可根据业务需求继续完善和优化。
