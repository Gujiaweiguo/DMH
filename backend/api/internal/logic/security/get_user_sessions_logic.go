package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserSessionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSessionsLogic {
	return &GetUserSessionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserSessionsLogic) GetUserSessions() (resp *types.UserSessionListResp, err error) {
	// 从查询参数获取分页条件
	page := 1
	pageSize := 20

	// 这里应该从HTTP请求中获取查询参数，暂时使用默认值
	// 在实际实现中，需要从gin.Context或其他HTTP上下文中获取参数

	sessions, total, err := l.svcCtx.SessionService.GetActiveSessions(page, pageSize)
	if err != nil {
		logx.Errorf("获取用户会话失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var sessionList []types.UserSessionResp
	for _, session := range sessions {
		sessionResp := types.UserSessionResp{
			Id:           session.ID,
			UserId:       session.UserID,
			ClientIp:     session.ClientIP,
			UserAgent:    session.UserAgent,
			LoginAt:      session.LoginAt.Format("2006-01-02 15:04:05"),
			LastActiveAt: session.LastActiveAt.Format("2006-01-02 15:04:05"),
			ExpiresAt:    session.ExpiresAt.Format("2006-01-02 15:04:05"),
			Status:       session.Status,
			CreatedAt:    session.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		
		sessionList = append(sessionList, sessionResp)
	}

	resp = &types.UserSessionListResp{
		Total:    total,
		Sessions: sessionList,
	}

	return resp, nil
}