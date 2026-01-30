// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDistributorByCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorByCodeLogic {
	return &GetDistributorByCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorByCodeLogic) GetDistributorByCode() (resp *types.DistributorResp, err error) {
	// todo: add your logic here and delete this line

	return
}
