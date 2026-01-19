package distributor

import (
	"context"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyLogic {
	return &ApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Apply 分销商申请
func (l *ApplyLogic) Apply(req *types.DistributorApplyReq) (*types.DistributorApplicationResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 检查品牌是否存在
	var brand model.Brand
	if err := l.svcCtx.DB.Where("id = ? AND status = ?", req.BrandId, "active").First(&brand).Error; err != nil {
		return nil, fmt.Errorf("品牌不存在或已停用")
	}

	// 检查是否已有待审批或已批准的申请
	var existingApp model.DistributorApplication
	if err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ? AND status IN (?)",
		userId, req.BrandId, []string{"pending", "approved"}).
		First(&existingApp).Error; err == nil {
		// 已有申请记录
		if existingApp.Status == "pending" {
			return nil, fmt.Errorf("您已提交过申请，请等待审核")
		}
		if existingApp.Status == "approved" {
			return nil, fmt.Errorf("您已经是该品牌的分销商")
		}
	}

	// 创建申请记录
	application := model.DistributorApplication{
		UserId:  userId,
		BrandId: req.BrandId,
		Status:  "pending",
		Reason:  req.Reason,
	}

	if err := l.svcCtx.DB.Create(&application).Error; err != nil {
		l.Logger.Errorf("创建分销商申请失败: %v", err)
		return nil, fmt.Errorf("提交申请失败")
	}

	// 查询用户信息
	var user model.User
	l.svcCtx.DB.Where("id = ?", userId).First(&user)

	return &types.DistributorApplicationResp{
		Id:        application.Id,
		UserId:    application.UserId,
		Username:  user.Username,
		BrandId:   application.BrandId,
		BrandName: brand.Name,
		Status:    application.Status,
		Reason:    application.Reason,
		CreatedAt: application.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetApplications 获取申请列表
func (l *ApplyLogic) GetApplications(status string) (*types.DistributorApplicationListResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	query := l.svcCtx.DB.Model(&model.DistributorApplication{}).Where("user_id = ?", userId)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var applications []model.DistributorApplication
	var total int64

	query.Count(&total)
	query.Order("created_at DESC").Find(&applications)

	resp := &types.DistributorApplicationListResp{
		Total:        total,
		Applications: make([]types.DistributorApplicationResp, 0),
	}

	for _, app := range applications {
		var brand model.Brand
		l.svcCtx.DB.Where("id = ?", app.BrandId).First(&brand)

		var reviewerName string
		if app.ReviewedBy != nil {
			var reviewer model.User
			l.svcCtx.DB.Where("id = ?", *app.ReviewedBy).First(&reviewer)
			reviewerName = reviewer.RealName
			if reviewerName == "" {
				reviewerName = reviewer.Username
			}
		}

		applicationResp := types.DistributorApplicationResp{
			Id:          app.Id,
			UserId:      app.UserId,
			BrandId:     app.BrandId,
			BrandName:   brand.Name,
			Status:      app.Status,
			Reason:      app.Reason,
			ReviewNotes: app.ReviewNotes,
			CreatedAt:   app.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if app.ReviewedBy != nil {
			applicationResp.ReviewedBy = *app.ReviewedBy
			applicationResp.Reviewer = reviewerName
		}

		if app.ReviewedAt != nil {
			applicationResp.ReviewedAt = app.ReviewedAt.Format("2006-01-02 15:04:05")
		}

		resp.Applications = append(resp.Applications, applicationResp)
	}

	return resp, nil
}

// GetApplication 获取申请详情
func (l *ApplyLogic) GetApplication(id int64) (*types.DistributorApplicationResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	var application model.DistributorApplication
	if err := l.svcCtx.DB.Where("id = ? AND user_id = ?", id, userId).First(&application).Error; err != nil {
		return nil, fmt.Errorf("申请不存在")
	}

	var brand model.Brand
	l.svcCtx.DB.Where("id = ?", application.BrandId).First(&brand)

	var reviewerName string
	if application.ReviewedBy != nil {
		var reviewer model.User
		l.svcCtx.DB.Where("id = ?", *application.ReviewedBy).First(&reviewer)
		reviewerName = reviewer.RealName
		if reviewerName == "" {
			reviewerName = reviewer.Username
		}
	}

	resp := &types.DistributorApplicationResp{
		Id:          application.Id,
		UserId:      application.UserId,
		BrandId:     application.BrandId,
		BrandName:   brand.Name,
		Status:      application.Status,
		Reason:      application.Reason,
		ReviewNotes: application.ReviewNotes,
		CreatedAt:   application.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if application.ReviewedBy != nil {
		resp.ReviewedBy = *application.ReviewedBy
		resp.Reviewer = reviewerName
	}

	if application.ReviewedAt != nil {
		resp.ReviewedAt = application.ReviewedAt.Format("2006-01-02 15:04:05")
	}

	return resp, nil
}

// CheckDistributorStatus 检查用户是否是分销商
func (l *ApplyLogic) CheckDistributorStatus(brandId int64) (bool, *model.Distributor, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return false, nil, fmt.Errorf("获取用户信息失败")
	}

	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ? AND status = ?", userId, brandId, "active").
		First(&distributor).Error; err != nil {
		return false, nil, nil
	}

	return true, &distributor, nil
}
