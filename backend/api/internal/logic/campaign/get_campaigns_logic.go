// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"fmt"
	"strings"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCampaignsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCampaignsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCampaignsLogic {
	return &GetCampaignsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCampaignsLogic) GetCampaigns(req *types.GetCampaignsReq) (resp *types.CampaignListResp, err error) {
	// 创建模拟FormField数据
	createFormFields := func(fieldTypes []string) []types.FormField {
		var fields []types.FormField
		for i, fieldType := range fieldTypes {
			switch fieldType {
			case "姓名":
				fields = append(fields, types.FormField{
					Type:        "text",
					Name:        "name",
					Label:       "姓名",
					Required:    true,
					Placeholder: "请输入姓名",
				})
			case "手机号":
				fields = append(fields, types.FormField{
					Type:        "phone",
					Name:        "phone",
					Label:       "手机号",
					Required:    true,
					Placeholder: "请输入手机号",
				})
			case "意向课程":
				fields = append(fields, types.FormField{
					Type:     "select",
					Name:     "course",
					Label:    "意向课程",
					Required: true,
					Options:  []string{"前端开发", "后端开发", "全栈开发"},
				})
			default:
				fields = append(fields, types.FormField{
					Type:        "text",
					Name:        fmt.Sprintf("field_%d", i),
					Label:       fieldType,
					Required:    false,
					Placeholder: fmt.Sprintf("请输入%s", fieldType),
				})
			}
		}
		return fields
	}

	// 模拟活动数据
	allCampaigns := []types.CampaignResp{
		{
			Id:          1,
			BrandId:     1,
			BrandName:   "示例品牌A",
			Name:        "春季新品推广活动",
			Description: "报名参与春季新品推广，完成任务可获得丰厚奖励",
			FormFields:  createFormFields([]string{"姓名", "手机号", "意向课程"}),
			RewardRule:  50.00,
			StartTime:   "2025-01-01 00:00:00",
			EndTime:     "2025-03-31 23:59:59",
			Status:      "active",
			CreatedAt:   "2025-01-01 00:00:00",
		},
		{
			Id:          2,
			BrandId:     1,
			BrandName:   "示例品牌A",
			Name:        "夏季促销推广",
			Description: "参与夏季促销活动推广，每成功推荐1人可获得30元奖励",
			FormFields:  createFormFields([]string{"姓名", "手机号"}),
			RewardRule:  30.00,
			StartTime:   "2025-04-01 00:00:00",
			EndTime:     "2025-06-30 23:59:59",
			Status:      "active",
			CreatedAt:   "2025-01-15 10:00:00",
		},
		{
			Id:          3,
			BrandId:     2,
			BrandName:   "示例品牌B",
			Name:        "会员招募计划",
			Description: "招募品牌会员，成功推荐可获得积分和现金奖励",
			FormFields:  createFormFields([]string{"姓名", "手机号", "所在城市"}),
			RewardRule:  80.00,
			StartTime:   "2025-01-01 00:00:00",
			EndTime:     "2025-12-31 23:59:59",
			Status:      "active",
			CreatedAt:   "2025-01-10 09:00:00",
		},
		{
			Id:          4,
			BrandId:     2,
			BrandName:   "示例品牌B",
			Name:        "双11狂欢预热",
			Description: "双11预热活动，提前报名享受专属优惠",
			FormFields:  createFormFields([]string{"姓名", "购买意向"}),
			RewardRule:  20.00,
			StartTime:   "2024-11-01 00:00:00",
			EndTime:     "2024-11-11 23:59:59",
			Status:      "ended",
			CreatedAt:   "2024-10-20 08:00:00",
		},
	}

	// 根据参数过滤数据
	filteredCampaigns := []types.CampaignResp{}
	for _, campaign := range allCampaigns {
		// 状态筛选
		if req.Status != "" && campaign.Status != req.Status {
			continue
		}
		// 关键词搜索
		if req.Keyword != "" {
			keyword := strings.ToLower(req.Keyword)
			if !strings.Contains(strings.ToLower(campaign.Name), keyword) &&
				!strings.Contains(strings.ToLower(campaign.Description), keyword) {
				continue
			}
		}
		filteredCampaigns = append(filteredCampaigns, campaign)
	}

	// 分页
	total := int64(len(filteredCampaigns))
	start := int64(0)
	if req.Page > 0 && req.PageSize > 0 {
		start = (req.Page - 1) * req.PageSize
		if start >= total {
			start = total
		}
	}
	end := start + req.PageSize
	if end > total {
		end = total
	}

	// 分页结果
	pagedCampaigns := []types.CampaignResp{}
	if start < total {
		pagedCampaigns = filteredCampaigns[start:end]
	}

	resp = &types.CampaignListResp{
		Total:     total,
		Campaigns: pagedCampaigns,
	}

	l.Logger.Infof("获取活动列表成功，共 %d 个活动，分页结果 %d 个", total, len(pagedCampaigns))

	return
}
