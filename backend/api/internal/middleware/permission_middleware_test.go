package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PermissionMiddlewareTestSuite struct {
	suite.Suite
	db         *gorm.DB
	middleware *PermissionMiddleware
}

func (suite *PermissionMiddlewareTestSuite) SetupSuite() {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// 自动迁移表结构
	err = db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
		&model.UserBrand{},
		&model.Brand{},
	)
	suite.Require().NoError(err)

	suite.db = db
	suite.middleware = NewPermissionMiddleware(db)
}

func (suite *PermissionMiddlewareTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *PermissionMiddlewareTestSuite) SetupTest() {
	// 清理测试数据
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM roles")
	suite.db.Exec("DELETE FROM permissions")
	suite.db.Exec("DELETE FROM user_roles")
	suite.db.Exec("DELETE FROM role_permissions")
	suite.db.Exec("DELETE FROM user_brands")
	suite.db.Exec("DELETE FROM brands")

	// 创建测试数据
	suite.createTestData()
}

func (suite *PermissionMiddlewareTestSuite) createTestData() {
	// 创建角色
	roles := []model.Role{
		{ID: 1, Name: "平台管理员", Code: "platform_admin"},
		{ID: 2, Name: "参与者", Code: "participant"},
	}
	for _, role := range roles {
		suite.db.Create(&role)
	}

	// 创建权限
	permissions := []model.Permission{
		{ID: 1, Name: "用户管理", Code: "user:manage", Resource: "user", Action: "manage"},
		{ID: 2, Name: "品牌管理", Code: "brand:manage", Resource: "brand", Action: "manage"},
		{ID: 3, Name: "活动管理", Code: "campaign:manage", Resource: "campaign", Action: "manage"},
		{ID: 4, Name: "订单查看", Code: "order:read", Resource: "order", Action: "read"},
	}
	for _, permission := range permissions {
		suite.db.Create(&permission)
	}

	// 创建角色权限关联
	rolePermissions := []model.RolePermission{
		{RoleID: 1, PermissionID: 1}, // platform_admin -> user:manage
		{RoleID: 1, PermissionID: 2}, // platform_admin -> brand:manage
		{RoleID: 1, PermissionID: 3}, // platform_admin -> campaign:manage
		{RoleID: 1, PermissionID: 4}, // platform_admin -> order:read
		{RoleID: 2, PermissionID: 4}, // participant -> order:read
	}
	for _, rp := range rolePermissions {
		suite.db.Create(&rp)
	}

	// 创建用户
	users := []model.User{
		{Id: 1, Username: "admin", Role: "platform_admin", Status: "active"},
		{Id: 2, Username: "brand_manager", Role: "participant", Status: "active"},
		{Id: 3, Username: "user", Role: "participant", Status: "active"},
	}
	for _, user := range users {
		suite.db.Create(&user)
	}

	// 创建用户角色关联
	userRoles := []model.UserRole{
		{UserID: 1, RoleID: 1}, // admin -> platform_admin
		{UserID: 2, RoleID: 2}, // brand_manager -> participant
		{UserID: 3, RoleID: 2}, // user -> participant
	}
	for _, ur := range userRoles {
		suite.db.Create(&ur)
	}

	// 创建品牌
	brands := []model.Brand{
		{Id: 1, Name: "品牌A", Status: "active"},
		{Id: 2, Name: "品牌B", Status: "active"},
	}
	for _, brand := range brands {
		suite.db.Create(&brand)
	}

	// 创建用户品牌关联
	userBrands := []model.UserBrand{
		{UserId: 2, BrandId: 1}, // brand_manager -> 品牌A
	}
	for _, ub := range userBrands {
		suite.db.Create(&ub)
	}
}

func (suite *PermissionMiddlewareTestSuite) TestCheckPermission() {
	// 测试平台管理员权限
	hasPermission, err := suite.middleware.CheckPermission(1, "user:manage", "", 0)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	// 测试品牌管理员权限
	hasPermission, err = suite.middleware.CheckPermission(2, "brand:manage", "", 0)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	// 测试品牌管理员无用户管理权限
	hasPermission, err = suite.middleware.CheckPermission(2, "user:manage", "", 0)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)

	// 测试参与者权限
	hasPermission, err = suite.middleware.CheckPermission(3, "order:read", "", 0)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	// 测试参与者无品牌管理权限
	hasPermission, err = suite.middleware.CheckPermission(3, "brand:manage", "", 0)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)
}

