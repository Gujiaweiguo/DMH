# 前端测试策略文档

> 文档版本: 1.0
> 创建日期: 2026-02-14
> 适用项目: DMH 数字营销中台

---

## 一、概述

### 1.1 目的

本文档定义 DMH 项目前端（Admin 管理后台 + H5 移动端）的测试策略、规范和最佳实践。

### 1.2 范围

| 端 | 技术栈 | 测试框架 |
|---|--------|---------|
| Admin | Vue 3 + Vite + TypeScript | Vitest + Playwright |
| H5 | Vue 3 + Vite + Vant | Vitest + Playwright |

### 1.3 测试层级

```
┌─────────────────────────────────────────┐
│           E2E 测试 (Playwright)          │  ← 用户流程验证
├─────────────────────────────────────────┤
│        组件测试 (Vitest + Vue Test Utils) │  ← 组件行为验证
├─────────────────────────────────────────┤
│        单元测试 (Vitest)                  │  ← 纯函数/Logic 验证
└─────────────────────────────────────────┘
```

---

## 二、测试分层策略

### 2.1 单元测试 (Unit Tests)

**定义**: 测试独立的函数、工具类、Logic 层

**范围**:
- `src/utils/` - 工具函数
- `src/logic/` - 业务逻辑函数
- `src/composables/` - 组合式函数
- `src/services/` - API 封装

**示例**:

```typescript
// tests/unit/string.logic.test.js
import { describe, it, expect } from 'vitest'
import { truncate, capitalize } from '@/logic/string.logic'

describe('string.logic', () => {
  describe('truncate', () => {
    it('should truncate long strings', () => {
      expect(truncate('hello world', 5)).toBe('he...')
    })
    
    it('should return original string if shorter than limit', () => {
      expect(truncate('hi', 5)).toBe('hi')
    })
  })
})
```

**覆盖率目标**: 80%+

### 2.2 组件测试 (Component Tests)

**定义**: 测试 Vue 组件的渲染、交互、状态

**范围**:
- `src/views/` - 页面组件
- `src/components/` - 通用组件
- `src/layouts/` - 布局组件

**示例**:

```typescript
// tests/components/LoginView.test.ts
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import LoginView from '@/views/LoginView.vue'

describe('LoginView', () => {
  it('should render login form', () => {
    const wrapper = mount(LoginView)
    expect(wrapper.find('input[type="text"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('should show error on invalid credentials', async () => {
    const wrapper = mount(LoginView)
    await wrapper.find('form').trigger('submit')
    expect(wrapper.text()).toContain('请输入用户名')
  })
})
```

**覆盖率目标**: 60%+

### 2.3 E2E 测试 (End-to-End Tests)

**定义**: 模拟真实用户操作流程

**范围**:
- 用户登录/登出
- 核心业务流程（活动创建、报名、核销）
- 关键用户路径

**示例**:

```typescript
// e2e/admin-flows.spec.ts
import { test, expect } from '@playwright/test'

test.describe('Admin 管理后台', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:3000')
  })

  test('管理员登录流程', async ({ page }) => {
    await page.fill('input[name="username"]', 'admin')
    await page.fill('input[name="password"]', '123456')
    await page.click('button[type="submit"]')
    await expect(page).toHaveURL(/.*dashboard/)
  })
})
```

**覆盖率目标**: 核心流程 100%

---

## 三、Admin 测试规范

### 3.1 优先级模块

| 优先级 | 模块 | 文件 | 测试类型 |
|--------|------|------|---------|
| P0 | 登录 | `LoginView.vue` | 组件 + E2E |
| P0 | 用户管理 | `UserManagementView.vue` | 组件 |
| P0 | 品牌管理 | `BrandManagementView.vue` | 组件 |
| P1 | 活动管理 | `CampaignListView.vue` | 组件 + E2E |
| P1 | Dashboard | `DashboardView.vue` | 组件 |
| P2 | 反馈管理 | `FeedbackManagementView.vue` | 组件 |

### 3.2 测试用例模板

