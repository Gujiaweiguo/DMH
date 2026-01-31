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

type UpdateDistributorStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDistributorStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDistributorStatusLogic {
	return &UpdateDistributorStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDistributorStatusLogic) UpdateDistributorStatus(req *types.UpdateDistributorStatusReq) (resp *types.CommonResp, err error) {
	distributorId := l.ctx.Value("distributorId").(int64)

	distributor := &model.Distributor{}
	if err := l.svcCtx.DB.First(distributor, distributorId).Error; err != nil {
		return nil, err
	}

	distributor.Status = req.Status
	if err := l.svcCtx.DB.Save(distributor).Error; err != nil {
		return nil, err
	}

	resp = &types.CommonResp{
		Message: "Distributor status updated successfully",
	}

	return resp, nil
}
