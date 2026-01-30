// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoginAttemptsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLoginAttemptsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginAttemptsLogic {
	return &GetLoginAttemptsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoginAttemptsLogic) GetLoginAttempts() (resp *types.LoginAttemptListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
