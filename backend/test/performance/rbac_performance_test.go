package performance

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"dmh/api/internal/middleware"
	"dmh/api/internal/service"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RBACPerformanceTestSuite struct {
	suite.Suite
	db                   *gorm.DB
	permissionMiddleware *middleware.PermissionMiddleware
	passwordService      *service.PasswordService
	auditService         *service.AuditService
	sessionService       *service.SessionService
}

func (suite *RBACPerformanceTestSuite) SetupSuite() {
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
		&model.PasswordPolicy{},
		&model.PasswordHistory{},
		&model.LoginAttempt{},
		&model.UserSession{},
		&model.AuditLog{},
		&model.SecurityEvent{},
	)
	suite.Require().NoError(err)

	suite.db = db
	suite.permissionMiddleware = middleware.NewPermissionMiddleware(db)
	suite.passwordService = service.NewPasswordService(db)
	suite.auditService = service.NewAuditService(db)
	suite.sessionService = service.NewSessionService(db, suite.passwordService)

	// 创建大量测试数据
	suite.createLargeTestDataset()
}

func (suite *RBACPerformanceTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *RBACPerformanceTestSuite) createLargeTestDataset() {
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
		{RoleID: 1, PermissionID: 1},
		{RoleID: 1, PermissionID: 2},
		{RoleID: 1, PermissionID: 3},
		{RoleID: 1, PermissionID: 4},
		{RoleID: 1, PermissionID: 5},
		{RoleID: 2, PermissionID: 2},
		{RoleID: 2, PermissionID: 3},
		{RoleID: 2, PermissionID: 4},
		{RoleID: 3, PermissionID: 4},
	}
	for _, rp := range rolePermissions {
		suite.db.Create(&rp)
	}

	// 创建1000个品牌
	fmt.Println("Creating 1000 brands...")
	brands := make([]model.Brand, 1000)
	for i := 0; i < 1000; i++ {
		brands[i] = model.Brand{
			Id:     int64(i + 1),
			Name:   fmt.Sprintf("品牌%d", i+1),
			Status: "active",
		}
	}
	suite.db.CreateInBatches(brands, 100)

	// 创建10000个用户
	fmt.Println("Creating 10000 users...")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	users := make([]model.User, 10000)
	for i := 0; i < 10000; i++ {
		roleType := "participant"
		if i < 100 {
			roleType = "platform_admin"
		} else if i < 1100 {
			roleType = "brand_admin"
		}

		users[i] = model.User{
			Id:       int64(i + 1),
			Username: fmt.Sprintf("user%d", i+1),
			Password: string(hashedPassword),
			Phone:    fmt.Sprintf("138%08d", i+1),
			Email:    fmt.Sprintf("user%d@test.com", i+1),
			RealName: fmt.Sprintf("用户%d", i+1),
			Role:     roleType,
			Status:   "active",
		}
	}
	suite.db.CreateInBatches(users, 100)

	// 创建用户角色关联
	fmt.Println("Creating user role associations...")
	userRoles := make([]model.UserRole, 10000)
	for i := 0; i < 10000; i++ {
		roleID := int64(3) // participant
		if i < 100 {
			roleID = 1 // platform_admin
		} else if i < 1100 {
			roleID = 2 // brand_admin
		}

		userRoles[i] = model.UserRole{
			UserID: int64(i + 1),
			RoleID: roleID,
		}
	}
	suite.db.CreateInBatches(userRoles, 100)

	// 为品牌管理员分配品牌（每个品牌管理员管理1-5个品牌）
	fmt.Println("Creating user brand associations...")
	var userBrands []model.UserBrand
	for i := 100; i < 1100; i++ { // 品牌管理员
		userID := int64(i + 1)
		brandCount := rand.Intn(5) + 1 // 1-5个品牌
		
		for j := 0; j < brandCount; j++ {
			brandID := int64(rand.Intn(1000) + 1)
			userBrands = append(userBrands, model.UserBrand{
				UserID:  userID,
				BrandId: brandID,
			})
		}
	}
	suite.db.CreateInBatches(userBrands, 100)

	fmt.Println("Test data creation completed!")
}

