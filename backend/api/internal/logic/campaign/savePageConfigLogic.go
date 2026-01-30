// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SavePageConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSavePageConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SavePageConfigLogic {
	return &SavePageConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SavePageConfigLogic) SavePageConfig(req *types.PageConfigReq) (resp *types.PageConfigResp, err error) {
	// todo: add your logic here and delete this line

	return
}
