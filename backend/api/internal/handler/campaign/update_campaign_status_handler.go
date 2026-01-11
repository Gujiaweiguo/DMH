package campaign

import (
	"fmt"
	"net/http"
	"strconv"

	"dmh/api/internal/svc"
	"dmh/model"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

type UpdateStatusReq struct {
	Status string `json:"status"`
}

func UpdateCampaignStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 使用go-zero的pathvar获取路径参数
		idStr := pathvar.Vars(r)["id"]
		if idStr == "" {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("活动ID不能为空"))
			return
		}
		
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 直接更新数据库
		result := svcCtx.DB.Model(&model.Campaign{}).Where("id = ?", id).Update("status", req.Status)
		if result.Error != nil {
			httpx.ErrorCtx(r.Context(), w, result.Error)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"success": true,
			"message": "状态更新成功",
		})
	}
}
