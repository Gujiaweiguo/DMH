package distributor

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type PromotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PromotionLogic {
	return &PromotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GenerateLink 生成推广链接
func (l *PromotionLogic) GenerateLink(req *types.GenerateLinkReq) (*types.GenerateLinkResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 检查活动是否存在
	var campaign model.Campaign
	if err := l.svcCtx.DB.Where("id = ? AND status = ?", req.CampaignId, "active").First(&campaign).Error; err != nil {
		return nil, fmt.Errorf("活动不存在或已结束")
	}

	// 检查用户是否是该品牌的分销商
	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ? AND status = ?", userId, campaign.BrandId, "active").
		First(&distributor).Error; err != nil {
		return nil, fmt.Errorf("您不是该品牌的分销商或分销商资格未激活")
	}

	linkCode := ""
	linkId := int64(0)

	// 检查是否已有该活动的推广链接
	var existingLink model.DistributorLink
	if err := l.svcCtx.DB.Where("distributor_id = ? AND campaign_id = ?", distributor.Id, req.CampaignId).
		First(&existingLink).Error; err == nil {
		// 已有链接，更新使用时间
		linkId = existingLink.Id
		linkCode = existingLink.LinkCode
	} else {
		// 生成新的推广码
		linkCode = l.generateLinkCode(distributor.Id, req.CampaignId)

		// 创建推广链接记录
		newLink := model.DistributorLink{
			DistributorId: distributor.Id,
			CampaignId:    req.CampaignId,
			LinkCode:      linkCode,
			Status:        "active",
		}

		if err := l.svcCtx.DB.Create(&newLink).Error; err != nil {
			l.Logger.Errorf("创建推广链接失败: %v", err)
			return nil, fmt.Errorf("创建推广链接失败")
		}
		linkId = newLink.Id
	}

	// 生成推广链接
	// 格式: {campaign_url}?distributor_code={linkCode}
	baseUrl := l.getCampaignBaseUrl()
	link := fmt.Sprintf("%s/campaign/%d?distributor_code=%s", baseUrl, req.CampaignId, linkCode)

	// 生成二维码URL
	qrcodeUrl := fmt.Sprintf("%s/api/v1/distributor/qrcode/%s", baseUrl, linkCode)

	return &types.GenerateLinkResp{
		LinkId:     linkId,
		Link:       link,
		LinkCode:   linkCode,
		QrcodeUrl:  qrcodeUrl,
		CampaignId: req.CampaignId,
	}, nil
}

// GetQrcode 获取推广二维码
func (l *PromotionLogic) GetQrcode(linkCode string) (*types.GetQrcodeResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 查找推广链接记录
	var link model.DistributorLink
	if err := l.svcCtx.DB.Where("link_code = ?", linkCode).First(&link).Error; err != nil {
		return nil, fmt.Errorf("推广链接不存在")
	}

	// 验证是否是自己的推广链接
	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("id = ? AND user_id = ?", link.DistributorId, userId).First(&distributor).Error; err != nil {
		return nil, fmt.Errorf("无权访问该推广链接")
	}

	// 生成二维码URL
	baseUrl := l.getCampaignBaseUrl()
	qrcodeUrl := fmt.Sprintf("%s/api/v1/distributor/qrcode/%s", baseUrl, linkCode)

	return &types.GetQrcodeResp{
		QrcodeUrl: qrcodeUrl,
		LinkCode:  linkCode,
	}, nil
}

// TrackLinkClick 记录推广链接点击（公开接口，不需要认证）
func (l *PromotionLogic) TrackLinkClick(linkCode string) (int64, error) {
	var link model.DistributorLink
	if err := l.svcCtx.DB.Where("link_code = ? AND status = ?", linkCode, "active").First(&link).Error; err != nil {
		return 0, fmt.Errorf("推广链接不存在")
	}

	// 更新点击次数
	if err := l.svcCtx.DB.Model(&link).UpdateColumn("click_count", link.ClickCount+1).Error; err != nil {
		return 0, fmt.Errorf("更新点击次数失败")
	}

	// 返回分销商ID，用于设置推荐关系
	return link.DistributorId, nil
}

// GetDistributorByCode 根据推广码获取分销商信息
func (l *PromotionLogic) GetDistributorByCode(linkCode string) (*model.Distributor, error) {
	var link model.DistributorLink
	if err := l.svcCtx.DB.Where("link_code = ? AND status = ?", linkCode, "active").First(&link).Error; err != nil {
		return nil, fmt.Errorf("推广链接不存在")
	}

	var distributor model.Distributor
	if err := l.svcCtx.DB.Preload("User").Preload("Brand").Where("id = ? AND status = ?", link.DistributorId, "active").First(&distributor).Error; err != nil {
		return nil, fmt.Errorf("分销商不存在或未激活")
	}

	return &distributor, nil
}

// GetMyLinks 获取我的推广链接列表
func (l *PromotionLogic) GetMyLinks() ([]types.GenerateLinkResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 查询用户的所有分销商记录
	var distributors []model.Distributor
	if err := l.svcCtx.DB.Where("user_id = ? AND status = ?", userId, "active").Find(&distributors).Error; err != nil {
		return nil, fmt.Errorf("查询分销商记录失败")
	}

	if len(distributors) == 0 {
		return []types.GenerateLinkResp{}, nil
	}

	distributorIds := make([]int64, len(distributors))
	for i, d := range distributors {
		distributorIds[i] = d.Id
	}

	// 查询推广链接
	var links []model.DistributorLink
	if err := l.svcCtx.DB.Where("distributor_id IN ? AND status = ?", distributorIds, "active").
		Preload("Campaign").
		Order("created_at DESC").
		Find(&links).Error; err != nil {
		return nil, fmt.Errorf("查询推广链接失败")
	}

	baseUrl := l.getCampaignBaseUrl()
	resp := make([]types.GenerateLinkResp, 0, len(links))

	for _, link := range links {
		campaignUrl := fmt.Sprintf("%s/campaign/%d?distributor_code=%s", baseUrl, link.CampaignId, link.LinkCode)
		qrcodeUrl := fmt.Sprintf("%s/api/v1/distributor/qrcode/%s", baseUrl, link.LinkCode)

		resp = append(resp, types.GenerateLinkResp{
			LinkId:     link.Id,
			Link:       campaignUrl,
			LinkCode:   link.LinkCode,
			QrcodeUrl:  qrcodeUrl,
			CampaignId: link.CampaignId,
		})
	}

	return resp, nil
}

// GetLinkStatistics 获取推广链接统计
func (l *PromotionLogic) GetLinkStatistics(linkId int64) (map[string]interface{}, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 查找推广链接记录
	var link model.DistributorLink
	if err := l.svcCtx.DB.Preload("Campaign").Where("id = ?", linkId).First(&link).Error; err != nil {
		return nil, fmt.Errorf("推广链接不存在")
	}

	// 验证是否是自己的推广链接
	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("id = ? AND user_id = ?", link.DistributorId, userId).First(&distributor).Error; err != nil {
		return nil, fmt.Errorf("无权访问该推广链接")
	}

	// 获取活动名称
	campaignName := ""
	if link.Campaign != nil {
		campaignName = link.Campaign.Name
	}

	return map[string]interface{}{
		"linkId":      link.Id,
		"campaignId":  link.CampaignId,
		"campaignName": campaignName,
		"linkCode":    link.LinkCode,
		"clickCount":  link.ClickCount,
		"orderCount":  link.OrderCount,
		"status":      link.Status,
		"createdAt":   link.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// generateLinkCode 生成推广码
func (l *PromotionLogic) generateLinkCode(distributorId, campaignId int64) string {
	// 使用 distributorId + campaignId + timestamp 生成唯一码
	data := fmt.Sprintf("%d-%d-%d", distributorId, campaignId, time.Now().UnixNano())
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])[:12] // 取前12位
}

// getCampaignBaseUrl 获取活动基础URL
func (l *PromotionLogic) getCampaignBaseUrl() string {
	// TODO: 从配置中获取基础URL
	// 暂时返回默认值
	return "https://dmh.example.com"
}

// RecordOrderFromLink 记录通过推广链接产生的订单
func (l *PromotionLogic) RecordOrderFromLink(linkCode string, orderId int64) error {
	var link model.DistributorLink
	if err := l.svcCtx.DB.Where("link_code = ?", linkCode).First(&link).Error; err != nil {
		return fmt.Errorf("推广链接不存在")
	}

	// 更新订单数量
	if err := l.svcCtx.DB.Model(&link).UpdateColumn("order_count", link.OrderCount+1).Error; err != nil {
		return fmt.Errorf("更新订单数量失败")
	}

	return nil
}
