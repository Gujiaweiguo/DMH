# Spec: 订单与支付系统

**Module**: Order & Payment System  
**Priority**: P0  
**Status**: ✅ Approved  
**Related Proposal**: [001-dmh-mvp-core-features](../changes/001-dmh-mvp-core-features.md)

---

## 📋 模块概述

订单与支付系统是 DMH 的核心交易模块，负责处理用户报名、支付、订单管理和实时奖励结算。该模块需要保证高可用性和数据一致性。

---

## 🎯 核心功能

### 1. 订单创建

#### 业务流程
```
1. 用户提交表单
   ↓
2. 防重复检查（campaign_id + phone）
   ↓
3. 验证活动有效性
   ↓
4. 创建订单记录（状态：pending）
   ↓
5. 返回订单信息
```

#### 防重复逻辑
- 数据库唯一索引：`UNIQUE KEY uk_campaign_phone (campaign_id, phone, deleted_at)`
- 业务逻辑校验：查询是否存在未删除的订单
- 提示信息："您已参与过该活动"

### 2. 支付发起

#### 微信支付集成
```go
1. 生成预支付订单
   - 调用微信统一下单 API
   - 传入：订单号、金额、商品描述、回调URL
   
2. 返回支付参数
   - 小程序：获取 prepay_id
   - H5：获取 mweb_url
   
3. 前端调起支付
   - wx.requestPayment()
```

#### 支付参数配置
```yaml
WeChat:
  AppID: wx1234567890abcdef
  MchID: 1234567890
  APIKey: your_api_key_here
  NotifyURL: https://api.dmh.com/api/v1/orders/payment/callback
```

### 3. 支付回调处理（核心业务逻辑）

#### 回调流程
```
1. 接收微信回调通知
   ↓
2. 验证签名（防篡改）
   ↓
3. 幂等性检查（trade_no）
   ↓
4. 开启数据库事务
   ├─ 更新订单状态为 paid
   ├─ 查询推荐人信息
   ├─ 计算奖励金额
   ├─ 创建奖励记录
   ├─ 更新用户余额（乐观锁）
   └─ 更新订单 sync_status
   ↓
5. 提交事务
   ↓
6. 异步同步到外部数据库
   ↓
7. 返回成功响应给微信
```

#### 关键代码逻辑（Go）
```go
func (l *PaymentCallbackLogic) PaymentCallback(req *types.PaymentCallbackReq) error {
    // 1. 验证签名
    if !l.verifySignature(req) {
        return errors.New("invalid signature")
    }
    
    // 2. 幂等性检查
    if l.isProcessed(req.TradeNo) {
        return nil // 已处理，直接返回成功
    }
    
    // 3. 开启事务
    tx := l.svcCtx.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 4. 更新订单状态
    order, err := l.updateOrderStatus(tx, req.OrderId, "paid", req.TradeNo)
    if err != nil {
        tx.Rollback()
        return err
    }
    
    // 5. 处理推荐奖励
    if order.ReferrerId > 0 {
        // 查询活动奖励规则
        campaign, _ := l.getCampaign(tx, order.CampaignId)
        rewardAmount := campaign.RewardRule
        
        // 创建奖励记录
        reward := &model.Reward{
            UserId:     order.ReferrerId,
            OrderId:    order.Id,
            CampaignId: order.CampaignId,
            Amount:     rewardAmount,
            Status:     "settled",
            SettledAt:  time.Now(),
        }
        if err := l.createReward(tx, reward); err != nil {
            tx.Rollback()
            return err
        }
        
        // 更新用户余额（乐观锁）
        if err := l.updateBalance(tx, order.ReferrerId, rewardAmount); err != nil {
            tx.Rollback()
            return err
        }
    }
    
    // 6. 提交事务
    if err := tx.Commit().Error; err != nil {
        return err
    }
    
    // 7. 异步同步到外部数据库
    if l.svcCtx.SyncAdapter != nil {
        l.svcCtx.SyncAdapter.AsyncSyncOrder(&syncadapter.SyncOrderData{
            OrderId:    order.Id,
            CampaignId: order.CampaignId,
            Phone:      order.Phone,
            Amount:     order.Amount,
            PayStatus:  "paid",
            CreatedAt:  order.CreatedAt,
        })
    }
    
    return nil
}
```

