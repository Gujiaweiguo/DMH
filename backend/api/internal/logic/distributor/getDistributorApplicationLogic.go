// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDistributorApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorApplicationLogic {
	return &GetDistributorApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorApplicationLogic) GetDistributorApplication() (resp *types.DistributorApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
