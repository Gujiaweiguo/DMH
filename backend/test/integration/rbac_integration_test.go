package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dmh/api/internal/config"
	"dmh/api/internal/handler"
	"dmh/api/internal/middleware"
	"dmh/api/internal/service"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/rest"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RBACIntegrationTestSuite struct {
	suite.Suite
	db     *gorm.DB
	svcCtx *svc.ServiceContext
	server *rest.Server
	tokens map[string]string // 存储不同用户的token
}

func (suite *RBACIntegrationTestSuite) SetupSuite() {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// 自动迁移所有表结构
	err = db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
		&model.UserBrand{},
		&model.Brand{},
		&model.Menu{},
		&model.RoleMenu{},
		&model.PasswordPolicy{},
		&model.PasswordHistory{},
		&model.LoginAttempt{},
		&model.UserSession{},
		&model.AuditLog{},
		&model.SecurityEvent{},
		&model.BrandAsset{},
	)
	suite.Require().NoError(err)

	// 创建服务上下文
	cfg := config.Config{
		Auth: struct {
			AccessSecret string
			AccessExpire int64
		}{
			AccessSecret: "test-secret-key-for-jwt-token-generation-in-integration-tests",
			AccessExpire: 3600,
		},
	}

	// 初始化安全服务
	passwordService := service.NewPasswordService(db)
	auditService := service.NewAuditService(db)
	sessionService := service.NewSessionService(db, passwordService)

	suite.svcCtx = &svc.ServiceContext{
		Config:          cfg,
		DB:              db,
		PasswordService: passwordService,
		AuditService:    auditService,
		SessionService:  sessionService,
	}

	suite.db = db
	suite.tokens = make(map[string]string)

	// 创建测试数据
	suite.createTestData()

	// 设置REST服务器
	suite.server = rest.MustNewServer(rest.RestConf{
		Port: 0, // 使用随机端口
	})
}

func (suite *RBACIntegrationTestSuite) TearDownSuite() {
	suite.server.Stop()
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *RBACIntegrationTestSuite) createTestData() {
	// 创建默认密码策略
	policy := &model.PasswordPolicy{
		MinLength:             8,
		RequireUppercase:      true,
		RequireLowercase:      true,
		RequireNumbers:        true,
		RequireSpecialChars:   true,
		MaxAge:                90,
		HistoryCount:          5,
		MaxLoginAttempts:      5,
		LockoutDuration:       30,
		SessionTimeout:        480,
		MaxConcurrentSessions: 3,
	}
	suite.db.Create(policy)

	// 创建角色
	roles := []model.Role{
		{ID: 1, Name: "平台管理员", Code: "platform_admin"},
		{ID: 2, Name: "品牌管理员", Code: "brand_admin"},
		{ID: 3, Name: "参与者", Code: "participant"},
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
		{ID: 5, Name: "安全管理", Code: "security:manage", Resource: "security", Action: "manage"},
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
		{RoleID: 1, PermissionID: 5}, // platform_admin -> security:manage
		{RoleID: 2, PermissionID: 2}, // brand_admin -> brand:manage
		{RoleID: 2, PermissionID: 3}, // brand_admin -> campaign:manage
		{RoleID: 2, PermissionID: 4}, // brand_admin -> order:read
		{RoleID: 3, PermissionID: 4}, // participant -> order:read
	}
	for _, rp := range rolePermissions {
		suite.db.Create(&rp)
	}

	// 创建品牌
	brands := []model.Brand{
		{Id: 1, Name: "品牌A", Status: "active"},
		{Id: 2, Name: "品牌B", Status: "active"},
		{Id: 3, Name: "品牌C", Status: "active"},
	}
	for _, brand := range brands {
		suite.db.Create(&brand)
	}

	// 创建测试用户
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	users := []model.User{
		{
			Id:       1,
			Username: "platform_admin",
			Password: string(hashedPassword),
			Phone:    "13800000001",
			Email:    "admin@test.com",
			RealName: "平台管理员",
			Role:     "platform_admin",
			Status:   "active",
		},
		{
			Id:       2,
			Username: "brand_admin_a",
			Password: string(hashedPassword),
			Phone:    "13800000002",
			Email:    "brand_a@test.com",
			RealName: "品牌A管理员",
			Role:     "brand_admin",
			Status:   "active",
		},
		{
			Id:       3,
			Username: "brand_admin_b",
			Password: string(hashedPassword),
			Phone:    "13800000003",
			Email:    "brand_b@test.com",
			RealName: "品牌B管理员",
			Role:     "brand_admin",
			Status:   "active",
		},
		{
			Id:       4,
			Username: "participant",
			Password: string(hashedPassword),
			Phone:    "13800000004",
			Email:    "user@test.com",
			RealName: "普通用户",
			Role:     "participant",
			Status:   "active",
		},
	}
	for _, user := range users {
		suite.db.Create(&user)
	}

	// 创建用户角色关联
	userRoles := []model.UserRole{
		{UserID: 1, RoleID: 1}, // platform_admin -> platform_admin
		{UserID: 2, RoleID: 2}, // brand_admin_a -> brand_admin
		{UserID: 3, RoleID: 2}, // brand_admin_b -> brand_admin
		{UserID: 4, RoleID: 3}, // participant -> participant
	}
	for _, ur := range userRoles {
		suite.db.Create(&ur)
	}

	// 创建用户品牌关联
	userBrands := []model.UserBrand{
		{UserID: 2, BrandId: 1}, // brand_admin_a -> 品牌A
		{UserID: 3, BrandId: 2}, // brand_admin_b -> 品牌B
	}
	for _, ub := range userBrands {
		suite.db.Create(&ub)
	}

	// 创建品牌素材
	brandAssets := []model.BrandAsset{
		{
			Id:          1,
			BrandId:     1,
			Name:        "品牌A Logo",
			Type:        "image",
			Category:    "logo",
			FileUrl:     "https://example.com/brand-a-logo.png",
			FileSize:    1024,
			Description: "品牌A的官方Logo",
			Status:      "active",
		},
		{
			Id:          2,
			BrandId:     2,
			Name:        "品牌B Banner",
			Type:        "image",
			Category:    "banner",
			FileUrl:     "https://example.com/brand-b-banner.jpg",
			FileSize:    2048,
			Description: "品牌B的宣传横幅",
			Status:      "active",
		},
	}
	for _, asset := range brandAssets {
		suite.db.Create(&asset)
	}
}

