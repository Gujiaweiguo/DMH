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

type GetBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandLogic {
	return &GetBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandLogic) GetBrand() (resp *types.BrandResp, err error) {
	// 从路径参数获取品牌ID
	brandId := l.ctx.Value("brandId").(string)
	if brandId == "" {
		return nil, fmt.Errorf("品牌ID不能为空")
	}

	userInfo := l.ctx.Value("userInfo").(map[string]interface{})
	userId := int64(userInfo["userId"].(float64))
	role := userInfo["role"].(string)

	var brand model.Brand
	if err = l.svcCtx.DB.Where("id = ?", brandId).First(&brand).Error; err != nil {
		l.Logger.Errorf("查询品牌失败: %v", err)
		return nil, fmt.Errorf("品牌不存在")
	}

	// 权限检查
	if role == "brand_admin" {
		// 品牌管理员只能查看自己管理的品牌
		var userBrand model.UserBrand
		if err = l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userId, brandId).First(&userBrand).Error; err != nil {
			return nil, fmt.Errorf("权限不足，无法查看该品牌")
		}
	} else if role == "participant" && brand.Status != "active" {
		return nil, fmt.Errorf("品牌不可用")
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
