package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

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
	// 从查询参数获取分页和过滤条件
	page := 1
	pageSize := 20
	filters := make(map[string]interface{})

	// 这里应该从HTTP请求中获取查询参数，暂时使用默认值
	// 在实际实现中，需要从gin.Context或其他HTTP上下文中获取参数

	events, total, err := l.svcCtx.AuditService.GetSecurityEvents(page, pageSize, filters)
	if err != nil {
		logx.Errorf("获取安全事件失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var eventList []types.SecurityEventResp
	for _, event := range events {
		eventResp := types.SecurityEventResp{
			Id:          event.ID,
			EventType:   event.EventType,
			Severity:    event.Severity,
			Username:    event.Username,
			ClientIp:    event.ClientIP,
			UserAgent:   event.UserAgent,
			Description: event.Description,
			Details:     event.Details,
			Handled:     event.Handled,
			CreatedAt:   event.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		
		if event.UserID != nil {
			eventResp.UserId = *event.UserID
		}
		
		if event.HandledBy != nil {
			eventResp.HandledBy = *event.HandledBy
		}
		
		if event.HandledAt != nil {
			eventResp.HandledAt = event.HandledAt.Format("2006-01-02 15:04:05")
		}
		
		eventList = append(eventList, eventResp)
	}

	resp = &types.SecurityEventListResp{
		Total:  total,
		Events: eventList,
	}

	return resp, nil
}