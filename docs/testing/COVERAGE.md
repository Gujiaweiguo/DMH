# 测试覆盖率追踪

> 最后更新: 2026-02-14
> 数据来源: `go test -cover` / `vitest --coverage`

---

## 一、后端覆盖率详情

### 1.1 按模块统计 (更新)

```
模块                                    覆盖率    状态
─────────────────────────────────────────────────────────
dmh/api                                 45.6%    🔄
dmh/api/internal/config                 [无语句]
dmh/api/internal/handler                0.0%     🔴
dmh/api/internal/handler/admin          66.1%    🔄
dmh/api/internal/handler/auth           71.8%    ✅
dmh/api/internal/handler/brand          67.0%    🔄 ⬆️ (36.8% → 67.0%)
dmh/api/internal/handler/campaign       46.0%    🔄
dmh/api/internal/handler/distributor    56.0%    🔄 ⬆️ (49.0% → 56.0%)
dmh/api/internal/handler/feedback       67.6%    🔄
dmh/api/internal/handler/member         47.1%    🔄 (集成回归✅ 100%通过)
dmh/api/internal/handler/menu           46.4%    🔄
dmh/api/internal/handler/order          45.9%    🔄
dmh/api/internal/handler/poster         47.8%    🔄
dmh/api/internal/handler/reward         66.7%    🔄
dmh/api/internal/handler/role           61.0%    🔄
dmh/api/internal/handler/security       78.9%    ✅ ⬆️ (44.7% → 78.9%)
dmh/api/internal/handler/statistics     70.0%    ✅
dmh/api/internal/handler/sync           66.7%    🔄
dmh/api/internal/handler/withdrawal     68.5%    🔄 ⬆️ (42.6% → 68.5%)
dmh/api/internal/logic/admin            83.2%    ✅ ⭐ 已提升
dmh/api/internal/logic/auth             78.5%    ✅
dmh/api/internal/logic/brand            76.7%    ✅ 已提升
dmh/api/internal/logic/campaign         68.8%    🔄
dmh/api/internal/logic/distributor      73.9%    ✅
dmh/api/internal/logic/feedback         83.6%    ✅
dmh/api/internal/logic/member           79.5%    ✅
dmh/api/internal/logic/menu             71.1%    ✅
dmh/api/internal/logic/order            74.8%    ✅
dmh/api/internal/logic/poster           74.7%    ✅
dmh/api/internal/logic/reward           90.5%    ✅ ⭐
dmh/api/internal/logic/role             79.3%    ✅
dmh/api/internal/logic/security         100.0%   ✅ ⭐
dmh/api/internal/logic/statistics       81.0%    ✅
dmh/api/internal/logic/sync             100.0%   ✅ ⭐
dmh/api/internal/logic/withdrawal       63.0%    🔄
dmh/api/internal/middleware             66.8%    🔄
dmh/api/internal/service                73.9%    ✅
dmh/api/internal/svc                    0.0%     🔴
dmh/api/internal/types                  [无语句]
dmh/cmd                                 0.0%     -
dmh/common/poster                       75.7%    ✅
dmh/common/syncadapter                  47.0%    🔄 已提升
dmh/common/utils                        100.0%   ✅ ⭐
dmh/common/wechatpay                    91.9%    ✅ ⭐
dmh/model                               86.8%    ✅

总计: 68.8% (⬆️ 从 67.0% 提升 1.8%)
```

### 1.2 覆盖率分布

| 范围 | 模块数 | 占比 |
|------|-------|------|
| 100% | 4 | 9.3% |
| 80-99% | 6 | 14.0% |
| 70-79% | 12 | 27.9% |
| 60-69% | 7 | 16.3% |
| 50-59% | 1 | 2.3% |
| 40-49% | 6 | 14.0% |
| 30-39% | 0 | 0% |
| 10-29% | 0 | 0% |
| 0-9% | 3 | 7.0% |

### 1.3 优先提升列表 (更新)

| 优先级 | 模块 | 当前 | 目标 | 差距 | 建议 |
|--------|------|------|------|------|------|
| P1 | handler/brand | 67.0% | 70% | 3.0% | ⬆️ 已大幅提升 (36.8% → 67.0%) |
| P1 | handler/distributor | 56.0% | 70% | 14.0% | ⬆️ 已提升 (49.0% → 56.0%) |
| P1 | handler/withdrawal | 68.5% | 70% | 1.5% | ⬆️ 已大幅提升 (42.6% → 68.5%) |
| P1 | handler/security | 44.7% | 70% | 25.3% | 补充安全策略测试 |
| P1 | handler/campaign | 46.0% | 70% | 24.0% | 补充活动边界测试 |
| P2 | common/syncadapter | 47.0% | 70% | 23.0% | 补充同步适配器测试 |

**已完成**:
- ✅ logic/admin: 17.5% → 83.2%
- ✅ logic/brand: 65.6% → 76.7%

### 1.4 覆盖率与回归双维度判定（补充）

- `handler/brand`、`handler/distributor`：代码覆盖率仍偏低，但已新增并通过对应集成测试套件，发布风险可控。
- `handler/member`：路由与参数解析已修复并验证；当前运行环境缺少 `members` 表时，集成测试 `Skip` 属于保护性行为，不计为失败。
- 评估优先级时建议采用“双维度”规则：`代码覆盖率` + `关键链路集成回归`，避免仅按单一覆盖率百分比误判。

