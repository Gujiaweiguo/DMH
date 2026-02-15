# DMH 测试改进计划

> 文档版本: 1.0
> 创建日期: 2026-02-14
> 最后更新: 2026-02-14

---

## 一、执行摘要

### 1.1 当前状态

| 模块 | 覆盖率 | 目标 | 状态 |
|------|--------|------|------|
| **后端 Go** | ~76.7% | 70% | ✅ 达标 |
| **前端 Admin** | ~4.4% | 70% | 🔴 差距 65.6% |
| **前端 H5** | ~1.1% | 70% | 🔴 差距 68.9% |
| **E2E 测试** | 96.4% | 95% | ✅ 达标 |

### 1.2 核心问题

1. **集成测试环境损坏** - 大量集成测试因登录失败/服务未启动被跳过
2. **前端测试严重不足** - Admin/H5 组件测试几乎空白
3. **测试文档不完整** - 缺少前端测试策略、性能测试报告等
4. **低覆盖模块** - `logic/admin` 仅 17.5%

---

## 二、改进路线图

### Phase 1: 紧急修复 (第1周)

**目标**: 恢复集成测试执行能力

| 任务 | 优先级 | 工时 | 负责人 |
|------|--------|------|--------|
| 修复后端服务启动脚本 | P0 | 2h | - |
| 修复测试账号密码问题 | P0 | 1h | - |
| 验证集成测试可执行 | P0 | 1h | - |
| 补充 logic/admin 测试 | P0 | 4h | - |

**验收标准**:
- [ ] 集成测试执行率 > 90%
- [ ] logic/admin 覆盖率 > 50%

### Phase 2: 前端覆盖提升 (第2-3周)

**目标**: Admin 覆盖率 30%+, H5 覆盖率 20%+

| 任务 | 优先级 | 工时 | 说明 |
|------|--------|------|------|
| Admin 核心页面测试 | P1 | 8h | Dashboard, 用户管理, 品牌管理 |
| Admin Services 补测 | P1 | 4h | distributorApi, mockApi |
| H5 核心组件测试 | P1 | 6h | 活动列表, 报名表单, 反馈中心 |
| 建立 PR 测试门禁 | P2 | 2h | 覆盖率阈值检查 |

**验收标准**:
- [ ] Admin 覆盖率 > 30%
- [ ] H5 覆盖率 > 20%

### Phase 3: 文档与基建 (第4周)

**目标**: 完善测试文档体系

| 任务 | 优先级 | 工时 | 说明 |
|------|--------|------|------|
| 前端测试策略文档 | P1 | 2h | 定义测试范围、用例规范 |
| 集成测试执行报告 | P2 | 1h | 记录执行结果、失败分析 |
| 性能测试报告模板 | P2 | 1h | 基准数据、瓶颈分析模板 |
| 回归测试清单 | P2 | 1h | 发布必测项 |

**验收标准**:
- [ ] 所有测试文档创建完成
- [ ] 文档评审通过

---

## 三、详细任务分解

### 3.1 修复集成测试环境

#### 问题诊断

```
当前错误:
1. login returned 400: 用户名或密码错误
2. dial tcp 127.0.0.1:8889: connect: connection refused
```

#### 修复步骤

```bash
# Step 1: 确认测试账号存在
cd backend
mysql -u root -p dmh -e "SELECT username FROM users WHERE username='admin';"

# Step 2: 重置测试账号密码
mysql -u root -p dmh -e "
UPDATE users SET password = '\$2a\$10\$...' 
WHERE username IN ('admin', 'brand_manager');
"

# Step 3: 启动后端服务
cd deployment && docker compose -f docker-compose-simple.yml up -d dmh-api

# Step 4: 验证服务可用
curl http://localhost:8889/api/v1/health

# Step 5: 运行集成测试
cd backend
DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
DMH_TEST_ADMIN_USERNAME=admin \
DMH_TEST_ADMIN_PASSWORD=123456 \
go test ./test/integration/... -v -count=1
```

#### 受影响测试文件

| 文件 | 当前状态 | 预期状态 |
|------|---------|---------|
| `order_verify_auth_guard_integration_test.go` | SKIP | PASS |
| `permission_test.go` | SKIP | PASS |
| `rate_limiting_test.go` | SKIP | PASS |
| `rbac_integration_test.go` | SKIP | PASS |
| `security_verification_test.go` | SKIP | PASS |
| `advanced_features_performance_test.go` | SKIP | PASS |
| `rbac_performance_test.go` | SKIP | PASS |

### 3.2 前端测试策略

#### Admin 优先测试模块

| 模块 | 文件 | 覆盖率目标 | 测试类型 |
|------|------|-----------|---------|
| 用户管理 | `UserManagementView.tsx` | 60% | 组件测试 |
| 品牌管理 | `BrandManagementView.tsx` | 60% | 组件测试 |
| 活动列表 | `CampaignListView.tsx` | 60% | 组件测试 |
| 登录页面 | `LoginView.tsx` | 80% | 组件测试 |
| Dashboard | `DashboardView.tsx` | 50% | 组件测试 |

#### H5 优先测试模块

| 模块 | 文件 | 覆盖率目标 | 测试类型 |
|------|------|-----------|---------|
| 活动列表 | `CampaignList.vue` | 60% | 组件测试 |
| 报名表单 | `CampaignForm.vue` | 60% | 组件测试 |
| 反馈中心 | `FeedbackCenter.vue` | 60% | 组件测试 |
| 订单列表 | `MyOrders.vue` | 50% | 组件测试 |

#### 测试用例设计原则

