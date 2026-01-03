// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"errors"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.AdminUpdateUserReq, userId int64) (resp *types.UserInfoResp, err error) {

	// 查找用户
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		l.Logger.Errorf("查询用户失败: %v", err)
		return nil, fmt.Errorf("查询用户失败")
	}

	// 开启事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新用户信息
	updates := make(map[string]interface{})
	
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Status != "" {
		// 验证状态值
		validStatuses := map[string]bool{
			"active":   true,
			"disabled": true,
			"locked":   true,
		}
		if !validStatuses[req.Status] {
			tx.Rollback()
			return nil, fmt.Errorf("无效的用户状态: %s", req.Status)
		}
		updates["status"] = req.Status
	}
	if req.Role != "" {
		// 验证角色
		validRoles := map[string]bool{
			"platform_admin": true,
			"brand_admin":    true,
			"participant":    true,
		}
		if !validRoles[req.Role] {
			tx.Rollback()
			return nil, fmt.Errorf("无效的用户角色: %s", req.Role)
		}
		updates["role"] = req.Role
	}

	// 更新用户基本信息
	if len(updates) > 0 {
		if err := tx.Model(&user).Updates(updates).Error; err != nil {
			tx.Rollback()
			l.Logger.Errorf("更新用户信息失败: %v", err)
			return nil, fmt.Errorf("更新用户信息失败")
		}
	}

	// 处理品牌关联更新
	if req.BrandIds != nil {
		// 删除现有的品牌关联
		if err := tx.Where("user_id = ?", userId).Delete(&model.UserBrand{}).Error; err != nil {
			tx.Rollback()
			l.Logger.Errorf("删除用户品牌关联失败: %v", err)
			return nil, fmt.Errorf("更新用户品牌关联失败")
		}

		// 创建新的品牌关联
		for _, brandId := range req.BrandIds {
			// 验证品牌是否存在
			var brand model.Brand
			if err := tx.Where("id = ?", brandId).First(&brand).Error; err != nil {
				tx.Rollback()
				l.Logger.Errorf("品牌不存在: %d", brandId)
				return nil, fmt.Errorf("品牌不存在: %d", brandId)
			}

			userBrand := model.UserBrand{
				UserId:  userId,
				BrandId: brandId,
			}
			if err := tx.Create(&userBrand).Error; err != nil {
				tx.Rollback()
				l.Logger.Errorf("创建用户品牌关联失败: %v", err)
				return nil, fmt.Errorf("创建用户品牌关联失败")
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交事务失败: %v", err)
		return nil, fmt.Errorf("更新用户失败")
	}

	// 重新查询用户信息
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		l.Logger.Errorf("查询更新后用户信息失败: %v", err)
		return nil, fmt.Errorf("查询用户信息失败")
	}

	// 获取用户品牌关联
	var userBrands []model.UserBrand
	l.svcCtx.DB.Where("user_id = ?", userId).Find(&userBrands)
	
	brandIds := make([]int64, len(userBrands))
	for i, ub := range userBrands {
		brandIds[i] = ub.BrandId
	}

	l.Logger.Infof("管理员更新用户成功: %s", user.Username)

	return &types.UserInfoResp{
		Id:        user.Id,
		Username:  user.Username,
		Phone:     user.Phone,
		Email:     user.Email,
		RealName:  user.RealName,
		Status:    user.Status,
		Roles:     []string{user.Role},
		BrandIds:  brandIds,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