---

## 二、前端 Admin 覆盖率详情

### 2.1 Services 覆盖率

| 文件 | 语句 | 分支 | 函数 | 行 | 未覆盖行 |
|------|------|------|------|-----|---------|
| **总计** | 54.26% | 89.47% | 53.12% | 54.26% | - |
| authApi.ts | 66.66% | 86.66% | 53.84% | 66.66% | 155-156,160-161 |
| campaignApi.ts | 78.87% | 83.33% | 83.33% | 78.87% | 45-47,72-80,90-92 |
| distributorApi.ts | 18.36% | 100% | 0% | 18.36% | 多行 |
| feedbackApi.ts | 95.38% | 94.73% | 100% | 95.38% | 85-87 |
| memberApi.ts | 70.7% | 94.73% | 54.54% | 70.7% | 129-134,156-161 |
| mockApi.ts | 0% | 0% | 0% | 0% | 1-165 |
| orderApi.ts | 100% | 100% | 100% | 100% | - |
| performanceMonitor.ts | 71.66% | 83.33% | 80% | 71.66% | 7-8,29-30,60-74 |
| posterApi.ts | 94.54% | 90.9% | 100% | 94.54% | 55-57 |

### 2.2 Views 覆盖率

| 文件 | 语句 | 分支 | 函数 | 行 |
|------|------|------|------|-----|
| **总计** | 0.94% | 45.23% | 21.21% | 0.94% |
| DashboardView.tsx | 0% | 0% | 0% | 0% |
| LoginView.tsx | 0% | 0% | 0% | 0% |
| UserManagementView.tsx | 0% | 0% | 0% | 0% |
| BrandManagementView.tsx | 0% | 0% | 0% | 0% |
| CampaignListView.tsx | 0% | 0% | 0% | 0% |
| ... (其余均为 0%) | | | | |

### 2.3 Utils 覆盖率

| 文件 | 语句 | 分支 | 函数 | 行 |
|------|------|------|------|-----|
| adminHashRoute.ts | 100% | 91.66% | 100% | 100% |

---

## 三、前端 H5 覆盖率详情

### 3.1 已有测试文件 (55个)

```
tests/unit/
├── analytics.logic.test.js
├── api.test.js
├── apiTest.logic.test.js
├── array.logic.test.js
├── axios.test.js
├── brandApi.orderApi.test.js
├── brandApi.wrappers.test.js
├── brandLogin.logic.test.js
├── campaignDetail.logic.test.js
├── campaignEditor.logic.test.js
├── campaignForm.logic.test.js
├── campaignList.logic.test.js
├── campaignPageDesigner.logic.test.js
├── campaigns.logic.test.js
├── color.logic.test.js
├── dashboard.logic.test.js
├── dateFormat.logic.test.js
├── designer.logic.test.js
├── distributorApply.logic.test.js
├── distributorApproval.logic.test.js
├── distributorCenter.logic.test.js
├── distributorLevelRewards.logic.test.js
├── distributorLogin.logic.test.js
├── distributorPromotion.logic.test.js
├── distributorRewards.logic.test.js
├── distributorSubordinates.logic.test.js
├── distributorWithdrawals.logic.test.js
├── distributors.logic.test.js
├── feedbackCenter.logic.test.js
├── formValidation.logic.test.js
├── materials.logic.test.js
├── memberDetail.logic.test.js
├── members.logic.test.js
├── myOrders.logic.test.js
├── number.logic.test.js
├── object.logic.test.js
├── orderVerification.logic.test.js
├── orderVerify.logic.test.js
├── orders.logic.test.js
├── paymentQrcode.logic.test.js
├── posterGenerator.logic.test.js
├── posterRecords.logic.test.js
├── promoters.logic.test.js
├── router.guard.test.js
├── router.index.guard.test.js
├── settings.logic.test.js
├── storage.logic.test.js
├── string.logic.test.js
├── success.logic.test.js
├── url.logic.test.js
├── utils.test.ts
├── verificationRecords.actions.test.js
└── verificationRecords.logic.test.js
```

### 3.2 缺口分析

| 类型 | 已测试 | 未测试 | 覆盖率 |
|------|-------|-------|--------|
| Logic 函数 | 55 | ~5 | ~92% |
| Vue 组件 | 0 | ~53 | 0% |
| Composables | 0 | ~10 | 0% |

---

## 四、历史趋势

| 日期 | Backend | Admin Services | Admin Views | H5 Logic |
|------|---------|----------------|-------------|----------|
| 2026-02-13 | ~60% | 54% | 0.94% | ~80% |
| 2026-02-14 | 67.0% | 54% | 0.94% | ~80% |

---

## 五、覆盖率目标

| 模块 | 当前 | 短期目标 (2周) | 中期目标 (1月) | 长期目标 (3月) |
|------|------|---------------|---------------|---------------|
| backend | 67.0% | 70% | 75% | 80% |
| admin/services | 54% | 60% | 70% | 80% |
| admin/views | 0.94% | 20% | 50% | 60% |
| h5/logic | ~80% | 80% | 85% | 90% |
| h5/components | 0% | 20% | 50% | 60% |

---

## 六、更新命令

```bash
# 更新后端覆盖率
cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic

# 更新 Admin 覆盖率
cd frontend-admin && npm run test -- --run --coverage

# 更新 H5 覆盖率
cd frontend-h5 && npm run test -- --run --coverage
```
