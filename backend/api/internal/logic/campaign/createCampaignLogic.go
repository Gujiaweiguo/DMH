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

type CreateCampaignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCampaignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCampaignLogic {
	return &CreateCampaignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCampaignLogic) CreateCampaign(req *types.CreateCampaignReq) (resp *types.CampaignResp, err error) {
	startTime, err1 := time.Parse("2006-01-02T15:04:05", req.StartTime)
	endTime, err2 := time.Parse("2006-01-02T15:04:05", req.EndTime)
	if err1 != nil || err2 != nil {
		l.Errorf("Time format error: startTime=%v, endTime=%v", err1, err2)
		return nil, fmt.Errorf("Time format error")
	}

	newCampaign := model.Campaign{
		BrandId:          req.BrandId,
		Name:             req.Name,
		Description:      req.Description,
		RewardRule:       req.RewardRule,
		StartTime:        startTime,
		EndTime:          endTime,
		Status:           "active",
		PosterTemplateId: 1,
	}

	if len(req.FormFields) > 0 {
		formFieldsJSON, err := json.Marshal(req.FormFields)
		if err == nil {
			newCampaign.FormFields = string(formFieldsJSON)
		}
	}

	if err := l.svcCtx.DB.Create(&newCampaign).Error; err != nil {
		l.Errorf("Failed to create campaign: %v", err)
		return nil, fmt.Errorf("Failed to create campaign: %w", err)
	}

	l.Infof("Campaign created successfully: campaignId=%d, name=%s", newCampaign.Id, newCampaign.Name)

	var formFieldsStr string
	if newCampaign.FormFields != "" {
		formFieldsStr = newCampaign.FormFields
	} else {
		formFieldsStr = "[]"
	}

	resp = &types.CampaignResp{
		Id:          newCampaign.Id,
		BrandId:     newCampaign.BrandId,
		Name:        newCampaign.Name,
		Description: newCampaign.Description,
		FormFields:  formFieldsStr,
		RewardRule:  newCampaign.RewardRule,
		StartTime:   newCampaign.StartTime.Format("2006-01-02T15:04:05"),
		EndTime:     newCampaign.EndTime.Format("2006-01-02T15:04:05"),
		Status:      newCampaign.Status,
		CreatedAt:   newCampaign.CreatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAt:   newCampaign.UpdatedAt.Format("2006-01-02T15:04:05"),
	}

	return resp, nil
}
