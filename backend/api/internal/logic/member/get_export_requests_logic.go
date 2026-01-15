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

type GetExportRequestsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetExportRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExportRequestsLogic {
	return &GetExportRequestsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetExportRequestsLogic) GetExportRequests(req *types.GetExportRequestsReq) (*types.ExportRequestListResp, error) {
	userId, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

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

	// 获取用户信息
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 构建查询
	query := l.svcCtx.DB.Model(&model.ExportRequest{})

	// 品牌管理员只能查看自己品牌的导出请求
	if user.Role != "platform_admin" {
		var userBrands []model.UserBrand
		if err := l.svcCtx.DB.Where("user_id = ?", userId).Find(&userBrands).Error; err != nil {
			return nil, fmt.Errorf("获取用户品牌关联失败: %w", err)
		}

		if len(userBrands) == 0 {
			return &types.ExportRequestListResp{Total: 0, Requests: []types.ExportRequestResp{}}, nil
		}

		brandIds := make([]int64, len(userBrands))
		for i, ub := range userBrands {
			brandIds[i] = ub.BrandId
		}

		query = query.Where("brand_id IN ?", brandIds)
	}

	// 品牌筛选
	if req.BrandId > 0 {
		query = query.Where("brand_id = ?", req.BrandId)
	}

	// 状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("查询导出请求总数失败: %w", err)
	}

	// 查询列表
	var requests []model.ExportRequest
	if err := query.Offset(int(offset)).Limit(int(pageSize)).
		Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, fmt.Errorf("查询导出请求列表失败: %w", err)
	}

	// 查询关联的品牌和用户信息
	brandIds := make([]int64, 0)
	userIds := make([]int64, 0)
	for _, r := range requests {
		brandIds = append(brandIds, r.BrandID)
		userIds = append(userIds, r.RequestedBy)
		if r.ApprovedBy != nil {
			userIds = append(userIds, *r.ApprovedBy)
		}
	}

	var brands []model.Brand
	if len(brandIds) > 0 {
		l.svcCtx.DB.Where("id IN ?", brandIds).Find(&brands)
	}

	var users []model.User
	if len(userIds) > 0 {
		l.svcCtx.DB.Where("id IN ?", userIds).Find(&users)
	}

	brandMap := make(map[int64]model.Brand)
	for _, b := range brands {
		brandMap[b.Id] = b
	}

	userMap := make(map[int64]model.User)
	for _, u := range users {
		userMap[u.Id] = u
	}

	// 构建响应
	respRequests := make([]types.ExportRequestResp, len(requests))
	for i, r := range requests {
		brand := brandMap[r.BrandID]
		requestedByUser := userMap[r.RequestedBy]

		respRequests[i] = types.ExportRequestResp{
			Id:              r.ID,
			BrandId:         r.BrandID,
			BrandName:       brand.Name,
			RequestedBy:     r.RequestedBy,
			RequestedByName: requestedByUser.Username,
			Reason:          r.Reason,
			Filters:         r.Filters,
			Status:          r.Status,
			RejectReason:    r.RejectReason,
			FileUrl:         r.FileUrl,
			RecordCount:     r.RecordCount,
			CreatedAt:       r.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if r.ApprovedBy != nil {
			respRequests[i].ApprovedBy = *r.ApprovedBy
			approvedByUser := userMap[*r.ApprovedBy]
			respRequests[i].ApprovedByName = approvedByUser.Username
		}

		if r.ApprovedAt != nil {
			respRequests[i].ApprovedAt = r.ApprovedAt.Format("2006-01-02 15:04:05")
		}
	}

	return &types.ExportRequestListResp{
		Total:    total,
		Requests: respRequests,
	}, nil
}
