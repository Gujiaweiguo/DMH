package auth

import (
	"context"
	"fmt"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenResp, err error) {
	// 参数验证
	if req.Token == "" {
		return nil, fmt.Errorf("token不能为空")
	}

	// 使用认证中间件刷新token
	authMiddleware := middleware.NewAuthMiddleware(l.svcCtx.Config.Auth.AccessSecret)
	newToken, err := authMiddleware.RefreshToken(req.Token)
	if err != nil {
		l.Logger.Errorf("刷新token失败: %v", err)
		return nil, fmt.Errorf("token刷新失败: %v", err)
	}

	l.Logger.Info("token刷新成功")

	return &types.RefreshTokenResp{
		Token: newToken,
	}, nil
}