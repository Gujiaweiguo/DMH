## Context

DMH 已具备多模块测试资产（后端单元/集成、前端 Vitest/Playwright、部分 CI 守卫），且 `test-plan.md` 已定义系统级测试范围与准入标准。

当前问题是执行治理层缺口：
- 测试入口分散，缺少单一规范约束“何时跑、跑哪些、如何判定通过”。
- 对异常场景与证据归档的执行一致性不足。
- 发布前 QA Gate 未以 OpenSpec requirement 形式固化。

## Goals / Non-Goals

- Goals:
  - 建立“基于 `test-plan.md` 的系统测试执行规范”，覆盖范围、策略、执行矩阵、门禁和证据归档。
  - 将正常与异常场景验证纳入统一可追溯执行闭环。
  - 与现有脚本和 CI 守卫对齐，不破坏既有工作流。
- Non-Goals:
  - 不新增业务功能。
  - 不更换现有测试框架或引入新外部测试平台。
  - 不在提案阶段实施具体代码改造。

## Decisions

- Decision 1: 新增独立 capability `system-test-execution`
  - Why: 这是跨模块的执行治理能力，不属于单一业务域（campaign/order/rbac/feedback），独立能力更利于长期维护。

- Decision 2: 采用“统一执行矩阵 + 质量门禁 + 证据归档”三层结构
  - Why: 既能约束可执行性（run），又能约束可发布性（gate），并保证可审计（evidence）。

- Decision 3: 以现有脚本和工作流为基础能力
  - Why: 降低变更风险，优先复用 `backend/scripts/run_order_mysql8_regression.sh`、`go test`、`npm run test`/`test:e2e`、现有 GitHub workflows。

## Alternatives Considered

- 方案 A：把测试执行要求拆分到每个业务 capability 中（campaign/order/rbac/feedback）
  - 优点：紧贴业务。
  - 缺点：重复定义执行规则，难统一 QA Gate，跨模块回归流程碎片化。
  - 结论：不采用。

- 方案 B：扩展 `spec-governance` 承担测试执行治理
  - 优点：复用治理类 capability。
  - 缺点：`spec-governance` 当前聚焦归档一致性与追溯索引，混入测试执行会造成职责混杂。
  - 结论：不采用。

## Risks / Trade-offs

- 风险：规范要求与当前 CI 覆盖不完全一致，短期可能出现“规范已定义但自动化尚未完全落地”的差距。
  - Mitigation：在 tasks 中显式定义最小必跑集和手动补充验证流程，并要求记录证据。

- 风险：P0/P1 门禁执行严格后，发布节奏可能放缓。
  - Mitigation：引入风险豁免记录模板，在可控风险前提下允许带条件发布。

## Migration Plan

1. 先通过本变更固化执行规范与验收标准。
2. 后续实现阶段按 tasks 逐项落地执行矩阵与证据模板。
3. 将执行结果纳入发布流程评审，形成常态化质量门禁。

## Open Questions

- 是否需要新增“每日定时回归（nightly）”作为强制 requirement（当前先定义为可选扩展）。
