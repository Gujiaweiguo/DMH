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

type ApproveDistributorApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveDistributorApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveDistributorApplicationLogic {
	return &ApproveDistributorApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveDistributorApplicationLogic) ApproveDistributorApplication(req *types.ApproveDistributorReq) (resp *types.DistributorApplicationResp, err error) {
	application := &model.DistributorApplication{}
	if err := l.svcCtx.DB.First(application, l.ctx.Value("applicationId")).Error; err != nil {
		return nil, err
	}

	reviewerId := l.ctx.Value("userId").(int64)

	application.Status = req.Action
	application.ReviewedBy = &reviewerId
	application.ReviewNotes = req.Reason

	now := l.svcCtx.DB.NowFunc()
	application.ReviewedAt = &now

	if err := l.svcCtx.DB.Save(application).Error; err != nil {
		return nil, err
	}

	if req.Action == "approved" {
		existingDistributor := &model.Distributor{}
		result := l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", application.UserId, application.BrandId).First(existingDistributor)

		if result.Error != nil {
			newDistributor := &model.Distributor{
				UserId:            application.UserId,
				BrandId:           application.BrandId,
				Level:             req.Level,
				Status:            "active",
				ApprovedBy:        &reviewerId,
				ApprovedAt:        &now,
				TotalEarnings:     0,
				SubordinatesCount: 0,
			}

			if err := l.svcCtx.DB.Create(newDistributor).Error; err != nil {
				return nil, err
			}
		} else {
			existingDistributor.Status = "active"
			existingDistributor.Level = req.Level
			existingDistributor.ApprovedBy = &reviewerId
			existingDistributor.ApprovedAt = &now

			if err := l.svcCtx.DB.Save(existingDistributor).Error; err != nil {
				return nil, err
			}
		}
	}

	if err := l.svcCtx.DB.Preload("User").Preload("Brand").Preload("Reviewer").First(application, application.Id).Error; err != nil {
		return nil, err
	}

	resp = &types.DistributorApplicationResp{
		Id:          application.Id,
		UserId:      application.UserId,
		BrandId:     application.BrandId,
		Status:      application.Status,
		Reason:      application.Reason,
		ReviewNotes: application.ReviewNotes,
		CreatedAt:   application.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if application.User != nil {
		resp.Username = application.User.Username
	}

	if application.Brand != nil {
		resp.BrandName = application.Brand.Name
	}

	if application.ReviewedBy != nil {
		resp.ReviewedBy = *application.ReviewedBy
	}

	if application.ReviewedAt != nil {
		resp.ReviewedAt = application.ReviewedAt.Format("2006-01-02 15:04:05")
	}

	if application.Reviewer != nil {
		resp.Reviewer = application.Reviewer.Username
	}

	return resp, nil
}
