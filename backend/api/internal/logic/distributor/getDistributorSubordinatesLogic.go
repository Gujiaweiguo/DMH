// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDistributorSubordinatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorSubordinatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorSubordinatesLogic {
	return &GetDistributorSubordinatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorSubordinatesLogic) GetDistributorSubordinates() (resp *types.SubordinateListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
