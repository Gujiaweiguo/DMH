// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPageConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPageConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPageConfigLogic {
	return &GetPageConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPageConfigLogic) GetPageConfig() (resp *types.PageConfigResp, err error) {
	// todo: add your logic here and delete this line

	return
}
