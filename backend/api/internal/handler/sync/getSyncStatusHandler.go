// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sync

import (
	"net/http"

	"dmh/api/internal/logic/sync"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetSyncStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sync.NewGetSyncStatusLogic(r.Context(), svcCtx)
		resp, err := l.GetSyncStatus()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
