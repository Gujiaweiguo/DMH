// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"net/http"

	"dmh/api/internal/logic/distributor"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GenerateDistributorLinkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenerateLinkReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewGenerateDistributorLinkLogic(r.Context(), svcCtx)
		resp, err := l.GenerateDistributorLink(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
