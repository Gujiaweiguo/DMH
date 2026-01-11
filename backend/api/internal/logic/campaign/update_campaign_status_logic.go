package campaign

import (
	"context"
	"net/http"
	"strconv"

	"dmh/api/internal/svc"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCampaignStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCampaignStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCampaignStatusLogic {
	return &UpdateCampaignStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCampaignStatusLogic) UpdateCampaignStatus(r *http.Request, status string) (resp interface{}, err error) {
	// 从URL获取活动ID
	idStr := r.PathValue("id")
	if idStr == "" {
		// 兼容旧版本
		idStr = r.URL.Query().Get("id")
	}
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	// 更新活动状态
	result := l.svcCtx.DB.Model(&model.Campaign{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return nil, result.Error
	}

	return map[string]interface{}{
		"success": true,
		"message": "状态更新成功",
	}, nil
}