// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandDistributorsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandDistributorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandDistributorsLogic {
	return &GetBrandDistributorsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandDistributorsLogic) GetBrandDistributors(req *types.GetDistributorsReq) (resp *types.DistributorListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
