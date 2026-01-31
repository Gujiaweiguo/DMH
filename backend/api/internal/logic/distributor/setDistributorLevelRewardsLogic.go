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

type SetDistributorLevelRewardsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetDistributorLevelRewardsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDistributorLevelRewardsLogic {
	return &SetDistributorLevelRewardsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetDistributorLevelRewardsLogic) SetDistributorLevelRewards(req *types.SetDistributorLevelRewardsReq) (resp *types.CommonResp, err error) {
	brandId := l.ctx.Value("brandId").(int64)

	for _, reward := range req.Rewards {
		var levelReward model.DistributorLevelReward
		result := l.svcCtx.DB.Where("brand_id = ? AND level = ?", brandId, reward.Level).First(&levelReward)

		if result.Error != nil {
			newReward := model.DistributorLevelReward{
				BrandId:          brandId,
				Level:            int(reward.Level),
				RewardPercentage: reward.RewardPercentage,
			}
			if err := l.svcCtx.DB.Create(&newReward).Error; err != nil {
				return nil, err
			}
		} else {
			levelReward.RewardPercentage = reward.RewardPercentage
			if err := l.svcCtx.DB.Save(&levelReward).Error; err != nil {
				return nil, err
			}
		}
	}

	resp = &types.CommonResp{
		Message: "Success",
	}

	return resp, nil
}
