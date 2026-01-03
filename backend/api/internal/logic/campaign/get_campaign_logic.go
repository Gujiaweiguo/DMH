// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCampaignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCampaignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCampaignLogic {
	return &GetCampaignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCampaignLogic) GetCampaign(id string) (resp *types.CampaignResp, err error) {
	// 模拟FormField数据
	formFields := []types.FormField{
		{
			Type:        "text",
			Name:        "name",
			Label:       "姓名",
			Required:    true,
			Placeholder: "请输入姓名",
		},
		{
			Type:        "phone",
			Name:        "phone",
			Label:       "手机号",
			Required:    true,
			Placeholder: "请输入手机号",
		},
		{
			Type:     "select",
			Name:     "course",
			Label:    "意向课程",
			Required: true,
			Options:  []string{"前端开发", "后端开发", "全栈开发"},
		},
	}
	
	// 返回模拟的活动详情
	return &types.CampaignResp{
		Id:          1,
		BrandId:     1,
		BrandName:   "示例品牌A",
		Name:        "春季新品推广活动",
		Description: "报名参与春季新品推广，完成任务可获得丰厚奖励",
		FormFields:  formFields,
		RewardRule:  50.0,
		StartTime:   "2025-01-01 00:00:00",
		EndTime:     "2025-03-31 23:59:59",
		Status:      "active",
		CreatedAt:   "2025-01-01 00:00:00",
	}, nil
}
