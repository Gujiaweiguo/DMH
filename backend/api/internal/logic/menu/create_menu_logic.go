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

type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuLogic) CreateMenu(req *types.CreateMenuReq) (resp *types.MenuResp, err error) {
	// 验证菜单编码唯一性
	var count int64
	err = l.svcCtx.DB.Model(&model.Menu{}).Where("code = ?", req.Code).Count(&count).Error
	if err != nil {
		l.Logger.Errorf("检查菜单编码唯一性失败: %v", err)
		return nil, fmt.Errorf("创建菜单失败")
	}
	if count > 0 {
		return nil, fmt.Errorf("菜单编码已存在")
	}

	// 验证父菜单是否存在
	if req.ParentId > 0 {
		var parentMenu model.Menu
		err = l.svcCtx.DB.Where("id = ?", req.ParentId).First(&parentMenu).Error
		if err != nil {
			l.Logger.Errorf("父菜单不存在: %v", err)
			return nil, fmt.Errorf("父菜单不存在")
		}
	}

	// 创建菜单
	menu := &model.Menu{
		Name:     req.Name,
		Code:     req.Code,
		Path:     req.Path,
		Icon:     req.Icon,
		Sort:     req.Sort,
		Type:     req.Type,
		Platform: req.Platform,
		Status:   "active",
	}

	if req.ParentId > 0 {
		menu.ParentID = &req.ParentId
	}

	err = l.svcCtx.DB.Create(menu).Error
	if err != nil {
		l.Logger.Errorf("创建菜单失败: %v", err)
		return nil, fmt.Errorf("创建菜单失败")
	}

	l.Logger.Infof("菜单创建成功: ID=%d, Code=%s", menu.ID, menu.Code)

	// 构建响应
	resp = &types.MenuResp{
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
		resp.ParentId = *menu.ParentID
	}

	return resp, nil
}
