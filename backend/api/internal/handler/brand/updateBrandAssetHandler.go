// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"net/http"
	"strconv"

	"dmh/api/internal/logic/brand"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateBrandAssetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateBrandAssetReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		brandIdStr := r.PathValue("brandId")
		idStr := r.PathValue("id")

		brandId, _ := strconv.ParseInt(brandIdStr, 10, 64)
		id, _ := strconv.ParseInt(idStr, 10, 64)

		req.BrandId = brandId
		req.Id = id

		l := brand.NewUpdateBrandAssetLogic(r.Context(), svcCtx)
		resp, err := l.UpdateBrandAsset(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
