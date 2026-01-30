// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCampaignsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCampaignsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCampaignsLogic {
	return &GetCampaignsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCampaignsLogic) GetCampaigns(req *types.GetCampaignsReq) (resp *types.CampaignListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
