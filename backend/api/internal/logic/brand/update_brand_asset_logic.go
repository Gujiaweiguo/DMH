// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"context"
	"fmt"

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

func (l *UpdateBrandAssetLogic) UpdateBrandAsset(req *types.BrandAssetReq) (resp *types.BrandAssetResp, err error) {
	// 从路径参数获取品牌ID和素材ID
	brandId := l.ctx.Value("brandId").(string)
	assetId := l.ctx.Value("assetId").(string)
	if brandId == "" || assetId == "" {
		return nil, fmt.Errorf("品牌ID和素材ID不能为空")
	}

	userInfo := l.ctx.Value("userInfo").(map[string]interface{})
	userId := int64(userInfo["userId"].(float64))
	role := userInfo["role"].(string)

	// 权限检查
	if role == "brand_admin" {
		// 品牌管理员只能更新自己管理的品牌素材
		var userBrand model.UserBrand
		if err = l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userId, brandId).First(&userBrand).Error; err != nil {
			return nil, fmt.Errorf("权限不足，无法更新该品牌素材")
		}
	} else if role != "platform_admin" {
		return nil, fmt.Errorf("权限不足，只有平台管理员和品牌管理员可以更新品牌素材")
	}

	// 查询素材
	var asset model.BrandAsset
	if err = l.svcCtx.DB.Where("id = ? AND brand_id = ?", assetId, brandId).First(&asset).Error; err != nil {
		l.Logger.Errorf("查询品牌素材失败: %v", err)
		return nil, fmt.Errorf("品牌素材不存在")
	}

	// 更新素材信息
	updates := make(map[string]interface{})
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

	if len(updates) > 0 {
		if err = l.svcCtx.DB.Model(&asset).Updates(updates).Error; err != nil {
			l.Logger.Errorf("更新品牌素材失败: %v", err)
			return nil, fmt.Errorf("更新品牌素材失败")
		}
	}

	// 重新查询更新后的素材
	if err = l.svcCtx.DB.Where("id = ?", assetId).First(&asset).Error; err != nil {
		l.Logger.Errorf("查询更新后的品牌素材失败: %v", err)
		return nil, fmt.Errorf("查询品牌素材失败")
	}

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
