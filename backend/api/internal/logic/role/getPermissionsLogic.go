// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPermissionsLogic {
	return &GetPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPermissionsLogic) GetPermissions() (resp []types.PermissionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
