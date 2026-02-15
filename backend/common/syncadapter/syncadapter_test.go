package syncadapter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SyncAdapterTestSuite struct {
	suite.Suite
}

func (suite *SyncAdapterTestSuite) TestNewFieldMapper() {
	mapper := NewFieldMapper()
	assert.NotNil(suite.T(), mapper)
}

func (suite *SyncAdapterTestSuite) TestNewSyncMetrics() {
	metrics := NewSyncMetrics()
	assert.NotNil(suite.T(), metrics)
}

func (suite *SyncAdapterTestSuite) TestExternalSyncConfig() {
	config := ExternalSyncConfig{
		Type:     "mysql",
		Host:     "localhost",
		Port:     3306,
		User:     "test",
		Password: "password",
		Database: "testdb",
		Charset:  "utf8mb4",
	}

	assert.Equal(suite.T(), "mysql", config.Type)
	assert.Equal(suite.T(), "localhost", config.Host)
	assert.Equal(suite.T(), 3306, config.Port)
}

func (suite *SyncAdapterTestSuite) TestSyncTask() {
	task := SyncTask{
		TaskId:   "test-task-001",
		Type:     "order",
		OrderId:  123,
		Attempts: 0,
	}

	assert.Equal(suite.T(), "test-task-001", task.TaskId)
	assert.Equal(suite.T(), "order", task.Type)
	assert.Equal(suite.T(), int64(123), task.OrderId)
}

func TestSyncAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(SyncAdapterTestSuite))
}

func TestFieldMapper_MapOrder(t *testing.T) {
	mapper := NewFieldMapper()
	now := time.Now()

	data := &SyncOrderData{
		OrderId:    123,
		CampaignId: 456,
		MemberId:   789,
		UnionID:    "wx_unionid",
		Phone:      "13800138000",
		FormData:   map[string]interface{}{"name": "test"},
		Amount:     100.00,
		PayStatus:  "paid",
		CreatedAt:  now,
	}

	result := mapper.MapOrder(data)

	assert.Equal(t, int64(123), result["order_id"])
	assert.Equal(t, int64(456), result["campaign_id"])
	assert.Equal(t, int64(789), result["member_id"])
	assert.Equal(t, "wx_unionid", result["unionid"])
	assert.Equal(t, "13800138000", result["phone"])
	assert.Equal(t, 100.00, result["amount"])
	assert.Equal(t, "paid", result["pay_status"])
	assert.Equal(t, now, result["created_at"])
}

func TestFieldMapper_MapReward(t *testing.T) {
	mapper := NewFieldMapper()
	now := time.Now()

	data := &SyncRewardData{
		RewardId:  123,
		UserId:    456,
		MemberId:  789,
		OrderId:   101112,
		Amount:    50.00,
		Status:    "settled",
		SettledAt: now,
	}

	result := mapper.MapReward(data)

	assert.Equal(t, int64(123), result["reward_id"])
	assert.Equal(t, int64(456), result["user_id"])
	assert.Equal(t, int64(789), result["member_id"])
	assert.Equal(t, int64(101112), result["order_id"])
	assert.Equal(t, 50.00, result["amount"])
	assert.Equal(t, "settled", result["status"])
	assert.Equal(t, now, result["settled_at"])
}

func TestSyncMetrics_RecordSync_Success(t *testing.T) {
	metrics := NewSyncMetrics()

	metrics.RecordSync("order", true, 100*time.Millisecond)

	assert.Equal(t, int64(1), metrics.TotalSyncs)
	assert.Equal(t, int64(1), metrics.SuccessSyncs)
	assert.Equal(t, int64(0), metrics.FailedSyncs)
	assert.Equal(t, 100*time.Millisecond, metrics.TotalTime)
}

func TestSyncMetrics_RecordSync_Failure(t *testing.T) {
	metrics := NewSyncMetrics()

	metrics.RecordSync("reward", false, 50*time.Millisecond)

	assert.Equal(t, int64(1), metrics.TotalSyncs)
	assert.Equal(t, int64(0), metrics.SuccessSyncs)
	assert.Equal(t, int64(1), metrics.FailedSyncs)
	assert.Equal(t, 50*time.Millisecond, metrics.TotalTime)
}

