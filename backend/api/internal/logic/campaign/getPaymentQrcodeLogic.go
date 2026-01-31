// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"fmt"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentQrcodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaymentQrcodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentQrcodeLogic {
	return &GetPaymentQrcodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaymentQrcodeLogic) GetPaymentQrcode(req types.GetPaymentQrcodeReq) (resp *types.PaymentQrcodeResp, err error) {
	var campaign model.Campaign
	if err := l.svcCtx.DB.First(&campaign, req.Id).Error; err != nil {
		l.Errorf("查询活动失败: %v", err)
		return nil, fmt.Errorf("活动不存在")
	}

	amount := campaign.RewardRule * 10

	qrcodeContent := fmt.Sprintf("weixin://wxpay/bizpayurl?pr=campaign_%d_%d", req.Id, time.Now().Unix())

	resp = &types.PaymentQrcodeResp{
		QrcodeUrl:    qrcodeContent,
		Amount:       amount,
		CampaignName: campaign.Name,
	}

	l.Infof("成功生成支付二维码: campaignId=%d, amount=%.2f", req.Id, amount)
	return resp, nil
}
