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

type GetBrandsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandsLogic {
	return &GetBrandsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandsLogic) GetBrands() (resp *types.BrandListResp, err error) {
	var brands []model.Brand

	if err := l.svcCtx.DB.Find(&brands).Error; err != nil {
		l.Errorf("Failed to query brands: %v", err)
		return nil, fmt.Errorf("Failed to query brands: %w", err)
	}

	l.Infof("Successfully queried brands: count=%d", len(brands))

	var brandResps []types.BrandResp
	for _, brand := range brands {
		brandResps = append(brandResps, types.BrandResp{
			Id:          brand.Id,
			Name:        brand.Name,
			Logo:        brand.Logo,
			Description: brand.Description,
			Status:      brand.Status,
			CreatedAt:   brand.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:   brand.UpdatedAt.Format("2006-01-02T15:04:05"),
		})
	}

	resp = &types.BrandListResp{
		Total:  int64(len(brandResps)),
		Brands: brandResps,
	}

	return resp, nil
}
