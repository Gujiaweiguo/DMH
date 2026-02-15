package admin

import (
	"context"
	"errors"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.AdminCreateUserReq) (resp *types.UserInfoResp, err error) {
	if req.Username == "" {
		return nil, errors.New("用户名不能为空")
	}

	if req.Password == "" {
		return nil, errors.New("密码不能为空")
	}

	if req.Phone == "" {
		return nil, errors.New("手机号不能为空")
	}

	var existingUser model.User
	if err := l.svcCtx.DB.Where("username = ? OR phone = ?", req.Username, req.Phone).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名或手机号已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Email:    req.Email,
		RealName: req.RealName,
		Status:   "active",
	}

	if err := l.svcCtx.DB.Create(user).Error; err != nil {
		l.Errorf("创建用户失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	if len(req.BrandIds) > 0 {
		for _, brandId := range req.BrandIds {
			userRole := &model.UserRole{
				UserID: user.Id,
				RoleID: brandId,
			}
			if err := l.svcCtx.DB.Create(userRole).Error; err != nil {
				l.Errorf("创建用户角色关系失败: %v", err)
			}
		}
	}

	resp = &types.UserInfoResp{
		Id:       user.Id,
		Username: user.Username,
		Phone:    user.Phone,
		RealName: user.RealName,
		Email:    user.Email,
		Status:   user.Status,
	}

	return resp, nil
}
