package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type RoleHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
}

func (suite *RoleHandlerIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
}

func (suite *RoleHandlerIntegrationTestSuite) loginAsAdmin() {
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

func (suite *RoleHandlerIntegrationTestSuite) doRequest(method, path string, payload interface{}, withAuth bool) (int, []byte) {
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

func (suite *RoleHandlerIntegrationTestSuite) Test_1_GetPermissions() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/permissions", nil, true)
	suite.Equal(http.StatusOK, status)

	var result []map[string]interface{}
	err := json.Unmarshal(body, &result)
	suite.NoError(err)
	suite.T().Logf("✓ 查询权限列表成功，条数: %d", len(result))
}

func (suite *RoleHandlerIntegrationTestSuite) Test_2_GetRoles() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/roles", nil, true)
	suite.Equal(http.StatusOK, status)

	var result []map[string]interface{}
	err := json.Unmarshal(body, &result)
	suite.NoError(err)
	suite.T().Logf("✓ 查询角色列表成功，条数: %d", len(result))
}

func (suite *RoleHandlerIntegrationTestSuite) Test_3_GetUserPermissions() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/users/1/permissions", nil, true)
	suite.Equal(http.StatusOK, status)

	var result map[string]interface{}
	err := json.Unmarshal(body, &result)
	suite.NoError(err)
	suite.T().Logf("✓ 查询用户权限成功，响应: %s", string(body))
}

func (suite *RoleHandlerIntegrationTestSuite) Test_4_GetUserPermissionsNonExistentUser() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/users/999999/permissions", nil, true)
	suite.Equal(http.StatusOK, status)

	var result struct {
		UserId int64 `json:"userId"`
	}
	err := json.Unmarshal(body, &result)
	suite.NoError(err)
	suite.Equal(int64(999999), result.UserId)
	suite.T().Logf("✓ 不存在用户权限查询完成，响应: %s", string(body))
}

func (suite *RoleHandlerIntegrationTestSuite) Test_5_GetUserPermissionsInvalidUserId() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/users/abc/permissions", nil, true)
	suite.NotEqual(http.StatusOK, status)
	suite.T().Logf("✓ 非法用户ID被拒绝，状态码: %d，响应: %s", status, string(body))
}

func (suite *RoleHandlerIntegrationTestSuite) Test_6_ConfigRolePermissions() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	payload := map[string]interface{}{
		"roleId":        1,
		"permissionIds": []int64{1, 2},
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/roles/permissions", payload, true)
	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 配置角色权限成功，响应: %s", string(body))
}

func (suite *RoleHandlerIntegrationTestSuite) Test_7_ConfigRolePermissionsInvalidJSON() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/roles/permissions", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	suite.NotEqual(http.StatusOK, resp.StatusCode)
	suite.T().Logf("✓ 非法 JSON 被拒绝，状态码: %d，响应: %s", resp.StatusCode, string(body))
}

func (suite *RoleHandlerIntegrationTestSuite) Test_8_UnauthorizedAccess() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/roles", nil, false)
	suite.Equal(http.StatusUnauthorized, status)
	suite.T().Logf("✓ 无 Token 访问被拒绝，状态码: %d，响应: %s", status, string(body))
}

func TestRoleHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RoleHandlerIntegrationTestSuite))
}
