// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDistributorLinksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorLinksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorLinksLogic {
	return &GetDistributorLinksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorLinksLogic) GetDistributorLinks() (resp []types.GenerateLinkResp, err error) {
	userId := l.ctx.Value("userId").(int64)

	distributor := &model.Distributor{}
	if err := l.svcCtx.DB.Where("user_id = ?", userId).First(distributor).Error; err != nil {
		return nil, err
	}

	var links []model.DistributorLink
	if err := l.svcCtx.DB.Where("distributor_id = ?", distributor.Id).Find(&links).Error; err != nil {
		return nil, err
	}

	resp = make([]types.GenerateLinkResp, 0, len(links))

	for _, link := range links {
		linkResp := types.GenerateLinkResp{
			LinkId:     link.Id,
			Link:       fmt.Sprintf("/distributor/%s", link.LinkCode),
			LinkCode:   link.LinkCode,
			CampaignId: link.CampaignId,
		}

		resp = append(resp, linkResp)
	}

	return resp, nil
}
