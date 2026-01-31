// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDistributorLevelRewardsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorLevelRewardsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorLevelRewardsLogic {
	return &GetDistributorLevelRewardsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorLevelRewardsLogic) GetDistributorLevelRewards() (resp *types.DistributorLevelRewardsResp, err error) {
	brandId := l.ctx.Value("brandId").(int64)

	var levelRewards []model.DistributorLevelReward
	if err := l.svcCtx.DB.Where("brand_id = ?", brandId).Find(&levelRewards).Error; err != nil {
		return nil, err
	}

	resp = &types.DistributorLevelRewardsResp{
		BrandId: brandId,
		Rewards: make([]types.DistributorLevelRewardResp, 0, len(levelRewards)),
	}

	for _, reward := range levelRewards {
		rewardResp := types.DistributorLevelRewardResp{
			Id:               reward.Id,
			BrandId:          reward.BrandId,
			Level:            reward.Level,
			RewardPercentage: reward.RewardPercentage,
		}

		resp.Rewards = append(resp.Rewards, rewardResp)
	}

	return resp, nil
}
