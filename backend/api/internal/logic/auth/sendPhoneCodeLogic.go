// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"errors"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendPhoneCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendPhoneCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPhoneCodeLogic {
	return &SendPhoneCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendPhoneCodeLogic) SendPhoneCode(req *types.SendCodeReq) (resp *types.CommonResp, err error) {
	if req.Target == "" {
		return nil, errors.New("手机号不能为空")
	}

	l.Infof("发送验证码到手机号: %s", req.Target)

	return &types.CommonResp{
		Message: "发送成功",
	}, nil
}
