// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

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
	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return &types.UserSessionListResp{Total: 0, Sessions: []types.UserSessionResp{}}, nil
	}

	var total int64
	if err := l.svcCtx.DB.Model(&model.UserSession{}).Count(&total).Error; err != nil {
		l.Errorf("查询用户会话总数失败: %v", err)
		return nil, err
	}

	if total == 0 {
		return &types.UserSessionListResp{Total: 0, Sessions: []types.UserSessionResp{}}, nil
	}

	var sessions []model.UserSession
	if err := l.svcCtx.DB.Order("last_active_at DESC, created_at DESC").Find(&sessions).Error; err != nil {
		l.Errorf("查询用户会话失败: %v", err)
		return nil, err
	}

	respSessions := make([]types.UserSessionResp, 0, len(sessions))
	for _, session := range sessions {
		respSessions = append(respSessions, types.UserSessionResp{
			Id:           session.ID,
			UserId:       session.UserID,
			ClientIp:     session.ClientIP,
			UserAgent:    session.UserAgent,
			LoginAt:      session.LoginAt.Format("2006-01-02 15:04:05"),
			LastActiveAt: session.LastActiveAt.Format("2006-01-02 15:04:05"),
			ExpiresAt:    session.ExpiresAt.Format("2006-01-02 15:04:05"),
			Status:       session.Status,
			CreatedAt:    session.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.UserSessionListResp{Total: total, Sessions: respSessions}, nil
}