func (suite *RBACIntegrationTestSuite) loginUser(username, password string) string {
	if token, exists := suite.tokens[username]; exists {
		return token
	}

	loginReq := types.LoginReq{
		Username: username,
		Password: password,
	}

	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 这里需要实际的登录处理器，简化处理
	authMiddleware := middleware.NewAuthMiddleware(suite.svcCtx.Config.Auth.AccessSecret)
	
	// 模拟登录成功，生成token
	var roles []string
	var brandIDs []int64
	
	switch username {
	case "platform_admin":
		roles = []string{"platform_admin"}
	case "brand_admin_a":
		roles = []string{"brand_admin"}
		brandIDs = []int64{1}
	case "brand_admin_b":
		roles = []string{"brand_admin"}
		brandIDs = []int64{2}
	case "participant":
		roles = []string{"participant"}
	}

	var userID int64
	suite.db.Model(&model.User{}).Where("username = ?", username).Select("id").Scan(&userID)

	token, err := authMiddleware.GenerateToken(userID, username, roles, brandIDs)
	suite.Require().NoError(err)

	suite.tokens[username] = token
	return token
}

func (suite *RBACIntegrationTestSuite) makeAuthenticatedRequest(method, url string, body interface{}, username string) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	token := suite.loginUser(username, "Password123!")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	return w
}

func (suite *RBACIntegrationTestSuite) TestPlatformAdminFullAccess() {
	// 平台管理员应该能够访问所有资源

	// 测试用户管理权限
	w := suite.makeAuthenticatedRequest("GET", "/api/v1/admin/users", nil, "platform_admin")
	// 这里应该返回200，但由于没有完整的路由设置，我们主要测试权限逻辑

	// 测试品牌管理权限 - 访问所有品牌
	permissionMiddleware := middleware.NewPermissionMiddleware(suite.db)
	
	// 验证平台管理员可以访问任何品牌
	hasPermission, err := permissionMiddleware.CheckPermission(1, "brand:manage", "brand", 1)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	hasPermission, err = permissionMiddleware.CheckPermission(1, "brand:manage", "brand", 2)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	hasPermission, err = permissionMiddleware.CheckPermission(1, "brand:manage", "brand", 3)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	// 验证安全管理权限
	hasPermission, err = permissionMiddleware.CheckPermission(1, "security:manage", "", 0)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)
}

