package brand

import (
	"context"
	"errors"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandStatsLogic {
	return &GetBrandStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandStatsLogic) GetBrandStats(req *types.GetBrandStatsReq) (resp *types.BrandStatsResp, err error) {
	if req.Id <= 0 {
		return nil, errors.New("品牌ID无效")
	}

	var brand model.Brand
	if err := l.svcCtx.DB.Where("id = ?", req.Id).First(&brand).Error; err != nil {
		return nil, errors.New("品牌不存在")
	}

	var totalCampaigns int64
	if err := l.svcCtx.DB.Model(&model.Campaign{}).Where("brand_id = ?", brand.Id).Count(&totalCampaigns).Error; err != nil {
		return nil, errors.New("查询活动统计失败")
	}

	var activeCampaigns int64
	if err := l.svcCtx.DB.Model(&model.Campaign{}).Where("brand_id = ? AND status = ?", brand.Id, "active").Count(&activeCampaigns).Error; err != nil {
		return nil, errors.New("查询活动统计失败")
	}

	var totalOrders int64
	if err := l.svcCtx.DB.Model(&model.Order{}).
		Joins("JOIN campaigns ON campaigns.id = orders.campaign_id").
		Where("campaigns.brand_id = ?", brand.Id).
		Count(&totalOrders).Error; err != nil {
		return nil, errors.New("查询订单统计失败")
	}

	var totalRevenue float64
	if err := l.svcCtx.DB.Model(&model.Order{}).
		Select("COALESCE(SUM(orders.amount), 0)").
		Joins("JOIN campaigns ON campaigns.id = orders.campaign_id").
		Where("campaigns.brand_id = ? AND orders.pay_status = ?", brand.Id, "paid").
		Scan(&totalRevenue).Error; err != nil {
		return nil, errors.New("查询订单金额失败")
	}

	var totalRewards float64
	if err := l.svcCtx.DB.Model(&model.Reward{}).
		Select("COALESCE(SUM(amount), 0)").
		Joins("JOIN orders ON orders.id = rewards.order_id").
		Joins("JOIN campaigns ON campaigns.id = orders.campaign_id").
		Where("campaigns.brand_id = ?", brand.Id).
		Scan(&totalRewards).Error; err != nil {
		l.Errorf("查询奖励金额失败: %v", err)
		totalRewards = 0
	}

	var participantCount int64
	if err := l.svcCtx.DB.Model(&model.Order{}).
		Select("COUNT(DISTINCT orders.member_id)").
		Joins("JOIN campaigns ON campaigns.id = orders.campaign_id").
		Where("campaigns.brand_id = ? AND orders.member_id IS NOT NULL", brand.Id).
		Scan(&participantCount).Error; err != nil {
		l.Errorf("查询参与人数失败: %v", err)
		participantCount = 0
	}

	conversionRate := float64(0)
	if totalOrders > 0 && participantCount > 0 {
		conversionRate = float64(participantCount) / float64(totalOrders) * 100
	}

	resp = &types.BrandStatsResp{
		BrandId:          brand.Id,
		TotalCampaigns:   totalCampaigns,
		ActiveCampaigns:  activeCampaigns,
		TotalOrders:      totalOrders,
		TotalRevenue:     totalRevenue,
		TotalRewards:     totalRewards,
		ParticipantCount: participantCount,
		ConversionRate:   conversionRate,
		LastUpdated:      brand.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}
