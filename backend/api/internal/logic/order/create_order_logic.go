package order

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.OrderResp, err error) {
	// 查询活动信息
	var campaign model.Campaign
	if err := l.svcCtx.DB.Where("id = ?", req.CampaignId).First(&campaign).Error; err != nil {
		return nil, fmt.Errorf("活动不存在: %w", err)
	}

	// 检查活动状态和时间
	now := time.Now()
	if campaign.Status != "active" {
		return nil, fmt.Errorf("活动未激活")
	}
	if now.Before(campaign.StartTime) {
		return nil, fmt.Errorf("活动未开始")
	}
	if now.After(campaign.EndTime) {
		return nil, fmt.Errorf("活动已结束")
	}

	// 将表单数据转换为 JSON
	formDataBytes, err := json.Marshal(req.FormData)
	if err != nil {
		return nil, fmt.Errorf("表单数据格式错误: %w", err)
	}

	// 处理会员关联
	var memberId *int64
	if req.UnionID != "" {
		// 查找或创建会员
		var member model.Member
		result := l.svcCtx.DB.Where("unionid = ?", req.UnionID).First(&member)
		
		if result.Error == gorm.ErrRecordNotFound {
			// 创建新会员
			member = model.Member{
				UnionID:  req.UnionID,
				Phone:    req.Phone,
				Status:   "active",
				Source:   "campaign", // 可以根据实际情况设置来源
			}
			
			if err := l.svcCtx.DB.Create(&member).Error; err != nil {
				return nil, fmt.Errorf("创建会员失败: %w", err)
			}

			// 创建会员画像
			profile := model.MemberProfile{
				MemberID: member.ID,
			}
			l.svcCtx.DB.Create(&profile)

			// 创建会员品牌关联
			brandLink := model.MemberBrandLink{
				MemberID:        member.ID,
				BrandID:         campaign.BrandId,
				FirstCampaignID: req.CampaignId,
			}
			l.svcCtx.DB.Create(&brandLink)
		} else if result.Error != nil {
			return nil, fmt.Errorf("查询会员失败: %w", result.Error)
		} else {
			// 会员已存在，检查品牌关联
			var count int64
			l.svcCtx.DB.Model(&model.MemberBrandLink{}).
				Where("member_id = ? AND brand_id = ?", member.ID, campaign.BrandId).
				Count(&count)
			
			if count == 0 {
				// 创建品牌关联
				brandLink := model.MemberBrandLink{
					MemberID:        member.ID,
					BrandID:         campaign.BrandId,
					FirstCampaignID: req.CampaignId,
				}
				l.svcCtx.DB.Create(&brandLink)
			}
		}

		memberId = &member.ID
	}

	// 创建订单
	order := model.Order{
		CampaignId: req.CampaignId,
		MemberID:   memberId,
		UnionID:    req.UnionID,
		Phone:      req.Phone,
		FormData:   string(formDataBytes),
		ReferrerId: req.ReferrerId,
		Status:     "pending",
		Amount:     req.Amount,
		PayStatus:  "unpaid",
		SyncStatus: "pending",
	}

	if err := l.svcCtx.DB.Create(&order).Error; err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	// 更新会员画像（如果有会员）
	if memberId != nil {
		l.updateMemberProfile(*memberId, order.Id)
	}

	return &types.OrderResp{
		Id:         order.Id,
		CampaignId: order.CampaignId,
		CampaignName: campaign.Name,
		Phone:      order.Phone,
		Status:     order.Status,
		Amount:     order.Amount,
		CreatedAt:  order.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (l *CreateOrderLogic) updateMemberProfile(memberId int64, orderId int64) {
	var profile model.MemberProfile
	if err := l.svcCtx.DB.Where("member_id = ?", memberId).First(&profile).Error; err != nil {
		return
	}

	// 更新订单数
	profile.TotalOrders++

	// 更新首次/最后下单时间
	now := time.Now()
	if profile.FirstOrderAt == nil {
		profile.FirstOrderAt = &now
	}
	profile.LastOrderAt = &now

	// 更新参与活动数（去重）
	var campaignCount int64
	l.svcCtx.DB.Model(&model.Order{}).
		Where("member_id = ?", memberId).
		Distinct("campaign_id").
		Count(&campaignCount)
	profile.ParticipatedCampaigns = int(campaignCount)

	l.svcCtx.DB.Save(&profile)
}
