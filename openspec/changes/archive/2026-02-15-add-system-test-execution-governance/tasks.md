## 1. Baseline Alignment

- [ ] 1.1 建立 `test-plan.md` 测试范围到执行入口映射表（模块 -> 命令/脚本/CI）
  - 执行命令：`rg -n "## 5\.|TC-" test-plan.md`、`rg -n "test|integration|e2e|regression" backend/scripts .github/workflows frontend-admin/package.json frontend-h5/package.json`
  - 前置条件：`test-plan.md` 存在；仓库内测试入口文件可读取
  - 通过标准：映射表覆盖 RBAC/活动/订单/反馈/规格治理 5 个模块，且每个模块至少 1 个执行入口
  - 证据输出：`docs/testing/execution/SCOPE_MAPPING.md`

- [ ] 1.2 固化 P0/P1 用例分层与执行顺序
  - 执行命令：`rg -n "P0|P1|RBAC|活动管理|订单|反馈|规格治理" test-plan.md`
  - 前置条件：1.1 已完成；`test-plan.md` 中包含优先级定义
  - 通过标准：执行顺序文档明确 `RBAC -> 活动 -> 订单 -> 反馈 -> 规格治理`，并标注 P0/P1 边界
  - 证据输出：`docs/testing/execution/EXECUTION_ORDER.md`

- [ ] 1.3 补齐测试前置条件说明（环境、账号、初始化、环境变量）
  - 执行命令：`rg -n "DMH_INTEGRATION_BASE_URL|DMH_TEST_ADMIN_USERNAME|DMH_TEST_ADMIN_PASSWORD|localhost:8889|localhost:3000|localhost:3100" AGENTS.md README.md test-plan.md`
  - 前置条件：环境基线信息可从 `AGENTS.md`、`README.md`、`test-plan.md` 获取
  - 通过标准：文档包含服务地址、测试账号、数据库初始化步骤、`DMH_INTEGRATION_BASE_URL` 等必需变量
  - 证据输出：`docs/testing/execution/PREREQUISITES.md`

## 2. Unified Execution Matrix

- [ ] 2.1 定义后端执行矩阵（单元/集成/订单专项回归）
  - 执行命令：`cd backend && go test ./... -v`、`cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 go test ./test/integration/... -v -count=1`、`backend/scripts/run_order_mysql8_regression.sh`
  - 前置条件：1.1~1.3 已完成；`backend/scripts/run_order_mysql8_regression.sh` 可用
  - 通过标准：矩阵至少包含 `go test ./...`、`go test ./test/integration/...`、`backend/scripts/run_order_mysql8_regression.sh`
  - 证据输出：`docs/testing/execution/BACKEND_MATRIX.md`

- [ ] 2.2 定义前端执行矩阵（admin/h5 单元 + E2E）
  - 执行命令：`cd frontend-admin && npm run test && npm run test:e2e`、`cd frontend-h5 && npm run test && npm run test:e2e`
  - 前置条件：`frontend-admin/package.json`、`frontend-h5/package.json` 包含 `test` 与 `test:e2e` 命令
  - 通过标准：矩阵明确 admin/h5 两端执行命令、依赖和最小通过条件
  - 证据输出：`docs/testing/execution/FRONTEND_MATRIX.md`

- [ ] 2.3 定义异常场景验证矩阵（鉴权/参数/超时/幂等/权限）
  - 执行命令：`rg -n "异常|超时|重复|权限不足|401|403|校验" test-plan.md backend/test/integration frontend-admin/tests frontend-h5/tests`
  - 前置条件：1.2 已完成；异常用例可从 `test-plan.md` 提取
  - 通过标准：每类异常至少 1 条可执行用例，并绑定到具体入口或测试文件
  - 证据输出：`docs/testing/execution/ABNORMAL_MATRIX.md`

## 3. QA Gate and Evidence

- [ ] 3.1 落地发布前质量门禁规则（P0=100%，P1>=95%，无 Blocker/Critical）
  - 执行命令：`rg -n "P0|P1|Blocker|Critical|通过率" test-plan.md`
  - 前置条件：2.1~2.3 已完成；缺陷分级口径已定义
  - 通过标准：规则文档明确门禁阈值、判定公式、失败处理分支
  - 证据输出：`docs/testing/execution/QUALITY_GATE.md`

- [ ] 3.2 制定测试结果记录模板（通过率/失败用例/耗时/日志截图）
  - 执行命令：`mkdir -p docs/testing/execution`
  - 前置条件：质量门禁规则已定义
  - 通过标准：模板字段完整，能覆盖后端、前端、CI 三类执行结果
  - 证据输出：`docs/testing/execution/TEST_RESULT_TEMPLATE.md`

- [ ] 3.3 制定风险豁免模板（风险评估与审批链）
  - 执行命令：`mkdir -p docs/testing/execution`
  - 前置条件：3.1 已完成
  - 通过标准：模板包含未通过项、影响评估、补救计划、审批人和有效期
  - 证据输出：`docs/testing/execution/RISK_WAIVER_TEMPLATE.md`

## 4. CI Integration and Traceability

- [ ] 4.1 对齐现有 CI 守卫与执行矩阵
  - 执行命令：`rg -n "name:|pull_request:|push:|run:" .github/workflows/order-mysql8-regression.yml .github/workflows/feedback-guard.yml`
  - 前置条件：2.1~2.2 已完成；CI 工作流文件可读取
  - 通过标准：文档明确 `order-mysql8-regression.yml` 与 `feedback-guard.yml` 覆盖范围及缺口
  - 证据输出：`docs/testing/execution/CI_ALIGNMENT.md`

- [ ] 4.2 定义主干合并前最小必跑集与证据归档位置
  - 执行命令：`rg -n "go test|run_order_mysql8_regression|npm run test|test:e2e" test-plan.md README.md AGENTS.md`
  - 前置条件：4.1 已完成
  - 通过标准：必跑集包含后端、前端、订单专项回归；证据目录和命名规则明确
  - 证据输出：`docs/testing/execution/PREMERGE_MIN_SUITE.md`

- [ ] 4.3 明确失败回滚与重试策略（定位 -> 修复 -> 重跑 -> 复验）
  - 执行命令：`rg -n "ERROR|retry|重跑|回归|失败" backend/scripts/run_order_mysql8_regression.sh backend/scripts/repair_login_and_run_order_regression.sh docs/testing/REGRESSION_CHECKLIST.md`
  - 前置条件：3.1、4.2 已完成
  - 通过标准：流程包含责任人、触发条件、重跑范围和复验门禁
  - 证据输出：`docs/testing/execution/FAILURE_PLAYBOOK.md`

## 5. Validation

- [ ] 5.1 运行 OpenSpec 严格校验
  - 执行命令：`openspec validate add-system-test-execution-governance --strict --no-interactive`
  - 前置条件：`proposal.md`、`tasks.md`、`design.md` 与 spec delta 已更新
  - 通过标准：命令返回 `Change 'add-system-test-execution-governance' is valid`
  - 证据输出：`openspec/changes/add-system-test-execution-governance/validation.log`
