// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"fmt"
	"time"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	var user model.User

	if err := l.svcCtx.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		l.Errorf("用户不存在或密码错误: %v", err)
		return nil, fmt.Errorf("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		l.Errorf("密码验证失败: %v", err)
		return nil, fmt.Errorf("用户名或密码错误")
	}

	if user.Status != "active" {
		l.Errorf("用户已被禁用: %s", user.Username)
		return nil, fmt.Errorf("账号已被禁用")
	}

	// 查询用户角色
	var roles []model.Role
	l.svcCtx.DB.Table("roles").
		Select("roles.*").
		Joins("INNER JOIN user_roles ur ON roles.id = ur.role_id").
		Where("ur.user_id = ?", user.Id).
		Find(&roles)

	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	// 查询用户关联的品牌
	var userBrands []model.UserBrand
	l.svcCtx.DB.Where("user_id = ?", user.Id).Find(&userBrands)

	brandIds := make([]int64, 0, len(userBrands))
	for _, ub := range userBrands {
		brandIds = append(brandIds, ub.BrandId)
	}

	// 生成JWT token
	token, err := l.generateToken(user.Id, user.Username, roleCodes)
	if err != nil {
		l.Errorf("生成token失败: %v", err)
		return nil, err
	}

	resp = &types.LoginResp{
		Token:    token,
		UserId:   user.Id,
		Username: user.Username,
		Phone:    user.Phone,
		RealName: user.RealName,
		Roles:    roleCodes,
		BrandIds: brandIds,
	}
	return resp, nil
}

func (l *LoginLogic) generateToken(userId int64, username string, roles []string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second)

	claims := &middleware.JWTClaims{
		UserID:   userId,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
	if err != nil {
		return "", fmt.Errorf("生成JWT token失败: %w", err)
	}

	return tokenString, nil
}
