// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"context"
	"errors"
	"strconv"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBrandAssetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBrandAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBrandAssetLogic {
	return &CreateBrandAssetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBrandAssetLogic) CreateBrandAsset(req *types.BrandAssetReq) (resp *types.BrandAssetResp, err error) {
	if req.Name == "" {
		return nil, errors.New("素材名称不能为空")
	}

	if req.FileUrl == "" {
		return nil, errors.New("文件URL不能为空")
	}

	brandId, err := strconv.ParseInt(req.BrandId, 10, 64)
	if err != nil || brandId <= 0 {
		return nil, errors.New("品牌ID无效")
	}

	var brand model.Brand
	if err := l.svcCtx.DB.First(&brand, brandId).Error; err != nil {
		return nil, errors.New("品牌不存在")
	}

	asset := &model.BrandAsset{
		BrandID:     brandId,
		Name:        req.Name,
		Type:        req.Type,
		Category:    req.Category,
		Tags:        req.Tags,
		FileUrl:     req.FileUrl,
		FileSize:    req.FileSize,
		Description: req.Description,
		Status:      "active",
	}

	if err := l.svcCtx.DB.Create(asset).Error; err != nil {
		l.Errorf("创建品牌素材失败: %v", err)
		return nil, errors.New("创建品牌素材失败")
	}

	resp = &types.BrandAssetResp{
		Id:          asset.ID,
		BrandId:     asset.BrandID,
		Name:        asset.Name,
		Type:        asset.Type,
		Category:    asset.Category,
		Tags:        asset.Tags,
		FileUrl:     asset.FileUrl,
		FileSize:    asset.FileSize,
		Description: asset.Description,
		CreatedAt:   asset.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   asset.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}
