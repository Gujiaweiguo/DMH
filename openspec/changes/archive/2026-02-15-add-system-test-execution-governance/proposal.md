# Change: 建立系统级全面测试执行治理流程

## Why
- `test-plan.md` 已定义完整范围、策略与用例，但当前 OpenSpec 尚缺少“如何统一执行全部测试任务”的规范化要求，导致执行口径可能不一致。
- 现有自动化能力分散在后端脚本、前端命令与部分 CI 工作流中，缺少统一 QA Gate 与结果追溯约束，影响回归稳定性与发布可控性。

## What Changes
- 新增 `system-test-execution` 能力规格，定义覆盖范围、执行编排、异常场景验证、质量门禁与证据归档要求。
- 明确以 `test-plan.md` 作为系统测试任务基线，要求执行时覆盖正常与异常场景并保留可追溯结果。
- 统一本地与 CI 的执行口径（后端单元/集成、订单回归、前端单元/E2E），形成可复现的最小执行矩阵。
- 定义发布前准入标准（P0/P1 通过率、阻塞缺陷门禁）与例外处理机制（风险豁免记录）。

## Impact
- Affected specs:
  - system-test-execution (new)
- Affected code:
  - `test-plan.md`
  - `backend/scripts/run_tests.sh`
  - `backend/scripts/run_order_mysql8_regression.sh`
  - `backend/scripts/run_order_logic_tests.sh`
  - `frontend-admin/package.json`
  - `frontend-h5/package.json`
  - `.github/workflows/order-mysql8-regression.yml`
  - `.github/workflows/feedback-guard.yml`
- Behavior changes: none（本提案聚焦测试执行治理与验收约束，不引入业务功能变更）

## Out of Scope
- 不在本提案中新增业务功能或修改业务接口行为。
- 不在本提案中引入新的第三方测试平台或替换现有测试框架。
