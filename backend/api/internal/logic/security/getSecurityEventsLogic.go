// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSecurityEventsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSecurityEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSecurityEventsLogic {
	return &GetSecurityEventsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSecurityEventsLogic) GetSecurityEvents() (resp *types.SecurityEventListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
