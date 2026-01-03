// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"dmh/api/internal/logic/admin"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func GetUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		l := admin.NewGetUserLogic(r.Context(), svcCtx)
		resp, err := l.GetUser(userId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
