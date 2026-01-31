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

type GetBrandDistributorApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandDistributorApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandDistributorApplicationsLogic {
	return &GetBrandDistributorApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandDistributorApplicationsLogic) GetBrandDistributorApplications(req *types.GetDistributorApplicationsReq) (resp *types.DistributorApplicationListResp, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	query := l.svcCtx.DB.Model(&model.DistributorApplication{})

	if req.BrandId > 0 {
		query = query.Where("brand_id = ?", req.BrandId)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var applications []model.DistributorApplication
	offset := int((req.Page - 1) * req.PageSize)
	if err := query.Offset(offset).Limit(int(req.PageSize)).
		Preload("User").
		Preload("Brand").
		Preload("Reviewer").
		Find(&applications).Error; err != nil {
		return nil, err
	}

	resp = &types.DistributorApplicationListResp{
		Total:        total,
		Applications: make([]types.DistributorApplicationResp, 0, len(applications)),
	}

	for _, app := range applications {
		appResp := types.DistributorApplicationResp{
			Id:          app.Id,
			UserId:      app.UserId,
			BrandId:     app.BrandId,
			Status:      app.Status,
			Reason:      app.Reason,
			ReviewNotes: app.ReviewNotes,
			CreatedAt:   app.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if app.User != nil {
			appResp.Username = app.User.Username
		}

		if app.Brand != nil {
			appResp.BrandName = app.Brand.Name
		}

		if app.ReviewedBy != nil {
			appResp.ReviewedBy = *app.ReviewedBy
		}

		if app.ReviewedAt != nil {
			appResp.ReviewedAt = app.ReviewedAt.Format("2006-01-02 15:04:05")
		}

		if app.Reviewer != nil {
			appResp.Reviewer = app.Reviewer.Username
		}

		resp.Applications = append(resp.Applications, appResp)
	}

	return resp, nil
}
