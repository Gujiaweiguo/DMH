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
)

type CreateBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBrandLogic {
	return &CreateBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBrandLogic) CreateBrand(req *types.CreateBrandReq) (resp *types.BrandResp, err error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("品牌名称不能为空")
	}

	brand := &model.Brand{
		Name:        name,
		Logo:        strings.TrimSpace(req.Logo),
		Description: strings.TrimSpace(req.Description),
		Status:      "active",
	}

	if err := l.svcCtx.DB.Create(brand).Error; err != nil {
		l.Errorf("创建品牌失败: %v", err)
		return nil, errors.New("创建品牌失败")
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
