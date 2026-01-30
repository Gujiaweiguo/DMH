// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDistributorStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDistributorStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDistributorStatusLogic {
	return &UpdateDistributorStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDistributorStatusLogic) UpdateDistributorStatus(req *types.UpdateDistributorStatusReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
