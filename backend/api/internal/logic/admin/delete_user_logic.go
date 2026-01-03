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

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(userId int64) (resp *types.CommonResp, err error) {

	// 查找用户
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		l.Logger.Errorf("查询用户失败: %v", err)
		return nil, fmt.Errorf("查询用户失败")
	}

	// 防止删除平台管理员
	if user.Role == "platform_admin" {
		return nil, fmt.Errorf("不能删除平台管理员")
	}

	// 开启事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除用户品牌关联
	if err := tx.Where("user_id = ?", userId).Delete(&model.UserBrand{}).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("删除用户品牌关联失败: %v", err)
		return nil, fmt.Errorf("删除用户失败")
	}

	// 删除用户
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("删除用户失败: %v", err)
		return nil, fmt.Errorf("删除用户失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交事务失败: %v", err)
		return nil, fmt.Errorf("删除用户失败")
	}

	l.Logger.Infof("管理员删除用户成功: %s", user.Username)

	return &types.CommonResp{
		Message: "用户删除成功",
	}, nil
}
