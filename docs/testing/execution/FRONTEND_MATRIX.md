# Frontend Matrix

## Admin 端

| 分类 | 命令 | 覆盖目标 | 前置条件 | 通过标准 | 结果 |
|---|---|---|---|---|---|
| 单元测试 | `cd frontend-admin && npm run test` | 视图/服务单测（14 文件 121 用例） | Node 20+、依赖已装 | Exit code = 0 | ✅ PASS (1.73s, 14/121) |
| 单元覆盖率 | `cd frontend-admin && npm run test:cov` | 单元测试 + 覆盖率报告 | Node 20+、依赖已装 | Exit code = 0 | 待填写 |
| E2E | `cd frontend-admin && npm run test:e2e` | 管理端关键流程（admin-flows、feedback、dashboard） | API/Admin 可访问 | Exit code = 0 | ✅ PASS (21/21, 16.2s) |
| E2E Headless | `cd frontend-admin && npm run test:e2e:headless` | CI 模式 E2E | API/Admin 可访问 | Exit code = 0 | 未执行 |

## H5 端

| 分类 | 命令 | 覆盖目标 | 前置条件 | 通过标准 | 结果 |
|---|---|---|---|---|---|
| 单元测试 | `cd frontend-h5 && npm run test` | 逻辑层单测（54 文件 985 用例） | Node 20+、依赖已装 | Exit code = 0 | ✅ PASS (4.38s, 54/985) |
| 单元覆盖率 | `cd frontend-h5 && npm run test:cov` | 单元测试 + 覆盖率报告 | Node 20+、依赖已装 | Exit code = 0 | 待填写 |
| E2E | `cd frontend-h5 && npm run test:e2e` | H5 关键流程（h5-flows、feedback） | API/H5 可访问 | Exit code = 0 | ✅ PASS (7/7, 7.0s) |
| E2E Headless | `cd frontend-h5 && npm run test:e2e:headless` | CI 模式 E2E | API/H5 可访问 | Exit code = 0 | 未执行 |

## 失败记录

| 端 | 失败用例 | 失败摘要 | 是否阻塞 | 处理状态 |
|---|---|---|---|---|
| Admin | 无 | - | 否 | 已完成 |
| H5 | 无 | - | 否 | 已完成 |

## 日志与产物

- Admin 日志: `docs/testing/execution/runs/<date>-<id>/frontend-admin/`
- H5 日志: `docs/testing/execution/runs/<date>-<id>/frontend-h5/`
