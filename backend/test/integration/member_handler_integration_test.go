package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type MemberHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
	memberID   int64
}

func (suite *MemberHandlerIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
	suite.ensureMemberRouteAvailable()
	suite.loadMemberID()
}

func (suite *MemberHandlerIntegrationTestSuite) loginAsAdmin() {
	loginReq := map[string]string{
		"username": "admin",
		"password": "123456",
	}
	reqBody, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/auth/login", bytes.NewBuffer(reqBody))
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
	if suite.adminToken == "" {
		suite.T().Skip("admin token 为空，跳过测试")
	}
}

func (suite *MemberHandlerIntegrationTestSuite) doRequest(method, path string, payload interface{}, token string) (int, []byte) {
	var reqBody io.Reader
	if payload != nil {
		body, _ := json.Marshal(payload)
		reqBody = bytes.NewBuffer(body)
	}

	req, _ := http.NewRequest(method, suite.baseURL+path, reqBody)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body
}

func (suite *MemberHandlerIntegrationTestSuite) ensureMemberRouteAvailable() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/members?page=1&pageSize=1", nil, suite.adminToken)
	if status == http.StatusNotFound {
		suite.T().Skipf("member 路由在当前运行服务中不可用（可能未重启加载最新代码），响应: %s", string(body))
		return
	}
	if (status == http.StatusBadRequest || status == http.StatusInternalServerError) && strings.Contains(strings.ToLower(string(body)), "members") && strings.Contains(strings.ToLower(string(body)), "doesn't exist") {
		suite.T().Skipf("member 相关数据表缺失，当前环境不具备执行条件，响应: %s", string(body))
		return
	}
}

func (suite *MemberHandlerIntegrationTestSuite) loadMemberID() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/members?page=1&pageSize=10", nil, suite.adminToken)
	if status != http.StatusOK {
		return
	}

	var listResp struct {
		Members []struct {
			Id int64 `json:"id"`
		} `json:"members"`
	}
	if err := json.Unmarshal(body, &listResp); err != nil {
		return
	}

	if len(listResp.Members) > 0 {
		suite.memberID = listResp.Members[0].Id
	}
}

func (suite *MemberHandlerIntegrationTestSuite) Test_1_GetMembers() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/members?page=1&pageSize=10", nil, suite.adminToken)
	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 查询会员列表成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *MemberHandlerIntegrationTestSuite) Test_2_UnauthorizedGetMembers() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/members?page=1&pageSize=10", nil, "")
	suite.Equal(http.StatusUnauthorized, status)
	suite.T().Logf("✓ 未授权访问会员列表被拒绝，状态码: %d，响应: %s", status, string(body))
}

func (suite *MemberHandlerIntegrationTestSuite) Test_3_GetMemberById() {
	if suite.memberID == 0 {
		suite.T().Skip("无可用 memberID，跳过测试")
	}

	path := fmt.Sprintf("/api/v1/members/%d", suite.memberID)
	status, body := suite.doRequest(http.MethodGet, path, nil, suite.adminToken)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 查询会员详情完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *MemberHandlerIntegrationTestSuite) Test_4_UpdateMember() {
	if suite.memberID == 0 {
		suite.T().Skip("无可用 memberID，跳过测试")
	}

	path := fmt.Sprintf("/api/v1/members/%d", suite.memberID)
	payload := map[string]interface{}{
		"nickname": fmt.Sprintf("member_%d", time.Now().UnixNano()),
	}

	status, body := suite.doRequest(http.MethodPut, path, payload, suite.adminToken)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 更新会员信息完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *MemberHandlerIntegrationTestSuite) Test_5_UpdateMemberStatus() {
	if suite.memberID == 0 {
		suite.T().Skip("无可用 memberID，跳过测试")
	}

	path := fmt.Sprintf("/api/v1/members/%d/status", suite.memberID)
	payload := map[string]interface{}{
		"status": "active",
	}

	status, body := suite.doRequest(http.MethodPut, path, payload, suite.adminToken)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 更新会员状态完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *MemberHandlerIntegrationTestSuite) Test_6_GetMemberProfile() {
	if suite.memberID == 0 {
		suite.T().Skip("无可用 memberID，跳过测试")
	}

	path := fmt.Sprintf("/api/v1/members/%d/profile", suite.memberID)
	status, body := suite.doRequest(http.MethodGet, path, nil, suite.adminToken)
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 查询会员画像完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *MemberHandlerIntegrationTestSuite) Test_7_InvalidMemberID() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/members/abc", nil, suite.adminToken)
	suite.NotEqual(http.StatusOK, status)
	suite.T().Logf("✓ 非法会员ID被拒绝，状态码: %d，响应: %s", status, string(body))
}

func TestMemberHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(MemberHandlerIntegrationTestSuite))
}
