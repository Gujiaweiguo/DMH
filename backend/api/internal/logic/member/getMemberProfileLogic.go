// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package member

import (
	"context"
	"fmt"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMemberProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMemberProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMemberProfileLogic {
	return &GetMemberProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMemberProfileLogic) GetMemberProfile(req *types.GetMemberReq) (resp *types.MemberProfileResp, err error) {
	var profile model.MemberProfile

	if err := l.svcCtx.DB.Where("member_id = ?", req.Id).First(&profile).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, fmt.Errorf("Member profile not found")
		}
		l.Errorf("Failed to query member profile: %v", err)
		return nil, fmt.Errorf("Failed to query member profile: %w", err)
	}

	resp = &types.MemberProfileResp{
		MemberId:              profile.MemberID,
		TotalOrders:           profile.TotalOrders,
		TotalPayment:          profile.TotalPayment,
		TotalReward:           profile.TotalReward,
		FirstOrderAt:          timeToStr(profile.FirstOrderAt),
		LastOrderAt:           timeToStr(profile.LastOrderAt),
		FirstPaymentAt:        timeToStr(profile.FirstPaymentAt),
		LastPaymentAt:         timeToStr(profile.LastPaymentAt),
		ParticipatedCampaigns: profile.ParticipatedCampaigns,
	}

	l.Infof("Successfully queried member profile: member_id=%d", profile.MemberID)
	return resp, nil
}

func timeToStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02T15:04:05")
}
