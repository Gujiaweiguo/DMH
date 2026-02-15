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

func GetBrandAssetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		brandIdStr := r.PathValue("brandId")
		idStr := r.PathValue("id")

		brandId, _ := strconv.ParseInt(brandIdStr, 10, 64)
		id, _ := strconv.ParseInt(idStr, 10, 64)

		req := &types.GetBrandAssetReq{
			BrandId: brandId,
			Id:      id,
		}

		l := brand.NewGetBrandAssetLogic(r.Context(), svcCtx)
		resp, err := l.GetBrandAsset(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
