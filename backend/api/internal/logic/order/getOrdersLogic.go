package order

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrdersLogic {
	return &GetOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrdersLogic) GetOrders() (resp *types.OrderListResp, err error) {
	var orders []types.OrderResp

	if err := l.svcCtx.DB.Order("created_at DESC").Find(&orders).Error; err != nil {
		l.Errorf("查询订单列表失败: %v", err)
		return nil, err
	}

	resp = &types.OrderListResp{
		Total:  int64(len(orders)),
		Orders: orders,
	}
	return resp, nil
}
