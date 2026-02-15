# Backend Matrix

## 执行矩阵

| 层级 | 命令/脚本 | 覆盖目标 | 前置条件 | 通过标准 | 结果 |
|---|---|---|---|---|---|
| RBAC 单元 | `cd backend && go test -v ./api/internal/service -run TestPasswordServiceTestSuite -timeout 30s` | 密码服务 | Go 依赖就绪 | Exit code = 0 | ✅ PASS (4.42s) |
| RBAC 单元 | `cd backend && go test -v ./api/internal/service -run TestAuditServiceTestSuite -timeout 30s` | 审计服务 | Go 依赖就绪 | Exit code = 0 | ✅ PASS (0.02s) |
| RBAC 单元 | `cd backend && go test -v ./api/internal/service -run TestSessionServiceTestSuite -timeout 30s` | 会话服务 | Go 依赖就绪 | Exit code = 0 | ✅ PASS (1.03s) |
| RBAC 中间件 | `cd backend && go test -v ./api/internal/middleware -run TestPermissionMiddlewareTestSuite -timeout 30s` | 权限中间件 | Go 依赖就绪 | Exit code = 0 | ✅ PASS |
| RBAC 集成 | `cd backend && go test -v ./test/integration -run TestRBACIntegrationTestSuite -timeout 60s` | RBAC 全链路 | 后端服务在线 + 环境变量 | Exit code = 0 | ✅ PASS |
| RBAC 性能 | `cd backend && go test -v ./test/performance -run TestRBACPerformanceTestSuite -timeout 300s` | RBAC 性能基线 | 后端服务在线 | Exit code = 0 | ✅ PASS (34.29s) |
| 全量单元 | `cd backend && go test ./... -v` | 逻辑层与基础服务 | Go 依赖就绪 | Exit code = 0 | ✅ PASS（性能项含预期 SKIP） |
| 全量集成 | `cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 go test ./test/integration/... -v -count=1` | 核心 API 关键链路 | 后端服务在线 + 环境变量 | Exit code = 0 | ✅ PASS (27 suites, 4.078s) |
| 订单回归 | `backend/scripts/run_order_mysql8_regression.sh` | 订单鉴权与重复报名回归 | API + 测试账号有效 | 脚本输出 PASS | ✅ PASS (0.15s) |
| 订单逻辑 | `backend/scripts/run_order_logic_tests.sh` | 订单创建/核销单元 + 冒烟 | Go 依赖就绪 | Exit code = 0 | ✅ PASS (2s) |
| 反馈单元 | `cd backend && go test ./api/internal/handler/feedback ./api/internal/logic/feedback` | 反馈 handler/logic | Go 依赖就绪 | Exit code = 0 | ✅ PASS (0.02s) |

## 环境变量

| 变量名 | 命令中默认值 | 说明 |
|---|---|---|
| `DMH_INTEGRATION_BASE_URL` | `http://localhost:8889` | 集成测试 API 基址 |
| `DMH_TEST_ADMIN_USERNAME` | `admin` | 集成测试管理员账号 |
| `DMH_TEST_ADMIN_PASSWORD` | `123456` | 集成测试管理员密码 |

## 日志与产物

- 控制台日志路径: `docs/testing/execution/runs/<date>-<id>/backend/`
- 失败报文路径: `docs/testing/execution/runs/<date>-<id>/backend/failures/`
