# DMH 活动高级功能回滚方案

## 回滚概述

本文档描述了 DMH 活动高级功能（海报生成、支付配置、订单核销）的回滚方案，用于在部署出现问题时快速恢复到稳定状态。

---

## 回滚触发条件

在以下情况下应执行回滚：
1. 部署后出现严重 Bug 导致核心功能不可用
2. 数据库迁移失败或导致数据不一致
3. 性能严重下降（API 响应时间 > 5 秒）
4. 安全漏洞暴露
5. 用户投诉超过阈值（> 10% 用户受影响）

---

## 快速回滚方案（5 分钟内）

### 1. 服务回滚

**目标：** 恢复到上一个稳定版本

**步骤：**

1. 停止当前服务
   ```bash
   cd deployment
   docker-compose -f docker-compose-dmh.yml down
   ```

2. 恢复后端二进制文件
   ```bash
   cp deployment/backups/{BACKUP_DATE}/bin/dmh-api backend/bin/dmh-api
   ```

3. 恢复前端构建产物
   ```bash
   # H5 前端
   cp -r deployment/backups/{BACKUP_DATE}/h5-dist frontend-h5/dist
   
   # 管理后台前端
   cp -r deployment/backups/{BACKUP_DATE}/admin-dist frontend-admin/dist
   ```

4. 重启服务
   ```bash
   docker-compose -f docker-compose-dmh.yml up -d
   ```

5. 验证服务状态
   ```bash
   # 检查后端健康状态
   curl http://localhost:8889/health
   
   # 检查前端访问
   curl -I http://localhost:3100
   curl -I http://localhost:3000
   ```

**预估时间：** 3-5 分钟

---

### 2. 数据库回滚

**目标：** 回滚数据库变更

**新增表：**
- `poster_template_configs` - 海报模板配置表

**新增字段：**
- `campaigns.payment_config` - 支付配置字段
- `campaigns.poster_template_id` - 海报模板ID字段
- `campaigns.enable_distribution` - 启用分销字段
- `campaigns.distribution_level` - 分销层级字段
- `campaigns.distribution_rewards` - 分销奖励配置字段
- `orders.verification_status` - 核销状态字段
- `orders.verified_at` - 核销时间字段
- `orders.verified_by` - 核销人ID字段
- `orders.verification_code` - 核销码字段

**回滚 SQL 脚本：**

```sql
-- 回滚活动高级功能数据库变更
USE dmh;

-- 1. 删除海报模板配置表
DROP TABLE IF EXISTS poster_template_configs;

-- 2. 删除 campaigns 表的新增字段
ALTER TABLE campaigns 
DROP COLUMN IF EXISTS payment_config,
DROP COLUMN IF EXISTS poster_template_id,
DROP COLUMN IF EXISTS enable_distribution,
DROP COLUMN IF EXISTS distribution_level,
DROP COLUMN IF EXISTS distribution_rewards;

-- 3. 删除 orders 表的新增字段和索引
ALTER TABLE orders
DROP COLUMN IF EXISTS verification_status,
DROP COLUMN IF EXISTS verified_at,
DROP COLUMN IF EXISTS verified_by,
DROP COLUMN IF EXISTS verification_code;

DROP INDEX IF EXISTS idx_verification_status ON orders;
DROP INDEX IF EXISTS idx_verified_at ON orders;

SELECT '数据库回滚完成！' as status, 
       NOW() as rollback_time;
```

**执行回滚：**
```bash
docker exec -i mysql8 mysql -uroot -p'#Admin168' dmh < deployment/scripts/rollback-advanced-features.sql
```

**预估时间：** 1-2 分钟

---

## 回滚验证清单

### 后端验证

- [ ] 后端服务正常启动
- [ ] 后端健康检查通过
- [ ] 登录接口正常
- [ ] 活动列表接口正常
- [ ] 订单创建接口正常
- [ ] 核销接口可以正常访问

### 前端验证

- [ ] H5 前端可以正常访问
- [ ] 管理后台前端可以正常访问
- [ ] 页面加载正常
- [ ] 登录功能正常
- [ ] 活动管理页面正常

### 数据库验证

- [ ] poster_template_configs 表已删除
- [ ] campaigns 表新增字段已删除
- [ ] orders 表新增字段已删除
- [ ] 现有数据完整性正常

---

## 回滚后通知事项

1. 立即通知开发团队
2. 通知用户
3. 监控日志
4. 记录回滚信息
