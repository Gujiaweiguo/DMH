// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DistributorApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDistributorApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DistributorApplyLogic {
	return &DistributorApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DistributorApplyLogic) DistributorApply(req *types.DistributorApplyReq) (resp *types.DistributorApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
