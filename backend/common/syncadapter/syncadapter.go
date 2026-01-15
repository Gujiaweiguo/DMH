package syncadapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"        // Oracle 驱动
	_ "github.com/microsoft/go-mssqldb" // SQL Server 驱动
	"github.com/zeromicro/go-zero/core/logx"
)

// ExternalSyncConfig 外部同步配置
type ExternalSyncConfig struct {
	Type     string // mysql, oracle, sqlserver
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Schema   string // SQL Server 专用
	Charset  string
}

// SyncAdapter 外部数据库同步适配器
type SyncAdapter struct {
	db      *sql.DB
	config  ExternalSyncConfig
	mapper  *FieldMapper
	logger  logx.Logger
	metrics *SyncMetrics
}

// NewSyncAdapter 创建同步适配器
func NewSyncAdapter(config ExternalSyncConfig) (*SyncAdapter, error) {
	// 1. 根据数据库类型连接
	db, err := connectDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to external database: %w", err)
	}

	// 2. 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(10 * time.Minute)

	// 3. 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping external database: %w", err)
	}

	// 4. 初始化字段映射器
	mapper := NewFieldMapper()

	// 5. 初始化指标收集
	metrics := NewSyncMetrics()

	return &SyncAdapter{
		db:      db,
		config:  config,
		mapper:  mapper,
		logger:  logx.WithContext(context.Background()),
		metrics: metrics,
	}, nil
}

// connectDatabase 根据类型连接不同的数据库
func connectDatabase(config ExternalSyncConfig) (*sql.DB, error) {
	var dsn string
	var driver string

	switch config.Type {
	case "mysql":
		driver = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
			config.User, config.Password, config.Host, config.Port, config.Database,
			getCharset(config.Charset))

	case "oracle":
		driver = "godror"
		// Oracle DSN 格式: user/password@host:port/service_name
		dsn = fmt.Sprintf("%s/%s@%s:%d/%s",
			config.User, config.Password, config.Host, config.Port, config.Database)

	case "sqlserver":
		driver = "sqlserver"
		// SQL Server DSN 格式: sqlserver://username:password@host:port?database=dbname
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			config.User, config.Password, config.Host, config.Port, config.Database)

	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getCharset(charset string) string {
	if charset == "" {
		return "utf8mb4"
	}
	return charset
}

// SyncOrderData 同步订单数据到外部数据库
type SyncOrderData struct {
	OrderId    int64
	CampaignId int64
	MemberId   int64  // 会员ID
	UnionID    string // 微信 unionid
	Phone      string
	FormData   map[string]interface{}
	Amount     float64
	PayStatus  string
	CreatedAt  time.Time
}

// SyncOrder 同步订单到外部数据库
func (s *SyncAdapter) SyncOrder(ctx context.Context, data *SyncOrderData) error {
	startTime := time.Now()

	// 1. 转换数据格式
	externalData := s.mapper.MapOrder(data)

	// 2. 将FormData转换为JSON字符串
	formDataJSON, err := json.Marshal(data.FormData)
	if err != nil {
		formDataJSON = []byte("{}")
	}

	// 3. 构建INSERT语句（支持幂等性）
	query := `
		INSERT INTO external_orders 
		(order_id, campaign_id, member_id, unionid, phone, form_data, amount, pay_status, created_at, synced_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
		member_id = VALUES(member_id),
		unionid = VALUES(unionid),
		amount = VALUES(amount),
		pay_status = VALUES(pay_status),
		synced_at = NOW()
	`

	// 4. 执行插入（带超时控制）
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = s.db.ExecContext(ctx, query,
		externalData["order_id"],
		externalData["campaign_id"],
		externalData["member_id"],
		externalData["unionid"],
		externalData["phone"],
		string(formDataJSON),
		externalData["amount"],
		externalData["pay_status"],
		externalData["created_at"],
	)

	if err != nil {
		s.metrics.RecordSync("order", false, time.Since(startTime))
		return fmt.Errorf("failed to sync order to external database: %w", err)
	}

	// 5. 记录指标
	s.metrics.RecordSync("order", true, time.Since(startTime))
	s.logger.Infof("order %d synced successfully in %v", data.OrderId, time.Since(startTime))

	return nil
}

// SyncRewardData 同步奖励数据
type SyncRewardData struct {
	RewardId  int64
	UserId    int64
	MemberId  int64 // 会员ID
	OrderId   int64
	Amount    float64
	Status    string
	SettledAt time.Time
}

