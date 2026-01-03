// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCampaignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCampaignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCampaignLogic {
	return &CreateCampaignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCampaignLogic) CreateCampaign(req *types.CreateCampaignReq) (resp *types.CampaignResp, err error) {
	// 验证FormFields结构
	for _, field := range req.FormFields {
		if field.Type == "" || field.Name == "" || field.Label == "" {
			return nil, fmt.Errorf("FormField必须包含type、name和label字段")
		}
		
		// 验证字段类型
		if field.Type != "text" && field.Type != "phone" && field.Type != "select" {
			return nil, fmt.Errorf("不支持的字段类型: %s", field.Type)
		}
		
		// select类型必须有选项
		if field.Type == "select" && len(field.Options) == 0 {
			return nil, fmt.Errorf("select类型字段必须提供选项列表")
		}
	}
	
	// 生成新的活动ID（模拟）
	newId := int64(100 + len(req.FormFields))
	
	// 返回创建成功的活动信息
	return &types.CampaignResp{
		Id:          newId,
		BrandId:     req.BrandId,
		BrandName:   "测试品牌",
		Name:        req.Name,
		Description: req.Description,
		FormFields:  req.FormFields,
		RewardRule:  req.RewardRule,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      "active",
		CreatedAt:   "2026-01-01 16:30:00",
	}, nil
}
