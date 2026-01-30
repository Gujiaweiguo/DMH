// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBrandAssetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBrandAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBrandAssetLogic {
	return &DeleteBrandAssetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBrandAssetLogic) DeleteBrandAsset() (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
