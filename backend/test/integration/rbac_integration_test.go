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

type RBACIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
}

func (suite *RBACIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
}

func (suite *RBACIntegrationTestSuite) loginAsAdmin() {
	loginReq := map[string]string{"username": "admin", "password": "123456"}
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("无法连接后端，跳过 RBAC 集成测试: %v", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var loginResp struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(respBody, &loginResp)
	if loginResp.Token == "" {
		suite.T().Skipf("登录失败，跳过 RBAC 集成测试: %s", string(respBody))
		return
	}
	suite.adminToken = loginResp.Token
}

func (suite *RBACIntegrationTestSuite) TestUnauthorizedAdminAPI() {
	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/admin/users", nil)
	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	suite.True(resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden)
}

func (suite *RBACIntegrationTestSuite) TestAdminTokenCanAccessAdminAPI() {
	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)
	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
}

func TestRBACIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RBACIntegrationTestSuite))
}
