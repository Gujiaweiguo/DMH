package distributor

import (
	"context"
	"errors"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyDistributorStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyDistributorStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyDistributorStatusLogic {
	return &GetMyDistributorStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyDistributorStatusLogic) GetMyDistributorStatus(req *types.GetMyDistributorStatusReq) (resp *types.DistributorResp, err error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId <= 0 {
		return nil, errors.New("用户未登录")
	}

	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("user_id = ? AND status = ?", userId, "active").
		First(&distributor).Error; err != nil {
		return nil, errors.New("您还不是分销商")
	}

	var brand model.Brand
	brandName := ""
	if err := l.svcCtx.DB.Where("id = ?", distributor.BrandId).First(&brand).Error; err == nil {
		brandName = brand.Name
	}

	resp = &types.DistributorResp{
		Id:                distributor.Id,
		UserId:            distributor.UserId,
		BrandId:           distributor.BrandId,
		BrandName:         brandName,
		Level:             distributor.Level,
		Status:            distributor.Status,
		TotalEarnings:     distributor.TotalEarnings,
		SubordinatesCount: distributor.SubordinatesCount,
		CreatedAt:         distributor.CreatedAt.Format(time.RFC3339),
	}

	if distributor.ParentId != nil {
		resp.ParentId = *distributor.ParentId
	}
	if distributor.ApprovedBy != nil {
		resp.ApprovedBy = *distributor.ApprovedBy
	}
	if distributor.ApprovedAt != nil {
		resp.ApprovedAt = distributor.ApprovedAt.Format(time.RFC3339)
	}

	return resp, nil
}
