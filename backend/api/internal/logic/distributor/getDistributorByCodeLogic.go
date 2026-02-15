package distributor

import (
	"context"
	"errors"
	"strconv"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDistributorByCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorByCodeLogic {
	return &GetDistributorByCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorByCodeLogic) GetDistributorByCode(req *types.GetDistributorByCodeReq) (resp *types.DistributorResp, err error) {
	if req.Code == "" {
		return nil, errors.New("分销商代码不能为空")
	}

	distributorId, err := strconv.ParseInt(req.Code, 10, 64)
	if err != nil {
		return nil, errors.New("无效的分销商代码")
	}

	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("id = ? AND status = ?", distributorId, "active").
		First(&distributor).Error; err != nil {
		return nil, errors.New("分销商不存在")
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
