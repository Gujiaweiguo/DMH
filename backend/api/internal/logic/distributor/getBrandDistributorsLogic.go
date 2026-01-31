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

type GetBrandDistributorsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandDistributorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandDistributorsLogic {
	return &GetBrandDistributorsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandDistributorsLogic) GetBrandDistributors(req *types.GetDistributorsReq) (resp *types.DistributorListResp, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	query := l.svcCtx.DB.Model(&model.Distributor{})

	if req.Level > 0 {
		query = query.Where("level = ?", req.Level)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		query = query.Joins("LEFT JOIN users ON users.id = distributors.user_id").
			Where("users.username LIKE ?", "%"+req.Keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var distributors []model.Distributor
	offset := int((req.Page - 1) * req.PageSize)
	if err := query.Offset(offset).Limit(int(req.PageSize)).
		Preload("User").
		Preload("Brand").
		Preload("Parent.User").
		Find(&distributors).Error; err != nil {
		return nil, err
	}

	resp = &types.DistributorListResp{
		Total:        total,
		Distributors: make([]types.DistributorResp, 0, len(distributors)),
	}

	for _, dist := range distributors {
		distributorResp := types.DistributorResp{
			Id:                dist.Id,
			UserId:            dist.UserId,
			BrandId:           dist.BrandId,
			Level:             dist.Level,
			Status:            dist.Status,
			TotalEarnings:     dist.TotalEarnings,
			SubordinatesCount: dist.SubordinatesCount,
			CreatedAt:         dist.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if dist.User != nil {
			distributorResp.Username = dist.User.Username
		}

		if dist.Brand != nil {
			distributorResp.BrandName = dist.Brand.Name
		}

		if dist.ParentId != nil {
			distributorResp.ParentId = *dist.ParentId
			if dist.Parent != nil && dist.Parent.User != nil {
				distributorResp.ParentName = dist.Parent.User.Username
			}
		}

		if dist.ApprovedBy != nil {
			distributorResp.ApprovedBy = *dist.ApprovedBy
		}

		if dist.ApprovedAt != nil {
			distributorResp.ApprovedAt = dist.ApprovedAt.Format("2006-01-02 15:04:05")
		}

		resp.Distributors = append(resp.Distributors, distributorResp)
	}

	return resp, nil
}
