# Prerequisites

> 来源：`test-plan.md` 第 3.3 章、`AGENTS.md`、`README.md`

## 环境地址

| 服务 | 地址 | 用途 | 检查命令 |
|---|---|---|---|
| 后端 API | `http://localhost:8889` | API 测试 | `curl -s -o /dev/null -w "%{http_code}" http://localhost:8889/api/v1/auth/login` |
| 管理后台 | `http://localhost:3000` | E2E 测试 | `curl -s -o /dev/null -w "%{http_code}" http://localhost:3000` |
| H5 前端 | `http://localhost:3100` | E2E 测试 | `curl -s -o /dev/null -w "%{http_code}" http://localhost:3100` |

## 测试账号

| 角色 | 用户名 | 密码 | 用途 |
|---|---|---|---|
| 平台管理员 | `admin` | `123456` | 后端集成测试、Admin E2E |
| 品牌管理员 | `brand_manager` | `123456` | 活动管理、核销测试 |

## 必要环境变量

| 变量名 | 默认值 | 说明 | 已确认 |
|---|---|---|---|
| `DMH_INTEGRATION_BASE_URL` | `http://localhost:8889` | 集成测试 API 基址 | [x] |
| `DMH_TEST_ADMIN_USERNAME` | `admin` | 集成测试管理员账号 | [x] |
| `DMH_TEST_ADMIN_PASSWORD` | `123456` | 集成测试管理员密码 | [x] |

## 初始化检查

### 服务层

- [x] MySQL 可连接且测试库已初始化
  - 检查：`mysql -h localhost -u root -p -e "SHOW DATABASES LIKE 'dmh'"`
  - 初始化：`mysql -h localhost -u root -p dmh < backend/scripts/init.sql`
- [x] Redis 可连接（可选）
  - 检查：`redis-cli ping`
- [x] 后端服务可访问
  - 检查：`curl http://localhost:8889/api/v1/auth/login`
- [x] 前端依赖已安装
  - Admin：`cd frontend-admin && npm install`
  - H5：`cd frontend-h5 && npm install`

### 数据层

- [x] 测试数据已准备
  - 跨品牌测试数据（品牌 A/B）
  - 活动有效/过期样本
  - 可核销/已核销订单样本
- [x] 测试账号密码与数据库一致
  - 若不一致，执行 `backend/scripts/repair_login_and_run_order_regression.sh`

### 工具层

- [x] Go 环境：`go version` >= 1.23 → go1.25.3 ✅
- [x] Node 环境：`node --version` >= 20 → v22.21.1 ✅
- [x] OpenSpec CLI：`openspec --version`（可选）→ 0.19.0 ✅
- [x] Playwright 浏览器：`npx playwright install`（前端 E2E）

## 快速初始化脚本

```bash
# Docker 环境启动
cd deploy && docker compose -f docker-compose-simple.yml up -d

# 后端服务启动（非 Docker）
cd backend && go run api/dmh.go -f api/etc/dmh-api.yaml &

# 前端依赖安装
cd frontend-admin && npm install
cd frontend-h5 && npm install

# 测试账号修复（如登录失败）
bash backend/scripts/repair_login_and_run_order_regression.sh
```
