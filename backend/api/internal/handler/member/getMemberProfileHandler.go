// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package member

import (
	"net/http"

	"dmh/api/internal/logic/member"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMemberProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memberId, err := parseMemberIDFromPath(r.URL.Path)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		req := &types.GetMemberReq{Id: memberId}

		l := member.NewGetMemberProfileLogic(r.Context(), svcCtx)
		resp, err := l.GetMemberProfile(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
