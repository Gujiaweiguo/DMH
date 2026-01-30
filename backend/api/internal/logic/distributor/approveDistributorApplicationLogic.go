// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveDistributorApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveDistributorApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveDistributorApplicationLogic {
	return &ApproveDistributorApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveDistributorApplicationLogic) ApproveDistributorApplication(req *types.ApproveDistributorReq) (resp *types.DistributorApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
