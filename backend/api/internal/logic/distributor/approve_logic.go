package distributor

import (
	"context"
	"fmt"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ApproveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLogic {
	return &ApproveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ApproveApplication 审批分销商申请
func (l *ApproveLogic) ApproveApplication(brandId, applicationId int64, req *types.ApproveDistributorReq) (*types.DistributorApplicationResp, error) {
	// 获取当前用户信息（品牌管理员或平台管理员）
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 检查申请是否存在
	var application model.DistributorApplication
	if err := l.svcCtx.DB.Where("id = ? AND brand_id = ?", applicationId, brandId).First(&application).Error; err != nil {
		return nil, fmt.Errorf("申请不存在")
	}

	// 检查申请状态
	if application.Status != "pending" {
		return nil, fmt.Errorf("申请已处理，当前状态：%s", application.Status)
	}

	now := time.Now()
	reviewedAt := now

	// 根据操作类型处理
	if req.Action == "reject" {
		// 拒绝申请
		application.Status = "rejected"
		application.ReviewedBy = &userId
		application.ReviewedAt = &reviewedAt
		application.ReviewNotes = req.Reason

		if err := l.svcCtx.DB.Save(&application).Error; err != nil {
			l.Logger.Errorf("更新申请状态失败: %v", err)
			return nil, fmt.Errorf("处理申请失败")
		}

		return l.buildApplicationResp(application)
	}

	// 批准申请
	if req.Action != "approve" {
		return nil, fmt.Errorf("无效的操作类型")
	}

	// 验证级别
	if req.Level < 1 || req.Level > 3 {
		return nil, fmt.Errorf("级别必须在1-3之间")
	}

	// 开始事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查是否已存在分销商记录
	var existingDistributor model.Distributor
	if err := tx.Where("user_id = ? AND brand_id = ?", application.UserId, application.BrandId).
		First(&existingDistributor).Error; err == nil {
		// 已存在分销商记录，更新状态
		existingDistributor.Status = "active"
		existingDistributor.Level = req.Level
		existingDistributor.ApprovedBy = &userId
		existingDistributor.ApprovedAt = &reviewedAt
		if err := tx.Save(&existingDistributor).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("更新分销商记录失败")
		}
	} else if err == gorm.ErrRecordNotFound {
		// 创建新的分销商记录
		distributor := model.Distributor{
			UserId:     application.UserId,
			BrandId:    application.BrandId,
			Level:      req.Level,
			Status:     "active",
			ApprovedBy: &userId,
			ApprovedAt: &reviewedAt,
		}

		if err := tx.Create(&distributor).Error; err != nil {
			tx.Rollback()
			l.Logger.Errorf("创建分销商记录失败: %v", err)
			return nil, fmt.Errorf("创建分销商记录失败")
		}
	} else {
		tx.Rollback()
		return nil, fmt.Errorf("查询分销商记录失败")
	}

	// 更新申请状态
	application.Status = "approved"
	application.ReviewedBy = &userId
	application.ReviewedAt = &reviewedAt
	application.ReviewNotes = req.Reason

	if err := tx.Save(&application).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新申请状态失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败")
	}

	return l.buildApplicationResp(application)
}

// GetPendingApplications 获取待审批申请列表（品牌管理员）
func (l *ApproveLogic) GetPendingApplications(brandId int64, page, pageSize int64, status string) (*types.DistributorApplicationListResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 检查用户是否有权限查看此品牌的申请
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 检查权限：平台管理员可以查看所有，品牌管理员只能查看自己的品牌
	if user.Role != "platform_admin" {
		// 检查是否是品牌管理员
		var count int64
		l.svcCtx.DB.Model(&model.UserBrand{}).
			Where("user_id = ? AND brand_id = ?", userId, brandId).
			Count(&count)
		if count == 0 {
			return nil, fmt.Errorf("无权限查看该品牌的申请")
		}
	}

	query := l.svcCtx.DB.Model(&model.DistributorApplication{}).Where("brand_id = ?", brandId)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var applications []model.DistributorApplication
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Limit(int(pageSize)).Offset(int(offset)).Find(&applications)

	resp := &types.DistributorApplicationListResp{
		Total:        total,
		Applications: make([]types.DistributorApplicationResp, 0),
	}

	for _, app := range applications {
		var brand model.Brand
		l.svcCtx.DB.Where("id = ?", app.BrandId).First(&brand)

		var applicant model.User
		l.svcCtx.DB.Where("id = ?", app.UserId).First(&applicant)

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
			Username:    applicant.Username,
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

// GetApplicationDetail 获取申请详情（管理员）
func (l *ApproveLogic) GetApplicationDetail(brandId, applicationId int64) (*types.DistributorApplicationResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 检查权限
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 平台管理员可以查看所有，品牌管理员只能查看自己的品牌
	if user.Role != "platform_admin" {
		var count int64
		l.svcCtx.DB.Model(&model.UserBrand{}).
			Where("user_id = ? AND brand_id = ?", userId, brandId).
			Count(&count)
		if count == 0 {
			return nil, fmt.Errorf("无权限查看该品牌的申请")
		}
	}

	var application model.DistributorApplication
	if err := l.svcCtx.DB.Where("id = ? AND brand_id = ?", applicationId, brandId).First(&application).Error; err != nil {
		return nil, fmt.Errorf("申请不存在")
	}

	return l.buildApplicationResp(application)
}

// buildApplicationResp 构建申请响应
func (l *ApproveLogic) buildApplicationResp(application model.DistributorApplication) (*types.DistributorApplicationResp, error) {
	var brand model.Brand
	l.svcCtx.DB.Where("id = ?", application.BrandId).First(&brand)

	var applicant model.User
	l.svcCtx.DB.Where("id = ?", application.UserId).First(&applicant)

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
		Username:    applicant.Username,
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
