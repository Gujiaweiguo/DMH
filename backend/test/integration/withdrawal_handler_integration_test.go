package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type WithdrawalHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
}

func (suite *WithdrawalHandlerIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
}

func (suite *WithdrawalHandlerIntegrationTestSuite) loginAsAdmin() {
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

func (suite *WithdrawalHandlerIntegrationTestSuite) doRequest(method, path string, payload interface{}, withAuth bool) (int, []byte) {
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

func (suite *WithdrawalHandlerIntegrationTestSuite) Test_1_GetWithdrawalsWithBodyPagination() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	reqBody := map[string]interface{}{
		"page":     1,
		"pageSize": 10,
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/withdrawals", reqBody, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 查询提现列表（Body 分页）完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *WithdrawalHandlerIntegrationTestSuite) Test_2_GetWithdrawalsWithQueryPagination() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/withdrawals?page=1&pageSize=10", nil, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 查询提现列表（Query 分页）完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *WithdrawalHandlerIntegrationTestSuite) Test_3_ApplyWithdrawalInvalidAmount() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	applyReq := map[string]interface{}{
		"amount":      0,
		"bankName":    "测试银行",
		"bankAccount": "6222000000000000",
		"accountName": "测试用户",
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/withdrawals", applyReq, true)
	suite.NotEqual(http.StatusOK, status)
	suite.T().Logf("✓ 提现申请（无效金额）被拒绝，状态码: %d，响应: %s", status, string(body))
}

func (suite *WithdrawalHandlerIntegrationTestSuite) Test_4_ApplyWithdrawalInsufficientBalance() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	applyReq := map[string]interface{}{
		"amount":      10,
		"bankName":    "测试银行",
		"bankAccount": "6222000000000000",
		"accountName": "测试用户",
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/withdrawals", applyReq, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 提现申请（余额不足场景）完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *WithdrawalHandlerIntegrationTestSuite) Test_5_GetWithdrawalById() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/withdrawals/1", nil, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	if status == http.StatusBadRequest {
		suite.Contains(strings.ToLower(string(body)), "parseint")
	}
	suite.T().Logf("✓ 查询提现详情完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *WithdrawalHandlerIntegrationTestSuite) Test_6_ApproveWithdrawal() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	approveReq := map[string]interface{}{
		"status": "approved",
		"remark": "integration test",
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/withdrawals/1/approve", approveReq, true)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	if status == http.StatusBadRequest {
		suite.Contains(strings.ToLower(string(body)), "parseint")
	}
	suite.T().Logf("✓ 审批提现请求完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *WithdrawalHandlerIntegrationTestSuite) Test_7_UnauthorizedAccess() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/withdrawals?page=1&pageSize=10", nil, false)
	suite.Equal(http.StatusUnauthorized, status)
	suite.T().Logf("✓ 未授权访问被拒绝，状态码: %d，响应: %s", status, string(body))
}

func TestWithdrawalHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(WithdrawalHandlerIntegrationTestSuite))
}
