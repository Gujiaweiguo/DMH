package health

import (
	"context"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
)

type GetHealthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHealthLogic {
	return &GetHealthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHealthLogic) GetHealth() (*types.HealthResp, error) {
	resp := &types.HealthResp{
		Status:    "ok",
		Timestamp: time.Now().Unix(),
	}

	if l.svcCtx.DB == nil {
		resp.Database = "unavailable"
		return resp, nil
	}

	sqlDB, err := l.svcCtx.DB.DB()
	if err != nil {
		resp.Database = "unavailable"
		return resp, nil
	}

	if err := sqlDB.PingContext(l.ctx); err != nil {
		resp.Database = "unhealthy"
		return resp, nil
	}

	resp.Database = "ok"
	return resp, nil
}
