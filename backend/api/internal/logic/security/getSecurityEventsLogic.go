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

type GetSecurityEventsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSecurityEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSecurityEventsLogic {
	return &GetSecurityEventsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSecurityEventsLogic) GetSecurityEvents() (resp *types.SecurityEventListResp, err error) {
	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return &types.SecurityEventListResp{Total: 0, Events: []types.SecurityEventResp{}}, nil
	}

	var total int64
	if err := l.svcCtx.DB.Model(&model.SecurityEvent{}).Count(&total).Error; err != nil {
		l.Errorf("查询安全事件总数失败: %v", err)
		return nil, err
	}

	if total == 0 {
		return &types.SecurityEventListResp{Total: 0, Events: []types.SecurityEventResp{}}, nil
	}

	var securityEvents []model.SecurityEvent
	if err := l.svcCtx.DB.Order("created_at DESC, id DESC").Find(&securityEvents).Error; err != nil {
		l.Errorf("查询安全事件失败: %v", err)
		return nil, err
	}

	events := make([]types.SecurityEventResp, 0, len(securityEvents))
	for _, event := range securityEvents {
		var userID int64
		if event.UserID != nil {
			userID = *event.UserID
		}

		var handledBy int64
		if event.HandledBy != nil {
			handledBy = *event.HandledBy
		}

		handledAt := ""
		if event.HandledAt != nil {
			handledAt = event.HandledAt.Format("2006-01-02 15:04:05")
		}

		events = append(events, types.SecurityEventResp{
			Id:          event.ID,
			EventType:   event.EventType,
			Severity:    event.Severity,
			UserId:      userID,
			Username:    event.Username,
			ClientIp:    event.ClientIP,
			UserAgent:   event.UserAgent,
			Description: event.Description,
			Details:     event.Details,
			Handled:     event.Handled,
			HandledBy:   handledBy,
			HandledAt:   handledAt,
			CreatedAt:   event.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.SecurityEventListResp{Total: total, Events: events}, nil
}