func (suite *RBACPerformanceTestSuite) TestPermissionCheckPerformance() {
	// 测试权限检查性能
	
	userIDs := []int64{1, 101, 1001, 5001, 9001} // 不同类型的用户
	permissions := []string{"user:manage", "brand:manage", "campaign:manage", "order:read", "security:manage"}
	
	iterations := 1000
	
	start := time.Now()
	
	for i := 0; i < iterations; i++ {
		userID := userIDs[i%len(userIDs)]
		permission := permissions[i%len(permissions)]
		
		_, err := suite.permissionMiddleware.CheckPermission(userID, permission, "", 0)
		assert.NoError(suite.T(), err)
	}
	
	duration := time.Since(start)
	avgTime := duration / time.Duration(iterations)
	
	fmt.Printf("Permission check performance:\n")
	fmt.Printf("Total time: %v\n", duration)
	fmt.Printf("Average time per check: %v\n", avgTime)
	fmt.Printf("Checks per second: %.2f\n", float64(iterations)/duration.Seconds())
	
	// 性能要求：每次权限检查应该在1ms以内
	assert.True(suite.T(), avgTime < time.Millisecond, "Permission check should be faster than 1ms")
}

func (suite *RBACPerformanceTestSuite) TestBrandPermissionCheckPerformance() {
	// 测试品牌权限检查性能
	
	brandAdminIDs := []int64{101, 201, 301, 401, 501} // 品牌管理员
	brandIDs := []int64{1, 100, 200, 500, 999}
	
	iterations := 1000
	
	start := time.Now()
	
	for i := 0; i < iterations; i++ {
		userID := brandAdminIDs[i%len(brandAdminIDs)]
		brandID := brandIDs[i%len(brandIDs)]
		
		_, err := suite.permissionMiddleware.CheckPermission(userID, "brand:manage", "brand", brandID)
		assert.NoError(suite.T(), err)
	}
	
	duration := time.Since(start)
	avgTime := duration / time.Duration(iterations)
	
	fmt.Printf("Brand permission check performance:\n")
	fmt.Printf("Total time: %v\n", duration)
	fmt.Printf("Average time per check: %v\n", avgTime)
	fmt.Printf("Checks per second: %.2f\n", float64(iterations)/duration.Seconds())
	
	// 性能要求：每次品牌权限检查应该在2ms以内
	assert.True(suite.T(), avgTime < 2*time.Millisecond, "Brand permission check should be faster than 2ms")
}

func (suite *RBACPerformanceTestSuite) TestConcurrentPermissionChecks() {
	// 测试并发权限检查性能
	
	concurrency := 100
	checksPerGoroutine := 100
	totalChecks := concurrency * checksPerGoroutine
	
	var wg sync.WaitGroup
	start := time.Now()
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			
			for j := 0; j < checksPerGoroutine; j++ {
				userID := int64(rand.Intn(10000) + 1)
				permissions := []string{"user:manage", "brand:manage", "campaign:manage", "order:read"}
				permission := permissions[rand.Intn(len(permissions))]
				
				_, err := suite.permissionMiddleware.CheckPermission(userID, permission, "", 0)
				assert.NoError(suite.T(), err)
			}
		}(i)
	}
	
	wg.Wait()
	duration := time.Since(start)
	avgTime := duration / time.Duration(totalChecks)
	
	fmt.Printf("Concurrent permission check performance:\n")
	fmt.Printf("Concurrency: %d goroutines\n", concurrency)
	fmt.Printf("Total checks: %d\n", totalChecks)
	fmt.Printf("Total time: %v\n", duration)
	fmt.Printf("Average time per check: %v\n", avgTime)
	fmt.Printf("Checks per second: %.2f\n", float64(totalChecks)/duration.Seconds())
	
	// 性能要求：并发情况下每次权限检查应该在5ms以内
	assert.True(suite.T(), avgTime < 5*time.Millisecond, "Concurrent permission check should be faster than 5ms")
}