```typescript
// tests/views/UserManagementView.test.ts
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import UserManagementView from '@/views/UserManagementView.vue'

// Mock API
vi.mock('@/services/userApi', () => ({
  userApi: {
    getUsers: vi.fn(() => Promise.resolve({ data: { list: [], total: 0 } })),
    createUser: vi.fn(() => Promise.resolve({ data: { id: 1 } })),
    updateUser: vi.fn(() => Promise.resolve({ data: {} })),
    deleteUser: vi.fn(() => Promise.resolve({ data: {} })),
  }
}))

describe('UserManagementView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('渲染测试', () => {
    it('应该显示用户列表表格', () => {
      // ...
    })

    it('应该显示操作按钮', () => {
      // ...
    })
  })

  describe('交互测试', () => {
    it('点击新增按钮应该打开弹窗', async () => {
      // ...
    })

    it('提交表单应该调用创建 API', async () => {
      // ...
    })
  })

  describe('错误处理', () => {
    it('API 失败时应该显示错误提示', async () => {
      // ...
    })
  })
})
```

### 3.3 Services 测试规范

```typescript
// tests/services/userApi.test.ts
import { describe, it, expect, vi } from 'vitest'
import { userApi } from '@/services/userApi'

describe('userApi', () => {
  describe('getUsers', () => {
    it('应该返回用户列表', async () => {
      const result = await userApi.getUsers({ page: 1, pageSize: 10 })
      expect(result.data.list).toBeInstanceOf(Array)
      expect(result.data.total).toBeGreaterThanOrEqual(0)
    })
  })
})
```

---

## 四、H5 测试规范

### 4.1 优先级模块

| 优先级 | 模块 | 文件 | 测试类型 |
|--------|------|------|---------|
| P0 | 活动列表 | `CampaignList.vue` | 组件 + E2E |
| P0 | 报名表单 | `CampaignForm.vue` | 组件 |
| P0 | 反馈中心 | `FeedbackCenter.vue` | 组件 + E2E |
| P1 | 我的订单 | `MyOrders.vue` | 组件 |
| P1 | 品牌登录 | `BrandLogin.vue` | 组件 + E2E |
| P2 | 分销中心 | `DistributorCenter.vue` | 组件 |

### 4.2 Logic 层测试（已覆盖良好）

H5 项目已有 55 个 Logic 层测试文件，覆盖率 ~80%。继续维护：

```
tests/unit/
├── analytics.logic.test.js
├── campaignList.logic.test.js
├── formValidation.logic.test.js
├── feedbackCenter.logic.test.js
├── ... (共 55 个)
```

### 4.3 组件测试模板

```typescript
// tests/components/CampaignList.test.ts
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import CampaignList from '@/views/CampaignList.vue'

vi.mock('@/api/campaign', () => ({
  getCampaigns: vi.fn(() => Promise.resolve({
    data: { list: [{ id: 1, name: '测试活动' }], total: 1 }
  }))
}))

describe('CampaignList', () => {
  describe('渲染测试', () => {
    it('应该显示活动卡片列表', async () => {
      const wrapper = mount(CampaignList)
      await wrapper.vm.$nextTick()
      expect(wrapper.find('.campaign-card').exists()).toBe(true)
    })
  })

  describe('交互测试', () => {
    it('点击活动应该跳转详情页', async () => {
      const mockRouter = { push: vi.fn() }
      const wrapper = mount(CampaignList, {
        global: { mocks: { $router: mockRouter } }
      })
      await wrapper.find('.campaign-card').trigger('click')
      expect(mockRouter.push).toHaveBeenCalled()
    })
  })
})
```

---

## 五、E2E 测试场景

### 5.1 Admin 核心场景

| 场景 ID | 场景名称 | 步骤数 | 优先级 |
|---------|---------|--------|--------|
| E2E-A01 | 管理员登录 | 4 | P0 |
| E2E-A02 | 创建用户 | 6 | P0 |
| E2E-A03 | 创建品牌 | 5 | P0 |
| E2E-A04 | 创建活动 | 8 | P1 |
| E2E-A05 | 查看统计 | 3 | P1 |
| E2E-A06 | 反馈处理 | 5 | P2 |

### 5.2 H5 核心场景

| 场景 ID | 场景名称 | 步骤数 | 优先级 |
|---------|---------|--------|--------|
| E2E-H01 | 浏览活动列表 | 3 | P0 |
| E2E-H02 | 活动报名 | 6 | P0 |
| E2E-H03 | 提交反馈 | 5 | P0 |
| E2E-H04 | 品牌管理员登录 | 4 | P1 |
| E2E-H05 | 查看我的订单 | 3 | P1 |

---

## 六、测试命名规范

### 6.1 文件命名

```
单元测试:    {filename}.test.ts      例: utils.test.ts
组件测试:    {ComponentName}.test.ts 例: LoginView.test.ts
E2E 测试:   {feature}.spec.ts       例: admin-flows.spec.ts
```

