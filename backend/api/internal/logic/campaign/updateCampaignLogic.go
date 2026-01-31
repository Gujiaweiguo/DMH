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
	startTime, err1 := time.Parse("2006-01-02T15:04:05", req.StartTime)
	endTime, err2 := time.Parse("2006-01-02T15:04:05", req.EndTime)
	if err1 != nil || err2 != nil {
		l.Errorf("Time format error: startTime=%v, endTime=%v", err1, err2)
		return nil, fmt.Errorf("Time format error")
	}

	var campaign model.Campaign
	if err := l.svcCtx.DB.Where("id = ?", req.BrandId).First(&campaign).Error; err != nil {
		l.Errorf("Failed to query campaign: %v", err)
		return nil, fmt.Errorf("Campaign not found")
	}

	campaign.Name = req.Name
	campaign.Description = req.Description
	campaign.StartTime = startTime
	campaign.EndTime = endTime

	if len(req.FormFields) > 0 {
		formFieldsJSON, _ := json.Marshal(req.FormFields)
		json.Unmarshal(formFieldsJSON, &campaign.FormFields)
	}

	if err := l.svcCtx.DB.Save(&campaign).Error; err != nil {
		l.Errorf("Failed to update campaign: %v", err)
		return nil, fmt.Errorf("Failed to update campaign: %w", err)
	}

	l.Infof("Campaign updated successfully: campaignId=%d, name=%s", campaign.Id, campaign.Name)

	var formFieldsStr string
	if campaign.FormFields != "" {
		formFieldsStr = campaign.FormFields
	} else {
		formFieldsStr = "[]"
	}

	resp = &types.CampaignResp{
		Id:          campaign.Id,
		BrandId:     campaign.BrandId,
		Name:        campaign.Name,
		Description: campaign.Description,
		FormFields:  formFieldsStr,
		RewardRule:  campaign.RewardRule,
		StartTime:   campaign.StartTime.Format("2006-01-02T15:04:05"),
		EndTime:     campaign.EndTime.Format("2006-01-02T15:04:05"),
		Status:      campaign.Status,
		CreatedAt:   campaign.CreatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAt:   campaign.UpdatedAt.Format("2006-01-02T15:04:05"),
	}

	return resp, nil
}
