// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package withdrawal

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyWithdrawalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyWithdrawalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyWithdrawalLogic {
	return &ApplyWithdrawalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyWithdrawalLogic) ApplyWithdrawal(req *types.WithdrawalApplyReq) (resp *types.WithdrawalResp, err error) {
	// todo: add your logic here and delete this line

	return
}
