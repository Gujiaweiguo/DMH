// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package withdrawal

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWithdrawalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWithdrawalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWithdrawalLogic {
	return &GetWithdrawalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWithdrawalLogic) GetWithdrawal() (resp *types.WithdrawalResp, err error) {
	// todo: add your logic here and delete this line

	return
}
