// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"
	"errors"
	"strings"

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
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}

	password := strings.TrimSpace(req.NewPassword)
	if password == "" {
		return nil, errors.New("新密码不能为空")
	}

	score := calculatePasswordStrength(password)
	level := strengthLevel(score)
	message := "密码强度较好"

	if l.svcCtx != nil && l.svcCtx.PasswordService != nil && l.svcCtx.DB != nil {
		score = l.svcCtx.PasswordService.GeneratePasswordStrengthScore(password)
		level = l.svcCtx.PasswordService.GetPasswordStrengthLevel(score)
		if validateErr := l.svcCtx.PasswordService.ValidatePassword(password, 0); validateErr != nil {
			message = validateErr.Error()
		} else {
			message = "密码符合当前安全策略"
		}
	} else if score < 60 {
		message = "建议增加长度并混合大小写、数字和特殊字符"
	}

	return &types.PasswordStrengthResp{
		Score:   score,
		Level:   level,
		Message: message,
	}, nil
}

func calculatePasswordStrength(password string) int {
	score := 0

	if len(password) >= 8 {
		score += 30
	}
	if len(password) >= 12 {
		score += 15
	}

	hasLower := false
	hasUpper := false
	hasNumber := false
	hasSpecial := false

	for _, ch := range password {
		switch {
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= '0' && ch <= '9':
			hasNumber = true
		default:
			hasSpecial = true
		}
	}

	if hasLower {
		score += 15
	}
	if hasUpper {
		score += 15
	}
	if hasNumber {
		score += 15
	}
	if hasSpecial {
		score += 10
	}

	if score > 100 {
		return 100
	}

	return score
}

func strengthLevel(score int) string {
	if score >= 80 {
		return "强"
	}
	if score >= 60 {
		return "中等"
	}
	if score >= 40 {
		return "弱"
	}
	return "很弱"
}
