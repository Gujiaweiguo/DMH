// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetDistributorStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorStatisticsLogic {
	return &GetDistributorStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorStatisticsLogic) GetDistributorStatistics(req *types.GetDistributorStatisticsReq) (resp *types.DistributorStatisticsResponse, err error) {
	userId, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	distributor := &model.Distributor{}
	if err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ? AND status = ?", userId, req.BrandId, "active").First(distributor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.DistributorStatisticsResponse{
				Code: 200,
				Data: types.DistributorStatisticsResp{},
			}, nil
		}
		return nil, fmt.Errorf("failed to query distributor: %w", err)
	}

	var totalOrders int64
	if err := l.svcCtx.DB.Model(&model.Order{}).
		Joins("JOIN campaigns ON campaigns.id = orders.campaign_id").
		Where("orders.referrer_id = ? AND campaigns.brand_id = ?", userId, req.BrandId).
		Count(&totalOrders).Error; err != nil {
		return nil, fmt.Errorf("failed to count orders: %w", err)
	}

	var balance float64
	userBalance := &model.UserBalance{}
	if err := l.svcCtx.DB.Where("user_id = ?", userId).First(userBalance).Error; err == nil {
		balance = userBalance.Balance
	}

	var todayEarnings float64
	var monthEarnings float64
	startOfDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())

	l.svcCtx.DB.Model(&model.DistributorReward{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("distributor_id = ? AND status = ? AND settled_at >= ?", distributor.Id, "settled", startOfDay).
		Scan(&todayEarnings)

	l.svcCtx.DB.Model(&model.DistributorReward{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("distributor_id = ? AND status = ? AND settled_at >= ?", distributor.Id, "settled", startOfMonth).
		Scan(&monthEarnings)

	resp = &types.DistributorStatisticsResponse{
		Code: 200,
		Data: types.DistributorStatisticsResp{
			DistributorId:     distributor.Id,
			TotalOrders:       totalOrders,
			TotalEarnings:     distributor.TotalEarnings,
			TodayEarnings:     todayEarnings,
			MonthEarnings:     monthEarnings,
			SubordinatesCount: distributor.SubordinatesCount,
			Balance:           balance,
			ClickCount:        0,
			ConversionRate:    0,
		},
	}

	return resp, nil
}
