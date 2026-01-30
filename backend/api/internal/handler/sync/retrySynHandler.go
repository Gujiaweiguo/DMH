// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sync

import (
	"net/http"

	"dmh/api/internal/logic/sync"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RetrySynHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sync.NewRetrySynLogic(r.Context(), svcCtx)
		resp, err := l.RetrySyn()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
