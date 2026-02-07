// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"
	"fmt"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyDistributorDashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyDistributorDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyDistributorDashboardLogic {
	return &GetMyDistributorDashboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyDistributorDashboardLogic) GetMyDistributorDashboard() (resp *types.DistributorDashboardResp, err error) {
	userId, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	var distributors []model.Distributor
	if err := l.svcCtx.DB.Where("user_id = ? AND status = ?", userId, "active").Find(&distributors).Error; err != nil {
		return nil, fmt.Errorf("failed to query distributors: %w", err)
	}

	if len(distributors) == 0 {
		return &types.DistributorDashboardResp{
			Code: 200,
			Data: types.DistributorDashboardData{
				HasDistributor: false,
				Brands:         []types.DistributorBrandOption{},
			},
		}, nil
	}

	uniqueBrandIDs := make(map[int64]struct{}, len(distributors))
	for _, dist := range distributors {
		uniqueBrandIDs[dist.BrandId] = struct{}{}
	}

	brandIds := make([]int64, 0, len(uniqueBrandIDs))
	for brandId := range uniqueBrandIDs {
		brandIds = append(brandIds, brandId)
	}

	var brands []model.Brand
	if err := l.svcCtx.DB.Where("id IN ?", brandIds).Find(&brands).Error; err != nil {
		return nil, fmt.Errorf("failed to query brands: %w", err)
	}

	respBrands := make([]types.DistributorBrandOption, 0, len(brands))
	for _, brand := range brands {
		respBrands = append(respBrands, types.DistributorBrandOption{
			BrandId:   brand.Id,
			BrandName: brand.Name,
		})
	}

	resp = &types.DistributorDashboardResp{
		Code: 200,
		Data: types.DistributorDashboardData{
			HasDistributor: true,
			Brands:         respBrands,
		},
	}

	return resp, nil
}
