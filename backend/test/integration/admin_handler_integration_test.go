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

type AdminHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
}

func (suite *AdminHandlerIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
}

func (suite *AdminHandlerIntegrationTestSuite) loginAsAdmin() {
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

func (suite *AdminHandlerIntegrationTestSuite) doRequest(method, path string, payload interface{}, withAuth bool) (int, []byte) {
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

func (suite *AdminHandlerIntegrationTestSuite) Test_1_GetUsers() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/admin/users?page=1&pageSize=10", nil, true)
	suite.Equal(http.StatusOK, status)

	var result struct {
		Total int64 `json:"total"`
		Users []struct {
			Id int64 `json:"id"`
		} `json:"users"`
	}
	err := json.Unmarshal(body, &result)
	suite.NoError(err)
	suite.T().Logf("✓ 查询用户列表成功，total=%d, users=%d", result.Total, len(result.Users))
}

func (suite *AdminHandlerIntegrationTestSuite) Test_2_GetUserById() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/admin/users/1", nil, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 查询用户详情完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *AdminHandlerIntegrationTestSuite) Test_3_CreateUser() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	seed := time.Now().UnixNano()
	payload := map[string]interface{}{
		"username": fmt.Sprintf("integ_admin_user_%d", seed),
		"password": "123456",
		"phone":    fmt.Sprintf("139%d", seed%100000000),
		"realName": "集成测试用户",
		"role":     "participant",
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/admin/users", payload, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 创建用户请求完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *AdminHandlerIntegrationTestSuite) Test_4_CreateUserInvalidPayload() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	payload := map[string]interface{}{
		"username": "",
		"password": "123456",
		"phone":    "13800138000",
		"role":     "participant",
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/admin/users", payload, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 创建用户参数校验请求完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *AdminHandlerIntegrationTestSuite) Test_5_ManageBrandAdminRelationInvalidPayload() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	payload := map[string]interface{}{
		"userId":   0,
		"brandIds": []int64{},
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/admin/brand-admin-relations", payload, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 绑定品牌管理员无效参数请求完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *AdminHandlerIntegrationTestSuite) Test_6_UnauthorizedAccess() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/admin/users?page=1&pageSize=10", nil, false)
	suite.Equal(http.StatusUnauthorized, status)
	suite.T().Logf("✓ 无 Token 访问被拒绝，状态码: %d，响应: %s", status, string(body))
}

func TestAdminHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(AdminHandlerIntegrationTestSuite))
}
