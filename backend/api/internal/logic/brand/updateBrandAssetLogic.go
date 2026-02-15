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

type UpdateBrandAssetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBrandAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBrandAssetLogic {
	return &UpdateBrandAssetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBrandAssetLogic) UpdateBrandAsset(req *types.UpdateBrandAssetReq) (resp *types.BrandAssetResp, err error) {
	if req.BrandId <= 0 || req.Id <= 0 {
		return nil, errors.New("参数无效")
	}

	var asset model.BrandAsset
	if err := l.svcCtx.DB.Where("id = ? AND brand_id = ?", req.Id, req.BrandId).First(&asset).Error; err != nil {
		return nil, errors.New("素材不存在")
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}
	if req.FileUrl != "" {
		updates["file_url"] = req.FileUrl
	}
	if req.FileSize > 0 {
		updates["file_size"] = req.FileSize
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if len(updates) == 0 {
		return &types.BrandAssetResp{
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
		}, nil
	}

	if err := l.svcCtx.DB.Model(&asset).Updates(updates).Error; err != nil {
		l.Errorf("更新品牌素材失败: %v", err)
		return nil, errors.New("更新品牌素材失败")
	}

	if err := l.svcCtx.DB.First(&asset, req.Id).Error; err != nil {
		return nil, errors.New("获取更新后的素材失败")
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
