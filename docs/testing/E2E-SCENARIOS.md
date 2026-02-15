# E2E 测试场景清单

> 最后更新: 2026-02-13

---

## 一、已有 E2E 测试

### 1.1 Frontend Admin

| 文件 | 场景 | 状态 |
|------|------|------|
| `e2e/admin-flows.spec.ts` | 登录流程 | ✅ |
| | - 管理员登录成功 | ✅ |
| | - 错误凭据显示错误 | ✅ |
| | 导航流程 | ✅ |
| | - 导航到活动监控 | ✅ |
| | - 导航到用户管理 | ✅ |
| | - 导航到分销监控 | ✅ |
| | - 导航到反馈管理 | ✅ |
| `e2e/admin-dashboard.spec.ts` | 仪表盘 | ✅ |
| `e2e/feedback.spec.ts` | 反馈流程 | ✅ |

### 1.2 Frontend H5

| 文件 | 场景 | 状态 |
|------|------|------|
| `e2e/h5-flows.spec.ts` | 活动流程 | ✅ |
| | - 查看活动列表 | ✅ |
| | - 导航到反馈 | ✅ |
| | 品牌管理流程 | ✅ |
| | - 访问品牌登录页 | ✅ |
| | 分销商流程 | ✅ |
| | - 访问分销商页 | ✅ |
| | 订单流程 | ✅ |
| | - 访问订单页 | ✅ |
| `e2e/feedback.spec.ts` | 反馈流程 | ✅ |

---

## 二、待补充 E2E 场景

### 2.1 Frontend Admin (优先级 P1)

#### 场景 1: 用户管理完整流程

```typescript
// e2e/user-management.spec.ts
describe('User Management Flow', () => {
  test('admin can create, edit, and delete user', async ({ page }) => {
    // 1. 登录
    // 2. 导航到用户管理
    // 3. 点击创建用户
    // 4. 填写表单并提交
    // 5. 验证用户出现在列表
    // 6. 编辑用户信息
    // 7. 验证更新成功
    // 8. 禁用用户
    // 9. 验证状态变更
  })
})
```

| 步骤 | 操作 | 验证点 |
|------|------|--------|
| 1 | 登录 admin | 进入仪表盘 |
| 2 | 点击用户管理 | 显示用户列表 |
| 3 | 点击创建用户 | 打开创建表单 |
| 4 | 填写用户信息并提交 | 成功提示，列表更新 |
| 5 | 搜索新用户 | 显示新创建用户 |
| 6 | 点击编辑 | 打开编辑表单 |
| 7 | 修改信息并保存 | 成功提示，信息更新 |
| 8 | 点击禁用 | 确认对话框 |
| 9 | 确认禁用 | 状态变为禁用 |

#### 场景 2: 品牌与活动创建流程

```typescript
// e2e/brand-campaign.spec.ts
describe('Brand and Campaign Flow', () => {
  test('admin can create brand and campaign', async ({ page }) => {
    // 1. 创建品牌
    // 2. 为品牌创建活动
    // 3. 配置活动页面
    // 4. 启用活动
    // 5. 验证活动状态
  })
})
```

| 步骤 | 操作 | 验证点 |
|------|------|--------|
| 1 | 创建品牌 | 品牌出现在列表 |
| 2 | 选择品牌创建活动 | 进入活动编辑 |
| 3 | 配置基本信息 | 表单验证通过 |
| 4 | 设计活动页面 | 组件渲染正确 |
| 5 | 保存并启用 | 活动状态为启用 |
| 6 | 查看活动详情 | 信息显示正确 |

#### 场景 3: 权限变更验证

```typescript
// e2e/permission-change.spec.ts
describe('Permission Change Flow', () => {
  test('role permission changes take effect immediately', async ({ page }) => {
    // 1. 创建测试角色
    // 2. 分配部分权限
    // 3. 创建用户并分配角色
    // 4. 以该用户登录
    // 5. 验证只能看到授权菜单
    // 6. 管理员修改权限
    // 7. 用户刷新后验证权限变更
  })
})
```

#### 场景 4: 成员数据管理

```typescript
// e2e/member-management.spec.ts
describe('Member Management Flow', () => {
  test('admin can manage member data', async ({ page }) => {
    // 1. 查看成员列表
    // 2. 搜索筛选成员
    // 3. 查看成员详情
    // 4. 导出成员数据
    // 5. 合并重复成员
  })
})
```

---

### 2.2 Frontend H5 (优先级 P1)

#### 场景 1: 活动浏览与报名完整流程

```typescript
// e2e/campaign-participation.spec.ts
describe('Campaign Participation Flow', () => {
  test('user can browse and join campaign', async ({ page }) => {
    // 1. 访问首页
    // 2. 浏览活动列表
    // 3. 点击活动查看详情
    // 4. 填写报名表单
    // 5. 提交报名
    // 6. 查看我的报名记录
  })
})
```

| 步骤 | 操作 | 验证点 |
|------|------|--------|
| 1 | 访问 / | 显示活动列表 |
| 2 | 滚动浏览 | 加载更多活动 |
| 3 | 点击活动卡片 | 进入详情页 |
| 4 | 查看活动信息 | 页面组件正确渲染 |
| 5 | 填写报名表单 | 表单验证 |
| 6 | 提交报名 | 成功提示 |
| 7 | 查看我的报名 | 显示刚才的报名 |

