// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePasswordPolicyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePasswordPolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordPolicyLogic {
	return &UpdatePasswordPolicyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePasswordPolicyLogic) UpdatePasswordPolicy(req *types.UpdatePasswordPolicyReq) (resp *types.PasswordPolicyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
