// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"net/http"

	"dmh/api/internal/logic/order"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ScanOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewScanOrderLogic(r.Context(), svcCtx)
		resp, err := l.ScanOrder()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