#### 场景 2: 品牌管理员操作流程

```typescript
// e2e/brand-admin.spec.ts
describe('Brand Admin Flow', () => {
  test('brand admin can manage campaigns', async ({ page }) => {
    // 1. 品牌管理员登录
    // 2. 查看活动列表
    // 3. 创建新活动
    // 4. 使用设计器编辑页面
    // 5. 查看报名数据
    // 6. 核销订单
  })
})
```

| 步骤 | 操作 | 验证点 |
|------|------|--------|
| 1 | 访问 /brand/login | 显示登录表单 |
| 2 | 输入凭据登录 | 进入品牌管理 |
| 3 | 点击创建活动 | 进入活动编辑 |
| 4 | 配置基本信息 | 表单验证 |
| 5 | 使用设计器 | 组件拖拽正常 |
| 6 | 保存活动 | 成功提示 |
| 7 | 查看报名 | 显示报名列表 |
| 8 | 扫码核销 | 核销成功 |

#### 场景 3: 分销商推广流程

```typescript
// e2e/distributor.spec.ts
describe('Distributor Flow', () => {
  test('user can become distributor and earn rewards', async ({ page }) => {
    // 1. 访问分销中心
    // 2. 申请成为分销商
    // 3. 获取推广链接
    // 4. 模拟推广注册
    // 5. 查看下级列表
    // 6. 查看佣金
    // 7. 申请提现
  })
})
```

| 步骤 | 操作 | 验证点 |
|------|------|--------|
| 1 | 访问 /distributor | 显示分销介绍 |
| 2 | 申请成为分销商 | 提交申请 |
| 3 | 审核通过后登录 | 进入分销中心 |
| 4 | 复制推广链接 | 链接复制成功 |
| 5 | 通过链接注册新用户 | 关联成功 |
| 6 | 查看我的下级 | 显示新用户 |
| 7 | 查看佣金明细 | 佣金计算正确 |
| 8 | 申请提现 | 提交成功 |

#### 场景 4: 订单核销流程

```typescript
// e2e/order-verification.spec.ts
describe('Order Verification Flow', () => {
  test('brand admin can verify orders', async ({ page }) => {
    // 1. 创建测试活动
    // 2. 用户报名
    // 3. 管理员扫码核销
    // 4. 验证核销状态
    // 5. 尝试重复核销
    // 6. 验证重复核销失败
  })
})
```

---

## 三、E2E 测试优先级矩阵

| 场景 | 影响 | 频率 | 复杂度 | 优先级 |
|------|------|------|--------|--------|
| 活动浏览报名 | 高 | 高 | 中 | P0 |
| 用户登录 | 高 | 高 | 低 | P0 ✅ |
| 品牌管理活动 | 高 | 中 | 高 | P1 |
| 订单核销 | 高 | 中 | 中 | P1 |
| 分销商推广 | 中 | 中 | 高 | P1 |
| 权限变更 | 高 | 低 | 中 | P2 |
| 成员管理 | 中 | 低 | 中 | P2 |
| 数据导出 | 低 | 低 | 低 | P3 |

---

## 四、测试环境配置

### 4.1 环境要求

```yaml
# 前置条件
- MySQL 8.0 运行中
- Redis 运行中（可选）
- 后端 API 运行在 localhost:8889
- 前端运行在 localhost:3000 (Admin) / localhost:3100 (H5)
- 测试账号: admin / 123456
```

### 4.2 测试数据

```sql
-- 测试账号（应存在于数据库）
-- 平台管理员
INSERT INTO users (username, password, status) VALUES ('admin', '$2a$10$...', 'active');
-- 品牌管理员
INSERT INTO users (username, password, status) VALUES ('brand_manager', '$2a$10$...', 'active');
-- 普通用户
INSERT INTO users (username, password, status) VALUES ('test_user', '$2a$10$...', 'active');
```

### 4.3 Playwright 配置

```typescript
// 已有配置
// frontend-admin/playwright.config.ts
// frontend-h5/playwright.config.ts

// 建议添加
export default defineConfig({
  // ...
  use: {
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    trace: 'retain-on-failure',
  },
})
```

---

## 五、E2E 测试最佳实践

### 5.1 选择器策略

```typescript
// ✅ 推荐: 使用 data-testid
await page.click('[data-testid="submit-button"]')

// ✅ 推荐: 使用语义化选择器
await page.click('button[type="submit"]')
await page.fill('input[name="username"]')

// ❌ 避免: 使用易变的类名
await page.click('.btn-primary-234')  // 可能因样式变更而失效
```

### 5.2 等待策略

```typescript
// ✅ 推荐: 等待特定状态
await expect(page.locator('.success-message')).toBeVisible()

// ❌ 避免: 固定等待时间
await page.waitForTimeout(2000)  // 仅在必要时使用
```

### 5.3 测试隔离

```typescript
// 每个测试应独立，不依赖其他测试的状态
test.beforeEach(async ({ page }) => {
  // 重置到已知状态
  await page.goto('/')
  // 清除 localStorage
  await page.evaluate(() => localStorage.clear())
})
```

---

## 六、更新日志

| 日期 | 变更内容 |
|------|---------|
| 2026-02-13 | 创建 E2E 场景清单文档 |
