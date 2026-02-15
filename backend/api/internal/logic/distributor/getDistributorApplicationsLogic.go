package distributor

import (
	"context"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDistributorApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDistributorApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDistributorApplicationsLogic {
	return &GetDistributorApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDistributorApplicationsLogic) GetDistributorApplications(req *types.GetDistributorApplicationsReq) (resp *types.DistributorApplicationListResp, err error) {
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

	query := l.svcCtx.DB.Model(&model.DistributorApplication{})
	if req.BrandId > 0 {
		query = query.Where("brand_id = ?", req.BrandId)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		l.Logger.Errorf("查询分销商申请总数失败: %v", err)
		return nil, err
	}

	var applications []model.DistributorApplication
	if err := query.Order("created_at DESC").Limit(int(pageSize)).Offset(int(offset)).Find(&applications).Error; err != nil {
		l.Logger.Errorf("查询分销商申请列表失败: %v", err)
		return nil, err
	}

	appList := make([]types.DistributorApplicationResp, 0, len(applications))
	for _, app := range applications {
		item := types.DistributorApplicationResp{
			Id:        app.Id,
			UserId:    app.UserId,
			BrandId:   app.BrandId,
			Status:    app.Status,
			Reason:    app.Reason,
			CreatedAt: app.CreatedAt.Format(time.RFC3339),
		}
		if app.ReviewedBy != nil {
			item.ReviewedBy = *app.ReviewedBy
		}
		if app.ReviewedAt != nil {
			item.ReviewedAt = app.ReviewedAt.Format(time.RFC3339)
		}
		item.ReviewNotes = app.ReviewNotes
		appList = append(appList, item)
	}

	return &types.DistributorApplicationListResp{
		Total:        total,
		Applications: appList,
	}, nil
}
