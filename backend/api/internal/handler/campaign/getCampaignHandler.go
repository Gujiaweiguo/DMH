// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCampaignHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 URL 路径中提取 ID
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/campaigns/")
		idStr := strings.TrimSuffix(path, "/")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("invalid campaign id: %w", err))
			return
		}

		var campaign model.Campaign

		if err := svcCtx.DB.First(&campaign, id).Error; err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("campaign not found: %w", err))
			return
		}

		// 解析 formFields JSON 字符串为对象
		var formFields []types.FormField
		if campaign.FormFields != "" {
			if err := json.Unmarshal([]byte(campaign.FormFields), &formFields); err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
		}

		resp := &types.CampaignResp{
			Id:                 campaign.Id,
			BrandId:            campaign.BrandId,
			Name:               campaign.Name,
			Description:        campaign.Description,
			FormFields:         formFields,
			RewardRule:         campaign.RewardRule,
			StartTime:          campaign.StartTime.Format("2006-01-02T15:04:05"),
			EndTime:            campaign.EndTime.Format("2006-01-02T15:04:05"),
			Status:             campaign.Status,
			EnableDistribution: campaign.EnableDistribution,
			DistributionLevel:  campaign.DistributionLevel,
			DistributionRewards: func() string {
				if campaign.DistributionRewards != nil {
					return *campaign.DistributionRewards
				}
				return ""
			}(),
			PaymentConfig: func() string {
				if campaign.PaymentConfig != nil {
					return *campaign.PaymentConfig
				}
				return ""
			}(),
			PosterTemplateId: campaign.PosterTemplateId,
			CreatedAt:        campaign.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:        campaign.UpdatedAt.Format("2006-01-02T15:04:05"),
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
