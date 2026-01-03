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
		// 品牌管理员只能为自己管理的品牌创建素材
		var userBrand model.UserBrand
		if err = l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userId, brandId).First(&userBrand).Error; err != nil {
			return nil, fmt.Errorf("权限不足，无法为该品牌创建素材")
		}
	} else if role != "platform_admin" {
		return nil, fmt.Errorf("权限不足，只有平台管理员和品牌管理员可以创建品牌素材")
	}

	// 验证品牌存在
	var brand model.Brand
	if err = l.svcCtx.DB.Where("id = ?", brandId).First(&brand).Error; err != nil {
		return nil, fmt.Errorf("品牌不存在")
	}

	// 创建品牌素材
	asset := &model.BrandAsset{
		BrandID:     brand.Id,
		Name:        req.Name,
		Type:        req.Type,
		Category:    req.Category,
		Tags:        req.Tags,
		FileUrl:     req.FileUrl,
		FileSize:    req.FileSize,
		Description: req.Description,
		Status:      "active",
	}

	if err = l.svcCtx.DB.Create(asset).Error; err != nil {
		l.Logger.Errorf("创建品牌素材失败: %v", err)
		return nil, fmt.Errorf("创建品牌素材失败")
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
