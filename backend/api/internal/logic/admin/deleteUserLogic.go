// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"errors"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *DeleteUserLogic) DeleteUser() (resp *types.CommonResp, err error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, errors.New("无法从context中获取用户ID")
	}

	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if err := l.svcCtx.DB.Delete(&user).Error; err != nil {
		l.Errorf("删除用户失败: %v", err)
		return nil, errors.New("删除用户失败")
	}

	return nil, nil
}
