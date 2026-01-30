// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"net/http"

	"dmh/api/internal/logic/brand"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetBrandHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := brand.NewGetBrandLogic(r.Context(), svcCtx)
		resp, err := l.GetBrand()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
