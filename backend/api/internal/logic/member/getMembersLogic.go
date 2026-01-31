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

type GetMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMembersLogic {
	return &GetMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMembersLogic) GetMembers(req *types.GetMembersReq) (resp *types.GetMembersResp, err error) {
	var members []model.Member
	query := l.svcCtx.DB.Model(&model.Member{}).Where("deleted_at IS NULL")

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		query = query.Where("nickname LIKE ? OR phone LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Source != "" {
		query = query.Where("source = ?", req.Source)
	}
	if req.Gender > 0 {
		query = query.Where("gender = ?", req.Gender)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		l.Errorf("Failed to count members: %v", err)
		return nil, fmt.Errorf("Failed to count members: %w", err)
	}

	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(int(offset)).Limit(int(req.PageSize))
	}

	if err := query.Find(&members).Error; err != nil {
		l.Errorf("Failed to query members: %v", err)
		return nil, fmt.Errorf("Failed to query members: %w", err)
	}

	var memberResps []types.MemberResp
	for _, member := range members {
		memberResps = append(memberResps, types.MemberResp{
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
		})
	}

	l.Infof("Successfully queried members: count=%d", len(members))

	resp = &types.GetMembersResp{
		Total:   total,
		Members: memberResps,
	}

	return resp, nil
}
