// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"errors"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProfileLogic) UpdateProfile(req *types.UpdateProfileReq) (resp *types.UserInfoResp, err error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, errors.New("无法从context中获取用户ID")
	}

	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		l.Errorf("查询用户失败: %v", err)
		return nil, errors.New("用户不存在")
	}

	if req.RealName != "" {
		user.RealName = req.RealName
	}

	if err := l.svcCtx.DB.Save(&user).Error; err != nil {
		l.Errorf("更新用户信息失败: %v", err)
		return nil, errors.New("更新失败")
	}

	return &types.UserInfoResp{
		Id:       user.Id,
		Username: user.Username,
		RealName: user.RealName,
	}, nil
}
