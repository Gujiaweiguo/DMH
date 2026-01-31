// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package member

import (
	"net/http"
	"strconv"
	"strings"

	"dmh/api/internal/logic/member"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateMemberStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateMemberStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		pathParts := strings.Split(r.URL.Path, "/")
		memberIdStr := pathParts[len(pathParts)-1]
		memberId, _ := strconv.ParseInt(memberIdStr, 10, 64)

		l := member.NewUpdateMemberStatusLogic(r.Context(), svcCtx)
		resp, err := l.UpdateMemberStatus(memberId, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
