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

type GetBrandAssetsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandAssetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandAssetsLogic {
	return &GetBrandAssetsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandAssetsLogic) GetBrandAssets() (resp *types.BrandAssetListResp, err error) {
	// 从路径参数获取品牌ID
	brandId := l.ctx.Value("brandId").(string)
	if brandId == "" {
		return nil, fmt.Errorf("品牌ID不能为空")
	}

	userInfo := l.ctx.Value("userInfo").(map[string]interface{})
	userId := int64(userInfo["userId"].(float64))
	role := userInfo["role"].(string)

	// 权限检查
	if role == "brand_admin" {
		// 品牌管理员只能查看自己管理的品牌素材
		var userBrand model.UserBrand
		if err = l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userId, brandId).First(&userBrand).Error; err != nil {
			return nil, fmt.Errorf("权限不足，无法查看该品牌素材")
		}
	} else if role != "platform_admin" {
		return nil, fmt.Errorf("权限不足，只有平台管理员和品牌管理员可以查看品牌素材")
	}

	// 验证品牌存在
	var brand model.Brand
	if err = l.svcCtx.DB.Where("id = ?", brandId).First(&brand).Error; err != nil {
		return nil, fmt.Errorf("品牌不存在")
	}

	// 查询品牌素材
	var assets []model.BrandAsset
	query := l.svcCtx.DB.Where("brand_id = ? AND status = ?", brandId, "active")
	
	if err = query.Order("created_at DESC").Find(&assets).Error; err != nil {
		l.Logger.Errorf("查询品牌素材失败: %v", err)
		return nil, fmt.Errorf("查询品牌素材失败")
	}

	var assetList []types.BrandAssetResp
	for _, asset := range assets {
		assetList = append(assetList, types.BrandAssetResp{
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
		})
	}

	return &types.BrandAssetListResp{
		Total:  int64(len(assetList)),
		Assets: assetList,
	}, nil
}
