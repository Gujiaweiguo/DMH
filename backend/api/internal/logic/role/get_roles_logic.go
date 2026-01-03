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

type GetRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolesLogic {
	return &GetRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRolesLogic) GetRoles() (resp []types.RoleResp, err error) {
	var roles []model.Role
	
	err = l.svcCtx.DB.Find(&roles).Error
	if err != nil {
		l.Logger.Errorf("查询角色列表失败: %v", err)
		return nil, fmt.Errorf("查询角色列表失败")
	}
	
	resp = make([]types.RoleResp, 0, len(roles))
	for _, role := range roles {
		// 查询角色权限
		var permissions []model.Permission
		err = l.svcCtx.DB.Table("permissions p").
			Joins("JOIN role_permissions rp ON rp.permission_id = p.id").
			Where("rp.role_id = ?", role.ID).
			Find(&permissions).Error
		
		var permissionCodes []string
		if err == nil {
			for _, perm := range permissions {
				permissionCodes = append(permissionCodes, perm.Code)
			}
		}
		
		resp = append(resp, types.RoleResp{
			Id:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Permissions: permissionCodes,
			CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	
	return resp, nil
}
