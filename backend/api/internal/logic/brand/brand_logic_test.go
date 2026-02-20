package brand

import (
	"context"
	"testing"

	"dmh/api/internal/svc"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupBrandTestDB(t *testing.T) *gorm.DB {
	dsn := "root:Admin168@tcp(127.0.0.1:3306)/dmh_test?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Brand{}, &model.BrandAsset{}, &model.Campaign{}, &model.Order{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func createTestBrand(t *testing.T, db *gorm.DB, name, status string) *model.Brand {
	brand := &model.Brand{
		Name:        name,
		Description: "Test brand description",
		Status:      status,
	}
	if err := db.Create(brand).Error; err != nil {
		t.Fatalf("Failed to create test brand: %v", err)
	}
	return brand
}

func TestGetBrandsLogic_GetBrands_Success(t *testing.T) {
	db := setupBrandTestDB(t)

	createTestBrand(t, db, "Brand 1", "active")
	createTestBrand(t, db, "Brand 2", "active")
	createTestBrand(t, db, "Brand 3", "disabled")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandsLogic(ctx, svcCtx)

	resp, err := logic.GetBrands()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(3), resp.Total)
	assert.Len(t, resp.Brands, 3)
}

func TestGetBrandsLogic_ReturnsCorrectData(t *testing.T) {
	db := setupBrandTestDB(t)

	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandsLogic(ctx, svcCtx)

	resp, err := logic.GetBrands()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Brands, 1)
	assert.Equal(t, brand.Id, resp.Brands[0].Id)
	assert.Equal(t, brand.Name, resp.Brands[0].Name)
	assert.Equal(t, brand.Status, resp.Brands[0].Status)
}

func TestGetBrandsLogic_EmptyResult(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandsLogic(ctx, svcCtx)

	resp, err := logic.GetBrands()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.Total)
	assert.Len(t, resp.Brands, 0)
}
