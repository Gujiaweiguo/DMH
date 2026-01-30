// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandDistributorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandDistributorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandDistributorLogic {
	return &GetBrandDistributorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandDistributorLogic) GetBrandDistributor() (resp *types.DistributorResp, err error) {
	// todo: add your logic here and delete this line

	return
}
