// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"net/http"

	"dmh/api/internal/logic/role"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ConfigRolePermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RolePermissionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewConfigRolePermissionsLogic(r.Context(), svcCtx)
		resp, err := l.ConfigRolePermissions(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
