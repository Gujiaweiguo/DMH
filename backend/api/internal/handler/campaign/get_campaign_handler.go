// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"net/http"

	"dmh/api/internal/logic/campaign"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCampaignHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := campaign.NewGetCampaignLogic(r.Context(), svcCtx)
		
		// 从URL路径中获取活动ID
		id := r.URL.Path[len("/api/v1/campaigns/"):]
		
		resp, err := l.GetCampaign(id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
