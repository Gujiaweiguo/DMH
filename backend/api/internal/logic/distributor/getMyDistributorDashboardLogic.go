// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyDistributorDashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyDistributorDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyDistributorDashboardLogic {
	return &GetMyDistributorDashboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyDistributorDashboardLogic) GetMyDistributorDashboard() (resp *types.DistributorStatisticsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
