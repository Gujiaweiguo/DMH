// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandStatsLogic {
	return &GetBrandStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandStatsLogic) GetBrandStats() (resp *types.BrandStatsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
