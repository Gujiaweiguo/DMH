// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPermissionsLogic {
	return &GetUserPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPermissionsLogic) GetUserPermissions() (resp *types.UserPermissionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
