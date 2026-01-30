// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDistributorQrcodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorQrcodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorQrcodeLogic {
	return &GetDistributorQrcodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorQrcodeLogic) GetDistributorQrcode() (resp *types.GetQrcodeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
