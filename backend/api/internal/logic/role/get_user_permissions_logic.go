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

func (l *GetUserPermissionsLogic) GetUserPermissions(userId int64) (resp *types.UserPermissionResp, err error) {
	// 验证用户是否存在
	var user model.User
	err = l.svcCtx.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		l.Logger.Errorf("用户不存在: %v", err)
		return nil, fmt.Errorf("用户不存在")
	}

	// 查询用户角色
	var roles []model.Role
	err = l.svcCtx.DB.Table("roles r").
		Joins("JOIN user_roles ur ON ur.role_id = r.id").
		Where("ur.user_id = ?", userId).
		Find(&roles).Error
	if err != nil {
		l.Logger.Errorf("查询用户角色失败: %v", err)
		return nil, fmt.Errorf("查询用户权限失败")
	}

	var roleCodes []string
	var roleIds []int64
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
		roleIds = append(roleIds, role.ID)
	}

	// 查询用户权限
	var permissions []model.Permission
	if len(roleIds) > 0 {
		err = l.svcCtx.DB.Table("permissions p").
			Joins("JOIN role_permissions rp ON rp.permission_id = p.id").
			Where("rp.role_id IN ?", roleIds).
			Group("p.id").
			Find(&permissions).Error
		if err != nil {
			l.Logger.Errorf("查询用户权限失败: %v", err)
			return nil, fmt.Errorf("查询用户权限失败")
		}
	}

	var permissionCodes []string
	for _, perm := range permissions {
		permissionCodes = append(permissionCodes, perm.Code)
	}

	// 查询品牌管理员的品牌ID
	var brandIDs []int64
	if l.containsRole(roleCodes, "brand_admin") {
		var userBrands []model.UserBrand
		err = l.svcCtx.DB.Where("user_id = ?", userId).Find(&userBrands).Error
		if err == nil {
			for _, ub := range userBrands {
				brandIDs = append(brandIDs, ub.BrandId)
			}
		}
	}

	return &types.UserPermissionResp{
		UserId:      userId,
		Roles:       roleCodes,
		Permissions: permissionCodes,
		BrandIds:    brandIDs,
	}, nil
}

// containsRole 检查是否包含指定角色
func (l *GetUserPermissionsLogic) containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
