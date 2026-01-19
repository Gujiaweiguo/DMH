package distributor

import (
	"context"
	"fmt"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatisticsLogic {
	return &StatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetStatistics 获取分销商统计数据
func (l *StatisticsLogic) GetStatistics(brandId int64) (*types.DistributorStatisticsResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 获取分销商信息
	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userId, brandId).First(&distributor).Error; err != nil {
		return nil, fmt.Errorf("分销商不存在")
	}

	// 获取总订单数
	var totalOrders int64
	l.svcCtx.DB.Model(&model.DistributorReward{}).
		Where("distributor_id = ?", distributor.Id).
		Count(&totalOrders)

	// 获取总收益
	var totalEarnings float64
	l.svcCtx.DB.Model(&model.DistributorReward{}).
		Where("distributor_id = ? AND status = ?", distributor.Id, "settled").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalEarnings)

	// 获取今日收益
	today := time.Now().Format("2006-01-02")
	var todayEarnings float64
	l.svcCtx.DB.Model(&model.DistributorReward{}).
		Where("distributor_id = ? AND status = ? AND DATE(created_at) = ?", distributor.Id, "settled", today).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&todayEarnings)

	// 获取本月收益
	monthStart := time.Now().Format("2006-01-01")
	var monthEarnings float64
	l.svcCtx.DB.Model(&model.DistributorReward{}).
		Where("distributor_id = ? AND status = ? AND created_at >= ?", distributor.Id, "settled", monthStart).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&monthEarnings)

	// 获取点击总数
	var clickCount int64
	l.svcCtx.DB.Model(&model.DistributorLink{}).
		Where("distributor_id = ?", distributor.Id).
		Select("COALESCE(SUM(click_count), 0)").
		Scan(&clickCount)

	// 计算转化率
	conversionRate := 0.0
	if clickCount > 0 {
		conversionRate = float64(totalOrders) / float64(clickCount) * 100
	}

	return &types.DistributorStatisticsResp{
		DistributorId:    distributor.Id,
		TotalOrders:      totalOrders,
		TotalEarnings:    totalEarnings,
		TodayEarnings:    todayEarnings,
		MonthEarnings:    monthEarnings,
		SubordinatesCount: distributor.SubordinatesCount,
		ClickCount:       int(clickCount),
		ConversionRate:   conversionRate,
	}, nil
}

// GetRewards 获取奖励明细
func (l *StatisticsLogic) GetRewards(brandId int64, req *types.GetDistributorRewardsReq) (*types.DistributorRewardListResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 获取分销商信息
	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userId, brandId).First(&distributor).Error; err != nil {
		return nil, fmt.Errorf("分销商不存在")
	}

	query := l.svcCtx.DB.Model(&model.DistributorReward{}).Where("distributor_id = ?", distributor.Id)

	// 筛选条件
	if req.Level > 0 {
		query = query.Where("level = ?", req.Level)
	}
	if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	var rewards []model.DistributorReward
	offset := (req.Page - 1) * req.PageSize
	query.Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&rewards)

	resp := &types.DistributorRewardListResp{
		Total:   total,
		Rewards: make([]types.DistributorRewardResp, 0),
	}

	for _, reward := range rewards {
		var fromUsername string
		if reward.FromUserId != nil {
			var fromUser model.User
			if l.svcCtx.DB.Where("id = ?", *reward.FromUserId).First(&fromUser).Error == nil {
				fromUsername = fromUser.RealName
				if fromUsername == "" {
					fromUsername = fromUser.Username
				}
			}
		}

		rewardResp := types.DistributorRewardResp{
			Id:         reward.Id,
			OrderId:    reward.OrderId,
			Amount:     reward.Amount,
			Level:      reward.Level,
			RewardRate: reward.RewardRate,
			Status:     reward.Status,
			CreatedAt:  reward.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if reward.FromUserId != nil {
			rewardResp.FromUserId = *reward.FromUserId
			rewardResp.FromUsername = fromUsername
		}

		if reward.SettledAt != nil {
			rewardResp.SettledAt = reward.SettledAt.Format("2006-01-02 15:04:05")
		}

		resp.Rewards = append(resp.Rewards, rewardResp)
	}

	return resp, nil
}