func (suite *RBACPerformanceTestSuite) TestPermissionCacheEffectiveness() {
	// 测试权限缓存效果
	
	userID := int64(101) // 品牌管理员
	permission := "brand:manage"
	
	// 第一次查询（无缓存）
	start := time.Now()
	_, err := suite.permissionMiddleware.GetUserPermissions(userID)
	assert.NoError(suite.T(), err)
	firstQueryTime := time.Since(start)
	
	// 第二次查询（有缓存）
	start = time.Now()
	_, err = suite.permissionMiddleware.GetUserPermissions(userID)
	assert.NoError(suite.T(), err)
	cachedQueryTime := time.Since(start)
	
	fmt.Printf("Permission cache effectiveness:\n")
	fmt.Printf("First query time (no cache): %v\n", firstQueryTime)
	fmt.Printf("Cached query time: %v\n", cachedQueryTime)
	fmt.Printf("Cache speedup: %.2fx\n", float64(firstQueryTime)/float64(cachedQueryTime))
	
	// 缓存应该至少提供2倍的性能提升
	assert.True(suite.T(), cachedQueryTime*2 < firstQueryTime, "Cache should provide at least 2x speedup")
	
	// 测试多次缓存查询的一致性能
	iterations := 100
	start = time.Now()
	
	for i := 0; i < iterations; i++ {
		_, err := suite.permissionMiddleware.GetUserPermissions(userID)
		assert.NoError(suite.T(), err)
	}
	
	avgCachedTime := time.Since(start) / time.Duration(iterations)
	fmt.Printf("Average cached query time (%d iterations): %v\n", iterations, avgCachedTime)
	
	// 缓存查询应该非常快
	assert.True(suite.T(), avgCachedTime < 100*time.Microsecond, "Cached queries should be very fast")
}

func (suite *RBACPerformanceTestSuite) TestPasswordServicePerformance() {
	// 测试密码服务性能
	
	passwords := []string{
		"SimplePass123!",
		"ComplexP@ssw0rd!",
		"VeryLongAndComplexPasswordWithManyCharacters123!@#",
		"Short1!",
		"MediumLength123!",
	}
	
	iterations := 1000
	
	// 测试密码强度检查性能
	start := time.Now()
	for i := 0; i < iterations; i++ {
		password := passwords[i%len(passwords)]
		suite.passwordService.GeneratePasswordStrengthScore(password)
	}
	strengthCheckTime := time.Since(start)
	
	// 测试密码验证性能
	start = time.Now()
	for i := 0; i < iterations; i++ {
		password := passwords[i%len(passwords)]
		suite.passwordService.ValidatePassword(password, 0)
	}
	validationTime := time.Since(start)
	
	// 测试密码加密性能
	start = time.Now()
	for i := 0; i < iterations; i++ {
		password := passwords[i%len(passwords)]
		suite.passwordService.HashPassword(password)
	}
	hashingTime := time.Since(start)
	
	fmt.Printf("Password service performance:\n")
	fmt.Printf("Strength check time (%d iterations): %v (avg: %v)\n", 
		iterations, strengthCheckTime, strengthCheckTime/time.Duration(iterations))
	fmt.Printf("Validation time (%d iterations): %v (avg: %v)\n", 
		iterations, validationTime, validationTime/time.Duration(iterations))
	fmt.Printf("Hashing time (%d iterations): %v (avg: %v)\n", 
		iterations, hashingTime, hashingTime/time.Duration(iterations))
	
	// 性能要求
	assert.True(suite.T(), strengthCheckTime/time.Duration(iterations) < time.Millisecond, 
		"Password strength check should be fast")
	assert.True(suite.T(), validationTime/time.Duration(iterations) < time.Millisecond, 
		"Password validation should be fast")
	// 密码加密可以稍慢，因为需要安全性
	assert.True(suite.T(), hashingTime/time.Duration(iterations) < 100*time.Millisecond, 
		"Password hashing should complete within reasonable time")
}

