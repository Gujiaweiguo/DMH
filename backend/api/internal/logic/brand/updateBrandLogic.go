// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"context"
	"errors"
	"strings"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBrandLogic {
	return &UpdateBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBrandLogic) UpdateBrand(req *types.UpdateBrandReq) (resp *types.BrandResp, err error) {
	if req.Id <= 0 {
		return nil, errors.New("品牌ID无效")
	}

	var brand model.Brand
	if err := l.svcCtx.DB.First(&brand, req.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("品牌不存在")
		}
		l.Errorf("查询品牌失败: %v", err)
		return nil, errors.New("查询品牌失败")
	}

	if name := strings.TrimSpace(req.Name); name != "" {
		brand.Name = name
	}
	if logo := strings.TrimSpace(req.Logo); logo != "" {
		brand.Logo = logo
	}
	if desc := strings.TrimSpace(req.Description); desc != "" {
		brand.Description = desc
	}
	if status := strings.TrimSpace(req.Status); status != "" {
		brand.Status = status
	}

	if err := l.svcCtx.DB.Save(&brand).Error; err != nil {
		l.Errorf("更新品牌失败: %v", err)
		return nil, errors.New("更新品牌失败")
	}

	resp = &types.BrandResp{
		Id:          brand.Id,
		Name:        brand.Name,
		Logo:        brand.Logo,
		Description: brand.Description,
		Status:      brand.Status,
		CreatedAt:   brand.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   brand.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}
