package member

import (
	"context"
	"fmt"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type MergeMemberPreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMergeMemberPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MergeMemberPreviewLogic {
	return &MergeMemberPreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MergeMemberPreviewLogic) MergeMemberPreview(req *types.MemberMergeReq) (*types.MemberMergePreviewResp, error) {
	// 只有平台管理员可以合并会员
	userId, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	if user.Role != "platform_admin" {
		return nil, fmt.Errorf("只有平台管理员可以合并会员")
	}

	// 查询源会员和目标会员
	var sourceMember, targetMember model.Member
	if err := l.svcCtx.DB.Where("id = ?", req.SourceMemberId).First(&sourceMember).Error; err != nil {
		return nil, fmt.Errorf("源会员不存在: %w", err)
	}
	if err := l.svcCtx.DB.Where("id = ?", req.TargetMemberId).First(&targetMember).Error; err != nil {
		return nil, fmt.Errorf("目标会员不存在: %w", err)
	}

	// 检查冲突
	conflicts := []string{}
	canMerge := true

	// UnionID 冲突检查
	if sourceMember.UnionID != "" && targetMember.UnionID != "" && sourceMember.UnionID != targetMember.UnionID {
		conflicts = append(conflicts, fmt.Sprintf("UnionID 冲突: 源会员(%s) vs 目标会员(%s)", sourceMember.UnionID, targetMember.UnionID))
		canMerge = false
	}

	// 手机号冲突检查
	if sourceMember.Phone != "" && targetMember.Phone != "" && sourceMember.Phone != targetMember.Phone {
		conflicts = append(conflicts, fmt.Sprintf("手机号冲突: 源会员(%s) vs 目标会员(%s)", sourceMember.Phone, targetMember.Phone))
	}

	// 查询关联数据统计
	var sourceOrderCount, targetOrderCount int64
	l.svcCtx.DB.Model(&model.Order{}).Where("member_id = ?", req.SourceMemberId).Count(&sourceOrderCount)
	l.svcCtx.DB.Model(&model.Order{}).Where("member_id = ?", req.TargetMemberId).Count(&targetOrderCount)

	if sourceOrderCount > 0 {
		conflicts = append(conflicts, fmt.Sprintf("源会员有 %d 个订单将被迁移", sourceOrderCount))
	}

	// 获取会员详情
	sourceLogic := NewGetMemberLogic(l.ctx, l.svcCtx)
	targetLogic := NewGetMemberLogic(l.ctx, l.svcCtx)

	sourceResp, err := sourceLogic.GetMember(req.SourceMemberId)
	if err != nil {
		return nil, err
	}

	targetResp, err := targetLogic.GetMember(req.TargetMemberId)
	if err != nil {
		return nil, err
	}

	return &types.MemberMergePreviewResp{
		SourceMember: *sourceResp,
		TargetMember: *targetResp,
		Conflicts:    conflicts,
		CanMerge:     canMerge,
	}, nil
}
