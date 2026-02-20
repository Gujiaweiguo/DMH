package brand

import (
	"context"
	"fmt"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBrandLogic_CreateBrand_Success(t *testing.T) {
	db := setupBrandTestDB(t)
	l := NewCreateBrandLogic(context.Background(), &svc.ServiceContext{DB: db})

	resp, err := l.CreateBrand(&types.CreateBrandReq{
		Name:        "  New Brand  ",
		Logo:        "https://img/logo.png",
		Description: "desc",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "New Brand", resp.Name)
	assert.Equal(t, "https://img/logo.png", resp.Logo)
	assert.Equal(t, "active", resp.Status)
}

func TestCreateBrandLogic_CreateBrand_EmptyName(t *testing.T) {
	db := setupBrandTestDB(t)
	l := NewCreateBrandLogic(context.Background(), &svc.ServiceContext{DB: db})

	resp, err := l.CreateBrand(&types.CreateBrandReq{Name: "   "})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "品牌名称不能为空")
}

func TestCreateBrandAssetLogic_CreateBrandAsset_Success(t *testing.T) {
	db := setupBrandTestDB(t)
	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateBrandAssetLogic(ctx, svcCtx)

	req := &types.BrandAssetReq{
		BrandId:  fmt.Sprintf("%d", brand.Id),
		Name:     "测试素材",
		Type:     "image",
		FileUrl:  "https://example.com/asset.jpg",
		FileSize: 1024,
	}

	resp, err := logic.CreateBrandAsset(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "测试素材", resp.Name)
	assert.Equal(t, brand.Id, resp.BrandId)
}

func TestDeleteBrandAssetLogic_DeleteBrandAsset_Success(t *testing.T) {
	db := setupBrandTestDB(t)
	brand := createTestBrand(t, db, "Test Brand", "active")
	asset := &model.BrandAsset{
		BrandID:  brand.Id,
		Name:     "测试素材",
		Type:     "image",
		FileUrl:  "https://example.com/asset.jpg",
		FileSize: 1024,
	}
	db.Create(asset)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewDeleteBrandAssetLogic(ctx, svcCtx)

	req := &types.DeleteBrandAssetReq{
		BrandId: brand.Id,
		Id:      asset.ID,
	}

	resp, err := logic.DeleteBrandAsset(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "删除成功", resp.Message)
}

func TestGetBrandAssetLogic_GetBrandAsset_Success(t *testing.T) {
	db := setupBrandTestDB(t)
	brand := createTestBrand(t, db, "Test Brand", "active")
	asset := &model.BrandAsset{
		BrandID:     brand.Id,
		Name:        "测试素材",
		Type:        "image",
		FileUrl:     "https://example.com/asset.jpg",
		FileSize:    1024,
		Description: "测试描述",
	}
	db.Create(asset)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandAssetLogic(ctx, svcCtx)

	req := &types.GetBrandAssetReq{
		BrandId: brand.Id,
		Id:      asset.ID,
	}

	resp, err := logic.GetBrandAsset(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "测试素材", resp.Name)
	assert.Equal(t, "测试描述", resp.Description)
}

func TestUpdateBrandAssetLogic_UpdateBrandAsset_Success(t *testing.T) {
	db := setupBrandTestDB(t)
	brand := createTestBrand(t, db, "Test Brand", "active")
	asset := &model.BrandAsset{
		BrandID:     brand.Id,
		Name:        "原名称",
		Type:        "image",
		FileUrl:     "https://example.com/asset.jpg",
		FileSize:    1024,
		Description: "原描述",
	}
	db.Create(asset)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateBrandAssetLogic(ctx, svcCtx)

	req := &types.UpdateBrandAssetReq{
		BrandId:     brand.Id,
		Id:          asset.ID,
		Name:        "新名称",
		Description: "新描述",
	}

	resp, err := logic.UpdateBrandAsset(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "新名称", resp.Name)
	assert.Equal(t, "新描述", resp.Description)
}

func TestUpdateBrandLogic_UpdateBrand_NotFound(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateBrandLogic(ctx, svcCtx)

	resp, err := logic.UpdateBrand(&types.UpdateBrandReq{Id: 9999, Name: "新名称"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "品牌不存在")
}

func TestGetBrandLogic_GetBrand_Success(t *testing.T) {
	db := setupBrandTestDB(t)
	brand := createTestBrand(t, db, "Brand A", "active")
	l := NewGetBrandLogic(context.Background(), &svc.ServiceContext{DB: db})

	resp, err := l.GetBrand(&types.GetBrandReq{Id: brand.Id})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, brand.Id, resp.Id)
	assert.Equal(t, "Brand A", resp.Name)
}

func TestGetBrandLogic_GetBrand_InvalidID(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandLogic(ctx, svcCtx)

	resp, err := logic.GetBrand(&types.GetBrandReq{Id: 9999})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "品牌不存在")
}

