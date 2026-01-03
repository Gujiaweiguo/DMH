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
	// 生成密码强度评分
	score := l.svcCtx.PasswordService.GeneratePasswordStrengthScore(req.NewPassword)
	level := l.svcCtx.PasswordService.GetPasswordStrengthLevel(score)

	// 生成建议信息
	var message string
	switch level {
	case "很弱":
		message = "密码强度很弱，建议使用包含大小写字母、数字和特殊字符的复杂密码"
	case "弱":
		message = "密码强度较弱，建议增加密码长度或使用更多字符类型"
	case "中等":
		message = "密码强度中等，可以考虑进一步增强"
	case "强":
		message = "密码强度良好"
	default:
		message = "密码强度评估完成"
	}

	resp = &types.PasswordStrengthResp{
		Score:   score,
		Level:   level,
		Message: message,
	}

	return resp, nil
}