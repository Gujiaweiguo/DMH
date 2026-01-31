// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package member

import (
	"context"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMemberLogic {
	return &UpdateMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMemberLogic) UpdateMember(memberId int64, req *types.UpdateMemberReq) (resp *types.MemberResp, err error) {
	var member model.Member

	if err := l.svcCtx.DB.Where("id = ? AND deleted_at IS NULL", memberId).First(&member).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, fmt.Errorf("Member not found")
		}
		l.Errorf("Failed to query member: %v", err)
		return nil, fmt.Errorf("Failed to query member: %w", err)
	}

	if req.Nickname != "" {
		member.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		member.Avatar = req.Avatar
	}
	if req.Gender > 0 {
		member.Gender = req.Gender
	}

	if err := l.svcCtx.DB.Save(&member).Error; err != nil {
		l.Errorf("Failed to update member: %v", err)
		return nil, fmt.Errorf("Failed to update member: %w", err)
	}

	resp = &types.MemberResp{
		Id:        member.ID,
		UnionID:   member.UnionID,
		Nickname:  member.Nickname,
		Avatar:    member.Avatar,
		Phone:     member.Phone,
		Gender:    member.Gender,
		Source:    member.Source,
		Status:    member.Status,
		CreatedAt: member.CreatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAt: member.UpdatedAt.Format("2006-01-02T15:04:05"),
	}

	l.Infof("Successfully updated member: id=%d", memberId)
	return resp, nil
}
