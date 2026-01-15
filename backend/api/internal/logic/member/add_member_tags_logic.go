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

type AddMemberTagsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddMemberTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMemberTagsLogic {
	return &AddMemberTagsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMemberTagsLogic) AddMemberTags(req *types.AddMemberTagsReq) error {
	// 只有平台管理员可以添加标签
	userId, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		return err
	}
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	if user.Role != "platform_admin" {
		return fmt.Errorf("只有平台管理员可以添加标签")
	}

	// 检查会员是否存在
	var member model.Member
	if err := l.svcCtx.DB.Where("id = ?", req.MemberId).First(&member).Error; err != nil {
		return fmt.Errorf("会员不存在: %w", err)
	}

	// 批量添加标签
	for _, tagId := range req.TagIds {
		// 检查是否已存在
		var count int64
		l.svcCtx.DB.Model(&model.MemberTagLink{}).
			Where("member_id = ? AND tag_id = ?", req.MemberId, tagId).
			Count(&count)

		if count == 0 {
			tagLink := model.MemberTagLink{
				MemberID:  req.MemberId,
				TagID:     tagId,
				CreatedBy: userId,
			}
			if err := l.svcCtx.DB.Create(&tagLink).Error; err != nil {
				return fmt.Errorf("添加标签失败: %w", err)
			}
		}
	}

	return nil
}
