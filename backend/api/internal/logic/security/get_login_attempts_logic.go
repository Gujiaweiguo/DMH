package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoginAttemptsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLoginAttemptsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginAttemptsLogic {
	return &GetLoginAttemptsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoginAttemptsLogic) GetLoginAttempts() (resp *types.LoginAttemptListResp, err error) {
	// 从查询参数获取分页和过滤条件
	page := 1
	pageSize := 20
	filters := make(map[string]interface{})

	// 这里应该从HTTP请求中获取查询参数，暂时使用默认值
	// 在实际实现中，需要从gin.Context或其他HTTP上下文中获取参数

	attempts, total, err := l.svcCtx.AuditService.GetLoginAttempts(page, pageSize, filters)
	if err != nil {
		logx.Errorf("获取登录尝试记录失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var attemptList []types.LoginAttemptResp
	for _, attempt := range attempts {
		attemptResp := types.LoginAttemptResp{
			Id:         attempt.ID,
			Username:   attempt.Username,
			ClientIp:   attempt.ClientIP,
			UserAgent:  attempt.UserAgent,
			Success:    attempt.Success,
			FailReason: attempt.FailReason,
			CreatedAt:  attempt.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		
		if attempt.UserID != nil {
			attemptResp.UserId = *attempt.UserID
		}
		
		attemptList = append(attemptList, attemptResp)
	}

	resp = &types.LoginAttemptListResp{
		Total:    total,
		Attempts: attemptList,
	}

	return resp, nil
}