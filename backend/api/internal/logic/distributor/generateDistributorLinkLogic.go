// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"
	"crypto/rand"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"encoding/hex"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateDistributorLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateDistributorLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateDistributorLinkLogic {
	return &GenerateDistributorLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateDistributorLinkLogic) GenerateDistributorLink(req *types.GenerateLinkReq) (resp *types.GenerateLinkResp, err error) {
	userId := l.ctx.Value("userId").(int64)

	distributor := &model.Distributor{}
	if err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userId, req.CampaignId).First(distributor).Error; err != nil {
		return nil, fmt.Errorf("distributor not found")
	}

	linkCode := generateLinkCode()

	link := &model.DistributorLink{
		DistributorId: distributor.Id,
		CampaignId:    req.CampaignId,
		LinkCode:      linkCode,
		ClickCount:    0,
		OrderCount:    0,
		Status:        "active",
	}

	if err := l.svcCtx.DB.Create(link).Error; err != nil {
		return nil, err
	}

	resp = &types.GenerateLinkResp{
		LinkId:     link.Id,
		Link:       fmt.Sprintf("/distributor/%s", linkCode),
		LinkCode:   linkCode,
		CampaignId: req.CampaignId,
	}

	return resp, nil
}

func generateLinkCode() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
