# Failure Playbook

> 测试执行失败时的标准处理流程与升级策略

## 标准流程

1. **定位失败**：用例/模块/环境/数据。
2. **判断是否阻塞**：是否触发 QA Gate（Blocker/Critical/P0 失败）。
3. **修复**：代码/测试/环境/数据。
4. **重跑**：按最小必跑集重跑相关范围。
5. **复验**：确认全绿后更新测试记录与结论。

## 处理矩阵

| 失败类型 | 首要动作 | 重跑范围 | 负责人 | 验收条件 |
|---|---|---|---|---|
| 后端单元失败 | 修复对应逻辑/测试代码 | 后端单元 + 相关集成 | 后端负责人 | 相关命令 Exit code = 0 |
| 集成测试失败 | 排查环境/接口/数据问题 | 集成 + 订单专项 | 后端负责人 | 关键链路全绿 |
| 订单回归失败 | 对照 `order-mysql8-regression.yml` 日志定位 | 订单回归脚本 | 后端负责人 | 脚本输出 PASS |
| 前端单元失败 | 修复组件/逻辑函数 | 对应端单元 + 该端 E2E | 前端负责人 | 该端 `npm run test` 全绿 |
| E2E 失败 | 区分脚本问题或真实缺陷 | 该端 E2E + 相关回归点 | 前端负责人 | 业务流程可通过 |
| CI 守卫失败 | 对照 workflow 日志定位 | 对应 workflow 全量重跑 | 责任模块 owner | CI 通过（绿色勾） |
| OpenSpec 校验失败 | 检查 Scenario 格式/Requirement 命名 | 修复后重跑校验 | 规格维护者 | `Change 'xxx' is valid` |

## 常见失败原因与排查

| 失败表现 | 可能原因 | 排查命令/动作 |
|---|---|---|
| `DMH_INTEGRATION_BASE_URL` 连接失败 | 后端服务未启动/端口被占用 | `lsof -i :8889`；`docker ps` |
| 登录返回 400 | 测试账号密码不匹配 | 执行 `backend/scripts/repair_login_and_run_order_regression.sh` |
| 集成测试 SKIP | 环境变量未设置 | 检查 `DMH_INTEGRATION_BASE_URL`/`DMH_TEST_ADMIN_USERNAME`/`DMH_TEST_ADMIN_PASSWORD` |
| 前端 E2E 超时 | 前端未启动/端口被占用 | `lsof -i :3000`/`lsof -i :3100` |
| `go test ./...` 失败 | 依赖未下载/代码编译错误 | `go mod download`；`go build ./...` |
| OpenSpec 校验失败 | Scenario 格式错误（非 `#### Scenario:`） | `openspec show <change> --json --deltas-only` |

## 升级与通报

| 缺陷级别 | 通报对象 | 时限 | 处理要求 |
|---|---|---|---|
| Blocker | 技术负责人 + 发布负责人 | 立即 | 阻塞发布，修复后重跑全量回归 |
| Critical | 技术负责人 + 发布负责人 | 4 小时内 | 阻塞发布，修复后重跑相关范围 |
| Major | 模块负责人 | 24 小时内 | 记录缺陷池，约定修复窗口 |
| Minor | 模块负责人 | 下个迭代 | 记录缺陷池，优先级较低 |

## 回滚策略

1. 若修复超过 2 小时未果，考虑回滚到上一稳定版本。
2. 回滚后重新执行预合并最小集。
3. 记录回滚原因与影响范围。
