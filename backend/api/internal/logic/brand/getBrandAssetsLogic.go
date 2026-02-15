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

type GetBrandAssetsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandAssetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandAssetsLogic {
	return &GetBrandAssetsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandAssetsLogic) GetBrandAssets(req *types.GetBrandAssetsReq) (resp *types.BrandAssetListResp, err error) {
	if req.Id <= 0 {
		return nil, errors.New("品牌ID无效")
	}

	var assets []model.BrandAsset
	if err := l.svcCtx.DB.Where("brand_id = ?", req.Id).Order("id DESC").Find(&assets).Error; err != nil {
		l.Errorf("查询品牌素材失败: %v", err)
		return nil, errors.New("查询品牌素材失败")
	}

	items := make([]types.BrandAssetResp, 0, len(assets))
	for _, asset := range assets {
		items = append(items, types.BrandAssetResp{
			Id:          asset.ID,
			BrandId:     asset.BrandID,
			Name:        asset.Name,
			Type:        asset.Type,
			Category:    asset.Category,
			Tags:        asset.Tags,
			FileUrl:     asset.FileUrl,
			FileSize:    asset.FileSize,
			Description: asset.Description,
			CreatedAt:   asset.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   asset.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp = &types.BrandAssetListResp{
		Total:  int64(len(items)),
		Assets: items,
	}

	return resp, nil
}