func TestGetBrandStatsLogic_Success(t *testing.T) {
	db := setupBrandTestDB(t)
	brand := createTestBrand(t, db, "Test Brand", "active")
	campaign := &model.Campaign{
		Name:       "Test Campaign",
		BrandId:    brand.Id,
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		RewardRule: 10,
		CreatedAt:  db.NowFunc(),
	}
	db.Create(campaign)

	order := &model.Order{
		CampaignId: campaign.Id,
		Amount:     100,
		PayStatus:  "paid",
		FormData:   "{}",
		CreatedAt:  db.NowFunc(),
	}
	db.Create(order)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandStatsLogic(ctx, svcCtx)

	req := &types.GetBrandStatsReq{
		Id: brand.Id,
	}

	resp, err := logic.GetBrandStats(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, brand.Id, resp.BrandId)
	assert.Equal(t, int64(1), resp.TotalCampaigns)
	assert.Equal(t, int64(1), resp.ActiveCampaigns)
	assert.Equal(t, int64(1), resp.TotalOrders)
}

func TestGetBrandStatsLogic_InvalidBrandID(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandStatsLogic(ctx, svcCtx)

	req := &types.GetBrandStatsReq{
		Id: 9999,
	}

	resp, err := logic.GetBrandStats(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "品牌不存在")
}

func TestGetBrandAssetsLogic_GetBrandAssets_Success(t *testing.T) {
	db := setupBrandTestDB(t)
	require.NoError(t, db.AutoMigrate(&model.BrandAsset{}))
	brand := createTestBrand(t, db, "Brand A", "active")

	require.NoError(t, db.Create(&model.BrandAsset{BrandID: brand.Id, Name: "A1", Type: "image", FileUrl: "https://f/1.png"}).Error)
	require.NoError(t, db.Create(&model.BrandAsset{BrandID: brand.Id, Name: "A2", Type: "image", FileUrl: "https://f/2.png"}).Error)

	l := NewGetBrandAssetsLogic(context.Background(), &svc.ServiceContext{DB: db})
	resp, err := l.GetBrandAssets(&types.GetBrandAssetsReq{Id: brand.Id})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(2), resp.Total)
	assert.Len(t, resp.Assets, 2)
	assert.Equal(t, brand.Id, resp.Assets[0].BrandId)
}

func TestGetBrandAssetsLogic_GetBrandAssets_InvalidID(t *testing.T) {
	db := setupBrandTestDB(t)
	l := NewGetBrandAssetsLogic(context.Background(), &svc.ServiceContext{DB: db})

	resp, err := l.GetBrandAssets(&types.GetBrandAssetsReq{Id: -1})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "品牌ID无效")
}

