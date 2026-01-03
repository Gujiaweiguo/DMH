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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ResetUserPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	return &ResetUserPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetUserPasswordLogic) ResetUserPassword(req *types.AdminResetPasswordReq, userId int64) (resp *types.CommonResp, err error) {
	// 参数验证
	if req.NewPassword == "" {
		return nil, fmt.Errorf("新密码不能为空")
	}

	if len(req.NewPassword) < 6 {
		return nil, fmt.Errorf("密码长度不能少于6位")
	}

	// 查找用户
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		l.Logger.Errorf("查询用户失败: %v", err)
		return nil, fmt.Errorf("查询用户失败")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		l.Logger.Errorf("密码加密失败: %v", err)
		return nil, fmt.Errorf("密码加密失败")
	}

	// 更新密码
	if err := l.svcCtx.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		l.Logger.Errorf("重置用户密码失败: %v", err)
		return nil, fmt.Errorf("重置密码失败")
	}

	l.Logger.Infof("管理员重置用户密码成功: %s", user.Username)

	return &types.CommonResp{
		Message: "密码重置成功",
	}, nil
}
