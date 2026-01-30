// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package poster

import (
	"net/http"

	"dmh/api/internal/logic/poster"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPosterRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := poster.NewGetPosterRecordsLogic(r.Context(), svcCtx)
		resp, err := l.GetPosterRecords()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
