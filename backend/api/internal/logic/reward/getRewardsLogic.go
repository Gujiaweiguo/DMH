// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package reward

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRewardsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRewardsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRewardsLogic {
	return &GetRewardsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRewardsLogic) GetRewards() (resp []types.RewardResp, err error) {
	// todo: add your logic here and delete this line

	return
}
