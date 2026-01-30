// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigRolePermissionsLogic {
	return &ConfigRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigRolePermissionsLogic) ConfigRolePermissions(req *types.RolePermissionReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
