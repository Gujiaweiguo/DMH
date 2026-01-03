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

type GetMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenusLogic {
	return &GetMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenusLogic) GetMenus() (resp []types.MenuResp, err error) {
	var menus []model.Menu
	
	err = l.svcCtx.DB.Where("status = ?", "active").Order("sort ASC, id ASC").Find(&menus).Error
	if err != nil {
		l.Logger.Errorf("查询菜单列表失败: %v", err)
		return nil, fmt.Errorf("查询菜单列表失败")
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
	
	return rootMenus, nil
}
