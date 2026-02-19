# Spec: 外网数据同步适配器

**Module**: External Sync Adapter  
**Priority**: P1  
**Status**: ✅ Approved  
**Related Proposal**: [001-dmh-mvp-core-features](../changes/001-dmh-mvp-core-features.md)

---

## 📋 模块概述

外网数据同步适配器负责将 DMH 系统的订单和奖励数据实时同步到客户既有的外部数据库系统（Oracle/SQL Server），替代传统 ESB 方案，实现低延迟的数据集成。

---

## 🎯 核心功能

### 1. 数据库连接管理

#### 支持的数据库类型
- Oracle 11g+
- SQL Server 2012+
- MySQL 5.7+ （可选）

#### 连接配置
```yaml
ExternalSync:
  Enabled: true
  Database:
    Type: oracle  # oracle | sqlserver | mysql
    Host: external-db.example.com
    Port: 1521
    User: sync_user
    Password: ${ENCRYPTED_PASSWORD}  # 加密存储
    Database: external_dmh
    Schema: dbo  # SQL Server 专用
    Charset: utf8mb4
```

#### 连接池配置
```go
db.SetMaxOpenConns(10)        // 最大连接数
db.SetMaxIdleConns(5)         // 最大空闲连接
db.SetConnMaxLifetime(1*time.Hour)  // 连接最大生命周期
db.SetConnMaxIdleTime(10*time.Minute)  // 空闲连接超时
```

#### 健康检查
```go
func (s *SyncAdapter) HealthCheck() error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    if err := s.db.PingContext(ctx); err != nil {
        return fmt.Errorf("external database unhealthy: %w", err)
    }
    
    return nil
}
```

### 2. 数据同步逻辑

#### 同步时机
```
1. 支付成功后立即同步订单数据
2. 奖励结算后立即同步奖励数据
3. 同步失败后自动重试（3次）
4. 支持手动重试
```

#### 异步同步流程
```
1. 支付回调处理完成
   ↓
2. 将同步任务放入队列
   ├─ Redis List / RabbitMQ / Kafka
   └─ 任务内容：OrderId, Type, Data
   ↓
3. 后台Worker消费队列
   ↓
4. 执行数据同步
   ├─ 查询DMH数据
   ├─ 转换数据格式
   ├─ 执行INSERT语句
   └─ 更新同步状态
   ↓
5. 失败处理
   ├─ 记录错误日志
   ├─ 更新sync_status为failed
   └─ 支持手动重试
```

### 3. 字段映射规则

#### 订单数据映射
```yaml
Mapping:
  Orders:
    # DMH字段 → 外部系统字段
    - source: id
      target: order_id
      type: bigint
    
    - source: phone
      target: student_phone
      type: varchar(20)
    
    - source: form_data.name
      target: student_name
      type: varchar(100)
      extract: json  # JSON字段提取
    
    - source: form_data.course
      target: course_name
      type: varchar(100)
      extract: json
    
    - source: amount
      target: order_amount
      type: decimal(10,2)
    
    - source: created_at
      target: register_time
      type: datetime
    
    - source: pay_status
      target: payment_status
      type: varchar(20)
```

#### 奖励数据映射
```yaml
Mapping:
  Rewards:
    - source: id
      target: reward_id
      type: bigint
    
    - source: user_id
      target: referrer_id
      type: bigint
    
    - source: order_id
      target: order_id
      type: bigint
    
    - source: amount
      target: reward_amount
      type: decimal(10,2)
    
    - source: settled_at
      target: settlement_time
      type: datetime
```

### 4. 同步状态监控

#### 状态枚举
- `pending` - 待同步
- `syncing` - 同步中
- `synced` - 已同步
- `failed` - 同步失败

#### 监控指标
- 待同步订单数
- 同步成功率
- 同步平均耗时
- 失败订单列表
- 网络连接状态

---

## 🔌 API 接口定义

### 1. 查询同步状态
```
GET /api/v1/sync/status/:orderId

Response:
{
  "orderId": 12345,
  "syncStatus": "synced",
  "syncTime": "2025-01-01T10:05:30Z",
  "attempts": 1,
  "errorMsg": null
}
```

### 2. 手动触发同步
```
POST /api/v1/sync/retry/:orderId

Response:
{
  "orderId": 12345,
  "message": "同步任务已加入队列",
  "taskId": "task_abc123"
}
```

### 3. 批量查询同步状态
```
GET /api/v1/sync/status/batch?orderIds=12345,12346,12347

Response:
{
  "results": [
    {
      "orderId": 12345,
      "syncStatus": "synced",
      "syncTime": "2025-01-01T10:05:30Z"
    }
  ]
}
```

