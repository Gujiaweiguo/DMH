// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandDistributorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandDistributorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandDistributorLogic {
	return &GetBrandDistributorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandDistributorLogic) GetBrandDistributor() (resp *types.DistributorResp, err error) {
	distributorId := l.ctx.Value("distributorId").(int64)

	distributor := &model.Distributor{}
	if err := l.svcCtx.DB.Preload("User").Preload("Brand").Preload("Parent.User").First(distributor, distributorId).Error; err != nil {
		return nil, err
	}

	resp = &types.DistributorResp{
		Id:                distributor.Id,
		UserId:            distributor.UserId,
		BrandId:           distributor.BrandId,
		Level:             distributor.Level,
		Status:            distributor.Status,
		TotalEarnings:     distributor.TotalEarnings,
		SubordinatesCount: distributor.SubordinatesCount,
		CreatedAt:         distributor.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if distributor.User != nil {
		resp.Username = distributor.User.Username
	}

	if distributor.Brand != nil {
		resp.BrandName = distributor.Brand.Name
	}

	if distributor.ParentId != nil {
		resp.ParentId = *distributor.ParentId
		if distributor.Parent != nil && distributor.Parent.User != nil {
			resp.ParentName = distributor.Parent.User.Username
		}
	}

	if distributor.ApprovedBy != nil {
		resp.ApprovedBy = *distributor.ApprovedBy
	}

	if distributor.ApprovedAt != nil {
		resp.ApprovedAt = distributor.ApprovedAt.Format("2006-01-02 15:04:05")
	}

	return resp, nil
}
