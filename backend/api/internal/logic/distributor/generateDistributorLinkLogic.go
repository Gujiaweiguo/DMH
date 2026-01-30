// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateDistributorLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateDistributorLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateDistributorLinkLogic {
	return &GenerateDistributorLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateDistributorLinkLogic) GenerateDistributorLink(req *types.GenerateLinkReq) (resp *types.GenerateLinkResp, err error) {
	// todo: add your logic here and delete this line

	return
}
