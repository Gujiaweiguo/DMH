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

type GetBrandDistributorApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandDistributorApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandDistributorApplicationLogic {
	return &GetBrandDistributorApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandDistributorApplicationLogic) GetBrandDistributorApplication(req *types.GetDistributorApplicationReq) (resp *types.DistributorApplicationResp, err error) {
	if req.Id <= 0 {
		return nil, errors.New("申请ID无效")
	}

	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId <= 0 {
		return nil, errors.New("用户未登录")
	}

	var application model.DistributorApplication
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&application).Error
	if err != nil {
		return nil, errors.New("申请记录不存在")
	}

	if application.UserId != userId {
		return nil, errors.New("无权查看此申请")
	}

	resp = &types.DistributorApplicationResp{
		Id:        application.Id,
		UserId:    application.UserId,
		BrandId:   application.BrandId,
		Status:    application.Status,
		Reason:    application.Reason,
		CreatedAt: application.CreatedAt.Format(time.RFC3339),
	}

	if application.ReviewedBy != nil {
		resp.ReviewedBy = *application.ReviewedBy
	}
	if application.ReviewedAt != nil {
		resp.ReviewedAt = application.ReviewedAt.Format(time.RFC3339)
	}
	resp.ReviewNotes = application.ReviewNotes

	return resp, nil
}
