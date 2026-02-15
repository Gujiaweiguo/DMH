# Scope Mapping

| 模块 | 来源用例范围 | 执行入口（命令/脚本/CI） | 优先级 | 当前状态 | 证据路径 |
|---|---|---|---|---|---|
| RBAC | `TC-RBAC-01` ~ `TC-RBAC-13` | `backend/scripts/run_tests.sh`；`cd backend && go test ./test/integration/... -v -count=1` | P0 | ✅ **已执行通过** | `docs/testing/execution/runs/2026-02-15-run001/rbac/` |
| 活动管理 | `TC-CAMP-01` ~ `TC-CAMP-11` | `cd backend && go test ./test/integration/campaign_* -v -count=1`；`cd frontend-h5 && npm run test && npm run test:e2e` | P0 | ✅ **已执行通过**（后端集成 + H5 单元/E2E） | `docs/testing/execution/runs/2026-02-15-run001/campaign/` |
| 订单支付 | `TC-ORDER-01` ~ `TC-ORDER-10` | `backend/scripts/run_order_mysql8_regression.sh`；`backend/scripts/run_order_logic_tests.sh`；`order-mysql8-regression.yml` | P0 | ✅ **已执行通过** | `docs/testing/execution/runs/2026-02-15-run001/order/` |
| 反馈系统 | `TC-FEED-01` ~ `TC-FEED-13` | `cd backend && go test ./api/internal/handler/feedback ./api/internal/logic/feedback`；`feedback-guard.yml`；前端 feedback E2E | P1 | ✅ **已执行通过**（后端单元 + Admin 单元/E2E + H5 E2E） | `docs/testing/execution/runs/2026-02-15-run001/feedback/` |
| 规格治理 | `TC-SG-01` ~ `TC-SG-02` | `openspec validate --strict --no-interactive`；归档索引检查 | P1 | ✅ **已执行通过** | `docs/testing/execution/runs/2026-02-15-run001/spec-governance/` |

## 说明

- 用例范围来自 `test-plan.md` 第 5 章。
- 执行时间：2026-02-15 09:55:00 - 10:55:00
- 本轮已补充完成前端 E2E 与后端全量集成测试（Docker Compose 环境）。
