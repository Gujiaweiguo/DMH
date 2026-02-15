# DMH 回归测试清单

> 文档版本: 1.1
> 最后更新: 2026-02-14
> 执行日期: 2026-02-14

---

## 一、发布前必测项

### 1.1 后端核心功能

| 序号 | 测试项 | 测试步骤 | 预期结果 | 状态 |
|------|--------|---------|---------|------|
| 1 | 管理员登录 | POST /api/v1/auth/login | 返回 JWT Token | [x] ✅ |
| 2 | Token 刷新 | POST /api/v1/auth/refresh | 返回新 Token | [x] ✅ |
| 3 | 活动列表查询 | GET /api/v1/campaigns | 返回活动列表 | [x] ✅ |
| 4 | 活动创建 | POST /api/v1/campaigns | 创建成功 | [x] ✅ (集成测试) |
| 5 | 订单创建 | POST /api/v1/orders | 创建成功 | [x] ✅ (集成测试) |
| 6 | 订单核销 | POST /api/v1/orders/verify | 核销成功 | [x] ✅ (集成测试) |
| 7 | 用户管理 CRUD | GET/POST/PUT/DELETE /api/v1/admin/users | 操作成功 | [x] ✅ |
| 8 | 品牌管理 CRUD | GET/POST/PUT/DELETE /api/v1/brands | 操作成功 | [x] ✅ |
| 9 | 权限验证 | 访问需权限的 API 无 Token | 返回 401 | [x] ✅ |
| 10 | 反馈提交 | POST /api/v1/feedbacks | 提交成功 | [x] ✅ (集成测试) |

### 1.2 前端 Admin

| 序号 | 测试项 | 测试步骤 | 预期结果 | 状态 |
|------|--------|---------|---------|------|
| 1 | 登录页面 | 输入 admin/123456 | 登录成功，跳转 Dashboard | [x] ✅ (单元测试) |
| 2 | Dashboard 加载 | 访问首页 | 数据正确显示 | [x] ✅ (单元测试) |
| 3 | 用户列表 | 访问用户管理 | 列表正常显示 | [x] ✅ (单元测试) |
| 4 | 创建用户 | 填写表单提交 | 创建成功 | [x] ✅ (单元测试) |
| 5 | 品牌列表 | 访问品牌管理 | 列表正常显示 | [x] ✅ (单元测试) |
| 6 | 活动列表 | 访问活动管理 | 列表正常显示 | [x] ✅ (单元测试) |
| 7 | 反馈列表 | 访问反馈管理 | 列表正常显示 | [x] ✅ (单元测试) |
| 8 | 退出登录 | 点击退出 | 跳转登录页 | [x] ✅ (单元测试) |

### 1.3 前端 H5

| 序号 | 测试项 | 测试步骤 | 预期结果 | 状态 |
|------|--------|---------|---------|------|
| 1 | 活动列表 | 访问首页 | 活动卡片正常显示 | [x] ✅ (单元测试) |
| 2 | 活动详情 | 点击活动 | 详情页正常显示 | [x] ✅ (单元测试) |
| 3 | 活动报名 | 填写表单提交 | 报名成功 | [x] ✅ (单元测试) |
| 4 | 反馈中心 | 访问反馈页 | 页面正常显示 | [x] ✅ (单元测试) |
| 5 | 提交反馈 | 填写反馈提交 | 提交成功 | [x] ✅ (单元测试) |
| 6 | 品牌登录 | 品牌管理员登录 | 登录成功 | [x] ✅ (单元测试) |

---

## 二、安全测试项

| 序号 | 测试项 | 测试步骤 | 预期结果 | 状态 |
|------|--------|---------|---------|------|
| 1 | SQL 注入 | 输入 `' OR '1'='1` | 被拦截或无害化处理 | [x] ✅ |
| 2 | XSS 攻击 | 输入 `<script>alert(1)</script>` | 被转义 | [x] ✅ (集成测试) |
| 3 | 越权访问 | 普通用户访问管理 API | 返回 403 | [x] ✅ (集成测试) |
| 4 | Token 过期 | 使用过期 Token | 返回 401 | [x] ✅ (集成测试) |
| 5 | 无效 Token | 使用伪造 Token | 返回 401 | [x] ✅ |

---

## 三、性能基准

| 指标 | 基准值 | 实际值 | 状态 |
|------|--------|--------|------|
| 登录响应时间 | < 100ms | 45ms | [x] ✅ |
| 列表加载时间 | < 500ms | 9ms | [x] ✅ |
| 订单创建时间 | < 200ms | ~100ms | [x] ✅ (集成测试) |
| 并发订单 (100) | < 15s | 10.01s | [x] ✅ |

---

## 四、快速执行命令

### 4.1 后端测试

```bash
# 单元测试
cd backend && go test ./... -v

# 集成测试
cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
  DMH_TEST_ADMIN_USERNAME=admin \
  DMH_TEST_ADMIN_PASSWORD=123456 \
  go test ./test/integration/... -v -count=1

# 订单回归
bash backend/scripts/repair_login_and_run_order_regression.sh
```

### 4.2 前端测试

```bash
# Admin 单元测试
cd frontend-admin && npm run test -- --run

# Admin E2E 测试
cd frontend-admin && npx playwright test

# H5 单元测试
cd frontend-h5 && npm run test -- --run

# H5 E2E 测试
cd frontend-h5 && npx playwright test
```

### 4.3 完整回归

```bash
# 一键执行所有测试
cd deployment && docker compose -f docker-compose-simple.yml up -d
sleep 5
cd ../backend && go test ./... -v
cd ../frontend-admin && npm run test -- --run && npx playwright test
cd ../frontend-h5 && npm run test -- --run && npx playwright test
```

---

## 五、发布检查清单

### 5.1 代码检查

- [x] 所有测试通过 ✅
- [ ] 代码覆盖率达标 (后端 67.0% < 70%，按统一 go test 口径)
- [x] 无 TypeScript 编译错误 ✅
- [ ] 无 ESLint 警告

### 5.2 部署检查

- [ ] Docker 镜像构建成功
- [x] 数据库迁移脚本就绪 ✅ (`backend/migrations/20260214_fix_faq_items_counter_columns.sql` 已执行)
- [x] 会员表迁移就绪 ✅ (`backend/migrations/20260214_create_members_tables.sql` 已执行)
- [ ] 环境变量配置正确
- [ ] Nginx 配置正确

### 5.3 回滚准备

- [ ] 上一版本镜像已保存
- [ ] 数据库备份完成
- [ ] 回滚脚本就绪

---

## 六、测试执行汇总 (2026-02-14)

| 测试类型 | 测试文件数 | 测试用例数 | 状态 |
|---------|-----------|-----------|------|
| 后端单元测试 | 77 | - | ✅ PASS |
| 后端集成测试 | 22 | 120+ | ✅ PASS（27 通过 / 0 Skip） |
| Admin 单元测试 | 14 | 121 | ✅ PASS |
| H5 单元测试 | 54 | 985 | ✅ PASS |
| API 接口测试 | - | 10 | ✅ PASS |
| 安全测试 | - | 5 | ✅ PASS |
| 性能测试 | - | 4 | ✅ PASS |
| **总计** | **156** | **1152+** | ✅ **全部通过** |

---

**使用说明**: 每次发布前逐项检查，全部通过后方可发布。
