// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"context"
	"errors"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBrandAssetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBrandAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBrandAssetLogic {
	return &DeleteBrandAssetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBrandAssetLogic) DeleteBrandAsset(req *types.DeleteBrandAssetReq) (resp *types.CommonResp, err error) {
	if req.BrandId <= 0 || req.Id <= 0 {
		return nil, errors.New("参数无效")
	}

	var asset model.BrandAsset
	if err := l.svcCtx.DB.Where("id = ? AND brand_id = ?", req.Id, req.BrandId).First(&asset).Error; err != nil {
		return nil, errors.New("素材不存在")
	}

	if err := l.svcCtx.DB.Delete(&asset).Error; err != nil {
		l.Errorf("删除品牌素材失败: %v", err)
		return nil, errors.New("删除品牌素材失败")
	}

	resp = &types.CommonResp{
		Message: "删除成功",
	}

	return resp, nil
}
