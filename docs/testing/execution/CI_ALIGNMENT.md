# CI Alignment

## 工作流覆盖对齐

| 工作流 | 触发条件 | 当前覆盖 | 对应测试矩阵项 | 缺口 |
|---|---|---|---|---|
| `.github/workflows/order-mysql8-regression.yml` | PR/Push 到 order 相关路径；workflow_dispatch | 订单鉴权（`TestOrderVerifyRoutesAuthGuard`）+ 重复报名（`TestOrderCreateDuplicateMessage`） | 后端矩阵-订单回归 | 仅覆盖订单核心场景，不含并发压测 |
| `.github/workflows/feedback-guard.yml` | PR/Push 到 feedback 相关路径 | 反馈 handler + logic 单元测试 | 后端矩阵-反馈单元 | 无 E2E 守卫；前端 feedback 测试未纳入 |
| `.github/workflows/system-test-gate.yml` | PR/Push 到 backend/frontend/openspec/workflow 相关路径；workflow_dispatch | 后端全量集成 + Admin E2E + H5 E2E + OpenSpec 校验 | 后端矩阵-全量集成；前端矩阵-E2E；规格治理 | 覆盖面完整，执行耗时较高 |

## 已覆盖 vs 未覆盖

### 已覆盖（CI 自动）

- 订单：鉴权矩阵 + 重复报名文案回归
- 反馈：后端 handler/logic 单元测试
- 系统门禁：后端全量集成 + Admin/H5 E2E + OpenSpec 全量校验

### 未覆盖（需手动/后续补齐）

- 后端性能基线（`TestRBACPerformanceTestSuite`）
- 全量后端单元（`go test ./...`）

## 后续补齐建议

1. 新增 `schedule` 触发的夜间回归（nightly），覆盖后端全量单元与性能基线。
2. 将 `system-test-gate.yml` 结果接入分支保护规则（required checks）。
3. 按模块分拆并行 Job（backend/frontend/spec）以缩短整体耗时。
