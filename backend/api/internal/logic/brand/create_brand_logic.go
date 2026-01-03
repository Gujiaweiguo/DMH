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

type CreateBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBrandLogic {
	return &CreateBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBrandLogic) CreateBrand(req *types.CreateBrandReq) (resp *types.BrandResp, err error) {
	// 只有平台管理员可以创建品牌
	userInfo := l.ctx.Value("userInfo").(map[string]interface{})
	role := userInfo["role"].(string)
	if role != "platform_admin" {
		return nil, fmt.Errorf("权限不足，只有平台管理员可以创建品牌")
	}

	// 创建品牌
	brand := &model.Brand{
		Name:        req.Name,
		Logo:        req.Logo,
		Description: req.Description,
		Status:      "active",
	}

	if err = l.svcCtx.DB.Create(brand).Error; err != nil {
		l.Logger.Errorf("创建品牌失败: %v", err)
		return nil, fmt.Errorf("创建品牌失败")
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
