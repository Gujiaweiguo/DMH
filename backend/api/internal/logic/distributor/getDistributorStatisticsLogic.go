// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDistributorStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorStatisticsLogic {
	return &GetDistributorStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorStatisticsLogic) GetDistributorStatistics() (resp *types.DistributorStatisticsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
