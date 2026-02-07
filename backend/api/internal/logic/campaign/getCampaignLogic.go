// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCampaignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCampaignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCampaignLogic {
	return &GetCampaignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCampaignLogic) GetCampaign(req *types.GetCampaignReq) (resp *types.CampaignResp, err error) {
	var campaign model.Campaign

	if err := l.svcCtx.DB.First(&campaign, req.Id).Error; err != nil {
		l.Errorf("Failed to query campaign: %v", err)
		return nil, fmt.Errorf("Failed to query campaign: %w", err)
	}

	l.Infof("Successfully queried campaign: id=%d, name=%s, form_fields=%s", campaign.Id, campaign.Name, campaign.FormFields)

	resp = &types.CampaignResp{
		Id:                  campaign.Id,
		BrandId:             campaign.BrandId,
		Name:                campaign.Name,
		Description:         campaign.Description,
		FormFields:          campaign.FormFields,
		RewardRule:          campaign.RewardRule,
		StartTime:           campaign.StartTime.Format("2006-01-02T15:04:05"),
		EndTime:             campaign.EndTime.Format("2006-01-02T15:04:05"),
		Status:              campaign.Status,
		EnableDistribution:  campaign.EnableDistribution,
		DistributionLevel:   campaign.DistributionLevel,
		DistributionRewards: "",
		PaymentConfig:       "",
		PosterTemplateId:    campaign.PosterTemplateId,
		CreatedAt:           campaign.CreatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAt:           campaign.UpdatedAt.Format("2006-01-02T15:04:05"),
	}

	if campaign.DistributionRewards != nil {
		resp.DistributionRewards = *campaign.DistributionRewards
	}
	if campaign.PaymentConfig != nil {
		resp.PaymentConfig = *campaign.PaymentConfig
	}

	return resp, nil
}
