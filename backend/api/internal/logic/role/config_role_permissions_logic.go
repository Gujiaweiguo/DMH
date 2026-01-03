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
	// 开启事务
	tx := l.svcCtx.DB.Begin()
	if tx.Error != nil {
		l.Logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, fmt.Errorf("配置角色权限失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 验证角色是否存在
	var role model.Role
	err = tx.Where("id = ?", req.RoleId).First(&role).Error
	if err != nil {
		tx.Rollback()
		l.Logger.Errorf("角色不存在: %v", err)
		return nil, fmt.Errorf("角色不存在")
	}

	// 删除现有的角色权限关联
	err = tx.Where("role_id = ?", req.RoleId).Delete(&model.RolePermission{}).Error
	if err != nil {
		tx.Rollback()
		l.Logger.Errorf("删除现有角色权限失败: %v", err)
		return nil, fmt.Errorf("配置角色权限失败")
	}

	// 添加新的角色权限关联
	for _, permissionId := range req.PermissionIds {
		// 验证权限是否存在
		var permission model.Permission
		err = tx.Where("id = ?", permissionId).First(&permission).Error
		if err != nil {
			tx.Rollback()
			l.Logger.Errorf("权限不存在: %v", err)
			return nil, fmt.Errorf("权限ID %d 不存在", permissionId)
		}

		// 创建角色权限关联
		rolePermission := &model.RolePermission{
			RoleID:       req.RoleId,
			PermissionID: permissionId,
		}
		err = tx.Create(rolePermission).Error
		if err != nil {
			tx.Rollback()
			l.Logger.Errorf("创建角色权限关联失败: %v", err)
			return nil, fmt.Errorf("配置角色权限失败")
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交事务失败: %v", err)
		return nil, fmt.Errorf("配置角色权限失败")
	}

	l.Logger.Infof("角色 %d 权限配置成功", req.RoleId)

	return &types.CommonResp{
		Message: "角色权限配置成功",
	}, nil
}