func (suite *RBACPerformanceTestSuite) TestSessionServicePerformance() {
	// 测试会话服务性能
	
	userIDs := []int64{1, 101, 1001, 5001, 9001}
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"
	
	iterations := 1000
	var sessionIDs []string
	
	// 测试会话创建性能
	start := time.Now()
	for i := 0; i < iterations; i++ {
		userID := userIDs[i%len(userIDs)]
		session, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
		assert.NoError(suite.T(), err)
		sessionIDs = append(sessionIDs, session.ID)
	}
	creationTime := time.Since(start)
	
	// 测试会话验证性能
	start = time.Now()
	for i := 0; i < iterations; i++ {
		sessionID := sessionIDs[i]
		_, err := suite.sessionService.ValidateSession(sessionID)
		assert.NoError(suite.T(), err)
	}
	validationTime := time.Since(start)
	
	// 测试会话更新性能
	start = time.Now()
	for i := 0; i < iterations; i++ {
		sessionID := sessionIDs[i]
		err := suite.sessionService.UpdateSessionActivity(sessionID)
		assert.NoError(suite.T(), err)
	}
	updateTime := time.Since(start)
	
	fmt.Printf("Session service performance:\n")
	fmt.Printf("Creation time (%d sessions): %v (avg: %v)\n", 
		iterations, creationTime, creationTime/time.Duration(iterations))
	fmt.Printf("Validation time (%d sessions): %v (avg: %v)\n", 
		iterations, validationTime, validationTime/time.Duration(iterations))
	fmt.Printf("Update time (%d sessions): %v (avg: %v)\n", 
		iterations, updateTime, updateTime/time.Duration(iterations))
	
	// 性能要求
	assert.True(suite.T(), creationTime/time.Duration(iterations) < 5*time.Millisecond, 
		"Session creation should be fast")
	assert.True(suite.T(), validationTime/time.Duration(iterations) < time.Millisecond, 
		"Session validation should be very fast")
	assert.True(suite.T(), updateTime/time.Duration(iterations) < 2*time.Millisecond, 
		"Session update should be fast")
}

func (suite *RBACPerformanceTestSuite) TestAuditServicePerformance() {
	// 测试审计服务性能
	
	userIDs := []int64{1, 101, 1001, 5001, 9001}
	actions := []string{"create_user", "update_user", "delete_user", "login", "logout"}
	resources := []string{"user", "brand", "campaign", "order"}
	
	iterations := 1000
	
	// 测试审计日志记录性能
	start := time.Now()
	for i := 0; i < iterations; i++ {
		userID := userIDs[i%len(userIDs)]
		action := actions[i%len(actions)]
		resource := resources[i%len(resources)]
		
		ctx := &service.AuditContext{
			UserID:    &userID,
			Username:  fmt.Sprintf("user%d", userID),
			ClientIP:  "192.168.1.100",
			UserAgent: "Mozilla/5.0",
		}
		
		err := suite.auditService.LogUserAction(ctx, action, resource, fmt.Sprintf("%d", i), nil)
		assert.NoError(suite.T(), err)
	}
	loggingTime := time.Since(start)
	
	// 测试审计日志查询性能
	start = time.Now()
	for i := 0; i < 100; i++ { // 较少的查询次数
		_, _, err := suite.auditService.GetAuditLogs(1, 20, map[string]interface{}{})
		assert.NoError(suite.T(), err)
	}
	queryTime := time.Since(start)
	
	fmt.Printf("Audit service performance:\n")
	fmt.Printf("Logging time (%d logs): %v (avg: %v)\n", 
		iterations, loggingTime, loggingTime/time.Duration(iterations))
	fmt.Printf("Query time (100 queries): %v (avg: %v)\n", 
		queryTime, queryTime/100)
	
	// 性能要求
	assert.True(suite.T(), loggingTime/time.Duration(iterations) < 2*time.Millisecond, 
		"Audit logging should be fast")
	assert.True(suite.T(), queryTime/100 < 10*time.Millisecond, 
		"Audit log queries should be reasonably fast")
}

