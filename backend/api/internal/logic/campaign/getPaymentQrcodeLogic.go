// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"encoding/json"
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

type PaymentConfig struct {
	DepositAmount  float64 `json:"depositAmount"`  // 订金金额
	FullAmount     float64 `json:"fullAmount"`     // 全款金额
	PaymentType    string  `json:"paymentType"`    // 支付类型：deposit/full
	WechatMerchant string  `json:"wechatMerchant"` // 微信商户号
	CallbackURL    string  `json:"callbackUrl"`    // 支付回调地址
}

func (l *GetPaymentQrcodeLogic) GetPaymentQrcode(req types.GetPaymentQrcodeReq) (resp *types.PaymentQrcodeResp, err error) {
	// 1. 查询活动信息
	var campaign model.Campaign
	if err := l.svcCtx.DB.First(&campaign, req.Id).Error; err != nil {
		l.Errorf("查询活动失败: %v", err)
		return nil, fmt.Errorf("活动不存在")
	}

	// 2. 解析支付配置
	var paymentConfig PaymentConfig
	if campaign.PaymentConfig != nil && *campaign.PaymentConfig != "" {
		if err := json.Unmarshal([]byte(*campaign.PaymentConfig), &paymentConfig); err != nil {
			l.Errorf("解析支付配置失败: %v", err)
			// 继续使用默认配置
		}
	}

	// 3. 确定支付金额
	var amount float64
	if paymentConfig.PaymentType == "deposit" && paymentConfig.DepositAmount > 0 {
		amount = paymentConfig.DepositAmount
	} else if paymentConfig.FullAmount > 0 {
		amount = paymentConfig.FullAmount
	} else {
		// 默认金额（使用活动奖励规则作为参考）
		amount = campaign.RewardRule * 10
	}

	// 4. 生成支付二维码
	// 注意：这里使用模拟的支付链接，实际应该调用微信支付 API
	qrcodeContent := fmt.Sprintf("weixin://wxpay/bizpayurl?pr=campaign_%d_%d", req.Id, time.Now().Unix())

	// 5. 返回响应
	resp = &types.PaymentQrcodeResp{
		QrcodeUrl:    qrcodeContent,
		Amount:       amount,
		CampaignName: campaign.Name,
	}

	l.Infof("成功生成支付二维码: campaignId=%d, amount=%.2f", req.Id, amount)
	return resp, nil
}