// GetSubordinates 获取下级分销商列表
func (l *StatisticsLogic) GetSubordinates(brandId int64, page, pageSize int64) (*types.SubordinateListResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 获取当前用户的分销商信息
	var currentDistributor model.Distributor
	if err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userId, brandId).First(&currentDistributor).Error; err != nil {
		return nil, fmt.Errorf("分销商不存在")
	}

	query := l.svcCtx.DB.Model(&model.Distributor{}).Where("parent_id = ? AND brand_id = ?", currentDistributor.Id, brandId)

	var total int64
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var subordinates []model.Distributor
	offset := (page - 1) * pageSize
	query.Preload("User").
		Order("created_at DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Find(&subordinates)

	resp := &types.SubordinateListResp{
		Total:        total,
		Subordinates: make([]types.SubordinateResp, 0),
	}

	for _, sub := range subordinates {
		username := ""
		if sub.User != nil {
			username = sub.User.RealName
			if username == "" {
				username = sub.User.Username
			}
		}

		// 获取下级的订单数
		var orderCount int64
		l.svcCtx.DB.Model(&model.DistributorReward{}).
			Where("distributor_id = ?", sub.Id).
			Count(&orderCount)

		// 获取下级的总收益
		var earnings float64
		l.svcCtx.DB.Model(&model.DistributorReward{}).
			Where("distributor_id = ? AND status = ?", sub.Id, "settled").
			Select("COALESCE(SUM(amount), 0)").
			Scan(&earnings)

		resp.Subordinates = append(resp.Subordinates, types.SubordinateResp{
			Id:           sub.Id,
			UserId:       sub.UserId,
			Username:     username,
			Level:        sub.Level,
			TotalOrders:  int(orderCount),
			TotalEarnings: earnings,
			CreatedAt:    sub.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

// GetDashboard 获取分销商仪表盘数据
func (l *StatisticsLogic) GetDashboard() (map[string]interface{}, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 获取用户的所有分销商记录
	var distributors []model.Distributor
	if err := l.svcCtx.DB.Preload("Brand").Where("user_id = ?", userId).Find(&distributors).Error; err != nil {
		return nil, fmt.Errorf("查询分销商记录失败")
	}

	if len(distributors) == 0 {
		return map[string]interface{}{
			"hasDistributor": false,
			"brands":         []map[string]interface{}{},
		}, nil
	}

	brands := make([]map[string]interface{}, 0, len(distributors))
	totalEarnings := 0.0

	for _, d := range distributors {
		// 获取该品牌的收益
		var earnings float64
		l.svcCtx.DB.Model(&model.DistributorReward{}).
			Where("distributor_id = ? AND status = ?", d.Id, "settled").
			Select("COALESCE(SUM(amount), 0)").
			Scan(&earnings)

		totalEarnings += earnings

		brandName := ""
		if d.Brand != nil {
			brandName = d.Brand.Name
		}

		brands = append(brands, map[string]interface{}{
			"brandId":       d.BrandId,
			"brandName":     brandName,
			"level":         d.Level,
			"status":        d.Status,
			"earnings":      earnings,
			"subordinates":  d.SubordinatesCount,
		})
	}

	return map[string]interface{}{
		"hasDistributor": true,
		"totalEarnings":  totalEarnings,
		"brands":         brands,
	}, nil
}

// GetLevelRewards 获取品牌级别奖励配置
func (l *StatisticsLogic) GetLevelRewards(brandId int64) (*types.DistributorLevelRewardsResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 检查权限
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 品牌管理员只能查看自己品牌的配置
	if user.Role != "platform_admin" {
		var count int64
		l.svcCtx.DB.Model(&model.UserBrand{}).
			Where("user_id = ? AND brand_id = ?", userId, brandId).
			Count(&count)
		if count == 0 {
			return nil, fmt.Errorf("无权限查看该品牌的配置")
		}
	}

	var rewards []model.DistributorLevelReward
	if err := l.svcCtx.DB.Where("brand_id = ?", brandId).Order("level ASC").Find(&rewards).Error; err != nil {
		return nil, fmt.Errorf("查询奖励配置失败")
	}

	rewardResps := make([]types.DistributorLevelRewardResp, 0, len(rewards))
	for _, r := range rewards {
		rewardResps = append(rewardResps, types.DistributorLevelRewardResp{
			Id:               r.Id,
			BrandId:          r.BrandId,
			Level:            r.Level,
			RewardPercentage: r.RewardPercentage,
		})
	}

	return &types.DistributorLevelRewardsResp{
		BrandId: brandId,
		Rewards: rewardResps,
	}, nil
}

// SetLevelRewards 设置品牌级别奖励配置
func (l *StatisticsLogic) SetLevelRewards(brandId int64, req *types.SetDistributorLevelRewardsReq) error {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return fmt.Errorf("获取用户信息失败")
	}

	// 检查权限
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return fmt.Errorf("用户不存在")
	}

	// 品牌管理员只能修改自己品牌的配置
	if user.Role != "platform_admin" {
		var count int64
		l.svcCtx.DB.Model(&model.UserBrand{}).
			Where("user_id = ? AND brand_id = ?", userId, brandId).
			Count(&count)
		if count == 0 {
			return fmt.Errorf("无权限修改该品牌的配置")
		}
	}

	// 验证级别和比例
	for _, r := range req.Rewards {
		if r.Level < 1 || r.Level > 3 {
			return fmt.Errorf("级别必须在1-3之间")
		}
		if r.RewardPercentage < 0 || r.RewardPercentage > 100 {
			return fmt.Errorf("奖励比例必须在0-100之间")
		}
	}

	// 开始事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除旧的配置
	if err := tx.Where("brand_id = ?", brandId).Delete(&model.DistributorLevelReward{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除旧配置失败")
	}

	// 创建新的配置
	for _, r := range req.Rewards {
		reward := model.DistributorLevelReward{
			BrandId:          brandId,
			Level:            int(r.Level),
			RewardPercentage: r.RewardPercentage,
		}
		if err := tx.Create(&reward).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建配置失败")
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败")
	}

	return nil
}
