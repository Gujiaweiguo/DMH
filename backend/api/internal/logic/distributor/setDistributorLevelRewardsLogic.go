// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetDistributorLevelRewardsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetDistributorLevelRewardsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDistributorLevelRewardsLogic {
	return &SetDistributorLevelRewardsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetDistributorLevelRewardsLogic) SetDistributorLevelRewards(req *types.SetDistributorLevelRewardsReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
