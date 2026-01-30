// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package reward

import (
	"net/http"

	"dmh/api/internal/logic/reward"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRewardsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := reward.NewGetRewardsLogic(r.Context(), svcCtx)
		resp, err := l.GetRewards()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
