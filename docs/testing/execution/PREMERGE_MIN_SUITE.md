# Pre-merge Minimum Suite

> 每次主干合并前必须执行的最小回归集，覆盖 P0 核心链路。

## 必跑项（按顺序）

| 序号 | 分类 | 命令/脚本 | 目标 | 前置条件 | 通过标准 | 结果 |
|---|---|---|---|---|---|---|
| 1 | 后端单元 | `cd backend && go test ./... -v` | 基线回归 | Go 依赖就绪 | Exit code = 0 | ✅ PASS |
| 2 | 后端集成 | `cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 DMH_TEST_ADMIN_USERNAME=admin DMH_TEST_ADMIN_PASSWORD=123456 go test ./test/integration/... -v -count=1` | 关键链路 | API 在线 + 环境变量 | Exit code = 0 | ✅ PASS (27 suites) |
| 3 | 订单专项 | `backend/scripts/run_order_mysql8_regression.sh` | 订单鉴权/重复报名回归 | API + 测试账号有效 | 脚本输出 PASS | ✅ PASS |
| 4 | Admin 单元 | `cd frontend-admin && npm run test` | 管理端单元（14 文件 121 用例） | Node 20+、依赖已装 | Exit code = 0 | ✅ PASS |
| 5 | H5 单元 | `cd frontend-h5 && npm run test` | H5 单元（54 文件 985 用例） | Node 20+、依赖已装 | Exit code = 0 | ✅ PASS |
| 6 | Admin E2E | `cd frontend-admin && npm run test:e2e` | 管理端关键流程 | Admin 可访问 | Exit code = 0 | ✅ PASS (21/21) |
| 7 | H5 E2E | `cd frontend-h5 && npm run test:e2e` | H5 关键流程 | H5 可访问 | Exit code = 0 | ✅ PASS (7/7) |
| 8 | OpenSpec 校验 | `openspec validate --strict --no-interactive` | 规格一致性 | OpenSpec CLI 可用 | 无错误 | ✅ PASS |

## 可选项（视变更范围决定）

| 分类 | 命令 | 目标 | 触发条件 |
|---|---|---|---|
| 后端性能 | `cd backend && go test -v ./test/performance -run TestRBACPerformanceTestSuite -timeout 300s` | RBAC 性能基线 | 涉及权限性能变更 |
| 前端覆盖率 | `cd frontend-admin && npm run test:cov` | Admin 覆盖率报告 | 需要覆盖率数据 |
| 前端覆盖率 | `cd frontend-h5 && npm run test:cov` | H5 覆盖率报告 | 需要覆盖率数据 |

## 证据归档规范

- 目录：`docs/testing/execution/runs/<yyyy-mm-dd>-<run-id>/`
- 必备文件：
  - `summary.md`：通过率汇总、阻塞项、结论
  - `raw.log`：原始控制台输出（或分段日志）
  - `failures.md`（若有失败）：失败用例清单与复现步骤
- 命名示例：`docs/testing/execution/runs/2026-02-15-run001/`
