// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package withdrawal

import (
	"net/http"

	"dmh/api/internal/logic/withdrawal"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetWithdrawalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := withdrawal.NewGetWithdrawalLogic(r.Context(), svcCtx)
		resp, err := l.GetWithdrawal()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
