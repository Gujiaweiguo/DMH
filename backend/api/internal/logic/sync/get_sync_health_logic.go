// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sync

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncHealthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSyncHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncHealthLogic {
	return &GetSyncHealthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncHealthLogic) GetSyncHealth() (resp *types.SyncHealthResp, err error) {
	// todo: add your logic here and delete this line

	return
}
