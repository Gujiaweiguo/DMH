package statistics

import (
	"dmh/api/internal/handler/testutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"dmh/api/internal/svc"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupStatisticsHandlerTestDB(t *testing.T) *gorm.DB {
	db := testutil.SetupGormTestDB(t)

	err := db.AutoMigrate(&model.Order{}, &model.Campaign{}, &model.Brand{}, &model.User{}, &model.Distributor{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestStatisticsHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetDashboardStatsHandler(nil))
}

func TestGetDashboardStatsHandler_InvalidBrandId(t *testing.T) {
	db := setupStatisticsHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDashboardStatsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/statistics/dashboard?brandId=0", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetDashboardStatsHandler_MissingBrandId(t *testing.T) {
	db := setupStatisticsHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDashboardStatsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/statistics/dashboard", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}
