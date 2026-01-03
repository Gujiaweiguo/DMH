// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"fmt"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserMenusLogic {
	return &GetUserMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserMenusLogic) GetUserMenus() (resp *types.UserMenuResp, err error) {
	// 从context中获取用户信息
	userID, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		l.Logger.Errorf("获取用户ID失败: %v", err)
		return nil, fmt.Errorf("获取用户信息失败")
	}

	roles, err := middleware.GetUserRolesFromContext(l.ctx)
	if err != nil {
		l.Logger.Errorf("获取用户角色失败: %v", err)
		return nil, fmt.Errorf("获取用户角色失败")
	}

	// 获取平台参数（从查询参数或默认为admin）
	platform := "admin" // 默认为admin平台

	// 查询用户角色ID
	var roleIDs []int64
	if len(roles) > 0 {
		var roleModels []model.Role
		err = l.svcCtx.DB.Where("code IN ?", roles).Find(&roleModels).Error
		if err != nil {
			l.Logger.Errorf("查询角色信息失败: %v", err)
			return nil, fmt.Errorf("查询用户权限失败")
		}
		
		for _, role := range roleModels {
			roleIDs = append(roleIDs, role.ID)
		}
	}

	// 查询用户有权限的菜单
	var menus []model.Menu
	if len(roleIDs) > 0 {
		err = l.svcCtx.DB.Table("menus m").
			Joins("JOIN role_menus rm ON rm.menu_id = m.id").
			Where("rm.role_id IN ? AND m.platform = ? AND m.status = ?", roleIDs, platform, "active").
			Group("m.id").
			Order("m.sort ASC, m.id ASC").
			Find(&menus).Error
	} else {
		// 如果没有角色，返回空菜单
		menus = []model.Menu{}
	}

	if err != nil {
		l.Logger.Errorf("查询用户菜单失败: %v", err)
		return nil, fmt.Errorf("查询用户菜单失败")
	}

	// 构建菜单树结构
	menuMap := make(map[int64]*types.MenuResp)
	var rootMenus []types.MenuResp

	// 先创建所有菜单项
	for _, menu := range menus {
		menuResp := types.MenuResp{
			Id:        menu.ID,
			Name:      menu.Name,
			Code:      menu.Code,
			Path:      menu.Path,
			Icon:      menu.Icon,
			Sort:      menu.Sort,
			Type:      menu.Type,
			Platform:  menu.Platform,
			Status:    menu.Status,
			CreatedAt: menu.CreatedAt.Format("2006-01-02 15:04:05"),
			Children:  []types.MenuResp{},
		}

		if menu.ParentID != nil {
			menuResp.ParentId = *menu.ParentID
		}

		menuMap[menu.ID] = &menuResp
	}

	// 构建树结构
	for _, menu := range menus {
		menuResp := menuMap[menu.ID]
		if menu.ParentID == nil {
			// 根菜单
			rootMenus = append(rootMenus, *menuResp)
		} else {
			// 子菜单
			if parent, exists := menuMap[*menu.ParentID]; exists {
				parent.Children = append(parent.Children, *menuResp)
			}
		}
	}

	return &types.UserMenuResp{
		UserId:   userID,
		Platform: platform,
		Menus:    rootMenus,
	}, nil
}
