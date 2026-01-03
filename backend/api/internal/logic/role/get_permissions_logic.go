// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

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
	var permissions []model.Permission
	
	err = l.svcCtx.DB.Find(&permissions).Error
	if err != nil {
		l.Logger.Errorf("查询权限列表失败: %v", err)
		return nil, fmt.Errorf("查询权限列表失败")
	}
	
	resp = make([]types.PermissionResp, 0, len(permissions))
	for _, perm := range permissions {
		resp = append(resp, types.PermissionResp{
			Id:          perm.ID,
			Name:        perm.Name,
			Code:        perm.Code,
			Resource:    perm.Resource,
			Action:      perm.Action,
			Description: perm.Description,
		})
	}
	
	return resp, nil
}
