// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package poster

import (
	"context"
	"fmt"

	"dmh/common/poster"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateCampaignPosterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateCampaignPosterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateCampaignPosterLogic {
	return &GenerateCampaignPosterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateCampaignPosterLogic) GenerateCampaignPoster(req *types.GeneratePosterReq) (resp *types.GeneratePosterResp, err error) {
	startTime := time.Now()

	// 1. 查询活动信息
	var campaign model.Campaign
	if err := l.svcCtx.DB.First(&campaign, req.TemplateId).Error; err != nil {
		l.Errorf("查询活动失败: %v", err)
		return nil, fmt.Errorf("活动不存在")
	}

	// 2. 调用海报服务生成海报
	campaignDesc := ""
	if campaign.Description != nil {
		campaignDesc = *campaign.Description
	}

	posterURL, err := l.svcCtx.PosterService.GenerateCampaignPoster(
		campaign.Name,
		campaignDesc,
		"",
	)
	if err != nil {
		l.Errorf("生成海报失败: %v", err)
		return nil, fmt.Errorf("生成海报失败: %w", err)
	}

	generationTime := time.Since(startTime).Milliseconds()

	// 3. 保存海报记录到数据库
	posterRecord := model.PosterRecord{
		RecordType:     "campaign",
		CampaignID:     &campaign.Id,
		DistributorID:  nil,
		TemplateName:   fmt.Sprintf("模板%d", req.TemplateId),
		PosterUrl:      posterURL,
		ThumbnailUrl:   "",
		FileSize:       "",
		GenerationTime: int(generationTime),
		DownloadCount:  0,
		ShareCount:     0,
		Status:         "success",
	}

	if err := l.svcCtx.DB.Create(&posterRecord).Error; err != nil {
		l.Errorf("保存海报记录失败: %v", err)
		return nil, fmt.Errorf("保存海报记录失败: %w", err)
	}

	l.Infof("海报生成成功: campaignId=%d, posterUrl=%s, time=%dms", campaign.Id, posterURL, generationTime)

	// 4. 返回响应
	resp = &types.GeneratePosterResp{
		PosterUrl:      posterURL,
		ThumbnailUrl:   "",
		FileSize:       "",
		GenerationTime: int(generationTime),
	}

	return resp, nil
}
