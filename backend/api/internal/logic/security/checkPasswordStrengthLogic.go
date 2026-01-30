// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckPasswordStrengthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckPasswordStrengthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPasswordStrengthLogic {
	return &CheckPasswordStrengthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckPasswordStrengthLogic) CheckPasswordStrength(req *types.ChangePasswordReq) (resp *types.PasswordStrengthResp, err error) {
	// todo: add your logic here and delete this line

	return
}
