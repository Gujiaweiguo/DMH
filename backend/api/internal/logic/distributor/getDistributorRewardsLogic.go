// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDistributorRewardsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorRewardsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorRewardsLogic {
	return &GetDistributorRewardsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorRewardsLogic) GetDistributorRewards(req *types.GetDistributorRewardsReq) (resp *types.DistributorRewardListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
