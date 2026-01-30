// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package poster

import (
	"net/http"

	"dmh/api/internal/logic/poster"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GenerateCampaignPosterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GeneratePosterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := poster.NewGenerateCampaignPosterLogic(r.Context(), svcCtx)
		resp, err := l.GenerateCampaignPoster(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
