package role

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dmh/api/internal/svc"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRoleHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Role{}, &model.Permission{}, &model.RolePermission{}, &model.User{}, &model.UserRole{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestRoleHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetRolesHandler(nil))
	assert.NotNil(t, GetPermissionsHandler(nil))
	assert.NotNil(t, ConfigRolePermissionsHandler(nil))
	assert.NotNil(t, GetUserPermissionsHandler(nil))
}

func TestGetRolesHandler_Success(t *testing.T) {
	db := setupRoleHandlerTestDB(t)

	role := &model.Role{Name: "Admin", Code: "admin"}
	db.Create(role)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetRolesHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/roles", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetPermissionsHandler_Success(t *testing.T) {
	db := setupRoleHandlerTestDB(t)

	permission := &model.Permission{Name: "View Users", Code: "users:view", Resource: "users", Action: "read"}
	db.Create(permission)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPermissionsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/permissions", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestConfigRolePermissionsHandler_ParseError(t *testing.T) {
	handler := ConfigRolePermissionsHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/roles/permissions", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetUserPermissionsHandler_ParseError(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserPermissionsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/invalid/permissions", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}
