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

type DeleteBrandAssetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBrandAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBrandAssetLogic {
	return &DeleteBrandAssetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBrandAssetLogic) DeleteBrandAsset() (resp *types.CommonResp, err error) {
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
		// 品牌管理员只能删除自己管理的品牌素材
		var userBrand model.UserBrand
		if err = l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userId, brandId).First(&userBrand).Error; err != nil {
			return nil, fmt.Errorf("权限不足，无法删除该品牌素材")
		}
	} else if role != "platform_admin" {
		return nil, fmt.Errorf("权限不足，只有平台管理员和品牌管理员可以删除品牌素材")
	}

	// 查询素材
	var asset model.BrandAsset
	if err = l.svcCtx.DB.Where("id = ? AND brand_id = ?", assetId, brandId).First(&asset).Error; err != nil {
		l.Logger.Errorf("查询品牌素材失败: %v", err)
		return nil, fmt.Errorf("品牌素材不存在")
	}

	// 软删除素材（更新状态为disabled）
	if err = l.svcCtx.DB.Model(&asset).Update("status", "disabled").Error; err != nil {
		l.Logger.Errorf("删除品牌素材失败: %v", err)
		return nil, fmt.Errorf("删除品牌素材失败")
	}

	return &types.CommonResp{
		Message: "品牌素材删除成功",
	}, nil
}
