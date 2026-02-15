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

type BindPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindPhoneLogic {
	return &BindPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindPhoneLogic) BindPhone(req *types.BindPhoneReq) (resp *types.UserInfoResp, err error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, errors.New("无法从context中获取用户ID")
	}

	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		l.Errorf("查询用户失败: %v", err)
		return nil, errors.New("用户不存在")
	}

	if req.Phone == "" {
		return nil, errors.New("手机号不能为空")
	}

	user.Phone = req.Phone
	if err := l.svcCtx.DB.Save(&user).Error; err != nil {
		l.Errorf("绑定手机号失败: %v", err)
		return nil, errors.New("绑定失败")
	}

	return &types.UserInfoResp{
		Id:       user.Id,
		Username: user.Username,
		Phone:    user.Phone,
	}, nil
}
