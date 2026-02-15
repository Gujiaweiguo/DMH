// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCampaignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCampaignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCampaignLogic {
	return &UpdateCampaignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCampaignLogic) UpdateCampaign(req *types.UpdateCampaignReq) (resp *types.CampaignResp, err error) {
	var campaign model.Campaign
	if err := l.svcCtx.DB.Where("id = ?", req.Id).First(&campaign).Error; err != nil {
		l.Errorf("Failed to query campaign: %v", err)
		return nil, fmt.Errorf("Campaign not found")
	}

	if req.BrandId != nil && *req.BrandId > 0 {
		campaign.BrandId = *req.BrandId
	}
	if req.Name != nil {
		campaign.Name = *req.Name
	}
	if req.Description != nil {
		campaign.Description = *req.Description
	}
	if req.StartTime != nil {
		startTime, err := time.Parse("2006-01-02T15:04:05", *req.StartTime)
		if err != nil {
			l.Errorf("Time format error: startTime=%v", err)
			return nil, fmt.Errorf("Time format error")
		}
		campaign.StartTime = startTime
	}
	if req.EndTime != nil {
		endTime, err := time.Parse("2006-01-02T15:04:05", *req.EndTime)
		if err != nil {
			l.Errorf("Time format error: endTime=%v", err)
			return nil, fmt.Errorf("Time format error")
		}
		campaign.EndTime = endTime
	}
	if req.RewardRule != nil {
		campaign.RewardRule = *req.RewardRule
	}
	if req.Status != nil && *req.Status != "" {
		campaign.Status = *req.Status
	}

	if req.DistributionLevel != nil {
		if *req.DistributionLevel < 1 || *req.DistributionLevel > 3 {
			return nil, fmt.Errorf("Distribution level must be between 1 and 3")
		}
		campaign.DistributionLevel = *req.DistributionLevel
	}
	if req.DistributionRewards != nil {
		if strings.TrimSpace(*req.DistributionRewards) == "" {
			campaign.DistributionRewards = nil
		} else {
			campaign.DistributionRewards = req.DistributionRewards
		}
	}
	if req.EnableDistribution != nil {
		campaign.EnableDistribution = *req.EnableDistribution
		if !campaign.EnableDistribution {
			campaign.DistributionRewards = nil
		} else if campaign.DistributionLevel == 0 {
			campaign.DistributionLevel = 1
		}
	}

	if req.PaymentConfig != nil {
		if strings.TrimSpace(*req.PaymentConfig) == "" {
			campaign.PaymentConfig = nil
		} else {
			campaign.PaymentConfig = req.PaymentConfig
		}
	}

	if req.PosterTemplateId != nil && *req.PosterTemplateId > 0 {
		campaign.PosterTemplateId = *req.PosterTemplateId
	}

	if req.FormFields != nil {
		formFieldsJSON, _ := json.Marshal(req.FormFields)
		json.Unmarshal(formFieldsJSON, &campaign.FormFields)
	}

	if err := l.svcCtx.DB.Save(&campaign).Error; err != nil {
		l.Errorf("Failed to update campaign: %v", err)
		return nil, fmt.Errorf("Failed to update campaign: %w", err)
	}

	l.Infof("Campaign updated successfully: campaignId=%d, name=%s", campaign.Id, campaign.Name)

	// 解析 formFields JSON 字符串为对象
	var formFields []types.FormField
	if campaign.FormFields != "" {
		if err := json.Unmarshal([]byte(campaign.FormFields), &formFields); err != nil {
			l.Errorf("Failed to parse formFields JSON: %v", err)
			// 继续使用空数组
		}
	}

	distributionRewardsResp := ""
	if campaign.DistributionRewards != nil {
		distributionRewardsResp = *campaign.DistributionRewards
	}
	paymentConfigResp := ""
	if campaign.PaymentConfig != nil {
		paymentConfigResp = *campaign.PaymentConfig
	}

	resp = &types.CampaignResp{
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
		DistributionRewards: distributionRewardsResp,
		PaymentConfig:       paymentConfigResp,
		PosterTemplateId:    campaign.PosterTemplateId,
		CreatedAt:           campaign.CreatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAt:           campaign.UpdatedAt.Format("2006-01-02T15:04:05"),
	}

	return resp, nil
}
