// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"dmh/api/internal/logic/admin"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func UpdateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminUpdateUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 提取路径参数
		userIdStr := pathvar.Vars(r)["id"]
		if userIdStr == "" {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("用户ID不能为空"))
			return
		}
		
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的用户ID"))
			return
		}

		l := admin.NewUpdateUserLogic(r.Context(), svcCtx)
		resp, err := l.UpdateUser(&req, userId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
