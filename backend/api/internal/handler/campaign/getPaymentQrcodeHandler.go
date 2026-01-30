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

func GetPaymentQrcodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPaymentQrcodeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := campaign.NewGetPaymentQrcodeLogic(r.Context(), svcCtx)
		resp, err := l.GetPaymentQrcode(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
