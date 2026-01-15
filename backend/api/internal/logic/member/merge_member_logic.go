package member

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type MergeMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMergeMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MergeMemberLogic {
	return &MergeMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MergeMemberLogic) MergeMember(req *types.MemberMergeReq) error {
	// 只有平台管理员可以合并会员
	userId, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		return err
	}
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	if user.Role != "platform_admin" {
		return fmt.Errorf("只有平台管理员可以合并会员")
	}

	// 先进行预检查
	previewLogic := NewMergeMemberPreviewLogic(l.ctx, l.svcCtx)
	preview, err := previewLogic.MergeMemberPreview(req)
	if err != nil {
		return err
	}

	if !preview.CanMerge {
		return fmt.Errorf("会员存在冲突，无法合并")
	}

	// 创建合并请求记录
	conflictInfoBytes, _ := json.Marshal(preview.Conflicts)
	mergeRequest := model.MemberMergeRequest{
		SourceMemberID: req.SourceMemberId,
		TargetMemberID: req.TargetMemberId,
		Status:         "pending",
		Reason:         req.Reason,
		ConflictInfo:   string(conflictInfoBytes),
		CreatedBy:      userId,
	}

	if err := l.svcCtx.DB.Create(&mergeRequest).Error; err != nil {
		return fmt.Errorf("创建合并请求失败: %w", err)
	}

	// 执行合并（使用事务）
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 迁移订单
		if err := tx.Model(&model.Order{}).
			Where("member_id = ?", req.SourceMemberId).
			Update("member_id", req.TargetMemberId).Error; err != nil {
			return fmt.Errorf("迁移订单失败: %w", err)
		}

		// 2. 迁移奖励
		if err := tx.Model(&model.Reward{}).
			Where("member_id = ?", req.SourceMemberId).
			Update("member_id", req.TargetMemberId).Error; err != nil {
			return fmt.Errorf("迁移奖励失败: %w", err)
		}

		// 3. 迁移品牌关联（去重）
		var sourceBrandLinks []model.MemberBrandLink
		if err := tx.Where("member_id = ?", req.SourceMemberId).Find(&sourceBrandLinks).Error; err != nil {
			return fmt.Errorf("查询源会员品牌关联失败: %w", err)
		}

		for _, link := range sourceBrandLinks {
			// 检查目标会员是否已有该品牌关联
			var count int64
			tx.Model(&model.MemberBrandLink{}).
				Where("member_id = ? AND brand_id = ?", req.TargetMemberId, link.BrandID).
				Count(&count)

			if count == 0 {
				// 不存在则迁移
				if err := tx.Model(&model.MemberBrandLink{}).
					Where("id = ?", link.ID).
					Update("member_id", req.TargetMemberId).Error; err != nil {
					return fmt.Errorf("迁移品牌关联失败: %w", err)
				}
			} else {
				// 已存在则删除源关联
				if err := tx.Delete(&model.MemberBrandLink{}, link.ID).Error; err != nil {
					return fmt.Errorf("删除重复品牌关联失败: %w", err)
				}
			}
		}

		// 4. 迁移标签（去重）
		var sourceTagLinks []model.MemberTagLink
		if err := tx.Where("member_id = ?", req.SourceMemberId).Find(&sourceTagLinks).Error; err != nil {
			return fmt.Errorf("查询源会员标签失败: %w", err)
		}

		for _, link := range sourceTagLinks {
			var count int64
			tx.Model(&model.MemberTagLink{}).
				Where("member_id = ? AND tag_id = ?", req.TargetMemberId, link.TagID).
				Count(&count)

			if count == 0 {
				if err := tx.Model(&model.MemberTagLink{}).
					Where("id = ?", link.ID).
					Update("member_id", req.TargetMemberId).Error; err != nil {
					return fmt.Errorf("迁移标签失败: %w", err)
				}
			} else {
				if err := tx.Delete(&model.MemberTagLink{}, link.ID).Error; err != nil {
					return fmt.Errorf("删除重复标签失败: %w", err)
				}
			}
		}

		// 5. 合并会员画像数据
		var sourceProfile, targetProfile model.MemberProfile
		tx.Where("member_id = ?", req.SourceMemberId).First(&sourceProfile)
		tx.Where("member_id = ?", req.TargetMemberId).First(&targetProfile)

		// 累加统计数据
		targetProfile.TotalOrders += sourceProfile.TotalOrders
		targetProfile.TotalPayment += sourceProfile.TotalPayment
		targetProfile.TotalReward += sourceProfile.TotalReward
		targetProfile.ParticipatedCampaigns += sourceProfile.ParticipatedCampaigns

		// 更新首次/最后时间
		if sourceProfile.FirstOrderAt != nil {
			if targetProfile.FirstOrderAt == nil || sourceProfile.FirstOrderAt.Before(*targetProfile.FirstOrderAt) {
				targetProfile.FirstOrderAt = sourceProfile.FirstOrderAt
			}
		}
		if sourceProfile.LastOrderAt != nil {
			if targetProfile.LastOrderAt == nil || sourceProfile.LastOrderAt.After(*targetProfile.LastOrderAt) {
				targetProfile.LastOrderAt = sourceProfile.LastOrderAt
			}
		}
		if sourceProfile.FirstPaymentAt != nil {
			if targetProfile.FirstPaymentAt == nil || sourceProfile.FirstPaymentAt.Before(*targetProfile.FirstPaymentAt) {
				targetProfile.FirstPaymentAt = sourceProfile.FirstPaymentAt
			}
		}
		if sourceProfile.LastPaymentAt != nil {
			if targetProfile.LastPaymentAt == nil || sourceProfile.LastPaymentAt.After(*targetProfile.LastPaymentAt) {
				targetProfile.LastPaymentAt = sourceProfile.LastPaymentAt
			}
		}

		if err := tx.Save(&targetProfile).Error; err != nil {
			return fmt.Errorf("更新目标会员画像失败: %w", err)
		}

		// 6. 删除源会员画像
		if err := tx.Delete(&model.MemberProfile{}, "member_id = ?", req.SourceMemberId).Error; err != nil {
			return fmt.Errorf("删除源会员画像失败: %w", err)
		}

		// 7. 软删除源会员
		if err := tx.Delete(&model.Member{}, req.SourceMemberId).Error; err != nil {
			return fmt.Errorf("删除源会员失败: %w", err)
		}

		// 8. 更新合并请求状态
		now := time.Now()
		if err := tx.Model(&model.MemberMergeRequest{}).
			Where("id = ?", mergeRequest.ID).
			Updates(map[string]interface{}{
				"status":      "completed",
				"executed_at": now,
			}).Error; err != nil {
			return fmt.Errorf("更新合并请求状态失败: %w", err)
		}

		return nil
	})

	if err != nil {
		// 更新合并请求为失败状态
		l.svcCtx.DB.Model(&model.MemberMergeRequest{}).
			Where("id = ?", mergeRequest.ID).
			Updates(map[string]interface{}{
				"status":    "failed",
				"error_msg": err.Error(),
			})
		return err
	}

	return nil
}
