package distributor

import (
	"context"
	"errors"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDistributorRewardsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorRewardsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorRewardsLogic {
	return &GetDistributorRewardsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorRewardsLogic) GetDistributorRewards(req *types.GetDistributorRewardsReq) (resp *types.DistributorRewardListResp, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId <= 0 {
		return nil, errors.New("用户未登录")
	}

	query := l.svcCtx.DB.Model(&model.DistributorReward{}).Where("user_id = ?", userId)

	if req.Level > 0 {
		query = query.Where("level = ?", req.Level)
	}

	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		l.Logger.Errorf("查询分销商奖励总数失败: %v", err)
		return nil, err
	}

	var rewards []model.DistributorReward
	if err := query.Order("created_at DESC").
		Limit(int(pageSize)).Offset(int(offset)).
		Find(&rewards).Error; err != nil {
		l.Logger.Errorf("查询分销商奖励列表失败: %v", err)
		return nil, err
	}

	rewardList := make([]types.DistributorRewardResp, 0, len(rewards))
	for _, reward := range rewards {
		rewardList = append(rewardList, types.DistributorRewardResp{
			Id:        reward.Id,
			OrderId:   reward.OrderId,
			Amount:    reward.Amount,
			Level:     reward.Level,
			CreatedAt: reward.CreatedAt.Format(time.RFC3339),
		})
	}

	return &types.DistributorRewardListResp{
		Total:   total,
		Rewards: rewardList,
	}, nil
}
