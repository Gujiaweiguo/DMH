// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"net/http"

	"dmh/api/internal/logic/distributor"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMyDistributorDashboardHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := distributor.NewGetMyDistributorDashboardLogic(r.Context(), svcCtx)
		resp, err := l.GetMyDistributorDashboard()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
