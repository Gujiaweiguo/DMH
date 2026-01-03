package auth

import (
	"context"
	"testing"
	"time"

	"dmh/api/internal/config"
	"dmh/api/internal/service"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type LoginLogicTestSuite struct {
	suite.Suite
	db         *gorm.DB
	svcCtx     *svc.ServiceContext
	loginLogic *LoginLogic
}

func (suite *LoginLogicTestSuite) SetupSuite() {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// 自动迁移表结构
	err = db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.UserBrand{},
		&model.Brand{},
		&model.PasswordPolicy{},
		&model.PasswordHistory{},
		&model.LoginAttempt{},
		&model.UserSession{},
		&model.AuditLog{},
		&model.SecurityEvent{},
	)
	suite.Require().NoError(err)

	// 创建服务上下文
	cfg := config.Config{
		Auth: struct {
			AccessSecret string
			AccessExpire int64
		}{
			AccessSecret: "test-secret-key-for-jwt-token-generation",
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
	suite.loginLogic = NewLoginLogic(context.Background(), suite.svcCtx)
}

func (suite *LoginLogicTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *LoginLogicTestSuite) SetupTest() {
	// 清理测试数据
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM roles")
	suite.db.Exec("DELETE FROM user_roles")
	suite.db.Exec("DELETE FROM user_brands")
	suite.db.Exec("DELETE FROM brands")
	suite.db.Exec("DELETE FROM password_policies")
	suite.db.Exec("DELETE FROM password_histories")
	suite.db.Exec("DELETE FROM login_attempts")
	suite.db.Exec("DELETE FROM user_sessions")
	suite.db.Exec("DELETE FROM audit_logs")
	suite.db.Exec("DELETE FROM security_events")

	// 创建测试数据
	suite.createTestData()
}

func (suite *LoginLogicTestSuite) createTestData() {
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

	// 创建品牌
	brands := []model.Brand{
		{Id: 1, Name: "品牌A", Status: "active"},
		{Id: 2, Name: "品牌B", Status: "active"},
	}
	for _, brand := range brands {
		suite.db.Create(&brand)
	}

	// 创建测试用户
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	users := []model.User{
		{
			Id:       1,
			Username: "admin",
			Password: string(hashedPassword),
			Phone:    "13800000001",
			Email:    "admin@test.com",
			RealName: "管理员",
			Role:     "platform_admin",
			Status:   "active",
		},
		{
			Id:       2,
			Username: "brand_manager",
			Password: string(hashedPassword),
			Phone:    "13800000002",
			Email:    "brand@test.com",
			RealName: "品牌经理",
			Role:     "brand_admin",
			Status:   "active",
		},
		{
			Id:       3,
			Username: "user",
			Password: string(hashedPassword),
			Phone:    "13800000003",
			Email:    "user@test.com",
			RealName: "普通用户",
			Role:     "participant",
			Status:   "active",
		},
		{
			Id:       4,
			Username: "disabled_user",
			Password: string(hashedPassword),
			Phone:    "13800000004",
			Email:    "disabled@test.com",
			RealName: "禁用用户",
			Role:     "participant",
			Status:   "disabled",
		},
	}
	for _, user := range users {
		suite.db.Create(&user)
	}

	// 创建用户角色关联
	userRoles := []model.UserRole{
		{UserID: 1, RoleID: 1},
		{UserID: 2, RoleID: 2},
		{UserID: 3, RoleID: 3},
		{UserID: 4, RoleID: 3},
	}
	for _, ur := range userRoles {
		suite.db.Create(&ur)
	}

	// 创建用户品牌关联
	userBrands := []model.UserBrand{
		{UserID: 2, BrandId: 1},
	}
	for _, ub := range userBrands {
		suite.db.Create(&ub)
	}
}

func (suite *LoginLogicTestSuite) TestSuccessfulLogin() {
	req := &types.LoginReq{
		Username: "admin",
		Password: "Password123!",
	}

	resp, err := suite.loginLogic.Login(req)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(1), resp.UserId)
	assert.Equal(suite.T(), "admin", resp.Username)
	assert.Equal(suite.T(), "13800000001", resp.Phone)
	assert.Equal(suite.T(), "管理员", resp.RealName)
	assert.Contains(suite.T(), resp.Roles, "platform_admin")
	assert.NotEmpty(suite.T(), resp.Token)

	// 验证登录尝试记录
	var attempt model.LoginAttempt
	err = suite.db.Where("username = ? AND success = ?", "admin", true).First(&attempt).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), *attempt.UserID)

	// 验证会话创建
	var session model.UserSession
	err = suite.db.Where("user_id = ? AND status = ?", 1, "active").First(&session).Error
	assert.NoError(suite.T(), err)

	// 验证审计日志
	var auditLog model.AuditLog
	err = suite.db.Where("user_id = ? AND action = ?", 1, "login").First(&auditLog).Error
	assert.NoError(suite.T(), err)
}

