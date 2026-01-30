// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDistributorLevelRewardsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorLevelRewardsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorLevelRewardsLogic {
	return &GetDistributorLevelRewardsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorLevelRewardsLogic) GetDistributorLevelRewards() (resp *types.DistributorLevelRewardsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
