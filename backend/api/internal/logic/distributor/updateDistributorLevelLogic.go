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

type UpdateDistributorLevelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDistributorLevelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDistributorLevelLogic {
	return &UpdateDistributorLevelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDistributorLevelLogic) UpdateDistributorLevel(req *types.UpdateDistributorLevelReq) (resp *types.CommonResp, err error) {
	distributorId := l.ctx.Value("distributorId").(int64)

	distributor := &model.Distributor{}
	if err := l.svcCtx.DB.First(distributor, distributorId).Error; err != nil {
		return nil, err
	}

	distributor.Level = req.Level
	if err := l.svcCtx.DB.Save(distributor).Error; err != nil {
		return nil, err
	}

	resp = &types.CommonResp{
		Message: "Distributor level updated successfully",
	}

	return resp, nil
}
