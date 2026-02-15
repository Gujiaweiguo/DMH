// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"context"
	"errors"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandAssetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandAssetLogic {
	return &GetBrandAssetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandAssetLogic) GetBrandAsset(req *types.GetBrandAssetReq) (resp *types.BrandAssetResp, err error) {
	if req.BrandId <= 0 || req.Id <= 0 {
		return nil, errors.New("参数无效")
	}

	var asset model.BrandAsset
	if err := l.svcCtx.DB.Where("id = ? AND brand_id = ?", req.Id, req.BrandId).First(&asset).Error; err != nil {
		return nil, errors.New("素材不存在")
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
