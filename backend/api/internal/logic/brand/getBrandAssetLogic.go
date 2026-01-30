// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandAssetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandAssetLogic {
	return &GetBrandAssetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandAssetLogic) GetBrandAsset() (resp *types.BrandAssetResp, err error) {
	// todo: add your logic here and delete this line

	return
}
