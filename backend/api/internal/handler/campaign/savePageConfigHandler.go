// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"net/http"

	"dmh/api/internal/logic/campaign"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SavePageConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageConfigReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := campaign.NewSavePageConfigLogic(r.Context(), svcCtx)
		resp, err := l.SavePageConfig(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
