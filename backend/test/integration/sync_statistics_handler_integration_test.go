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

type SyncStatisticsHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
}

func (suite *SyncStatisticsHandlerIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
}

func (suite *SyncStatisticsHandlerIntegrationTestSuite) loginAsAdmin() {
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

func (suite *SyncStatisticsHandlerIntegrationTestSuite) doRequest(method, path string, payload interface{}, withAuth bool) (int, []byte) {
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

func (suite *SyncStatisticsHandlerIntegrationTestSuite) Test_1_GetSyncHealth() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/sync/health", nil, true)
	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 同步健康检查成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *SyncStatisticsHandlerIntegrationTestSuite) Test_2_GetSyncStatistics() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/sync/statistics", nil, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 同步统计查询完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *SyncStatisticsHandlerIntegrationTestSuite) Test_3_GetSyncStatus() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/sync/status/1", nil, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 同步状态查询完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *SyncStatisticsHandlerIntegrationTestSuite) Test_4_RetrySync() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/sync/retry/1", nil, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 重试同步请求完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *SyncStatisticsHandlerIntegrationTestSuite) Test_5_GetDashboardStatsWithoutBrandId() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/dashboard-stats", nil, true)
	suite.NotEqual(http.StatusOK, status)
	suite.T().Logf("✓ 无 brandId 查询被拒绝，状态码: %d，响应: %s", status, string(body))
}

func (suite *SyncStatisticsHandlerIntegrationTestSuite) Test_6_GetDashboardStatsWithBrandId() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/dashboard-stats?brandId=1", nil, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 带 brandId 查询完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *SyncStatisticsHandlerIntegrationTestSuite) Test_7_UnauthorizedSyncAccess() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/sync/health", nil, false)
	suite.Equal(http.StatusUnauthorized, status)
	suite.T().Logf("✓ 未授权访问 sync 被拒绝，状态码: %d，响应: %s", status, string(body))
}

func (suite *SyncStatisticsHandlerIntegrationTestSuite) Test_8_UnauthorizedDashboardAccess() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/dashboard-stats", nil, false)
	suite.Equal(http.StatusUnauthorized, status)
	suite.T().Logf("✓ 未授权访问 dashboard 被拒绝，状态码: %d，响应: %s", status, string(body))
}

func TestSyncStatisticsHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(SyncStatisticsHandlerIntegrationTestSuite))
}
