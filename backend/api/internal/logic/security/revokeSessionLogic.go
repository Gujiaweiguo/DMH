// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRevokeSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeSessionLogic {
	return &RevokeSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeSessionLogic) RevokeSession() (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
