package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPasswordPolicyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPasswordPolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPasswordPolicyLogic {
	return &GetPasswordPolicyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPasswordPolicyLogic) GetPasswordPolicy() (resp *types.PasswordPolicyResp, err error) {
	policy, err := l.svcCtx.PasswordService.GetPasswordPolicy()
	if err != nil {
		logx.Errorf("获取密码策略失败: %v", err)
		return nil, err
	}

	resp = &types.PasswordPolicyResp{
		Id:                    policy.ID,
		MinLength:             policy.MinLength,
		RequireUppercase:      policy.RequireUppercase,
		RequireLowercase:      policy.RequireLowercase,
		RequireNumbers:        policy.RequireNumbers,
		RequireSpecialChars:   policy.RequireSpecialChars,
		MaxAge:                policy.MaxAge,
		HistoryCount:          policy.HistoryCount,
		MaxLoginAttempts:      policy.MaxLoginAttempts,
		LockoutDuration:       policy.LockoutDuration,
		SessionTimeout:        policy.SessionTimeout,
		MaxConcurrentSessions: policy.MaxConcurrentSessions,
		CreatedAt:             policy.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:             policy.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}