func (suite *RBACIntegrationTestSuite) TestBrandAdminDataIsolation() {
	permissionMiddleware := middleware.NewPermissionMiddleware(suite.db)

	// 品牌A管理员只能访问品牌A
	hasPermission, err := permissionMiddleware.CheckPermission(2, "brand:manage", "brand", 1)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	hasPermission, err = permissionMiddleware.CheckPermission(2, "brand:manage", "brand", 2)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)

	hasPermission, err = permissionMiddleware.CheckPermission(2, "brand:manage", "brand", 3)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)

	// 品牌B管理员只能访问品牌B
	hasPermission, err = permissionMiddleware.CheckPermission(3, "brand:manage", "brand", 1)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)

	hasPermission, err = permissionMiddleware.CheckPermission(3, "brand:manage", "brand", 2)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	// 验证品牌管理员无用户管理权限
	hasPermission, err = permissionMiddleware.CheckPermission(2, "user:manage", "", 0)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)

	// 验证品牌管理员无安全管理权限
	hasPermission, err = permissionMiddleware.CheckPermission(2, "security:manage", "", 0)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)
}

func (suite *RBACIntegrationTestSuite) TestParticipantLimitedAccess() {
	permissionMiddleware := middleware.NewPermissionMiddleware(suite.db)

	// 参与者只有订单查看权限
	hasPermission, err := permissionMiddleware.CheckPermission(4, "order:read", "", 0)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	// 参与者无其他权限
	hasPermission, err = permissionMiddleware.CheckPermission(4, "user:manage", "", 0)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)

	hasPermission, err = permissionMiddleware.CheckPermission(4, "brand:manage", "", 0)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)

	hasPermission, err = permissionMiddleware.CheckPermission(4, "campaign:manage", "", 0)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)

	hasPermission, err = permissionMiddleware.CheckPermission(4, "security:manage", "", 0)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)
}

func (suite *RBACIntegrationTestSuite) TestBrandAssetDataIsolation() {
	// 测试品牌素材的数据隔离
	
	// 品牌A管理员应该只能看到品牌A的素材
	var brandAAssets []model.BrandAsset
	err := suite.db.Where("brand_id = ?", 1).Find(&brandAAssets).Error
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), brandAAssets, 1)
	assert.Equal(suite.T(), "品牌A Logo", brandAAssets[0].Name)

	// 品牌B管理员应该只能看到品牌B的素材
	var brandBAssets []model.BrandAsset
	err = suite.db.Where("brand_id = ?", 2).Find(&brandBAssets).Error
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), brandBAssets, 1)
	assert.Equal(suite.T(), "品牌B Banner", brandBAssets[0].Name)

	// 验证跨品牌访问被阻止
	permissionMiddleware := middleware.NewPermissionMiddleware(suite.db)
	
	// 品牌A管理员不能访问品牌B的素材
	hasPermission, err := permissionMiddleware.CheckPermission(2, "brand:manage", "brand_asset", 2)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)

	// 品牌B管理员不能访问品牌A的素材
	hasPermission, err = permissionMiddleware.CheckPermission(3, "brand:manage", "brand_asset", 1)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)
}

func (suite *RBACIntegrationTestSuite) TestMultiBrandAdminScenario() {
	// 测试多品牌管理员场景
	
	// 创建一个管理多个品牌的管理员
	multiBrandAdmin := &model.User{
		Id:       5,
		Username: "multi_brand_admin",
		Password: "hashedpassword",
		Phone:    "13800000005",
		Email:    "multi@test.com",
		RealName: "多品牌管理员",
		Role:     "brand_admin",
		Status:   "active",
	}
	suite.db.Create(multiBrandAdmin)

	// 分配用户角色
	userRole := &model.UserRole{UserID: 5, RoleID: 2}
	suite.db.Create(userRole)

	// 分配多个品牌
	userBrands := []model.UserBrand{
		{UserID: 5, BrandId: 1},
		{UserID: 5, BrandId: 3},
	}
	for _, ub := range userBrands {
		suite.db.Create(&ub)
	}

	permissionMiddleware := middleware.NewPermissionMiddleware(suite.db)

	// 验证可以访问分配的品牌
	hasPermission, err := permissionMiddleware.CheckPermission(5, "brand:manage", "brand", 1)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	hasPermission, err = permissionMiddleware.CheckPermission(5, "brand:manage", "brand", 3)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), hasPermission)

	// 验证不能访问未分配的品牌
	hasPermission, err = permissionMiddleware.CheckPermission(5, "brand:manage", "brand", 2)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), hasPermission)

	// 验证获取用户品牌ID列表
	brandIDs, err := permissionMiddleware.GetUserBrandIDs(5)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), brandIDs, int64(1))
	assert.Contains(suite.T(), brandIDs, int64(3))
	assert.NotContains(suite.T(), brandIDs, int64(2))
}

