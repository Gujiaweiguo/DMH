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

type UpdateBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBrandLogic {
	return &UpdateBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBrandLogic) UpdateBrand(req *types.UpdateBrandReq) (resp *types.BrandResp, err error) {
	// 从路径参数获取品牌ID
	brandId := l.ctx.Value("brandId").(string)
	if brandId == "" {
		return nil, fmt.Errorf("品牌ID不能为空")
	}

	userInfo := l.ctx.Value("userInfo").(map[string]interface{})
	userId := int64(userInfo["userId"].(float64))
	role := userInfo["role"].(string)

	// 查询品牌
	var brand model.Brand
	if err = l.svcCtx.DB.Where("id = ?", brandId).First(&brand).Error; err != nil {
		l.Logger.Errorf("查询品牌失败: %v", err)
		return nil, fmt.Errorf("品牌不存在")
	}

	// 权限检查
	if role == "brand_admin" {
		// 品牌管理员只能更新自己管理的品牌
		var userBrand model.UserBrand
		if err = l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userId, brandId).First(&userBrand).Error; err != nil {
			return nil, fmt.Errorf("权限不足，无法更新该品牌")
		}
	} else if role != "platform_admin" {
		return nil, fmt.Errorf("权限不足，只有平台管理员和品牌管理员可以更新品牌")
	}

	// 更新品牌信息
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Logo != "" {
		updates["logo"] = req.Logo
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != "" {
		// 只有平台管理员可以更改状态
		if role == "platform_admin" {
			updates["status"] = req.Status
		}
	}

	if len(updates) > 0 {
		if err = l.svcCtx.DB.Model(&brand).Updates(updates).Error; err != nil {
			l.Logger.Errorf("更新品牌失败: %v", err)
			return nil, fmt.Errorf("更新品牌失败")
		}
	}

	// 重新查询更新后的品牌
	if err = l.svcCtx.DB.Where("id = ?", brandId).First(&brand).Error; err != nil {
		l.Logger.Errorf("查询更新后的品牌失败: %v", err)
		return nil, fmt.Errorf("查询品牌失败")
	}

	return &types.BrandResp{
		Id:          brand.Id,
		Name:        brand.Name,
		Logo:        brand.Logo,
		Description: brand.Description,
		Status:      brand.Status,
		CreatedAt:   brand.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   brand.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
