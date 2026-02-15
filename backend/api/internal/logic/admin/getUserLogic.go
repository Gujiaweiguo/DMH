// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"errors"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser() (resp *types.UserInfoResp, err error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, errors.New("无法从context中获取用户ID")
	}

	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	var roles []model.Role
	err = l.svcCtx.DB.Table("user_roles").
		Select("roles.*").
		Joins("LEFT JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userId).
		Find(&roles).Error
	if err != nil {
		l.Errorf("查询用户角色失败: %v", err)
		return nil, errors.New("查询用户角色失败")
	}

	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	resp = &types.UserInfoResp{
		Id:        user.Id,
		Username:  user.Username,
		Phone:     user.Phone,
		RealName:  user.RealName,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Roles:     roleCodes,
		Status:    user.Status,
		BrandIds:  []int64{},
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}