func (suite *PermissionMiddlewareTestSuite) TestCheckBrandPermission() {
	// 测试平台管理员可以访问任何品牌
	hasPermission, err := suite.middleware.CheckPermission(1, "brand:manage", "brand", 1)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	hasPermission, err = suite.middleware.CheckPermission(1, "brand:manage", "brand", 2)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	// 测试品牌管理员只能访问分配的品牌
	hasPermission, err = suite.middleware.CheckPermission(2, "brand:manage", "brand", 1)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	hasPermission, err = suite.middleware.CheckPermission(2, "brand:manage", "brand", 2)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)

	// 测试参与者无品牌访问权限
	hasPermission, err = suite.middleware.CheckPermission(3, "brand:manage", "brand", 1)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)
}

func (suite *PermissionMiddlewareTestSuite) TestGetUserPermissions() {
	// 测试获取平台管理员权限
	permissions, err := suite.middleware.GetUserPermissions(1)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), permissions, "user:manage")
	assert.Contains(suite.T(), permissions, "brand:manage")
	assert.Contains(suite.T(), permissions, "campaign:manage")
	assert.Contains(suite.T(), permissions, "order:read")

	// 测试获取品牌管理员权限
	permissions, err = suite.middleware.GetUserPermissions(2)
	assert.NoError(suite.T(), err)
	assert.NotContains(suite.T(), permissions, "user:manage")
	assert.Contains(suite.T(), permissions, "brand:manage")
	assert.Contains(suite.T(), permissions, "campaign:manage")
	assert.Contains(suite.T(), permissions, "order:read")

	// 测试获取参与者权限
	permissions, err = suite.middleware.GetUserPermissions(3)
	assert.NoError(suite.T(), err)
	assert.NotContains(suite.T(), permissions, "user:manage")
	assert.NotContains(suite.T(), permissions, "brand:manage")
	assert.NotContains(suite.T(), permissions, "campaign:manage")
	assert.Contains(suite.T(), permissions, "order:read")
}

func (suite *PermissionMiddlewareTestSuite) TestGetUserBrandIDs() {
	// 测试获取平台管理员品牌ID（应该返回所有品牌）
	brandIDs, err := suite.middleware.GetUserBrandIDs(1)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), brandIDs, int64(1))
	assert.Contains(suite.T(), brandIDs, int64(2))

	// 测试获取品牌管理员品牌ID
	brandIDs, err = suite.middleware.GetUserBrandIDs(2)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), brandIDs, int64(1))
	assert.NotContains(suite.T(), brandIDs, int64(2))

	// 测试获取参与者品牌ID（应该为空）
	brandIDs, err = suite.middleware.GetUserBrandIDs(3)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), brandIDs)
}

func (suite *PermissionMiddlewareTestSuite) TestRequirePermissionMiddleware() {
	// 创建测试处理器
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// 测试有权限的请求
	middleware := suite.middleware.RequirePermission("order:read")
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userId", int64(3)))
	req = req.WithContext(context.WithValue(req.Context(), "username", "user"))
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "success", w.Body.String())

	// 测试无权限的请求
	middleware = suite.middleware.RequirePermission("user:manage")
	wrappedHandler = middleware(handler)

	req = httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userId", int64(3)))
	req = req.WithContext(context.WithValue(req.Context(), "username", "user"))
	w = httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func (suite *PermissionMiddlewareTestSuite) TestRequireBrandPermissionMiddleware() {
	// 创建测试处理器
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// 测试品牌管理员访问分配的品牌
	middleware := suite.middleware.RequireBrandPermission("brand:manage", "brandId")
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test?brandId=1", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userId", int64(2)))
	req = req.WithContext(context.WithValue(req.Context(), "username", "brand_manager"))
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// 测试品牌管理员访问未分配的品牌
	req = httptest.NewRequest("GET", "/test?brandId=2", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userId", int64(2)))
	req = req.WithContext(context.WithValue(req.Context(), "username", "brand_manager"))
	w = httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func (suite *PermissionMiddlewareTestSuite) TestCachePermissions() {
	// 第一次获取权限（应该从数据库查询）
	permissions1, err := suite.middleware.GetUserPermissions(1)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), permissions1)

	// 第二次获取权限（应该从缓存获取）
	permissions2, err := suite.middleware.GetUserPermissions(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), permissions1, permissions2)

	// 清除缓存
	suite.middleware.ClearUserPermissionCache(1)

	// 再次获取权限（应该重新从数据库查询）
	permissions3, err := suite.middleware.GetUserPermissions(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), permissions1, permissions3)
}

func TestPermissionMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionMiddlewareTestSuite))
}
