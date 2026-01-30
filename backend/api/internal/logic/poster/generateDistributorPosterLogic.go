// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package poster

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateDistributorPosterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateDistributorPosterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateDistributorPosterLogic {
	return &GenerateDistributorPosterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateDistributorPosterLogic) GenerateDistributorPoster(req *types.GeneratePosterReq) (resp *types.GeneratePosterResp, err error) {
	// todo: add your logic here and delete this line

	return
}
