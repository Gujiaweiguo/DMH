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

type UpdateMemberStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMemberStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMemberStatusLogic {
	return &UpdateMemberStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMemberStatusLogic) UpdateMemberStatus(memberId int64, req *types.UpdateMemberStatusReq) (resp *types.CommonResp, err error) {
	var member model.Member

	if err := l.svcCtx.DB.Where("id = ? AND deleted_at IS NULL", memberId).First(&member).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, fmt.Errorf("Member not found")
		}
		l.Errorf("Failed to query member: %v", err)
		return nil, fmt.Errorf("Failed to query member: %w", err)
	}

	member.Status = req.Status

	if err := l.svcCtx.DB.Save(&member).Error; err != nil {
		l.Errorf("Failed to update member status: %v", err)
		return nil, fmt.Errorf("Failed to update member status: %w", err)
	}

	resp = &types.CommonResp{
		Message: "Member status updated successfully",
	}

	l.Infof("Successfully updated member status: id=%d, status=%s", memberId, req.Status)
	return resp, nil
}
