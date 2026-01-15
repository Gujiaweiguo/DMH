package order

import (
	"context"
	"fmt"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type PaymentCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentCallbackLogic {
	return &PaymentCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentCallbackLogic) PaymentCallback(req *types.PaymentCallbackReq) error {
	// 查询订单
	var order model.Order
	if err := l.svcCtx.DB.Where("id = ?", req.OrderId).First(&order).Error; err != nil {
		return fmt.Errorf("订单不存在: %w", err)
	}

	// 检查订单状态
	if order.PayStatus == "paid" {
		return fmt.Errorf("订单已支付")
	}

	// 查询活动信息
	var campaign model.Campaign
	if err := l.svcCtx.DB.Where("id = ?", order.CampaignId).First(&campaign).Error; err != nil {
		return fmt.Errorf("活动不存在: %w", err)
	}

	// 使用事务处理支付和奖励
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		order.PayStatus = "paid"
		order.Status = "paid"
		order.TradeNo = req.TradeNo
		if err := tx.Save(&order).Error; err != nil {
			return fmt.Errorf("更新订单状态失败: %w", err)
		}

		// 更新会员画像（如果有会员）
		if order.MemberID != nil {
			l.updateMemberPaymentProfile(tx, *order.MemberID, order.Amount)
		}

		// 如果有推荐人，创建奖励
		if order.ReferrerId > 0 && campaign.RewardRule > 0 {
			// 查询推荐人
			var referrer model.User
			if err := tx.Where("id = ?", order.ReferrerId).First(&referrer).Error; err != nil {
				logx.Errorf("查询推荐人失败: %v", err)
			} else {
				// 查询推荐人的会员ID
				var referrerMemberId *int64
				if order.MemberID != nil {
					// 尝试通过推荐人的手机号查找会员
					var referrerMember model.Member
					result := tx.Where("phone = ?", referrer.Phone).First(&referrerMember)
					if result.Error == nil {
						referrerMemberId = &referrerMember.ID
					}
				}

				// 创建奖励记录
				reward := model.Reward{
					UserId:     order.ReferrerId,
					MemberID:   referrerMemberId,
					OrderId:    order.Id,
					CampaignId: order.CampaignId,
					Amount:     campaign.RewardRule,
					Status:     "pending",
				}

				if err := tx.Create(&reward).Error; err != nil {
					return fmt.Errorf("创建奖励失败: %w", err)
				}

				// 更新推荐人会员画像（如果有）
				if referrerMemberId != nil {
					l.updateMemberRewardProfile(tx, *referrerMemberId, campaign.RewardRule)
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (l *PaymentCallbackLogic) updateMemberPaymentProfile(tx *gorm.DB, memberId int64, amount float64) {
	var profile model.MemberProfile
	if err := tx.Where("member_id = ?", memberId).First(&profile).Error; err != nil {
		return
	}

	// 更新支付金额
	profile.TotalPayment += amount

	// 更新首次/最后支付时间
	now := time.Now()
	if profile.FirstPaymentAt == nil {
		profile.FirstPaymentAt = &now
	}
	profile.LastPaymentAt = &now

	tx.Save(&profile)
}

func (l *PaymentCallbackLogic) updateMemberRewardProfile(tx *gorm.DB, memberId int64, amount float64) {
	var profile model.MemberProfile
	if err := tx.Where("member_id = ?", memberId).First(&profile).Error; err != nil {
		return
	}

	// 更新奖励金额
	profile.TotalReward += amount

	tx.Save(&profile)
}

