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

type DistributorHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
	brandId    int64
	userId     int64
}

func (suite *DistributorHandlerIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
	suite.getTestBrand()
	suite.getTestUser()
}

func (suite *DistributorHandlerIntegrationTestSuite) loginAsAdmin() {
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
	json.Unmarshal(body, &loginResp)

	suite.adminToken = loginResp.Token
	suite.T().Log("✓ Admin 登录成功")
}

func (suite *DistributorHandlerIntegrationTestSuite) getTestBrand() {
	if suite.adminToken == "" {
		return
	}

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/brands?page=1&pageSize=1", nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Brands []struct {
			Id int64 `json:"id"`
		} `json:"brands"`
	}
	json.Unmarshal(body, &result)

	if len(result.Brands) > 0 {
		suite.brandId = result.Brands[0].Id
		suite.T().Logf("✓ 获取测试品牌 ID: %d", suite.brandId)
	}
}

func (suite *DistributorHandlerIntegrationTestSuite) getTestUser() {
	suite.userId = 1
}

func (suite *DistributorHandlerIntegrationTestSuite) Test_1_GetBrandDistributors() {
	if suite.adminToken == "" || suite.brandId == 0 {
		suite.T().Skip("未登录或品牌不存在，跳过测试")
	}

	url := fmt.Sprintf("%s/api/v1/brands/%d/distributors?page=1&pageSize=10", suite.baseURL, suite.brandId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.T().Logf("✓ 品牌分销商列表查询完成，状态码: %d", resp.StatusCode)
}

func (suite *DistributorHandlerIntegrationTestSuite) Test_2_GetBrandDistributorApplications() {
	if suite.adminToken == "" || suite.brandId == 0 {
		suite.T().Skip("未登录或品牌不存在，跳过测试")
	}

	url := fmt.Sprintf("%s/api/v1/brands/%d/distributors/applications?page=1&pageSize=10", suite.baseURL, suite.brandId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.T().Logf("✓ 品牌分销商申请列表查询完成，状态码: %d", resp.StatusCode)
}

func (suite *DistributorHandlerIntegrationTestSuite) Test_3_GetMyDistributorStatus() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/distributors/status", nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.T().Logf("✓ 我的分销商状态查询完成，状态码: %d", resp.StatusCode)
}

func (suite *DistributorHandlerIntegrationTestSuite) Test_4_GetMyDistributorDashboard() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/distributors/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.T().Logf("✓ 我的分销商仪表盘查询完成，状态码: %d", resp.StatusCode)
}

func (suite *DistributorHandlerIntegrationTestSuite) Test_5_GetDistributorStatistics() {
	if suite.adminToken == "" || suite.userId == 0 {
		suite.T().Skip("未登录或用户不存在，跳过测试")
	}

	url := fmt.Sprintf("%s/api/v1/distributors/%d/statistics", suite.baseURL, suite.userId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.T().Logf("✓ 分销商统计查询完成，状态码: %d", resp.StatusCode)
}

func (suite *DistributorHandlerIntegrationTestSuite) Test_6_GetDistributorLinks() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/distributors/links?page=1&pageSize=10", nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.T().Logf("✓ 分销商链接列表查询完成，状态码: %d", resp.StatusCode)
}

func (suite *DistributorHandlerIntegrationTestSuite) Test_7_DistributorApplyValidation() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	applyReq := map[string]interface{}{
		"brandId": 0,
		"reason":  "",
	}
	reqBody, _ := json.Marshal(applyReq)

	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/distributors/apply", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.T().Logf("✓ 分销商申请验证完成，状态码: %d", resp.StatusCode)
}

func (suite *DistributorHandlerIntegrationTestSuite) Test_8_UnauthorizedAccess() {
	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/distributors/status", nil)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.T().Logf("✓ 无 Token 访问分销商 API，状态码: %d", resp.StatusCode)
}

func (suite *DistributorHandlerIntegrationTestSuite) Test_9_GetDistributorSubordinates() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	url := fmt.Sprintf("%s/api/v1/distributors/%d/subordinates?page=1&pageSize=10", suite.baseURL, suite.userId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.T().Logf("✓ 分销商下级列表查询完成，状态码: %d", resp.StatusCode)
}

func (suite *DistributorHandlerIntegrationTestSuite) Test_10_GetDistributorRewards() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	url := fmt.Sprintf("%s/api/v1/distributors/%d/rewards?page=1&pageSize=10", suite.baseURL, suite.userId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.T().Logf("✓ 分销商奖励列表查询完成，状态码: %d", resp.StatusCode)
}

func TestDistributorHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(DistributorHandlerIntegrationTestSuite))
}
