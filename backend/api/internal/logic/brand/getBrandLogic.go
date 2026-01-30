// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandLogic {
	return &GetBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandLogic) GetBrand() (resp *types.BrandResp, err error) {
	// todo: add your logic here and delete this line

	return
}
