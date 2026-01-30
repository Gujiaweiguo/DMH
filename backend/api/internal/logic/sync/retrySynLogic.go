// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sync

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RetrySynLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRetrySynLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RetrySynLogic {
	return &RetrySynLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RetrySynLogic) RetrySyn() (resp *types.RetrySyncResp, err error) {
	// todo: add your logic here and delete this line

	return
}