func (suite *LoginLogicTestSuite) TestBrandAdminLogin() {
	req := &types.LoginReq{
		Username: "brand_manager",
		Password: "Password123!",
	}

	resp, err := suite.loginLogic.Login(req)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(2), resp.UserId)
	assert.Contains(suite.T(), resp.Roles, "brand_admin")
	assert.Contains(suite.T(), resp.BrandIds, int64(1))
	assert.NotContains(suite.T(), resp.BrandIds, int64(2))
}

func (suite *LoginLogicTestSuite) TestInvalidUsername() {
	req := &types.LoginReq{
		Username: "nonexistent",
		Password: "Password123!",
	}

	resp, err := suite.loginLogic.Login(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "用户名或密码错误")

	// 验证登录失败记录
	var attempt model.LoginAttempt
	err = suite.db.Where("username = ? AND success = ?", "nonexistent", false).First(&attempt).Error
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), attempt.UserID)
}

func (suite *LoginLogicTestSuite) TestInvalidPassword() {
	req := &types.LoginReq{
		Username: "admin",
		Password: "wrongpassword",
	}

	resp, err := suite.loginLogic.Login(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "用户名或密码错误")

	// 验证登录失败记录
	var attempt model.LoginAttempt
	err = suite.db.Where("username = ? AND success = ?", "admin", false).First(&attempt).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), *attempt.UserID)

	// 验证用户登录失败次数增加
	var user model.User
	err = suite.db.First(&user, 1).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, user.LoginAttempts)
}

func (suite *LoginLogicTestSuite) TestDisabledUser() {
	req := &types.LoginReq{
		Username: "disabled_user",
		Password: "Password123!",
	}

	resp, err := suite.loginLogic.Login(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "账号已被禁用")

	// 验证登录失败记录
	var attempt model.LoginAttempt
	err = suite.db.Where("username = ? AND success = ?", "disabled_user", false).First(&attempt).Error
	assert.NoError(suite.T(), err)
}

func (suite *LoginLogicTestSuite) TestAccountLockout() {
	// 模拟多次登录失败
	req := &types.LoginReq{
		Username: "user",
		Password: "wrongpassword",
	}

	// 连续5次失败登录
	for i := 0; i < 5; i++ {
		resp, err := suite.loginLogic.Login(req)
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), resp)
	}

	// 第6次登录应该触发账户锁定
	resp, err := suite.loginLogic.Login(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "账号已被锁定")

	// 验证用户被锁定
	var user model.User
	err = suite.db.First(&user, 3).Error
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user.LockedUntil)
	assert.True(suite.T(), user.LockedUntil.After(time.Now()))

	// 验证安全事件记录
	var event model.SecurityEvent
	err = suite.db.Where("event_type = ? AND user_id = ?", "account_locked", 3).First(&event).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "medium", event.Severity)
}

func (suite *LoginLogicTestSuite) TestValidationErrors() {
	// 测试空用户名
	req := &types.LoginReq{
		Username: "",
		Password: "Password123!",
	}
	resp, err := suite.loginLogic.Login(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "用户名不能为空")

	// 测试空密码
	req = &types.LoginReq{
		Username: "admin",
		Password: "",
	}
	resp, err = suite.loginLogic.Login(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "密码不能为空")

	// 测试用户名长度
	req = &types.LoginReq{
		Username: "ab",
		Password: "Password123!",
	}
	resp, err = suite.loginLogic.Login(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "用户名长度应在3-50个字符之间")

	// 测试密码长度
	req = &types.LoginReq{
		Username: "admin",
		Password: "12345",
	}
	resp, err = suite.loginLogic.Login(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "密码长度不能少于6位")
}

func (suite *LoginLogicTestSuite) TestPasswordExpiration() {
	// 创建过期密码历史
	hashedPassword, _ := suite.svcCtx.PasswordService.HashPassword("Password123!")
	history := &model.PasswordHistory{
		UserID:       1,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now().AddDate(0, 0, -100), // 100天前
	}
	suite.db.Create(history)

	req := &types.LoginReq{
		Username: "admin",
		Password: "Password123!",
	}

	resp, err := suite.loginLogic.Login(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "密码已过期")
}

func (suite *LoginLogicTestSuite) TestConcurrentSessions() {
	req := &types.LoginReq{
		Username: "admin",
		Password: "Password123!",
	}

	// 创建3个会话（达到限制）
	for i := 0; i < 3; i++ {
		resp, err := suite.loginLogic.Login(req)
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), resp)
	}

	// 验证有3个活跃会话
	var count int64
	err := suite.db.Model(&model.UserSession{}).
		Where("user_id = ? AND status = ?", 1, "active").
		Count(&count).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), count)

	// 第4次登录应该强制下线最旧的会话
	resp, err := suite.loginLogic.Login(req)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)

	// 验证仍然只有3个活跃会话
	err = suite.db.Model(&model.UserSession{}).
		Where("user_id = ? AND status = ?", 1, "active").
		Count(&count).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), count)
}

func TestLoginLogicTestSuite(t *testing.T) {
	suite.Run(t, new(LoginLogicTestSuite))
}