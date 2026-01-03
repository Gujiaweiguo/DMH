// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sync

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSyncStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncStatsLogic {
	return &GetSyncStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncStatsLogic) GetSyncStats() (resp *types.SyncStatsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
