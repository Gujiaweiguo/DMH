package admin

import (
	"bytes"
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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAdminHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Role{}, &model.UserRole{}, &model.UserBrand{}, &model.Brand{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestAdminHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetUsersHandler(nil))
	assert.NotNil(t, GetUserHandler(nil))
	assert.NotNil(t, CreateUserHandler(nil))
	assert.NotNil(t, UpdateUserHandler(nil))
	assert.NotNil(t, DeleteUserHandler(nil))
	assert.NotNil(t, ResetUserPasswordHandler(nil))
	assert.NotNil(t, ManageBrandAdminRelationHandler(nil))
}

func TestGetUsersHandler_Success(t *testing.T) {
	db := setupAdminHandlerTestDB(t)

	user := &model.User{Username: "admin", Password: "pass", Phone: "13800138000", Status: "active", Role: "platform_admin"}
	db.Create(user)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUsersHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCreateUserHandler_Success(t *testing.T) {
	db := setupAdminHandlerTestDB(t)

	role := &model.Role{Name: "Participant", Code: "participant"}
	db.Create(role)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := CreateUserHandler(svcCtx)

	reqBody := types.AdminCreateUserReq{
		Username: "newuser",
		Password: "password123",
		Phone:    "13800138111",
		Role:     "participant",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCreateUserHandler_ParseError(t *testing.T) {
	handler := CreateUserHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetUserHandler_ParseError(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestUpdateUserHandler_ParseError(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateUserHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/users/1", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestDeleteUserHandler_ParseError(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := DeleteUserHandler(svcCtx)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/users/invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestResetUserPasswordHandler_ParseError(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := ResetUserPasswordHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/1/reset-password", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestManageBrandAdminRelationHandler_ParseError(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := ManageBrandAdminRelationHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/brand-relations", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}
