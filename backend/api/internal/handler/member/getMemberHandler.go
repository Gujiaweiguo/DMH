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

func GetMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		memberIdStr := pathParts[len(pathParts)-1]
		memberId, _ := strconv.ParseInt(memberIdStr, 10, 64)

		req := &types.GetMemberReq{Id: memberId}

		l := member.NewGetMemberLogic(r.Context(), svcCtx)
		resp, err := l.GetMember(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