#### 余额更新（乐观锁）
```sql
UPDATE user_balances 
SET 
    balance = balance + ?,
    total_reward = total_reward + ?,
    version = version + 1
WHERE 
    user_id = ? 
    AND version = ?
```

如果更新影响行数为 0，说明版本号冲突，需要重试。

### 4. 订单查询

#### 查询维度
- 按订单 ID 查询
- 按用户手机号查询
- 按活动 ID 查询（管理后台）
- 按推荐人 ID 查询

#### 返回信息
- 订单基本信息
- 活动信息
- 表单数据
- 支付状态
- 同步状态

---

## 🔌 API 接口定义

### 1. 创建订单
```
POST /api/v1/orders
Content-Type: application/json

Request:
{
  "campaignId": 1,
  "phone": "13800138000",
  "formData": {
    "name": "张三",
    "phone": "13800138000",
    "course": "前端开发"
  },
  "referrerId": 100  // 推荐人ID，可选
}

Response:
{
  "id": 12345,
  "campaignId": 1,
  "phone": "13800138000",
  "formData": {...},
  "referrerId": 100,
  "status": "pending",
  "amount": 99.00,
  "createdAt": "2025-01-01T10:00:00Z"
}

Error Response:
{
  "code": 40001,
  "message": "您已参与过该活动"
}
```

### 2. 发起支付
```
POST /api/v1/orders/:id/pay
Content-Type: application/json

Request:
{
  "payType": "wechat",  // wechat | alipay
  "clientType": "h5"    // h5 | miniprogram
}

Response:
{
  "prepayId": "wx20250101100000abcdef",
  "payParams": {
    "appId": "wx1234567890",
    "timeStamp": "1704096000",
    "nonceStr": "abc123",
    "package": "prepay_id=wx20250101100000abcdef",
    "signType": "RSA",
    "paySign": "signature_here"
  }
}
```

### 3. 支付回调
```
POST /api/v1/orders/payment/callback
Content-Type: application/json

Request (微信回调):
{
  "orderId": 12345,
  "payStatus": "paid",
  "amount": 99.00,
  "tradeNo": "4200001234567890",
  "signature": "..."
}

Response:
{
  "code": "SUCCESS",
  "message": "OK"
}
```

### 4. 查询订单
```
GET /api/v1/orders/:id

Response:
{
  "id": 12345,
  "campaignId": 1,
  "campaignName": "新年促销活动",
  "phone": "13800138000",
  "formData": {...},
  "referrerId": 100,
  "status": "paid",
  "amount": 99.00,
  "payStatus": "paid",
  "tradeNo": "4200001234567890",
  "syncStatus": "synced",
  "createdAt": "2025-01-01T10:00:00Z",
  "paidAt": "2025-01-01T10:05:00Z"
}
```

### 5. 查询我的订单
```
GET /api/v1/orders/my?phone=13800138000

Response:
{
  "total": 5,
  "orders": [
    {
      "id": 12345,
      "campaignName": "新年促销活动",
      "status": "paid",
      "amount": 99.00,
      "createdAt": "2025-01-01T10:00:00Z"
    }
  ]
}
```

---

## 💾 数据存储

