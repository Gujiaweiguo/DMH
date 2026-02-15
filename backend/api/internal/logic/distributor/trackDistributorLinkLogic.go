package distributor

import (
	"context"
	"errors"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type TrackDistributorLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTrackDistributorLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrackDistributorLinkLogic {
	return &TrackDistributorLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TrackDistributorLinkLogic) TrackDistributorLink(req *types.TrackDistributorLinkReq) (resp *types.CommonResp, err error) {
	if req.Code == "" {
		return nil, errors.New("分销商代码不能为空")
	}

	// Find the distributor by code (ID)
	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("id = ? AND status = ?", req.Code, "active").First(&distributor).Error; err != nil {
		return nil, errors.New("无效的分销商")
	}

	// Record link tracking (can be extended to store tracking info)
	l.Logger.Infof("分销商链接追踪: code=%s, distributor_id=%d", req.Code, distributor.Id)

	return &types.CommonResp{
		Message: "追踪成功",
	}, nil
}
