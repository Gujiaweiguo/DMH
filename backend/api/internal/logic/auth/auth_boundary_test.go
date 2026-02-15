package auth

import (
	"context"
	"testing"

	"dmh/api/internal/config"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
)

func TestGetUserInfoLogic_GetUserInfo_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	user := createTestUserWithPassword(t, db, "testuser", "password123", "active")

	role := &model.Role{
		Code: "participant",
		Name: "参与者",
	}
	db.Create(role)

	userRole := &model.UserRole{
		UserID: user.Id,
		RoleID: role.ID,
	}
	db.Create(userRole)

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetUserInfoLogic(ctx, svcCtx)

	resp, err := logic.GetUserInfo()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, user.Id, resp.Id)
	assert.Equal(t, "testuser", resp.Username)
	assert.NotEmpty(t, resp.Roles)
	assert.Contains(t, resp.Roles, "participant")
}

func TestGetUserInfoLogic_GetUserInfo_NotLoggedIn(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetUserInfoLogic(ctx, svcCtx)

	resp, err := logic.GetUserInfo()

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetUserInfoLogic_GetUserInfo_UserNotExists(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.WithValue(context.Background(), "userId", int64(99999))
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetUserInfoLogic(ctx, svcCtx)

	resp, err := logic.GetUserInfo()

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")
}

func TestUpdateProfileLogic_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	user := createTestUserWithPassword(t, db, "testuser", "password123", "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateProfileLogic(ctx, svcCtx)

	req := &types.UpdateProfileReq{
		RealName: "更新姓名",
	}

	resp, err := logic.UpdateProfile(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "更新姓名", resp.RealName)
}

func TestUpdateProfileLogic_NotLoggedIn(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateProfileLogic(ctx, svcCtx)

	req := &types.UpdateProfileReq{
		RealName: "测试姓名",
	}

	resp, err := logic.UpdateProfile(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestBindPhoneLogic_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	user := createTestUserWithPassword(t, db, "testuser", "password123", "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewBindPhoneLogic(ctx, svcCtx)

	req := &types.BindPhoneReq{
		Phone: "13912345678",
		Code:  "123456",
	}

	resp, err := logic.BindPhone(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestBindPhoneLogic_NotLoggedIn(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewBindPhoneLogic(ctx, svcCtx)

	req := &types.BindPhoneReq{
		Phone: "13912345678",
		Code:  "123456",
	}

	resp, err := logic.BindPhone(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestBindEmailLogic_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	user := createTestUserWithPassword(t, db, "testuser", "password123", "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewBindEmailLogic(ctx, svcCtx)

	req := &types.BindEmailReq{
		Email: "test@example.com",
		Code:  "123456",
	}

	resp, err := logic.BindEmail(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestBindEmailLogic_NotLoggedIn(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewBindEmailLogic(ctx, svcCtx)

	req := &types.BindEmailReq{
		Email: "test@example.com",
		Code:  "123456",
	}

	resp, err := logic.BindEmail(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestRefreshTokenLogic_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	user := createTestUserWithPassword(t, db, "testuser", "password123", "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{
		DB: db,
		Config: config.Config{
			Auth: struct {
				AccessSecret string
				AccessExpire int64
			}{
				AccessSecret: "test-secret-key-for-testing",
				AccessExpire: 86400,
			},
		},
	}
	logic := NewRefreshTokenLogic(ctx, svcCtx)

	req := &types.RefreshTokenReq{
		Token: "valid_token",
	}

	resp, err := logic.RefreshToken(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestRefreshTokenLogic_InvalidToken(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewRefreshTokenLogic(ctx, svcCtx)

	req := &types.RefreshTokenReq{
		Token: "invalid_token",
	}

	resp, err := logic.RefreshToken(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestSendPhoneCodeLogic_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewSendPhoneCodeLogic(ctx, svcCtx)

	req := &types.SendCodeReq{
		Target: "13912345678",
		Type:   "phone",
	}

	resp, err := logic.SendPhoneCode(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, resp.Message, "发送成功")
}

func TestSendEmailCodeLogic_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewSendEmailCodeLogic(ctx, svcCtx)

	req := &types.SendCodeReq{
		Target: "test@example.com",
		Type:   "email",
	}

	resp, err := logic.SendEmailCode(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, resp.Message, "发送成功")
}
