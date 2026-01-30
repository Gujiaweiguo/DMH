// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManageBrandAdminRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManageBrandAdminRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageBrandAdminRelationLogic {
	return &ManageBrandAdminRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManageBrandAdminRelationLogic) ManageBrandAdminRelation(req *types.BrandAdminRelationReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
