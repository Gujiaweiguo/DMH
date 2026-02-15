# Execution Order

## 执行顺序（默认）

1. RBAC（P0）
2. 活动管理（P0）
3. 订单支付（P0）
4. 反馈系统（P1）
5. 规格治理（P1）

## 分层策略

- P0: 必须 100% 通过
- P1: 通过率 >= 95%

## 阶段执行建议

- 阶段 1（RBAC）
  - `backend/scripts/run_tests.sh`
  - 失败即阻塞，先修复再进入下一阶段
- 阶段 2（活动管理）
  - `cd backend && go test ./test/integration/campaign_* -v -count=1`
  - `cd frontend-h5 && npm run test && npm run test:e2e`
- 阶段 3（订单支付）
  - `backend/scripts/run_order_mysql8_regression.sh`
  - `backend/scripts/run_order_logic_tests.sh`
- 阶段 4（反馈系统）
  - `cd backend && go test ./api/internal/handler/feedback ./api/internal/logic/feedback`
  - `cd frontend-admin && npm run test && npm run test:e2e`
- 阶段 5（规格治理）
  - `openspec validate --strict --no-interactive`
  - 人工核查 `openspec/changes/archive/ARCHIVE_STATUS_INDEX.md`

## 每阶段记录

| 阶段 | 开始时间 | 结束时间 | 执行人 | 结果 | 备注 |
|---|---|---|---|---|---|
| RBAC | 2026-02-15T09:55:00Z | 2026-02-15T09:56:00Z | AI Agent | ✅ PASS | 单元/中间件/集成/性能/竞态全部通过 |
| 活动管理 | 2026-02-15T09:56:30Z | 2026-02-15T10:54:00Z | AI Agent | ✅ PASS | 后端集成 PASS；H5 单元 54/985 PASS；H5 E2E 7/7 PASS |
| 订单支付 | 2026-02-15T09:57:20Z | 2026-02-15T09:57:45Z | AI Agent | ✅ PASS | 订单回归 PASS；订单逻辑 PASS |
| 反馈系统 | 2026-02-15T09:57:50Z | 2026-02-15T10:54:30Z | AI Agent | ✅ PASS | 反馈单元 PASS；Admin 单元 14/121 PASS；Admin E2E 21/21 PASS |
| 规格治理 | 2026-02-15T09:59:05Z | 2026-02-15T10:55:00Z | AI Agent | ✅ PASS | OpenSpec specs 5/5 PASS；变更 valid；补充后端全量集成 27 suites PASS |
