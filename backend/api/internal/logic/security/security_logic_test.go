package security

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSecurityLogicTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	err = db.AutoMigrate(&model.PasswordPolicy{}, &model.UserSession{}, &model.User{}, &model.SecurityEvent{}, &model.AuditLog{})
	if err != nil {
		t.Fatalf("failed to migrate sqlite db: %v", err)
	}

	return db
}

func TestSecurityLogicConstructors(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}

	assert.NotNil(t, NewCheckPasswordStrengthLogic(ctx, svcCtx))
	assert.NotNil(t, NewForceLogoutUserLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetAuditLogsLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetLoginAttemptsLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetPasswordPolicyLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetSecurityEventsLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetUserSessionsLogic(ctx, svcCtx))
	assert.NotNil(t, NewHandleSecurityEventLogic(ctx, svcCtx))
	assert.NotNil(t, NewRevokeSessionLogic(ctx, svcCtx))
	assert.NotNil(t, NewUpdatePasswordPolicyLogic(ctx, svcCtx))
}

func TestSecurityLogicCurrentMethodBehavior(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}

	check := NewCheckPasswordStrengthLogic(ctx, svcCtx)
	checkResp, checkErr := check.CheckPasswordStrength(&types.ChangePasswordReq{NewPassword: "Password123!"})
	assert.NoError(t, checkErr)
	assert.NotNil(t, checkResp)
	assert.NotZero(t, checkResp.Score)

	force := NewForceLogoutUserLogic(ctx, svcCtx)
	forceResp, forceErr := force.ForceLogoutUser(1, &types.ForceLogoutReq{})
	assert.Error(t, forceErr)
	assert.Nil(t, forceResp)

	audit := NewGetAuditLogsLogic(ctx, svcCtx)
	auditResp, auditErr := audit.GetAuditLogs()
	assert.NoError(t, auditErr)
	assert.NotNil(t, auditResp)
	assert.Equal(t, int64(0), auditResp.Total)
	assert.Len(t, auditResp.Logs, 0)

	attempts := NewGetLoginAttemptsLogic(ctx, svcCtx)
	attemptsResp, attemptsErr := attempts.GetLoginAttempts()
	assert.NoError(t, attemptsErr)
	assert.NotNil(t, attemptsResp)
	assert.Equal(t, int64(0), attemptsResp.Total)
	assert.Len(t, attemptsResp.Attempts, 0)

	policy := NewGetPasswordPolicyLogic(ctx, svcCtx)
	policyResp, policyErr := policy.GetPasswordPolicy()
	assert.NoError(t, policyErr)
	assert.NotNil(t, policyResp)
	assert.Equal(t, 8, policyResp.MinLength)
	assert.Equal(t, 5, policyResp.MaxLoginAttempts)

	events := NewGetSecurityEventsLogic(ctx, svcCtx)
	eventsResp, eventsErr := events.GetSecurityEvents()
	assert.NoError(t, eventsErr)
	assert.NotNil(t, eventsResp)
	assert.Equal(t, int64(0), eventsResp.Total)
	assert.Len(t, eventsResp.Events, 0)

	sessions := NewGetUserSessionsLogic(ctx, svcCtx)
	sessionsResp, sessionsErr := sessions.GetUserSessions()
	assert.NoError(t, sessionsErr)
	assert.NotNil(t, sessionsResp)
	assert.Equal(t, int64(0), sessionsResp.Total)
	assert.Len(t, sessionsResp.Sessions, 0)

	handle := NewHandleSecurityEventLogic(ctx, svcCtx)
	handleResp, handleErr := handle.HandleSecurityEvent(1, &types.HandleSecurityEventReq{})
	assert.Error(t, handleErr)
	assert.Nil(t, handleResp)

	revoke := NewRevokeSessionLogic(ctx, svcCtx)
	revokeResp, revokeErr := revoke.RevokeSession("session-id")
	assert.Error(t, revokeErr)
	assert.Nil(t, revokeResp)

	update := NewUpdatePasswordPolicyLogic(ctx, svcCtx)
	updateResp, updateErr := update.UpdatePasswordPolicy(&types.UpdatePasswordPolicyReq{MinLength: 8})
	assert.Error(t, updateErr)
	assert.Nil(t, updateResp)
}

func TestUpdatePasswordPolicyLogic_Success(t *testing.T) {
	db := setupSecurityLogicTestDB(t)

	logic := NewUpdatePasswordPolicyLogic(context.WithValue(context.Background(), "userId", int64(1)), &svc.ServiceContext{DB: db})
	resp, err := logic.UpdatePasswordPolicy(&types.UpdatePasswordPolicyReq{
		MinLength:             10,
		RequireUppercase:      true,
		RequireLowercase:      true,
		RequireNumbers:        true,
		RequireSpecialChars:   false,
		MaxAge:                120,
		HistoryCount:          6,
		MaxLoginAttempts:      6,
		LockoutDuration:       20,
		SessionTimeout:        180,
		MaxConcurrentSessions: 4,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 10, resp.MinLength)

	var saved model.PasswordPolicy
	err = db.First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, 10, saved.MinLength)
	assert.Equal(t, 180, saved.SessionTimeout)
}

func TestForceLogoutUserLogic_Success(t *testing.T) {
	db := setupSecurityLogicTestDB(t)

	user := &model.User{Username: "target", Password: "hashed", Phone: "13800138000", Status: "active"}
	assert.NoError(t, db.Create(user).Error)
	assert.NoError(t, db.Create(&model.UserSession{ID: "s1", UserID: user.Id, ClientIP: "127.0.0.1", UserAgent: "ua", Status: "active"}).Error)
	assert.NoError(t, db.Create(&model.UserSession{ID: "s2", UserID: user.Id, ClientIP: "127.0.0.1", UserAgent: "ua", Status: "active"}).Error)

	logic := NewForceLogoutUserLogic(context.WithValue(context.Background(), "userId", int64(99)), &svc.ServiceContext{DB: db})
	resp, err := logic.ForceLogoutUser(user.Id, &types.ForceLogoutReq{Reason: "security check"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "用户已强制下线", resp.Message)

	var activeCount int64
	err = db.Model(&model.UserSession{}).Where("user_id = ? AND status = ?", user.Id, "active").Count(&activeCount).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), activeCount)
}
