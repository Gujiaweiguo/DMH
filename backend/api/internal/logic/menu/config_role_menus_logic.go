// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigRoleMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigRoleMenusLogic {
	return &ConfigRoleMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigRoleMenusLogic) ConfigRoleMenus(req *types.RoleMenuReq) (resp *types.CommonResp, err error) {
	// 开启事务
	tx := l.svcCtx.DB.Begin()
	if tx.Error != nil {
		l.Logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, fmt.Errorf("配置角色菜单失败")
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

	// 删除现有的角色菜单关联
	err = tx.Where("role_id = ?", req.RoleId).Delete(&model.RoleMenu{}).Error
	if err != nil {
		tx.Rollback()
		l.Logger.Errorf("删除现有角色菜单失败: %v", err)
		return nil, fmt.Errorf("配置角色菜单失败")
	}

	// 添加新的角色菜单关联
	for _, menuId := range req.MenuIds {
		// 验证菜单是否存在
		var menu model.Menu
		err = tx.Where("id = ?", menuId).First(&menu).Error
		if err != nil {
			tx.Rollback()
			l.Logger.Errorf("菜单不存在: %v", err)
			return nil, fmt.Errorf("菜单ID %d 不存在", menuId)
		}

		// 创建角色菜单关联
		roleMenu := &model.RoleMenu{
			RoleID: req.RoleId,
			MenuID: menuId,
		}
		err = tx.Create(roleMenu).Error
		if err != nil {
			tx.Rollback()
			l.Logger.Errorf("创建角色菜单关联失败: %v", err)
			return nil, fmt.Errorf("配置角色菜单失败")
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交事务失败: %v", err)
		return nil, fmt.Errorf("配置角色菜单失败")
	}

	l.Logger.Infof("角色 %d 菜单权限配置成功", req.RoleId)

	return &types.CommonResp{
		Message: "角色菜单权限配置成功",
	}, nil
}
