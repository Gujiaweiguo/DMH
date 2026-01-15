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

type CreateExportRequestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateExportRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateExportRequestLogic {
	return &CreateExportRequestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateExportRequestLogic) CreateExportRequest(req *types.ExportRequestReq) (*types.ExportRequestResp, error) {
	userId, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}

	// 检查用户是否有该品牌的权限
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 品牌管理员需要检查品牌权限
	if user.Role != "platform_admin" {
		var count int64
		l.svcCtx.DB.Model(&model.UserBrand{}).
			Where("user_id = ? AND brand_id = ?", userId, req.BrandId).
			Count(&count)

		if count == 0 {
			return nil, fmt.Errorf("无权为该品牌申请导出")
		}
	}

	// 创建导出请求
	exportReq := model.ExportRequest{
		BrandID:     req.BrandId,
		RequestedBy: userId,
		Reason:      req.Reason,
		Filters:     req.Filters,
		Status:      "pending",
	}

	if err := l.svcCtx.DB.Create(&exportReq).Error; err != nil {
		return nil, fmt.Errorf("创建导出请求失败: %w", err)
	}

	// 查询品牌名称
	var brand model.Brand
	l.svcCtx.DB.Where("id = ?", req.BrandId).First(&brand)

	return &types.ExportRequestResp{
		Id:              exportReq.ID,
		BrandId:         exportReq.BrandID,
		BrandName:       brand.Name,
		RequestedBy:     exportReq.RequestedBy,
		RequestedByName: user.Username,
		Reason:          exportReq.Reason,
		Filters:         exportReq.Filters,
		Status:          exportReq.Status,
		RecordCount:     exportReq.RecordCount,
		CreatedAt:       exportReq.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