### 4. 同步统计
```
GET /api/v1/sync/statistics?startDate=2025-01-01&endDate=2025-01-31

Response:
{
  "totalOrders": 1000,
  "syncedOrders": 980,
  "failedOrders": 20,
  "successRate": 0.98,
  "avgSyncTime": "1.2s",
  "period": {
    "start": "2025-01-01",
    "end": "2025-01-31"
  }
}
```

### 5. 健康检查
```
GET /api/v1/sync/health

Response:
{
  "status": "healthy",
  "database": {
    "connected": true,
    "type": "oracle",
    "host": "external-db.example.com",
    "latency": "15ms"
  },
  "queue": {
    "pending": 5,
    "processing": 2
  }
}
```

---

## 💾 数据存储

### sync_logs 表
```sql
CREATE TABLE sync_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL COMMENT '订单ID',
    sync_type VARCHAR(20) NOT NULL COMMENT '同步类型: order/reward',
    sync_status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '同步状态',
    attempts INT NOT NULL DEFAULT 0 COMMENT '尝试次数',
    error_msg TEXT COMMENT '错误信息',
    synced_at DATETIME NULL COMMENT '同步成功时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_order_id (order_id),
    INDEX idx_sync_status (sync_status),
    INDEX idx_sync_type (sync_type),
    INDEX idx_synced_at (synced_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## 🔧 核心实现

### SyncAdapter 结构
```go
type SyncAdapter struct {
    db      *sql.DB
    config  ExternalSyncConfig
    mapper  *FieldMapper
    logger  *logx.Logger
    metrics *SyncMetrics
}

func NewSyncAdapter(config ExternalSyncConfig) (*SyncAdapter, error) {
    // 1. 连接外部数据库
    db, err := connectDatabase(config)
    if err != nil {
        return nil, err
    }
    
    // 2. 初始化字段映射器
    mapper := NewFieldMapper(config.Mapping)
    
    // 3. 初始化指标收集
    metrics := NewSyncMetrics()
    
    return &SyncAdapter{
        db:      db,
        config:  config,
        mapper:  mapper,
        metrics: metrics,
    }, nil
}
```

### 订单同步实现
```go
func (s *SyncAdapter) SyncOrder(ctx context.Context, orderId int64) error {
    // 1. 查询DMH订单数据
    order, err := s.getOrderFromDMH(orderId)
    if err != nil {
        return fmt.Errorf("get order failed: %w", err)
    }
    
    // 2. 转换数据格式
    externalData := s.mapper.MapOrder(order)
    
    // 3. 构建INSERT语句
    query := s.buildInsertQuery("external_orders", externalData)
    
    // 4. 执行插入（带超时控制）
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    _, err = s.db.ExecContext(ctx, query, externalData.Values()...)
    if err != nil {
        return fmt.Errorf("insert to external db failed: %w", err)
    }
    
    // 5. 记录指标
    s.metrics.RecordSync("order", true, time.Since(startTime))
    
    return nil
}
```

### 字段映射器
```go
type FieldMapper struct {
    orderMapping  []FieldMapping
    rewardMapping []FieldMapping
}

type FieldMapping struct {
    Source string        // DMH字段
    Target string        // 外部系统字段
    Type   string        // 数据类型
    Extract string       // JSON提取规则
    Transform func(interface{}) interface{}  // 数据转换函数
}

func (m *FieldMapper) MapOrder(order *model.Order) map[string]interface{} {
    result := make(map[string]interface{})
    
    for _, mapping := range m.orderMapping {
        value := m.extractValue(order, mapping.Source, mapping.Extract)
        if mapping.Transform != nil {
            value = mapping.Transform(value)
        }
        result[mapping.Target] = value
    }
    
    return result
}
```

### 异步同步队列
```go
type SyncQueue struct {
    redis *redis.Client
    key   string
}

func (q *SyncQueue) Enqueue(task *SyncTask) error {
    data, _ := json.Marshal(task)
    return q.redis.RPush(context.Background(), q.key, data).Err()
}

