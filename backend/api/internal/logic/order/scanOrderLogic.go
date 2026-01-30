// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScanOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScanOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScanOrderLogic {
	return &ScanOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScanOrderLogic) ScanOrder() (resp *types.ScanOrderResp, err error) {
	// todo: add your logic here and delete this line

	return
}