func TestSyncMetrics_GetStats(t *testing.T) {
	metrics := NewSyncMetrics()

	metrics.RecordSync("order", true, 100*time.Millisecond)
	metrics.RecordSync("order", false, 50*time.Millisecond)
	metrics.RecordSync("reward", true, 200*time.Millisecond)

	stats := metrics.GetStats()

	assert.Equal(t, int64(3), stats["total_syncs"])
	assert.Equal(t, int64(2), stats["success_syncs"])
	assert.Equal(t, int64(1), stats["failed_syncs"])
	assert.Equal(t, 0.6666666666666666, stats["success_rate"])
	assert.Equal(t, "116.666666ms", stats["avg_time"])
}

func TestGetCharset(t *testing.T) {
	assert.Equal(t, "utf8mb4", getCharset(""))
	assert.Equal(t, "utf8", getCharset("utf8"))
	assert.Equal(t, "latin1", getCharset("latin1"))
}

func TestConnectDatabase_UnsupportedType(t *testing.T) {
	config := ExternalSyncConfig{
		Type: "unsupported",
	}

	db, err := connectDatabase(config)

	assert.Nil(t, db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database type")
}

func TestConnectDatabase_MySQL(t *testing.T) {
	config := ExternalSyncConfig{
		Type:     "mysql",
		Host:     "localhost",
		Port:     3306,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
		Charset:  "utf8mb4",
	}

	db, err := connectDatabase(config)

	assert.NotNil(t, db)
	assert.NoError(t, err)
	if db != nil {
		db.Close()
	}
}

func TestConnectDatabase_Oracle(t *testing.T) {
	config := ExternalSyncConfig{
		Type:     "oracle",
		Host:     "localhost",
		Port:     1521,
		User:     "testuser",
		Password: "testpass",
		Database: "ORCL",
	}

	db, err := connectDatabase(config)

	assert.NotNil(t, db)
	assert.NoError(t, err)
	if db != nil {
		db.Close()
	}
}

func TestConnectDatabase_SQLServer(t *testing.T) {
	config := ExternalSyncConfig{
		Type:     "sqlserver",
		Host:     "localhost",
		Port:     1433,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
	}

	db, err := connectDatabase(config)

	assert.NotNil(t, db)
	assert.NoError(t, err)
	if db != nil {
		db.Close()
	}
}

func TestSyncAdapter_Close(t *testing.T) {
	config := ExternalSyncConfig{
		Type:     "mysql",
		Host:     "localhost",
		Port:     3306,
		User:     "test",
		Password: "test",
		Database: "test",
	}

	adapter := &SyncAdapter{
		config: config,
	}

	err := adapter.Close()
	assert.NoError(t, err)
}

func TestSyncAdapter_Close_Nil(t *testing.T) {
	adapter := &SyncAdapter{
		db: nil,
	}

	err := adapter.Close()
	assert.NoError(t, err)
}

func TestSyncMetrics_MultipleRecords(t *testing.T) {
	metrics := NewSyncMetrics()

	for i := 0; i < 100; i++ {
		metrics.RecordSync("order", i%2 == 0, 10*time.Millisecond)
	}

	stats := metrics.GetStats()

	assert.Equal(t, int64(100), stats["total_syncs"])
	assert.Equal(t, int64(50), stats["success_syncs"])
	assert.Equal(t, int64(50), stats["failed_syncs"])
	assert.Equal(t, 0.5, stats["success_rate"])
}

func TestSyncMetrics_ZeroSyncs(t *testing.T) {
	metrics := NewSyncMetrics()

	stats := metrics.GetStats()

	assert.Equal(t, int64(0), stats["total_syncs"])
	assert.Equal(t, int64(0), stats["success_syncs"])
	assert.Equal(t, int64(0), stats["failed_syncs"])
	assert.Equal(t, 0.0, stats["success_rate"])
}

func TestFieldMapper_MapOrder_NilFormData(t *testing.T) {
	mapper := NewFieldMapper()

	data := &SyncOrderData{
		OrderId:    123,
		CampaignId: 456,
		MemberId:   0,
		UnionID:    "",
		Phone:      "",
		FormData:   nil,
		Amount:     0,
		PayStatus:  "",
		CreatedAt:  time.Time{},
	}

	result := mapper.MapOrder(data)

	assert.Equal(t, int64(0), result["member_id"])
	assert.Equal(t, "", result["unionid"])
}

func TestSyncMetrics_EmptyStats(t *testing.T) {
	metrics := &SyncMetrics{}

	stats := metrics.GetStats()

	assert.Equal(t, int64(0), stats["total_syncs"])
}
