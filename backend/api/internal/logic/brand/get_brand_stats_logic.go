// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"context"
	"fmt"
	"time"

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

func (l *GetBrandStatsLogic) GetBrandStats() (resp *types.BrandStatsResp, err error) {
	// 从路径参数获取品牌ID
	brandId := l.ctx.Value("brandId").(string)
	if brandId == "" {
		return nil, fmt.Errorf("品牌ID不能为空")
	}

	userInfo := l.ctx.Value("userInfo").(map[string]interface{})
	userId := int64(userInfo["userId"].(float64))
	role := userInfo["role"].(string)

	// 权限检查
	if role == "brand_admin" {
		// 品牌管理员只能查看自己管理的品牌统计
		var userBrand model.UserBrand
		if err = l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userId, brandId).First(&userBrand).Error; err != nil {
			return nil, fmt.Errorf("权限不足，无法查看该品牌统计")
		}
	} else if role != "platform_admin" {
		return nil, fmt.Errorf("权限不足，只有平台管理员和品牌管理员可以查看品牌统计")
	}

	// 验证品牌存在
	var brand model.Brand
	if err = l.svcCtx.DB.Where("id = ?", brandId).First(&brand).Error; err != nil {
		return nil, fmt.Errorf("品牌不存在")
	}

	// 统计活动数量
	var totalCampaigns, activeCampaigns int64
	l.svcCtx.DB.Model(&model.Campaign{}).Where("brand_id = ?", brandId).Count(&totalCampaigns)
	l.svcCtx.DB.Model(&model.Campaign{}).Where("brand_id = ? AND status = ?", brandId, "active").Count(&activeCampaigns)

	// 统计订单数量和收入
	var totalOrders int64
	var totalRevenue float64
	l.svcCtx.DB.Model(&model.Order{}).
		Joins("JOIN campaigns ON orders.campaign_id = campaigns.id").
		Where("campaigns.brand_id = ? AND orders.pay_status = ?", brandId, "paid").
		Count(&totalOrders)

	// 计算总收入
	var result struct {
		TotalRevenue float64
	}
	l.svcCtx.DB.Model(&model.Order{}).
		Select("COALESCE(SUM(orders.amount), 0) as total_revenue").
		Joins("JOIN campaigns ON orders.campaign_id = campaigns.id").
		Where("campaigns.brand_id = ? AND orders.pay_status = ?", brandId, "paid").
		Scan(&result)
	totalRevenue = result.TotalRevenue

	// 统计奖励发放
	var totalRewards float64
	var rewardResult struct {
		TotalRewards float64
	}
	l.svcCtx.DB.Model(&model.Reward{}).
		Select("COALESCE(SUM(amount), 0) as total_rewards").
		Joins("JOIN campaigns ON rewards.campaign_id = campaigns.id").
		Where("campaigns.brand_id = ? AND rewards.status = ?", brandId, "settled").
		Scan(&rewardResult)
	totalRewards = rewardResult.TotalRewards

	// 统计参与用户数量（去重）
	var participantCount int64
	l.svcCtx.DB.Model(&model.Order{}).
		Select("COUNT(DISTINCT orders.phone)").
		Joins("JOIN campaigns ON orders.campaign_id = campaigns.id").
		Where("campaigns.brand_id = ?", brandId).
		Count(&participantCount)

	// 计算转化率（付费订单/总订单）
	var conversionRate float64
	if totalOrders > 0 {
		var paidOrders int64
		l.svcCtx.DB.Model(&model.Order{}).
			Joins("JOIN campaigns ON orders.campaign_id = campaigns.id").
			Where("campaigns.brand_id = ? AND orders.pay_status = ?", brandId, "paid").
			Count(&paidOrders)
		conversionRate = float64(paidOrders) / float64(totalOrders) * 100
	}

	return &types.BrandStatsResp{
		BrandId:          int64(brand.Id),
		TotalCampaigns:   totalCampaigns,
		ActiveCampaigns:  activeCampaigns,
		TotalOrders:      totalOrders,
		TotalRevenue:     totalRevenue,
		TotalRewards:     totalRewards,
		ParticipantCount: participantCount,
		ConversionRate:   conversionRate,
		LastUpdated:      time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
