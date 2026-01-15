package member

import (
	"context"
	"fmt"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMembersLogic {
	return &GetMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMembersLogic) GetMembers(req *types.GetMembersReq) (*types.MemberListResp, error) {
	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 获取当前用户信息（从 context 中）
	userId, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	
	// 查询用户角色
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 构建查询
	query := l.svcCtx.DB.Model(&model.Member{})

	// 品牌管理员只能查看关联品牌的会员
	if user.Role != "platform_admin" {
		// 获取用户关联的品牌
		var userBrands []model.UserBrand
		if err := l.svcCtx.DB.Where("user_id = ?", userId).Find(&userBrands).Error; err != nil {
			return nil, fmt.Errorf("获取用户品牌关联失败: %w", err)
		}
		
		if len(userBrands) == 0 {
			return &types.MemberListResp{Total: 0, Members: []types.MemberResp{}}, nil
		}

		brandIds := make([]int64, len(userBrands))
		for i, ub := range userBrands {
			brandIds[i] = ub.BrandId
		}

		// 只查询关联品牌的会员
		query = query.Joins("INNER JOIN member_brand_links ON members.id = member_brand_links.member_id").
			Where("member_brand_links.brand_id IN ?", brandIds)
	}

	// 品牌筛选（平台管理员可用）
	if req.BrandId > 0 {
		query = query.Joins("INNER JOIN member_brand_links mbl ON members.id = mbl.member_id").
			Where("mbl.brand_id = ?", req.BrandId)
	}

	// 关键词搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("nickname LIKE ? OR phone LIKE ? OR unionid LIKE ?", keyword, keyword, keyword)
	}

	// 来源筛选
	if req.Source != "" {
		query = query.Where("source = ?", req.Source)
	}

	// 状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 标签筛选
	if len(req.TagIds) > 0 {
		query = query.Joins("INNER JOIN member_tag_links mtl ON members.id = mtl.member_id").
			Where("mtl.tag_id IN ?", req.TagIds)
	}

	// 时间范围筛选
	if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate)
	}

	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("查询会员总数失败: %w", err)
	}

	// 查询会员列表
	var members []model.Member
	if err := query.Offset(int(offset)).Limit(int(pageSize)).
		Order("created_at DESC").Find(&members).Error; err != nil {
		return nil, fmt.Errorf("查询会员列表失败: %w", err)
	}

	// 查询会员画像
	memberIds := make([]int64, len(members))
	for i, m := range members {
		memberIds[i] = m.ID
	}

	var profiles []model.MemberProfile
	if len(memberIds) > 0 {
		l.svcCtx.DB.Where("member_id IN ?", memberIds).Find(&profiles)
	}

	profileMap := make(map[int64]model.MemberProfile)
	for _, p := range profiles {
		profileMap[p.MemberID] = p
	}

	// 构建响应
	respMembers := make([]types.MemberResp, len(members))
	for i, m := range members {
		profile := profileMap[m.ID]
		
		respMembers[i] = types.MemberResp{
			Id:                    m.ID,
			UnionID:               m.UnionID,
			Nickname:              m.Nickname,
			Avatar:                m.Avatar,
			Phone:                 m.Phone,
			Gender:                m.Gender,
			Source:                m.Source,
			Status:                m.Status,
			TotalOrders:           profile.TotalOrders,
			TotalPayment:          profile.TotalPayment,
			TotalReward:           profile.TotalReward,
			ParticipatedCampaigns: profile.ParticipatedCampaigns,
			CreatedAt:             m.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if profile.FirstOrderAt != nil {
			respMembers[i].FirstOrderAt = profile.FirstOrderAt.Format("2006-01-02 15:04:05")
		}
		if profile.LastOrderAt != nil {
			respMembers[i].LastOrderAt = profile.LastOrderAt.Format("2006-01-02 15:04:05")
		}
	}

	return &types.MemberListResp{
		Total:   total,
		Members: respMembers,
	}, nil
}
