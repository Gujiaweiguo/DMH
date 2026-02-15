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

type DistributorApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDistributorApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DistributorApplyLogic {
	return &DistributorApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DistributorApplyLogic) DistributorApply(req *types.DistributorApplyReq) (resp *types.DistributorApplicationResp, err error) {
	if req.BrandId <= 0 {
		return nil, errors.New("品牌ID无效")
	}

	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId <= 0 {
		return nil, errors.New("用户未登录")
	}

	var existingApp model.DistributorApplication
	err = l.svcCtx.DB.Where("user_id = ? AND brand_id = ? AND status = ?", userId, req.BrandId, "pending").
		First(&existingApp).Error
	if err == nil {
		return nil, errors.New("您已提交申请，请勿重复申请")
	}

	var existingDistributor model.Distributor
	err = l.svcCtx.DB.Where("user_id = ? AND brand_id = ? AND status = ?", userId, req.BrandId, "active").
		First(&existingDistributor).Error
	if err == nil {
		return nil, errors.New("您已经是该品牌的分销商")
	}

	application := model.DistributorApplication{
		UserId:    userId,
		BrandId:   req.BrandId,
		Status:    "pending",
		Reason:    req.Reason,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := l.svcCtx.DB.Create(&application).Error; err != nil {
		l.Logger.Errorf("创建分销商申请失败: %v", err)
		return nil, errors.New("申请提交失败")
	}

	return &types.DistributorApplicationResp{
		Id:        application.Id,
		UserId:    application.UserId,
		BrandId:   application.BrandId,
		Status:    application.Status,
		Reason:    application.Reason,
		CreatedAt: application.CreatedAt.Format(time.RFC3339),
	}, nil
}
