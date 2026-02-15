package statistics

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dmh/api/internal/svc"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupStatisticsHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Order{}, &model.Campaign{}, &model.Brand{}, &model.User{}, &model.Distributor{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestStatisticsHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetDashboardStatsHandler(nil))
}

func TestGetDashboardStatsHandler_Success(t *testing.T) {
	db := setupStatisticsHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{Name: "Test Campaign", BrandId: brand.Id, Status: "active"}
	db.Create(campaign)

	user := &model.User{Username: "testuser", Password: "pass", Phone: "13800138000", Status: "active"}
	db.Create(user)

	order := &model.Order{CampaignId: campaign.Id, Phone: "13800138000", Amount: 100.00, PayStatus: "paid", Status: "active"}
	db.Create(order)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDashboardStatsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/statistics/dashboard?brandId=1", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
