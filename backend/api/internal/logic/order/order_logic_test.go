// Unit tests for order logic
package order

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(4)
	sqlDB.SetMaxIdleConns(4)

	// Auto migrate all models
	err = db.AutoMigrate(
		&model.Campaign{},
		&model.Order{},
		&model.VerificationRecord{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Mirror production duplicate guard for deterministic duplicate-create tests.
	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS uk_orders_campaign_phone ON orders(campaign_id, phone)").Error
	if err != nil {
		t.Fatalf("Failed to create unique index: %v", err)
	}

	return db
}

func createTestCampaign(t *testing.T, db *gorm.DB) *model.Campaign {
	campaign := &model.Campaign{
		Name:        "测试活动",
		Description: "这是一个测试活动",
		FormFields:  `[{"type":"text","name":"name","label":"姓名","required":true,"placeholder":"请输入姓名"}]`,
		RewardRule:  10.00,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(24 * time.Hour),
		Status:      "active",
		BrandId:     1,
	}

	if err := db.Create(campaign).Error; err != nil {
		t.Fatalf("Failed to create test campaign: %v", err)
	}

	return campaign
}

func cleanupTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
}

// CreateOrderLogic tests
func TestCreateOrderLogic_CreateOrder_Success(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	campaign := createTestCampaign(t, db)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateOrderLogic(ctx, svcCtx)

	req := &types.CreateOrderReq{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData: map[string]string{
			"name": "张三",
		},
		ReferrerId: 0,
	}

	resp, err := logic.CreateOrder(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, campaign.Id, resp.CampaignId)
	assert.Equal(t, "13800138000", resp.Phone)
	assert.Equal(t, "pending", resp.Status)
	assert.NotZero(t, resp.Id)
}

func TestCreateOrderLogic_InvalidPhone(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	campaign := createTestCampaign(t, db)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateOrderLogic(ctx, svcCtx)

	req := &types.CreateOrderReq{
		CampaignId: campaign.Id,
		Phone:      "invalid",
		FormData: map[string]string{
			"name": "张三",
		},
	}

	resp, err := logic.CreateOrder(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "手机号格式错误")
}

func TestCreateOrderLogic_DuplicateOrder(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	campaign := createTestCampaign(t, db)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateOrderLogic(ctx, svcCtx)

	req := &types.CreateOrderReq{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData: map[string]string{
			"name": "张三",
		},
	}

	// Create first order
	resp1, err1 := logic.CreateOrder(req)
	assert.NoError(t, err1)
	assert.NotNil(t, resp1)

	// Try to create second order with same phone and campaign
	resp2, err2 := logic.CreateOrder(req)

	assert.Error(t, err2)
	assert.Nil(t, resp2)
	assert.Contains(t, err2.Error(), "重复")
}

// VerifyOrderLogic tests
func TestVerifyOrderLogic_VerifyOrder_Success(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create a test order
	createLogic := NewCreateOrderLogic(context.Background(), &svc.ServiceContext{DB: db})
	orderId := int64(1)
	timestamp := time.Now().Unix()
	verificationCode := createLogic.generateVerificationCode(orderId, "13800138000", timestamp)

	order := &model.Order{
		Id:                 orderId,
		CampaignId:         1,
		Phone:              "13800138000",
		FormData:           `{"name":"张三"}`,
		Status:             "paid",
		PayStatus:          "paid",
		VerificationStatus: "unverified",
		VerificationCode:   verificationCode,
	}
	db.Create(order)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "userId", int64(100))
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewVerifyOrderLogic(ctx, svcCtx)

	req := &types.VerifyOrderReq{
		Code: verificationCode,
	}

	resp, err := logic.VerifyOrder(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, order.Id, resp.OrderId)
	assert.Equal(t, "verified", resp.Status)
}

func TestVerifyOrderLogic_VerifyOrder_AlreadyVerified(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create a verified order
	createLogic := NewCreateOrderLogic(context.Background(), &svc.ServiceContext{DB: db})
	orderId := int64(1)
	timestamp := time.Now().Unix()
	verificationCode := createLogic.generateVerificationCode(orderId, "13800138000", timestamp)

	order := &model.Order{
		Id:                 orderId,
		CampaignId:         1,
		Phone:              "13800138000",
		FormData:           `{"name":"张三"}`,
		Status:             "paid",
		PayStatus:          "paid",
		VerificationStatus: "verified",
		VerificationCode:   verificationCode,
	}
	db.Create(order)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "userId", int64(100))
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewVerifyOrderLogic(ctx, svcCtx)

	req := &types.VerifyOrderReq{
		Code: verificationCode,
	}

	resp, err := logic.VerifyOrder(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "已核销")
}

// Test verification code generation
func TestCreateOrderLogic_GenerateVerificationCode(t *testing.T) {
	logic := &CreateOrderLogic{}

	orderId := int64(1)
	timestamp := int64(1234567890)
	code := logic.generateVerificationCode(orderId, "13800138000", timestamp)

	assert.NotEmpty(t, code)
	assert.Contains(t, code, "1_13800138000_")
	assert.Contains(t, code, "_1234567890_")

	// Parse and verify
	parts := len(code)
	assert.Greater(t, parts, 0)
}
