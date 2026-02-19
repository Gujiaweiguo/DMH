# Spec: 实时奖励系统

**Module**: Reward System  
**Priority**: P0  
**Status**: ✅ Approved  
**Related Proposal**: [001-dmh-mvp-core-features](../changes/001-dmh-mvp-core-features.md)

---

## 📋 模块概述

实时奖励系统负责在用户完成支付后，自动计算并结算推荐奖励。核心目标是在 2 秒内完成奖励结算，并确保数据一致性和高并发场景下的正确性。

---

## 🎯 核心功能

### 1. 奖励计算

#### 计算规则（MVP 阶段）
```
奖励金额 = 活动奖励规则（固定金额）
```

#### 触发条件
- 推荐用户完成支付
- 订单状态为 `paid`
- 推荐人 ID 有效（referrer_id > 0）

#### 计算流程
```
1. 支付回调成功
   ↓
2. 查询订单推荐人ID
   ↓
3. 查询活动奖励规则
   ↓
4. 计算奖励金额
   ↓
5. 在事务中执行：
   ├─ 创建奖励记录
   ├─ 更新用户余额
   └─ 记录结算时间
   ↓
6. 提交事务（确保原子性）
```

### 2. 余额管理

#### 余额字段
- `balance` - 当前可用余额
- `total_reward` - 累计获得奖励
- `version` - 乐观锁版本号

#### 余额更新（乐观锁）
```sql
-- 第1步：查询当前余额和版本号
SELECT balance, total_reward, version 
FROM user_balances 
WHERE user_id = ?;

-- 第2步：计算新余额
newBalance = balance + rewardAmount
newTotalReward = total_reward + rewardAmount

-- 第3步：更新（带版本号检查）
UPDATE user_balances 
SET 
    balance = ?,
    total_reward = ?,
    version = version + 1,
    updated_at = NOW()
WHERE 
    user_id = ? 
    AND version = ?;

-- 第4步：检查影响行数
-- 如果影响行数 = 0，说明版本冲突，需要重试
```

#### 并发控制
```go
func (l *RewardLogic) UpdateBalance(userId int64, amount float64) error {
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        // 查询当前版本
        balance, err := l.getUserBalance(userId)
        if err != nil {
            return err
        }
        
        // 尝试更新
        result := l.svcCtx.DB.Model(&model.UserBalance{}).
            Where("user_id = ? AND version = ?", userId, balance.Version).
            Updates(map[string]interface{}{
                "balance":      balance.Balance + amount,
                "total_reward": balance.TotalReward + amount,
                "version":      balance.Version + 1,
            })
        
        if result.RowsAffected > 0 {
            return nil // 更新成功
        }
        
        // 版本冲突，等待后重试
        time.Sleep(time.Millisecond * 10)
    }
    
    return errors.New("update balance failed after retries")
}
```

### 3. 奖励记录

#### 记录内容
- 奖励 ID
- 用户 ID（推荐人）
- 订单 ID（被推荐人的订单）
- 活动 ID
- 奖励金额
- 状态（pending/settled/cancelled）
- 结算时间

#### 状态流转
```
pending → settled    // 正常结算
pending → cancelled  // 订单取消/退款
```

### 4. 奖励查询

#### 查询维度
- 按用户 ID 查询奖励列表
- 按订单 ID 查询奖励
- 按活动 ID 统计奖励
- 按时间范围统计

#### 统计信息
- 总奖励金额
- 已结算奖励
- 待结算奖励
- 奖励笔数

---

## 🔌 API 接口定义

### 1. 查询用户余额
```
GET /api/v1/rewards/balance/:userId

Response:
{
  "userId": 100,
  "balance": 156.50,
  "totalReward": 200.00,
  "updatedAt": "2025-01-01T15:30:00Z"
}
```

### 2. 查询奖励列表
```
GET /api/v1/rewards/:userId?page=1&pageSize=20

Response:
{
  "total": 10,
  "rewards": [
    {
      "id": 1,
      "userId": 100,
      "orderId": 12345,
      "campaignId": 1,
      "campaignName": "新年促销活动",
      "amount": 10.00,
      "status": "settled",
      "settledAt": "2025-01-01T10:05:02Z",
      "createdAt": "2025-01-01T10:05:00Z"
    }
  ]
}
```

### 3. 查询奖励详情
```
GET /api/v1/rewards/detail/:rewardId

Response:
{
  "id": 1,
  "userId": 100,
  "orderId": 12345,
  "orderPhone": "13800138000",
  "campaignId": 1,
  "campaignName": "新年促销活动",
  "amount": 10.00,
  "status": "settled",
  "settledAt": "2025-01-01T10:05:02Z",
  "createdAt": "2025-01-01T10:05:00Z"
}
```

### 4. 奖励统计
```
GET /api/v1/rewards/statistics/:userId?startDate=2025-01-01&endDate=2025-01-31

Response:
{
  "userId": 100,
  "totalAmount": 200.00,
  "settledAmount": 156.50,
  "pendingAmount": 43.50,
  "rewardCount": 10,
  "period": {
    "start": "2025-01-01",
    "end": "2025-01-31"
  }
}
```