1. **Happy Path**: 正常流程必须覆盖
2. **Error Handling**: 错误提示、边界条件
3. **Loading States**: 加载状态展示
4. **User Interactions**: 按钮、表单交互
5. **Data Binding**: 数据绑定正确性

### 3.3 低覆盖模块补测

#### logic/admin (当前 17.5%)

**未覆盖功能**:
- 用户 CRUD 操作
- 权限分配逻辑
- 批量操作处理

**补测用例**:

```go
// 优先级 P0
func TestAdminLogic_CreateUser(t *testing.T) {}
func TestAdminLogic_UpdateUser(t *testing.T) {}
func TestAdminLogic_DeleteUser(t *testing.T) {}
func TestAdminLogic_ListUsers(t *testing.T) {}
func TestAdminLogic_ResetPassword(t *testing.T) {}

// 优先级 P1
func TestAdminLogic_AssignRole(t *testing.T) {}
func TestAdminLogic_BatchOperation(t *testing.T) {}
```

---

## 四、测试文档清单

### 4.1 已有文档

| 文档 | 路径 | 状态 |
|------|------|------|
| 测试计划 | `docs/testing/TEST-PLAN.md` | ✅ |
| E2E 场景 | `docs/testing/E2E-SCENARIOS.md` | ✅ |
| 覆盖率追踪 | `docs/testing/COVERAGE.md` | ✅ |
| 覆盖率差距 | `docs/COVERAGE_GAP_REPORT.md` | ✅ |
| 执行进度 | `docs/TEST_EXECUTION_PROGRESS.md` | ✅ |
| 反馈测试场景 | `docs/feedback-testing-scenarios.md` | ✅ |

### 4.2 待创建文档

| 文档 | 路径 | 优先级 | 内容 |
|------|------|--------|------|
| 前端测试策略 | `docs/testing/FRONTEND_TEST_STRATEGY.md` | P1 | Admin/H5 测试规范 |
| 集成测试报告 | `docs/testing/INTEGRATION_TEST_REPORT.md` | P2 | 执行结果、失败分析 |
| 性能测试报告 | `docs/testing/PERFORMANCE_TEST_REPORT.md` | P2 | 基准数据、瓶颈 |
| 回归测试清单 | `docs/testing/REGRESSION_CHECKLIST.md` | P2 | 发布必测项 |

---

## 五、CI/CD 集成建议

### 5.1 PR 门禁

```yaml
# .github/workflows/test.yml
name: Test Coverage Check

on: [pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Run tests with coverage
        run: |
          cd backend
          go test ./... -coverprofile=coverage.out
      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$COVERAGE < 70" | bc -l) )); then
            echo "Coverage $COVERAGE% is below threshold 70%"
            exit 1
          fi

  frontend-admin:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - name: Install dependencies
        run: cd frontend-admin && npm ci
      - name: Run tests with coverage
        run: cd frontend-admin && npm run test -- --coverage
      - name: Check coverage threshold
        run: |
          # TODO: 添加覆盖率阈值检查
          echo "Frontend coverage check placeholder"
```

### 5.2 夜间 E2E 测试

```yaml
# .github/workflows/e2e-nightly.yml
name: E2E Nightly Tests

on:
  schedule:
    - cron: '0 2 * * *'  # 每天凌晨 2 点

jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run E2E tests
        run: |
          cd frontend-admin && npx playwright test
          cd frontend-h5 && npx playwright test
      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: playwright-report
          path: |
            frontend-admin/playwright-report/
            frontend-h5/playwright-report/
```

---

## 六、度量指标

### 6.1 覆盖率目标

| 阶段 | 后端 | Admin | H5 | 时间 |
|------|------|-------|-----|------|
| 当前 | 76.7% | 4.4% | 1.1% | - |
| Phase 1 | 78% | 10% | 10% | +1周 |
| Phase 2 | 80% | 30% | 20% | +3周 |
| 目标 | 80% | 70% | 70% | +8周 |

### 6.2 测试质量指标

| 指标 | 当前 | 目标 | 说明 |
|------|------|------|------|
| 测试通过率 | ~95% | 99%+ | CI 构建成功率 |
| 集成测试执行率 | ~30% | 95%+ | 非 SKIP 比例 |
| E2E 稳定性 | 96.4% | 99%+ | 减少假失败 |
| 缺陷逃逸率 | - | <5% | 生产发现的缺陷占比 |

---

## 七、风险与缓解

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|---------|
| 前端测试工时超预期 | 延期 | 高 | 分阶段交付，优先核心模块 |
| 集成测试环境不稳定 | 阻塞 | 中 | 容器化测试环境 |
| E2E 测试假失败 | 信任度下降 | 中 | 增加重试机制，优化等待策略 |
| 测试维护成本高 | 技术债 | 中 | 测试代码评审，重构优化 |

---

## 八、附录

### A. 测试命令速查

```bash
# 后端全量测试
cd backend && go test ./... -v

# 后端覆盖率
cd backend && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

# 集成测试
cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 go test ./test/integration/... -v

# Admin 单元测试
cd frontend-admin && npm run test

# Admin E2E 测试
cd frontend-admin && npx playwright test

# H5 单元测试
cd frontend-h5 && npm run test

# H5 E2E 测试
cd frontend-h5 && npx playwright test
```

### B. 参考文档

- [Go Testing 最佳实践](https://go.dev/doc/tutorial/add-a-test)
- [Vitest 文档](https://vitest.dev/)
- [Playwright 文档](https://playwright.dev/)
- [Vue 测试指南](https://vuejs.org/guide/scaling-up/testing.html)

---

**文档维护**: 每周更新覆盖率数据，每阶段结束时更新进度
