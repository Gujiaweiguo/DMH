package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type SecurityHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
}

func (suite *SecurityHandlerIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
}

func (suite *SecurityHandlerIntegrationTestSuite) loginAsAdmin() {
	loginReq := map[string]string{
		"username": "admin",
		"password": "123456",
	}
	reqBody, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("无法连接到后端服务: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var loginResp struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(body, &loginResp)

	if loginResp.Token == "" {
		suite.T().Skipf("登录失败: %s", string(body))
		return
	}

	suite.adminToken = loginResp.Token
	suite.T().Log("✓ Admin 登录成功")
}

func (suite *SecurityHandlerIntegrationTestSuite) doRequest(method, path string, payload interface{}, withAuth bool) (int, []byte) {
	var reqBody io.Reader
	if payload != nil {
		body, _ := json.Marshal(payload)
		reqBody = bytes.NewBuffer(body)
	}

	req, _ := http.NewRequest(method, suite.baseURL+path, reqBody)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if withAuth {
		req.Header.Set("Authorization", "Bearer "+suite.adminToken)
	}

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body
}

func (suite *SecurityHandlerIntegrationTestSuite) Test_1_GetPasswordPolicy() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/security/password-policy", nil, true)
	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 获取密码策略成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *SecurityHandlerIntegrationTestSuite) Test_2_UpdatePasswordPolicy() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	updateReq := map[string]interface{}{
		"minLength":             8,
		"requireUppercase":      true,
		"requireLowercase":      true,
		"requireNumbers":        true,
		"requireSpecialChars":   false,
		"maxAge":                90,
		"historyCount":          5,
		"maxLoginAttempts":      5,
		"lockoutDuration":       30,
		"sessionTimeout":        120,
		"maxConcurrentSessions": 3,
	}

	status, body := suite.doRequest(http.MethodPut, "/api/v1/security/password-policy", updateReq, true)
	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 更新密码策略完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *SecurityHandlerIntegrationTestSuite) Test_3_CheckPasswordStrength() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	checkReq := map[string]string{
		"oldPassword": "123456",
		"newPassword": "Password123!",
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/security/check-password-strength", checkReq, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 检查密码强度完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *SecurityHandlerIntegrationTestSuite) Test_4_GetLoginAttempts() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/security/login-attempts?page=1&pageSize=10", nil, true)
	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 查询登录尝试记录成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *SecurityHandlerIntegrationTestSuite) Test_5_GetAuditLogs() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/security/audit-logs?page=1&pageSize=10", nil, true)
	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 查询审计日志成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *SecurityHandlerIntegrationTestSuite) Test_6_GetSecurityEvents() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/security/events?page=1&pageSize=10", nil, true)
	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 查询安全事件成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *SecurityHandlerIntegrationTestSuite) Test_7_GetUserSessions() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/security/sessions?page=1&pageSize=10", nil, true)
	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 查询用户会话成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *SecurityHandlerIntegrationTestSuite) Test_8_RevokeSession() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodDelete, "/api/v1/security/sessions/integration-test-session", nil, true)
	allowedStatus := []int{http.StatusOK, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 撤销会话请求完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *SecurityHandlerIntegrationTestSuite) Test_9_ForceLogoutUser() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/security/force-logout/1", nil, true)
	allowedStatus := []int{http.StatusOK, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 强制登出请求完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *SecurityHandlerIntegrationTestSuite) Test_10_HandleSecurityEvent() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	handleReq := map[string]string{
		"note": fmt.Sprintf("integration test handle at %d", time.Now().Unix()),
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/security/events/1/handle", handleReq, true)
	allowedStatus := []int{http.StatusOK, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 处理安全事件请求完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *SecurityHandlerIntegrationTestSuite) Test_11_UnauthorizedAccess() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/security/password-policy", nil, false)
	suite.Equal(http.StatusUnauthorized, status)
	suite.T().Logf("✓ 未授权访问被拒绝，状态码: %d，响应: %s", status, string(body))
}

func TestSecurityHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(SecurityHandlerIntegrationTestSuite))
}
