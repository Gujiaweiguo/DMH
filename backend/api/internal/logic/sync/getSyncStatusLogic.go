// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sync

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSyncStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncStatusLogic {
	return &GetSyncStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncStatusLogic) GetSyncStatus() (resp *types.SyncStatusResp, err error) {
	// todo: add your logic here and delete this line

	return
}
