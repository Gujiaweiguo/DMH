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

type GetDistributorSubordinatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorSubordinatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorSubordinatesLogic {
	return &GetDistributorSubordinatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorSubordinatesLogic) GetDistributorSubordinates(req *types.GetDistributorSubordinatesReq) (resp *types.SubordinateListResp, err error) {
	if req.DistributorId <= 0 {
		return nil, errors.New("分销商ID无效")
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	var total int64
	if err := l.svcCtx.DB.Model(&model.Distributor{}).
		Where("parent_id = ? AND status = ?", req.DistributorId, "active").
		Count(&total).Error; err != nil {
		l.Logger.Errorf("查询下级分销商总数失败: %v", err)
		return nil, err
	}

	var subordinates []model.Distributor
	if err := l.svcCtx.DB.Where("parent_id = ? AND status = ?", req.DistributorId, "active").
		Order("created_at DESC").
		Limit(int(pageSize)).Offset(int(offset)).
		Find(&subordinates).Error; err != nil {
		l.Logger.Errorf("查询下级分销商列表失败: %v", err)
		return nil, err
	}

	subList := make([]types.SubordinateResp, 0, len(subordinates))
	for _, sub := range subordinates {
		subList = append(subList, types.SubordinateResp{
			Id:            sub.Id,
			UserId:        sub.UserId,
			Level:         sub.Level,
			TotalEarnings: sub.TotalEarnings,
			CreatedAt:     sub.CreatedAt.Format(time.RFC3339),
		})
	}

	return &types.SubordinateListResp{
		Total:        total,
		Subordinates: subList,
	}, nil
}
