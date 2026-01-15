package member

import (
	"context"
	"encoding/json"
	"fmt"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMemberLogic {
	return &GetMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMemberLogic) GetMember(memberId int64) (*types.MemberResp, error) {
	// 获取当前用户信息
	userId, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 查询会员
	var member model.Member
	if err := l.svcCtx.DB.Where("id = ?", memberId).First(&member).Error; err != nil {
		return nil, fmt.Errorf("会员不存在: %w", err)
	}

	// 品牌管理员权限检查
	var allowedBrandIds []int64
	if user.Role != "platform_admin" {
		// 检查会员是否属于该品牌管理员的品牌
		var userBrands []model.UserBrand
		if err := l.svcCtx.DB.Where("user_id = ?", userId).Find(&userBrands).Error; err != nil {
			return nil, fmt.Errorf("获取用户品牌关联失败: %w", err)
		}

		brandIds := make([]int64, len(userBrands))
		for i, ub := range userBrands {
			brandIds[i] = ub.BrandId
		}
		allowedBrandIds = brandIds

		var count int64
		l.svcCtx.DB.Model(&model.MemberBrandLink{}).
			Where("member_id = ? AND brand_id IN ?", memberId, brandIds).
			Count(&count)

		if count == 0 {
			return nil, fmt.Errorf("无权查看该会员")
		}
	}

	// 查询会员画像
	var profile model.MemberProfile
	l.svcCtx.DB.Where("member_id = ?", memberId).First(&profile)

	// 查询会员标签
	var tagLinks []model.MemberTagLink
	l.svcCtx.DB.Where("member_id = ?", memberId).Find(&tagLinks)

	tagIds := make([]int64, len(tagLinks))
	for i, tl := range tagLinks {
		tagIds[i] = tl.TagID
	}

	var tags []model.MemberTag
	if len(tagIds) > 0 {
		l.svcCtx.DB.Where("id IN ?", tagIds).Find(&tags)
	}

	respTags := make([]types.MemberTagResp, len(tags))
	for i, tag := range tags {
		respTags[i] = types.MemberTagResp{
			Id:          tag.ID,
			Name:        tag.Name,
			Category:    tag.Category,
			Color:       tag.Color,
			Description: tag.Description,
		}
	}

	// 查询会员品牌关联
	var brandLinks []model.MemberBrandLink
	l.svcCtx.DB.Where("member_id = ?", memberId).Find(&brandLinks)

	brandIds := make([]int64, len(brandLinks))
	for i, bl := range brandLinks {
		brandIds[i] = bl.BrandID
	}

	var brands []model.Brand
	if len(brandIds) > 0 {
		l.svcCtx.DB.Where("id IN ?", brandIds).Find(&brands)
	}

	brandMap := make(map[int64]model.Brand)
	for _, b := range brands {
		brandMap[b.Id] = b
	}

	respBrands := make([]types.MemberBrandResp, len(brandLinks))
	for i, bl := range brandLinks {
		brand := brandMap[bl.BrandID]
		respBrands[i] = types.MemberBrandResp{
			BrandId:         bl.BrandID,
			BrandName:       brand.Name,
			FirstCampaignId: bl.FirstCampaignID,
			CreatedAt:       bl.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// 构建响应
	resp := &types.MemberResp{
		Id:                    member.ID,
		UnionID:               member.UnionID,
		Nickname:              member.Nickname,
		Avatar:                member.Avatar,
		Phone:                 member.Phone,
		Gender:                member.Gender,
		Source:                member.Source,
		Status:                member.Status,
		TotalOrders:           profile.TotalOrders,
		TotalPayment:          profile.TotalPayment,
		TotalReward:           profile.TotalReward,
		ParticipatedCampaigns: profile.ParticipatedCampaigns,
		CreatedAt:             member.CreatedAt.Format("2006-01-02 15:04:05"),
		Tags:                  respTags,
		Brands:                respBrands,
	}

	if profile.FirstOrderAt != nil {
		resp.FirstOrderAt = profile.FirstOrderAt.Format("2006-01-02 15:04:05")
	}
	if profile.LastOrderAt != nil {
		resp.LastOrderAt = profile.LastOrderAt.Format("2006-01-02 15:04:05")
	}

	// 查询会员订单
	var orders []model.Order
	orderQuery := l.svcCtx.DB.Model(&model.Order{}).Where("member_id = ?", memberId)
	if user.Role != "platform_admin" && len(allowedBrandIds) > 0 {
		orderQuery = orderQuery.Joins("INNER JOIN campaigns ON campaigns.id = orders.campaign_id").
			Where("campaigns.brand_id IN ?", allowedBrandIds)
	}
	if err := orderQuery.Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("查询会员订单失败: %w", err)
	}

	campaignIds := make([]int64, 0, len(orders))
	for _, order := range orders {
		campaignIds = append(campaignIds, order.CampaignId)
	}
	campaignMap := make(map[int64]string)
	if len(campaignIds) > 0 {
		var campaigns []model.Campaign
		l.svcCtx.DB.Where("id IN ?", campaignIds).Find(&campaigns)
		for _, c := range campaigns {
			campaignMap[c.Id] = c.Name
		}
	}

	respOrders := make([]types.OrderResp, 0, len(orders))
	for _, order := range orders {
		formData := make(map[string]string)
		if order.FormData != "" {
			var raw map[string]interface{}
			if err := json.Unmarshal([]byte(order.FormData), &raw); err == nil {
				for key, value := range raw {
					formData[key] = fmt.Sprint(value)
				}
			}
		}
		respOrders = append(respOrders, types.OrderResp{
			Id:         order.Id,
			CampaignId: order.CampaignId,
			CampaignName: campaignMap[order.CampaignId],
			Phone:      order.Phone,
			FormData:   formData,
			ReferrerId: order.ReferrerId,
			Status:     order.Status,
			Amount:     order.Amount,
			CreatedAt:  order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	resp.Orders = respOrders

	return resp, nil
}
