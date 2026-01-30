// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"net/http"

	"dmh/api/internal/logic/menu"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ConfigRoleMenusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RoleMenuReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := menu.NewConfigRoleMenusLogic(r.Context(), svcCtx)
		resp, err := l.ConfigRoleMenus(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
