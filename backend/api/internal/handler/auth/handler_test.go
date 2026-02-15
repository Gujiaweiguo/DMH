package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Role{}, &model.UserRole{}, &model.UserBrand{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func createTestUserForHandler(t *testing.T, db *gorm.DB, username, password string) *model.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Phone:    "13800138000",
		Email:    username + "@test.com",
		Status:   "active",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func TestLoginHandler_Success(t *testing.T) {
	db := setupHandlerTestDB(t)
	createTestUserForHandler(t, db, "testuser", "password123")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := LoginHandler(svcCtx)

	reqBody := types.LoginReq{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var result types.LoginResp
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Token)
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	db := setupHandlerTestDB(t)
	createTestUserForHandler(t, db, "testuser", "password123")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := LoginHandler(svcCtx)

	reqBody := types.LoginReq{
		Username: "testuser",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	db := setupHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := LoginHandler(svcCtx)

	reqBody := types.LoginReq{
		Username: "nonexistent",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	handler := LoginHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestRegisterHandler_Success(t *testing.T) {
	db := setupHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := RegisterHandler(svcCtx)

	reqBody := types.RegisterReq{
		Username: "newuser",
		Password: "password123",
		Phone:    "13800138111",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var result types.LoginResp
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Token)
}

func TestRegisterHandler_DuplicateUsername(t *testing.T) {
	db := setupHandlerTestDB(t)
	createTestUserForHandler(t, db, "existinguser", "password123")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := RegisterHandler(svcCtx)

	reqBody := types.RegisterReq{
		Username: "existinguser",
		Password: "password123",
		Phone:    "13800138111",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	handler := RegisterHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestRefreshTokenHandler_InvalidToken(t *testing.T) {
	db := setupHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := RefreshTokenHandler(svcCtx)

	reqBody := types.RefreshTokenReq{
		Token: "invalid-token",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestRefreshTokenHandler_InvalidJSON(t *testing.T) {
	handler := RefreshTokenHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetUserInfoHandler_Success(t *testing.T) {
	db := setupHandlerTestDB(t)
	user := createTestUserForHandler(t, db, "testuser", "password123")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserInfoHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/userinfo", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, "userId", user.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetUserInfoHandler_NoUserId(t *testing.T) {
	db := setupHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserInfoHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/userinfo", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestChangePasswordHandler_Success(t *testing.T) {
	db := setupHandlerTestDB(t)
	user := createTestUserForHandler(t, db, "testuser", "oldpassword")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ChangePasswordHandler(svcCtx)

	reqBody := types.ChangePasswordReq{
		OldPassword: "oldpassword",
		NewPassword: "newpassword123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/auth/password", bytes.NewReader(body))
	ctx := req.Context()
	ctx = context.WithValue(ctx, "userId", user.Id)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestChangePasswordHandler_WrongOldPassword(t *testing.T) {
	db := setupHandlerTestDB(t)
	user := createTestUserForHandler(t, db, "testuser", "oldpassword")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ChangePasswordHandler(svcCtx)

	reqBody := types.ChangePasswordReq{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/auth/password", bytes.NewReader(body))
	ctx := req.Context()
	ctx = context.WithValue(ctx, "userId", user.Id)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestLogoutHandler_Success(t *testing.T) {
	db := setupHandlerTestDB(t)
	user := createTestUserForHandler(t, db, "testuser", "password123")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := LogoutHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, "userId", user.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestAuthHandlersConstruct(t *testing.T) {
	assert.NotNil(t, BindEmailHandler(nil))
	assert.NotNil(t, BindPhoneHandler(nil))
	assert.NotNil(t, SendEmailCodeHandler(nil))
	assert.NotNil(t, SendPhoneCodeHandler(nil))
	assert.NotNil(t, UpdateProfileHandler(nil))
}

func TestBindEmailHandler_ParseError(t *testing.T) {
	handler := BindEmailHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/bind-email", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestBindPhoneHandler_ParseError(t *testing.T) {
	handler := BindPhoneHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/bind-phone", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestSendEmailCodeHandler_ParseError(t *testing.T) {
	handler := SendEmailCodeHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/send-code", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestSendPhoneCodeHandler_ParseError(t *testing.T) {
	handler := SendPhoneCodeHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/send-code", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestUpdateProfileHandler_ParseError(t *testing.T) {
	handler := UpdateProfileHandler(nil)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/auth/profile", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}
