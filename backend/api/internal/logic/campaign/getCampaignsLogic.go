// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"encoding/json"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCampaignsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCampaignsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCampaignsLogic {
	return &GetCampaignsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCampaignsLogic) GetCampaigns(req *types.GetCampaignsReq) (resp *types.CampaignListResp, err error) {
	var campaigns []model.Campaign

	query := l.svcCtx.DB.Model(&model.Campaign{})

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(int(offset)).Limit(int(req.PageSize))
	}

	if err := query.Find(&campaigns).Error; err != nil {
		l.Errorf("Failed to query campaigns: %v", err)
		return nil, fmt.Errorf("Failed to query campaigns: %w", err)
	}

	l.Infof("Successfully queried campaigns: count=%d", len(campaigns))

	var campaignResps []types.CampaignResp
	for _, campaign := range campaigns {
		distributionRewards := ""
		if campaign.DistributionRewards != nil {
			distributionRewards = *campaign.DistributionRewards
		}
		paymentConfig := ""
		if campaign.PaymentConfig != nil {
			paymentConfig = *campaign.PaymentConfig
		}

		// 解析 formFields JSON 字符串为对象
		var formFields []types.FormField
		if campaign.FormFields != "" {
			if err := json.Unmarshal([]byte(campaign.FormFields), &formFields); err != nil {
				l.Errorf("Failed to parse formFields JSON: %v", err)
			}
		}

		campaignResps = append(campaignResps, types.CampaignResp{
			Id:                  campaign.Id,
			BrandId:             campaign.BrandId,
			Name:                campaign.Name,
			Description:         campaign.Description,
			FormFields:          formFields,
			RewardRule:          campaign.RewardRule,
			StartTime:           campaign.StartTime.Format("2006-01-02T15:04:05"),
			EndTime:             campaign.EndTime.Format("2006-01-02T15:04:05"),
			Status:              campaign.Status,
			EnableDistribution:  campaign.EnableDistribution,
			DistributionLevel:   campaign.DistributionLevel,
			DistributionRewards: distributionRewards,
			PaymentConfig:       paymentConfig,
			PosterTemplateId:    campaign.PosterTemplateId,
			CreatedAt:           campaign.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:           campaign.UpdatedAt.Format("2006-01-02T15:04:05"),
		})
	}

	resp = &types.CampaignListResp{
		Total:     int64(len(campaigns)),
		Campaigns: campaignResps,
	}

	return resp, nil
}