### orders 表
```sql
CREATE TABLE orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    campaign_id BIGINT NOT NULL COMMENT '活动ID',
    phone VARCHAR(20) NOT NULL COMMENT '用户手机号',
    form_data JSON COMMENT '表单数据',
    referrer_id BIGINT DEFAULT 0 COMMENT '推荐人ID',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '订单状态',
    amount DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '订单金额',
    pay_status VARCHAR(20) NOT NULL DEFAULT 'unpaid' COMMENT '支付状态',
    trade_no VARCHAR(100) DEFAULT '' COMMENT '交易流水号',
    sync_status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '同步状态',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL,
    paid_at DATETIME NULL COMMENT '支付时间',
    
    INDEX idx_campaign_id (campaign_id),
    INDEX idx_phone (phone),
    INDEX idx_referrer_id (referrer_id),
    INDEX idx_status (status),
    INDEX idx_pay_status (pay_status),
    INDEX idx_sync_status (sync_status),
    INDEX idx_trade_no (trade_no),
    UNIQUE KEY uk_campaign_phone (campaign_id, phone, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 状态枚举

#### order status
- `pending` - 待支付
- `paid` - 已支付
- `cancelled` - 已取消
- `refunded` - 已退款

#### pay_status
- `unpaid` - 未支付
- `paid` - 已支付
- `refunded` - 已退款

#### sync_status
- `pending` - 待同步
- `syncing` - 同步中
- `synced` - 已同步
- `failed` - 同步失败

---

## 🔐 安全措施

### 1. 签名验证
```go
func (l *PaymentCallbackLogic) verifySignature(req *types.PaymentCallbackReq) bool {
    // 按微信规定的规则生成签名
    signStr := fmt.Sprintf("orderId=%d&payStatus=%s&amount=%.2f&tradeNo=%s&key=%s",
        req.OrderId, req.PayStatus, req.Amount, req.TradeNo, l.svcCtx.Config.WeChat.APIKey)
    
    expectedSign := md5.Sum([]byte(signStr))
    return req.Signature == hex.EncodeToString(expectedSign[:])
}
```

### 2. 幂等性检查
```go
func (l *PaymentCallbackLogic) isProcessed(tradeNo string) bool {
    var count int64
    l.svcCtx.DB.Model(&model.Order{}).
        Where("trade_no = ? AND pay_status = 'paid'", tradeNo).
        Count(&count)
    return count > 0
}
```

### 3. 防重复报名
- 数据库唯一索引约束
- 创建订单前检查
- 友好的错误提示

### 4. 金额校验
```go
// 校验回调金额与订单金额是否一致
if req.Amount != order.Amount {
    return errors.New("amount mismatch")
}
```

---

## 📊 性能优化

### 1. 数据库优化
- 为常用查询字段添加索引
- 使用预编译语句
- 连接池配置：
  ```yaml
  MaxOpenConns: 100
  MaxIdleConns: 20
  ConnMaxLifetime: 1h
  ```

### 2. 缓存策略
- 活动信息缓存（Redis）
  - Key: `campaign:{id}`
  - TTL: 5 分钟
- 用户订单缓存
  - Key: `user:orders:{phone}`
  - TTL: 1 分钟

### 3. 异步处理
- 外部数据库同步使用异步队列
- 支付回调返回要快（< 500ms）
- 通知推送使用异步任务

---

## ✅ 验收标准

### 功能验收
- [ ] 可以成功创建订单
- [ ] 防重复报名生效
- [ ] 支付流程完整
- [ ] 支付回调正确处理
- [ ] 奖励实时结算
- [ ] 订单查询正常
- [ ] 同步状态更新

### 性能验收
- [ ] 创建订单响应 < 300ms
- [ ] 支付回调处理 < 500ms
- [ ] 奖励结算延迟 < 2s
- [ ] 支持 100 QPS 并发

### 安全验收
- [ ] 签名验证生效
- [ ] 幂等性检查生效
- [ ] 防重复报名生效
- [ ] 金额校验生效

---

## 🧪 测试用例

### 单元测试
1. 创建订单 - 正常流程
2. 创建订单 - 重复报名检查
3. 创建订单 - 活动不存在
4. 创建订单 - 活动已结束
5. 支付回调 - 正常流程
6. 支付回调 - 签名错误
7. 支付回调 - 重复回调
8. 支付回调 - 金额不匹配
9. 奖励结算 - 正常流程
10. 奖励结算 - 余额并发更新

### 集成测试
1. 完整的支付流程测试
2. 推荐奖励完整流程
3. 数据同步完整流程
4. 并发订单创建测试

### 压力测试
1. 100 QPS 订单创建
2. 50 QPS 支付回调
3. 数据库连接池测试
4. 内存占用监控

---

## 📝 开发清单

### 后端开发
- [ ] 实现订单创建逻辑
- [ ] 实现防重复检查
- [ ] 集成微信支付 SDK
- [ ] 实现支付回调处理
- [ ] 实现奖励结算逻辑
- [ ] 实现余额更新（乐观锁）
- [ ] 实现订单查询接口
- [ ] 集成外部同步适配器
- [ ] 编写单元测试
- [ ] 编写集成测试
- [ ] 性能测试

### 前端开发
- [ ] 实现订单提交页面
- [ ] 集成支付 SDK
- [ ] 实现支付成功页面
- [ ] 实现订单查询页面
- [ ] 错误处理
- [ ] 支付状态轮询

---

## 🔗 相关文档
- [Proposal: DMH MVP 核心功能](../changes/001-dmh-mvp-core-features.md)
- [Spec: 实时奖励系统](./003-reward-system.md)
- [Spec: 外网数据同步](./004-sync-adapter.md)
