// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"errors"
	"time"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/golang-jwt/jwt/v4"
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
	if req.Token == "" {
		return nil, errors.New("token不能为空")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(l.svcCtx.Config.Auth.AccessSecret), nil
	})

	if err != nil || !token.Valid {
		l.Errorf("token验证失败: %v", err)
		return nil, errors.New("token无效或已过期")
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second)

	newClaims := &middleware.JWTClaims{
		UserID:   int64(claims["userId"].(float64)),
		Username: claims["username"].(string),
		Roles:    claims["roles"].([]string),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := newToken.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
	if err != nil {
		l.Errorf("生成新token失败: %v", err)
		return nil, errors.New("生成token失败")
	}

	return &types.RefreshTokenResp{
		Token: tokenString,
	}, nil
}
