## ADDED Requirements

### Requirement: Test plan aligned execution scope
系统 SHALL 以 `test-plan.md` 作为系统测试执行基线，并确保测试范围覆盖 RBAC、活动管理、订单支付、反馈系统和规格治理。

#### Scenario: Execute planned module coverage
- **WHEN** 测试负责人发起一次完整系统测试执行
- **THEN** 执行范围 SHALL 覆盖 `test-plan.md` 中定义的全部 In Scope 模块
- **AND** 每个模块 SHALL 绑定至少一个可执行入口（命令、脚本或 CI 工作流）

### Requirement: Unified multi-layer test execution matrix
系统 SHALL 提供统一的分层执行矩阵，覆盖后端单元测试、后端集成测试、订单专项回归、前端单元测试和前端 E2E 测试。

#### Scenario: Run backend and frontend baseline suites
- **WHEN** 变更进入回归验证阶段
- **THEN** 系统 SHALL 执行后端 `go test` 基线与订单回归脚本
- **AND** 系统 SHALL 执行 `frontend-admin` 与 `frontend-h5` 的单元测试和 E2E 测试

### Requirement: Normal and abnormal scenario verification
系统 SHALL 在测试执行中同时覆盖正常与异常场景，异常场景至少包括鉴权失败、参数非法、权限不足、超时重试与重复请求幂等。

#### Scenario: Verify abnormal protections
- **WHEN** 测试执行到异常场景验证阶段
- **THEN** 系统 SHALL 对每类异常场景执行至少一个可重复验证的测试用例
- **AND** 失败结果 SHALL 包含可定位信息（错误码、日志或失败步骤）

### Requirement: Release QA gate for comprehensive testing
系统 SHALL 在发布前应用统一质量门禁：P0 用例通过率必须为 100%，P1 用例通过率必须不低于 95%，且不得存在 Blocker/Critical 未关闭缺陷。

#### Scenario: Enforce release gate
- **WHEN** 测试结果用于发布决策
- **THEN** 若 P0/P1 通过率或缺陷级别未满足门禁条件，系统 SHALL 阻止发布通过
- **AND** 若采用风险豁免，系统 SHALL 记录豁免理由、风险评估与审批信息

### Requirement: Traceable evidence and execution records
系统 SHALL 为每次系统测试执行产出可追溯记录，至少包含执行时间、执行范围、通过率、失败用例、关键日志/截图/报文引用。

#### Scenario: Archive execution evidence
- **WHEN** 一轮系统测试执行完成
- **THEN** 系统 SHALL 产出标准化测试结果记录
- **AND** 记录 SHALL 可用于回归审计与发布复盘