### 6.2 测试用例命名

```typescript
// 使用 describe/it 结构，描述行为而非实现
describe('LoginView', () => {
  describe('表单验证', () => {
    it('空用户名时应该显示错误', () => {})
    it('空密码时应该显示错误', () => {})
    it('无效凭证时应该显示错误', () => {})
  })

  describe('登录流程', () => {
    it('成功登录后应该跳转到 Dashboard', () => {})
    it('记住密码功能应该正常工作', () => {})
  })
})
```

---

## 七、Mock 策略

### 7.1 API Mock

```typescript
// tests/__mocks__/api.ts
export const mockApiResponses = {
  login: {
    success: { code: 0, data: { token: 'mock-token' } },
    failure: { code: 401, message: '用户名或密码错误' }
  },
  campaigns: {
    list: { code: 0, data: { list: [], total: 0 } }
  }
}

// 使用 MSW (推荐)
import { setupServer } from 'msw/node'
import { rest } from 'msw'

export const server = setupServer(
  rest.post('/api/v1/auth/login', (req, res, ctx) => {
    return res(ctx.json(mockApiResponses.login.success))
  })
)
```

### 7.2 Router Mock

```typescript
// tests/__mocks__/router.ts
import { vi } from 'vitest'

export const createRouterMock = () => ({
  push: vi.fn(),
  replace: vi.fn(),
  go: vi.fn(),
  back: vi.fn(),
  currentRoute: { value: { path: '/', params: {} } }
})
```

### 7.3 Store Mock

```typescript
// tests/__mocks__/store.ts
import { vi } from 'vitest'
import { createStore } from 'vuex'

export const createTestStore = (initialState = {}) => {
  return createStore({
    state: {
      user: null,
      ...initialState
    },
    mutations: {
      setUser: vi.fn()
    }
  })
}
```

---

## 八、测试执行命令

### 8.1 Admin 命令

```bash
# 运行所有单元测试
cd frontend-admin && npm run test

# 运行特定文件
npm run test -- LoginView.test.ts

# 生成覆盖率报告
npm run test -- --coverage

# 运行 E2E 测试
npx playwright test

# 运行特定 E2E 测试
npx playwright test admin-flows.spec.ts
```

### 8.2 H5 命令

```bash
# 运行所有单元测试
cd frontend-h5 && npm run test

# 运行特定文件
npm run test -- campaignList.logic.test.js

# 生成覆盖率报告
npm run test -- --coverage

# 运行 E2E 测试
npx playwright test
```

---

## 九、CI 集成

### 9.1 PR 检查

```yaml
# .github/workflows/frontend-test.yml
name: Frontend Tests

on: [pull_request]

jobs:
  admin-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - name: Install dependencies
        run: cd frontend-admin && npm ci
      - name: Run unit tests
        run: cd frontend-admin && npm run test -- --run --coverage
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          files: frontend-admin/coverage/coverage-final.json
          flags: admin

  h5-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - name: Install dependencies
        run: cd frontend-h5 && npm ci
      - name: Run unit tests
        run: cd frontend-h5 && npm run test -- --run --coverage
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          files: frontend-h5/coverage/coverage-final.json
          flags: h5
```

---

## 十、测试最佳实践

### 10.1 DO

✅ 测试行为，不测试实现
✅ 使用有意义的测试描述
✅ 保持测试独立性
✅ Mock 外部依赖
✅ 测试边界条件
✅ 测试错误处理

### 10.2 DON'T

❌ 测试中硬编码敏感数据
❌ 测试之间共享状态
❌ 过度 Mock
❌ 忽略异步操作
❌ 测试私有方法

---

## 十一、附录

### A. 测试覆盖率目标

| 层级 | 当前 | 短期目标 | 长期目标 |
|------|------|---------|---------|
| Admin Services | 54% | 70% | 80% |
| Admin Views | 0.94% | 30% | 60% |
| H5 Logic | ~80% | 85% | 90% |
| H5 Components | 0% | 30% | 60% |

### B. 参考资源

- [Vue Test Utils 文档](https://test-utils.vuejs.org/)
- [Vitest 文档](https://vitest.dev/)
- [Playwright 文档](https://playwright.dev/)
- [Testing Library](https://testing-library.com/docs/vue-testing-library/intro/)

---

**文档维护**: 每季度评审更新