// SyncReward 同步奖励到外部数据库
func (s *SyncAdapter) SyncReward(ctx context.Context, data *SyncRewardData) error {
	startTime := time.Now()

	// 1. 转换数据格式
	externalData := s.mapper.MapReward(data)

	// 2. 构建INSERT语句
	query := `
		INSERT INTO external_rewards
		(reward_id, user_id, member_id, order_id, amount, status, settled_at, synced_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
		member_id = VALUES(member_id),
		amount = VALUES(amount),
		status = VALUES(status),
		settled_at = VALUES(settled_at),
		synced_at = NOW()
	`

	// 3. 执行插入（带超时控制）
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query,
		externalData["reward_id"],
		externalData["user_id"],
		externalData["member_id"],
		externalData["order_id"],
		externalData["amount"],
		externalData["status"],
		externalData["settled_at"],
	)

	if err != nil {
		s.metrics.RecordSync("reward", false, time.Since(startTime))
		return fmt.Errorf("failed to sync reward to external database: %w", err)
	}

	// 4. 记录指标
	s.metrics.RecordSync("reward", true, time.Since(startTime))
	s.logger.Infof("reward %d synced successfully in %v", data.RewardId, time.Since(startTime))

	return nil
}

// AsyncSyncOrder 异步同步订单（放入队列）
func (s *SyncAdapter) AsyncSyncOrder(data *SyncOrderData) {
	// 这里应该将同步任务放入消息队列（如Redis、RabbitMQ等）
	// MVP阶段可以使用goroutine简单实现
	go func() {
		ctx := context.Background()
		if err := s.SyncOrder(ctx, data); err != nil {
			// 记录日志或重试逻辑
			fmt.Printf("async sync order failed: %v\n", err)
		}
	}()
}

// Close 关闭连接
func (s *SyncAdapter) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// HealthCheck 健康检查
func (s *SyncAdapter) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return s.db.PingContext(ctx)
}

// =============================================================================
// FieldMapper - 字段映射器
// =============================================================================

type FieldMapper struct{}

func NewFieldMapper() *FieldMapper {
	return &FieldMapper{}
}

// MapOrder 映射订单数据
func (m *FieldMapper) MapOrder(data *SyncOrderData) map[string]interface{} {
	return map[string]interface{}{
		"order_id":    data.OrderId,
		"campaign_id": data.CampaignId,
		"member_id":   data.MemberId,
		"unionid":     data.UnionID,
		"phone":       data.Phone,
		"amount":      data.Amount,
		"pay_status":  data.PayStatus,
		"created_at":  data.CreatedAt,
	}
}

// MapReward 映射奖励数据
func (m *FieldMapper) MapReward(data *SyncRewardData) map[string]interface{} {
	return map[string]interface{}{
		"reward_id":  data.RewardId,
		"user_id":    data.UserId,
		"member_id":  data.MemberId,
		"order_id":   data.OrderId,
		"amount":     data.Amount,
		"status":     data.Status,
		"settled_at": data.SettledAt,
	}
}

// =============================================================================
// SyncMetrics - 同步指标收集
// =============================================================================

type SyncMetrics struct {
	TotalSyncs   int64
	SuccessSyncs int64
	FailedSyncs  int64
	TotalTime    time.Duration
}

func NewSyncMetrics() *SyncMetrics {
	return &SyncMetrics{}
}

// RecordSync 记录同步指标
func (m *SyncMetrics) RecordSync(syncType string, success bool, duration time.Duration) {
	m.TotalSyncs++
	m.TotalTime += duration

	if success {
		m.SuccessSyncs++
	} else {
		m.FailedSyncs++
	}
}

// GetStats 获取统计信息
func (m *SyncMetrics) GetStats() map[string]interface{} {
	avgTime := time.Duration(0)
	if m.TotalSyncs > 0 {
		avgTime = m.TotalTime / time.Duration(m.TotalSyncs)
	}

	successRate := float64(0)
	if m.TotalSyncs > 0 {
		successRate = float64(m.SuccessSyncs) / float64(m.TotalSyncs)
	}

	return map[string]interface{}{
		"total_syncs":   m.TotalSyncs,
		"success_syncs": m.SuccessSyncs,
		"failed_syncs":  m.FailedSyncs,
		"avg_time":      avgTime.String(),
		"success_rate":  successRate,
	}
}