---

## 💾 数据存储

### rewards 表
```sql
CREATE TABLE rewards (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID（推荐人）',
    order_id BIGINT NOT NULL COMMENT '关联订单ID',
    campaign_id BIGINT NOT NULL COMMENT '活动ID',
    amount DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '奖励金额',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '奖励状态',
    settled_at DATETIME NULL COMMENT '结算时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_user_id (user_id),
    INDEX idx_order_id (order_id),
    INDEX idx_campaign_id (campaign_id),
    INDEX idx_status (status),
    INDEX idx_settled_at (settled_at),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### user_balances 表
```sql
CREATE TABLE user_balances (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    balance DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '当前余额',
    total_reward DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '累计奖励',
    version BIGINT NOT NULL DEFAULT 0 COMMENT '版本号（乐观锁）',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_user_id (user_id),
    INDEX idx_balance (balance)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## ⚡ 性能优化

### 1. 2秒结算目标

#### 性能分析
```
支付回调接收       : 50ms
数据库事务处理     : 200ms
├─ 更新订单状态    : 50ms
├─ 查询推荐人      : 30ms
├─ 创建奖励记录    : 40ms
└─ 更新余额        : 80ms
异步同步入队       : 20ms
返回响应           : 10ms
------------------------
总计               : 280ms (< 2000ms ✓)
```

#### 优化措施
1. **数据库索引优化**
   - 为所有查询字段添加索引
   - 使用复合索引优化联表查询

2. **事务优化**
   - 最小化事务范围
   - 避免在事务中执行外部调用
   - 使用批量操作

3. **连接池配置**
   ```yaml
   MaxOpenConns: 100
   MaxIdleConns: 20
   ConnMaxLifetime: 1h
   ```

### 2. 并发控制

#### 乐观锁重试策略
```go
retryConfig := RetryConfig{
    MaxRetries: 3,
    InitialDelay: 10 * time.Millisecond,
    MaxDelay: 100 * time.Millisecond,
    Multiplier: 2.0,
}
```

#### 性能监控
```go
// 记录奖励结算耗时
startTime := time.Now()
err := l.settleReward(order)
duration := time.Since(startTime)

// 告警：如果超过2秒
if duration > 2*time.Second {
    l.Logger.Warnf("Reward settlement too slow: %v", duration)
}
```

---

## 🔐 安全措施

### 1. 数据一致性
- 使用数据库事务保证原子性
- 乐观锁防止并发问题
- 记录操作日志便于审计

### 2. 防作弊
- 同一活动同一手机号只能报名一次
- 推荐关系验证（referrer_id 有效性）
- 订单金额校验

### 3. 异常处理
```go
defer func() {
    if r := recover(); r != nil {
        l.Logger.Errorf("Reward settlement panic: %v", r)
        // 回滚事务
        tx.Rollback()
        // 记录错误日志
        l.logError(order.Id, fmt.Sprintf("panic: %v", r))
    }
}()
```

---

## ✅ 验收标准

### 功能验收
- [ ] 支付成功后自动创建奖励记录
- [ ] 余额正确更新
- [ ] 奖励列表正确显示
- [ ] 统计数据准确

### 性能验收
- [ ] 奖励结算延迟 < 2 秒
- [ ] 支持 100 QPS 并发
- [ ] 乐观锁重试成功率 > 99%

### 安全验收
- [ ] 并发场景余额正确
- [ ] 不会重复结算奖励
- [ ] 异常情况数据不丢失

---

## 🧪 测试用例

### 单元测试
1. 计算奖励金额 - 正常流程
2. 创建奖励记录 - 正常流程
3. 更新余额 - 正常流程
4. 更新余额 - 乐观锁冲突重试
5. 查询余额 - 正常流程
6. 查询奖励列表 - 分页
7. 查询奖励统计 - 按时间范围

### 并发测试
1. 100 个并发奖励结算
2. 同一用户多笔奖励同时结算
3. 乐观锁重试验证
4. 数据库连接池压力测试

### 性能测试
1. 单笔奖励结算耗时
2. 批量奖励结算吞吐量
3. 查询接口响应时间
4. 数据库CPU和内存占用

---

## 📝 开发清单

### 后端开发
- [ ] 实现奖励计算逻辑
- [ ] 实现余额更新（乐观锁）
- [ ] 实现奖励记录创建
- [ ] 实现查询余额接口
- [ ] 实现查询奖励列表接口
- [ ] 实现奖励统计接口
- [ ] 性能优化
- [ ] 编写单元测试
- [ ] 编写并发测试
- [ ] 性能压测

### 前端开发
- [ ] 我的余额页面
- [ ] 奖励记录列表页
- [ ] 奖励统计图表
- [ ] 实时余额刷新

---

## 🔗 相关文档
- [Proposal: DMH MVP 核心功能](../changes/001-dmh-mvp-core-features.md)
- [Spec: 订单与支付系统](./002-order-payment-system.md)
