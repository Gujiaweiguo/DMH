// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"
	"errors"
	"strconv"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}

	if err := validatePasswordPolicy(req); err != nil {
		return nil, err
	}

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	var policy model.PasswordPolicy
	err = l.svcCtx.DB.Order("id ASC").First(&policy).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Errorf("查询密码策略失败: %v", err)
			return nil, errors.New("查询密码策略失败")
		}

		policy = model.PasswordPolicy{}
	}

	policy.MinLength = req.MinLength
	policy.RequireUppercase = req.RequireUppercase
	policy.RequireLowercase = req.RequireLowercase
	policy.RequireNumbers = req.RequireNumbers
	policy.RequireSpecialChars = req.RequireSpecialChars
	policy.MaxAge = req.MaxAge
	policy.HistoryCount = req.HistoryCount
	policy.MaxLoginAttempts = req.MaxLoginAttempts
	policy.LockoutDuration = req.LockoutDuration
	policy.SessionTimeout = req.SessionTimeout
	policy.MaxConcurrentSessions = req.MaxConcurrentSessions

	if policy.ID == 0 {
		if err := l.svcCtx.DB.Create(&policy).Error; err != nil {
			l.Errorf("创建密码策略失败: %v", err)
			return nil, errors.New("更新密码策略失败")
		}
	} else {
		if err := l.svcCtx.DB.Save(&policy).Error; err != nil {
			l.Errorf("更新密码策略失败: %v", err)
			return nil, errors.New("更新密码策略失败")
		}
	}

	operatorID, hasOperator := userIDFromContext(l.ctx)
	username, _ := l.ctx.Value("username").(string)
	var userID *int64
	if hasOperator {
		userID = &operatorID
	}

	if err := l.svcCtx.DB.Create(&model.AuditLog{
		UserID:     userID,
		Username:   username,
		Action:     "update_password_policy",
		Resource:   "security_policy",
		ResourceID: strconv.FormatInt(policy.ID, 10),
		Details:    "管理员更新密码策略",
		Status:     "success",
	}).Error; err != nil {
		l.Errorf("记录密码策略审计日志失败: %v", err)
	}

	return &types.PasswordPolicyResp{
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
	}, nil
}

func validatePasswordPolicy(req *types.UpdatePasswordPolicyReq) error {
	if req.MinLength < 6 || req.MinLength > 50 {
		return errors.New("最小密码长度必须在6到50之间")
	}

	if req.MaxAge < 0 || req.MaxAge > 365 {
		return errors.New("密码有效期必须在0到365天之间")
	}

	if req.HistoryCount < 0 || req.HistoryCount > 20 {
		return errors.New("历史密码记录数量必须在0到20之间")
	}

	if req.MaxLoginAttempts < 1 || req.MaxLoginAttempts > 20 {
		return errors.New("最大登录尝试次数必须在1到20之间")
	}

	if req.LockoutDuration < 1 || req.LockoutDuration > 1440 {
		return errors.New("锁定时长必须在1到1440分钟之间")
	}

	if req.SessionTimeout < 5 || req.SessionTimeout > 1440 {
		return errors.New("会话超时时间必须在5到1440分钟之间")
	}

	if req.MaxConcurrentSessions < 1 || req.MaxConcurrentSessions > 10 {
		return errors.New("最大并发会话数必须在1到10之间")
	}

	return nil
}

func userIDFromContext(ctx context.Context) (int64, bool) {
	switch value := ctx.Value("userId").(type) {
	case int64:
		return value, true
	case int:
		return int64(value), true
	case float64:
		return int64(value), true
	case string:
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}