func (suite *RBACPerformanceTestSuite) TestLargeDatasetQueries() {
	// 测试大数据集查询性能
	
	// 测试获取用户权限（在大量用户中）
	start := time.Now()
	for i := 0; i < 100; i++ {
		userID := int64(rand.Intn(10000) + 1)
		_, err := suite.permissionMiddleware.GetUserPermissions(userID)
		assert.NoError(suite.T(), err)
	}
	userPermissionTime := time.Since(start)
	
	// 测试获取用户品牌（在大量品牌中）
	start = time.Now()
	for i := 0; i < 100; i++ {
		userID := int64(rand.Intn(1000) + 101) // 品牌管理员范围
		_, err := suite.permissionMiddleware.GetUserBrandIDs(userID)
		assert.NoError(suite.T(), err)
	}
	userBrandTime := time.Since(start)
	
	// 测试品牌权限检查（在大量品牌中）
	start = time.Now()
	for i := 0; i < 100; i++ {
		userID := int64(rand.Intn(1000) + 101) // 品牌管理员范围
		brandID := int64(rand.Intn(1000) + 1)
		_, err := suite.permissionMiddleware.CheckPermission(userID, "brand:manage", "brand", brandID)
		assert.NoError(suite.T(), err)
	}
	brandPermissionTime := time.Since(start)
	
	fmt.Printf("Large dataset query performance:\n")
	fmt.Printf("User permission queries (100 queries): %v (avg: %v)\n", 
		userPermissionTime, userPermissionTime/100)
	fmt.Printf("User brand queries (100 queries): %v (avg: %v)\n", 
		userBrandTime, userBrandTime/100)
	fmt.Printf("Brand permission checks (100 checks): %v (avg: %v)\n", 
		brandPermissionTime, brandPermissionTime/100)
	
	// 性能要求：即使在大数据集中，查询也应该保持合理的性能
	assert.True(suite.T(), userPermissionTime/100 < 10*time.Millisecond, 
		"User permission queries should scale well")
	assert.True(suite.T(), userBrandTime/100 < 10*time.Millisecond, 
		"User brand queries should scale well")
	assert.True(suite.T(), brandPermissionTime/100 < 10*time.Millisecond, 
		"Brand permission checks should scale well")
}

func (suite *RBACPerformanceTestSuite) TestMemoryUsage() {
	// 测试内存使用情况
	
	// 这是一个简化的内存使用测试
	// 在实际环境中，可以使用更专业的内存分析工具
	
	iterations := 10000
	
	// 测试权限检查的内存使用
	for i := 0; i < iterations; i++ {
		userID := int64(rand.Intn(10000) + 1)
		_, err := suite.permissionMiddleware.GetUserPermissions(userID)
		assert.NoError(suite.T(), err)
		
		// 每1000次迭代清理一次缓存，模拟实际使用场景
		if i%1000 == 0 {
			suite.permissionMiddleware.ClearUserPermissionCache(userID)
		}
	}
	
	fmt.Printf("Memory usage test completed with %d iterations\n", iterations)
	// 在实际测试中，这里应该检查内存使用情况
}

func TestRBACPerformanceTestSuite(t *testing.T) {
	suite.Run(t, new(RBACPerformanceTestSuite))
}

// 基准测试
func BenchmarkPermissionCheck(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{}, &model.Role{}, &model.Permission{}, &model.UserRole{}, &model.RolePermission{})
	
	// 创建基本测试数据
	role := &model.Role{ID: 1, Code: "test_role"}
	db.Create(role)
	
	permission := &model.Permission{ID: 1, Code: "test:permission"}
	db.Create(permission)
	
	user := &model.User{Id: 1, Username: "testuser", Role: "test_role"}
	db.Create(user)
	
	db.Create(&model.UserRole{UserID: 1, RoleID: 1})
	db.Create(&model.RolePermission{RoleID: 1, PermissionID: 1})
	
	middleware := middleware.NewPermissionMiddleware(db)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		middleware.CheckPermission(1, "test:permission", "", 0)
	}
}

func BenchmarkPasswordStrengthCheck(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.PasswordPolicy{})
	
	service := service.NewPasswordService(db)
	password := "TestPassword123!"
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		service.GeneratePasswordStrengthScore(password)
	}
}

func BenchmarkSessionValidation(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.UserSession{}, &model.PasswordPolicy{})
	
	passwordService := service.NewPasswordService(db)
	sessionService := service.NewSessionService(db, passwordService)
	
	// 创建测试会话
	session, _ := sessionService.CreateSession(1, "192.168.1.1", "test")
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		sessionService.ValidateSession(session.ID)
	}
}