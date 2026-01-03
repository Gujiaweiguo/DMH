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

type ManageBrandAdminRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManageBrandAdminRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageBrandAdminRelationLogic {
	return &ManageBrandAdminRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManageBrandAdminRelationLogic) ManageBrandAdminRelation(req *types.BrandAdminRelationReq) (resp *types.CommonResp, err error) {
	// 参数验证
	if req.UserId <= 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 验证用户是否存在且为品牌管理员
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", req.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		l.Logger.Errorf("查询用户失败: %v", err)
		return nil, fmt.Errorf("查询用户失败")
	}

	if user.Role != "brand_admin" {
		return nil, fmt.Errorf("只能为品牌管理员分配品牌权限")
	}

	// 验证所有品牌是否存在
	if len(req.BrandIds) > 0 {
		var brandCount int64
		if err := l.svcCtx.DB.Model(&model.Brand{}).Where("id IN ?", req.BrandIds).Count(&brandCount).Error; err != nil {
			l.Logger.Errorf("验证品牌失败: %v", err)
			return nil, fmt.Errorf("验证品牌失败")
		}
		if brandCount != int64(len(req.BrandIds)) {
			return nil, fmt.Errorf("存在无效的品牌ID")
		}
	}

	// 开启事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除用户现有的品牌关联
	if err := tx.Where("user_id = ?", req.UserId).Delete(&model.UserBrand{}).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("删除用户品牌关联失败: %v", err)
		return nil, fmt.Errorf("更新品牌关联失败")
	}

	// 创建新的品牌关联
	for _, brandId := range req.BrandIds {
		userBrand := model.UserBrand{
			UserId:  req.UserId,
			BrandId: brandId,
		}
		if err := tx.Create(&userBrand).Error; err != nil {
			tx.Rollback()
			l.Logger.Errorf("创建用户品牌关联失败: %v", err)
			return nil, fmt.Errorf("创建品牌关联失败")
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交事务失败: %v", err)
		return nil, fmt.Errorf("更新品牌关联失败")
	}

	l.Logger.Infof("管理员更新品牌管理员关系成功: 用户ID=%d, 品牌数量=%d", req.UserId, len(req.BrandIds))

	return &types.CommonResp{
		Message: "品牌管理员关系更新成功",
	}, nil
}