func TestUpdateBrandLogic_UpdateBrand_Success(t *testing.T) {
	db := setupBrandTestDB(t)
	brand := createTestBrand(t, db, "Original Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateBrandLogic(ctx, svcCtx)

	req := &types.UpdateBrandReq{
		Id:          brand.Id,
		Name:        "Updated Brand",
		Logo:        "https://example.com/new-logo.png",
		Description: "Updated description",
		Status:      "disabled",
	}

	resp, err := logic.UpdateBrand(req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, brand.Id, resp.Id)
	assert.Equal(t, "Updated Brand", resp.Name)
	assert.Equal(t, "https://example.com/new-logo.png", resp.Logo)
	assert.Equal(t, "Updated description", resp.Description)
	assert.Equal(t, "disabled", resp.Status)
}

func TestUpdateBrandLogic_UpdateBrand_InvalidID(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateBrandLogic(ctx, svcCtx)

	req := &types.UpdateBrandReq{
		Id:   0,
		Name: "Test",
	}

	resp, err := logic.UpdateBrand(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "品牌ID无效")
}

func TestUpdateBrandLogic_UpdateBrand_PartialUpdate(t *testing.T) {
	db := setupBrandTestDB(t)
	brand := createTestBrand(t, db, "Original Brand", "active")

	originalDescription := brand.Description
	originalStatus := brand.Status

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateBrandLogic(ctx, svcCtx)

	req := &types.UpdateBrandReq{
		Id:   brand.Id,
		Name: "  Updated Brand  ",
	}

	resp, err := logic.UpdateBrand(req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "Updated Brand", resp.Name)
	assert.Equal(t, originalDescription, resp.Description)
	assert.Equal(t, originalStatus, resp.Status)
}

func TestCreateBrandAssetLogic_CreateBrandAsset_InvalidBrandId(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateBrandAssetLogic(ctx, svcCtx)

	req := &types.BrandAssetReq{
		BrandId:  "invalid",
		Name:     "Test Asset",
		Type:     "image",
		FileUrl:  "https://example.com/asset.jpg",
		FileSize: 1024,
	}

	resp, err := logic.CreateBrandAsset(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCreateBrandAssetLogic_CreateBrandAsset_BrandNotFound(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateBrandAssetLogic(ctx, svcCtx)

	req := &types.BrandAssetReq{
		BrandId:  "9999",
		Name:     "Test Asset",
		Type:     "image",
		FileUrl:  "https://example.com/asset.jpg",
		FileSize: 1024,
	}

	resp, err := logic.CreateBrandAsset(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestDeleteBrandAssetLogic_DeleteBrandAsset_NotFound(t *testing.T) {
	db := setupBrandTestDB(t)
	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewDeleteBrandAssetLogic(ctx, svcCtx)

	req := &types.DeleteBrandAssetReq{
		BrandId: brand.Id,
		Id:      9999,
	}

	resp, err := logic.DeleteBrandAsset(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetBrandAssetLogic_GetBrandAsset_NotFound(t *testing.T) {
	db := setupBrandTestDB(t)
	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandAssetLogic(ctx, svcCtx)

	req := &types.GetBrandAssetReq{
		BrandId: brand.Id,
		Id:      9999,
	}

	resp, err := logic.GetBrandAsset(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpdateBrandAssetLogic_UpdateBrandAsset_NotFound(t *testing.T) {
	db := setupBrandTestDB(t)
	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateBrandAssetLogic(ctx, svcCtx)

	req := &types.UpdateBrandAssetReq{
		BrandId: brand.Id,
		Id:      9999,
		Name:    "New Name",
	}

	resp, err := logic.UpdateBrandAsset(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetBrandStatsLogic_InvalidID(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandStatsLogic(ctx, svcCtx)

	req := &types.GetBrandStatsReq{
		Id: -1,
	}

	resp, err := logic.GetBrandStats(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "品牌ID无效")
}

func TestGetBrandStatsLogic_WithMultipleCampaigns(t *testing.T) {
	db := setupBrandTestDB(t)
	require.NoError(t, db.AutoMigrate(&model.Reward{}, &model.Member{}))
	brand := createTestBrand(t, db, "Test Brand", "active")

	campaign1 := &model.Campaign{Name: "Campaign 1", BrandId: brand.Id, Status: "active", StartTime: time.Now(), EndTime: time.Now().Add(24 * time.Hour), RewardRule: 10}
	campaign2 := &model.Campaign{Name: "Campaign 2", BrandId: brand.Id, Status: "active", StartTime: time.Now(), EndTime: time.Now().Add(24 * time.Hour), RewardRule: 10}
	campaign3 := &model.Campaign{Name: "Campaign 3", BrandId: brand.Id, Status: "disabled", StartTime: time.Now(), EndTime: time.Now().Add(24 * time.Hour), RewardRule: 10}
	db.Create(campaign1)
	db.Create(campaign2)
	db.Create(campaign3)

	order1 := &model.Order{CampaignId: campaign1.Id, Amount: 100, PayStatus: "paid", FormData: "{}"}
	order2 := &model.Order{CampaignId: campaign2.Id, Amount: 200, PayStatus: "paid", FormData: "{}"}
	order3 := &model.Order{CampaignId: campaign3.Id, Amount: 50, PayStatus: "unpaid", FormData: "{}"}
	db.Create(order1)
	db.Create(order2)
	db.Create(order3)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandStatsLogic(ctx, svcCtx)

	resp, err := logic.GetBrandStats(&types.GetBrandStatsReq{Id: brand.Id})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(3), resp.TotalCampaigns)
	assert.Equal(t, int64(2), resp.ActiveCampaigns)
	assert.Equal(t, int64(3), resp.TotalOrders)
	assert.Equal(t, 300.0, resp.TotalRevenue)
}
