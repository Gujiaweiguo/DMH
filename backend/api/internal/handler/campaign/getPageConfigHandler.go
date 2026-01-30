// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"net/http"

	"dmh/api/internal/logic/campaign"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPageConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := campaign.NewGetPageConfigLogic(r.Context(), svcCtx)
		resp, err := l.GetPageConfig()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
