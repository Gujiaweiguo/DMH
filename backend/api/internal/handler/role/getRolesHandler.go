// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"net/http"

	"dmh/api/internal/logic/role"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewGetRolesLogic(r.Context(), svcCtx)
		resp, err := l.GetRoles()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