func (q *SyncQueue) Dequeue() (*SyncTask, error) {
    result, err := q.redis.BLPop(context.Background(), 0, q.key).Result()
    if err != nil {
        return nil, err
    }
    
    var task SyncTask
    json.Unmarshal([]byte(result[1]), &task)
    return &task, nil
}
```

### Worker 消费者
```go
func (w *SyncWorker) Run() {
    for {
        // 1. 从队列取任务
        task, err := w.queue.Dequeue()
        if err != nil {
            w.logger.Error("dequeue failed:", err)
            time.Sleep(time.Second)
            continue
        }
        
        // 2. 执行同步
        err = w.syncAdapter.SyncOrder(context.Background(), task.OrderId)
        
        // 3. 更新状态
        if err != nil {
            w.updateSyncStatus(task.OrderId, "failed", err.Error())
            
            // 重试逻辑
            if task.Attempts < 3 {
                task.Attempts++
                w.queue.Enqueue(task)  // 重新入队
            }
        } else {
            w.updateSyncStatus(task.OrderId, "synced", "")
        }
    }
}
```

---

## 🔐 安全措施

### 1. 密码加密存储
```go
// 使用AES加密存储数据库密码
encryptedPassword := encrypt(config.Password, secretKey)
config.Password = encryptedPassword

// 使用时解密
realPassword := decrypt(config.Password, secretKey)
```

### 2. IP白名单
```sql
-- 在外部数据库配置DMH服务器IP白名单
-- Oracle示例
CREATE USER sync_user IDENTIFIED BY password;
GRANT CONNECT, RESOURCE TO sync_user;
-- 限制IP访问（需网络层配置）
```

### 3. 只读权限（可选）
```sql
-- 仅授予INSERT权限，防止误操作
GRANT INSERT ON external_orders TO sync_user;
GRANT INSERT ON external_rewards TO sync_user;
```

### 4. SQL注入防护
```go
// 使用预编译语句
stmt, err := db.Prepare("INSERT INTO external_orders (order_id, phone, name) VALUES (?, ?, ?)")
defer stmt.Close()
_, err = stmt.ExecContext(ctx, orderId, phone, name)
```

---

## ⚡ 性能优化

### 1. 批量同步
```go
// 批量插入优化
func (s *SyncAdapter) SyncOrdersBatch(orderIds []int64) error {
    // 构建批量INSERT语句
    query := "INSERT INTO external_orders VALUES "
    values := []interface{}{}
    
    for i, orderId := range orderIds {
        order := s.getOrder(orderId)
        query += "(?, ?, ?)"
        if i < len(orderIds)-1 {
            query += ", "
        }
        values = append(values, order.Id, order.Phone, order.Amount)
    }
    
    _, err := s.db.Exec(query, values...)
    return err
}
```

### 2. 连接复用
- 使用连接池避免频繁建立连接
- 设置合理的空闲连接数
- 定期清理过期连接

### 3. 异步处理
- 支付回调不等待同步完成
- 使用消息队列解耦
- Worker数量可配置

### 4. 监控告警
```go
// 同步延迟告警
if pendingCount > 100 {
    alert("同步任务堆积超过100条")
}

// 同步成功率告警
if successRate < 0.95 {
    alert("同步成功率低于95%")
}
```

---

## ✅ 验收标准

### 功能验收
- [ ] 支持 Oracle 数据库连接
- [ ] 支持 SQL Server 数据库连接
- [ ] 订单数据正确同步
- [ ] 奖励数据正确同步
- [ ] 字段映射正确
- [ ] 同步失败自动重试
- [ ] 手动重试功能正常

### 性能验收
- [ ] 单条同步耗时 < 1 秒
- [ ] 同步延迟 < 1 分钟
- [ ] 同步成功率 > 95%
- [ ] 支持 50 QPS 同步

### 安全验收
- [ ] 密码加密存储
- [ ] 防SQL注入
- [ ] 连接超时控制
- [ ] 错误日志记录

---

## 🧪 测试用例

### 单元测试
1. 连接外部数据库 - Oracle
2. 连接外部数据库 - SQL Server
3. 字段映射 - 订单数据
4. 字段映射 - 奖励数据
5. JSON字段提取
6. 同步单条订单
7. 同步失败重试
8. 健康检查

### 集成测试
1. 完整的同步流程
2. 异步队列测试
3. Worker消费测试
4. 批量同步测试
5. 网络异常恢复测试

---

## 📝 开发清单

### 后端开发
- [ ] 实现 SyncAdapter 核心类
- [ ] 实现 Oracle 驱动集成
- [ ] 实现 SQL Server 驱动集成
- [ ] 实现字段映射器
- [ ] 实现异步队列
- [ ] 实现 Worker 消费者
- [ ] 实现重试机制
- [ ] 实现同步状态API
- [ ] 实现监控指标
- [ ] 编写单元测试
- [ ] 编写集成测试

### 运维配置
- [ ] 外部数据库账号配置
- [ ] IP白名单配置
- [ ] 监控告警配置
- [ ] 日志收集配置

---

## 🔗 相关文档
- [Proposal: DMH MVP 核心功能](../changes/001-dmh-mvp-core-features.md)
- [Spec: 订单与支付系统](./002-order-payment-system.md)