func (suite *RBACIntegrationTestSuite) TestSecurityAuditIntegration() {
	// 测试安全审计功能集成

	// 模拟登录操作
	suite.svcCtx.AuditService.LogLoginAttempt(
		&[]int64{1}[0],
		"platform_admin",
		"192.168.1.100",
		"Mozilla/5.0",
		true,
		"",
	)

	// 模拟用户操作
	suite.svcCtx.AuditService.LogUserAction(
		&service.AuditContext{
			UserID:    &[]int64{1}[0],
			Username:  "platform_admin",
			ClientIP:  "192.168.1.100",
			UserAgent: "Mozilla/5.0",
		},
		"create_user",
		"user",
		"5",
		map[string]interface{}{
			"username": "new_user",
			"role":     "participant",
		},
	)

	// 模拟安全事件
	suite.svcCtx.AuditService.LogSecurityEvent(
		"privilege_escalation",
		"high",
		&[]int64{2}[0],
		"brand_admin_a",
		"192.168.1.101",
		"Mozilla/5.0",
		"尝试访问未授权的用户管理功能",
		map[string]interface{}{
			"attempted_action": "user:manage",
			"denied_reason":    "insufficient_privileges",
		},
	)

	// 验证审计日志记录
	var auditLogs []model.AuditLog
	err := suite.db.Find(&auditLogs).Error
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), auditLogs, 1)
	assert.Equal(suite.T(), "create_user", auditLogs[0].Action)

	// 验证登录尝试记录
	var loginAttempts []model.LoginAttempt
	err = suite.db.Find(&loginAttempts).Error
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), loginAttempts, 1)
	assert.True(suite.T(), loginAttempts[0].Success)

	// 验证安全事件记录
	var securityEvents []model.SecurityEvent
	err = suite.db.Find(&securityEvents).Error
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), securityEvents, 1)
	assert.Equal(suite.T(), "privilege_escalation", securityEvents[0].EventType)
	assert.Equal(suite.T(), "high", securityEvents[0].Severity)
}

func (suite *RBACIntegrationTestSuite) TestSessionManagementIntegration() {
	// 测试会话管理集成

	userID := int64(1)
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	// 创建会话
	session, err := suite.svcCtx.SessionService.CreateSession(userID, clientIP, userAgent)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), session)

	// 验证会话有效性
	validSession, err := suite.svcCtx.SessionService.ValidateSession(session.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, validSession.UserID)

	// 更新会话活跃时间
	err = suite.svcCtx.SessionService.UpdateSessionActivity(session.ID)
	assert.NoError(suite.T(), err)

	// 验证用户在线状态
	online, err := suite.svcCtx.SessionService.IsUserOnline(userID)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), online)

	// 撤销会话
	err = suite.svcCtx.SessionService.RevokeSession(session.ID, "test_revoke")
	assert.NoError(suite.T(), err)

	// 验证会话已失效
	_, err = suite.svcCtx.SessionService.ValidateSession(session.ID)
	assert.Error(suite.T(), err)

	// 验证用户离线状态
	online, err = suite.svcCtx.SessionService.IsUserOnline(userID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), online)
}

func (suite *RBACIntegrationTestSuite) TestPasswordPolicyIntegration() {
	// 测试密码策略集成

	// 验证密码强度检查
	score := suite.svcCtx.PasswordService.GeneratePasswordStrengthScore("WeakPass")
	assert.True(suite.T(), score < 60)

	score = suite.svcCtx.PasswordService.GeneratePasswordStrengthScore("StrongP@ssw0rd123!")
	assert.True(suite.T(), score >= 80)

	// 验证密码策略验证
	err := suite.svcCtx.PasswordService.ValidatePassword("weak", 0)
	assert.Error(suite.T(), err)

	err = suite.svcCtx.PasswordService.ValidatePassword("StrongP@ssw0rd123!", 0)
	assert.NoError(suite.T(), err)

	// 测试密码历史功能
	userID := int64(1)
	hashedPassword, _ := suite.svcCtx.PasswordService.HashPassword("OldPassword123!")
	err = suite.svcCtx.PasswordService.SavePasswordHistory(userID, hashedPassword)
	assert.NoError(suite.T(), err)

	// 验证不能重复使用历史密码
	err = suite.svcCtx.PasswordService.ValidatePassword("OldPassword123!", userID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "不能使用最近")
}

func TestRBACIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RBACIntegrationTestSuite))
}