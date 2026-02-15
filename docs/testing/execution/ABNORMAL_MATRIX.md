# Abnormal Scenario Matrix

> 用例来源：`test-plan.md` 第 5 章（类型=异常）

| 异常类型 | 代表用例ID | 触发方式 | 执行入口 | 期望结果 | 实际结果 |
|---|---|---|---|---|---|
| 鉴权失败（401） | TC-RBAC-06 | 使用过期/无效 token 访问受保护 API | `backend/test/integration/rbac_integration_test.go` | 返回 401，拒绝访问 | 待填写 |
| 权限不足（403） | TC-RBAC-02/04/07/10 | H5 尝试创建高权限角色 / 越权访问 / participant 审批提现 | `backend/test/integration/rbac_integration_test.go`；`backend/test/integration/order_verify_auth_guard_integration_test.go` | 返回 403，数据不泄露 | 待填写 |
| 参数非法（400） | TC-CAMP-05 | 缺少必填字段直接保存 | `frontend-h5/tests/unit/campaignEditor.logic.test.js` | 返回校验错误，禁止保存 | 待填写 |
| 字段格式错误 | TC-ORDER-04 | 提交非法手机号/邮箱 | `backend/api/internal/logic/order/create_order_logic_test.go` | 返回字段校验错误并指明字段 | 待填写 |
| 重复请求/幂等 | TC-ORDER-02/07 | 重复手机号报名 / 重复核销 | `backend/test/integration/order_verify_auth_guard_integration_test.go` | 幂等防护生效，返回业务错误 | 待填写 |
| 状态冲突 | TC-ORDER-03 | 活动无效/已过期时下单 | `backend/api/internal/logic/order/create_order_logic_test.go` | 返回业务错误，不创建订单 | 待填写 |
| 请求超时 | TC-RBAC-13 | 模拟后端超时 | 前端分销监控页面手动验证 | 前端显示明确超时提示并提供重试入口 | 待填写 |
| 评分超范围 | TC-FEED-03 | rating=0 或 6 | `backend/api/internal/logic/feedback/feedbacklogic_test.go` | 返回 400，评分范围错误 | 待填写 |
| 越权访问 | TC-FEED-05/07/12 | 普通用户访问他人反馈 / 更新状态 / 访问统计 | `backend/test/integration/feedback_handler_integration_test.go` | 返回 403 | 待填写 |

## 执行建议

- 优先执行后端集成测试中的异常分支（RBAC、订单）。
- 前端异常场景多为表单校验，可在 E2E 中覆盖或在单元测试中验证逻辑函数返回值。
